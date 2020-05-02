package rotate

import (
	"context"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"io"
	_ "log"
)

func init() {
	exif.RegisterParsers(mknote.All...)
}

func GetImageOrientation(ctx context.Context, r io.Reader) (string, error) {

	x, err := exif.Decode(r)

	// I'm not sure this is the best error handling but I am not sure
	// how else to trap "there is no EXIF" style errors...
	// (20200428/straup)

	if err == io.EOF {
		return "0", nil
	}

	if err != nil {
		return "", err
	}

	o, err := x.Get(exif.Orientation)

	if exif.IsTagNotPresentError(err) == true {
		return "0", nil
	}

	if err != nil {
		return "", err
	}

	orientation := o.String()
	return orientation, nil
}
