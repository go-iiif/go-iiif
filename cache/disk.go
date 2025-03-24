package cache

import (
	"fmt"

	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
)

// NewDiskCache returns a NewBlobCacheFromURI for a local files system location.
func NewDiskCache(cfg iiifconfig.CacheConfig) (Cache, error) {

	return NewBlobCacheFromURI(cfg.URI)
}

func NewDiskCacheFromURI(uri string) (Cache, error) {
	return NewBlobCacheFromURI(uri)
}
