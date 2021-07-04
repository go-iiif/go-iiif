package uri

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

const REWRITE_SCHEME string = "rewrite"

type RewriteURI struct {
	URI
	origin string
	target string
}

func init() {
	ctx := context.Background()
	RegisterURI(ctx, REWRITE_SCHEME, NewRewriteURI)
}

func NewRewriteURI(ctx context.Context, str_uri string) (URI, error) {

	u, err := url.Parse(str_uri)

	if err != nil {
		return nil, err
	}

	origin := strings.TrimLeft(u.Path, "/")

	if origin == "" {
		return nil, errors.New("Invalid path")
	}

	q := u.Query()

	target := q.Get("target")

	if target == "" {
		return nil, errors.New("Missing rewrite target")
	}

	if target == origin {
		return nil, errors.New("Invalid rewrite target")
	}

	rw := RewriteURI{
		origin: origin,
		target: target,
	}

	return &rw, nil
}

func (u *RewriteURI) Origin() string {
	return u.origin
}

func (u *RewriteURI) Target(opts *url.Values) (string, error) {
	return u.target, nil
}

func (u *RewriteURI) String() string {

	q := url.Values{}
	q.Set("target", u.target)

	raw_uri := fmt.Sprintf("%s?%s", u.origin, q.Encode())
	return fmt.Sprintf("rewrite:///%s", raw_uri)
}

func (u *RewriteURI) Scheme() string {
	return REWRITE_SCHEME
}
