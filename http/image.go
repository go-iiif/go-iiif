package http

import (
	_ "log"
	gohttp "net/http"
	"sync/atomic"
	"time"

	iiifcache "github.com/go-iiif/go-iiif/v6/cache"
	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
	iiifdriver "github.com/go-iiif/go-iiif/v6/driver"
	iiifimage "github.com/go-iiif/go-iiif/v6/image"
	iiiflevel "github.com/go-iiif/go-iiif/v6/level"
	iiifsource "github.com/go-iiif/go-iiif/v6/source"
)

func ImageHandler(config *iiifconfig.Config, driver iiifdriver.Driver, images_cache iiifcache.Cache, derivatives_cache iiifcache.Cache) (gohttp.HandlerFunc, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		logger := LoggerWithRequest(req, nil)

		params, err := GetIIIFParameters(req)

		if err != nil {
			logger.Error("Failed to derive IIIF parameters", "error", err)
			gohttp.Error(rsp, "Bad request", gohttp.StatusBadRequest)
			return
		}

		logger = logger.With("identifier", params.Identifier)

		endpoint := EndpointFromRequest(req)
		logger = logger.With("endpoint", endpoint)

		level, err := iiiflevel.NewLevelFromConfig(config, endpoint)

		if err != nil {
			logger.Error("Failed to derive level from config", "error", err)
			gohttp.Error(rsp, "Internal server error", gohttp.StatusInternalServerError)
			return
		}

		compliance := level.Compliance()

		transformation, err := iiifimage.NewTransformation(compliance, params.Region, params.Size, params.Rotation, params.Quality, params.Format)

		if err != nil {
			logger.Error("Failed to derive new transformation", "error", err)
			gohttp.Error(rsp, "Bad request", gohttp.StatusBadRequest)
			return
		}

		uri, err := transformation.ToURI(params.Identifier)

		if err != nil {
			logger.Error("Failed to derive transformation URI", "error", err)
			gohttp.Error(rsp, "Bad request", gohttp.StatusBadRequest)
			return
		}

		body, err := derivatives_cache.Get(uri)

		if err == nil {

			logger.Debug("Cache hit for image")
			cacheHit.Add(1)

			source, _ := iiifsource.NewMemorySource(body)
			image, _ := driver.NewImageFromConfigWithSource(config, source, "cache")

			rsp.Header().Set("Content-Type", image.ContentType())
			rsp.Write(image.Body())
			return
		}

		image, err := driver.NewImageFromConfigWithCache(config, images_cache, params.Identifier)

		if err != nil {
			logger.Error("Failed to retrieve image from cache", "error", err)
			gohttp.Error(rsp, "Internal server error", gohttp.StatusInternalServerError)
			return
		}

		/*
			something something something maybe sendfile something something
			(20160901/thisisaaronland)
		*/

		if transformation.HasTransformation() {

			cacheMiss.Add(1)

			t1 := time.Now()
			err = image.Transform(transformation)
			t2 := time.Since(t1)

			if err != nil {
				logger.Error("Failed to apply transformation", "error", err)
				gohttp.Error(rsp, "Internal server error", gohttp.StatusInternalServerError)
				return
			}

			go func(t time.Duration) {

				ns := t.Nanoseconds()
				ms := ns / (int64(time.Millisecond) / int64(time.Nanosecond))

				timers_mu.Lock()

				counter := atomic.AddInt64(&transforms_counter, 1)
				timer := atomic.AddInt64(&transforms_timer, ms)

				avg := float64(timer) / float64(counter)

				transformsCount.Add(1)
				transformsAvgTime.Set(avg)

				timers_mu.Unlock()
			}(t2)

			go func(k string, im iiifimage.Image) {
				logger.Debug("Set cache for image")
				derivatives_cache.Set(k, im.Body())
				cacheSet.Add(1)

			}(uri, image)
		}

		rsp.Header().Set("Content-Type", image.ContentType())
		rsp.Write(image.Body())
		return
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
