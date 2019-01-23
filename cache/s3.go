package cache

import (
	iiifaws "github.com/thisisaaronland/go-iiif/aws"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	"github.com/whosonfirst/go-whosonfirst-aws/s3"
)

type S3Cache struct {
	S3 *s3.S3Connection
}

func NewS3Cache(cfg iiifconfig.CacheConfig) (*S3Cache, error) {

	bucket := cfg.Path
	prefix := cfg.Prefix
	region := cfg.Region
	creds := cfg.Credentials

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

	c := S3Cache{
		S3: s3conn,
	}

	return &c, nil
}

func (c *S3Cache) Exists(key string) bool {

	_, err := c.S3.Head(key)

	if err != nil {
		return false
	}

	return true
}

func (c *S3Cache) Get(key string) ([]byte, error) {

	return iiifaws.S3GetWrapper(c.S3, key)
}

func (c *S3Cache) Set(key string, body []byte) error {

	return iiifaws.S3SetWrapper(c.S3, key, body)
}

func (c *S3Cache) Unset(key string) error {

	return c.S3.Delete(key)
}
