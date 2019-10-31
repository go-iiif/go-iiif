package uri

import (
)

type URI interface {
	String() string
	Origin() string
	Target() string
}

func NewURI(str_uri string) (URI, error) {

	return NewURIWithDriver(str_uri)
}
