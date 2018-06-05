package service

// https://groups.google.com/forum/#!topic/iiif-discuss/sPU5BvSWEOo
// http://palette.davidnewbury.com/

import (
       "github.com/aaronland/go-colours"
       "github.com/aaronland/go-colours/extruder"
       "github.com/aaronland/go-colours/grid"
       "github.com/aaronland/go-colours/palette"
	"image"
	_ "log"
)

type PaletteService struct {
	Service   `json:",omitempty"`
	Context   string              `json:"@context"`
	Profile   string              `json:"profile"`
	Label     string              `json:"label"`
	Palette   []colours.Colour `json:"palette,omitempty"`
}

// THIS SIGNATURE WILL ALMOST CERTAINLY CHANGE (20180601/thisisaaronland)

func NewPaletteService(im image.Image) (Service, error) {

     	// please read me from the config...
	// (20180605/thisisaaronland)

     	use_extruder := "vibrant"
	use_grid := "euclidian"
	use_palette := "css4"
	count_colours := 5

	ex, err := extruder.NewNamedExtruder(use_extruder)

	if err != nil {
		return nil, err
	}
     
	gr, err := grid.NewNamedGrid(use_grid)

	if err != nil {
		return nil, err
	}

	pl, err := palette.NewNamedPalette(use_palette)

	if err != nil {
		return nil, err
	}

	plts := []colours.Palette{
		pl,
	}

	has_colours, err := ex.Colours(im, count_colours)

	if err != nil {
		return nil, err
	}

	for _, c := range has_colours {

		for _, p := range plts {

			cl, err := gr.Closest(c, p)

			if err != nil {
				return nil, err
			}

			err = c.AppendClosest(cl)

			if err != nil {
				return nil, err
			}
		}

	}

	s := PaletteService{
		Context: "x-urn:service:palette",
		Profile: "x-urn:service:palette",
		Label:   "x-urn:service:palette",
		Palette: has_colours,
	}

	return &s, nil
}
