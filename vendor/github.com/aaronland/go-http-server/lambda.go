package server

import (
	"context"
	"github.com/whosonfirst/algnhsa"
	_ "log"
	"net/http"
	"net/url"
)

func init() {
	ctx := context.Background()
	RegisterServer(ctx, "lambda", NewLambdaServer)
}

type LambdaServer struct {
	Server
	url          *url.URL
	binary_types []string
}

func NewLambdaServer(ctx context.Context, uri string) (Server, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	server := LambdaServer{
		url: u,
	}

	q := u.Query()

	binary_types, ok := q["binary_type"]

	if ok {
		server.binary_types = binary_types
	}

	return &server, nil
}

func (s *LambdaServer) Address() string {
	return s.url.String()
}

func (s *LambdaServer) ListenAndServe(ctx context.Context, mux http.Handler) error {

	lambda_opts := new(algnhsa.Options)

	if len(s.binary_types) > 0 {
		lambda_opts.BinaryContentTypes = s.binary_types
	}

	algnhsa.ListenAndServe(mux, lambda_opts)
	return nil
}
