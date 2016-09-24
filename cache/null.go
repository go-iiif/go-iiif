package cache

import (
	"errors"
	"github.com/thisisaaronland/go-iiif/config"
)

type NullCache struct {
	Cache
}

func NewNullCache(cfg config.CacheConfig) (*NullCache, error) {

	c := NullCache{}

	return &c, nil
}

func (c *NullCache) Exists(rel_path string) bool {
	return false
}

func (c *NullCache) Get(rel_path string) ([]byte, error) {

	err := errors.New("null cache is null")
	return nil, err
}

func (c *NullCache) Set(rel_path string, body []byte) error {

	return nil
}

func (c *NullCache) Unset(rel_path string) error {

	return nil
}
