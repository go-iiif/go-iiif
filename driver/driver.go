package driver

import (
	iiifcache "github.com/go-iiif/go-iiif/cache"
	iiifconfig "github.com/go-iiif/go-iiif/config"
	iiifimage "github.com/go-iiif/go-iiif/image"
	iiifsource "github.com/go-iiif/go-iiif/source"
)

type Driver interface {
	NewImageFromConfigWithSource(*iiifconfig.Config, iiifsource.Source, string) (iiifimage.Image, error)
	NewImageFromConfigWithCache(*iiifconfig.Config, iiifcache.Cache, string) (iiifimage.Image, error)
	NewImageFromConfig(*iiifconfig.Config, string) (iiifimage.Image, error)
}
