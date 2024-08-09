package http

import (
	"expvar"
	"fmt"
	"log/slog"
	gohttp "net/http"
	"strings"
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
func ExpvarHandler(host string) (gohttp.HandlerFunc, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		remote := strings.Split(req.RemoteAddr, ":")

		if remote[0] != "127.0.0.1" && remote[0] != host {
			slog.Debug("Invalid host", "remote", remote[0], "allowed", host)
			gohttp.Error(rsp, "No soup for you!", gohttp.StatusForbidden)
			return
		}

		// This is copied wholesale from
		// https://golang.org/src/expvar/expvar.go

		rsp.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprintf(rsp, "{\n")

		first := true

		expvar.Do(func(kv expvar.KeyValue) {
			if !first {
				fmt.Fprintf(rsp, ",\n")
			}

			first = false
			fmt.Fprintf(rsp, "%q: %s", kv.Key, kv.Value)
		})

		fmt.Fprintf(rsp, "\n}\n")
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
