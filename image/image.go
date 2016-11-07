package image

import (
	"errors"
	iiifcache "github.com/thisisaaronland/go-iiif/cache"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	iiifsource "github.com/thisisaaronland/go-iiif/source"
	_ "log"
	"os"
)

type Image interface {
	Identifier() string
	Rename(string) error
	Transform(*Transformation) error // http://iiif.io/api/image/2.1/#order-of-implementation
	Update([]byte) error
	Body() []byte
	Format() string
	ContentType() string
	Dimensions() (Dimensions, error)
}

type Dimensions interface {
	Height() int
	Width() int
}

func NewImageFromConfigWithCache(config *iiifconfig.Config, cache iiifcache.Cache, id string) (Image, error) {

	var image Image

	body, err := cache.Get(id)

	if err == nil {

		source, err := iiifsource.NewMemorySource(body)

		if err != nil {
			return nil, err
		}

		image, err = NewImageFromConfigWithSource(config, source, id)

		if err != nil {
			return nil, err
		}

	} else {

		image, err = NewImageFromConfig(config, id)

		if err != nil {
			return nil, err
		}

		go func() {
			cache.Set(id, image.Body())
		}()
	}

	return image, nil

}

func NewImageFromConfig(config *iiifconfig.Config, id string) (Image, error) {

	source, err := iiifsource.NewSourceFromConfig(config)

	if err != nil {
		return nil, err
	}

	return NewImageFromConfigWithSource(config, source, id)
}

func NewImageFromConfigWithSource(config *iiifconfig.Config, source iiifsource.Source, id string) (Image, error) {

	if config.Graphics.Source.Name == "VIPS" {

		/*
			http://www.vips.ecs.soton.ac.uk/supported/7.42/doc/html/libvips/VipsImage.html#vips-image-new-temp-file
		*/

		if config.Graphics.Source.Tmpdir != "" {

			tmpdir := config.Graphics.Source.Tmpdir

			_, err := os.Stat(tmpdir)

			if os.IsNotExist(err) {
				return nil, err
			}

			os.Setenv("TMPDIR", tmpdir)
		}

		return NewVIPSImageFromConfigWithSource(config, source, id)
	}

	return nil, errors.New("Unknown graphics source")
}
