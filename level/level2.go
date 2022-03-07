package level

import (
	_ "fmt"
	iiifcompliance "github.com/go-iiif/go-iiif/v5/compliance"
	iiifconfig "github.com/go-iiif/go-iiif/v5/config"
	iiifprofile "github.com/go-iiif/go-iiif/v5/profile"
	iiifservice "github.com/go-iiif/go-iiif/v5/service"
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

func (l *Level2) Profile() (*iiifprofile.Profile, error) {

	p := iiifprofile.Profile{
		Context:  "http://iiif.io/api/image/2/context.json",
		Id:       "",
		Type:     "iiif:Image",
		Protocol: "http://iiif.io/api/image",
		Profile: []interface{}{
			"http://iiif.io/api/image/2/level2.json",
			l,
		},
		Services: []iiifservice.Service{},
	}

	return &p, nil
}
