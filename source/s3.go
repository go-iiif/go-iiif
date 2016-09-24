package source

import (
	iiifaws "github.com/thisisaaronland/go-iiif/aws"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	_ "log"
)

type S3Source struct {
	S3 *iiifaws.S3Connection
}

func NewS3Source(cfg *iiifconfig.Config) (*S3Source, error) {

	src := cfg.Images.Source

	bucket := src.Path
	prefix := ""

	region := "us-east-1"
	creds := "default"

	if src.Prefix == "" {
		prefix = src.Prefix
	}

	if src.Region == "" {
		region = src.Region
	}

	if src.Credentials == "" {
		creds = src.Credentials
	}

	s3cfg := iiifaws.S3Config{
		Bucket:      bucket,
		Prefix:      prefix,
		Region:      region,
		Credentials: creds,
	}

	s3, err := iiifaws.NewS3Connection(s3cfg)

	if err != nil {
		return nil, err
	}

	c := S3Source{
		S3: s3,
	}

	return &c, nil
}

func (c *S3Source) Read(id string) ([]byte, error) {

	return c.S3.Get(id)
}
