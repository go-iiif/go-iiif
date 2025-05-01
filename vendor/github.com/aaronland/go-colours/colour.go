package colours

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-roster"
)

type Colour interface {
	Name() string
	Hex() string
	Reference() string
	Closest() []Colour
	AppendClosest(Colour) error // I don't love this... (20180605/thisisaaronland)
	String() string
}

var colour_roster roster.Roster

// ColourInitializationFunc is a function defined by individual colour package and used to create
// an instance of that colour
type ColourInitializationFunc func(ctx context.Context, uri string) (Colour, error)

// RegisterColour registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `Colour` instances by the `NewColour` method.
func RegisterColour(ctx context.Context, scheme string, init_func ColourInitializationFunc) error {

	err := ensureColourRoster()

	if err != nil {
		return err
	}

	return colour_roster.Register(ctx, scheme, init_func)
}

func ensureColourRoster() error {

	if colour_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		colour_roster = r
	}

	return nil
}

// NewColour returns a new `Colour` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `ColourInitializationFunc`
// function used to instantiate the new `Colour`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterColour` method.
func NewColour(ctx context.Context, uri string) (Colour, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := colour_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(ColourInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func ColourSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureColourRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range colour_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
