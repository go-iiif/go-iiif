package encode

import (
	"context"
	"image"
	"io"
	"log/slog"

	"github.com/dsoprea/go-exif/v3"
	"golang.org/x/image/bmp"
)

func EncodeBMP(ctx context.Context, wr io.Writer, im image.Image, ib *exif.IfdBuilder) error {

	if ib == nil {
		return bmp.Encode(wr, im)
	}

	slog.Debug("WriteBMP method does not support writing EXIF data.")
	return bmp.Encode(wr, im)
}
