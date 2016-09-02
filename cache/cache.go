package cache

import (
	"github.com/thisisaaronland/go-iiif/config"
)

type Cache interface {
	Get(string) ([]byte, error)
	Set(string, []byte) error
	Unset(string) error
}

func NewCacheFromConfig(cfg config.CacheConfig) (Cache, error) {

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
