# go-colours

Go package for working with colours, principally colour extraction and "snap to grid"

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

	ctx := context.Background()
	
	ex, _ := extruder.NewExtruder(ctx, "vibrant://")
	gr, _ := grid.NewGrid(ctx, "euclidian://")
	p, _ := palette.NewPalette(ctx, "css4://")

	for _, path := range flag.Args() {

		fh, _ := os.Open(path)
		im, _, _ := image.Decode(fh)

		colours, _ := ex.Colours(ctx, im, 5)

		for _, c := range colours {

			closest, _ := gr.Closest(ctx, c, p)

			for _, cl := range closest {
				log.Println(c, cl)
			}
		}

	}
}
```

_Note that error handling has been removed for the sake of brevity._

## Tools

```
$> make cli
go build -mod vendor -ldflags="-s -w" -o bin/extrude cmd/extrude/main.go
go build -mod vendor -ldflags="-s -w" -o bin/inspect cmd/inspect/main.go
go build -mod vendor -ldflags="-s -w" -o bin/snap cmd/snap/main.go
go build -mod vendor -ldflags="-s -w" -o bin/review cmd/review/main.go
```

### extrude

Command line tool to extrude (derive) dominant colours from one or more images as well as closest matches colours using zero or more "snap-to-grid" colour palettes as JSON-encoded data written to STDOUT.

```
$> ./bin/extrude -h
Command line tool to extrude (derive) dominant colours from one or more images as well as closest matches colours using zero or more "snap-to-grid" colour palettes as JSON-encoded data written to STDOUT.
Usage:
	 ./bin/extrude [options] uri(N) uri(N)
  -allow-remote
    	Allow fetching remote images (HTTP(S)). (default true)
  -extruder-uri value
    	Zero or more aaronland/go-colours/extruder.Extruder URIs. Default is to use all registered extruder schemes.
  -palette-uri value
    	Zero or more aaronland/go-colours/palette.Palette URIs. Default is to use all registered palette schemes.
  -verbose
    	Enable verbose (debug) logging.
```	

#### Example

```
$> ./bin/extrude  https://static.sfomuseum.org/media/176/270/453/3/1762704533_jnxsOwjYqsa8RyGsJrYFJvAjnQMe1Nqv_z.jpg | jq
[
  {
    "uri": "1762704533_jnxsOwjYqsa8RyGsJrYFJvAjnQMe1Nqv_z.png",
    "extrusions": [
      {
        "extruder": "marekm4",
        "palettes": [
          "crayola",
          "css3",
          "css4"
        ],
        "swatches": [
          {
            "colour": {
              "name": "marekm4",
              "hex": "#b6baa1",
              "reference": "b6baa1"
            },
            "closest": [
              {
                "palette": "crayola",
                "colour": {
                  "name": "Cadet Blue",
                  "hex": "#b0b7c6",
                  "reference": "crayola"
                }
              },
              {
                "palette": "css3",
                "colour": {
                  "name": "darkgray",
                  "hex": "#a9a9a9",
                  "reference": "css3"
                }
              },
              {
                "palette": "css4",
                "colour": {
                  "name": "darkgrey",
                  "hex": "#a9a9a9",
                  "reference": "css4"
                }
              }
            ]
          },
          {
            "colour": {
              "name": "marekm4",
              "hex": "#728c9a",
              "reference": "728c9a"
            },
	    ...and so on
```

### review

Command line tool to perform colour extraction and "snap-to-grid" matching with one or more colour palettes for images, emitting the results as an HTML page.

```
$> ./bin/review -h
Command line tool to perform colour extraction and "snap-to-grid" matching with one or more colour palettes for images, emitting the results as an HTML page.
Usage:
	 ./bin/review [options] uri(N) uri(N)
  -allow-remote
    	Allow fetching remote images (HTTP(S)). (default true)
  -extruder-uri value
    	Zero or more aaronland/go-colours/extruder.Extruder URIs. Default is to use all registered extruder schemes.
  -palette-uri value
    	Zero or more aaronland/go-colours/palette.Palette URIs. Default is to use all registered palette schemes.
  -root string
    	The path to a directory where images and HTML files associated with the review should be stored. If empty a new temporary directory will be created (and deleted when the application exits).
  -verbose
    	Enable verbose (debug) logging.
```

#### Example

```
$> ./bin/review  https://static.sfomuseum.org/media/176/270/453/3/1762704533_jnxsOwjYqsa8RyGsJrYFJvAjnQMe1Nqv_z.jpg
2025/05/01 17:18:28 INFO Server is ready and features are viewable url=http://localhost:50530
```

And then when you open your browser to `http://localhost:50530` (or whatever address the `review` tool chooses) you'd see something like this:

![](docs/images/go-colours-review.png)

Unless you specify a custom `-root` flag all the images used by the web application (excluding the source images themselves) will be automatically be deleted when you stop the tool.

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

### marekm4://

This returns colours using the [marekm4/color-extractor](https://github.com/marekm4/color-extractor) package.

### quant://

This returns colours using the [soniakeys/quant](https://github.com/soniakeys/quant) package (specifically the [mean.Quantizer](https://pkg.go.dev/github.com/soniakeys/quant@v1.0.0/mean#Quantizer)).

### vibrant://

This returns colours using the [vibrant](github.com/RobCherry/vibrant) package.

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
* https://github.com/soniakeys/quant
* https://github.com/givp/RoyGBiv
