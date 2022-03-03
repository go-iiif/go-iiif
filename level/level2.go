package level

import (
	_ "fmt"
	iiifcompliance "github.com/go-iiif/go-iiif/v5/compliance"
	iiifconfig "github.com/go-iiif/go-iiif/v5/config"
	_ "log"
)

type Level2 struct {
	Level      `json:"-"`
	Formats    []string `json:"formats"`
	Qualities  []string `json:"qualities"`
	Supports   []string `json:"supports"`
	compliance iiifcompliance.Compliance
}

func NewLevel2(config *iiifconfig.Config, endpoint string) (*Level2, error) {

	compliance, err := iiifcompliance.NewLevel2Compliance(config)

	if err != nil {
		return nil, err
	}

	l := Level2{
		Formats:    compliance.Formats(),
		Qualities:  compliance.Qualities(),
		Supports:   compliance.Supports(),
		compliance: compliance,
	}

	return &l, nil
}

func (l *Level2) Compliance() iiifcompliance.Compliance {
	return l.compliance
}

func (l *Level2) Profile(endpoint string, image iiifimage.Image) iiifprofile.Profile {

	dims, err := image.Dimensions()

	if err != nil {
		return nil, fmt.Errorf("Failed to derive dimensions for image, %w", err)
	}

	p := iiifprofile.Profile{
		Context:  "http://iiif.io/api/image/2/context.json",
		Id:       fmt.Sprintf("%s/%s", endpoint, image.Identifier()),
		Type:     "iiif:Image",
		Protocol: "http://iiif.io/api/image",
		Width:    dims.Width(),
		Height:   dims.Height(),
		Profile: []interface{}{
			"http://iiif.io/api/image/2/level2.json",
			level,
		},
		Services: []iiifservice.Service{},
	}

	return &p, nil
}
