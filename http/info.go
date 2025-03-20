package http

import (
	"fmt"
	gohttp "net/http"

	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
	iiifdriver "github.com/go-iiif/go-iiif/v6/driver"
	iiifinfo "github.com/go-iiif/go-iiif/v6/info"
	iiiflevel "github.com/go-iiif/go-iiif/v6/level"
	iiifservice "github.com/go-iiif/go-iiif/v6/service"
)

func InfoHandler(config *iiifconfig.Config, driver iiifdriver.Driver) (gohttp.HandlerFunc, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		ctx := req.Context()
		logger := LoggerWithRequest(req, nil)

		id, err := GetIIIFParameter(req, "identifier")

		if err != nil {
			logger.Error("Failed to derive identifier", "error", err)
			gohttp.Error(rsp, "Bad request", gohttp.StatusBadRequest)
			return
		}

		logger = logger.With("identifier", id)

		image, err := driver.NewImageFromConfig(config, id)

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

			services := make([]iiifservice.Service, count_services)

			for idx, service_name := range config.Profile.Services.Enable {

				service_uri := fmt.Sprintf("%s://", service_name)
				service, err := iiifservice.NewService(ctx, service_uri, config, image)

				if err != nil {
					logger.Error("Failed to derive service", "service", service_uri, "error", err)
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
		rsp.Header().Set("Access-Control-Allow-Origin", "*")
		rsp.Write(b)

	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
