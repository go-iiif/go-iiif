package main

import (
	"flag"
	iiifcache "github.com/thisisaaronland/go-iiif/cache"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	iiifimage "github.com/thisisaaronland/go-iiif/image"
	iiiflevel "github.com/thisisaaronland/go-iiif/level"
	iiifsource "github.com/thisisaaronland/go-iiif/source"
	iiiftile "github.com/thisisaaronland/go-iiif/tile"
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
	var refresh = flag.Bool("refresh", false, "Refresh a tile even if already exists (default false)")

	flag.Parse()

	if *cfg == "" {
		log.Fatal("Missing config file")
	}

	config, err := iiifconfig.NewConfigFromFile(*cfg)

	if err != nil {
		log.Fatal(err)
	}

	level, err := iiiflevel.NewLevelFromConfig(config, "example.com")

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

	ids := flag.Args()

	ts, err := iiiftile.NewTileSeed(level, 256, 256)

	if err != nil {
		log.Fatal(err)
	}

	for _, s := range strings.Split(*sf, ",") {

		s = strings.Trim(s, " ")
		scale, err := strconv.Atoi(s)

		if err != nil {
			log.Fatal(err)
		}

		for _, id := range ids {

			log.Println(id)

			image, err := iiifimage.NewImageFromConfigWithCache(config, images_cache, id)

			if err != nil {
				log.Fatal(err)
			}

			crops, err := ts.TileSizes(image, scale)

			if err != nil {
				log.Fatal(err)
			}

			source, err := iiifsource.NewMemorySource(image.Body())

			if err != nil {
				log.Fatal(err)
			}

			procs := runtime.NumCPU() * 2

			ch := make(chan bool, procs)

			for i := 0; i < procs; i++ {
				ch <- true
			}

			wg := new(sync.WaitGroup)
			ta := time.Now()

			for _, transformation := range crops {

				wg.Add(1)

				go func(im iiifimage.Image, tr *iiifimage.Transformation, wg *sync.WaitGroup) {

					<-ch

					defer func() {
						wg.Done()
						ch <- true
					}()

					uri, _ := tr.ToURI(im.Identifier())

					if !*refresh {
						_, err := derivatives_cache.Get(uri)

						if err == nil {
							return
						}
					}

					tmp, _ := iiifimage.NewImageFromConfigWithSource(config, source, "cache")

					t1 := time.Now()
					err = tmp.Transform(tr)
					t2 := time.Since(t1)

					log.Println(uri, t2, err)

					if err == nil {
						derivatives_cache.Set(uri, tmp.Body())
					}

				}(image, transformation, wg)
			}

			wg.Wait()

			// generate info.json here – https://github.com/thisisaaronland/go-iiif/issues/7
			tb := time.Since(ta)
			log.Printf("generated %d crops in %v", len(crops), tb)
		}
	}
}
