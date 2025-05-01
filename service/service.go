package service

// http://iiif.io/api/image/2.1/#related-services
// http://iiif.io/api/annex/services/

import (
	"context"
	"net/url"

	"github.com/aaronland/go-roster"
	iiifconfig "github.com/go-iiif/go-iiif/v8/config"
	iiifimage "github.com/go-iiif/go-iiif/v8/image"
)

var service_roster roster.Roster

type ServiceInitializationFunc func(ctx context.Context, config *iiifconfig.Config, im iiifimage.Image) (Service, error)

type Service interface {
	Context() string
	Profile() string
	Label() string
	Value() interface{}
}

func NewService(ctx context.Context, uri string, cfg *iiifconfig.Config, im iiifimage.Image) (Service, error) {

	err := ensureServiceRoster()

	if err != nil {
		return nil, err
	}

	parsed, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := parsed.Scheme

	i, err := service_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(ServiceInitializationFunc)
	return init_func(ctx, cfg, im)
}

func RegisterService(ctx context.Context, scheme string, init_func ServiceInitializationFunc) error {

	err := ensureServiceRoster()

	if err != nil {
		return err
	}

	return service_roster.Register(ctx, scheme, init_func)
}

func ensureServiceRoster() error {

	if service_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		service_roster = r
	}

	return nil
}
