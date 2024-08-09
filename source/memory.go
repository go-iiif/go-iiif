package source

import (
	"crypto/sha256"
	"fmt"
	_ "log/slog"

	"github.com/dgraph-io/ristretto"
)

var memory_cache *ristretto.Cache

type MemorySource struct {
	Source
	key string
}

func NewMemorySource(body []byte) (Source, error) {

	sum := sha256.Sum256(body)
	key := fmt.Sprintf("%x", sum)

	return NewMemorySourceWithKey(key, body)
}

func NewMemorySourceWithKey(key string, body []byte) (Source, error) {

	if memory_cache == nil {

		cache, err := ristretto.NewCache(&ristretto.Config{
			NumCounters: 1e7,     // number of keys to track frequency of (10M).
			MaxCost:     1 << 30, // maximum cost of cache (1GB).
			BufferItems: 64,      // number of keys per Get buffer.
		})

		if err != nil {
			return nil, fmt.Errorf("Failed to create source memory cache, %w", err)
		}

		memory_cache = cache
	}

	memory_cache.Set(key, body, 1)
	memory_cache.Wait()

	bs := &MemorySource{
		key: key,
	}

	return bs, nil
}

func (bs *MemorySource) String() string {
	return fmt.Sprintf("memory://%s", bs.key)
}

func (bs *MemorySource) Read(key string) ([]byte, error) {

	v, exists := memory_cache.Get(key)

	if !exists {
		return nil, fmt.Errorf("%s not found", bs.key)
	}

	return v.([]byte), nil
}
