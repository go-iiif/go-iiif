package source

import (
	"errors"
	"github.com/thisisaaronland/go-iiif/config"
)

type Source interface {
	Read(uri string) ([]byte, error)
}

func NewSourceFromConfig(cfg config.ImagesConfig) (Source, error) {

	if cfg.Source.Name == "Disk" {
		cache, err := NewDiskSource(cfg)
		return cache, err
	} else {
		err := errors.New("Unknown source type")
		return nil, err
	}
}
