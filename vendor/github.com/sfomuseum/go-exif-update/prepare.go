package update

import (
	"fmt"
	"github.com/dsoprea/go-exif/v3"
	"github.com/dsoprea/go-exif/v3/common"
	_ "log"
	"math"
)

// Translate a decimal latitude value into a list of exifcommon.Rational instances
// (which in turn are an  an exifcommon.EncodedData instance)
func PrepareDecimalGPSLatitudeTag(lat float64) (interface{}, error) {

	var orientation string

	if lat > 0 {
		orientation = "N"
	} else {
		orientation = "S"
	}

	d, m, s := toDMS(lat)

	exif_gps := exif.GpsDegrees{
		Orientation: orientation[0],
		Degrees:     float64(d),
		Minutes:     float64(m),
		Seconds:     s,
	}

	return exif_gps.Raw(), nil
}

func PrepareDecimalGPSLatitudeRefTag(lat float64) (interface{}, error) {

	var orientation string

	if lat > 0 {
		orientation = "N"
	} else {
		orientation = "S"
	}

	return orientation, nil
}

// Translate a decimal longitude value into a list of exifcommon.Rational instances
// (which in turn are an  an exifcommon.EncodedData instance)
func PrepareDecimalGPSLongitudeTag(lon float64) (interface{}, error) {

	var orientation string

	if lon > 0 {
		orientation = "E"
	} else {
		orientation = "W"
	}

	d, m, s := toDMS(lon)

	exif_gps := exif.GpsDegrees{
		Orientation: orientation[0],
		Degrees:     float64(d),
		Minutes:     float64(m),
		Seconds:     s,
	}

	return exif_gps.Raw(), nil
}

func PrepareDecimalGPSLongitudeRefTag(lon float64) (interface{}, error) {

	var orientation string

	if lon > 0 {
		orientation = "E"
	} else {
		orientation = "W"
	}

	return orientation, nil
}

// PrepareTag attempts to translate and prepare an EXIF tag value into an exifcommon.EncodedData
// instance.
func PrepareTag(k string, v string) (interface{}, error) {

	_, t, err := GetIndexedTagFromName(k)

	if err != nil {
		return nil, fmt.Errorf("Failed to get indexed tag for '%s', %v", k, err)
	}

	for _, pr := range t.SupportedTypes {

		v2, err := exifcommon.TranslateStringToType(pr, v)

		if err != nil {
			return nil, err
		}

		// https://github.com/dsoprea/go-exif/blob/db167117f4830a268022c953f0f521fcc83d031e/v3/common/value_encoder.go#L225

		switch v2.(type) {
		case []byte, string:
			return v2, nil
		case exifcommon.Rational:
			return []exifcommon.Rational{v2.(exifcommon.Rational)}, nil
		case exifcommon.SignedRational:
			return []exifcommon.SignedRational{v2.(exifcommon.SignedRational)}, nil
		default:
			return nil, fmt.Errorf("Unsupported type '%T'", v2)
		}
	}

	return nil, fmt.Errorf("Unsupported type(s).")
}

// Cribbed from https://github.com/go-spatial/geom/blob/0e06498b336286332f6742fb898fb7118ede4005/planar/coord/coord.go

// Convert the given lng or lat value to the degree minute seconds values
func toDMS(v float64) (d int64, m int64, s float64) {
	var frac float64
	df, frac := math.Modf(v)
	mf, frac := math.Modf(60.0 * frac)
	s = 60.0 * frac
	return int64(math.Abs(df)), int64(math.Abs(mf)), math.Abs(s)
}
