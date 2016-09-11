package cache

import (
	"errors"
	gocache "github.com/patrickmn/go-cache"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	_ "log"
	"sync"
	"time"
)

type MemoryCache struct {
	Cache
	provider *gocache.Cache
	size     int
	maxsize  int
	sizemap  map[string]int
	lock     *sync.Mutex
}

func NewMemoryCache(cfg iiifconfig.CacheConfig) (*MemoryCache, error) {

	ttl := cfg.TTL
	limit := cfg.Limit
	window := time.Duration(ttl) * time.Second

	gc := gocache.New(window, 30*time.Second)

	size := 0
	maxsize := limit * 1024 * 1024

	sizemap := make(map[string]int)

	lock := new(sync.Mutex)

	mc := MemoryCache{
		provider: gc,
		size:     size,
		maxsize:  maxsize,
		sizemap:  sizemap,
		lock:     lock,
	}

	return &mc, nil
}

func (mc *MemoryCache) Get(key string) ([]byte, error) {

	data, ok := mc.provider.Get(key)

	if !ok {
		return nil, errors.New("cache miss")
	}

	return data.([]byte), nil
}

func (mc *MemoryCache) Set(key string, data []byte) error {

	mc.lock.Lock()
	defer mc.lock.Unlock()

	size := len(data)

	if size > mc.maxsize {
		return errors.New("Key is too big!")
	}

	if size+mc.size > mc.maxsize {

		// please prune me here...
		return errors.New("No more space!")
	}

	mc.size += size
	mc.sizemap[key] = size

	mc.provider.Set(key, data, gocache.DefaultExpiration)
	return nil
}

func (mc *MemoryCache) Unset(key string) error {

	mc.lock.Lock()
	defer mc.lock.Unlock()

	size, _ := mc.sizemap[key]
	mc.size -= size

	delete(mc.sizemap, key)

	mc.provider.Delete(key)
	return nil
}

func (mc *MemoryCache) Prune() error {

	return nil
}
