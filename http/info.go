package http

import (
	"encoding/json"
	iiifconfig "github.com/go-iiif/go-iiif/config"
	iiifdriver "github.com/go-iiif/go-iiif/driver"
	iiiflevel "github.com/go-iiif/go-iiif/level"
	iiifprofile "github.com/go-iiif/go-iiif/profile"
	iiifservice "github.com/go-iiif/go-iiif/service"
	gohttp "net/http"
)

func InfoHandler(config *iiifconfig.Config, driver iiifdriver.Driver) (gohttp.HandlerFunc, error) {

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

		profile, err := iiifprofile.NewProfile(endpoint, image, level)

		if err != nil {
			gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		for _, service_name := range config.Profile.Services.Enable {

			// this is dumb (20200224/thisisaaronland)
			// https://github.com/go-iiif/go-iiif/issues/71
			
			switch service_name {

			case "blurhash":

				service, err := iiifservice.NewBlurHashService(config.BlurHash, image)

				if err != nil {
					gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
					return
				}

				profile.AddService(service)

			case "imagehash":

				service, err := iiifservice.NewImageHashService(config.ImageHash, image)

				if err != nil {
					gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
					return
				}

				profile.AddService(service)

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
