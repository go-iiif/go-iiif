package source

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

// Source is an interface representing a primary image source.
type Source interface {
	// Read returns the body of the file located at 'uri'.
	Read(uri string) ([]byte, error)
	// Close performs any final operations specific to a data source.
	Close() error
}

var source_roster roster.Roster

// SourceInitializationFunc is a function defined by individual source package and used to create
// an instance of that source
type SourceInitializationFunc func(context.Context, string) (Source, error)

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
	return init_func(ctx, uri)
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
