package cache

import (
	"fmt"

	"github.com/go-iiif/go-iiif/v6/config"
)

// NewDiskCache returns a NewBlobCacheFromURI for a local files system location.
func NewDiskCache(cfg config.CacheConfig) (Cache, error) {

	root := cfg.Path
	uri := fmt.Sprintf("file://%s", root)

	return NewBlobCacheFromURI(uri)
}
