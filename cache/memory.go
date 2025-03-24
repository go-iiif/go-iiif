package cache

// This is nearly indistinguishable from cache/memory.go and maybe they should
// be reconciled but not today...

import (
	"context"
	"fmt"

	"github.com/dgraph-io/ristretto/v2"
)

// A MemoryCache represents a cache location in memory.
type MemoryCache struct {
	Cache
	cache *ristretto.Cache[string, []byte]
	uri   string
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

		err := RegisterCache(ctx, scheme, NewMemoryCache)

		if err != nil {
			return fmt.Errorf("Failed to register blob cache for '%s', %w", scheme, err)
		}

		register_map[scheme] = true
	}

	return nil
}

// NewMemoryCache returns a new `MemoryCache` instance derived from 'uri'
func NewMemoryCache(ctx context.Context, uri string) (Cache, error) {

	cache, err := ristretto.NewCache(&ristretto.Config[string, []byte]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to create source memory cache, %w", err)
	}

	mc := MemoryCache{
		cache: cache,
		uri:   uri,
	}

	return &mc, nil
}

func (mc *MemoryCache) String() string {
	return mc.uri
}

// Exists returns a bool set to true if the configured memory location exists.
func (mc *MemoryCache) Exists(key string) bool {

	_, ok := mc.cache.Get(key)

	return ok
}

// Get reads data from a memory location.
func (mc *MemoryCache) Get(key string) ([]byte, error) {

	data, ok := mc.cache.Get(key)

	if !ok {
		return nil, fmt.Errorf("Cache miss")
	}

	return data, nil
}

// Set writes data to a memory location.
func (mc *MemoryCache) Set(key string, data []byte) error {

	ok := mc.cache.Set(key, data, 1)

	if !ok {
		return fmt.Errorf("Failed to set key")
	}

	return nil
}

// Unset deletes data from a memory location.
func (mc *MemoryCache) Unset(key string) error {
	mc.cache.Del(key)
	return nil
}

func (mc *MemoryCache) Close() error {
	mc.cache.Close()
	return nil
}
