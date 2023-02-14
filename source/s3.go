package source

import (
	"fmt"

	_ "github.com/aaronland/gocloud-blob-s3"
	iiifconfig "github.com/go-iiif/go-iiif/v5/config"
)

func NewS3Source(cfg *iiifconfig.Config) (Source, error) {

	src := cfg.Images.Source

	bucket := src.Path
	prefix := src.Prefix
	region := src.Region
	creds := src.Credentials

	uri := fmt.Sprintf("s3blob://%s?region=%s&credentials=%s&prefix=%s", bucket, region, creds, prefix)
	return NewBlobSourceFromURI(uri)
}
