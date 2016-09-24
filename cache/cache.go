package cache

import (
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
)

type Cache interface {
	Exists(string) bool
	Get(string) ([]byte, error)
	Set(string, []byte) error
	Unset(string) error
}

func NewImagesCacheFromConfig(config *iiifconfig.Config) (Cache, error) {

	cfg := config.Images.Cache
	return NewCacheFromConfig(cfg)
}

func NewDerivativesCacheFromConfig(config *iiifconfig.Config) (Cache, error) {

	cfg := config.Derivatives.Cache
	return NewCacheFromConfig(cfg)
}

func NewCacheFromConfig(cfg iiifconfig.CacheConfig) (Cache, error) {

	if cfg.Name == "Disk" {
		cache, err := NewDiskCache(cfg)
		return cache, err
	} else if cfg.Name == "Memory" {
		cache, err := NewMemoryCache(cfg)
		return cache, err
	} else if cfg.Name == "S3" {
		cache, err := NewS3Cache(cfg)
		return cache, err
	} else {
		cache, err := NewNullCache(cfg)
		return cache, err
	}
}
