package tile

import (
	"errors"
	"fmt"
	iiifimage "github.com/thisisaaronland/go-iiif/image"
	iiiflevel "github.com/thisisaaronland/go-iiif/level"
	_ "log"
	"math"
)

type TileSeed struct {
	level  iiiflevel.Level
	height int
	width  int
}

func NewTileSeed(level iiiflevel.Level, h int, w int) (*TileSeed, error) {

	ts := TileSeed{
		level:  level,
		height: h,
		width:  w,
	}

	return &ts, nil
}

func (ts *TileSeed) TileSizes(im iiifimage.Image, sf int) ([]*iiifimage.Transformation, error) {

	dims, err := im.Dimensions()

	if err != nil {
		return nil, err
	}

	w := dims.Width()
	h := dims.Height()

	if sf*ts.width >= w && sf*ts.height >= h {
		return nil, errors.New("E_EXCESSIVE_SCALEFACTOR")
	}

	crops := make([]*iiifimage.Transformation, 0)

	// what follows was copied from
	// https://github.com/cmoa/iiif_s3/blob/master/lib/iiif_s3/builder.rb#L165-L199

	ty := int(math.Floor(float64(h) / float64(ts.height*sf)))
	tx := int(math.Floor(float64(w) / float64(ts.width*sf)))

	for xpos := 0; xpos < ty; xpos++ {

		for ypos := 0; ypos < tx; ypos++ {

			/*
				this is the data structure used by iiif_s3 and it's not
				clear how much of it is actually necessary here
				(20160911/thisisaaronland)
			*/

			foo := make(map[string]int)

			foo["scale_factor"] = sf
			foo["xpos"] = xpos
			foo["ypos"] = ypos
			foo["x"] = xpos * ts.width * sf
			foo["y"] = ypos * ts.height * sf
			foo["width"] = ts.width * sf
			foo["height"] = ts.height * sf
			foo["xsize"] = ts.width
			foo["ysize"] = ts.height

			if (foo["x"] + ts.width) > w {
				foo["width"] = w - foo["x"]
				foo["xsize"] = int(math.Ceil(float64(foo["width"]) / float64(sf)))
			}

			if (foo["y"] + ts.height) > h {

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

			_s := ts.width

			if _x+_w > w {
				_s = w - _x
			}

			if _y+_h > h {
				_s = h - _y
			}

			// fmt.Printf("%d,%d,%d,%d\t%d\n", _x, _y, _w, _h, _s)

			region := fmt.Sprintf("%d,%d,%d,%d", _x, _y, _w, _h)
			size := fmt.Sprintf("%d,", _s) // but maybe some client will send 'full' or what...?
			rotation := "0"
			quality := "color" // but maybe some client will send 'default'?
			format := "jpg"

			transformation, err := iiifimage.NewTransformation(ts.level, region, size, rotation, quality, format)

			if err != nil {
				return nil, err
			}

			crops = append(crops, transformation)
		}

	}

	return crops, nil
}
