# go-http-server

## Example

_Error hanldling has been removed for the sake of brevity._

### Using a server

```
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aaronland/go-http-server"
	"net/http"
)

func NewHandler() http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {
		msg := fmt.Sprintf("Hello, %s", req.Host)
		rsp.Write([]byte(msg))
	}

	h := http.HandlerFunc(fn)
	return h
}

func main() {

	server_uri := flag.String("server-uri", "http://localhost:8080", "...")

	flag.Parse()

	ctx := context.Background()

	s, _ := server.NewServer(ctx, *server_uri)

	mux := http.NewServeMux()
	mux.Handle("/", NewHandler())

	log.Printf("Listening on %s", s.Address())
	s.ListenAndServe(ctx, mux)
}

```

### Writing a server

```
package server

import (
	"context"
	"net/http"
	"net/url"
)

func init() {
	ctx := context.Background()
	RegisterServer(ctx, "http", NewHTTPServer)
}

type HTTPServer struct {
	Server
	url *url.URL
}

func NewHTTPServer(ctx context.Context, uri string) (Server, error) {

	u, _ := url.Parse(uri)

	u.Scheme = "http"

	server := HTTPServer{
		url: u,
	}

	return &server, nil
}

func (s *HTTPServer) Address() string {
	return s.url.String()
}

func (s *HTTPServer) ListenAndServe(ctx context.Context, mux *http.ServeMux) error {
	return http.ListenAndServe(s.url.Host, mux)
}
```

## Interfaces

### Server

```
type Server interface {
	ListenAndServe(context.Context, *http.ServeMux) error
	Address() string
}

type ServerInitializeFunc func(context.Context, string) (Server, error)
```

## Server schemes

The following schemes/implementations are included by default with this package.

### http://{HOST}

A standard, plain-vanilla, HTTP server.

### https://{HOST}?cert={TLS_CERTIFICATE}&key={TLS_KEY}

This is an alias to the `tls://` scheme.

### lambda://

An AWS-compatible HTTP server.

### mkcert://{HOST}

A thin wrapper to invoke the [mkcert](https://github.com/FiloSottile/mkcert) tool to generate locally signed TLS certificate and key files. Once created this implementation will invoke the `tls://` scheme with the files create by `mkcert`. It is hoped this will be a short-lived scheme but it is necessary in the absence of an [ACME](https://github.com/go-acme/lego) compatibility with the `mkcert` tool.

### tls://{HOST}?cert={TLS_CERTIFICATE}&key={TLS_KEY}

A standard, plain-vanilla, HTTPS/TLS server. You must provide TLS certificate and key files.

## See also

* https://github.com/FiloSottile/mkcert