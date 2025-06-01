//go:build libheif

package encode

import (
	"context"
	"fmt"
	"image"
	"io"
	"log/slog"

	"github.com/dsoprea/go-exif/v3"
	"github.com/strukturag/libheif-go"
)

func EncodeHEIC(ctx context.Context, wr io.Writer, im image.Image, ib *exif.IfdBuilder) error {

	slog.Debug("WriteHEIC method does not support writing EXIF data.")

	codec := libheif.CompressionHEVC

	heic_ctx, _, err := libheif.EncodeFromImage(im, codec,
		libheif.SetEncoderQuality(100),
		libheif.SetEncoderLossless(libheif.LosslessModeEnabled),
	)

	if err != nil {
		return fmt.Errorf("Failed to encode image, %w", err)
	}

	err = heic_ctx.Write(wr)

	if err != nil {
		return fmt.Errorf("Failed to write image, %w", err)
	}

	return nil
}
