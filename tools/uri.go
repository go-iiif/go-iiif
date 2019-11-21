package tools

import (
	iiifuri "github.com/go-iiif/go-iiif-uri"
)

type URIFunc func(string) (iiifuri.URI, error)

func DefaultURIFunc() URIFunc {

	fn := func(raw_uri string) (iiifuri.URI, error) {
		return iiifuri.NewURI(raw_uri)
	}

	return fn
}
