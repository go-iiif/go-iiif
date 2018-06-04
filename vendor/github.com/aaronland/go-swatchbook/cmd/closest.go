package main

import (
	"flag"
	"github.com/aaronland/go-swatchbook"
	"log"
)

func main() {

	flag.Parse()

	p, err := swatchbook.NewNamedPalette("css4")

	if err != nil {
		log.Fatal(err)
	}

	s, err := swatchbook.NewSwatchbookFromPalette(p)

	if err != nil {
		log.Fatal(err)
	}

	for _, h := range flag.Args() {

		target := &swatchbook.Color{
			Name: h,
			Hex:  h,
		}

		match := s.Closest(target)

		log.Printf("%s snaps to %s\n", target, match)
	}

}
