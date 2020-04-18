package process

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	iiifuri "github.com/go-iiif/go-iiif-uri"
	iiifconfig "github.com/go-iiif/go-iiif/config"
	iiifdriver "github.com/go-iiif/go-iiif/driver"
	iiifservice "github.com/go-iiif/go-iiif/service"
	"log"
	"net/url"
	"sync"
)

func ParallelProcessURIWithInstructionSet(cfg *iiifconfig.Config, driver iiifdriver.Driver, pr Processor, instruction_set IIIFInstructionSet, u iiifuri.URI) (map[string]interface{}, error) {

	done_ch := make(chan bool)
	err_ch := make(chan error)

	remaining := len(instruction_set) + 1 // + 1 because we fetch the palette for the source image

	results := make(map[string]interface{})
	results["origin"] = u.Origin()
	results["origin_uri"] = u.String()

	uris := make(map[Label]string)
	dimensions := make(map[Label][]int)

	mu := new(sync.RWMutex)

	go func() {

		defer func() {
			done_ch <- true
		}()

		origin := u.Origin()

		im, err := driver.NewImageFromConfig(cfg, origin)

		if err != nil {
			msg := fmt.Sprintf("failed to load image (%s) for processing profile.services : %s", u, err)
			err_ch <- errors.New(msg)
			return
		}

		hash := sha1.Sum(im.Body())
		fingerprint := hex.EncodeToString(hash[:])

		mu.Lock()
		results["origin_fingerprint"] = fingerprint
		mu.Unlock()

		ctx := context.Background()

		for _, service_name := range cfg.Profile.Services.Enable {

			service_uri := fmt.Sprintf("%s://", service_name)
			service, err := iiifservice.NewService(ctx, service_uri, cfg, im)

			if err != nil {
				msg := fmt.Sprintf("failed to create service for %s : %v", service_name, err)
				err_ch <- errors.New(msg)
				return
			}

			mu.Lock()
			results[service_name] = service.Value()
			mu.Unlock()
		}

	}()

	for label, i := range instruction_set {

		i = EnsureInstructions(i)

		go func(u iiifuri.URI, label Label, i IIIFInstructions) {

			defer func() {
				done_ch <- true
			}()

			var process_uri iiifuri.URI

			switch u.Driver() {
			case "idsecret":

				str_label := fmt.Sprintf("%s", label)

				opts := &url.Values{}
				opts.Set("label", str_label)
				opts.Set("format", i.Format)

				if str_label == "o" {
					opts.Set("original", "1")
				}

				target_str, err := u.Target(opts)

				if err != nil {
					msg := fmt.Sprintf("failed to derive target %s (%s) : %s", u, label, err)
					err_ch <- errors.New(msg)
					return
				}

				origin := u.Origin()

				rw_str := fmt.Sprintf("%s?target=%s", origin, target_str)
				rw_str = iiifuri.NewRewriteURIString(rw_str)

				rw_uri, err := iiifuri.NewURI(rw_str)

				if err != nil {
					msg := fmt.Sprintf("failed to generate rewrite URL %s (%s) : %s", u, label, err)
					err_ch <- errors.New(msg)
					return
				}

				process_uri = rw_uri

			default:
				process_uri = u
			}

			new_uri, im, err := pr.ProcessURIWithInstructions(process_uri, label, i)

			if err != nil {
				msg := fmt.Sprintf("failed to process %s (%s) : %s", u.String(), label, err)
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
