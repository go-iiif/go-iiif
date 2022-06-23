package server

// https://medium.com/@simonfrey/go-as-in-golang-standard-net-http-config-will-break-your-production-environment-1360871cb72b

// https://ieftimov.com/post/make-resilient-golang-net-http-servers-using-timeouts-deadlines-context-cancellation/

// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
// https://blog.cloudflare.com/exposing-go-on-the-internet/

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func init() {
	ctx := context.Background()
	RegisterServer(ctx, "http", NewHTTPServer)
	RegisterServer(ctx, "https", NewHTTPServer)	
}

// HTTPServer implements the `Server` interface for a basic `net/http` server.
type HTTPServer struct {
	Server
	url         *url.URL
	http_server *http.Server
	cert        string
	key         string
}

// NewHTTPServer returns a new `HTTPServer` instance configured by 'uri' which is
// expected to be defined in the form of:
//
//	{SCHEME}://{ADDRESS}:{PORT}?{PARAMETERS}
//
// Where {SCHEME} is either 'http' or 'https'; {ADDRESS} and {PORT} are the address
// and port to listen for requests on. Valid parameters are:
// * `tls_cert={CERTIFICATE}` The path for a TLS certificate to use; required if {SCHEME} is 'https'.
// * `tls_key={KEY}` The path for a TLS key to use; required if {SCHEME} is 'https'
// * `read_timeout={SECONDS}` A custom setting for HTTP read timeouts. Default is 2 seconds.
// * `write_timeout={SECONDS}` A custom setting for HTTP write timeouts. Default is 10 seconds.
// * `idle_timeout={SECONDS}` A custom setting for HTTP idle timeouts. Default is 15 seconds.
// * `header_timeout={SECONDS}` A custom setting for HTTP header timeouts. Default is 2 seconds.
func NewHTTPServer(ctx context.Context, uri string) (Server, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	u.Scheme = "http"

	read_timeout := 2 * time.Second
	write_timeout := 10 * time.Second
	idle_timeout := 15 * time.Second
	header_timeout := 2 * time.Second

	q := u.Query()

	if q.Get("read_timeout") != "" {

		to, err := strconv.Atoi(q.Get("read_timeout"))

		if err != nil {
			return nil, err
		}

		read_timeout = time.Duration(to) * time.Second
	}

	if q.Get("write_timeout") != "" {

		to, err := strconv.Atoi(q.Get("write_timeout"))

		if err != nil {
			return nil, err
		}

		write_timeout = time.Duration(to) * time.Second
	}

	if q.Get("idle_timeout") != "" {

		to, err := strconv.Atoi(q.Get("idle_timeout"))

		if err != nil {
			return nil, err
		}

		idle_timeout = time.Duration(to) * time.Second
	}

	if q.Get("header_timeout") != "" {

		to, err := strconv.Atoi(q.Get("header_timeout"))

		if err != nil {
			return nil, err
		}

		header_timeout = time.Duration(to) * time.Second
	}

	tls_cert := q.Get("cert")
	tls_key := q.Get("key")

	if (tls_cert != "") && (tls_key != "") {

		_, err = os.Stat(tls_cert)

		if err != nil {
			return nil, err
		}

		_, err = os.Stat(tls_key)

		if err != nil {
			return nil, err
		}

		u.Scheme = "https"

	} else if (tls_cert != "") && (tls_key == "") {
		return nil, errors.New("Missing TLS key parameter")
	} else if (tls_key != "") && (tls_key == "") {
		return nil, errors.New("Missing TLS cert parameter")
	} else {
		// pass
	}

	srv := &http.Server{
		Addr:              u.Host,
		ReadTimeout:       read_timeout,
		WriteTimeout:      write_timeout,
		IdleTimeout:       idle_timeout,
		ReadHeaderTimeout: header_timeout,
	}

	server := HTTPServer{
		url:         u,
		http_server: srv,
		cert:        tls_cert,
		key:         tls_key,
	}

	return &server, nil
}

// Address returns the fully-qualified URI where the server instance can be contacted.
func (s *HTTPServer) Address() string {

	u, _ := url.Parse(s.url.String())
	u.RawQuery = ""

	return u.String()
}

// ListenAndServe starts the server and listens for requests using 'mux' for routing.
func (s *HTTPServer) ListenAndServe(ctx context.Context, mux http.Handler) error {

	idleConnsClosed := make(chan struct{})

	go func() {

		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.

		err := s.http_server.Shutdown(context.Background())

		if err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}

		close(idleConnsClosed)
	}()

	s.http_server.Handler = mux

	var err error

	if s.cert != "" && s.key != "" {
		err = s.http_server.ListenAndServeTLS(s.cert, s.key)
	} else {
		err = s.http_server.ListenAndServe()
	}

	if err != nil && err != http.ErrServerClosed {
		return err
	}

	<-idleConnsClosed
	return nil
}
