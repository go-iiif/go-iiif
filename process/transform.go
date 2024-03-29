package process

import (
	"context"
	"fmt"

	iiifuri "github.com/go-iiif/go-iiif-uri"
	iiifcache "github.com/go-iiif/go-iiif/v6/cache"
	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
	iiifdriver "github.com/go-iiif/go-iiif/v6/driver"
	iiifimage "github.com/go-iiif/go-iiif/v6/image"
	iiiflevel "github.com/go-iiif/go-iiif/v6/level"
)

func TransformURIWithInstructions(u iiifuri.URI, i IIIFInstructions, config *iiifconfig.Config, driver iiifdriver.Driver, source_cache iiifcache.Cache, dest_cache iiifcache.Cache) (iiifuri.URI, iiifimage.Image, error) {

	origin := u.Origin()
	target, err := u.Target(nil)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to derive target for origin '%s', %w", origin, err)
	}

	level, err := iiiflevel.NewLevelFromConfig(config, "http://localhost")

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to derive level from config for origin '%s', %w", origin, err)
	}

	compliance := level.Compliance()

	transformation, err := iiifimage.NewTransformation(compliance, i.Region, i.Size, i.Rotation, i.Quality, i.Format)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to create new transformation for origin '%s', %w", origin, err)
	}

	// I do not love this...

	switch u.Scheme() {

	case "rewrite":
		// pass
	default:

		tr_uri, err := transformation.ToURI(target)

		if err != nil {
			return nil, nil, fmt.Errorf("Failed to create transformation URI for origin '%s' with target '%s', %w", origin, target, err)
		}

		ctx := context.Background()

		str_uri := fmt.Sprintf("%s://%s", iiifuri.FILE_SCHEME, tr_uri)
		new_uri, err := iiifuri.NewURI(ctx, str_uri)

		if err != nil {
			return nil, nil, fmt.Errorf("Failed to create new IIIF URI for origin '%s' from string '%s', %w", origin, str_uri, err)
		}

		new_target, err := new_uri.Target(nil)

		if err != nil {
			return nil, nil, fmt.Errorf("Failed to create target for origin '%s' from URI '%s', %w", origin, new_uri, err)
		}

		target = new_target
	}

	im, err := driver.NewImageFromConfigWithCache(config, source_cache, origin)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to create new image for origin '%s', %w", origin, err)
	}

	err = im.Transform(transformation)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to perform transformation for origin '%s', %w", origin, err)
	}

	err = dest_cache.Set(target, im.Body())

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to store new derivative for origin '%s' at '%s', %w", origin, target, err)
	}

	ctx := context.Background()

	str_uri := fmt.Sprintf("%s://%s", iiifuri.FILE_SCHEME, target)
	new_uri, err := iiifuri.NewURI(ctx, str_uri)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to derive new IIIF URI for origin '%s' from string '%s', %w", origin, str_uri, err)
	}

	return new_uri, im, nil
}
