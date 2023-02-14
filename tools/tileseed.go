package tools

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	aws_lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/fsnotify/fsnotify"
	iiifuri "github.com/go-iiif/go-iiif-uri"
	iiifconfig "github.com/go-iiif/go-iiif/v5/config"
	iiiftile "github.com/go-iiif/go-iiif/v5/tile"
	"github.com/sfomuseum/go-csvdict"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/lookup"
	"github.com/whosonfirst/go-whosonfirst-log"
	"gocloud.dev/blob"
)

type Seed struct {
	Source string
	Target string
}

type TileSeedOnCompleteFunc func(*iiifconfig.Config, string, string, int, error)

type TileSeedToolOptions struct {
	URIFunc        URIFunc
	OnCompleteFunc TileSeedOnCompleteFunc
}

type TileSeedTool struct {
	Tool
	uriFunc        URIFunc
	onCompleteFunc TileSeedOnCompleteFunc
}

func SeedFromString(str_uri string, no_extension bool) (*Seed, error) {

	ctx := context.Background()
	u, err := iiifuri.NewURI(ctx, str_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive seed from URI %s, %v", str_uri, err)
	}

	return SeedFromURI(u, no_extension)
}

func SeedFromURI(u iiifuri.URI, no_extension bool) (*Seed, error) {

	origin := u.Origin()
	target, err := u.Target(nil)

	if err != nil {
		return nil, err
	}

	if no_extension {
		target = strings.TrimSuffix(target, filepath.Ext(target))
	}

	seed := &Seed{
		Source: origin,
		Target: target,
	}

	return seed, nil
}

func NewTileSeedTool() (Tool, error) {

	uri_func := DefaultURIFunc()
	return NewTileSeedToolWithURIFunc(uri_func)
}

// this maintains parity with tools/process.go

func NewTileSeedToolWithURIFunc(uri_func URIFunc) (Tool, error) {

	opts := &TileSeedToolOptions{
		URIFunc: uri_func,
	}

	return NewTileSeedToolWithOptions(opts)
}

// this is the kitchen sink

func NewTileSeedToolWithOptions(opts *TileSeedToolOptions) (Tool, error) {

	t := &TileSeedTool{
		uriFunc:        opts.URIFunc,
		onCompleteFunc: opts.OnCompleteFunc,
	}

	return t, nil
}

func TileSeedToolFlagSet(ctx context.Context) (*flag.FlagSet, error) {

	fs := flag.NewFlagSet("tileseed", flag.ExitOnError)

	err := AppendCommonTileSeedToolFlags(ctx, fs)

	if err != nil {
		return nil, err
	}

	err = AppendTileSeedToolFlags(ctx, fs)

	if err != nil {
		return nil, err
	}

	return fs, nil
}

func AppendCommonTileSeedToolFlags(ctx context.Context, fs *flag.FlagSet) error {

	err := AppendCommonConfigFlags(ctx, fs)

	if err != nil {
		return err
	}

	err = AppendCommonToolModeFlags(ctx, fs)

	if err != nil {
		return err
	}

	return nil
}

func AppendTileSeedToolFlags(ctx context.Context, fs *flag.FlagSet) error {

	fs.String("csv-source", "A valid Go Cloud bucket URI where your CSV tileseed files are located.", "")

	fs.String("scale-factors", "4", "A comma-separated list of scale factors to seed tiles with")
	fs.String("quality", "default", "A valid IIIF quality parameter - if \"default\" then the code will try to determine which format you've set as the default")
	fs.String("format", "jpg", "A valid IIIF format parameter")

	fs.String("logfile", "", "Write logging information to this file")
	fs.String("loglevel", "info", "The amount of logging information to include, valid options are: debug, info, status, warning, error, fatal")

	fs.Int("processes", runtime.NumCPU(), "The number of concurrent processes to use when tiling images")

	fs.Bool("noextension", false, "Remove any extension from destination folder name.")

	fs.Bool("refresh", false, "Refresh a tile even if already exists (default false)")
	fs.String("endpoint", "http://localhost:8080", "The endpoint (scheme, host and optionally port) that will serving these tiles, used for generating an 'info.json' for each source image")

	fs.Bool("verbose", false, "Write logging to STDOUT in addition to any other log targets that may have been defined")

	return nil
}

func (t *TileSeedTool) Run(ctx context.Context) error {

	fs, err := TileSeedToolFlagSet(ctx)

	if err != nil {
		return fmt.Errorf("Failed to create tileseed tool flagset, %w", err)
	}

	flagset.Parse(fs)

	err = flagset.SetFlagsFromEnvVars(fs, "IIIF")

	if err != nil {
		return fmt.Errorf("Failed to assign tileseed tool flags from environment variables, %w", err)
	}

	return t.RunWithFlagSet(ctx, fs)
}

func (t *TileSeedTool) RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	paths := fs.Args()
	return t.RunWithFlagSetAndPaths(ctx, fs, paths...)
}

func (t *TileSeedTool) RunWithFlagSetAndPaths(ctx context.Context, fs *flag.FlagSet, paths ...string) error {

	config_source, err := lookup.StringVar(fs, "config-source")

	if err != nil {
		return fmt.Errorf("Failed to determine config-source flag, %w", err)
	}

	config_name, err := lookup.StringVar(fs, "config-name")

	if err != nil {
		return fmt.Errorf("Failed to determine config-name flag, %w", err)
	}

	csv_source, err := lookup.StringVar(fs, "csv-source")

	if err != nil {
		return fmt.Errorf("Failed to determine csv-source flag, %w", err)
	}

	scale_factors, err := lookup.StringVar(fs, "scale-factors")

	if err != nil {
		return fmt.Errorf("Failed to determine scale-factors flag, %w", err)
	}

	quality, err := lookup.StringVar(fs, "quality")

	if err != nil {
		return fmt.Errorf("Failed to determine quality flag, %w", err)
	}

	format, err := lookup.StringVar(fs, "format")

	if err != nil {
		return fmt.Errorf("Failed to determine format flag, %w", err)
	}

	logfile, err := lookup.StringVar(fs, "logfile")

	if err != nil {
		return fmt.Errorf("Failed to determine logfile flag, %w", err)
	}

	loglevel, err := lookup.StringVar(fs, "loglevel")

	if err != nil {
		return fmt.Errorf("Failed to determine loglevel flag, %w", err)
	}

	processes, err := lookup.IntVar(fs, "processes")

	if err != nil {
		return fmt.Errorf("Failed to determine processes flag, %w", err)
	}

	mode, err := lookup.StringVar(fs, "mode")

	if err != nil {
		return fmt.Errorf("Failed to determine mode flag, %w", err)
	}

	noextension, err := lookup.BoolVar(fs, "noextension")

	if err != nil {
		return fmt.Errorf("Failed to determine noextension flag, %w", err)
	}

	refresh, err := lookup.BoolVar(fs, "refresh")

	if err != nil {
		return err
	}

	endpoint, err := lookup.StringVar(fs, "endpoint")

	if err != nil {
		return err
	}

	verbose, err := lookup.BoolVar(fs, "verbose")

	if err != nil {
		return err
	}

	if config_source == "" {
		return fmt.Errorf("Required -config-source flag is empty.")
	}

	config_bucket, err := blob.OpenBucket(ctx, config_source)

	if err != nil {
		return fmt.Errorf("Failed to open bucket for config source, %w", err)
	}

	config, err := iiifconfig.NewConfigFromBucket(ctx, config_bucket, config_name)

	if err != nil {
		return fmt.Errorf("Failed to create new config from bucket, %w", err)
	}

	ts, err := iiiftile.NewTileSeed(config, 256, 256, endpoint, quality, format)

	if err != nil {
		return fmt.Errorf("Failed to create tileseed(er), %w", err)
	}

	writers := make([]io.Writer, 0)

	if verbose {
		writers = append(writers, os.Stdout)
	}

	if logfile != "" {

		fh, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)

		if err != nil {
			return fmt.Errorf("Failed to open logfile, %w", err)
		}

		writers = append(writers, fh)
	}

	writer := io.MultiWriter(writers...)

	logger := log.NewWOFLogger("")
	logger.AddLogger(writer, loglevel)

	scales := make([]int, 0)

	for _, s := range strings.Split(scale_factors, ",") {

		s = strings.Trim(s, " ")
		scale, err := strconv.Atoi(s)

		if err != nil {
			return fmt.Errorf("Failed to parse scale factor, %w", err)
		}

		scales = append(scales, scale)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	throttle := make(chan bool, processes)

	for i := 0; i < processes; i++ {
		throttle <- true
	}

	tile_func := func(seed *Seed, wg *sync.WaitGroup) error {

		wg.Add(1)
		t1 := time.Now()

		<-throttle

		logger.Debug("Tile waiting to seed '%s': %v", seed.Source, time.Since(t1))

		go func(seed *Seed, wg *sync.WaitGroup) {

			t1 := time.Now()

			src_id := seed.Source
			alt_id := seed.Target

			defer func() {
				logger.Debug("Time to seed tiles for '%s': %v", seed.Source, time.Since(t1))
				throttle <- true
				wg.Done()
			}()

			count, err := ts.SeedTiles(src_id, alt_id, scales, refresh)

			if t.onCompleteFunc != nil {
				t.onCompleteFunc(config, src_id, alt_id, count, err)
			}

			if err != nil {
				logger.Warning("Failed to seed tiles for '%s', %s", src_id, err)
			} else {
				logger.Debug("Seeded %d tiles for '%s'", count, src_id)
			}

		}(seed, wg)

		return nil
	}

	switch mode {
	case "cli", "-":

		wg := new(sync.WaitGroup)

		for _, id := range paths {

			u, err := t.uriFunc(id)

			if err != nil {
				return fmt.Errorf("Failed to derive URI from path '%s', %w", id, err)
			}

			seed, err := SeedFromURI(u, noextension)

			if err != nil {
				return fmt.Errorf("Failed to derive seed from URI '%s', %w", u, err)
			}

			tile_func(seed, wg)
		}

		wg.Wait()

	case "csv":

		csv_bucket, err := blob.OpenBucket(ctx, csv_source)

		if err != nil {
			return fmt.Errorf("Failed to open bucket from CSV source, %w", err)
		}

		wg := new(sync.WaitGroup)

		for _, path := range paths {

			fh, err := csv_bucket.NewReader(ctx, path, nil)

			if err != nil {
				return fmt.Errorf("Failed to open reader from %s, %w", path, err)
			}

			defer fh.Close()

			reader, err := csvdict.NewReader(fh)

			if err != nil {
				return fmt.Errorf("Failed to open CSV reader, %w", err)
			}

			counter := 0

			for {

				row, err := reader.Read()
				counter += 1

				if err == io.EOF {
					break
				}

				if err != nil {
					return err
				}

				src_id, ok := row["source_id"]

				if !ok {
					logger.Warning("Unable to determine source ID", row)
					continue
				}

				alt_id, ok := row["alternate_id"]

				if !ok {
					logger.Warning("Unable to determine alternate ID", row)
					continue
				}

				seed := &Seed{
					Source: src_id,
					Target: alt_id,
				}

				tile_func(seed, wg)
			}

		}

		wg.Wait()

	case "fsnotify":

		images_source := config.Images.Source.Path

		u, err := url.Parse(images_source)

		if err != nil {
			return fmt.Errorf("Failed to parse images source, %w", err)
		}

		if u.Scheme != "file" {
			return fmt.Errorf("Invalid image source for -mode fsnotify")
		}

		root := u.Path
		logger.Info("Watching %s", root)

		watcher, err := fsnotify.NewWatcher()

		if err != nil {
			return fmt.Errorf("Failed to create fsnotify watcher, %w", err)
		}

		defer watcher.Close()

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

						u, err := t.uriFunc(rel_path)

						if err != nil {
							logger.Warning("Failed to run URI function from path '%s' (%s), %s", rel_path, abs_path, err)
							continue
						}

						seed, err := SeedFromURI(u, noextension)

						if err != nil {
							logger.Warning("Failed to determine seed from path '%s' (%s), %s", rel_path, abs_path, err)
							continue
						}

						err = tile_func(seed, wg)

						if err != nil {
							logger.Warning("Failed to generate tiles for path '%s', %s", rel_path, err)
							continue
						}
					}

				case err, ok := <-watcher.Errors:

					if !ok {
						return
					}

					logger.Warning("fsnotify error: %s", err)
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

		handler := func(ctx context.Context, ev Event) error {

			wg := new(sync.WaitGroup)

			for _, r := range ev.Records {

				s3_entity := r.S3
				s3_obj := s3_entity.Object
				s3_key := s3_obj.Key

				s3_fname := filepath.Base(s3_key)

				u, err := t.uriFunc(s3_fname)

				if err != nil {
					return err
				}

				seed, err := SeedFromURI(u, noextension)

				if err != nil {
					return fmt.Errorf("Failed to seed tiles from %s, %w", u, err)
				}

				tile_func(seed, wg)
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
