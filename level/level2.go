package level

import (
	_ "fmt"
	iiifcompliance "github.com/thisisaaronland/go-iiif/compliance"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	_ "log"
)

type Level2 struct {
	Level      `json:"-"`
	Formats    []string                  `json:"formats"`
	Qualities  []string                  `json:"qualities"`
	Supports   []string                  `json:"supports"`
	compliance iiifcompliance.Compliance `json:"-"`
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
