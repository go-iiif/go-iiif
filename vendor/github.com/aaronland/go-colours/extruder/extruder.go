package extruder

import (
	"context"
	"fmt"
	"image"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-colours"
	"github.com/aaronland/go-roster"
)

type Extruder interface {
	Colours(image.Image, int) ([]colours.Colour, error)
	Name() string
}

var extruder_roster roster.Roster

// ExtruderInitializationFunc is a function defined by individual extruder package and used to create
// an instance of that extruder
type ExtruderInitializationFunc func(ctx context.Context, uri string) (Extruder, error)

// RegisterExtruder registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `Extruder` instances by the `NewExtruder` method.
func RegisterExtruder(ctx context.Context, scheme string, init_func ExtruderInitializationFunc) error {

	err := ensureExtruderRoster()

	if err != nil {
		return err
	}

	return extruder_roster.Register(ctx, scheme, init_func)
}

func ensureExtruderRoster() error {

	if extruder_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		extruder_roster = r
	}

	return nil
}

// NewExtruder returns a new `Extruder` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `ExtruderInitializationFunc`
// function used to instantiate the new `Extruder`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterExtruder` method.
func NewExtruder(ctx context.Context, uri string) (Extruder, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := extruder_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(ExtruderInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func ExtruderSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureExtruderRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range extruder_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
