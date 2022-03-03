package process

import (
	"context"
	"fmt"
	iiifuri "github.com/go-iiif/go-iiif-uri"
	iiifcache "github.com/go-iiif/go-iiif/v5/cache"
	iiifconfig "github.com/go-iiif/go-iiif/v5/config"
	iiifdriver "github.com/go-iiif/go-iiif/v5/driver"
	iiifimage "github.com/go-iiif/go-iiif/v5/image"
	iiiflevel "github.com/go-iiif/go-iiif/v5/level"
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

	compliance := level.Compliance()

	transformation, err := iiifimage.NewTransformation(compliance, i.Region, i.Size, i.Rotation, i.Quality, i.Format)

	if err != nil {
		return nil, nil, err
	}

	// I do not love this...

	switch u.Scheme() {

	case "rewrite":
		// pass
	default:

		tr_uri, err := transformation.ToURI(target)

		if err != nil {
			return nil, nil, err
		}

		ctx := context.Background()

		str_uri := fmt.Sprintf("%s://%s", iiifuri.FILE_SCHEME, tr_uri)
		new_uri, err := iiifuri.NewURI(ctx, str_uri)

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

	ctx := context.Background()

	str_uri := fmt.Sprintf("%s://%s", iiifuri.FILE_SCHEME, target)
	new_uri, err := iiifuri.NewURI(ctx, str_uri)

	if err != nil {
		return nil, nil, err
	}

	return new_uri, im, nil
}
