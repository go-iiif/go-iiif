package example

import (
	"embed"
)

//go:embed index.html *.css *.js
var FS embed.FS
