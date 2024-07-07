package rotate

import (
	"context"
	"fmt"
	"io"
	_ "log"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)

func init() {
	exif.RegisterParsers(mknote.All...)
}

// GetImageOrientation returns the string representation of the image orientation (0-8) in 'r'.
func GetImageOrientation(ctx context.Context, r io.Reader) (string, error) {

	x, err := exif.Decode(r)

	// I'm not sure this is the best error handling but I am not sure
	// how else to trap "there is no EXIF" style errors...
	// (20200428/straup)

	if err == io.EOF {
		return "0", nil
	}

	if err != nil {
		return "", fmt.Errorf("Failed to decode EXIF data, %w", err)
	}

	o, err := x.Get(exif.Orientation)

	if exif.IsTagNotPresentError(err) == true {
		return "0", nil
	}

	if err != nil {
		return "", fmt.Errorf("Failed to derive Orientation EXIF flag, %w", err)
	}

	orientation := o.String()
	return orientation, nil
}
