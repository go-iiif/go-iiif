package uri

import (
	"errors"
)

type URI interface {
	URL() string
	String() string
}

func NewURIWithType(str_uri string, str_type string) (URI, error) {

	var u URI
	var e error

	switch str_type {

	case "string":
		u, e = NewStringURI(str_uri)
	case "idsecret":
		u, e = NewIdSecretURI(str_uri)
	default:
		e = errors.New("Unknown URI type")
	}

	return u, e
}
