package ping

import (
	"net/http"
)

func PingHandler() (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {
		rsp.Header().Set("Content-Type", "text/plain")
		rsp.Write([]byte("PONG"))
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
