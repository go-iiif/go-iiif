package main

import (
	"flag"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	iiiftile "github.com/thisisaaronland/go-iiif/tile"
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

func main() {

	var cfg = flag.String("config", "", "Path to a valid go-iiif config file")
	var sf = flag.String("scale-factors", "4", "A comma-separated list of scale factors to seed tiles with")
	var quality = flag.String("quality", "default", "A valid IIIF quality parameter - if \"default\" then the code will try to determine which format you've set as the default")
	var format = flag.String("format", "jpg", "A valid IIIF format parameter")
	var logfile = flag.String("logfile", "", "Write logging information to this file")
	var loglevel = flag.String("loglevel", "info", "The amount of logging information to include, valid options are: debug, info, status, warning, error, fatal")
	var processes = flag.Int("processes", runtime.NumCPU(), "The number of concurrent processes to use when tiling images")
	var mode = flag.String("mode", "-", "Whether to read input as a CSV file or from STDIN which can be represented as \"-\"")
	var noextension = flag.Bool("noextension", false, "Remove any extension from destination folder name.")
	var refresh = flag.Bool("refresh", false, "Refresh a tile even if already exists (default false)")
	var endpoint = flag.String("endpoint", "http://localhost:8080", "The endpoint (scheme, host and optionally port) that will serving these tiles, used for generating an 'info.json' for each source image")
	var verbose = flag.Bool("verbose", false, "Write logging to STDOUT in addition to any other log targets that may have been defined")

	flag.Parse()

	if *cfg == "" {
		golog.Fatal("Missing config file")
	}

	config, err := iiifconfig.NewConfigFromFlag(*cfg)

	if err != nil {
		golog.Fatal(err)
	}

	ts, err := iiiftile.NewTileSeed(config, 256, 256, *endpoint, *quality, *format)

	if err != nil {
		golog.Fatal(err)
	}

	writers := make([]io.Writer, 0)

	if *verbose {
		writers = append(writers, os.Stdout)
	}

	if *logfile != "" {

		fh, err := os.OpenFile(*logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)

		if err != nil {
			golog.Fatal(err)
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
			logger.Fatal(err.Error())
		}

		scales = append(scales, scale)
	}

	if *mode == "csv" {

		throttle := make(chan bool, *processes)

		for i := 0; i < *processes; i++ {
			throttle <- true
		}

		for _, path := range flag.Args() {

			reader, err := csv.NewDictReaderFromPath(path)

			if err != nil {
				logger.Fatal(err.Error())
			}

			wg := new(sync.WaitGroup)

			counter := 0

			for {

				row, err := reader.Read()
				counter += 1

				if err == io.EOF {
					break
				}

				if err != nil {
					logger.Fatal(err.Error())
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

				t1 := time.Now()

				<-throttle

				t2 := time.Since(t1)
				logger.Debug("%d time spent waiting to parse %s, %v", counter, src_id, t2)

				wg.Add(1)

				go func(throttle chan bool, src_id string, alt_id string) {

					defer wg.Done()

					t1 := time.Now()

					count, err := ts.SeedTiles(src_id, alt_id, scales, *refresh)

					t2 := time.Since(t1)

					logger.Debug("%d time to process %s (%d tiles), %v", counter, src_id, count, t2)

					if err != nil {
						logger.Error("FAILED to tile %s (%d) because %s, in %v", src_id, counter, err, t2)
					} else {
						logger.Status("SUCCESS tiling %s (%d) %d tiles in %v", src_id, counter, count, t2)
					}

					throttle <- true

				}(throttle, src_id, alt_id)

			}

			wg.Wait()
		}

	} else {

		for _, id := range flag.Args() {

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

			if *noextension {
				alt_id = strings.TrimSuffix(alt_id, filepath.Ext(alt_id))
			}

			t1 := time.Now()

			count, err := ts.SeedTiles(src_id, alt_id, scales, *refresh)

			t2 := time.Since(t1)

			if err != nil {
				logger.Fatal(err.Error())
			}

			logger.Debug("%s time to process %d tiles: %v", id, count, t2)
		}
	}
}
