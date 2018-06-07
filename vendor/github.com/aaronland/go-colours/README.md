# go-colours

Go package for working with colours, principally colour extraction and "snap to grid"

## Install

You will need to have both `Go` (specifically a version of Go more recent than 1.6 so let's just assume you need [Go 1.8](https://golang.org/dl/) or higher) and the `make` programs installed on your computer. Assuming you do just type:

```
make bin
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Important

This is work in progress. Eventually it will be a complete port of the [py-cooperhewitt-swatchbook](https://github.com/aaronland/py-cooperhewitt-swatchbook) and [py-cooperhewitt-roboteyes-colors](https://github.com/aaronland/py-cooperhewitt-roboteyes-colors) (and by extension [RoyGBiv](https://github.com/givp/RoyGBiv)) packages, but today it is only a partial implementation.

Also, this documentation is incomplete.

## Example

```
package main

import (
	"flag"
	"github.com/aaronland/go-colours/extruder"
	"github.com/aaronland/go-colours/grid"
	"github.com/aaronland/go-colours/palette"
	"image"
	_ "image/jpeg"
	"log"
	"os"
)

func main() {

	flag.Parse()

	ex, _ := extruder.NewNamedExtruder("vibrant")

	gr, _ := grid.NewNamedGrid("euclidian")

	p, _ := palette.NewNamedPalette("css4")

	for _, path := range flag.Args() {

		fh, _ := os.Open(path)
		im, _, _ := image.Decode(fh)

		colours, _ := ex.Colours(im, 5)

		for _, c := range colours {

			closest, _ := gr.Closest(c, p)

			for _, cl := range closest {
				log.Println(c, cl)
			}
		}

	}
}
```

_Note that error handling has been removed for the sake of brevity._

## Interfaces

### Colour

```
type Colour interface {
	Name() string
	Hex() string
	Reference() string
	Closest() []Colour
	AppendClosest(Colour) error // I don't love this... (20180605/thisisaaronland)
	String() string
}
```

### Extruder

```
type Extruder interface {
	Colours(image.Image, int) ([]Colour, error)
}
```

### Grid

```
type Grid interface {
	Closest(Colour, Palette) (Colour, error)
}
```

### Palette

```
type Palette interface {
	Reference() string
	Colours() []Colour
}
```

## Extruders

Extruders are the things that generate a palette of colours for an `image.Image`.

### vibrant

This returns colours using the [vibrant](github.com/RobCherry/vibrant) package but rather than ranking colours using a particular metric it returns specific named "swatches" that are recast as `colours.Colour` interfaces. They are: `VibrantSwatch, LightVibrantSwatch, DarkVibrantSwatch, MutedSwatch, LightMutedSwatch, DarkMutedSwatch`.

### Grids

Grids are the things that perform operations or compare colours.

### euclidian

### Palettes

Palettes are a fixed set of colours.

### crayola

### css3

### css4

## See also

* https://github.com/RobCherry/vibrant
* https://github.com/lucasb-eyer/go-colorful

* https://github.com/aaronland/py-cooperhewitt-swatchbook
* https://github.com/aaronland/py-cooperhewitt-roboteyes-colors
* https://github.com/givp/RoyGBiv
