package grid

import (
	"context"
	"fmt"
	"math"
	"net/url"
	"sort"

	"github.com/aaronland/go-colours"
	"github.com/aaronland/go-colours/palette"
	"github.com/lucasb-eyer/go-colorful"
)

type EuclidianGrid struct {
	Grid
}

func init() {
	ctx := context.Background()
	err := RegisterGrid(ctx, "euclidian", NewEuclidianGrid)
	if err != nil {
		panic(err)
	}
}

func NewEuclidianGrid(ctx context.Context, uri string) (Grid, error) {
	eu := EuclidianGrid{}
	return &eu, nil
}

func (eu *EuclidianGrid) Closest(ctx context.Context, target colours.Colour, plt palette.Palette) (colours.Colour, error) {

	// http://stackoverflow.com/questions/9694165/convert-rgb-color-to-english-color-name-like-green
	// https://github.com/ubernostrum/webcolors/blob/master/webcolors.py#L473-L485

	cl, err := colorful.Hex(target.Hex())

	if err != nil {
		return nil, fmt.Errorf("Failed to derive hex, %w", err)
	}

	r1, g1, b1 := cl.RGB255()

	lookup := make(map[int]colours.Colour)
	keys := make([]int, 0)

	for _, candidate := range plt.Colours() {

		cl, err := colorful.Hex(candidate.Hex())

		if err != nil {
			return nil, err
		}

		r2, g2, b2 := cl.RGB255()

		r := math.Pow(float64(int32(r2)-int32(r1)), 2.0)
		g := math.Pow(float64(int32(g2)-int32(g1)), 2.0)
		b := math.Pow(float64(int32(b2)-int32(b1)), 2.0)

		k := int(r + g + b)
		lookup[k] = candidate

		keys = append(keys, k)
	}

	sort.Ints(keys)

	if len(keys) == 0 {
		return nil, fmt.Errorf("Nothing found")
	}

	match := lookup[keys[0]]

	q := url.Values{}
	q.Set("hex", match.Hex())
	q.Set("name", match.Name())
	q.Set("ref", plt.Reference())

	u := url.URL{}
	u.Scheme = "common"
	u.RawQuery = q.Encode()

	return colours.NewColour(ctx, u.String())
}
