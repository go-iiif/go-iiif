package source

import (
	"errors"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
)

type Source interface {
	Read(uri string) ([]byte, error)
}

func NewSourceFromConfig(config *iiifconfig.Config) (Source, error) {

	cfg := config.Images

	if cfg.Source.Name == "Disk" {
		cache, err := NewDiskSource(config)
		return cache, err
	} else {
		err := errors.New("Unknown source type")
		return nil, err
	}
}
