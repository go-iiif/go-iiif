package http

import (
	"expvar"
	"fmt"
	"log/slog"
	gohttp "net/http"
	gourl "net/url"
	"path/filepath"
	"strings"
	"sync"

	"github.com/whosonfirst/go-sanitize"
)

var cacheHit *expvar.Int
var cacheMiss *expvar.Int
var cacheSet *expvar.Int

var transformsCount *expvar.Int
var transformsAvgTime *expvar.Float

var transforms_counter int64
var transforms_timer int64

var timers_mu *sync.Mutex

var sanitize_opts = sanitize.DefaultOptions()

type IIIFParameters struct {
	Identifier string
	Region     string
	Size       string
	Rotation   string
	Quality    string
	Format     string
}

func init() {

	cacheHit = expvar.NewInt("CacheHit")
	cacheMiss = expvar.NewInt("CacheMiss")
	cacheSet = expvar.NewInt("CacheSet")

	transformsCount = expvar.NewInt("TransformsCount")
	transformsAvgTime = expvar.NewFloat("TransformsAvgTimeMS")

	transforms_counter = 0
	transforms_timer = 0

	timers_mu = new(sync.Mutex)
}

func LoggerWithRequest(req *gohttp.Request, logger *slog.Logger) *slog.Logger {

	if logger == nil {
		logger = slog.Default()
	}

	logger = logger.With("method", req.Method)
	// logger = logger.With("user agent", req.Header.Get("User-Agent"))
	logger = logger.With("path", req.URL.Path)
	logger = logger.With("remote addr", req.RemoteAddr)

	return logger
}

func EndpointFromRequest(req *gohttp.Request) string {

	scheme := "http"

	if req.TLS != nil {
		scheme = "https"
	}

	endpoint := fmt.Sprintf("%s://%s", scheme, req.Host)
	return endpoint
}

func GetIIIFParameter(req *gohttp.Request, key string) (string, error) {

	raw := req.PathValue(key)

	value, err := sanitize.SanitizeString(raw, sanitize_opts)

	if err != nil {
		return "", err
	}

	value, err = gourl.QueryUnescape(value)

	if err != nil {
		return "", err
	}

	// This should be already be stripped out by the time we get here but just
	// in case... (20160926/thisisaaronland)

	if strings.Contains(value, "../") {
		return "", fmt.Errorf("Invalid key %s", key)
	}

	return value, nil
}

func GetIIIFParameters(req *gohttp.Request) (*IIIFParameters, error) {

	id, err := GetIIIFParameter(req, "identifier")

	if err != nil {
		return nil, err
	}

	region, err := GetIIIFParameter(req, "region")

	if err != nil {
		return nil, err
	}

	size, err := GetIIIFParameter(req, "size")

	if err != nil {
		return nil, err
	}

	rotation, err := GetIIIFParameter(req, "rotation")

	if err != nil {
		return nil, err
	}

	// START OF UNFORTUNATE but net/http URL pattern-mux stuff doesn't allow {foo}.{bar}

	quality_dot_format, err := GetIIIFParameter(req, "quality_dot_format")

	if err != nil {
		return nil, err
	}

	dot_format := filepath.Ext(quality_dot_format)
	format := strings.TrimLeft(dot_format, ".")

	quality := strings.Replace(quality_dot_format, dot_format, "", 1)

	// END OF UNFORTUNATE but net/http URL pattern-mux stuff doesn't allow {foo}.{bar}

	params := IIIFParameters{
		Identifier: id,
		Region:     region,
		Size:       size,
		Rotation:   rotation,
		Quality:    quality,
		Format:     format,
	}

	return &params, nil
}
