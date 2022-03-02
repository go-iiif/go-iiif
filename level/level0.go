package level

// https://iiif.io/api/image/1.1/compliance/#level0
// https://github.com/glenrobson/glenrobson.github.io/blob/master/iiif/welsh_book/page001/info.json

import (
	"fmt"
	iiifcompliance "github.com/go-iiif/go-iiif/v4/compliance"
	iiifconfig "github.com/go-iiif/go-iiif/v4/config"
	_ "log"
)

type Level0Profile struct {
	Formats   []string `json:"formats"`
	Qualities []string `json:"qualities"`
}

type Level0Tile struct {
	Width        int   `json:"width"`
	Height       int   `json:"height"`
	ScaleFactors []int `json:"scaleFactors"`
}

type Level0Size struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// compliance iiifcompliance.Compliance

type Level0 struct {
	Level      `json:"-"`
	Profile    []*Level0Profile `json:"profile"`
	Tiles      []*Level0Tile    `json:"tiles"`
	Sizes      []*Level0Size    `json:"sizes"`
	Protocol   string           `json:"protocol"`
	Context    string           `json:"@context"`
	Id         string           `json:"@id"`
	Width      int              `json:"width"`
	Height     int              `json:"height"`
	compliance iiifcompliance.Compliance
}

func NewLevel0(config *iiifconfig.Config, endpoint string) (*Level0, error) {

	compliance, err := iiifcompliance.NewLevel0Compliance(config)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new level 0 compliance, %w", err)
	}

	p := &Level0Profile{
		Formats:   compliance.Formats(),
		Qualities: compliance.Qualities(),
	}

	l := Level0{
		Protocol:   "http://iiif.io/api/image",
		Profile:    []*Level0Profile{p},
		compliance: compliance,
	}

	return &l, nil
}

func (l *Level0) Compliance() iiifcompliance.Compliance {
	return l.compliance
}
