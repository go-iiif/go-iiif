package source

import (
	"fmt"
	_ "log"

	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
)

func NewDiskSourceURIFromConfig(config *iiifconfig.Config) (string, error) {

	uri := cfg.Source.URI

	if uri == "" {
		cfg := config.Images
		uri = fmt.Sprintf("file://%s", cfg.Source.Path)
	}

	return uri, nil
}

func NewDiskSource(config *iiifconfig.Config) (Source, error) {

	uri, err := NewDiskSourceURIFromConfig(config)

	if err != nil {
		return nil, err
	}

	return NewBlobSourceFromURI(uri)
}

func NewDiskSourceFromURI(uri string) (Source, error) {
	return NewBlobSourceFromURI(uri)
}
