package seed

/*

go run cmd/iiif-tile-seed/main.go \
	-config-images-source-uri file:///usr/local/src/go-iiif/static/example/images \
	-config-derivatives-cache-uri file:///usr/local/src/go-iiif/work \
	-scale-factors '8,4,2,1' \
	-verbose \
	'rewrite:///spanking-cat.jpg?target=spank'

*/

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"log/slog"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/aaronland/gocloud-blob/bucket"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/fsnotify/fsnotify"
	iiifhtml "github.com/go-iiif/go-iiif/v8/app/seed/html"
	iiifaws "github.com/go-iiif/go-iiif/v8/aws"
	iiifcache "github.com/go-iiif/go-iiif/v8/cache"
	iiifconfig "github.com/go-iiif/go-iiif/v8/config"
	iiiftile "github.com/go-iiif/go-iiif/v8/tile"
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

	if opts.GenerateHTML {

		// Note the use of a locally-scoped *iiifconfig.Config (which may have
		// had its defaults changed if running in CSV mode)

		html_cb := func(cfg *iiifconfig.Config, src_id string, alt_id string, count int, err error) error {

			logger := slog.Default()
			logger = logger.With("source", src_id)
			logger = logger.With("alt", alt_id)

			if err != nil {
				logger.Warn("Skipping on complete func because error present", "error", err)
				return nil
			}

			logger.Info("Generate HTML index page for tiles")

			derivatives_cache, err := iiifcache.NewCache(ctx, cfg.Derivatives.Cache.URI)

			if err != nil {
				logger.Error("Failed to load derivatives cache from config", "error", err)
				return err
			}

			assets := []string{
				"leaflet.iiif.bundle.js",
				"leaflet.css",
				"index.html",
			}

			write_assets := func(assets_fs embed.FS, assets []string) error {

				for _, fname := range assets {

					body, err := assets_fs.ReadFile(fname)

					if err != nil {
						return fmt.Errorf("Failed to read %s, %w", fname, err)
					}

					path := filepath.Join(alt_id, fname)

					err = derivatives_cache.Set(path, body)

					if err != nil {
						return fmt.Errorf("Failed to write %s, %w", path, err)
					}
				}

				return nil
			}

			err = write_assets(iiifhtml.FS, assets)

			if err != nil {
				logger.Error("Failed to write assets", "error", err)
				return err
			}

			return nil
		}

		opts.OnCompleteFuncs = append(opts.OnCompleteFuncs, html_cb)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	throttle := make(chan bool, opts.Workers)

	for i := 0; i < opts.Workers; i++ {
		throttle <- true
	}

	// This is scoped locally so we can access the global opts struct but
	// cfg is passed in because we might be doing path/URI wrangling on a per-image basis (and this run in Go routines)

	tile_func := func(ctx context.Context, cfg *iiifconfig.Config, tiled_im *TiledImage, wg *sync.WaitGroup) error {

		select {
		case <-ctx.Done():
			return nil
		default:
			// pass
		}

		wg.Add(1)
		t1 := time.Now()

		<-throttle

		ts, err := iiiftile.NewTileSeed(ctx, cfg, 256, 256, opts.Endpoint, opts.Quality, opts.Format)

		if err != nil {
			return fmt.Errorf("Failed to create tileseed(er), %w", err)
		}

		slog.Debug("Tile waiting to seed", "source", tiled_im.Source, "time", time.Since(t1))

		go func(ctx context.Context, ts *iiiftile.TileSeed, tiled_im *TiledImage, wg *sync.WaitGroup) {

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

			if err != nil {
				logger.Warn("Failed to seed tiles", "error", err)
			} else {
				logger.Debug("Tile seeding complete", "count", count)
			}

			for _, cb := range opts.OnCompleteFuncs {

				cb_err := cb(cfg, src_id, alt_id, count, err)

				if cb_err != nil {
					logger.Warn("Callback function failed", "error", cb_err)
				}
			}

		}(ctx, ts, tiled_im, wg)

		return nil
	}

	switch opts.Mode {
	case "cli", "-":

		wg := new(sync.WaitGroup)

		for _, id := range opts.Paths {

			tiled_im, err := opts.TiledImageFromString(ctx, id)

			if err != nil {
				return fmt.Errorf("Failed to derive seed from URI '%s', %w", id, err)
			}

			tile_func(ctx, opts.Config, tiled_im, wg)
		}

		wg.Wait()

	case "csv":

		wg := new(sync.WaitGroup)

		for _, path := range opts.Paths {

			bucket_uri, csv_path, err := bucket.ParseURI(path)

			if err != nil {
				return fmt.Errorf("Failed to parse CSV source '%s', %w", path, err)
			}

			csv_bucket, err := bucket.OpenBucket(ctx, bucket_uri)

			if err != nil {
				return fmt.Errorf("Failed to open bucket from CSV source, %w", err)
			}

			defer csv_bucket.Close()

			logger := slog.Default()
			logger = logger.With("path", path)

			r, err := csv_bucket.NewReader(ctx, csv_path, nil)

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

				ts_cfg, err := opts.Config.Clone()

				if err != nil {
					return err
				}

				var source_id string
				var target_id string

				src, ok := row["source_filename"]

				if !ok {
					logger.Warn("Unable to determine source filename for row", "row", row)
					continue
				}

				source_id = src

				src_root, ok := row["source_root"]

				if ok && src_root != "" {

					u, err := url.Parse(src_root)

					if err != nil {
						logger.Warn("Failed to parse source root, skipping", "root", src_root, "error", err)
						continue
					}

					if u.Scheme == "" {
						u.Scheme = "file"
					}

					logger.Info("Assign new source URI", "uri", u.String())
					ts_cfg.Images.Source.URI = u.String()
				}

				target, _ := row["target_filename"]

				if target == "" {
					target_id = source_id
				} else {
					target_id = target
				}

				target_root, ok := row["target_root"]

				if ok && target_root != "" {

					u, err := url.Parse(target_root)

					if err != nil {
						logger.Warn("Failed to parse target root, skipping", "root", target_root, "error", err)
						continue
					}

					if u.Scheme == "" {
						u.Scheme = "file"
					}

					logger.Info("Assign new cache URI", "uri", u.String())
					ts_cfg.Derivatives.Cache.URI = u.String()
				}

				tiled_im := &TiledImage{
					Source: source_id,
					Target: target_id,
				}

				logger.Info("Seed tiles", "source", source_id, "target", target_id)
				err = tile_func(ctx, ts_cfg, tiled_im, wg)

				if err != nil {
					logger.Error("Failed to tile image", "error", err)
				}
			}

		}

		wg.Wait()

	case "fsnotify":

		images_source := opts.Config.Images.Source.URI

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

						tiled_im, err := opts.TiledImageFromString(ctx, rel_path)

						if err != nil {
							logger.Warn("Failed to determine seed from path", "rel_path", rel_path, "abs_path", abs_path, "error", err)
							continue
						}

						err = tile_func(ctx, opts.Config, tiled_im, wg)

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

				tiled_im, err := opts.TiledImageFromString(ctx, s3_fname)

				if err != nil {
					return fmt.Errorf("Failed to seed tiles from %s, %w", s3_fname, err)
				}

				tile_func(ctx, opts.Config, tiled_im, wg)
			}

			wg.Wait()
			return nil
		}

		lambda.Start(handler)

	default:
		return fmt.Errorf("Invalid -mode")
	}

	return nil
}
