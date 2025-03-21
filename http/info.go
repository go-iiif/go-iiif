package http

import (
	"fmt"
	gohttp "net/http"
	"sync"
	"time"

	iiifcache "github.com/go-iiif/go-iiif/v6/cache"
	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
	iiifdriver "github.com/go-iiif/go-iiif/v6/driver"
	iiifinfo "github.com/go-iiif/go-iiif/v6/info"
	iiiflevel "github.com/go-iiif/go-iiif/v6/level"
	iiifservice "github.com/go-iiif/go-iiif/v6/service"
)

func InfoHandler(config *iiifconfig.Config, driver iiifdriver.Driver, images_cache iiifcache.Cache) (gohttp.HandlerFunc, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		ctx := req.Context()
		logger := LoggerWithRequest(req, nil)

		t1 := time.Now()

		defer func() {
			logger.Debug("Time to process request", "time", time.Since(t1))
		}()

		id, err := GetIIIFParameter(req, "identifier")

		if err != nil {
			logger.Error("Failed to derive identifier", "error", err)
			gohttp.Error(rsp, "Bad request", gohttp.StatusBadRequest)
			return
		}

		logger = logger.With("identifier", id)

		image, err := driver.NewImageFromConfigWithCache(config, images_cache, id)

		logger.Debug("WUT 0", "time", time.Since(t1))

		if err != nil {
			logger.Error("Failed to derive image from config", "error", err)
			gohttp.Error(rsp, "Internal server error", gohttp.StatusInternalServerError)
			return
		}

		endpoint := EndpointFromRequest(req)

		level, err := iiiflevel.NewLevelFromConfig(config, endpoint)

		if err != nil {
			logger.Error("Failed to derive level from config", "error", err)
			gohttp.Error(rsp, "Internal server error", gohttp.StatusInternalServerError)
			return
		}

		info, err := iiifinfo.New(iiifinfo.IMAGE_V2_CONTEXT, level, image)

		if err != nil {
			logger.Error("Failed to derive info", "error", err)
			gohttp.Error(rsp, "Internal server error", gohttp.StatusInternalServerError)
			return
		}

		count_services := len(config.Profile.Services.Enable)

		if count_services > 0 {

			ts1 := time.Now()

			services := make([]iiifservice.Service, count_services)

			mu := new(sync.RWMutex)
			done_ch := make(chan bool)
			err_ch := make(chan error)

			for idx, service_name := range config.Profile.Services.Enable {

				go func(idx int, service_name string) {

					ts2 := time.Now()

					defer func() {
						logger.Debug("Time to load service", "service", service_name, "time", time.Since(ts2))
						done_ch <- true
					}()

					service_uri := fmt.Sprintf("%s://", service_name)
					service, err := iiifservice.NewService(ctx, service_uri, config, image)

					if err != nil {
						err_ch <- fmt.Errorf("Failed to derive service for '%s', %w", service_uri, err)
						return
					}

					mu.Lock()
					services[idx] = service
					mu.Unlock()

				}(idx, service_name)
			}

			remaining := count_services

			for remaining > 0 {
				select {
				case <-done_ch:
					remaining -= 1
				case err := <-err_ch:
					logger.Error("Failed to derive service", "error", err)
					gohttp.Error(rsp, "Internal server error", gohttp.StatusInternalServerError)
					return
				}
			}

			info.Services = services

			logger.Debug("Time to load services", "count", count_services, "time", time.Since(ts1))
		}

		rsp.Header().Set("Content-Type", "application/json")
		rsp.Header().Set("Access-Control-Allow-Origin", "*")

		err = info.MarshalJSON(rsp)

		if err != nil {
			logger.Error("Failed to marshal info", "error", err)
			gohttp.Error(rsp, "Internal server error", gohttp.StatusInternalServerError)
			return
		}
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
