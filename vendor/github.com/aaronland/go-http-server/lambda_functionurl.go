package server

// https://github.com/aws/aws-lambda-go/blob/main/events/README_Lambda.md

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	ctx := context.Background()
	RegisterServer(ctx, "functionurl", NewLambdaFunctionURLServer)
}

// LambdaFunctionURLServer implements the `Server` interface for a use in a AWS LambdaFunctionURL + API Gateway context.
type LambdaFunctionURLServer struct {
	Server
	handler http.Handler
	binaryContentTypes map[string]bool
}

// NewLambdaFunctionURLServer returns a new `LambdaFunctionURLServer` instance configured by 'uri' which is
// expected to be defined in the form of:
//
//	functionurl://?{PARAMETERS}
//
// Valid parameters are:
// * `binary_type={MIMETYPE}` One or more mimetypes to be served by AWS FunctionURLs as binary content types.
func NewLambdaFunctionURLServer(ctx context.Context, uri string) (Server, error) {
	
	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	q := u.Query()

	binary_types := make(map[string]bool)
		
	for _, t := range q["binary_type"] {
		binary_types[t] = true
	}
	
	server := LambdaFunctionURLServer{
		binaryContentTypes: binary_types,
	}
	
	return &server, nil
}

// Address returns the fully-qualified URL used to instantiate 's'.
func (s *LambdaFunctionURLServer) Address() string {
	return "functionurl://"
}

// ListenAndServe starts the serve and listens for requests using 'mux' for routing.
func (s *LambdaFunctionURLServer) ListenAndServe(ctx context.Context, mux http.Handler) error {
	s.handler = mux
	lambda.Start(s.handleRequest)
	return nil
}

func (s *LambdaFunctionURLServer) handleRequest(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	req, err := newHTTPRequest(ctx, request)

	if err != nil {
		return events.LambdaFunctionURLResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	rec := httptest.NewRecorder()
	s.handler.ServeHTTP(rec, req)

	rsp := rec.Result()

	event_rsp_headers := make(map[string]string)

	for k, v := range rsp.Header {
		event_rsp_headers[k] = strings.Join(v, ",")
	}

	event_rsp := events.LambdaFunctionURLResponse{
		StatusCode: rsp.StatusCode,
		Headers: event_rsp_headers,
	}

	content_type := rsp.Header.Get("Content-Type")

	if s.binaryContentTypes[content_type] {
		event_rsp.Body = base64.StdEncoding.EncodeToString(rec.Body.Bytes())
		event_rsp.IsBase64Encoded = true
	} else {
		event_rsp.Body = rec.Body.String()
	}
	
	return event_rsp, nil
}

// This was clone and modified as necessary from https://github.com/akrylysov/algnhsa/blob/master/request.go#L30
// so there may still be issues.

// https://docs.aws.amazon.com/lambda/latest/dg/urls-invocation.html

func newHTTPRequest(ctx context.Context, event events.LambdaFunctionURLRequest) (*http.Request, error) {

	// https://pkg.go.dev/github.com/aws/aws-lambda-go/events#LambdaFunctionURLRequest
	// https://pkg.go.dev/github.com/aws/aws-lambda-go/events#LambdaFunctionURLRequestContextHTTPDescription

	rawQuery := event.RawQueryString

	if len(rawQuery) == 0 {

		params := url.Values{}

		for k, v := range event.QueryStringParameters {
			params.Set(k, v)
		}

		rawQuery = params.Encode()
	}

	headers := make(http.Header)

	for k, v := range event.Headers {
		headers.Set(k, v)
	}

	unescapedPath, err := url.PathUnescape(event.RawPath)

	if err != nil {
		return nil, err
	}
	u := url.URL{
		Host:     headers.Get("Host"),
		Path:     unescapedPath,
		RawQuery: rawQuery,
	}

	// Handle base64 encoded body.

	var body io.Reader = strings.NewReader(event.Body)

	if event.IsBase64Encoded {
		body = base64.NewDecoder(base64.StdEncoding, body)
	}

	req_context := event.RequestContext

	r, err := http.NewRequestWithContext(ctx, req_context.HTTP.Method, u.String(), body)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new HTTP request, %w", err)
	}

	r.RemoteAddr = req_context.HTTP.SourceIP
	r.RequestURI = u.RequestURI()

	r.Header = headers
	return r, nil
}
