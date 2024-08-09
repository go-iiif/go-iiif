// package update provides methods for updating EXIF data in JPEG files.
package update

import (
	"fmt"
	"github.com/dsoprea/go-exif/v3"
	"github.com/dsoprea/go-exif/v3/common"
	"github.com/dsoprea/go-jpeg-image-structure/v2"
	"github.com/sfomuseum/go-exif-update/tags"
	"io"
	_ "log"
	"strconv"
)

var ti *exif.TagIndex

var tag_paths []*exifcommon.IfdIdentity

func init() {

	ti = exif.NewTagIndex()

	// https://github.com/dsoprea/go-exif/blob/de2141190595193aa097a2bf3205ba0cf76dc14b/tags_data.go

	tag_paths = []*exifcommon.IfdIdentity{
		exifcommon.IfdStandardIfdIdentity,
		exifcommon.IfdExifStandardIfdIdentity,
		exifcommon.IfdExifIopStandardIfdIdentity,
		exifcommon.IfdGpsInfoStandardIfdIdentity,
		exifcommon.Ifd1StandardIfdIdentity,
	}

}

// PrepareAndUpdateExif attempts to prepare and translate EXIF properties in defined in exif_props
// in to their appropriate dsoprea/go-exif types and format and then updates the EXIF data encoded in
// r and writes that data to wr. This method also supports a handful of custom properties that are
// prefixed with X- and used to populate known EXIF properties. These are:
//   - X-Latitude and X-Longitude, which convert and assign decimal latitude and longitude coordinates
//     in to their GPSLatitude/Longitude and GPSLatitude/LongitudeRef EXIF properties.
func PrepareAndUpdateExif(r io.Reader, wr io.Writer, exif_props map[string]interface{}) error {

	prepared := make(map[string]interface{})

	var lat float64
	var lon float64

	x_lat, has_lat := exif_props["X-Latitude"]
	x_lon, has_lon := exif_props["X-Longitude"]

	if has_lat && !has_lon {
		return fmt.Errorf("Missing X-Longitude property (X-Latitude is set)")
	}

	if has_lon && !has_lat {
		return fmt.Errorf("Missing X-Latitude property (X-Longitude is set)")
	}

	if has_lat && has_lon {

		switch x_lat.(type) {
		case float64:
			lat = x_lat.(float64)
		case string:

			l, err := strconv.ParseFloat(x_lat.(string), 64)

			if err != nil {
				return err
			}

			lat = l
		default:
			return fmt.Errorf("Invalid type for latitude, %T", lat)
		}

		switch x_lon.(type) {
		case float64:
			lon = x_lon.(float64)
		case string:

			l, err := strconv.ParseFloat(x_lon.(string), 64)

			if err != nil {
				return err
			}

			lon = l

		default:
			return fmt.Errorf("Invalid type for longitude")
		}

		err := AppendGPSPropertiesWithLatitudeAndLongitude(prepared, lat, lon)

		if err != nil {
			return err
		}

		delete(exif_props, "X-Latitude")
		delete(exif_props, "X-Longitude")
	}

	for k, v := range exif_props {

		str_v := fmt.Sprintf("%v", v)
		v2, err := PrepareTag(k, str_v)

		if err != nil {
			return fmt.Errorf("Failed to prepare tag '%s', %v", k, err)
		}

		prepared[k] = v2
	}

	return UpdateExif(r, wr, prepared)
}

// This is really nothing more than a thin wrapper around the example code in
// dsoprea's go-jpeg-image-structure package.

// UpdateExif updates the EXIF data encoded in r and writes that data to wr.
func UpdateExif(r io.Reader, wr io.Writer, exif_props map[string]interface{}) error {

	img_data, err := io.ReadAll(r)

	// https://pkg.go.dev/github.com/dsoprea/go-jpeg-image-structure/v2?utm_source=godoc#example-SegmentList.SetExif

	jmp := jpegstructure.NewJpegMediaParser()

	intfc, err := jmp.ParseBytes(img_data)

	if err != nil {
		return err
	}

	sl := intfc.(*jpegstructure.SegmentList)

	rootIb, err := sl.ConstructExifBuilder()

	if err != nil {
		return err
	}

	for k, v := range exif_props {

		ok, err := tags.IsSupported(k)

		if err != nil {
			return err
		}

		if !ok {
			return fmt.Errorf("Tag '%s' is not supported at this time", k)
		}

		id, _, err := GetIndexedTagFromName(k)

		if err != nil {
			return err
		}

		err = setExifTag(rootIb, id.UnindexedString(), k, v)

		if err != nil {
			return err
		}

	}

	// Update the exif segment.

	err = sl.SetExif(rootIb)

	if err != nil {
		return err
	}

	return sl.Write(wr)
}

// Return the *exifcommon.IfdIdentity and *exif.IndexedTag instances associated
// with a given EXIF string tag name.
func GetIndexedTagFromName(k string) (*exifcommon.IfdIdentity, *exif.IndexedTag, error) {

	for _, id := range tag_paths {

		t, err := ti.GetWithName(id, k)

		if err != nil {
			continue
		}

		return id, t, nil
	}

	return nil, nil, fmt.Errorf("Unrecognized tag, %s", k)
}

// Cribbed from https://github.com/dsoprea/go-exif/issues/11

func setExifTag(rootIB *exif.IfdBuilder, ifdPath string, tagName string, tagValue interface{}) error {

	// log.Printf("setTag(): ifdPath: %v, tagName: %v, tagValue: %v", ifdPath, tagName, tagValue)

	ifdIb, err := exif.GetOrCreateIbFromRootIb(rootIB, ifdPath)

	if err != nil {
		return fmt.Errorf("Failed to get or create IB for %s: %v", ifdPath, err)
	}

	err = ifdIb.SetStandardWithName(tagName, tagValue)

	if err != nil {
		return fmt.Errorf("failed to set %s tag: %v", tagName, err)
	}

	return nil
}
