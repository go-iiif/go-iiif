package driver

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/aaronland/go-roster"
	iiifcache "github.com/go-iiif/go-iiif/v8/cache"
	iiifconfig "github.com/go-iiif/go-iiif/v8/config"
	iiifimage "github.com/go-iiif/go-iiif/v8/image"
	iiifsource "github.com/go-iiif/go-iiif/v8/source"
)

type Driver interface {
	NewImageFromConfigWithSource(context.Context, *iiifconfig.Config, iiifsource.Source, string) (iiifimage.Image, error)
	NewImageFromConfigWithCache(context.Context, *iiifconfig.Config, iiifcache.Cache, string) (iiifimage.Image, error)
	NewImageFromConfig(context.Context, *iiifconfig.Config, string) (iiifimage.Image, error)
}

type DriverInitializeFunc func(ctx context.Context, uri string) (Driver, error)

var driver_roster roster.Roster

func ensureSpatialRoster() error {

	if driver_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		driver_roster = r
	}

	return nil
}

func RegisterDriver(ctx context.Context, scheme string, f DriverInitializeFunc) error {

	err := ensureSpatialRoster()

	if err != nil {
		return err
	}

	return driver_roster.Register(ctx, scheme, f)
}

func Schemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureSpatialRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range driver_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}

func NewDriver(ctx context.Context, driver_uri string) (Driver, error) {

	u, err := url.Parse(driver_uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	if scheme == "" {

		scheme = "file"
		u.Scheme = scheme

		if !strings.HasPrefix(u.Path, "/") {
			u.Path = fmt.Sprintf("/%s", u.Path)
		}

		driver_uri = u.String()
	}

	i, err := driver_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	f := i.(DriverInitializeFunc)
	return f(ctx, driver_uri)
}
