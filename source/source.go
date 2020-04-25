package source

import (
	"errors"
	iiifconfig "github.com/go-iiif/go-iiif/v4/config"
	"strings"
)

type Source interface {
	Read(uri string) ([]byte, error)
}

func NewSourceFromConfig(config *iiifconfig.Config) (Source, error) {

	cfg := config.Images

	// note that there is no "Memory" source or at least not yet
	// since it assumes you're passing it []bytes and not a config
	// file (20160907/thisisaaronland)

	var source Source
	var err error

	switch strings.ToLower(cfg.Source.Name) {
	case "blob":
		source, err = NewBlobSource(config)
	case "disk":
		source, err = NewDiskSource(config)
	case "flickr":
		source, err = NewFlickrSource(config)
	case "s3":
		source, err = NewS3Source(config)
	case "s3blob":
		source, err = NewS3Source(config)
	case "uri":
		source, err = NewURISource(config)
	default:
		err = errors.New("Unknown source type")
	}

	return source, err
}
