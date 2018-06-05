package main

import (
	"flag"
	_ "github.com/aaronland/go-colours"
	"github.com/aaronland/go-colours/extruder"
	"github.com/aaronland/go-colours/grid"
	"github.com/aaronland/go-colours/palette"
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

	gr, err := grid.NewNamedGrid("euclidian")

	if err != nil {
		log.Fatal(err)
	}

	p, err := palette.NewNamedPalette("css4")

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

		log.Println(path)

		c, err := ex.Colours(im, 5)

		if err != nil {
			log.Fatal(err)
		}

		for _, c := range c {
			log.Println(c)

			cl, _ := gr.Closest(c, p)

			log.Println(cl)
		}

	}
}
