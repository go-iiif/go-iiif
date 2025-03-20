package seed

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/aaronland/gocloud-blob/bucket"
	aws_lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/fsnotify/fsnotify"
	iiifaws "github.com/go-iiif/go-iiif/v6/aws"
	iiiftile "github.com/go-iiif/go-iiif/v6/tile"
	"github.com/sfomuseum/go-csvdict/v2"
)

func Run(ctx context.Context) error {

	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	opts, err := RunOptionsFromFlagSet(ctx, fs)

	if err != nil {
		return err
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	if opts.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	ts, err := iiiftile.NewTileSeed(ctx, opts.Config, 256, 256, opts.Endpoint, opts.Quality, opts.Format)

	if err != nil {
		return fmt.Errorf("Failed to create tileseed(er), %w", err)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	throttle := make(chan bool, opts.Workers)

	for i := 0; i < opts.Workers; i++ {
		throttle <- true
	}

	tile_func := func(ctx context.Context, tiled_im *TiledImage, wg *sync.WaitGroup) error {

		select {
		case <-ctx.Done():
			return nil
		default:
			// pass
		}

		wg.Add(1)
		t1 := time.Now()

		<-throttle

		slog.Debug("Tile waiting to seed", "source", tiled_im.Source, "time", time.Since(t1))

		go func(ctx context.Context, tiled_im *TiledImage, wg *sync.WaitGroup) {

			t1 := time.Now()

			src_id := tiled_im.Source
			alt_id := tiled_im.Target

			logger := slog.Default()
			logger = logger.With("source", src_id)
			logger = logger.With("target", alt_id)

			defer func() {
				logger.Debug("Time to seed tiles", "time", time.Since(t1))
				throttle <- true
				wg.Done()
			}()

			select {
			case <-ctx.Done():
				return
			default:
				// pass
			}

			count, err := ts.SeedTiles(ctx, src_id, alt_id, opts.ScaleFactors, opts.Refresh)

			// FIX ME

			/*
				if t.onCompleteFunc != nil {
					t.onCompleteFunc(config, src_id, alt_id, count, err)
				}
			*/

			if err != nil {
				logger.Warn("Failed to seed tiles", "error", err)
			} else {
				logger.Debug("Seeded tiles complete", "count", count)
			}

		}(ctx, tiled_im, wg)

		return nil
	}

	switch opts.Mode {
	case "cli", "-":

		wg := new(sync.WaitGroup)

		for _, id := range opts.Paths {

			tiled_im, err := TiledImageFromString(id, opts.NoExtension)

			if err != nil {
				return fmt.Errorf("Failed to derive seed from URI '%s', %w", id, err)
			}

			tile_func(ctx, tiled_im, wg)
		}

		wg.Wait()

	case "csv":

		csv_bucket, err := bucket.OpenBucket(ctx, csv_source)

		if err != nil {
			return fmt.Errorf("Failed to open bucket from CSV source, %w", err)
		}

		defer csv_bucket.Close()

		wg := new(sync.WaitGroup)

		for _, path := range opts.Paths {

			logger := slog.Default()
			logger = logger.With("path", path)

			r, err := csv_bucket.NewReader(ctx, path, nil)

			if err != nil {
				return fmt.Errorf("Failed to open reader from %s, %w", path, err)
			}

			defer r.Close()

			csv_r, err := csvdict.NewReader(r)

			if err != nil {
				return fmt.Errorf("Failed to open CSV reader, %w", err)
			}

			for row, err := range csv_r.Iterate() {

				if err != nil {
					return err
				}

				src_id, ok := row["source_id"]

				if !ok {
					logger.Warn("Unable to determine source ID for row", "row", row)
					continue
				}

				alt_id, ok := row["alternate_id"]

				if !ok {
					logger.Warn("Unable to determine alternate ID for row", "row", row)
					continue
				}

				tiled_im := &TiledImage{
					Source: src_id,
					Target: alt_id,
				}

				tile_func(ctx, tiled_im, wg)
			}

		}

		wg.Wait()

	case "fsnotify":

		images_source := opts.Config.Images.Source.Path

		u, err := url.Parse(images_source)

		if err != nil {
			return fmt.Errorf("Failed to parse images source, %w", err)
		}

		if u.Scheme != "file" {
			return fmt.Errorf("Invalid image source for -mode fsnotify")
		}

		root := u.Path

		watcher, err := fsnotify.NewWatcher()

		if err != nil {
			return fmt.Errorf("Failed to create fsnotify watcher, %w", err)
		}

		defer watcher.Close()

		logger := slog.Default()
		logger = logger.With("root", root)

		logger.Info("Watching filesystem")

		done := make(chan bool)
		wg := new(sync.WaitGroup)

		go func() {

			for {
				select {
				case event, ok := <-watcher.Events:

					if !ok {
						return
					}

					if event.Op == fsnotify.Create {

						abs_path := event.Name

						rel_path := strings.Replace(abs_path, root, "", 1)
						rel_path = strings.TrimLeft(rel_path, "/")

						tiled_im, err := TiledImageFromString(rel_path, opts.NoExtension)

						if err != nil {
							logger.Warn("Failed to determine seed from path", "rel_path", rel_path, "abs_path", abs_path, "error", err)
							continue
						}

						err = tile_func(ctx, tiled_im, wg)

						if err != nil {
							logger.Error("Failed to generate tiles from path", "rel_path", rel_path, "error", err)
							continue
						}
					}

				case err, ok := <-watcher.Errors:

					if !ok {
						return
					}

					logger.Warn("fsnotify error", "error", err)
				}
			}
		}()

		err = watcher.Add(root)

		if err != nil {
			return fmt.Errorf("Failed to add '%s' to fsnotify watcher, %w", root, err)
		}

		<-done

		wg.Wait()

	case "lambda":

		handler := func(ctx context.Context, ev iiifaws.Event) error {
			wg := new(sync.WaitGroup)

			for _, r := range ev.Records {

				s3_entity := r.S3
				s3_obj := s3_entity.Object
				s3_key := s3_obj.Key

				s3_fname := filepath.Base(s3_key)

				tiled_im, err := TiledImageFromString(s3_fname, opts.NoExtension)

				if err != nil {
					return fmt.Errorf("Failed to seed tiles from %s, %w", s3_fname, err)
				}

				tile_func(ctx, tiled_im, wg)
			}

			wg.Wait()
			return nil
		}

		aws_lambda.Start(handler)

	default:
		return fmt.Errorf("Invalid -mode")
	}

	return nil
}
