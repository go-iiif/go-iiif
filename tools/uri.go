package tools

import (
	iiifuri "github.com/go-iiif/go-iiif-uri"
	"net/url"
)

type URIFunc func(string) (iiifuri.URI, error)

func DefaultURIFunc() URIFunc {

	fn := func(raw_uri string) (iiifuri.URI, error) {
		return iiifuri.NewURI(raw_uri)
	}

	return fn
}

func ReprocessIdSecretURIFunc(id string, secret_o string, secret string) URIFunc {

	fn := func(raw_uri string) (iiifuri.URI, error) {

		secret_q := url.Values{}
		secret_q.Set("id", id)
		secret_q.Set("secret", secret)
		secret_q.Set("secret_o", secret_o)

		secret_u := url.URL{}
		secret_u.Scheme = "idsecret"
		secret_u.RawQuery = secret_q.Encode()

		return iiifuri.NewURI(secret_u.String())
	}

	return fn
}
