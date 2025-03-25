package process

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	iiifuri "github.com/go-iiif/go-iiif-uri"
	iiifcache "github.com/go-iiif/go-iiif/v7/cache"
	iiifconfig "github.com/go-iiif/go-iiif/v7/config"
	iiifdriver "github.com/go-iiif/go-iiif/v7/driver"
	iiifimage "github.com/go-iiif/go-iiif/v7/image"
	iiiflevel "github.com/go-iiif/go-iiif/v7/level"
)

func TransformURIWithInstructions(ctx context.Context, u iiifuri.URI, i IIIFInstructions, config *iiifconfig.Config, driver iiifdriver.Driver, source_cache iiifcache.Cache, dest_cache iiifcache.Cache) (iiifuri.URI, iiifimage.Image, error) {

	origin := u.Origin()
	target, err := u.Target(nil)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to derive target for origin '%s', %w", origin, err)
	}

	logger := slog.Default()
	logger = logger.With("uri", u)
	logger = logger.With("origin", origin)
	logger = logger.With("target", target)
	logger = logger.With("source cache", source_cache)
	logger = logger.With("destination cache", dest_cache)

	logger.Debug("Transform with instructions")

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
		// logger.Warn("Origin has rewrite scheme, do not rewrite")
	default:

		tr_uri, err := transformation.ToURI(target)

		if err != nil {
			return nil, nil, fmt.Errorf("Failed to create transformation URI for origin '%s' with target '%s', %w", origin, target, err)
		}

		logger.Debug("Transformation URI", "uri", tr_uri)

		str_uri := fmt.Sprintf("%s://%s", iiifuri.FILE_SCHEME, tr_uri)
		new_uri, err := iiifuri.NewURI(ctx, str_uri)

		if err != nil {
			return nil, nil, fmt.Errorf("Failed to create new IIIF URI for origin '%s' from string '%s', %w", origin, str_uri, err)
		}

		new_target, err := new_uri.Target(nil)

		if err != nil {
			return nil, nil, fmt.Errorf("Failed to create target for origin '%s' from URI '%s', %w", origin, new_uri, err)
		}

		logger.Debug("Rewrite target", "new target", new_target)
		target = new_target

		logger = logger.With("new target", new_target)
	}

	logger.Debug("Derive image from source cache", "source cache", source_cache)
	im, err := driver.NewImageFromConfigWithCache(ctx, config, source_cache, origin)

	if err != nil {
		logger.Debug("Failed to derive image from source cache", "error", err)
		return nil, nil, fmt.Errorf("Failed to create new image for origin '%s', %w", origin, err)
	}

	err = im.Transform(transformation)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to perform transformation for origin '%s', %w", origin, err)
	}

	logger.Debug("Write target to destination cache", "destination cache", dest_cache)
	err = dest_cache.Set(target, im.Body())

	if err != nil {
		logger.Debug("Failed to write target to destination cache", "error", err)
		return nil, nil, fmt.Errorf("Failed to store new derivative for origin '%s' at '%s', %w", origin, target, err)
	}

	uri_target := target

	if !strings.HasPrefix(uri_target, "/") {
		uri_target = fmt.Sprintf("/%s", target)
	}

	str_uri := fmt.Sprintf("%s://%s", iiifuri.FILE_SCHEME, uri_target)
	new_uri, err := iiifuri.NewURI(ctx, str_uri)

	if err != nil {
		return nil, nil, fmt.Errorf("Failed to derive new IIIF URI for origin '%s' from string '%s', %w", origin, str_uri, err)
	}

	logger.Debug("Return transformation", "new uri", new_uri)
	return new_uri, im, nil
}
