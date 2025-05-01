package service

// https://groups.google.com/forum/#!topic/iiif-discuss/sPU5BvSWEOo
// http://palette.davidnewbury.com/

import (
	"context"
	"fmt"

	"github.com/aaronland/go-colours"
	"github.com/aaronland/go-colours/extruder"
	"github.com/aaronland/go-colours/grid"
	"github.com/aaronland/go-colours/palette"
	iiifconfig "github.com/go-iiif/go-iiif/v8/config"
	iiifimage "github.com/go-iiif/go-iiif/v8/image"
)

func init() {

	ctx := context.Background()
	err := RegisterService(ctx, "palette", initPaletteService)

	if err != nil {
		panic(err)
	}
}

func initPaletteService(ctx context.Context, cfg *iiifconfig.Config, im iiifimage.Image) (Service, error) {
	return NewPaletteService(cfg.PaletteServiceConfig, im)
}

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

func NewPaletteService(cfg iiifconfig.PaletteServiceConfig, image iiifimage.Image) (Service, error) {

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
		return nil, fmt.Errorf("Failed to derive image, %w", err)
	}

	ex, err := extruder.NewExtruder(ctx, extruder_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new extruder (%s), %w", extruder_uri, err)
	}

	gr, err := grid.NewGrid(ctx, grid_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new grid (%s), %w", grid_uri, err)
	}

	palettes := make([]palette.Palette, len(palette_uris))

	for i, uri := range palette_uris {

		p, err := palette.NewPalette(ctx, uri)

		if err != nil {
			return nil, fmt.Errorf("Failed to create new palette (%s), %w", uri, err)
		}

		palettes[i] = p
	}

	has_colours, err := ex.Colours(ctx, im, extruder_count)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive colours, %w", err)
	}

	for _, c := range has_colours {

		for _, pl := range palettes {

			cl, err := gr.Closest(ctx, c, pl)

			if err != nil {
				return nil, fmt.Errorf("Failed to derive closest match for '%s' from '%s', %w", c, pl.Reference(), err)
			}

			err = c.AppendClosest(cl)

			if err != nil {
				return nil, fmt.Errorf("Failed to append closest colour to '%s', %w", c, err)
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
