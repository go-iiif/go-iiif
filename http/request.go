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

	params := &IIIFParameters{}

	// DO THIS HERE

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
	return logger
}
