package uri

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"regexp"
)

type URI interface {
	Driver() string
	String() string
	Origin() string
	Target(*url.Values) (string, error)
}

func NewURI(str_uri string) (URI, error) {

	u, err := NewURIWithDriver(str_uri)

	if err == nil {
		return u, nil
	}

	re, re_err := regexp.Compile(`^\w+\:\/\/`)

	if re_err != nil {
		return nil, re_err
	}

	if re.MatchString(str_uri) {
		msg := fmt.Sprintf("Invalid or unsupported URI string: %s", err)
		return nil, errors.New(msg)
	}

	file_uri := NewFileURIString(str_uri)
	log.Println(str_uri, file_uri)

	return NewURIWithDriver(file_uri)
}
