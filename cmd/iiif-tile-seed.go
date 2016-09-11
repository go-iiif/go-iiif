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
	"sync"
	"time"
)

func main() {

	var cfg = flag.String("config", ".", "config")
	var sf = flag.Int("scale-factor", 4, "...")
	var refresh = flag.Bool("refresh", false, "...")

	flag.Parse()

	config, err := iiifconfig.NewConfigFromFile(*cfg)

	if err != nil {
		log.Fatal(err)
	}

	level, err := iiiflevel.NewLevelFromConfig(config, "example.com")

	if err != nil {
		log.Fatal(err)
	}

	cache, err := iiifcache.NewDerivativesCacheFromConfig(config)

	if err != nil {
		log.Fatal(err)
	}

	cache.Get("foo")

	ids := flag.Args()

	ts, err := iiiftile.NewTileSeed(level, 256, 256)

	if err != nil {
		log.Fatal(err)
	}

	for _, id := range ids {

		log.Println(id)

		image, err := iiifimage.NewImageFromConfig(config, id)

		if err != nil {
			log.Fatal(err)
		}

		crops, err := ts.TileSizes(image, *sf)

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
					_, err := cache.Get(uri)

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
					cache.Set(uri, tmp.Body())
				}

			}(image, transformation, wg)
		}

		wg.Wait()

		tb := time.Since(ta)
		log.Printf("generated %d crops in %v", len(crops), tb)
	}
}
