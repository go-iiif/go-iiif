package level

import (
	"fmt"

	iiifcompliance "github.com/go-iiif/go-iiif/v7/compliance"
	iiifconfig "github.com/go-iiif/go-iiif/v7/config"
)

type Level interface {
	Compliance() iiifcompliance.Compliance
	Endpoint() string
	Profile() string
}

func NewLevelFromConfig(config *iiifconfig.Config, endpoint string) (Level, error) {

	compliance := config.Level.Compliance

	switch compliance {
	case "0":
		return NewLevel0(config, endpoint)
	case "2":
		return NewLevel2(config, endpoint)
	default:
		return nil, fmt.Errorf("Invalid compliance level '%s'", compliance)
	}
}
