package image

import (
	"errors"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	iiifsource "github.com/thisisaaronland/go-iiif/source"
)

type Image interface {
	Identifier() string
	Transform(*Transformation) error
	Body() []byte
	Format() string
	ContentType() string
	Dimensions() (Dimensions, error)
}

type Dimensions interface {
	Height() int
	Width() int
}

func NewImageFromConfig(config *iiifconfig.Config, id string) (Image, error) {

	source, err := iiifsource.NewSourceFromConfig(config.Images)

	if err != nil {
		return nil, err
	}

	return NewImageFromConfigWithSource(config, source, id)
}

func NewImageFromConfigWithSource(config *iiifconfig.Config, source iiifsource.Source, id string) (Image, error) {

	if config.Graphics.Source.Name == "VIPS" {
		return NewVIPSImageFromSource(source, id)
	}

	return nil, errors.New("Unknown graphics source")
}
