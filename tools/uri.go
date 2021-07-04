package tools

import (
	"context"
	iiifuri "github.com/go-iiif/go-iiif-uri"
)

type URIFunc func(string) (iiifuri.URI, error)

func DefaultURIFunc() URIFunc {

	fn := func(raw_uri string) (iiifuri.URI, error) {
		ctx := context.Background()
		return iiifuri.NewURI(ctx, raw_uri)
	}

	return fn
}
