package cache

import (
	_ "github.com/aaronland/gocloud-blob/s3"
	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
)

// NewS3Cache returns a NewBlobCacheFromURI with a constructed blob uri.
func NewS3Cache(cfg iiifconfig.CacheConfig) (Cache, error) {

	return NewBlobCache(cfg)
}

func NewS3CacheFromURI(uri string) (Cache, error) {
	return NewBlobCacheFromURI(uri)
}
