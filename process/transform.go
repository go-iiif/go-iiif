package process

import (
	"fmt"
	iiifuri "github.com/go-iiif/go-iiif-uri"
	iiifcache "github.com/go-iiif/go-iiif/cache"
	iiifconfig "github.com/go-iiif/go-iiif/config"
	iiifdriver "github.com/go-iiif/go-iiif/driver"
	iiifimage "github.com/go-iiif/go-iiif/image"
	iiiflevel "github.com/go-iiif/go-iiif/level"
)

func TransformURIWithInstructions(u iiifuri.URI, i IIIFInstructions, config *iiifconfig.Config, driver iiifdriver.Driver, source_cache iiifcache.Cache, dest_cache iiifcache.Cache) (iiifuri.URI, iiifimage.Image, error) {

	origin := u.Origin()
	target := u.Target()

	level, err := iiiflevel.NewLevelFromConfig(config, "http://localhost")

	if err != nil {
		return nil, nil, err
	}

	transformation, err := iiifimage.NewTransformation(level, i.Region, i.Size, i.Rotation, i.Quality, i.Format)

	if err != nil {
		return nil, nil, err
	}

	new_uri, err := transformation.ToURI(target)

	if err != nil {
		return nil, nil, err
	}

	im, err := driver.NewImageFromConfigWithCache(config, source_cache, origin)

	if err != nil {
		return nil, nil, err
	}

	err = im.Transform(transformation)

	if err != nil {
		return nil, nil, err
	}

	err = dest_cache.Set(new_uri, im.Body())

	if err != nil {
		return nil, nil, err
	}

	new_uri = fmt.Sprintf("%s://", u.Driver(), new_uri)

	new_u, err := iiifuri.NewURI(new_uri)

	if err != nil {
		return nil, nil, err
	}

	return new_u, im, nil
}
