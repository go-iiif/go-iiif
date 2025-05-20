//go:build !libheif

package decode

import (
	"fmt"
	"image"
)

func ImageFromHEIC(body []byte) (image.Image, error) {

	return nil, fmt.Errorf("Not implemented")
}
