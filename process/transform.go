package process

import (
	iiifcache "github.com/thisisaaronland/go-iiif/cache"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	iiifimage "github.com/thisisaaronland/go-iiif/image"
	iiiflevel "github.com/thisisaaronland/go-iiif/level"
)

func TransformURIWithInstructions(uri string, i IIIFInstructions, config *iiifconfig.Config, source_cache iiifcache.Cache, dest_cache iiifcache.Cache) (string, error) {

	level, err := iiiflevel.NewLevelFromConfig(config, "http://localhost")

	if err != nil {
		return "", err
	}

	transformation, err := iiifimage.NewTransformation(level, i.Region, i.Size, i.Rotation, i.Quality, i.Format)

	if err != nil {
		return "", err
	}

	new_uri, err := transformation.ToURI(uri)

	if err != nil {
		return "", err
	}

	// BUT WHAT IF WE WANT TO FORCE THE TRANSFORM...
	
	_, err = dest_cache.Get(new_uri)

	if err == nil {
		return new_uri, nil
	}

	im, err := iiifimage.NewImageFromConfigWithCache(config, source_cache, uri)

	if err != nil {
		return "", err
	}
	
	err = im.Transform(transformation)

	if err != nil {
		return "", err
	}

	err = dest_cache.Set(new_uri, im.Body())

	if err != nil {
		return "", err
	}

	return new_uri, nil
}
