package cache

import (
	"fmt"
	_ "github.com/aaronland/go-cloud-s3blob"
	iiifconfig "github.com/go-iiif/go-iiif/v2/config"
	_ "log"
)

func NewS3Cache(cfg iiifconfig.CacheConfig) (Cache, error) {

	bucket := cfg.Path
	prefix := cfg.Prefix
	region := cfg.Region
	creds := cfg.Credentials

	uri := fmt.Sprintf("s3blob://%s?region=%s&credentials=%s&prefix=%s", bucket, region, creds, prefix)
	return NewBlobCacheFromURI(uri)
}
