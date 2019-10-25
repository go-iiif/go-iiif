package image

import (
	"errors"
	"strconv"
	"strings"
)

// see this? it's a iiifimage.Image not a (Go) image.Image

func ApplyCustomTransformations(t *Transformation, im Image) error {

	// None of what follows is part of the IIIF spec so it's not clear
	// to me yet how to make this in to a sane interface. For the time
	// being since there is only lipvips we'll just take the opportunity
	// to think about it... (20160917/thisisaaronland)

	// Also note the way we are diligently setting in `im.isgif` in each
	// of the features below. That's because this is a img/libvips-ism
	// and we assume that any of these can encode GIFs because pure-Go and
	// the rest of the code does need to know about it...
	// (20160922/thisisaaronland)

	if t.Quality == "dither" {

		err := DitherImage(im)

		if err != nil {
			return err
		}

	} else if strings.HasPrefix(t.Quality, "primitive:") {

		/*

			    "features": {
				"append": {
				    "quality": {
					"primitive": { "primitive": "dither", "required": false, "supported": true, "match": "^primitive:\\d,\\d+,\\d+$" }
				    }
				}
			    },

		*/

		fi, err := t.FormatInstructions(im)

		if err != nil {
			return err
		}

		parts := strings.Split(t.Quality, ":")
		parts = strings.Split(parts[1], ",")

		mode, err := strconv.Atoi(parts[0])

		if err != nil {
			return err
		}

		iters, err := strconv.Atoi(parts[1])

		if err != nil {
			return err
		}

		max_iters := 40 // FIX ME... config.Primitive.MaxIterations

		if max_iters > 0 && iters > max_iters {
			return errors.New("Invalid primitive iterations")
		}

		alpha, err := strconv.Atoi(parts[2])

		if err != nil {
			return err
		}

		if alpha > 255 {
			return errors.New("Invalid primitive alpha")
		}

		animated := false

		if fi.Format == "gif" {
			animated = true
		}

		opts := PrimitiveOptions{
			Alpha:      alpha,
			Mode:       mode,
			Iterations: iters,
			Size:       0,
			Animated:   animated,
		}

		err = PrimitiveImage(im, opts)

		if err != nil {
			return err
		}

		/*
			if fi.Format == "gif" {
				im.isgif = true
			}
		*/

	}

	return nil
}
