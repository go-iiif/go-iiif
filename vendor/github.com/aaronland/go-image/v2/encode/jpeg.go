package encode

import (
	"bufio"
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"io"
	_ "log/slog"

	"github.com/dsoprea/go-exif/v3"
	"github.com/dsoprea/go-jpeg-image-structure/v2"
)

func EncodeJPEG(ctx context.Context, wr io.Writer, im image.Image, ib *exif.IfdBuilder, jpeg_opts *jpeg.Options) error {

	if jpeg_opts == nil {

		jpeg_opts = &jpeg.Options{
			Quality: 100,
		}
	}

	if ib == nil {
		return jpeg.Encode(wr, im, jpeg_opts)
	}

	// Do EXIF dance

	var im_buf bytes.Buffer
	im_wr := bufio.NewWriter(&im_buf)

	err := jpeg.Encode(im_wr, im, jpeg_opts)

	if err != nil {
		return err
	}

	im_wr.Flush()

	// Write EXIF back to JPEG

	jmp := jpegstructure.NewJpegMediaParser()

	mp, err := jmp.ParseBytes(im_buf.Bytes())

	if err != nil {
		return err
	}

	sl := mp.(*jpegstructure.SegmentList)

	err = sl.SetExif(ib)

	if err != nil {
		return err
	}

	err = sl.Write(wr)

	if err != nil {
		return err
	}

	return nil
}
