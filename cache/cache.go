// package cache provides functions for saving IIIF files to supported locations, including disk, memory, blob, and Amazon S3.
package cache

import (
	"strings"

	iiifconfig "github.com/go-iiif/go-iiif/v5/config"
)

// A Cache is a representation of a cache location.
type Cache interface {
	Exists(string) bool
	Get(string) ([]byte, error)
	Set(string, []byte) error
	Unset(string) error
}

// NewImagesCacheFromConfig returns a NewCacheFromConfig.
func NewImagesCacheFromConfig(config *iiifconfig.Config) (Cache, error) {

	cfg := config.Images.Cache
	return NewCacheFromConfig(cfg)
}

// NewDerivativesCacheFromConfig returns a NewCacheFromConfig.
func NewDerivativesCacheFromConfig(config *iiifconfig.Config) (Cache, error) {

	cfg := config.Derivatives.Cache
	return NewCacheFromConfig(cfg)
}

// NewCacheFromConfig returns a Cache object depending on the type of cache requested. Cache types can be blob, disk, memory, s3 or s3blob.
func NewCacheFromConfig(config iiifconfig.CacheConfig) (Cache, error) {

	var cache Cache
	var err error

	switch strings.ToLower(config.Name) {
	case "blob":
		cache, err = NewBlobCache(config)
	case "disk":
		cache, err = NewDiskCache(config)
	case "memory":
		cache, err = NewMemoryCache(config)
	case "s3":
		cache, err = NewS3Cache(config)
	case "s3blob":
		cache, err = NewS3Cache(config)
	default:
		cache, err = NewNullCache(config)
	}

	return cache, err
}
