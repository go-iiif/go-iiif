package extruder

import (
	"context"
	"image"
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
	q.Set("name", str_hex)
	q.Set("ref", MAREKM4)

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
