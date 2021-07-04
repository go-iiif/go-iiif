package server

// https://github.com/FiloSottile/mkcert/issues/154
// https://smallstep.com/blog/private-acme-server/

import (
	"context"
	"errors"
	"net/url"
)

func init() {
	ctx := context.Background()
	RegisterServer(ctx, "https", NewTLSServer)
	RegisterServer(ctx, "tls", NewTLSServer)
}

func NewTLSServer(ctx context.Context, uri string) (Server, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	tls_cert := q.Get("cert")
	tls_key := q.Get("key")

	if tls_cert == "" {
		return nil, errors.New("Missing TLS cert parameter")
	}

	if tls_key == "" {
		return nil, errors.New("Missing TLS key parameter")
	}

	return NewHTTPServer(ctx, uri)
}
