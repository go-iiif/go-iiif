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

	fn := func(w gohttp.ResponseWriter, r *gohttp.Request) {

		ctx := r.Context()

		parser, err := NewIIIFQueryParser(r)

		if err != nil {
			gohttp.Error(w, err.Error(), gohttp.StatusBadRequest)
			return
		}

		id, err := parser.GetIIIFParameter("identifier")

		if err != nil {
			gohttp.Error(w, err.Error(), gohttp.StatusBadRequest)
			return
		}

		image, err := driver.NewImageFromConfig(config, id)

		if err != nil {
			gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		endpoint := EndpointFromRequest(r)

		level, err := iiiflevel.NewLevelFromConfig(config, endpoint)

		if err != nil {
			gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		info, err := iiifinfo.New(iiifinfo.IMAGE_V2_CONTEXT, level, image)

		if err != nil {
			gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		count_services := len(config.Profile.Services.Enable)

		if count_services > 0 {

			services := make([]iiifservice.Service, count_services)

			for idx, service_name := range config.Profile.Services.Enable {

				service_uri := fmt.Sprintf("%s://", service_name)
				service, err := iiifservice.NewService(ctx, service_uri, config, image)

				if err != nil {
					gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
					return
				}

				services[idx] = service
			}

			info.Services = services
		}

		b, err := iiifinfo.MarshalJSON(info)

		if err != nil {
			gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(b)

	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
