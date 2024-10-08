package process

import (
	"log/slog"

	iiifuri "github.com/go-iiif/go-iiif-uri"
	iiifcache "github.com/go-iiif/go-iiif/v6/cache"
	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
	iiifdriver "github.com/go-iiif/go-iiif/v6/driver"
	iiifimage "github.com/go-iiif/go-iiif/v6/image"
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

	logger := slog.Default()
	logger.Debug("New IIIF processor with caches", "source", source_cache, "destination", dest_cache)

	if source_cache == nil {

		c, err := iiifcache.NewImagesCacheFromConfig(config)

		if err != nil {
			return nil, err
		}

		logger.Debug("Created new IIIF processor source cache", "source", c)
		source_cache = c
	}

	if dest_cache == nil {

		c, err := iiifcache.NewDerivativesCacheFromConfig(config)

		if err != nil {
			return nil, err
		}

		logger.Debug("Created new IIIF processor destination cache", "destination", c)
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
