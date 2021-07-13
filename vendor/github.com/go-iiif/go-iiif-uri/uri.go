package uri

import (
	"context"
	"fmt"
	"github.com/aaronland/go-roster"
	_ "log"
	"net/url"
	"sort"
	"strings"
)

type URI interface {
	Scheme() string
	String() string
	Origin() string
	Target(*url.Values) (string, error)
}

type URIInitializeFunc func(ctx context.Context, uri string) (URI, error)

var uri_roster roster.Roster

func ensureSpatialRoster() error {

	if uri_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		uri_roster = r
	}

	return nil
}

func RegisterURI(ctx context.Context, scheme string, f URIInitializeFunc) error {

	err := ensureSpatialRoster()

	if err != nil {
		return err
	}

	return uri_roster.Register(ctx, scheme, f)
}

func Schemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureSpatialRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range uri_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}

func NewURI(ctx context.Context, uri string) (URI, error) {

	u, err := url.Parse(uri)

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

		uri = u.String()
	}

	i, err := uri_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	f := i.(URIInitializeFunc)
	return f(ctx, uri)
}
