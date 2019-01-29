package process

import (
	"errors"
	"fmt"
	iiifcache "github.com/thisisaaronland/go-iiif/cache"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	iiifimage "github.com/thisisaaronland/go-iiif/image"
	iiifservice "github.com/thisisaaronland/go-iiif/service"
	"log"
	"sync"
)

type ParallelProcessor struct {
	Processor
	config       *iiifconfig.Config
	source_cache iiifcache.Cache
	dest_cache   iiifcache.Cache
}

func NewParallelProcessor(config *iiifconfig.Config) (Processor, error) {

	return NewParallelProcessorWithCaches(config, nil, nil)
}

func NewParallelProcessorWithCaches(config *iiifconfig.Config, source_cache iiifcache.Cache, dest_cache iiifcache.Cache) (Processor, error) {

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

	pr := ParallelProcessor{
		config:       config,
		source_cache: source_cache,
		dest_cache:   dest_cache,
	}

	return &pr, nil
}

func (pr *ParallelProcessor) ProcessURIWithInstructionSet(uri string, instruction_set IIIFInstructionSet) (map[string]interface{}, error) {

	done_ch := make(chan bool)
	err_ch := make(chan error)

	remaining := len(instruction_set) + 1 // + 1 because we fetch the palette for the source image

	results := make(map[string]interface{})

	uris := make(map[string]string)
	dimensions := make(map[string][]int)

	mu := new(sync.RWMutex)

	go func() {

		defer func() {
			done_ch <- true
		}()

		im, err := iiifimage.NewImageFromConfig(pr.config, uri)

		if err != nil {
			msg := fmt.Sprintf("failed to derive palette for %s : %s", uri, err)
			err_ch <- errors.New(msg)
			return
		}

		for _, service_name := range pr.config.Profile.Services.Enable {

			if service_name == "palette" {

				s, err := iiifservice.NewPaletteService(pr.config.Palette, im)

				if err != nil {
					msg := fmt.Sprintf("failed to derive palette for %s : %s", uri, err)
					err_ch <- errors.New(msg)
					return
				}

				results["palette"] = s.Value()
				break
			}
		}

	}()

	for label, i := range instruction_set {

		i = EnsureInstructions(i)

		go func(uri string, label string, i IIIFInstructions) {

			defer func() {
				done_ch <- true
			}()

			new_uri, im, err := pr.ProcessURIWithInstructions(uri, i)

			if err != nil {
				msg := fmt.Sprintf("failed to process %s (%s) : %s", uri, label, err)
				err_ch <- errors.New(msg)
				return
			}

			dims, err := im.Dimensions()

			if err != nil {
				msg := fmt.Sprintf("failed to process %s (%s) : %s", uri, label, err)
				err_ch <- errors.New(msg)
				return
			}

			mu.Lock()

			uris[label] = new_uri

			dimensions[label] = []int{
				dims.Width(),
				dims.Height(),
			}

			mu.Unlock()

		}(uri, label, i)
	}

	for remaining > 0 {

		select {
		case <-done_ch:
			remaining -= 1
		case e := <-err_ch:
			log.Println(e)
		default:
			//
		}
	}

	results["uris"] = uris
	results["dimensions"] = dimensions

	return results, nil
}

func (pr *ParallelProcessor) ProcessURIWithInstructions(uri string, i IIIFInstructions) (string, iiifimage.Image, error) {

	return TransformURIWithInstructions(uri, i, pr.config, pr.source_cache, pr.dest_cache)
}
