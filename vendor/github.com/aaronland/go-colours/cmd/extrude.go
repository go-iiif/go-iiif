package main

import (
	"flag"
	"github.com/aaronland/go-colours/extruder"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
)

func main() {

	flag.Parse()

	ex, err := extruder.NewNamedExtruder("vibrant")

	if err != nil {
		log.Fatal(err)
	}

	for _, path := range flag.Args() {

		f, err := os.Open(path)

		if err != nil {
			log.Fatal(err)
		}

		im, _, err := image.Decode(f)

		if err != nil {
			log.Fatal(err)
		}

		c, err := ex.Colours(im, 5)

		if err != nil {
			log.Fatal(err)
		}

		for _, c := range c {
			log.Println(c)
		}
	}
}
