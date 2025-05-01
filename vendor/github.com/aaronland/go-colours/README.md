# go-colours

Go package for working with colours, principally colour extraction and "snap to grid"

## Important

This is work in progress.

## Documentation

Documentation is incomplete.

## Example

```
package main

import (
	"context"
	"flag"
	"image"
	_ "image/jpeg"
	"log"
	"os"

	"github.com/aaronland/go-colours/extruder"
	"github.com/aaronland/go-colours/grid"
	"github.com/aaronland/go-colours/palette"
)

func main() {

	flag.Parse()

	ex, _ := extruder.NewExtruder(ctx, "vibrant://")
	gr, _ := grid.NewGrid(ctx, "euclidian://")
	p, _ := palette.NewPalette(ctx, "css4://")

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
	Name() string
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

### vibrant://

This returns colours using the [vibrant](github.com/RobCherry/vibrant) package.

Importantly, this uses the [sfomuseum/vibrant](https://github.com/sfomuseum/vibrant) fork of the package to enable the filtering out of transparent pixels.

### marekm4://

This returns colours using the [marekm4/color-extractor](https://github.com/marekm4/color-extractor) package.

### Grids

Grids are the things that perform operations or compare colours.

### euclidian://

### Palettes

Palettes are a fixed set of colours.

### crayola://

### css3://

### css4://

## See also

* https://github.com/RobCherry/vibrant
* https://github.com/marekm4/color-extractor
* https://github.com/givp/RoyGBiv
