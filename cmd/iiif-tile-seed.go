package main

import (
	"flag"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	iiifimage "github.com/thisisaaronland/go-iiif/image"
	iiiflevel "github.com/thisisaaronland/go-iiif/level"
	iiiftile "github.com/thisisaaronland/go-iiif/tile"
	"log"
)

func main() {

	var cfg = flag.String("config", ".", "config")
	var sf = flag.Int("scale-factor", 4, "...")

	flag.Parse()

	config, err := iiifconfig.NewConfigFromFile(*cfg)

	if err != nil {
		log.Fatal(err)
	}

	level, err := iiiflevel.NewLevelFromConfig(config, "example.com")

	if err != nil {
		log.Fatal(err)
	}

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

		_, err = ts.TileSizes(image, *sf)

		if err != nil {
			log.Fatal(err)
		}
	}
}
