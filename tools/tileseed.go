package tools

import (
	"context"
	"errors"
	"flag"
	"fmt"
	aws_events "github.com/aws/aws-lambda-go/events"
	aws_lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/fsnotify/fsnotify"
	iiifuri "github.com/go-iiif/go-iiif-uri"
	iiifconfig "github.com/go-iiif/go-iiif/config"
	iiiftile "github.com/go-iiif/go-iiif/tile"
	"github.com/whosonfirst/go-whosonfirst-cli/flags"
	"github.com/whosonfirst/go-whosonfirst-csv"
	"github.com/whosonfirst/go-whosonfirst-log"
	"gocloud.dev/blob"
	"io"
	golog "log"
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

type TileSeedTool struct {
	Tool
	URIFunc URIFunc
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

func NewTileSeedToolWithURIFunc(uri_func URIFunc) (Tool, error) {
	t := &TileSeedTool{
		URIFunc: uri_func,
	}
	return t, nil
}

func (t *TileSeedTool) Run(ctx context.Context) error {

	var cfg = flag.String("config", "", "Path to a valid go-iiif config file. DEPRECATED - please use -config-source and -config name.")

	var config_source = flag.String("config-source", "", "A valid Go Cloud bucket URI where your go-iiif config file is located.")
	var config_name = flag.String("config-name", "config.json", "The name of your go-iiif config file.")

	var csv_source = flag.String("csv-source", "A valid Go Cloud bucket URI where your CSV tileseed files are located.", "")

	var sf = flag.String("scale-factors", "4", "A comma-separated list of scale factors to seed tiles with")
	var quality = flag.String("quality", "default", "A valid IIIF quality parameter - if \"default\" then the code will try to determine which format you've set as the default")
	var format = flag.String("format", "jpg", "A valid IIIF format parameter")
	var logfile = flag.String("logfile", "", "Write logging information to this file")
	var loglevel = flag.String("loglevel", "info", "The amount of logging information to include, valid options are: debug, info, status, warning, error, fatal")
	var processes = flag.Int("processes", runtime.NumCPU(), "The number of concurrent processes to use when tiling images")
	var mode = flag.String("mode", "cli", "Valid modes are: cli, csv, fsnotify, lambda.")

	var noextension = flag.Bool("noextension", false, "Remove any extension from destination folder name.")

	var refresh = flag.Bool("refresh", false, "Refresh a tile even if already exists (default false)")
	var endpoint = flag.String("endpoint", "http://localhost:8080", "The endpoint (scheme, host and optionally port) that will serving these tiles, used for generating an 'info.json' for each source image")
	var verbose = flag.Bool("verbose", false, "Write logging to STDOUT in addition to any other log targets that may have been defined")

	flag.Parse()

	err := flags.SetFlagsFromEnvVars("IIIF_TILESEED")

	if err != nil {
		return err
	}

	if *cfg != "" {

		golog.Println("-config flag is deprecated. Please use -config-source and -config-name (setting them now).")

		abs_config, err := filepath.Abs(*cfg)

		if err != nil {
			return err
		}

		*config_name = filepath.Base(abs_config)
		*config_source = fmt.Sprintf("file://%s", filepath.Dir(abs_config))
	}

	if *config_source == "" {
		return errors.New("Required -config-source flag is empty.")
	}

	config_bucket, err := blob.OpenBucket(ctx, *config_source)

	if err != nil {
		return err
	}

	config, err := iiifconfig.NewConfigFromBucket(ctx, config_bucket, *config_name)

	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	ts, err := iiiftile.NewTileSeed(config, 256, 256, *endpoint, *quality, *format)

	if err != nil {
		return err
	}

	writers := make([]io.Writer, 0)

	if *verbose {
		writers = append(writers, os.Stdout)
	}

	if *logfile != "" {

		fh, err := os.OpenFile(*logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)

		if err != nil {
			return err
		}

		writers = append(writers, fh)
	}

	writer := io.MultiWriter(writers...)

	logger := log.NewWOFLogger("")
	logger.AddLogger(writer, *loglevel)

	scales := make([]int, 0)

	for _, s := range strings.Split(*sf, ",") {

		s = strings.Trim(s, " ")
		scale, err := strconv.Atoi(s)

		if err != nil {
			return err
		}

		scales = append(scales, scale)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	throttle := make(chan bool, *processes)

	for i := 0; i < *processes; i++ {
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

			count, err := ts.SeedTiles(src_id, alt_id, scales, *refresh)

			if err != nil {
				logger.Warning("Failed to seed tiles for '%s', %s", src_id, err)
			} else {
				logger.Debug("Seeded %d tiles for '%s'", count, src_id)
			}

		}(seed, wg)

		return nil
	}

	switch *mode {
	case "cli", "-":

		wg := new(sync.WaitGroup)

		for _, id := range flag.Args() {

			u, err := t.URIFunc(id)

			if err != nil {
				logger.Fatal(err)
			}

			seed, err := SeedFromURI(u, *noextension)

			if err != nil {
				logger.Fatal(err)
			}

			tile_func(seed, wg)
		}

		wg.Wait()

	case "csv":

		csv_bucket, err := blob.OpenBucket(ctx, *csv_source)

		if err != nil {
			return err
		}

		wg := new(sync.WaitGroup)

		for _, path := range flag.Args() {

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

						u, err := t.URIFunc(rel_path)

						if err != nil {
							logger.Warning("Failed to run URI function from path '%s' (%s), %s", rel_path, abs_path, err)
							continue
						}

						seed, err := SeedFromURI(u, *noextension)

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

				// TILE FUNC HERE...

				seed, err := SeedFromString(s3_fname, *noextension)

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
