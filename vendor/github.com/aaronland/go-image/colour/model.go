package colour

import (
	"fmt"
	"io"
	"log/slog"
)

// UNKNOWN_MODEL defines an unknown or unspecified colour model.
const UNKNOWN_MODEL string = "unknown"

// SRGB_MODEL defines the sRGB colour space/model.
const SRGB_MODEL string = "sRGB"

// DISPLAYP3_MODEL defines the Apple DisplayP3 colour model
const DISPLAYP3_MODEL string = "DisplayP3"

// ARGB_MODEL defines the Adobe RGB colour model.
const ARGB_MODEL string = "Adobe RGB"

const (
	// The unknown or undefined colour model.
	UnknownModel Model = iota
	// The colour model corresponding to the sRGB colour space.
	SRGBModel
	// The colour model corresponding to the Adobe RGB colour space.
	AdobeRGBModel
	// The colour model corresponding to the Apple Display P3 colour space.
	AppleDisplayP3Model
)

// type Model defines an internal catalog of colour "models" for Go language image pixels.
// Models are derived from ICC profile descriptions and EXIF ColorSpace definitions in that
// order. The goal of these models is to provide hints to other applications which may need
// or want to recast the colour space of individual pixels in Go language `image.Image` instances.
// Currently there are only four "known" models: Apple's Display P3, Adobe RGB, sRGB and unknown.
type Model uint8

func (p Model) String() string {

	switch p {
	case SRGBModel:
		return SRGB_MODEL
	case AdobeRGBModel:
		return ARGB_MODEL
	case AppleDisplayP3Model:
		return DISPLAYP3_MODEL
	default:
		return UNKNOWN_MODEL
	}
}

// StringToModel returns the `Model` instance matching 'str_model'.
func StringToModel(str_model string) Model {

	switch str_model {
	case SRGB_MODEL:
		return SRGBModel
	case ARGB_MODEL:
		return AdobeRGBModel
	case DISPLAYP3_MODEL:
		return AppleDisplayP3Model
	default:
		return UnknownModel
	}
}

// DeriveModel attempts to derive a `Model` instance from the body of 'r' by checking
// for an ICC profile description or a EXIF ColorSpace definition, in that order.
func DeriveModel(r io.ReadSeeker) (Model, error) {

	pr, _ := ICCProfileDescription(r)

	if pr != "" {

		switch pr {
		case ICC_DISPLAY_P3:
			return AppleDisplayP3Model, nil
		case ICC_EPSON_RGB_G18:
			return SRGBModel, nil
		case ICC_SRGB_21:
			return SRGBModel, nil
		case ICC_CAMERA_RGB:
			return SRGBModel, nil
		case ICC_ADOBE_RGB_1998:
			return AdobeRGBModel, nil
		default:
			slog.Warn("Unknown or unsupported ICC profile", "description", pr)
		}
	}

	_, err := r.Seek(0, 0)

	if err != nil {
		return UnknownModel, fmt.Errorf("Failed to rewind reader after checking ICC profile, %w", err)
	}

	colorspace, err := ColorSpace(r)

	if err != nil {
		// slog.Warn("Failed to derive colorspace, returning unknown", "error", err)
		return UnknownModel, nil
	}

	switch colorspace {
	case COLORSPACE_SRGB:
		return SRGBModel, nil
	case COLORSPACE_ARGB:
		return AdobeRGBModel, nil
	case COLORSPACE_UNKNOWN:
		return UnknownModel, nil
	default:
		slog.Warn("Unknown or unsuported colorspace, returning unknown", "colorspace", colorspace)
		return UnknownModel, nil
	}
}
