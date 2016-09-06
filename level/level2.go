package level

import (
	"fmt"
	iiifcompliance "github.com/thisisaaronland/go-iiif/compliance"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	"log"
)

type Level2 struct {
	Level      `json:"-"`
	Context    string                    `json:@profile`
	Id         string                    `json:"@id"`
	Type       string                    `json:"@type"` // Optional or iiif:Image
	Formats    []string                  `json:"formats"`
	Qualities  []string                  `json:"qualities"`
	Supports   []string                  `json:"supports"`
	compliance iiifcompliance.Compliance `json:"-"`
}

func NewLevel2(config *iiifconfig.Config, host string) (*Level2, error) {

	compliance, err := iiifcompliance.NewLevel2Compliance(config)

	if err != nil {
		log.Fatal(err)
	}

	id := fmt.Sprintf("http://%s/level2.json", host)

	l := Level2{
		Context:    "http://iiif.io/api/image/2/context.json",
		Id:         id,
		Type:       "iiif:ImageProfile",
		Formats:    compliance.Formats(),
		Qualities:  compliance.Qualities(),
		Supports:   []string{},
		compliance: compliance,
	}

	return &l, nil
}

func (l *Level2) Compliance() iiifcompliance.Compliance {
	return l.compliance
}