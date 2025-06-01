package encode

import (
	"context"
	"image"
	"image/gif"
	"io"
	"log/slog"

	"github.com/dsoprea/go-exif/v3"
)

func EncodeGIF(ctx context.Context, wr io.Writer, im image.Image, ib *exif.IfdBuilder, gif_opts *gif.Options) error {

	slog.Debug("WriteGIF method does not support writing EXIF data.")
	return gif.Encode(wr, im, gif_opts)
}
