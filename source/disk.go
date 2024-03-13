package source

import (
	"fmt"
	_ "log"

	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
)

func NewDiskSource(config *iiifconfig.Config) (Source, error) {

	cfg := config.Images
	uri := fmt.Sprintf("file://%s", cfg.Source.Path)

	return NewBlobSourceFromURI(uri)
}
