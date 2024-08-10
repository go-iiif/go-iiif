package http

import (
	"fmt"
	"log/slog"
	gohttp "net/http"
	"regexp"
)

var re_quality_dot_format = regexp.MustCompile(`^([^\.]+)\.([^\.]+)$`)

type IIIFParameters struct {
	Identifier string
	Region     string
	Size       string
	Rotation   string
	Quality    string
	Format     string
}

func IIIFParamtersFromRequest(req *gohttp.Request) (*IIIFParameters, error) {

	id := req.PathValue("identifier")
	region := req.PathValue("region")
	size := req.PathValue("size")
	rotation := req.PathValue("rotation")

	// START OF because Go's net/http wildcard stuff
	// can't deal with {foo}.{bar} URIs

	quality_dot_format := req.PathValue("quality_dot_format")

	var quality string
	var format string

	if re_quality_dot_format.MatchString(quality_dot_format) {
		m := re_quality_dot_format.FindStringSubmatch(quality_dot_format)
		quality = m[1]
		format = m[2]
	}

	// END OF because Go's net/http wildcard stuff

	params := &IIIFParameters{
		Identifier: id,
		Region:     region,
		Size:       size,
		Rotation:   rotation,
		Quality:    quality,
		Format:     format,
	}

	return params, nil
}

func EndpointFromRequest(r *gohttp.Request) string {

	scheme := "http"

	if r.TLS != nil {
		scheme = "https"
	}

	endpoint := fmt.Sprintf("%s://%s", scheme, r.Host)
	return endpoint
}

func LoggerForRequest(req *gohttp.Request) *slog.Logger {
	logger := slog.Default()

	logger = logger.With("method", req.Method)
	logger = logger.With("path", req.URL.Path)
	logger = logger.With("remote addr", req.RemoteAddr)
	logger = logger.With("user ip", ReadUserIP(req))

	return logger
}

func ReadUserIP(req *gohttp.Request) string {

	addr := req.Header.Get("X-Real-Ip")

	if addr == "" {
		addr = req.Header.Get("X-Forwarded-For")
	}

	if addr == "" {
		addr = req.RemoteAddr
	}

	return addr
}
