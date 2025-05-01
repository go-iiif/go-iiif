package level

import (
	iiifcompliance "github.com/go-iiif/go-iiif/v8/compliance"
	iiifconfig "github.com/go-iiif/go-iiif/v8/config"
)

type Level2 struct {
	Level      `json:"-"`
	Formats    []string `json:"formats"`
	Qualities  []string `json:"qualities"`
	Supports   []string `json:"supports"`
	compliance iiifcompliance.Compliance
	endpoint   string
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
		endpoint:   endpoint,
	}

	return &l, nil
}

func (l *Level2) Endpoint() string {
	return l.endpoint
}

func (l *Level2) Compliance() iiifcompliance.Compliance {
	return l.compliance
}

func (l *Level2) Profile() string {
	return "http://iiif.io/api/image/2/level2.json"
}
