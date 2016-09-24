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

	// note that there is no "Memory" source or at least not yet
	// since it assumes you're passing it []bytes and not a config
	// file (20160907/thisisaaronland)

	if cfg.Source.Name == "Disk" {
		cache, err := NewDiskSource(config)
		return cache, err
	} else if cfg.Source.Name == "Flickr" {
		cache, err := NewFlickrSource(config)
		return cache, err
	} else if cfg.Source.Name == "S3" {
		cache, err := NewS3Source(config)
		return cache, err
	} else if cfg.Source.Name == "URI" {
		cache, err := NewURISource(config)
		return cache, err
	} else {
		err := errors.New("Unknown source type")
		return nil, err
	}
}
