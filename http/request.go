package http

import (
	"fmt"
	"log/slog"
	gohttp "net/http"
)

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
	quality := req.PathValue("quality")
	format := req.PathValue("format")

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
