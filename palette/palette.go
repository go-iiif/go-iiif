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