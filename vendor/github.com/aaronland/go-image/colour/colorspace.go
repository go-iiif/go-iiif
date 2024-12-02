package colour

import (
	"fmt"
	"io"

	"github.com/aaronland/go-image/exif"
)

// COLORSPACE_UNKNOWN is the EXIF value for unknown or undefined colorspaces.
const COLORSPACE_UNKNOWN uint16 = 65535

// COLORSPACE_SRGB is the EXIF value for the sRGB colorspace.
const COLORSPACE_SRGB uint16 = 1

// COLORSPACE_ARGB is the EXIF value for the Adobe RGB colorspace.
const COLORSPACE_ARGB uint16 = 2

// ColorSpace attempts to derive the EXIF ColorSpace value from the body of 'r'.
func ColorSpace(r io.Reader) (uint16, error) {

	tag, err := exif.TagValue(r, "ColorSpace", "IFD/Exif")

	if err != nil {
		return COLORSPACE_UNKNOWN, fmt.Errorf("Failed to determine tag value, %w", err)
	}

	v, err := tag.Value()

	if err != nil {
		return COLORSPACE_UNKNOWN, fmt.Errorf("Failed to derive tag value, %w", err)
	}

	colorspace := v.([]uint16)

	if len(colorspace) != 1 {
		return COLORSPACE_UNKNOWN, fmt.Errorf("Multiple values for colorspace")
	}

	return colorspace[0], nil
}
