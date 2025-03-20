package example

import (
	"embed"
)

//go:embed *.js *.html *.css images/* cache/*
var FS embed.FS
