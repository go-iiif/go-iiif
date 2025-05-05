package extruder

import (
	"context"
	"image"
	"net/url"

	"github.com/aaronland/go-colours"
	"github.com/soniakeys/quant/mean"
)

const QUANT string = "quant"

type QuantExtruder struct {
	Extruder
}

func init() {
	ctx := context.Background()
	err := RegisterExtruder(ctx, QUANT, NewQuantExtruder)
	if err != nil {
		panic(err)
	}
}

func NewQuantColour(ctx context.Context, str_hex string) (colours.Colour, error) {

	u := url.URL{}
	u.Scheme = "common"

	q := url.Values{}
	q.Set("hex", str_hex)
	q.Set("name", QUANT)
	q.Set("ref", str_hex)

	u.RawQuery = q.Encode()

	return colours.NewColour(ctx, u.String())
}

func NewQuantExtruder(ctx context.Context, uri string) (Extruder, error) {

	ex := QuantExtruder{}
	return &ex, nil
}

func (ex *QuantExtruder) Name() string {
	return QUANT
}

func (ex *QuantExtruder) Colours(ctx context.Context, im image.Image, limit int) ([]colours.Colour, error) {

	sz := limit
	q := mean.Quantizer(sz)

	pal := q.Palette(im)

	results := make([]colours.Colour, 0)

	for _, c := range pal.ColorPalette() {

		hex_value := toHexColor(c)

		colour, err := NewQuantColour(ctx, hex_value)

		if err != nil {
			return nil, err
		}

		results = append(results, colour)

		if len(results) == limit {
			break
		}
	}

	return results, nil
}
