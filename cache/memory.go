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
	keys     []string
	lock     *sync.Mutex
}

func NewMemoryCache(cfg iiifconfig.CacheConfig) (*MemoryCache, error) {

	ttl := cfg.TTL
	limit := cfg.Limit
	window := time.Duration(ttl) * time.Second

	gc := gocache.New(window, 30*time.Second)

	size := 0
	maxsize := limit * 1024 * 1024

	keys := make([]string, 0)
	sizemap := make(map[string]int)

	lock := new(sync.Mutex)

	mc := MemoryCache{
		provider: gc,
		size:     size,
		keys:     keys,
		maxsize:  maxsize,
		sizemap:  sizemap,
		lock:     lock,
	}

	gc.OnEvicted(mc.OnEvicted)

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

	_, ok := mc.sizemap[key]

	if ok {
		return nil
	}

	size := len(data)

	if size > mc.maxsize {
		return errors.New("Key is too big!")
	}

	if size+mc.size > mc.maxsize {

		for mc.size > mc.maxsize-size {

			for _, k := range mc.keys {
				mc.Unset(k)
			}
		}

	}

	mc.size += size
	mc.sizemap[key] = size
	mc.keys = append(mc.keys, key)

	mc.provider.Set(key, data, gocache.DefaultExpiration)

	return nil
}

func (mc *MemoryCache) Unset(key string) error {

	mc.provider.Delete(key)
	return nil
}

func (mc *MemoryCache) OnEvicted(key string, value interface{}) {

	mc.lock.Lock()
	defer mc.lock.Unlock()

	size, _ := mc.sizemap[key]
	mc.size -= size

	delete(mc.sizemap, key)

	new_keys := make([]string, 0)

	for _, k := range mc.keys {

		if k != key {
			new_keys = append(new_keys, k)
		}
	}

	mc.keys = new_keys
}
