package tile

import (
	"encoding/json"
	"errors"
	"fmt"
	iiifcache "github.com/thisisaaronland/go-iiif/cache"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	iiifimage "github.com/thisisaaronland/go-iiif/image"
	iiiflevel "github.com/thisisaaronland/go-iiif/level"
	iiifprofile "github.com/thisisaaronland/go-iiif/profile"
	iiifsource "github.com/thisisaaronland/go-iiif/source"
	"log"
	"math"
	_ "path/filepath"
	"runtime"
	_ "strings"
	"sync"
)

type TileSeed struct {
	config            *iiifconfig.Config
	level             iiiflevel.Level
	images_cache      iiifcache.Cache
	derivatives_cache iiifcache.Cache
	Endpoint          string
	Height            int
	Width             int
	Quality           string
	Format            string
}

func NewTileSeed(config *iiifconfig.Config, h int, w int, endpoint string, quality string, format string) (*TileSeed, error) {

	level, err := iiiflevel.NewLevelFromConfig(config, endpoint)

	if err != nil {
		log.Fatal(err)
	}

	images_cache, err := iiifcache.NewImagesCacheFromConfig(config)

	if err != nil {
		log.Fatal(err)
	}

	derivatives_cache, err := iiifcache.NewDerivativesCacheFromConfig(config)

	if err != nil {
		log.Fatal(err)
	}

	compliance := level.Compliance()
	_, err = compliance.DefaultQuality()

	if err != nil {
		return nil, err
	}

	ts := TileSeed{
		config:            config,
		level:             level,
		images_cache:      images_cache,
		derivatives_cache: derivatives_cache,
		Endpoint:          endpoint,
		Height:            h,
		Width:             w,
		Quality:           quality,
		Format:            format,
	}

	return &ts, nil
}

func (ts *TileSeed) SeedTiles(src_id string, alt_id string, scales []int, refresh bool) (int, error) {

	count := 0

	image, err := iiifimage.NewImageFromConfigWithCache(ts.config, ts.images_cache, src_id)

	if err != nil {
		return count, err
	}

	// https://github.com/thisisaaronland/go-iiif/issues/25
	// https://github.com/thisisaaronland/go-iiif/issues/27

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
			return count, err
		}
	}

	source, err := iiifsource.NewMemorySource(image.Body())

	if err != nil {
		return count, err
	}

	for _, scale := range scales {

		crops, err := ts.TileSizes(image, scale)

		if err != nil {
			log.Println(err)
			continue
		}

		procs := runtime.NumCPU() * 2 // move me in to the constructor...

		ch := make(chan bool, procs)

		for i := 0; i < procs; i++ {
			ch <- true
		}

		wg := new(sync.WaitGroup)

		for _, transformation := range crops {

			wg.Add(1)

			go func(im iiifimage.Image, tr *iiifimage.Transformation, wg *sync.WaitGroup) {

				<-ch

				defer func() {
					wg.Done()
					ch <- true
				}()

				uri, _ := tr.ToURI(alt_id)

				if !refresh {

					_, err := ts.derivatives_cache.Get(uri)

					if err == nil {
						return
					}
				}

				tmp, _ := iiifimage.NewImageFromConfigWithSource(ts.config, source, im.Identifier())

				err = tmp.Transform(tr)

				if err == nil {
					ts.derivatives_cache.Set(uri, tmp.Body())
				}

			}(image, transformation, wg)
		}

		wg.Wait()

		// something something something using the channel above to increment count...

		count += len(crops)
	}

	level, err := iiiflevel.NewLevelFromConfig(ts.config, ts.Endpoint)

	if err != nil {
		return count, err
	}

	profile, err := iiifprofile.NewProfile(ts.Endpoint, image, level)

	if err != nil {
		return count, err
	}

	body, err := json.Marshal(profile)

	if err != nil {
		return count, err
	}

	uri := fmt.Sprintf("%s/info.json", alt_id)
	ts.derivatives_cache.Set(uri, body)

	return count, nil
}

func (ts *TileSeed) TileSizes(im iiifimage.Image, sf int) ([]*iiifimage.Transformation, error) {

	dims, err := im.Dimensions()

	if err != nil {
		return nil, err
	}

	w := dims.Width()
	h := dims.Height()

	if sf*ts.Width >= w && sf*ts.Height >= h {
		msg := fmt.Sprintf("E_EXCESSIVE_SCALEFACTOR %d (%d,%d) (%d,%d)", sf, w, h, sf*ts.Width, sf*ts.Height)
		return nil, errors.New(msg)
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

	// fmt.Printf("%d / %d * %d\n", w, ts.Width, sf)
	// fmt.Printf("tx %d ty %d\n", tx, ty)

	for ypos := 0; ypos < ty; ypos++ {

		for xpos := 0; xpos < tx; xpos++ {

			/*
				this is the data structure used by iiif_s3 and it's not
				clear how much of it is actually necessary here
				(20160911/thisisaaronland)
			*/

			foo := make(map[string]int)

			foo["scale_factor"] = sf
			foo["xpos"] = xpos
			foo["ypos"] = ypos
			foo["x"] = xpos * ts.Width * sf
			foo["y"] = ypos * ts.Height * sf
			foo["width"] = ts.Width * sf
			foo["height"] = ts.Height * sf
			foo["xsize"] = ts.Width
			foo["ysize"] = ts.Height

			if (foo["x"] + ts.Width) > w {
				foo["width"] = w - foo["x"]
				foo["xsize"] = int(math.Ceil(float64(foo["width"]) / float64(sf)))
			}

			if (foo["y"] + ts.Height) > h {

				foo["height"] = h - foo["y"]
				foo["ysize"] = int(math.Ceil(float64(foo["height"]) / float64(sf)))
			}

			/*

				this logic is repeated again in image/transformation.go but unlike there
				we are simply doing it here in order to generate a cache URI that works
				with the leaflet plugin... I think? (20160911/thisisaaronland)

			*/

			_x := foo["x"]
			_y := foo["y"]
			_w := foo["width"]
			_h := foo["height"]

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

			// fmt.Printf("GO %d,%d,%d,%d\tsize:%d @ %d\n", _x, _y, _w, _h, _s, sf)

			region := fmt.Sprintf("%d,%d,%d,%d", _x, _y, _w, _h)
			size := fmt.Sprintf("%d,", _s) // but maybe some client will send 'full' or what...?
			rotation := "0"
			quality := quality
			format := format

			transformation, err := iiifimage.NewTransformation(ts.level, region, size, rotation, quality, format)

			if err != nil {
				return nil, err
			}

			crops = append(crops, transformation)
		}

	}

	return crops, nil
}
