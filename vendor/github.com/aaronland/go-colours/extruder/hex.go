package extruder

import (
	"fmt"
	"image/color"
)

func toHexColor(c color.Color) string {

	r, g, b, _ := c.RGBA()

	toS := func(i uint8) string {
		h := fmt.Sprintf("%x", i)
		if len(h) == 1 {
			h = "0" + h
		}
		return h
	}

	hex := toS(uint8(r)) + toS(uint8(g)) + toS(uint8(b))
	return fmt.Sprintf("#%s", hex)
}
