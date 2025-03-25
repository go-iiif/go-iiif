// package cache provides functions for saving IIIF files to supported locations, including disk, memory, blob, and Amazon S3.
package cache

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"sync"

	"github.com/aaronland/go-roster"
)

// In principle this could also be done with a sync.OnceFunc call but that will
// require that everyone uses Go 1.21 (whose package import changes broke everything)
// which is literally days old as I write this. So maybe a few releases after 1.21.
//
// Also, _not_ using a sync.OnceFunc means we can call RegisterSchemes multiple times
// if and when multiple gomail-sender instances register themselves.

var register_mu = new(sync.RWMutex)
var register_map = map[string]bool{}

// A Cache is a representation of a cache provider.
type Cache interface {
	// Exists returns a boolean value indicating whether a key exists in the cache.
	Exists(string) bool
	// Get returns the value for a specific key in the cache.
	Get(string) ([]byte, error)
	// Set assigns the value for a specific key in the cache.
	Set(string, []byte) error
	// Unset removes a specific key from the cache.
	Unset(string) error
	// Close performs any final operations specific to a cache provider.
	Close() error
}

var cache_roster roster.Roster

// CacheInitializationFunc is a function defined by individual cache package and used to create
// an instance of that cache
type CacheInitializationFunc func(context.Context, string) (Cache, error)

// RegisterCache registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `Cache` instances by the `NewCache` method.
func RegisterCache(ctx context.Context, scheme string, init_func CacheInitializationFunc) error {

	err := ensureCacheRoster()

	if err != nil {
		return err
	}

	return cache_roster.Register(ctx, scheme, init_func)
}

func ensureCacheRoster() error {

	if cache_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		cache_roster = r
	}

	return nil
}

// NewCache returns a new `Cache` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `CacheInitializationFunc`
// function used to instantiate the new `Cache`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterCache` method.
func NewCache(ctx context.Context, uri string) (Cache, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := cache_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(CacheInitializationFunc)
	return init_func(ctx, uri)
}

// CacheSchemes returns the list of schemes that have been registered.
func CacheSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureCacheRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range cache_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
