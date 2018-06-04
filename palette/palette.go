package palette

// MOST OF THIS CODE WILL BE MOVED IN TO ANOTHER NON-IIIF PACKAGE SHORTLY
// AND WHAT IS LEFT WILL BE THE USUAL MakePaletteFromConfig TYPE OF THING
// RIGHT NOW I AM JUST TRYING TO MAKE IT ALL WORK (20180601/thisisaaronland)

import (
	"image"
)

type Palette interface {
	Extract(im image.Image, limit int) ([]Color, error)
}

type Color struct {
	Color   string `json:"color"`
	Closest string `json:"closest,omitempty"`
}

func Closest(hex_color string, possible []PColor) PColor {

     r, g,b, _ := hexcolor.HexToRGBA(hex)

     lookup := make(map[int]PColor)
     keys := make([]int, 0)

     for _, c := range possible {

     	 rc, gc, bc, _ := hexcolor.HexToRGBA(c.Hex)

	 rd := (rc - r) ** 2
	 gd := (gc - g) ** 2
	 bd := (bc - b) ** 2

	 k := rd + gd + bd
	 lookup[k] = c

	 keys = append(keys, k)
     }

     sort.Ints(keys)
     return lookup[keys[0]]
}