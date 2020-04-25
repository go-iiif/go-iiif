package process

import (
	iiifuri "github.com/go-iiif/go-iiif-uri"
	iiifcache "github.com/go-iiif/go-iiif/v4/cache"
	iiifconfig "github.com/go-iiif/go-iiif/v4/config"
	iiifdriver "github.com/go-iiif/go-iiif/v4/driver"
	iiifimage "github.com/go-iiif/go-iiif/v4/image"
	iiiflevel "github.com/go-iiif/go-iiif/v4/level"
	_ "log"
)

func TransformURIWithInstructions(u iiifuri.URI, i IIIFInstructions, config *iiifconfig.Config, driver iiifdriver.Driver, source_cache iiifcache.Cache, dest_cache iiifcache.Cache) (iiifuri.URI, iiifimage.Image, error) {

	origin := u.Origin()
	target, err := u.Target(nil)

	if err != nil {
		return nil, nil, err
	}

	level, err := iiiflevel.NewLevelFromConfig(config, "http://localhost")

	if err != nil {
		return nil, nil, err
	}

	transformation, err := iiifimage.NewTransformation(level, i.Region, i.Size, i.Rotation, i.Quality, i.Format)

	if err != nil {
		return nil, nil, err
	}

	// I do not love this...

	switch u.Driver() {

	case "rewrite":
		// pass
	default:

		tr_uri, err := transformation.ToURI(target)

		if err != nil {
			return nil, nil, err
		}

		str_uri := iiifuri.NewFileURIString(tr_uri)
		new_uri, err := iiifuri.NewURI(str_uri)

		if err != nil {
			return nil, nil, err
		}

		new_target, err := new_uri.Target(nil)

		if err != nil {
			return nil, nil, err
		}

		target = new_target
	}

	im, err := driver.NewImageFromConfigWithCache(config, source_cache, origin)

	if err != nil {
		return nil, nil, err
	}

	err = im.Transform(transformation)

	if err != nil {
		return nil, nil, err
	}

	err = dest_cache.Set(target, im.Body())

	if err != nil {
		return nil, nil, err
	}

	str_uri := iiifuri.NewFileURIString(target)
	new_uri, err := iiifuri.NewURI(str_uri)

	if err != nil {
		return nil, nil, err
	}

	return new_uri, im, nil
}
