package source

import (
	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
)

func NewDiskSource(cfg *iiifconfig.Config) (Source, error) {

	return NewBlobSourceFromURI(cfg.Images.Source.URI)
}

func NewDiskSourceFromURI(uri string) (Source, error) {
	return NewBlobSourceFromURI(uri)
}
