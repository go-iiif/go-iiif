package level

// https://iiif.io/api/image/1.1/compliance/#level0
// https://github.com/glenrobson/glenrobson.github.io/blob/master/iiif/welsh_book/page001/info.json

import (
	"fmt"

	iiifcompliance "github.com/go-iiif/go-iiif/v7/compliance"
	iiifconfig "github.com/go-iiif/go-iiif/v7/config"
)

type Level0 struct {
	Level      `json:"-"`
	Formats    []string `json:"formats"`
	Qualities  []string `json:"qualities"`
	Supports   []string `json:"supports"`
	compliance iiifcompliance.Compliance
	endpoint   string
}

func NewLevel0(config *iiifconfig.Config, endpoint string) (Level, error) {

	compliance, err := iiifcompliance.NewLevel0Compliance(config)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new level 0 compliance, %w", err)
	}

	l := Level0{
		Formats:    compliance.Formats(),
		Qualities:  compliance.Qualities(),
		Supports:   compliance.Supports(),
		compliance: compliance,
		endpoint:   endpoint,
	}

	return &l, nil
}

func (l *Level0) Compliance() iiifcompliance.Compliance {
	return l.compliance
}

func (l *Level0) Endpoint() string {
	return l.endpoint
}

func (l *Level0) Profile() string {
	return "http://iiif.io/api/image/2/level0.json"
}
