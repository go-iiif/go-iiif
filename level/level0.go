package level

// https://iiif.io/api/image/1.1/compliance/#level0
// https://github.com/glenrobson/glenrobson.github.io/blob/master/iiif/welsh_book/page001/info.json

import (
	"fmt"
	iiifcompliance "github.com/go-iiif/go-iiif/v5/compliance"
	iiifconfig "github.com/go-iiif/go-iiif/v5/config"
	iiifprofile "github.com/go-iiif/go-iiif/v5/profile"
	iiifservice "github.com/go-iiif/go-iiif/v5/service"
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
	Tiles      []*Level0Tile `json:"tiles"`
	Sizes      []*Level0Size `json:"sizes"`
	Protocol   string        `json:"protocol"`
	Context    string        `json:"@context"`
	Id         string        `json:"@id"`
	Width      int           `json:"width"`
	Height     int           `json:"height"`
	compliance iiifcompliance.Compliance
}

func NewLevel0(config *iiifconfig.Config, endpoint string) (Level, error) {

	compliance, err := iiifcompliance.NewLevel0Compliance(config)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new level 0 compliance, %w", err)
	}

	l := Level0{
		Protocol:   "http://iiif.io/api/image",
		compliance: compliance,
	}

	return &l, nil
}

func (l *Level0) Compliance() iiifcompliance.Compliance {
	return l.compliance
}

func (l *Level0) Profile() (*iiifprofile.Profile, error) {

	p := iiifprofile.Profile{
		Context:  "http://iiif.io/api/image/2/context.json",
		Id:       "",
		Type:     "iiif:Image",
		Protocol: "http://iiif.io/api/image",
		Profile: []interface{}{
			"http://iiif.io/api/image/2/level0.json",
			l,
		},
		Services: []iiifservice.Service{},
	}

	return &p, nil
}
