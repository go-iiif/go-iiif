package cache

import (
	"errors"

	"github.com/go-iiif/go-iiif/v6/config"
)

// A NullCache represents a Cache.
type NullCache struct {
	Cache
}

func NewNullCacheURIFromConfig(cfg config.CacheConfig) (string, error) {
	return "null://", nil
}

// NewNullCache returns a pointer to a NullCache.
func NewNullCache(cfg config.CacheConfig) (Cache, error) {
	uri, _ := NewNullCacheURIFromConfig(cfg)
	return NewNullCacheFromURI(uri)
}

func NewNullCacheFromURI(uri string) (Cache, error) {
	c := NullCache{}
	return &c, nil
}

// Exists is always false for a NullCache.
func (c *NullCache) Exists(rel_path string) bool {
	return false
}

// Get returns nil and an error message.
func (c *NullCache) Get(rel_path string) ([]byte, error) {

	err := errors.New("null cache is null")
	return nil, err
}

// Set returns nil.
func (c *NullCache) Set(rel_path string, body []byte) error {

	return nil
}

// Unset returns nil.
func (c *NullCache) Unset(rel_path string) error {

	return nil
}
