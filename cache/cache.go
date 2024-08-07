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
	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
)

// In principle this could also be done with a sync.OnceFunc call but that will
// require that everyone uses Go 1.21 (whose package import changes broke everything)
// which is literally days old as I write this. So maybe a few releases after 1.21.
//
// Also, _not_ using a sync.OnceFunc means we can call RegisterSchemes multiple times
// if and when multiple gomail-sender instances register themselves.

var register_mu = new(sync.RWMutex)
var register_map = map[string]bool{}

// A Cache is a representation of a cache location.
type Cache interface {
	// Exists returns a boolean value indicating whether a key exists in the cache.
	Exists(string) bool
	// Get returns the value for a specific key in the cache.
	Get(string) ([]byte, error)
	// Set assigns the value for a specific key in the cache.
	Set(string, []byte) error
	// Unset removes a specific key from the cache.
	Unset(string) error
}

// NewImagesCacheFromConfig returns a NewCacheFromConfig.
func NewImagesCacheFromConfig(config *iiifconfig.Config) (Cache, error) {

	cfg := config.Images.Cache
	return NewCacheFromConfig(cfg)
}

// NewDerivativesCacheFromConfig returns a NewCacheFromConfig.
func NewDerivativesCacheFromConfig(config *iiifconfig.Config) (Cache, error) {

	cfg := config.Derivatives.Cache
	return NewCacheFromConfig(cfg)
}

// NewCacheFromConfig returns a Cache object depending on the type of cache requested. Cache types can be blob, disk, memory, s3 or s3blob.
func NewCacheFromConfig(config iiifconfig.CacheConfig) (Cache, error) {

	var cache_uri string
	var err error

	switch config.URI {
	case "":

		switch strings.ToLower(config.Name) {
		case "blob":
			cache_uri, err = NewBlobCacheURIFromConfig(config)
		case "disk":
			cache_uri, err = NewDiskCacheURIFromConfig(config)
		case "memory":
			cache_uri, err = NewMemoryCacheURIFromConfig(config)
		case "s3":
			cache_uri, err = NewS3CacheURIFromConfig(config)
		case "s3blob":
			cache_uri, err = NewS3CacheURIFromConfig(config)
		default:
			cache_uri, err = NewNullCacheURIFromConfig(config)
		}
	default:
		cache_uri = config.URI
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to derive cache URI, %w", err)
	}

	ctx := context.Background()
	return NewCache(ctx, cache_uri)
}

//

var cache_roster roster.Roster

// CacheInitializationFunc is a function defined by individual cache package and used to create
// an instance of that cache
type CacheInitializationFunc func(uri string) (Cache, error)

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
	return init_func(uri)
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
