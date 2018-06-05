package grid

import (
	"errors"
	"github.com/aaronland/go-colours"
	"github.com/lucasb-eyer/go-colorful"
	_ "log"
	"math"
	"sort"
)

type EuclidianGrid struct {
	colours.Grid
}

func NewEuclidianGrid(args ...interface{}) (colours.Grid, error) {
	eu := EuclidianGrid{}
	return &eu, nil
}

func (eu *EuclidianGrid) Closest(target colours.Colour, palette colours.Palette) (colours.Colour, error) {

	// https://github.com/ubernostrum/webcolors/blob/master/webcolors.py#L473-L485

	cl, err := colorful.Hex(target.Hex())

	if err != nil {
		return nil, err
	}

	r1, g1, b1 := cl.RGB255()

	lookup := make(map[int]colours.Colour)
	keys := make([]int, 0)

	for _, candidate := range palette.Colours() {

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
		return nil, errors.New("Nothing found")
	}

	match := lookup[keys[0]]

	return colours.NewColour(match.Hex(), match.Name(), palette.Reference())
}
