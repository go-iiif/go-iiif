package html

import (
	"embed"
)

//go:embed *.html *.css *.js
var FS embed.FS
