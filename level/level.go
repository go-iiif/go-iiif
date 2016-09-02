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

func NewLevelFromConfig(config iiifconfig.IIIFConfig, host string) (Level, error) {

	version := config.Level

	if version == "0" {

		message := fmt.Sprintf("Unsupported level '%s'", version)
		return nil, errors.New(message)

	} else if version == "1" {

		message := fmt.Sprintf("Unsupported level '%s'", version)

		return nil, errors.New(message)
	} else if version == "2" {

		return NewLevel2(config, host)

	} else {

		message := fmt.Sprintf("Invalid level '%s'", version)
		return nil, errors.New(message)

	}
}
