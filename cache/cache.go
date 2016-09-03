package cache

import (
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
)

type Cache interface {
	Get(string) ([]byte, error)
	Set(string, []byte) error
	Unset(string) error
}

func NewCacheFromConfig(config *iiifconfig.Config) (Cache, error) {

	cfg := config.Derivatives.Cache

	if cfg.Name == "Disk" {
		cache, err := NewDiskCache(cfg)
		return cache, err
	} else if cfg.Name == "Memory" {
		cache, err := NewMemoryCache(cfg)
		return cache, err
	} else {
		cache, err := NewNullCache(cfg)
		return cache, err
	}
}
