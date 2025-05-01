package service

// https://groups.google.com/forum/#!topic/iiif-discuss/sPU5BvSWEOo
// http://palette.davidnewbury.com/

import (
	"context"

	"github.com/aaronland/go-colours"
	"github.com/aaronland/go-colours/extruder"
	"github.com/aaronland/go-colours/grid"
	"github.com/aaronland/go-colours/palette"
	iiifconfig "github.com/go-iiif/go-iiif/v8/config"
	iiifimage "github.com/go-iiif/go-iiif/v8/image"
)

func init() {

	ctx := context.Background()
	err := RegisterService(ctx, "palette", initColourService)

	if err != nil {
		panic(err)
	}
}

func initColourService(ctx context.Context, cfg *iiifconfig.Config, im iiifimage.Image) (Service, error) {
	return NewColourService(cfg.ColourServiceConfig, im)
}

type ColourService struct {
	Service        `json:",omitempty"`
	PaletteContext string           `json:"@context"`
	PaletteProfile string           `json:"profile"`
	PaletteLabel   string           `json:"label"`
	Palette        []colours.Colour `json:"palette,omitempty"`
}

func (s *ColourService) Context() string {
	return s.PaletteContext
}

func (s *ColourService) Profile() string {
	return s.PaletteProfile
}

func (s *ColourService) Label() string {
	return s.PaletteLabel
}

func (s *ColourService) Value() interface{} {
	return s.Palette
}

func NewColourService(cfg iiifconfig.ColourServiceConfig, image iiifimage.Image) (Service, error) {

	ctx := context.Background()

	extruder_uri := cfg.Extruder.URI
	extruder_count := cfg.Extruder.Count

	grid_uri := cfg.Grid.URI
	palette_uris := make([]string, 0)

	for _, p := range cfg.Palettes {
		palette_uris = append(palette_uris, p.URI)
	}

	im, err := iiifimage.IIIFImageToGolangImage(image)

	if err != nil {
		return nil, err
	}

	ex, err := extruder.NewExtruder(ctx, extruder_uri)

	if err != nil {
		return nil, err
	}

	gr, err := grid.NewGrid(ctx, grid_uri)

	if err != nil {
		return nil, err
	}

	plts := make([]palette.Palette, len(palette_uris))

	for i, p := range palette_uris {

		pl, err := palette.NewPalette(ctx, p)

		if err != nil {
			return nil, err
		}

		plts[i] = pl
	}

	has_colours, err := ex.Colours(im, extruder_count)

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

	s := ColourService{
		PaletteContext: "x-urn:service:go-iiif#palette",
		PaletteProfile: "x-urn:service:go-iiif#palette",
		PaletteLabel:   "x-urn:service:go-iiif#palette",
		Palette:        has_colours,
	}

	return &s, nil
}
