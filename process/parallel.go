package process

import (
	"errors"
	"fmt"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	iiifimage "github.com/thisisaaronland/go-iiif/image"
	iiifservice "github.com/thisisaaronland/go-iiif/service"
	"log"
	"sync"
)

func ParallelProcessURIWithInstructionSet(cfg *iiifconfig.Config, pr Processor, instruction_set IIIFInstructionSet, uri string) (map[string]interface{}, error) {

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

		im, err := iiifimage.NewImageFromConfig(cfg, uri)

		if err != nil {
			msg := fmt.Sprintf("failed to derive palette for %s : %s", uri, err)
			err_ch <- errors.New(msg)
			return
		}

		for _, service_name := range cfg.Profile.Services.Enable {

			if service_name == "palette" {

				s, err := iiifservice.NewPaletteService(cfg.Palette, im)

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

			new_uri, im, err := pr.ProcessURIWithInstructions(uri, label, i)

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
