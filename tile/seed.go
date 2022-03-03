package tile

import (
	"encoding/json"
	"fmt"
	iiifcache "github.com/go-iiif/go-iiif/v4/cache"
	iiifconfig "github.com/go-iiif/go-iiif/v4/config"
	iiifdriver "github.com/go-iiif/go-iiif/v4/driver"
	iiifimage "github.com/go-iiif/go-iiif/v4/image"
	iiiflevel "github.com/go-iiif/go-iiif/v4/level"
	iiifprofile "github.com/go-iiif/go-iiif/v4/profile"
	iiifsource "github.com/go-iiif/go-iiif/v4/source"
	"github.com/tidwall/pretty"
	"log"
	"math"
	"runtime"
	"sync"
)

type TileSeed struct {
	config            *iiifconfig.Config
	driver            iiifdriver.Driver
	level             iiiflevel.Level
	images_cache      iiifcache.Cache
	derivatives_cache iiifcache.Cache
	Endpoint          string
	Height            int
	Width             int
	Quality           string
	Format            string
	procs             int
}

func NewTileSeed(config *iiifconfig.Config, h int, w int, endpoint string, quality string, format string) (*TileSeed, error) {

	driver, err := iiifdriver.NewDriverFromConfig(config)

	if err != nil {
		return nil, fmt.Errorf("Failed to create tileseed driver from config, %w", err)
	}

	// level, err := iiiflevel.NewLevelFromConfig(config, endpoint)

	level, err := iiiflevel.NewLevel0(config, endpoint)

	if err != nil {
		return nil, fmt.Errorf("Failed to create level0 instance, %w", err)
	}

	images_cache, err := iiifcache.NewImagesCacheFromConfig(config)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive image cache from config, %w", err)
	}

	derivatives_cache, err := iiifcache.NewDerivativesCacheFromConfig(config)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive derivatives cache from config, %w", err)
	}

	compliance := level.Compliance()
	_, err = compliance.DefaultQuality()

	if err != nil {
		return nil, fmt.Errorf("Failed to ensure default quality for compliance, %w", err)
	}

	procs := runtime.NumCPU()

	ts := TileSeed{
		config:            config,
		driver:            driver,
		level:             level,
		images_cache:      images_cache,
		derivatives_cache: derivatives_cache,
		Endpoint:          endpoint,
		Height:            h,
		Width:             w,
		Quality:           quality,
		Format:            format,
		procs:             procs,
	}

	return &ts, nil
}

func (ts *TileSeed) SeedTiles(src_id string, alt_id string, scales []int, refresh bool) (int, error) {

	count := 0

	image, err := ts.driver.NewImageFromConfigWithCache(ts.config, ts.images_cache, src_id)

	if err != nil {
		return count, fmt.Errorf("Failed to create image for %s, %w", src_id, err)
	}

	// https://github.com/go-iiif/go-iiif/issues/25
	// https://github.com/go-iiif/go-iiif/issues/27

	// 191733_5755a1309e4d66a7_k.jpg,191733_5755a1309e4d66a7
	// means
	// store '191733_5755a1309e4d66a7_k.jpg' as 'CACHEROOT/191733_5755a1309e4d66a7'

	// 191733_5755a1309e4d66a7_k.jpg,191/733/191733_5755a1309e4d66a7_k.jpg
	// means
	// store '191733_5755a1309e4d66a7_k.jpg' as 'CACHEROOT/191/733/191733_5755a1309e4d66a7_k.jpg'

	// 191733_5755a1309e4d66a7_k.jpg,191/733/191733_5755a1309e4d66a7
	// means
	// store '191733_5755a1309e4d66a7_k.jpg' as 'CACHEROOT/191/733/191733_5755a1309e4d66a7'

	// the relevant part being that if basename(ALT_ID) != src_id then we need to signal
	// to iiifimage.Image that its Identifier() method needs to return basename(ALT_ID)
	// (20160925/thisisaaronland)

	if src_id != alt_id {

		err = image.Rename(alt_id)

		if err != nil {
			return count, fmt.Errorf("Failed to rename (%s as %s), %w", src_id, alt_id, err)
		}
	}

	source, err := iiifsource.NewMemorySource(image.Body())

	if err != nil {
		return count, fmt.Errorf("Failed to create image from memory, %w", err)
	}

	throttle := make(chan bool, ts.procs)

	for i := 0; i < ts.procs; i++ {
		throttle <- true
	}

	for _, scale := range scales {

		crops, err := ts.TileSizes(image, scale)

		if err != nil {
			log.Printf("Unable to determine tile sizes for %d (%s), %s\n", scale, image.Identifier(), err)
			continue
		}

		wg := new(sync.WaitGroup)

		for _, transformation := range crops {

			<-throttle

			wg.Add(1)

			go func(throttle chan bool, im iiifimage.Image, tr *iiifimage.Transformation, wg *sync.WaitGroup) {

				defer func() {
					wg.Done()
					throttle <- true
				}()

				// log.SetOutput(os.Stdout)

				uri, _ := tr.ToURI(alt_id)

				if !refresh {

					_, err := ts.derivatives_cache.Get(uri)

					if err == nil {
						return
					}
				}

				tmp, err := ts.driver.NewImageFromConfigWithSource(ts.config, source, im.Identifier())

				if err != nil {
					log.Println(err)
					return
				}

				err = tmp.Transform(tr)

				if err != nil {
					log.Printf("[%s] transform failed: %s\n", im.Identifier(), err)
					return
				}

				err = ts.derivatives_cache.Set(uri, tmp.Body())

				if err != nil {
					log.Printf("[%s] cache set (%s) failed: %s\n", im.Identifier(), uri, err)
					return
				}

				return

			}(throttle, image, transformation, wg)
		}

		wg.Wait()

		// something something something using the channel above to increment count...

		count += len(crops)
	}

	level, err := iiiflevel.NewLevel0(ts.config, ts.Endpoint)

	if err != nil {
		return count, fmt.Errorf("Failed to create new level 0, %w", err)
	}

	profile, err := level.Profile(ts.Endpoint, image)

	if err != nil {
		return count, fmt.Errorf("Failed to create new profile for level, %w", err)
	}

	body, err := json.Marshal(profile)

	if err != nil {
		return count, fmt.Errorf("Failed to marshal profile, %w", err)
	}

	body = pretty.Pretty(body)

	uri := fmt.Sprintf("%s/info.json", alt_id)
	ts.derivatives_cache.Set(uri, body)

	return count, nil
}

func (ts *TileSeed) TileSizes(im iiifimage.Image, sf int) ([]*iiifimage.Transformation, error) {

	dims, err := im.Dimensions()

	if err != nil {
		return nil, fmt.Errorf("Failed to derive dimensions for image, %w", err)
	}

	w := dims.Width()
	h := dims.Height()

	if sf*ts.Width >= w && sf*ts.Height >= h {
		return nil, fmt.Errorf("E_EXCESSIVE_SCALEFACTOR %d (%d,%d) (%d,%d)", sf, w, h, sf*ts.Width, sf*ts.Height)
	}

	quality := ts.Quality

	if quality == "default" {
		compliance := ts.level.Compliance()
		quality, _ = compliance.DefaultQuality()
	}

	format := ts.Format

	crops := make([]*iiifimage.Transformation, 0)

	// what follows was copied from
	// https://github.com/cmoa/iiif_s3/blob/master/lib/iiif_s3/builder.rb#L165-L199

	ty := int(math.Ceil(float64(h) / float64(ts.Height*sf)))
	tx := int(math.Ceil(float64(w) / float64(ts.Width*sf)))

	for ypos := 0; ypos < ty; ypos++ {

		for xpos := 0; xpos < tx; xpos++ {

			crop := make(map[string]int)

			crop["x"] = xpos * ts.Width * sf
			crop["y"] = ypos * ts.Height * sf
			crop["width"] = ts.Width * sf
			crop["height"] = ts.Height * sf

			if (crop["x"] + ts.Width) > w {
				crop["width"] = w - crop["x"]
				crop["xsize"] = int(math.Ceil(float64(crop["width"]) / float64(sf)))
			}

			if (crop["y"] + ts.Height) > h {

				crop["height"] = h - crop["y"]
				crop["ysize"] = int(math.Ceil(float64(crop["height"]) / float64(sf)))
			}

			/*

				this logic is repeated again in image/transformation.go but unlike there
				we are simply doing it here in order to generate a cache URI that works
				with the leaflet plugin... I think? (20160911/thisisaaronland)

			*/

			_x := crop["x"]
			_y := crop["y"]
			_w := crop["width"]
			_h := crop["height"]

			_s := ts.Width

			if _x+_w > w {
				_w = w - _x
			}

			if _y+_h > h {
				_h = h - _y
			}

			// this bit is cribbed from leaflet-iiif.js

			base := float64(ts.Width * sf)

			minx := float64(xpos) * base
			maxx := math.Min(minx+base, float64(w))

			diff := maxx - minx
			sz := math.Ceil(diff / float64(sf))

			_s = int(sz)

			region := fmt.Sprintf("%d,%d,%d,%d", _x, _y, _w, _h)
			size := fmt.Sprintf("%d,", _s) // but maybe some client will send 'full' or what...?
			rotation := "0"
			quality := quality
			format := format

			transformation, err := iiifimage.NewTransformation(ts.level, region, size, rotation, quality, format)

			if err != nil {
				return nil, fmt.Errorf("Failed to create new transformation, %w", err)
			}

			crops = append(crops, transformation)
		}

	}

	return crops, nil
}
