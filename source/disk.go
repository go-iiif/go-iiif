package source

import (
	"fmt"
	iiifconfig "github.com/go-iiif/go-iiif/v3/config"
	_ "log"
)

func NewDiskSource(config *iiifconfig.Config) (Source, error) {

	cfg := config.Images
	uri := fmt.Sprintf("file://%s", cfg.Source.Path)

	return NewBlobSourceFromURI(uri)
}
