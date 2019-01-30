package process

import (
	iiifcache "github.com/thisisaaronland/go-iiif/cache"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	iiifimage "github.com/thisisaaronland/go-iiif/image"
)

type IIIFURI struct {
	URI
	raw string
}

func NewIIIFURI(raw string) (URI, error) {

	u := IIIFURI{
		raw: raw,
	}

	return &u, nil
}

func (u *IIIFURI) URL() string {
	return u.raw
}

func (u *IIIFURI) String() string {
	return u.raw
}

type IIIFProcessor struct {
	Processor
	config       *iiifconfig.Config
	source_cache iiifcache.Cache
	dest_cache   iiifcache.Cache
}

func NewIIIFProcessor(config *iiifconfig.Config) (Processor, error) {
	return NewIIIFProcessorWithCaches(config, nil, nil)
}

func NewIIIFProcessorWithCaches(config *iiifconfig.Config, source_cache iiifcache.Cache, dest_cache iiifcache.Cache) (Processor, error) {

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
		source_cache: source_cache,
		dest_cache:   dest_cache,
	}

	return &pr, nil
}

func (pr *IIIFProcessor) ProcessURIWithInstructions(u URI, label Label, i IIIFInstructions) (URI, iiifimage.Image, error) {

	return TransformURIWithInstructions(u, i, pr.config, pr.source_cache, pr.dest_cache)
}
