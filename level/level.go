package level

import (
	"errors"
	"fmt"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
)

type Level interface {
	IsValidImageRegion(string) (bool, error)
	IsValidImageSize(string) (bool, error)
	IsValidImageRotation(string) (bool, error)
	IsValidImageQuality(string) (bool, error)
	IsValidImageFormat(string) (bool, error)
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
