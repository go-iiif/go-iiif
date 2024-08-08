package defaults

import (
	"embed"
)

//go:embed *.json
var FS embed.FS

const URI string = "defaults://"
