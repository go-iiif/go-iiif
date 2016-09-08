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

			region := fmt.Sprintf("%d,%d,%d,%d", foo["x"], foo["y"], foo["width"], foo["height"])
			size := "full"
			rotation := "0"
			quality := "default"
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
