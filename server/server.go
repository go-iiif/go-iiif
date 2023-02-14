package server

import (
	"errors"
	_ "log"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
)

type Server interface {
	ListenAndServe(*mux.Router) error
	Address() string
}

func NewServer(proto string, u *url.URL, args ...interface{}) (Server, error) {

	var svr Server
	var err error

	switch strings.ToUpper(proto) {

	case "HTTP":

		svr, err = NewHTTPServer(u, args...)

	case "LAMBDA":

		svr, err = NewLambdaServer(u, args...)

	default:
		return nil, errors.New("Invalid server protocol")
	}

	if err != nil {
		return nil, err
	}

	return svr, nil
}
