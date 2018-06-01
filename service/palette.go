package service

// https://groups.google.com/forum/#!topic/iiif-discuss/sPU5BvSWEOo
// http://palette.davidnewbury.com/

import (
	iiifpalette "github.com/thisisaaronland/go-iiif/palette"
)

type PaletteColor struct {
	Color   string `json:"colour"`
	Closest string `json:"closest"`
}

type PaletteService struct {
	Service   `json:",omitempty"`
	Context   string              `json:"@context"`
	Profile   string              `json:"profile"`
	Label     string              `json:"label"`
	Average   PaletteColor        `json:"average,omitempty"`
	Palette   []iiifpalette.Color `json:"palette,omitempty"`
	Reference string              `json:"reference-closest,omitempty"`
}

// THIS SIGNATURE WILL ALMOST CERTAINLY CHANGE (20180601/thisisaaronland)

func NewPaletteService(endpoint string, palette []iiifpalette.Color) (Service, error) {

	s := PaletteService{
		Context: "x-urn:service:palette",
		Profile: "x-urn:service:palette",
		Label:   "x-urn:service:palette",
		Palette: palette,
	}

	return &s, nil
}
