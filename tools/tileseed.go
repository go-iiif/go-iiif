package tools

import (
	"context"
	"errors"
	"flag"
	aws_events "github.com/aws/aws-lambda-go/events"
	aws_lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/fsnotify/fsnotify"
	iiifuri "github.com/go-iiif/go-iiif-uri"
	iiifconfig "github.com/go-iiif/go-iiif/v4/config"
	iiiftile "github.com/go-iiif/go-iiif/v4/tile"
	"github.com/sfomuseum/go-flags"
	"github.com/whosonfirst/go-whosonfirst-csv"
	"github.com/whosonfirst/go-whosonfirst-log"
	"gocloud.dev/blob"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
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

	u, err := iiifuri.NewURI(str_uri)

	if err != nil {
		return nil, err
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
		return err
	}

	flags.Parse(fs)

	err = flags.SetFlagsFromEnvVars(fs, "IIIF_TILESEED")

	if err != nil {
		return err
	}

	return t.RunWithFlagSet(ctx, fs)
}

func (t *TileSeedTool) RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	paths := fs.Args()
	return t.RunWithFlagSetAndPaths(ctx, fs, paths...)
}

func (t *TileSeedTool) RunWithFlagSetAndPaths(ctx context.Context, fs *flag.FlagSet, paths ...string) error {

	config_source, err := flags.StringVar(fs, "config-source")

	if err != nil {
		return err
	}

	config_name, err := flags.StringVar(fs, "config-name")

	if err != nil {
		return err
	}

	csv_source, err := flags.StringVar(fs, "csv-source")

	if err != nil {
		return err
	}

	scale_factors, err := flags.StringVar(fs, "scale-factors")

	if err != nil {
		return err
	}

	quality, err := flags.StringVar(fs, "quality")

	if err != nil {
		return err
	}

	format, err := flags.StringVar(fs, "format")

	if err != nil {
		return err
	}

	logfile, err := flags.StringVar(fs, "logfile")

	if err != nil {
		return err
	}

	loglevel, err := flags.StringVar(fs, "loglevel")

	if err != nil {
		return err
	}

	processes, err := flags.IntVar(fs, "processes")

	if err != nil {
		return err
	}

	mode, err := flags.StringVar(fs, "mode")

	if err != nil {
		return err
	}

	noextension, err := flags.BoolVar(fs, "noextension")

	if err != nil {
		return err
	}

	refresh, err := flags.BoolVar(fs, "refresh")

	if err != nil {
		return err
	}

	endpoint, err := flags.StringVar(fs, "endpoint")

	if err != nil {
		return err
	}

	verbose, err := flags.BoolVar(fs, "verbose")

	if err != nil {
		return err
	}

	if config_source == "" {
		return errors.New("Required -config-source flag is empty.")
	}

	config_bucket, err := blob.OpenBucket(ctx, config_source)

	if err != nil {
		return err
	}

	config, err := iiifconfig.NewConfigFromBucket(ctx, config_bucket, config_name)

	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	ts, err := iiiftile.NewTileSeed(config, 256, 256, endpoint, quality, format)

	if err != nil {
		return err
	}

	writers := make([]io.Writer, 0)

	if verbose {
		writers = append(writers, os.Stdout)
	}

	if logfile != "" {

		fh, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)

		if err != nil {
			return err
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
			return err
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
				return err
			}

			seed, err := SeedFromURI(u, noextension)

			if err != nil {
				return err
			}

			tile_func(seed, wg)
		}

		wg.Wait()

	case "csv":

		csv_bucket, err := blob.OpenBucket(ctx, csv_source)

		if err != nil {
			return err
		}

		wg := new(sync.WaitGroup)

		for _, path := range paths {

			fh, err := csv_bucket.NewReader(ctx, path, nil)

			if err != nil {
				return err
			}

			defer fh.Close()

			reader, err := csv.NewDictReader(fh)

			if err != nil {
				return err
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
			return err
		}

		if u.Scheme != "file" {
			return errors.New("Invalid image source for -mode fsnotify")
		}

		root := u.Path
		logger.Info("Watching %s", root)

		watcher, err := fsnotify.NewWatcher()

		if err != nil {
			return err
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
			return err
		}

		<-done

		wg.Wait()

	case "lambda":

		handler := func(ctx context.Context, ev aws_events.S3Event) error {

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
					return err
				}

				tile_func(seed, wg)
			}

			wg.Wait()
			return nil
		}

		aws_lambda.Start(handler)

	default:
		return errors.New("Invalid -mode")
	}

	return nil
}
