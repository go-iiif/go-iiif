package cache

import (
	"errors"
	gocache "github.com/patrickmn/go-cache"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	"log"
	"time"
)

type MemoryCache struct {
	Cache
	provider *gocache.Cache
}

func NewMemoryCache(cfg iiifconfig.CacheConfig) (*MemoryCache, error) {

	// ttl := cfg.TTL
	// limit := cfg.Limit
	// window := time.Duration(ttl) * time.Second

	gc := gocache.New(5*time.Minute, 30*time.Second)

	mc := MemoryCache{
		provider: gc,
	}

	return &mc, nil
}

func (mc *MemoryCache) Get(key string) ([]byte, error) {

	log.Println("GET", key)

	data, ok := mc.provider.Get(key)

	if !ok {

		log.Println("MISS", key)
		return nil, errors.New("cache miss")
	}

	return data.([]byte), nil
}

func (mc *MemoryCache) Set(key string, data []byte) error {

	log.Println("SET", key)
	mc.provider.Set(key, data, gocache.DefaultExpiration)

	return nil
}

func (mc *MemoryCache) Unset(key string) error {

	mc.provider.Delete(key)
	return nil
}
