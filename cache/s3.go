package cache

import (
	"fmt"
	iiifconfig "github.com/go-iiif/go-iiif/config"
	_ "log"
)

func NewS3Cache(cfg iiifconfig.CacheConfig) (Cache, error) {

	bucket := cfg.Path
	prefix := cfg.Prefix
	region := cfg.Region
	creds := cfg.Credentials

	uri := fmt.Sprintf("s3://%s?region=%s&credentials=%s&prefix=%s", bucket, region, creds, prefix)
	return NewBlobCacheFromURI(uri)
}
