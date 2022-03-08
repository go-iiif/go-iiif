package image

import (
	"errors"
	"fmt"
	_ "log"
	"strconv"
	"strings"
)

// see this? it's a Image not a (Go) image.Image

// Apply non-standard transformations to a go-iiif/image.Image instance.
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

	// This (splitting on ';' doesn't work yet because it fails on the
	// earlier pattern matching for the quality parameter. For example:
	// 2021/07/04 10:35:49 Invalid IIIF 2.1 feature property quality crisp:2.0,1.0,0.025;dither
	// https://github.com/go-iiif/go-iiif/issues/89

	for _, q := range strings.Split(t.Quality, ";") {

		if q == "dither" {

			err := DitherImage(im)

			if err != nil {
				return err
			}

		} else if q == "crisp" || strings.HasPrefix(q, "crisp:") {

			/*

				    "features": {
					"append": {
					    "quality": {
						"crisp": { "syntax": "crisp", "required": false, "supported": true, "match": "^crisp(?:\\:(\\d+\\.\\d+),(\\d+\\.\\d+),(\\d+\\.\\d+))?$" }
					    }
					}
				    }

			*/

			opts := DefaultCrispImageOptions()

			parts := strings.Split(q, ":")

			if len(parts) == 2 {

				str_opts := strings.Split(parts[1], ",")

				if len(str_opts) != 3 {
					return fmt.Errorf("Invalid crisp parameters '%s'", parts[1])
				}

				r, err := strconv.ParseFloat(str_opts[0], 64)

				if err != nil {
					return fmt.Errorf("Invalid radius parameter, %w", err)
				}

				a, err := strconv.ParseFloat(str_opts[1], 64)

				if err != nil {
					return fmt.Errorf("Invalid amount parameter, %w", err)
				}

				m, err := strconv.ParseFloat(str_opts[2], 64)

				if err != nil {
					return fmt.Errorf("Invalid median parameter, %w", err)
				}

				opts.Radius = r
				opts.Amount = a
				opts.Median = m
			}

			err := CrispImage(im, opts)

			if err != nil {
				return fmt.Errorf("Failed to crispen image, %w", err)
			}

		} else if strings.HasPrefix(q, "primitive:") {

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

	}

	return nil
}
