package cache

import (
	"fmt"

	_ "github.com/aaronland/gocloud-blob-s3"
	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
)

// NewS3Cache returns a NewBlobCacheFromURI with a constructed blob uri.
func NewS3Cache(cfg iiifconfig.CacheConfig) (Cache, error) {

	bucket := cfg.Path
	prefix := cfg.Prefix
	region := cfg.Region
	creds := cfg.Credentials

	uri := fmt.Sprintf("s3blob://%s?region=%s&credentials=%s&prefix=%s", bucket, region, creds, prefix)
	return NewBlobCacheFromURI(uri)
}
