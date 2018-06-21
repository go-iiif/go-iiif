package http

import (
	"expvar"
	"fmt"
	"log"
	gohttp "net/http"
	"strings"
)

func ExpvarHandler(host string) (gohttp.HandlerFunc, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		remote := strings.Split(req.RemoteAddr, ":")

		if remote[0] != "127.0.0.1" && remote[0] != host {

			log.Printf("host '%s' remote '%s'\n", remote[0], host)
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
