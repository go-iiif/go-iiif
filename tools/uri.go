package tools

import (
	"context"

	iiifuri "github.com/go-iiif/go-iiif-uri"
)

// URIFunc is custom function to derive a `iiifuri.URI` from a string.
type URIFunc func(string) (iiifuri.URI, error)

// DefaultURIFunc returns a function to create a `iiifuri.URI` instance from a string
// using the `iiifuri.NewURI` function.
func DefaultURIFunc() URIFunc {

	fn := func(raw_uri string) (iiifuri.URI, error) {
		ctx := context.Background()
		return iiifuri.NewURI(ctx, raw_uri)
	}

	return fn
}
