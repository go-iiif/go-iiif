package main

import (
	"encoding/json"
	"flag"
	"fmt"
	iiifcache "github.com/thisisaaronland/go-iiif/cache"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	iiifimage "github.com/thisisaaronland/go-iiif/image"
	iiiflevel "github.com/thisisaaronland/go-iiif/level"
	iiifprofile "github.com/thisisaaronland/go-iiif/profile"
	iiifsource "github.com/thisisaaronland/go-iiif/source"
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
	var mode = flag.String("mode", "-", "...")
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

	level, err := iiiflevel.NewLevelFromConfig(config, *endpoint)

	if err != nil {
		log.Fatal(err)
	}

	images_cache, err := iiifcache.NewImagesCacheFromConfig(config)

	if err != nil {
		log.Fatal(err)
	}

	derivatives_cache, err := iiifcache.NewDerivativesCacheFromConfig(config)

	if err != nil {
		log.Fatal(err)
	}

	ts, err := iiiftile.NewTileSeed(level, 256, 256)

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

				src_id, _ := row["src"]
				dest_id, _ := row["dest"]

				wg.Add(1)

				go func() {

					<-ch

					defer wg.Done()

					t1 := time.Now()

					count, err := SeedTiles(ts, src_id, dest_id, config, images_cache, derivatives_cache, scales, *endpoint, *refresh)

					t2 := time.Since(t1)

					if err != nil {
						fmt.Println("ERROR", src_id, t2, err)
					} else {
						fmt.Println("OKAY", src_id, count, t2)
					}

					ch <- true

				}()
			}

			wg.Wait()
		}

	} else {

		for _, id := range flag.Args() {

			var src_id string
			var dest_id string

			pointers := strings.Split(id, ",")

			if len(pointers) == 2 {
				src_id = pointers[0]
				dest_id = pointers[1]
			} else {
				src_id = pointers[0]
				dest_id = pointers[0]
			}

			_, err = SeedTiles(ts, src_id, dest_id, config, images_cache, derivatives_cache, scales, *endpoint, *refresh)

			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

/*

See this - the function signature is complete madness. We'll figure it out...
(20160913/thisisaaronland)

*/

func SeedTiles(ts *iiiftile.TileSeed, src_id string, dest_id string, config *iiifconfig.Config, images_cache iiifcache.Cache, derivatives_cache iiifcache.Cache, scales []int, endpoint string, refresh bool) (int, error) {

	count := 0

	image, err := iiifimage.NewImageFromConfigWithCache(config, images_cache, src_id)

	if err != nil {
		return count, err
	}

	for _, scale := range scales {

		crops, err := ts.TileSizes(image, scale)

		if err != nil {
			log.Println(err)
			continue
		}

		source, err := iiifsource.NewMemorySource(image.Body())

		if err != nil {
			return count, err
		}

		procs := runtime.NumCPU() * 2

		ch := make(chan bool, procs)

		for i := 0; i < procs; i++ {
			ch <- true
		}

		wg := new(sync.WaitGroup)

		for _, transformation := range crops {

			wg.Add(1)

			go func(im iiifimage.Image, tr *iiifimage.Transformation, wg *sync.WaitGroup) {

				<-ch

				defer func() {
					wg.Done()
					ch <- true
				}()

				uri, _ := tr.ToURI(dest_id)

				if !refresh {

					_, err := derivatives_cache.Get(uri)

					if err == nil {
						return
					}
				}

				tmp, _ := iiifimage.NewImageFromConfigWithSource(config, source, im.Identifier())

				err = tmp.Transform(tr)

				if err == nil {
					derivatives_cache.Set(uri, tmp.Body())
				}

			}(image, transformation, wg)
		}

		wg.Wait()

		// something something something using the channel above to increment count...

		count += len(crops)
	}

	level, err := iiiflevel.NewLevelFromConfig(config, endpoint)

	if err != nil {
		return count, err
	}

	profile, err := iiifprofile.NewProfile(endpoint, image, level)

	if err != nil {
		return count, err
	}

	body, err := json.Marshal(profile)

	if err != nil {
		return count, err
	}

	uri := fmt.Sprintf("%s/info.json", dest_id)
	derivatives_cache.Set(uri, body)

	return count, nil
}
