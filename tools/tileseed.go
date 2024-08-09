package tools

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aaronland/gocloud-blob/bucket"
	aws_lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/fsnotify/fsnotify"
	iiifuri "github.com/go-iiif/go-iiif-uri"
	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
	iiifcache "github.com/go-iiif/go-iiif/v6/cache"	
	iiiftile "github.com/go-iiif/go-iiif/v6/tile"
	"github.com/go-iiif/go-iiif/v6/static/javascript"
	"github.com/go-iiif/go-iiif/v6/static/css"
	"github.com/go-iiif/go-iiif/v6/static/html"
	"github.com/sfomuseum/go-csvdict"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/lookup"
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

	slog.Debug("New tile seed", "origin", origin, "target", target)

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

	err := AppendCommonFlags(ctx, fs)

	if err != nil {
		return err
	}
	
	err = AppendCommonConfigFlags(ctx, fs)

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

	fs.Bool("generate-tiles-html", false, "If true then the tiles directory will be updated to include HTML/JavaScript/CSS assets to display tiles as a \"slippy\" map (using the leaflet-iiif.js library.")
	
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

	// START OF logging

	logfile, err := lookup.StringVar(fs, "logfile")
	
	if err != nil {
		return fmt.Errorf("Failed to determine logfile flag, %w", err)
	}
	
	loglevel, err := lookup.StringVar(fs, "loglevel")
	
	if err != nil {
		return fmt.Errorf("Failed to determine loglevel flag, %w", err)
	}

	// This is mostly because I still can't figure out how
	// the get the current log level from log/slog and we need
	// to assign a default slog.Level to the logfile handler
	// below.
	var sloglevel slog.Level
	
	switch loglevel {
	case "debug":
		slog.SetLogLoggerLevel(slog.LevelDebug)
		sloglevel = slog.LevelDebug		
	case "info":
		slog.SetLogLoggerLevel(slog.LevelInfo)
		sloglevel = slog.LevelInfo
	case "status":
		slog.Warn("Log level is no longer supported, defaulting to INFO", "level", "loglevel")		
		slog.SetLogLoggerLevel(slog.LevelInfo)
		sloglevel = slog.LevelInfo		
	case "warning":
		slog.SetLogLoggerLevel(slog.LevelWarn)
		sloglevel = slog.LevelWarn		
	case "error", "fatal":
		slog.Warn("Log level is no longer supported, defaulting to WARN", "level", "loglevel")
		slog.SetLogLoggerLevel(slog.LevelWarn)
		sloglevel = slog.LevelWarn				
	default:
		sloglevel = slog.LevelInfo
	}

	verbose, err := lookup.BoolVar(fs, "verbose")

	if err != nil {
		return err
	}

	if verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		sloglevel = slog.LevelDebug
		slog.Debug("Verbose logging enabled")
	}

	if logfile != "" {

		wr, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)

		if err != nil {
			return fmt.Errorf("Failed to open '%s' for writing (logs), %w", logfile, err)
		}

		defer wr.Close()
		
		mw := io.MultiWriter(os.Stderr, wr)
		
		h := slog.NewTextHandler(mw, &slog.HandlerOptions{Level: sloglevel})
		slog.SetDefault(slog.New(h))
	}	
	
	// END OF logging

	config, err := iiifconfig.LoadConfigWithFlagSet(ctx, fs)

	if err != nil {
		return err
	}

	// END OF generate HTML for tiles
	
	generate_tiles_html, err := lookup.BoolVar(fs, "generate-tiles-html")
	
	if err != nil {
		return fmt.Errorf("Failed to determine generate-tiles-html flag, %w", err)
	}
	
	if generate_tiles_html {
		
		t.onCompleteFunc = func(cfg *iiifconfig.Config, src_id string, alt_id string, count int, err error){

			logger := slog.Default()
			logger = logger.With("source", src_id)
			logger = logger.With("alt", alt_id)
			
			if err != nil {
				logger.Warn("Skipping on complete func because error present", "error", err)
				return
			}

			logger.Info("Generate HTML index page for tiles")
			
			derivatives_cache, err := iiifcache.NewDerivativesCacheFromConfig(config)
			
			if err != nil {
				logger.Error("Failed to load derivatives cache from config", "error", err)
				return
			}
			
			write_assets := func(assets_fs embed.FS, assets []string) error {

				for _, fname := range assets {
					
					root := filepath.Join(alt_id, "assets")
					path := filepath.Join(root, fname)
					
					body, err := assets_fs.ReadFile(fname)
					
					if err != nil {
						return fmt.Errorf("Failed to read %s, %w", fname, err)
					}
					
					err = derivatives_cache.Set(path, body)
					
					if err != nil {
						return fmt.Errorf("Failed to write %s, %w", path, err)
					}
				}
				
				return nil
			}
			
			js_assets := []string{
				"leaflet.js",
				"leaflet-iiif.js",
			}
			
			css_assets := []string{
				"leaflet.css",
			}
			
			err = write_assets(javascript.FS, js_assets)
			
			if err != nil {
				logger.Error("Failed to write JS assets", "error", err)
				return
			}
			
			err = write_assets(css.FS, css_assets)
			
			if err != nil {
				logger.Error("Failed to write CSS assets", "error", err)
				return
			}
			
			root := filepath.Join(alt_id)
			path := filepath.Join(root, "index.html")
				
			body, err := html.FS.ReadFile("tiles.html")
			
			if err != nil {
				logger.Error("Failed to read tiles HTML", "error", err)
				return
			}
			
			err = derivatives_cache.Set(path, body)
			
			if err != nil {
				logger.Error("Failed to write HTML", "path", path, "error", err)
				return
			}
		}
		
	}

	// END OF generate HTML for tiles
	
	ts, err := iiiftile.NewTileSeed(config, 256, 256, endpoint, quality, format)

	if err != nil {
		return fmt.Errorf("Failed to create tileseed(er), %w", err)
	}

	// Something something something writers for slog stuff...

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

		slog.Debug("Tile waiting to seed", "source", seed.Source, "time", time.Since(t1))

		go func(seed *Seed, wg *sync.WaitGroup) {

			t1 := time.Now()

			src_id := seed.Source
			alt_id := seed.Target

			defer func() {
				slog.Debug("Time to seed tiles", "source", seed.Source, "time", time.Since(t1))
				throttle <- true
				wg.Done()
			}()

			count, err := ts.SeedTiles(src_id, alt_id, scales, refresh)

			if t.onCompleteFunc != nil {
				t.onCompleteFunc(config, src_id, alt_id, count, err)
			}

			if err != nil {
				slog.Warn("Failed to seed tiles", "id", src_id, "error", err)
			} else {
				slog.Debug("Seeded tiles complete", "id", src_id, "count", count)
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

		csv_bucket, err := bucket.OpenBucket(ctx, csv_source)

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
					slog.Warn("Unable to determine source ID for row", "row", row)
					continue
				}

				alt_id, ok := row["alternate_id"]

				if !ok {
					slog.Warn("Unable to determine alternate ID for row", "row", row)
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
		slog.Info("Watching filesystem", "root", root)

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
							slog.Warn("Failed to run URI function from path", "rel_path", rel_path, "abs_path", abs_path, "error", err)
							continue
						}

						seed, err := SeedFromURI(u, noextension)

						if err != nil {
							slog.Warn("Failed to determine seed from path", "rel_path", rel_path, "abs_path", abs_path, "error", err)
							continue
						}

						err = tile_func(seed, wg)

						if err != nil {
							slog.Warn("Failed to generate tiles from path", "rel_path", rel_path, "error", err)
							continue
						}
					}

				case err, ok := <-watcher.Errors:

					if !ok {
						return
					}

					slog.Warn("fsnotify error", "error", err)
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
