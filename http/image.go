package http

import (
	gohttp "net/http"
	"sync/atomic"
	"time"

	iiifcache "github.com/go-iiif/go-iiif/v7/cache"
	iiifconfig "github.com/go-iiif/go-iiif/v7/config"
	iiifdriver "github.com/go-iiif/go-iiif/v7/driver"
	iiifimage "github.com/go-iiif/go-iiif/v7/image"
	iiiflevel "github.com/go-iiif/go-iiif/v7/level"
	iiifsource "github.com/go-iiif/go-iiif/v7/source"
)

func ImageHandler(config *iiifconfig.Config, driver iiifdriver.Driver) (gohttp.HandlerFunc, error) {

	images_cache, err := iiifcache.NewImagesCacheFromConfig(config)

	if err != nil {
		return nil, err
	}

	derivatives_cache, err := iiifcache.NewDerivativesCacheFromConfig(config)

	if err != nil {
		return nil, err
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		logger := LoggerForRequest(req)

		params, err := IIIFParamtersFromRequest(req)

		if err != nil {
			logger.Error("Failed to derive IIIF parameters for request", "error", err)
			gohttp.Error(rsp, "Bad request", gohttp.StatusBadRequest)
			return
		}

		endpoint := EndpointFromRequest(req)

		level, err := iiiflevel.NewLevelFromConfig(config, endpoint)

		if err != nil {
			logger.Error("Failed to derive level from config", "error", err)
			gohttp.Error(rsp, "Internal server error", gohttp.StatusInternalServerError)
			return
		}

		compliance := level.Compliance()

		transformation, err := iiifimage.NewTransformation(compliance, params.Region, params.Size, params.Rotation, params.Quality, params.Format)

		if err != nil {
			logger.Error("Failed to derive transformation from params", "error", err)
			gohttp.Error(rsp, "Bad request", gohttp.StatusBadRequest)
			return
		}

		uri, err := transformation.ToURI(params.Identifier)

		if err != nil {
			logger.Error("Failed to derive URI from transmation", "id", params.Identifier, "error", err)
			gohttp.Error(rsp, "Bad request", gohttp.StatusBadRequest)
			return
		}

		body, err := derivatives_cache.Get(uri)

		if err == nil {

			cacheHit.Add(1)

			source, err := iiifsource.NewMemorySource(body)

			if err != nil {
				logger.Warn("Failed to derive image from cache body", "error", err)
			} else {

				image, err := driver.NewImageFromConfigWithSource(config, source, "cache")

				if err != nil {
					logger.Warn("Failed to derive new image from source (cache body)", "error", err)
				} else {
					rsp.Header().Set("Content-Type", image.ContentType())
					rsp.Write(image.Body())
					return
				}
			}
		}

		image, err := driver.NewImageFromConfigWithCache(config, images_cache, params.Identifier)

		if err != nil {
			logger.Error("Failed to return image from cache", "cache", images_cache, "id", params.Identifier, "error", err)
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
