package http

import (
	"expvar"
	"fmt"
	gohttp "net/http"
	"sync"
)

var cacheHit *expvar.Int
var cacheMiss *expvar.Int
var cacheSet *expvar.Int

var transformsCount *expvar.Int
var transformsAvgTime *expvar.Float

var transforms_counter int64
var transforms_timer int64

var timers_mu *sync.Mutex

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

func EndpointFromRequest(r *gohttp.Request) string {

	scheme := "http"

	if r.TLS != nil {
		scheme = "https"
	}

	endpoint := fmt.Sprintf("%s://%s", scheme, r.Host)
	return endpoint
}
