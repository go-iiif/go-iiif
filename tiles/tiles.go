package tiles

import (
	"errors"
	iiifimage "github.com/thisisaaronland/go-iiif/image"
)

type ImageTile struct {
	height int
	width  int
}

func NewImageTile(h int, w int) (*ImageTile, error) {

	i := ImageTile{
		height: h,
		width:  w,
	}

	return &i, nil
}

func (t *ImageTile) TileSizes(im *iiifimage.Image, sf int) ([]iiifimage.Transformation, error) {

	dims, err := im.Dimensions()

	if err != nil {
		return err
	}

	w := dims.Width()
	h := dims.Height()

	if sf*t.width >= w && sf*t.height >= h {
		return nil, errors.New("E_EXCESSIVE_SCALEFACTOR")
	}

	crops := make([]iiifimage.Transformation, 0)

	x := 0
	y := 0

	// https://github.com/zimeon/iiif/blob/master/iiif/static.py#L21

	return crops, nil
}
