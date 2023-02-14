package server

import (
	_ "log"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/whosonfirst/algnhsa"
)

type LambdaServer struct {
	Server
	url *url.URL
}

func NewLambdaServer(u *url.URL, args ...interface{}) (Server, error) {

	server := LambdaServer{
		url: u,
	}

	return &server, nil
}

func (s *LambdaServer) Address() string {
	return s.url.String()
}

func (s *LambdaServer) ListenAndServe(router *mux.Router) error {

	// this cr^H^H^H stuff is important (20180713/thisisaaronland)
	// go-rasterzen/README.md#lambda-api-gateway-and-images#lambda-api-gateway-and-images

	lambda_opts := new(algnhsa.Options)
	lambda_opts.BinaryContentTypes = []string{
		"image/jpeg",
		"image/gif",
		"image/png",
	}

	algnhsa.ListenAndServe(router, lambda_opts)
	return nil
}
