package process

import (
	"errors"
	"fmt"
	iiifcache "github.com/thisisaaronland/go-iiif/cache"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
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

	source_cache, err := iiifcache.NewImagesCacheFromConfig(config)

	if err != nil {
		return nil, err
	}

	dest_cache, err := iiifcache.NewDerivativesCacheFromConfig(config)

	if err != nil {
		return nil, err
	}

	pr := ParallelProcessor{
		config:       config,
		source_cache: source_cache,
		dest_cache:   dest_cache,
	}

	return &pr, nil
}

func (pr *ParallelProcessor) ProcessURIWithInstructionSet(uri string, instruction_set IIIFInstructionSet) (map[string]string, error) {

	done_ch := make(chan bool)
	err_ch := make(chan error)

	remaining := len(instruction_set)

	results := make(map[string]string)
	mu := new(sync.RWMutex)

	for label, i := range instruction_set {

		i = EnsureInstructions(i)

		go func(uri string, label string, i IIIFInstructions) {

			defer func() {
				done_ch <- true
			}()

			new_uri, err := pr.ProcessURIWithInstructions(uri, i)

			if err != nil {
				msg := fmt.Sprintf("failed to process %s (%s) : %s", uri, label, err)
				err_ch <- errors.New(msg)
			}

			mu.Lock()
			results[label] = new_uri
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

	return results, nil
}

func (pr *ParallelProcessor) ProcessURIWithInstructions(uri string, i IIIFInstructions) (string, error) {

	return TransformURIWithInstructions(uri, i, pr.config, pr.source_cache, pr.dest_cache)
}
