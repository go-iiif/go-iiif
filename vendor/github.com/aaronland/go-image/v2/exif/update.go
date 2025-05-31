// Package exif provides methods for manipulating EXIF data in images.
package exif

import (
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"

	"github.com/sfomuseum/go-exif-update"
)

// UpdateExif will encode 'im' as a JPEG image and write it to a temporary file and then
// copy that (temporary) file to 'wr' with EXIF properties defined by 'exif_props'.
func UpdateExif(im image.Image, wr io.Writer, exif_props map[string]interface{}) error {

	temp_wr, err := os.CreateTemp("", "exif.*.jpg")

	if err != nil {
		return fmt.Errorf("Failed to create temp file, %w", err)
	}

	defer os.Remove(temp_wr.Name())

	jpeg_opts := &jpeg.Options{
		Quality: 100,
	}

	err = jpeg.Encode(temp_wr, im, jpeg_opts)

	if err != nil {
		return fmt.Errorf("Failed to write JPEG, %w", err)
	}

	err = temp_wr.Close()

	if err != nil {
		return fmt.Errorf("Failed to close, %w", err)
	}

	jpeg_r, err := os.Open(temp_wr.Name())

	if err != nil {
		return fmt.Errorf("Failed to open %s, %v", temp_wr.Name(), err)
	}

	defer jpeg_r.Close()

	err = update.UpdateExif(jpeg_r, wr, exif_props)

	if err != nil {
		return fmt.Errorf("Failed to update EXIF data, %w", err)
	}

	return nil
}
