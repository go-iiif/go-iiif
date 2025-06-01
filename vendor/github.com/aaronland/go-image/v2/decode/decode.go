package decode

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"io"
	"log/slog"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/aaronland/go-image/v2/rotate"
	"github.com/dsoprea/go-exif/v3"
	"github.com/dsoprea/go-heic-exif-extractor/v2"
	"github.com/dsoprea/go-jpeg-image-structure/v2"
	"github.com/dsoprea/go-png-image-structure/v2"
	"github.com/dsoprea/go-tiff-image-structure/v2"
	"github.com/gabriel-vasile/mimetype"
)

type DecodeImageOptions struct {
	Rotate bool
}

func DecodeImage(ctx context.Context, im_r io.ReadSeeker) (image.Image, string, *exif.Ifd, error) {

	opts := &DecodeImageOptions{
		Rotate: true,
	}

	return DecodeImageWithOptions(ctx, im_r, opts)
}

func DecodeImageWithOptions(ctx context.Context, im_r io.ReadSeeker, opts *DecodeImageOptions) (image.Image, string, *exif.Ifd, error) {

	var ifd *exif.Ifd
	var im image.Image

	im_body, err := io.ReadAll(im_r)

	if err != nil {
		return nil, "", nil, err
	}

	br := bytes.NewReader(im_body)

	im, im_fmt, err := image.Decode(br)

	if err != nil {
		// Check error here...
		slog.Debug("Failed to decode image natively", "error", err)
	}

	mtype := mimetype.Detect(im_body)

	switch im_fmt {
	case "gif", "webp", "bmp":
		// pass
	case "jpeg":

		jmp := jpegstructure.NewJpegMediaParser()
		mc, err := jmp.ParseBytes(im_body)

		if err != nil {
			return nil, "", nil, err
		}

		jpg_ifd, _, err := mc.Exif()

		if err != nil {
			slog.Debug("Failed to derive EXIF", "error", err)
		} else {
			ifd = jpg_ifd
		}

	case "png":

		mp := pngstructure.NewPngMediaParser()

		mc, err := mp.ParseBytes(im_body)

		if err != nil {
			return nil, "", nil, err
		}

		png_ifd, _, err := mc.Exif()

		if err != nil {
			slog.Debug("Failed to derive EXIF", "error", err)
		} else {
			ifd = png_ifd
		}

	case "tiff":

		mp := tiffstructure.NewTiffMediaParser()

		mc, err := mp.ParseBytes(im_body)

		if err != nil {
			return nil, "", nil, err
		}

		tiff_ifd, _, err := mc.Exif()

		if err != nil {
			slog.Debug("Failed to derive EXIF", "error", err)
		} else {
			ifd = tiff_ifd
		}

	default:

		switch mtype.String() {
		case "image/heic":

			heic_im, err := ImageFromHEIC(im_body)

			if err != nil {
				return nil, "", nil, err
			}

			im = heic_im
			im_fmt = "heic"

			mp := heicexif.NewHeicExifMediaParser()
			mc, err := mp.ParseBytes(im_body)

			if err != nil {
				return nil, "", nil, err
			}

			heic_ifd, _, err := mc.Exif()

			if err != nil {
				slog.Debug("Failed to derive EXIF", "error", err)
			} else {
				ifd = heic_ifd
			}

			// Note: We are NOT removing or updating the Orientation tag
			// (which is assigned but incorrect) in libheif because I can
			// not figure out hwo to do that using the dsoprea packages
			// without causing everything to panic later in the code. Instead
			// we are accounting for this in RotateFromOrientation
			// https://github.com/strukturag/libheif/issues/227

		default:
			return nil, "", nil, fmt.Errorf("Unsupported media type")
		}
	}

	if opts.Rotate {

		_, r_im, err := rotateFromOrientation(ctx, im, mtype, ifd)

		if err != nil {
			return nil, "", nil, fmt.Errorf("Failed to rotate image, %w", err)
		}

		im = r_im
	}

	return im, mtype.String(), ifd, nil
}

func rotateFromOrientation(ctx context.Context, im image.Image, mtype *mimetype.MIME, ifd *exif.Ifd) (bool, image.Image, error) {

	if ifd == nil {
		return false, im, nil
	}

	// Ignore EXIF Orientation tags in libheif, kthxbye...
	// https://github.com/strukturag/libheif/issues/227

	if mtype.String() == "image/heic" {
		return true, im, nil
	}

	results, err := ifd.FindTagWithName("Orientation")

	if err != nil {

		if errors.Is(err, exif.ErrTagNotFound) {
			return false, im, nil
		}

		return false, nil, err
	}

	ite := results[0]
	orientation, err := ite.FormatFirst()

	if err != nil {
		return false, nil, err
	}

	// Rotate

	if orientation == "1" {
		return false, im, nil
	}

	r_im, err := rotate.RotateImageWithOrientation(ctx, im, orientation)

	if err != nil {
		return false, nil, err
	}

	return true, r_im, nil
}
