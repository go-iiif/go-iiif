# go-swatchbook

Go tools for wrangling palette-based colours.

## Important 

This is a port of the [py-cooperhewitt-swatchbook](https://github.com/cooperhewitt/py-cooperhewitt-swatchbook) package, but it is not feature complete yet. Under the hood this package uses the [go-hexcolor](github.com/pwaller/go-hexcolor) package but that might change, yet...

## Install

You will need to have both `Go` (specifically a version of Go more recent than 1.6 so let's just assume you need [Go 1.8](https://golang.org/dl/) or higher) and the `make` programs installed on your computer. Assuming you do just type:

```
make bin
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Example

```
package main

import (
	"flag"
	"github.com/aaronland/go-swatchbook"
	"log"
)

func main() {

	flag.Parse()

	p, _ := swatchbook.NewNamedPalette("css4")
	s, _ := swatchbook.NewSwatchbookFromPalette(p)

	for _, h := range flag.Args() {

		target := &swatchbook.Color{
			Name: h,
			Hex:  h,
		}

		match := s.Closest(target)
		log.Printf("%s snaps to %s\n", target, match)
	}
}
```

_Error handling has been removed for the sake of brevity._

## See also

* https://github.com/cooperhewitt/py-cooperhewitt-swatchbook
* github.com/pwaller/go-hexcolor