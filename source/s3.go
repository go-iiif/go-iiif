package source

import (
	"fmt"
	iiifconfig "github.com/go-iiif/go-iiif/config"
	_ "log"
)

func NewS3Source(cfg *iiifconfig.Config) (Source, error) {

	src := cfg.Images.Source

	bucket := src.Path
	prefix := src.Prefix
	region := src.Region
	creds := src.Credentials

	uri := fmt.Sprintf("s3://%s?region=%s&credentials=%s&prefix=%s", bucket, region, creds, prefix)
	return NewBlobSourceFromURI(uri)
}
