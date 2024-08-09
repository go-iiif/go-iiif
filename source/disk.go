package source

import (
	"fmt"
	_ "log"

	iiifconfig "github.com/go-iiif/go-iiif/v7/config"
)

func NewDiskSourceURIFromConfig(cfg *iiifconfig.Config) (string, error) {

	uri := cfg.Images.Source.URI

	if uri == "" {
		uri = fmt.Sprintf("file://%s", cfg.Images.Source.Path)
	}

	return uri, nil
}

func NewDiskSource(cfg *iiifconfig.Config) (Source, error) {

	uri, err := NewDiskSourceURIFromConfig(cfg)

	if err != nil {
		return nil, err
	}

	return NewBlobSourceFromURI(uri)
}

func NewDiskSourceFromURI(uri string) (Source, error) {
	return NewBlobSourceFromURI(uri)
}
