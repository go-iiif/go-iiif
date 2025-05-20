package encode

import (
	"bufio"
	"bytes"
	"context"
	"image"
	"image/png"
	"io"
	_ "log/slog"

	"github.com/dsoprea/go-exif/v3"
	"github.com/dsoprea/go-png-image-structure/v2"
)

func EncodePNG(ctx context.Context, wr io.Writer, im image.Image, ib *exif.IfdBuilder) error {

	if ib == nil {
		return png.Encode(wr, im)
	}

	var im_buf bytes.Buffer
	im_wr := bufio.NewWriter(&im_buf)

	err := png.Encode(im_wr, im)

	if err != nil {
		return err
	}

	im_wr.Flush()

	// Write EXIF back to PNG

	png_parser := pngstructure.NewPngMediaParser()

	mp, err := png_parser.ParseBytes(im_buf.Bytes())

	if err != nil {
		return err
	}

	cs := mp.(*pngstructure.ChunkSlice)

	err = cs.SetExif(ib)

	if err != nil {
		return err
	}

	err = cs.WriteTo(wr)

	if err != nil {
		return err
	}

	return nil
}
