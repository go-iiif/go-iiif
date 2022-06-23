package server

import (
	"context"
	"fmt"
	"github.com/aaronland/go-roster"
	_ "log"
	"net/http"
	"net/url"
	"sort"
)

// type Server is an interface for creating server instances that serve requests using a `http.Handler` router.
type Server interface {
	// ListenAndServe starts the server and listens for requests using a `http.Handler` instance for routing.
	ListenAndServe(context.Context, http.Handler) error
	// Address returns the fully-qualified URI that the server is listening for requests on.
	Address() string
}

// ServeritializeFunc is a function used to initialize an implementation of the `Server` interface.
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

// RegisterServer() associates 'scheme' with 'f' in an internal list of avilable `Server` implementations.
func RegisterServer(ctx context.Context, scheme string, f ServerInitializeFunc) error {

	err := ensureServers()

	if err != nil {
		return err
	}

	return servers.Register(ctx, scheme, f)
}

// NewServer() returns a new instance of `Server` for the scheme associated with 'uri'. It is assumed that this scheme
// will have previously been "registered" with the `RegisterServer` method.
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

// Schemes() returns the list of schemes that have been "registered".
func Schemes() []string {
	ctx := context.Background()
	drivers := servers.Drivers(ctx)

	schemes := make([]string, len(drivers))

	for idx, dr := range drivers {
		schemes[idx] = fmt.Sprintf("%s://", dr)
	}

	sort.Strings(schemes)
	return schemes
}
