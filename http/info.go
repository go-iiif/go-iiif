package http

import (
	"fmt"
	gohttp "net/http"

	iiifconfig "github.com/go-iiif/go-iiif/v7/config"
	iiifdriver "github.com/go-iiif/go-iiif/v7/driver"
	iiifinfo "github.com/go-iiif/go-iiif/v7/info"
	iiiflevel "github.com/go-iiif/go-iiif/v7/level"
	iiifservice "github.com/go-iiif/go-iiif/v7/service"
)

func InfoHandler(config *iiifconfig.Config, driver iiifdriver.Driver) (gohttp.HandlerFunc, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		logger := LoggerForRequest(req)

		ctx := req.Context()

		params, err := IIIFParamtersFromRequest(req)

		if err != nil {
			logger.Error("Failed to derive IIIF parameters for request", "error", err)
			gohttp.Error(rsp, "Bad request", gohttp.StatusBadRequest)
			return
		}

		id := params.Identifier
		logger = logger.With("id", id)

		image, err := driver.NewImageFromConfig(config, id)

		if err != nil {
			logger.Error("Failed to derive image from id", "error", err)
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

			services := make([]iiifservice.Service, count_services)

			for idx, service_name := range config.Profile.Services.Enable {

				service_uri := fmt.Sprintf("%s://", service_name)
				service, err := iiifservice.NewService(ctx, service_uri, config, image)

				if err != nil {
					logger.Error("Failed to instantiate server", "service", service_uri, "error", err)
					gohttp.Error(rsp, "Internal server error", gohttp.StatusInternalServerError)
					return
				}

				services[idx] = service
			}

			info.Services = services
		}

		b, err := iiifinfo.MarshalJSON(info)

		if err != nil {
			logger.Error("Failed to marshal info", "error", err)
			gohttp.Error(rsp, "Internal server error", gohttp.StatusInternalServerError)
			return
		}

		rsp.Header().Set("Content-Type", "application/json")
		rsp.Write(b)

	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
