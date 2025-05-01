package grid

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-colours"
	"github.com/aaronland/go-colours/palette"
	"github.com/aaronland/go-roster"
)

type Grid interface {
	Closest(context.Context, colours.Colour, palette.Palette) (colours.Colour, error)
}

var grid_roster roster.Roster

// GridInitializationFunc is a function defined by individual grid package and used to create
// an instance of that grid
type GridInitializationFunc func(ctx context.Context, uri string) (Grid, error)

// RegisterGrid registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `Grid` instances by the `NewGrid` method.
func RegisterGrid(ctx context.Context, scheme string, init_func GridInitializationFunc) error {

	err := ensureGridRoster()

	if err != nil {
		return err
	}

	return grid_roster.Register(ctx, scheme, init_func)
}

func ensureGridRoster() error {

	if grid_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		grid_roster = r
	}

	return nil
}

// NewGrid returns a new `Grid` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `GridInitializationFunc`
// function used to instantiate the new `Grid`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterGrid` method.
func NewGrid(ctx context.Context, uri string) (Grid, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := grid_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(GridInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func GridSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureGridRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range grid_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
