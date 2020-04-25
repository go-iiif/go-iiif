package process

import (
	iiifuri "github.com/go-iiif/go-iiif-uri"
	iiifcache "github.com/go-iiif/go-iiif/v4/cache"
	iiifconfig "github.com/go-iiif/go-iiif/v4/config"
	iiifdriver "github.com/go-iiif/go-iiif/v4/driver"
	iiifimage "github.com/go-iiif/go-iiif/v4/image"
)

type IIIFProcessor struct {
	Processor
	config       *iiifconfig.Config
	driver       iiifdriver.Driver
	source_cache iiifcache.Cache
	dest_cache   iiifcache.Cache
}

func NewIIIFProcessor(config *iiifconfig.Config, driver iiifdriver.Driver) (Processor, error) {
	return NewIIIFProcessorWithCaches(config, driver, nil, nil)
}

func NewIIIFProcessorWithCaches(config *iiifconfig.Config, driver iiifdriver.Driver, source_cache iiifcache.Cache, dest_cache iiifcache.Cache) (Processor, error) {

	if source_cache == nil {

		c, err := iiifcache.NewImagesCacheFromConfig(config)

		if err != nil {
			return nil, err
		}

		source_cache = c
	}

	if dest_cache == nil {

		c, err := iiifcache.NewDerivativesCacheFromConfig(config)

		if err != nil {
			return nil, err
		}

		dest_cache = c
	}

	pr := IIIFProcessor{
		config:       config,
		driver:       driver,
		source_cache: source_cache,
		dest_cache:   dest_cache,
	}

	return &pr, nil
}

func (pr *IIIFProcessor) ProcessURIWithInstructions(u iiifuri.URI, label Label, i IIIFInstructions) (iiifuri.URI, iiifimage.Image, error) {

	return TransformURIWithInstructions(u, i, pr.config, pr.driver, pr.source_cache, pr.dest_cache)
}
