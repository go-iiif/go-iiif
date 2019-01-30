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

func ParallelProcessURIWithInstructionSet(cfg *iiifconfig.Config, pr Processor, instruction_set IIIFInstructionSet, u URI) (map[string]interface{}, error) {

	done_ch := make(chan bool)
	err_ch := make(chan error)

	remaining := len(instruction_set) + 1 // + 1 because we fetch the palette for the source image

	results := make(map[string]interface{})

	uris := make(map[Label]string)
	dimensions := make(map[Label][]int)

	mu := new(sync.RWMutex)

	go func() {

		defer func() {
			done_ch <- true
		}()

		im, err := iiifimage.NewImageFromConfig(cfg, u.URL())

		if err != nil {
			msg := fmt.Sprintf("failed to derive palette for %s : %s", u, err)
			err_ch <- errors.New(msg)
			return
		}

		for _, service_name := range cfg.Profile.Services.Enable {

			if service_name == "palette" {

				s, err := iiifservice.NewPaletteService(cfg.Palette, im)

				if err != nil {
					msg := fmt.Sprintf("failed to derive palette for %s : %s", u, err)
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

		go func(u URI, label Label, i IIIFInstructions) {

			defer func() {
				done_ch <- true
			}()

			new_uri, im, err := pr.ProcessURIWithInstructions(u, label, i)

			if err != nil {
				msg := fmt.Sprintf("failed to process %s (%s) : %s", u, label, err)
				err_ch <- errors.New(msg)
				return
			}

			dims, err := im.Dimensions()

			if err != nil {
				msg := fmt.Sprintf("failed to process %s (%s) : %s", u, label, err)
				err_ch <- errors.New(msg)
				return
			}

			mu.Lock()

			uris[label] = new_uri.String()

			dimensions[label] = []int{
				dims.Width(),
				dims.Height(),
			}

			mu.Unlock()

		}(u, label, i)
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
