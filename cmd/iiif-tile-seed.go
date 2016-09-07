package main

import (
	"flag"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	iiifimage "github.com/thisisaaronland/go-iiif/image"
	iiiftile "github.com/thisisaaronland/go-iiif/tile"
	"log"
)

func main() {

	var cfg = flag.String("config", ".", "config")

	flag.Parse()

	config, err := iiifconfig.NewConfigFromFile(*cfg)

	if err != nil {
		log.Fatal(err)
	}

	ids := flag.Args()

	ts, err := iiiftile.NewTileSeed(256, 256)

	if err != nil {
		log.Fatal(err)
	}

	for _, id := range ids {

		image, err := iiifimage.NewImageFromConfig(config, id)

		if err != nil {
			log.Fatal(err)
		}

		ts.TileSizes(image, 4)
	}
}
