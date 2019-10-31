package process

import (
	"fmt"
	iiifuri "github.com/go-iiif/go-iiif-uri"
	iiifcache "github.com/go-iiif/go-iiif/cache"
	iiifconfig "github.com/go-iiif/go-iiif/config"
	iiifdriver "github.com/go-iiif/go-iiif/driver"
	iiifimage "github.com/go-iiif/go-iiif/image"
	iiiflevel "github.com/go-iiif/go-iiif/level"
	_ "log"
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

	tr_uri, err := transformation.ToURI(target)

	if err != nil {
		return nil, nil, err
	}

	str_uri := fmt.Sprintf("%s:///%s", u.Driver(), tr_uri)

	new_uri, err := iiifuri.NewURI(str_uri)

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

	new_target := new_uri.Target()
	err = dest_cache.Set(new_target, im.Body())

	if err != nil {
		return nil, nil, err
	}

	return new_uri, im, nil
}
