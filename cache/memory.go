package cache

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"sync"
	"time"

	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
	gocache "github.com/patrickmn/go-cache"
)

// A MemoryCache represents a cache location in memory.
type MemoryCache struct {
	Cache
	provider      *gocache.Cache
	size          int
	maxsize       int
	sizemap       map[string]int
	keys          []string
	lock          *sync.Mutex
	eviction_lock *sync.Mutex
	uri           string
}

func init() {
	ctx := context.Background()
	err := RegisterMemoryCacheSchemes(ctx)
	if err != nil {
		panic(err)
	}
}

// RegisterMemoryCacheSchemes will...
func RegisterMemoryCacheSchemes(ctx context.Context) error {

	register_mu.Lock()
	defer register_mu.Unlock()

	schemes := []string{
		"memory",
	}

	for _, scheme := range schemes {

		_, exists := register_map[scheme]

		if exists {
			continue
		}

		err := RegisterCache(ctx, scheme, NewMemoryCacheFromURI)

		if err != nil {
			return fmt.Errorf("Failed to register blob cache for '%s', %w", scheme, err)
		}

		register_map[scheme] = true
	}

	return nil
}

// NewMemoryCacheURIFromConfig returns a valid cache.Cache URI derived from 'config'.
func NewMemoryCacheURIFromConfig(cfg iiifconfig.CacheConfig) (string, error) {

	q := url.Values{}

	if cfg.TTL > 0 {
		q.Set("ttl", strconv.Itoa(cfg.TTL))
	}

	if cfg.TTL > 0 {
		q.Set("limit", strconv.Itoa(cfg.Limit))
	}

	u := url.URL{}
	u.Scheme = "memory"
	u.RawQuery = q.Encode()

	return u.String(), nil
}

// NewMemoryCache returns a new `MemoryCache` instance derived from 'cfg'.
func NewMemoryCache(cfg iiifconfig.CacheConfig) (Cache, error) {

	uri := cfg.URI

	if uri == "" {

		v, err := NewMemoryCacheURIFromConfig(cfg)

		if err != nil {
			return nil, err
		}

		uri = v
	}

	return NewMemoryCacheFromURI(uri)
}

// NewMemoryCacheFromURI returns a new `MemoryCache` instance derived from 'uri'
func NewMemoryCacheFromURI(uri string) (Cache, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	ttl := 300
	limit := 100

	if q.Has("ttl") {

		v, err := strconv.Atoi(q.Get("ttl"))

		if err != nil {
			return nil, fmt.Errorf("Invalid ?ttl= parameter, %w", err)
		}

		ttl = v
	}

	if q.Has("limit") {

		v, err := strconv.Atoi(q.Get("limit"))

		if err != nil {
			return nil, fmt.Errorf("Invalid ?limit= parameter, %w", err)
		}

		limit = v
	}

	window := time.Duration(ttl) * time.Second

	gc := gocache.New(window, 30*time.Second)

	size := 0
	maxsize := limit * 1024 * 1024

	keys := make([]string, 0)
	sizemap := make(map[string]int)

	/*

		see this - it's two separate locking mechanisms - that's because we need to account
		for the sizemap and maxsize properties being updated during multiple Set events, one
		of which may be trying purge old records to make room and/or the normal gocache janitor
		cleaning up expired documents (20160911/thisisaaronland)

	*/

	lock := new(sync.Mutex)
	ev_lock := new(sync.Mutex)

	mc := MemoryCache{
		provider:      gc,
		size:          size,
		keys:          keys,
		maxsize:       maxsize,
		sizemap:       sizemap,
		lock:          lock,
		eviction_lock: ev_lock,
		uri:           uri,
	}

	gc.OnEvicted(mc.OnEvicted)

	return &mc, nil
}

func (mc *MemoryCache) String() string {
	return mc.uri
}

// Exists returns a bool set to true if the configured memory location exists.
func (mc *MemoryCache) Exists(key string) bool {

	_, ok := mc.provider.Get(key)

	return ok
}

// Get reads data from a memory location.
func (mc *MemoryCache) Get(key string) ([]byte, error) {

	data, ok := mc.provider.Get(key)

	if !ok {
		slog.Debug("Get cache (MISS)", "key", key)
		return nil, errors.New("cache miss")
	}

	slog.Debug("Get cache (HIT)", "key", key)
	return data.([]byte), nil
}

// Set writes data to a memory location.
func (mc *MemoryCache) Set(key string, data []byte) error {

	mc.lock.Lock()
	defer mc.lock.Unlock()

	_, ok := mc.sizemap[key]

	if ok {
		return nil
	}

	size := len(data)

	if size > mc.maxsize {
		return errors.New("key is too big")
	}

	if size+mc.size > mc.maxsize {

		for mc.size > mc.maxsize-size {

			for _, k := range mc.keys {
				mc.Unset(k)
			}
		}

	}

	mc.eviction_lock.Lock()
	defer mc.eviction_lock.Unlock()

	mc.size += size
	mc.sizemap[key] = size
	mc.keys = append(mc.keys, key)

	mc.provider.Set(key, data, gocache.DefaultExpiration)

	slog.Debug("Set cache (OK)", "key", key)
	return nil
}

// Unset deletes data from a memory location.
func (mc *MemoryCache) Unset(key string) error {

	slog.Debug("Unset cache", "key", key)
	mc.provider.Delete(key)
	return nil
}

func (mc *MemoryCache) OnEvicted(key string, value interface{}) {

	mc.eviction_lock.Lock()
	defer mc.eviction_lock.Unlock()

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
