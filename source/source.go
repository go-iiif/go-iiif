package source

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
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

// Source is an interface representing a primary image source.
type Source interface {
	// Read returns the body of the file located at 'uri'.
	Read(uri string) ([]byte, error)
}

func NewSourceFromConfig(config *iiifconfig.Config) (Source, error) {

	cfg := config.Images

	var source_uri string
	var err error

	switch cfg.Source.URI {
	case "":

		switch strings.ToLower(cfg.Source.Name) {
		case "blob":
			source_uri, err = NewBlobSourceURIFromConfig(config)
		// case "memory":
		// 	source_uri, err = NewMemorySourceURIFromConfig(config)
		case "disk":
			source_uri, err = NewDiskSourceURIFromConfig(config)
		case "flickr":
			source_uri, err = NewFlickrSourceURIFromConfig(config)
		case "s3":
			source_uri, err = NewS3SourceURIFromConfig(config)
		case "s3blob":
			source_uri, err = NewS3SourceURIFromConfig(config)
		case "uri":
			source_uri, err = NewURISourceURIFromConfig(config)
		default:
			err = errors.New("Unknown source type")
		}
	default:
		source_uri = config.Images.Source.URI
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to derive source URI, %w", err)
	}

	slog.Debug("Create new source", "uri", source_uri)

	ctx := context.Background()
	return NewSource(ctx, source_uri)
}

//

var source_roster roster.Roster

// SourceInitializationFunc is a function defined by individual source package and used to create
// an instance of that source
type SourceInitializationFunc func(uri string) (Source, error)

// RegisterSource registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `Source` instances by the `NewSource` method.
func RegisterSource(ctx context.Context, scheme string, init_func SourceInitializationFunc) error {

	err := ensureSourceRoster()

	if err != nil {
		return err
	}

	return source_roster.Register(ctx, scheme, init_func)
}

func ensureSourceRoster() error {

	if source_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		source_roster = r
	}

	return nil
}

// NewSource returns a new `Source` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `SourceInitializationFunc`
// function used to instantiate the new `Source`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterSource` method.
func NewSource(ctx context.Context, uri string) (Source, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := source_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(SourceInitializationFunc)
	return init_func(uri)
}

// SourceSchemes returns the list of schemes that have been registered.
func SourceSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureSourceRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range source_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
