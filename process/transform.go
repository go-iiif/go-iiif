package process

import (
	iiifcache "github.com/thisisaaronland/go-iiif/cache"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	iiifimage "github.com/thisisaaronland/go-iiif/image"
	iiiflevel "github.com/thisisaaronland/go-iiif/level"
)

func TransformURIWithInstructions(u URI, i IIIFInstructions, config *iiifconfig.Config, source_cache iiifcache.Cache, dest_cache iiifcache.Cache) (URI, iiifimage.Image, error) {

	level, err := iiiflevel.NewLevelFromConfig(config, "http://localhost")

	if err != nil {
		return nil, nil, err
	}

	transformation, err := iiifimage.NewTransformation(level, i.Region, i.Size, i.Rotation, i.Quality, i.Format)

	if err != nil {
		return nil, nil, err
	}

	new_uri, err := transformation.ToURI(u.URL())

	if err != nil {
		return nil, nil, err
	}

	/*

		cached_im, err := dest_cache.Get(new_uri)

		if err == nil {
			return new_uri, cached_im, nil
		}

	*/

	im, err := iiifimage.NewImageFromConfigWithCache(config, source_cache, u.URL())

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

	new_u, err := NewIIIFURI(new_uri)

	if err != nil {
		return nil, nil, err
	}

	return new_u, im, nil
}
