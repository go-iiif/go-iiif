package extruder

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"strings"
	_ "log/slog"
	
	"github.com/RobCherry/vibrant"
	"github.com/aaronland/go-colours"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/draw"
)

// Important: IsTransparentFilter relies on the sfomuseum/vibrant fork
// which exposes Filter.IsAllowed as a public method
// github.com/sfomuseum/vibrant v0.0.0-20250430212339-abb21560aa26

type IsTransparentFilter struct {
	vibrant.Filter
}

func (f *IsTransparentFilter) IsAllowed(c color.Color) bool {
	_, _, _, a := c.RGBA()
	return a > 0.0
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
	return "vibrant"
}

func (v *VibrantExtruder) Colours(im image.Image, limit int) ([]colours.Colour, error) {

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
		hex = strings.TrimLeft(hex, "#")

		ctx := context.Background()

		c_uri := fmt.Sprintf("common://?hex=%s&name=%s&ref=vibrant", hex, hex)
		c, err := colours.NewColour(ctx, c_uri)

		if err != nil {
			return nil, fmt.Errorf("Failed to create new color '%s', %w", c_uri, err)
		}

		results = append(results, c)

		if limit > 0 && len(results) == limit {
			break
		}
	}

	return results, nil
}
