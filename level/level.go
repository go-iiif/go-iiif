package level

import (
	"errors"
	"fmt"
	iiifcompliance "github.com/thisisaaronland/go-iiif/compliance"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
)

type Level interface {
	Compliance() iiifcompliance.Compliance
}

func NewLevelFromConfig(config *iiifconfig.Config, host string) (Level, error) {

	compliance := config.Level.Compliance

	if compliance == "0" {

		message := fmt.Sprintf("Unsupported compliance level '%s'", compliance)
		return nil, errors.New(message)

	} else if compliance == "1" {

		message := fmt.Sprintf("Unsupported compliance level '%s'", compliance)
		return nil, errors.New(message)
	} else if compliance == "2" {

		return NewLevel2(config, host)

	} else {

		message := fmt.Sprintf("Invalid compliance level '%s'", compliance)
		return nil, errors.New(message)

	}
}
