package tiles

import (
	"errors"
	"fmt"
	iiifimage "github.com/thisisaaronland/go-iiif/image"
	_ "math"
)

type TileSeed struct {
	height int
	width  int
}

func NewTileSeed(h int, w int) (*TileSeed, error) {

	ts := TileSeed{
		height: h,
		width:  w,
	}

	return &ts, nil
}

func (ts *TileSeed) TileSizes(im iiifimage.Image, sf int) ([]iiifimage.Transformation, error) {

	dims, err := im.Dimensions()

	if err != nil {
		return nil, err
	}

	w := dims.Width()
	h := dims.Height()

	if sf*ts.width >= w && sf*ts.height >= h {
		return nil, errors.New("E_EXCESSIVE_SCALEFACTOR")
	}

	crops := make([]iiifimage.Transformation, 0)

	// https://github.com/zimeon/iiif/blob/master/iiif/static.py#L21
	// https://github.com/cmoa/iiif_s3/blob/master/lib/iiif_s3/builder.rb#L165-L199

	// I AM PRETTY SURE THIS IS WRONG (20160907/thisisaaronland)

	ty := h / ts.height * sf
	tx := w / ts.width * sf

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
				foo["xsize"] = foo["width"] / sf // math.Ceil
			}

			if (foo["y"] + ts.height) > h {

				foo["height"] = h - foo["y"]
				foo["ysize"] = foo["height"] / sf // math.Ceil
			}

			fmt.Println(foo)
		}

	}

	/*

		      config.tile_scale_factors.each do |s|
		        (0..(height*1.0/(tile_width*s)).floor).each do |tileY|
		          (0..(width*1.0/(tile_width*s)).floor).each do |tileX|
		            tile = {
		              scale_factor: s,
		              xpos: tileX,
		              ypos: tileY,
		              x: tileX * tile_width * s,
		              y: tileY * tile_width * s,
		              width: tile_width * s,
		              height: tile_width * s,
		              xSize: tile_width,
		              ySize: tile_width
		            }
		            if (tile[:x] + tile[:width]  > width)
		              tile[:width]  = width  - tile[:x]
		              tile[:xSize]  = (tile[:width]/(s*1.0)).ceil
		            end
		            if (tile[:y] + tile[:height] > height)
		              tile[:height] = height - tile[:y]
		              tile[:ySize]  = (tile[:height]/(s*1.0)).ceil
		            end
		            tiles.push(tile)
		          end
		        end
		end
	*/

	return crops, nil
}
