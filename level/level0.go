package level

import (
	_ "fmt"
	iiifcompliance "github.com/go-iiif/go-iiif/v4/compliance"
	iiifconfig "github.com/go-iiif/go-iiif/v4/config"
	_ "log"
)

type Level0 struct {
	Level      `json:"-"`
	Formats    []string `json:"formats"`
	Qualities  []string `json:"qualities"`
	Supports   []string `json:"supports"`
	compliance iiifcompliance.Compliance
}

func NewLevel0(config *iiifconfig.Config, endpoint string) (*Level0, error) {

	compliance, err := iiifcompliance.NewLevel2Compliance(config)

	if err != nil {
		return nil, err
	}

	l := Level0{
		Formats:    compliance.Formats(),
		Qualities:  compliance.Qualities(),
		Supports:   compliance.Supports(),
		compliance: compliance,
	}

	return &l, nil
}

func (l *Level0) Compliance() iiifcompliance.Compliance {
	return l.compliance
}
