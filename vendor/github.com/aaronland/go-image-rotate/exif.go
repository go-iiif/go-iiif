package rotate

import (
	"context"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"io"
)

func init() {
	exif.RegisterParsers(mknote.All...)
}

func GetImageOrientation(ctx context.Context, r io.Reader) (string, error) {

	x, err := exif.Decode(r)

	if err != nil {
		return "", err
	}

	o, err := x.Get(exif.Orientation)

	if err != nil {
		return "", err
	}

	orientation := o.String()
	return orientation, nil
}
