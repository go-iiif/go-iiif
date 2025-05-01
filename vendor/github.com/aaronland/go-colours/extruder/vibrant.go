package extruder

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"net/url"

	"github.com/aaronland/go-colours"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/sfomuseum/vibrant"
	"golang.org/x/image/draw"
)

const VIBRANT string = "vibrant"

// Important: IsTransparentFilter relies on the sfomuseum/vibrant fork
// which exposes Filter.IsAllowed as a public method
// github.com/sfomuseum/vibrant

type IsTransparentFilter struct {
	vibrant.Filter
}

func (f *IsTransparentFilter) IsAllowed(c color.Color) bool {
	_, _, _, a := c.RGBA()
	return a > 0.0
}

func NewVibrantColour(ctx context.Context, str_hex string) (colours.Colour, error) {

	u := url.URL{}
	u.Scheme = "common"

	q := url.Values{}
	q.Set("hex", str_hex)
	q.Set("name", VIBRANT)
	q.Set("ref", str_hex)

	u.RawQuery = q.Encode()

	return colours.NewColour(ctx, u.String())
}

type VibrantExtruder struct {
	Extruder
}

func init() {
	ctx := context.Background()
	err := RegisterExtruder(ctx, "vibrant", NewVibrantExtruder)
	if err != nil {
		panic(err)
	}
}

func NewVibrantExtruder(ctx context.Context, uri string) (Extruder, error) {

	v := VibrantExtruder{}
	return &v, nil
}

func (ex *VibrantExtruder) Name() string {
	return VIBRANT
}

func (v *VibrantExtruder) Colours(ctx context.Context, im image.Image, limit int) ([]colours.Colour, error) {

	results := make([]colours.Colour, 0)

	pb := vibrant.NewPaletteBuilder(im)
	pb = pb.MaximumColorCount(uint32(limit))
	pb = pb.Scaler(draw.ApproxBiLinear)

	f := new(IsTransparentFilter)
	pb = pb.AddFilter(f)

	palette := pb.Generate()

	for _, sw := range palette.Swatches() {

		if sw == nil {
			continue
		}

		cl, ok := colorful.MakeColor(sw.Color())

		if !ok {
			return nil, fmt.Errorf("Unable to make color, %v", sw.Color())
		}

		hex := cl.Hex()
		c, err := NewVibrantColour(ctx, hex)

		if err != nil {
			return nil, fmt.Errorf("Failed to create new color '%s', %w", hex, err)
		}

		results = append(results, c)

		if limit > 0 && len(results) == limit {
			break
		}
	}

	return results, nil
}
