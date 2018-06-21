package http

import (
	"fmt"
	gohttp "net/http"
	"path/filepath"
	"strings"
)

func ExampleHandler(root string) (gohttp.HandlerFunc, error) {

	base := fmt.Sprintf("/%s", filepath.Base(root))

	fs := gohttp.FileServer(gohttp.Dir(root))

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {
		req.URL.Path = strings.Replace(req.URL.Path, base, "", 1)
		fs.ServeHTTP(rsp, req)
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
