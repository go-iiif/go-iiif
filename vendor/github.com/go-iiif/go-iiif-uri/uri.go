package uri

import (
)

type URI interface {
	Driver() string
	String() string
	Origin() string
	Target(...interface{}) string
}

func NewURI(str_uri string) (URI, error) {

	return NewURIWithDriver(str_uri)
}
