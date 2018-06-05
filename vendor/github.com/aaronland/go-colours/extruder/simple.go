package extruder

// this is basically a clone of go-colorweave

import (
	"github.com/aaronland/go-colours"
	"github.com/lucasb-eyer/go-colorful"
	"image"
	_ "log"
	"sort"
	"sync"
	_ "time"
)

type SimpleExtruder struct {
	colours.Extruder
}

func NewSimpleExtruder(args ...interface{}) (colours.Extruder, error) {

	ex := SimpleExtruder{}

	return &ex, nil
}

func (ex *SimpleExtruder) Colours(im image.Image, limit int) ([]colours.Colour, error) {

	im, err := PrepareImage(im)

	if err != nil {
		return nil, err
	}

	bounds := im.Bounds()

	// pixels := bounds.Max.X * bounds.Max.Y

	mu := new(sync.Mutex)

	lookup := make(map[string]int)

	for i := 0; i <= bounds.Max.X; i++ {

		for j := 0; j <= bounds.Max.Y; j++ {

			pixel := im.At(i, j)
			red, green, blue, _ := pixel.RGBA()

			c := colorful.Color{
				float64(red) / 255.0,
				float64(green) / 255.0,
				float64(blue) / 255.0,
			}

			h := c.Hex()

			mu.Lock()

			count, ok := lookup[h]

			if ok {
				count += 1
			} else {
				count = 1
			}

			lookup[h] = count
			mu.Unlock()
		}
	}

	reverse_lookup := reverse_map(lookup)

	keys := make([]int, 0)

	for _, count := range lookup {
		keys = append(keys, count)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(keys)))

	results := make([]colours.Colour, 0)

	for _, count := range keys {

		for _, hex_value := range reverse_lookup[count] {

			// pct := (float64(count) / float64(pixels)) * 100.0
			// c, _ := colorful.Hex(hex_value)

			colour, err := colours.NewColour(hex_value)

			if err != nil {
				return nil, err
			}

			results = append(results, colour)

			if limit > 0 && len(results) >= limit {
				break
			}
		}

		if limit > 0 && len(results) >= limit {
			break
		}
	}

	return results, nil
}

func reverse_map(hex_map map[string]int) map[int][]string {

	count_map := make(map[int][]string)

	for hex_colour, count := range hex_map {

		colours, ok := count_map[count]

		if !ok {
			colours = make([]string, 0)
		}

		colours = append(colours, hex_colour)
		count_map[count] = colours
	}

	return count_map
}
