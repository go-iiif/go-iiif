package source

import (
	_ "github.com/aaronland/gocloud-blob/s3"
	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
)

func NewS3Source(cfg *iiifconfig.Config) (Source, error) {

	return NewBlobSourceFromURI(cfg.Images.Source.URI)
}

func NewS3SourceFromURI(uri string) (Source, error) {
	return NewBlobSourceFromURI(uri)
}
