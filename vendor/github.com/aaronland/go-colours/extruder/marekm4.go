package extruder

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"net/url"

	"github.com/aaronland/go-colours"
	"github.com/marekm4/color-extractor"
)

const MAREKM4 string = "marekm4"

type Marekm4Extruder struct {
	Extruder
}

func init() {
	ctx := context.Background()
	err := RegisterExtruder(ctx, MAREKM4, NewMarekm4Extruder)
	if err != nil {
		panic(err)
	}
}

func NewMarekm4Colour(ctx context.Context, str_hex string) (colours.Colour, error) {

	u := url.URL{}
	u.Scheme = "common"

	q := url.Values{}
	q.Set("hex", str_hex)
	q.Set("name", MAREKM4)
	q.Set("ref", str_hex)

	u.RawQuery = q.Encode()

	return colours.NewColour(ctx, u.String())
}

func NewMarekm4Extruder(ctx context.Context, uri string) (Extruder, error) {

	ex := Marekm4Extruder{}
	return &ex, nil
}

func (ex *Marekm4Extruder) Name() string {
	return MAREKM4
}

func (ex *Marekm4Extruder) Colours(ctx context.Context, im image.Image, limit int) ([]colours.Colour, error) {

	rsp := color_extractor.ExtractColors(im)

	if len(rsp) < limit {
		limit = len(rsp)
	}

	results := make([]colours.Colour, limit)

	for i := 0; i < limit; i++ {

		c := rsp[i]

		hex_value := toHexColor(c)
		colour, err := NewMarekm4Colour(ctx, hex_value)

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
