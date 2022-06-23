package server

import (
	"context"
	"fmt"
	"github.com/akrylysov/algnhsa"
	_ "log"
	"net/http"
	"net/url"
)

func init() {
	ctx := context.Background()
	RegisterServer(ctx, "lambda", NewLambdaServer)
}

// LambdaServer implements the `Server` interface for a use in a AWS Lambda + API Gateway context.
type LambdaServer struct {
	Server
	url          *url.URL
	binary_types []string
}

// NewLambdaServer returns a new `LambdaServer` instance configured by 'uri' which is
// expected to be defined in the form of:
//
//	lambda://?{PARAMETERS}
//
// Valid parameters are:
// * `binary_type={MIMETYPE}` One or more mimetypes to be served by AWS API Gateway as binary content types.
func NewLambdaServer(ctx context.Context, uri string) (Server, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
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

// Address returns the fully-qualified URL used to instantiate 's'.
func (s *LambdaServer) Address() string {
	return s.url.String()
}

// ListenAndServe starts the serve and listens for requests using 'mux' for routing.
func (s *LambdaServer) ListenAndServe(ctx context.Context, mux http.Handler) error {

	lambda_opts := new(algnhsa.Options)

	if len(s.binary_types) > 0 {
		lambda_opts.BinaryContentTypes = s.binary_types
	}

	algnhsa.ListenAndServe(mux, lambda_opts)
	return nil
}
