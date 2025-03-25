package cache

import (
	"context"
	"fmt"
)

// A NullCache represents a Cache.
type NullCache struct {
	Cache
}

func NewNullCache(ctx context.Context, uri string) (Cache, error) {
	c := NullCache{}
	return &c, nil
}

func (c *NullCache) String() string {
	return "null://"
}

// Exists is always false for a NullCache.
func (c *NullCache) Exists(rel_path string) bool {
	return false
}

// Get returns nil and an error message.
func (c *NullCache) Get(rel_path string) ([]byte, error) {
	return nil, fmt.Errorf("Null cache is null")
}

// Set returns nil.
func (c *NullCache) Set(rel_path string, body []byte) error {

	return nil
}

// Unset returns nil.
func (c *NullCache) Unset(rel_path string) error {
	return nil
}

func (c *NullCache) Close() error {
	return nil
}
