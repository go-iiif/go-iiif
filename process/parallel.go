package process

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/url"
	"sync"

	iiifuri "github.com/go-iiif/go-iiif-uri"
	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
	iiifdriver "github.com/go-iiif/go-iiif/v6/driver"
	iiifservice "github.com/go-iiif/go-iiif/v6/service"
)

// ParallelProcessURIWithInstructionSet processes 'u' according to each instruction in 'instruction_set' in concurrent processes
func ParallelProcessURIWithInstructionSet(ctx context.Context, cfg *iiifconfig.Config, driver iiifdriver.Driver, pr Processor, instruction_set IIIFInstructionSet, u iiifuri.URI) (map[string]interface{}, error) {

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

		im, err := driver.NewImageFromConfig(ctx, cfg, origin)

		if err != nil {
			err_ch <- fmt.Errorf("Failed to load image (%s) for processing profile.services : %w", u, err)
			return
		}

		// PLEASE PUT THIS IN A FUNCTION SOMEWHERE

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
				err_ch <- fmt.Errorf("Failed to create service for %s : %w", service_name, err)
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

			switch u.Scheme() {
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
					err_ch <- fmt.Errorf("Failed to derive target %s (%s) : %w", u, label, err)
					return
				}

				origin := u.Origin()

				rw_str := fmt.Sprintf("%s:///%s?target=%s", iiifuri.REWRITE_SCHEME, origin, target_str)

				rw_uri, err := iiifuri.NewURI(ctx, rw_str)

				if err != nil {
					err_ch <- fmt.Errorf("Failed to generate rewrite URL %s (%s) : %w", u, label, err)
					return
				}

				process_uri = rw_uri

			default:
				process_uri = u
			}

			new_uri, im, err := pr.ProcessURIWithInstructions(ctx, process_uri, label, i)

			if err != nil {
				err_ch <- fmt.Errorf("Failed to process %s (%s) : %w", u.String(), label, err)
				return
			}

			dims, err := im.Dimensions()

			if err != nil {
				err_ch <- fmt.Errorf("Failed to process %s (%s) : %w", u, label, err)
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
		case err := <-err_ch:
			slog.Error(err.Error())
		default:
			//
		}
	}

	results["uris"] = uris
	results["dimensions"] = dimensions

	return results, nil
}
