package server

import (
	"context"
	"github.com/aaronland/go-roster"
	_ "log"
	"net/http"
	"net/url"
	"strings"
)

type Server interface {
	ListenAndServe(context.Context, http.Handler) error
	Address() string
}

type ServerInitializeFunc func(context.Context, string) (Server, error)

var servers roster.Roster

func ensureServers() error {

	if servers == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		servers = r
	}

	return nil
}

func RegisterServer(ctx context.Context, scheme string, f ServerInitializeFunc) error {

	err := ensureServers()

	if err != nil {
		return err
	}

	return servers.Register(ctx, scheme, f)
}

func NewServer(ctx context.Context, uri string) (Server, error) {

	err := ensureServers()

	if err != nil {
		return nil, err
	}

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := servers.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	f := i.(ServerInitializeFunc)
	return f(ctx, uri)
}

func Schemes() []string {
	ctx := context.Background()
	return servers.Drivers(ctx)
}

func SchemesAsString() string {
	schemes := Schemes()
	return strings.Join(schemes, ",")
}
