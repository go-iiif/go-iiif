package swatchbook

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aaronland/go-swatchbook/palettes"
	"github.com/pwaller/go-hexcolor"
	// "image/color"
	_ "log"
	"math"
	"sort"
	"strings"
)

// I tried making this be all interface{} - y but that turned in to
// yak-shaving debugging Go/JSON unmarshaling (20180604/thisisaaronland)

type Color struct {
	Name string `json:"name"`
	Hex  string `json:"hex"`
}

type Palette struct {
	Source string  `json:"source"`
	Colors []*Color `json:"colors"`
}

type Swatchbook struct {
	Palette *Palette `json:"palette"`
}

func (c *Color) String() string {
     return fmt.Sprintf("%s (%s)", c.Hex, c.Name)
}

func (s *Swatchbook) Closest(target *Color) *Color {

     	// https://github.com/pwaller/go-hexcolor/blob/master/hexcolor.go
        // https://github.com/ubernostrum/webcolors/blob/master/webcolors.py#L473-L485

	r , g, b, _ := hexcolor.HexToRGBA(hexcolor.Hex(target.Hex))

	lookup := make(map[int]*Color)
	keys := make([]int, 0)

	p := s.Palette

	for _, c := range p.Colors {

		 rc, gc, bc, _ := hexcolor.HexToRGBA(hexcolor.Hex(c.Hex))

		rd := math.Pow(float64(int32(rc)-int32(r)), 2.0)
		gd := math.Pow(float64(int32(gc)-int32(g)), 2.0)
		bd := math.Pow(float64(int32(bc)-int32(b)), 2.0)

		k := int(rd + gd + bd)
		lookup[k] = c

		keys = append(keys, k)
	}

	sort.Ints(keys)

	/*
	for i, idx := range keys {
		log.Println(i, idx, lookup[idx])
	}
	*/

	return lookup[keys[0]]
}

func NewNamedPalette(name string) (*Palette, error) {

	var data []byte
	var err error

	switch strings.ToUpper(name) {
	case "CSS4":
		data = palettes.CSS4
	default:
		err = errors.New("Invalid or unknown palette")
	}

	if err != nil {
		return nil, err
	}

	var p Palette

	err = json.Unmarshal(data, &p)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func NewSwatchbookFromPalette(p *Palette) (*Swatchbook, error) {

	sb := Swatchbook{
		Palette: p,
	}

	return &sb, nil
}
