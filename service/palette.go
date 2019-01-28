package service

// https://groups.google.com/forum/#!topic/iiif-discuss/sPU5BvSWEOo
// http://palette.davidnewbury.com/

import (
	"github.com/aaronland/go-colours"
	"github.com/aaronland/go-colours/extruder"
	"github.com/aaronland/go-colours/grid"
	"github.com/aaronland/go-colours/palette"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	iiifimage "github.com/thisisaaronland/go-iiif/image"
	_ "log"
)

type PaletteService struct {
	Service        `json:",omitempty"`
	PaletteContext string           `json:"@context"`
	PaletteProfile string           `json:"profile"`
	PaletteLabel   string           `json:"label"`
	Palette        []colours.Colour `json:"palette,omitempty"`
}

func (s *PaletteService) Context() string {
	return s.PaletteContext
}

func (s *PaletteService) Profile() string {
	return s.PaletteProfile
}

func (s *PaletteService) Label() string {
	return s.PaletteLabel
}

func (s *PaletteService) Value() interface{} {
	return s.Palette
}

func NewPaletteService(cfg iiifconfig.PaletteConfig, image iiifimage.Image) (Service, error) {

	// b, _ := json.Marshal(cfg)
	// log.Println(string(b))

	use_extruder := cfg.Extruder.Name
	count_colours := cfg.Extruder.Count

	use_grid := cfg.Grid.Name
	use_palette := make([]string, 0)

	for _, p := range cfg.Palettes {
		use_palette = append(use_palette, p.Name)
	}

	im, err := iiifimage.IIIFImageToGolangImage(image)

	if err != nil {
		return nil, err
	}

	ex, err := extruder.NewNamedExtruder(use_extruder)

	if err != nil {
		return nil, err
	}

	gr, err := grid.NewNamedGrid(use_grid)

	if err != nil {
		return nil, err
	}

	plts := make([]colours.Palette, len(use_palette))

	for i, p := range use_palette {

		pl, err := palette.NewNamedPalette(p)

		if err != nil {
			return nil, err
		}

		plts[i] = pl
	}

	has_colours, err := ex.Colours(im, count_colours)

	if err != nil {
		return nil, err
	}

	for _, c := range has_colours {

		for _, pl := range plts {

			cl, err := gr.Closest(c, pl)

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
		PaletteContext: "x-urn:service:go-iiif#palette",
		PaletteProfile: "x-urn:service:go-iiif#palette",
		PaletteLabel:   "x-urn:service:go-iiif#palette",
		Palette:        has_colours,
	}

	return &s, nil
}
