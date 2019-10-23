package tools

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/aaronland/gocloud-blob-bucket"
	aws_events "github.com/aws/aws-lambda-go/events"
	aws_lambda "github.com/aws/aws-lambda-go/lambda"
	iiifconfig "github.com/go-iiif/go-iiif/config"
	iiifdriver "github.com/go-iiif/go-iiif/driver"
	iiiftile "github.com/go-iiif/go-iiif/tile"
	"github.com/whosonfirst/go-whosonfirst-csv"
	"github.com/whosonfirst/go-whosonfirst-log"
	"io"
	golog "log"
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
}

func SeedFromString(id string, no_extension bool) *Seed {

	var src_id string
	var alt_id string

	pointers := strings.Split(id, ",")

	if len(pointers) == 2 {
		src_id = pointers[0]
		alt_id = pointers[1]
	} else {
		src_id = pointers[0]
		alt_id = pointers[0]
	}

	if no_extension {
		alt_id = strings.TrimSuffix(alt_id, filepath.Ext(alt_id))
	}

	seed := &Seed{
		Source: src_id,
		Target: alt_id,
	}

	return seed
}

func NewTileSeedTool() (Tool, error) {

	t := &TileSeedTool{}
	return t, nil
}

func (t *TileSeedTool) Run(ctx context.Context) error {

	var cfg = flag.String("config", "", "Path to a valid go-iiif config file. DEPRECATED - please use -config_source and -config name.")

	var config_source = flag.String("config-source", "", "")
	var config_name = flag.String("config-name", "config.json", "")

	var sf = flag.String("scale-factors", "4", "A comma-separated list of scale factors to seed tiles with")
	var quality = flag.String("quality", "default", "A valid IIIF quality parameter - if \"default\" then the code will try to determine which format you've set as the default")
	var format = flag.String("format", "jpg", "A valid IIIF format parameter")
	var logfile = flag.String("logfile", "", "Write logging information to this file")
	var loglevel = flag.String("loglevel", "info", "The amount of logging information to include, valid options are: debug, info, status, warning, error, fatal")
	var processes = flag.Int("processes", runtime.NumCPU(), "The number of concurrent processes to use when tiling images")
	var mode = flag.String("mode", "cli", "...")

	var noextension = flag.Bool("noextension", false, "Remove any extension from destination folder name.")

	var refresh = flag.Bool("refresh", false, "Refresh a tile even if already exists (default false)")
	var endpoint = flag.String("endpoint", "http://localhost:8080", "The endpoint (scheme, host and optionally port) that will serving these tiles, used for generating an 'info.json' for each source image")
	var verbose = flag.Bool("verbose", false, "Write logging to STDOUT in addition to any other log targets that may have been defined")

	flag.Parse()

	if *cfg != "" {

		golog.Println("-config flag is deprecated. Please use -config-source and -config-name (setting them now).")

		abs_config, err := filepath.Abs(*cfg)

		if err != nil {
			return err
		}

		*config_name = filepath.Base(abs_config)
		*config_source = fmt.Sprintf("file://%s", filepath.Dir(abs_config))
	}

	config_bucket, err := bucket.OpenBucket(ctx, *config_source)

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

	driver, err := iiifdriver.NewDriverFromConfig(config)

	if err != nil {
		return err
	}

	ts, err := iiiftile.NewTileSeed(config, driver, 256, 256, *endpoint, *quality, *format)

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
			seed := SeedFromString(id, *noextension)
			tile_func(seed, wg)
		}

		wg.Wait()

	case "csv":

		wg := new(sync.WaitGroup)

		for _, path := range flag.Args() {

			// READ FROM BUCKET...
			reader, err := csv.NewDictReaderFromPath(path)

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

	case "lambda":

		handler := func(ctx context.Context, ev aws_events.S3Event) error {

			wg := new(sync.WaitGroup)

			for _, r := range ev.Records {

				s3_entity := r.S3
				s3_obj := s3_entity.Object
				s3_key := s3_obj.Key

				s3_fname := filepath.Base(s3_key)

				seed := SeedFromString(s3_fname, *noextension)
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
