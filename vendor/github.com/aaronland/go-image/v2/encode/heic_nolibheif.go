//go:build !libheif

package encode

import (
	"context"
	"fmt"
	"image"
	"io"
	_ "log/slog"

	"github.com/dsoprea/go-exif/v3"
)

// https://github.com/strukturag/libheif-go/blob/master/encoder_test.go#L64

func EncodeHEIC(ctx context.Context, wr io.Writer, im image.Image, ib *exif.IfdBuilder) error {
	return fmt.Errorf("Not implemented.")
}
