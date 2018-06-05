package main

import (
	"flag"
	"github.com/aaronland/go-colours"
	"github.com/aaronland/go-colours/grid"
	"github.com/aaronland/go-colours/palette"
	"log"
)

func main() {

	flag.Parse()

	gr, err := grid.NewNamedGrid("euclidian")

	if err != nil {
		log.Fatal(err)
	}

	p, err := palette.NewNamedPalette("css4")

	if err != nil {
		log.Fatal(err)
	}

	for _, hex := range flag.Args() {

		target, err := colours.NewColour(hex)

		if err != nil {
			log.Fatal(err)
		}

		match, err := gr.Closest(target, p)

		if err != nil {
			log.Fatal(err)
		}

		log.Printf("%s SNAPS TO %s\n", target, match)
	}
}
