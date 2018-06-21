package http

import (
	"encoding/json"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	iiifimage "github.com/thisisaaronland/go-iiif/image"
	iiiflevel "github.com/thisisaaronland/go-iiif/level"
	iiifprofile "github.com/thisisaaronland/go-iiif/profile"
	iiifservice "github.com/thisisaaronland/go-iiif/service"
	gohttp "net/http"
)

func InfoHandler(config *iiifconfig.Config) (gohttp.HandlerFunc, error) {

	fn := func(w gohttp.ResponseWriter, r *gohttp.Request) {

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

		image, err := iiifimage.NewImageFromConfig(config, id)

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

		profile, err := iiifprofile.NewProfile(endpoint, image, level)

		if err != nil {
			gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		for _, service_name := range config.Profile.Services.Enable {

			switch service_name {
			case "palette":

				service, err := iiifservice.NewPaletteService(config.Palette, image)

				if err != nil {
					gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
					return
				}

				profile.AddService(service)
			default:
				gohttp.Error(w, "Unsupported service", gohttp.StatusInternalServerError)
				return
			}
		}

		b, err := json.Marshal(profile)

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
