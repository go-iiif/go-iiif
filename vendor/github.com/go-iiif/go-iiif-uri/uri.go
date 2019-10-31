package uri

import (
	"net/url"
)

type URI interface {
	Driver() string
	String() string
	Origin() string
	Target(*url.Values) (string, error)
}

func NewURI(str_uri string) (URI, error) {

	return NewURIWithDriver(str_uri)
}
