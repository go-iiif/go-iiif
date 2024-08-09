package funcs

import (
	"github.com/tidwall/gjson"
)

func GjsonGet(body string, path string) interface{} {
	rsp := gjson.Get(body, path)
	return rsp.Value()
}
