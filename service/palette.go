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

	im, err := iiifimage.IIIFImageToGolangImage(image)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive image, %w", err)
	}

	extruders := make([]extruder.Extruder, len(cfg.Extruders))
	palettes := make([]palette.Palette, len(cfg.Palettes))

	extruder_counts := make(map[string]int)

	for idx, e_cfg := range cfg.Extruders {

		e, err := extruder.NewExtruder(ctx, e_cfg.URI)

		if err != nil {
			return nil, fmt.Errorf("Failed to create new extruder (%s), %w", e_cfg.URI, err)
		}

		extruder_counts[e.Name()] = e_cfg.Count
		extruders[idx] = e
	}

	for idx, p_cfg := range cfg.Palettes {

		p, err := palette.NewPalette(ctx, p_cfg.URI)

		if err != nil {
			return nil, fmt.Errorf("Failed to create new palette (%s), %w", p_cfg.URI, err)
		}

		palettes[idx] = p
	}

	gr, err := grid.NewGrid(ctx, cfg.Grid.URI)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new grid (%s), %w", cfg.Grid.URI, err)
	}

	all_colours := make([]colours.Colour, 0)

	for _, ex := range extruders {

		max_colours := extruder_counts[ex.Name()]
		has_colours, err := ex.Colours(ctx, im, max_colours)

		if err != nil {
			return nil, fmt.Errorf("Failed to derive colours for extruder %s, %w", ex.Name(), err)
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

			all_colours = append(all_colours, c)
		}
	}

	s := PaletteService{
		PaletteContext: "x-urn:service:go-iiif#palette",
		PaletteProfile: "x-urn:service:go-iiif#palette",
		PaletteLabel:   "x-urn:service:go-iiif#palette",
		Palette:        all_colours,
	}

	return &s, nil
}
