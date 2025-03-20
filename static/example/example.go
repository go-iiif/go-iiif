package example

import (
	"embed"
)

//go:embed tiled/* images/* cache/*
var FS embed.FS
