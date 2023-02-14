package server

import (
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	Server
	url *url.URL
}

func NewHTTPServer(u *url.URL, args ...interface{}) (Server, error) {

	u.Scheme = "http"

	server := HTTPServer{
		url: u,
	}

	return &server, nil
}

func (s *HTTPServer) Address() string {
	return s.url.String()
}

func (s *HTTPServer) ListenAndServe(router *mux.Router) error {
	return http.ListenAndServe(s.url.Host, router)
}
