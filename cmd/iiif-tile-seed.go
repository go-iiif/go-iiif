package main

import (
	"flag"
	"fmt"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	iiiftile "github.com/thisisaaronland/go-iiif/tile"
	csv "github.com/whosonfirst/go-whosonfirst-csv"
	"io"
	"log"
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
	var mode = flag.String("mode", "-", "Whether to read input as a CSV file or from STDIN which can be represented as \"-\"")
	var refresh = flag.Bool("refresh", false, "Refresh a tile even if already exists (default false)")
	var endpoint = flag.String("endpoint", "http://localhost:8080", "The endpoint (scheme, host and optionally port) that will serving these tiles, used for generating an 'info.json' for each source image")

	flag.Parse()

	if *cfg == "" {
		log.Fatal("Missing config file")
	}

	config, err := iiifconfig.NewConfigFromFile(*cfg)

	if err != nil {
		log.Fatal(err)
	}

	ts, err := iiiftile.NewTileSeed(config, 256, 256, *endpoint, *quality, *format)

	if err != nil {
		log.Fatal(err)
	}

	scales := make([]int, 0)

	for _, s := range strings.Split(*sf, ",") {

		s = strings.Trim(s, " ")
		scale, err := strconv.Atoi(s)

		if err != nil {
			log.Fatal(err)
		}

		scales = append(scales, scale)
	}

	if *mode == "csv" {

		procs := runtime.NumCPU()

		ch := make(chan bool, procs)

		for i := 0; i < procs; i++ {
			ch <- true
		}

		for _, path := range flag.Args() {

			reader, err := csv.NewDictReaderFromPath(path)

			if err != nil {
				log.Fatal(err)
			}

			wg := new(sync.WaitGroup)

			for {
				row, err := reader.Read()

				if err == io.EOF {
					break
				}

				if err != nil {
					log.Fatal(err)
				}

				src_id, ok := row["source_id"]

				if !ok {
					log.Println("Unable to determine source ID", row)
					continue
				}

				alt_id, ok := row["alternate_id"]

				if !ok {
					log.Println("Unable to determine alternate ID", row)
					continue
				}

				wg.Add(1)

				go func(src_id string, alt_id string) {

					<-ch

					defer wg.Done()

					t1 := time.Now()

					count, err := ts.SeedTiles(src_id, alt_id, scales, *refresh)

					t2 := time.Since(t1)

					if err != nil {
						fmt.Println("ERROR", src_id, t2, err)
					} else {
						fmt.Println("OKAY", src_id, count, t2)
					}

					ch <- true

				}(src_id, alt_id)
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

			t1 := time.Now()

			count, err := ts.SeedTiles(src_id, alt_id, scales, *refresh)

			t2 := time.Since(t1)

			if err != nil {
				log.Fatal(err)
			}

			log.Printf("[%s] time to process %d tiles: %v\n", id, count, t2)
		}
	}
}
