package cache

import (
	iiifaws "github.com/thisisaaronland/go-iiif/aws"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
)

type S3Cache struct {
	S3 *iiifaws.S3Connection
}

func NewS3Cache(cfg iiifconfig.CacheConfig) (*S3Cache, error) {

	bucket := cfg.Path
	prefix := cfg.Prefix
	region := cfg.Region
	creds := cfg.Credentials

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

	c := S3Cache{
		S3: s3,
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

	return c.S3.Get(key)
}

func (c *S3Cache) Set(key string, body []byte) error {

	return c.S3.Put(key, body)
}

func (c *S3Cache) Unset(key string) error {

	return c.S3.Delete(key)
}
