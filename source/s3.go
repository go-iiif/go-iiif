package source

import (
	"fmt"
	iiifaws "github.com/go-iiif/go-iiif/aws"
	iiifconfig "github.com/go-iiif/go-iiif/config"
	"github.com/whosonfirst/go-whosonfirst-aws/s3"
	_ "log"
)

type S3Source struct {
	S3 *s3.S3Connection
}

func NewS3Source(cfg *iiifconfig.Config) (Source, error) {

	src := cfg.Images.Source

	bucket := src.Path
	prefix := src.Prefix
	region := src.Region
	creds := src.Credentials

	uri := fmt.Sprintf("s3://%s?region=%s&credentials=%s&prefix=%s", bucket, region, creds, prefix)
	return NewBlobSourceFromURI(uri)

	// PLEASE REMOVE EVERYTHING ELSE AS SOON AS POSSIBLE

	s3cfg := &s3.S3Config{
		Bucket:      bucket,
		Prefix:      prefix,
		Region:      region,
		Credentials: creds,
	}

	s3cfg = iiifaws.S3ConfigWrapper(s3cfg)

	s3conn, err := s3.NewS3Connection(s3cfg)

	if err != nil {
		return nil, err
	}

	c := S3Source{
		S3: s3conn,
	}

	return &c, nil
}

func (c *S3Source) Read(key string) ([]byte, error) {

	return iiifaws.S3GetWrapper(c.S3, key)
}
