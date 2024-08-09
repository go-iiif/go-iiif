package cache

import (
	"fmt"

	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
)

// NewDiskCacheURIFromConfig returns a valid cache.Cache URI derived from 'config'.
func NewDiskCacheURIFromConfig(cfg iiifconfig.CacheConfig) (string, error) {

	root := cfg.Path
	uri := fmt.Sprintf("file://%s", root)
	return uri, nil
}

// NewDiskCache returns a NewBlobCacheFromURI for a local files system location.
func NewDiskCache(cfg iiifconfig.CacheConfig) (Cache, error) {

	uri := cfg.URI

	if uri == "" {

		v, err := NewDiskCacheURIFromConfig(cfg)

		if err != nil {
			return nil, err
		}

		uri = v
	}

	return NewBlobCacheFromURI(uri)
}

func NewDiskCacheFromURI(uri string) (Cache, error) {
	return NewBlobCacheFromURI(uri)
}
