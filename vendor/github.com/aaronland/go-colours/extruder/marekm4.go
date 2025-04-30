package extruder

import (
	"context"
	"fmt"
	"image"
	"image/color"
	_ "log/slog"

	"github.com/aaronland/go-colours"
	"github.com/marekm4/color-extractor"
)

type Marekm4Extruder struct {
	Extruder
}

func init() {
	ctx := context.Background()
	err := RegisterExtruder(ctx, "marekm4", NewMarekm4Extruder)
	if err != nil {
		panic(err)
	}
}

func NewMarekm4Extruder(ctx context.Context, uri string) (Extruder, error) {

	ex := Marekm4Extruder{}

	return &ex, nil
}

func (ex *Marekm4Extruder) Name() string {
	return "marekm4"
}

func (ex *Marekm4Extruder) Colours(im image.Image, limit int) ([]colours.Colour, error) {

	ctx := context.Background()
	rsp := color_extractor.ExtractColors(im)

	if len(rsp) < limit {
		limit = len(rsp)
	}

	results := make([]colours.Colour, limit)

	for i := 0; i < limit; i++ {

		c := rsp[i]

		hex_value := toHexColor(c)
		c_uri := fmt.Sprintf("common://?hex=%s", hex_value)

		colour, err := colours.NewColour(ctx, c_uri)

		if err != nil {
			return nil, err
		}

		results[i] = colour
	}

	return results, nil
}

func toHexColor(c color.Color) string {

	r, g, b, _ := c.RGBA()

	toS := func(i uint8) string {
		h := fmt.Sprintf("%x", i)
		if len(h) == 1 {
			h = "0" + h
		}
		return h
	}

	return toS(uint8(r)) + toS(uint8(g)) + toS(uint8(b))
}
