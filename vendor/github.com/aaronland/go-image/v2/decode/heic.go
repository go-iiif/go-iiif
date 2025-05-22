//go:build libheif

package decode

import (
	"fmt"
	"image"

	"github.com/strukturag/libheif-go"
)

func ImageFromHEIC(body []byte) (image.Image, error) {

	// https://github.com/spacestation93/heif_howto

	// First decode the HEIC image

	im_ctx, err := libheif.NewContext()

	if err != nil {
		return nil, fmt.Errorf("Failed to create new libheif context, %w", err)
	}

	err = im_ctx.ReadFromMemory(body)

	if err != nil {
		return nil, fmt.Errorf("Failed to read input data, %w", err)
	}

	im_handle, err := im_ctx.GetPrimaryImageHandle()

	if err != nil {
		return nil, fmt.Errorf("Failed to derive primary image handler, %w", err)
	}

	h_im, err := im_handle.DecodeImage(libheif.ColorspaceUndefined, libheif.ChromaUndefined, nil)

	if err != nil {
		return nil, fmt.Errorf("Failed to decode image, %w", err)
	}

	im, err := h_im.GetImage()

	if err != nil {
		return nil, fmt.Errorf("Failed to create image.Image, %w", err)
	}

	return im, nil
}
