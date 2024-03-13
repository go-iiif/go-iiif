package cache

import (
	"errors"

	"github.com/go-iiif/go-iiif/v6/config"
)

// A NullCache represents a Cache.
type NullCache struct {
	Cache
}

// NewNullCache returns a pointer to a NullCache.
func NewNullCache(cfg config.CacheConfig) (*NullCache, error) {

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
