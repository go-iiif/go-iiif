package cache

import (
	"fmt"
	"github.com/go-iiif/go-iiif/v5/config"
	_ "log"
)

func NewDiskCache(cfg config.CacheConfig) (Cache, error) {

	root := cfg.Path
	uri := fmt.Sprintf("file://%s", root)

	return NewBlobCacheFromURI(uri)
}
