package palette

import (
	"context"
	"embed"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-colours"
	"github.com/aaronland/go-roster"
)

//go:embed *.json
var FS embed.FS

type Palette interface {
	Reference() string
	Colours() []colours.Colour
}

var palette_roster roster.Roster

// PaletteInitializationFunc is a function defined by individual palette package and used to create
// an instance of that palette
type PaletteInitializationFunc func(ctx context.Context, uri string) (Palette, error)

// RegisterPalette registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `Palette` instances by the `NewPalette` method.
func RegisterPalette(ctx context.Context, scheme string, init_func PaletteInitializationFunc) error {

	err := ensurePaletteRoster()

	if err != nil {
		return err
	}

	return palette_roster.Register(ctx, scheme, init_func)
}

func ensurePaletteRoster() error {

	if palette_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		palette_roster = r
	}

	return nil
}

// NewPalette returns a new `Palette` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `PaletteInitializationFunc`
// function used to instantiate the new `Palette`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterPalette` method.
func NewPalette(ctx context.Context, uri string) (Palette, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := palette_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(PaletteInitializationFunc)
	return init_func(ctx, uri)
}

// Schemes returns the list of schemes that have been registered.
func PaletteSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensurePaletteRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range palette_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
