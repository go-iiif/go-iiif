package exif

/*

 	(the value of 0x2 is not standard EXIF. Instead, an Adobe RGB image is indicated by "Uncalibrated" with an InteropIndex of "R03". The values 0xfffd and 0xfffe are also non-standard, and are used by some Sony cameras)
0x1 = sRGB
0x2 = Adobe RGB
0xfffd = Wide Gamut RGB
0xfffe = ICC Profile
0xffff = Uncalibrated

// https://www.exiftool.org/TagNames/EXIF.html

// https://exiv2.org/makernote.html

*/

import (
	"fmt"
	"io"
	"log/slog"
	"slices"

	go_exif "github.com/dsoprea/go-exif/v3"
	go_exifcommon "github.com/dsoprea/go-exif/v3/common"
)

func TagIndex(r io.Reader) (*go_exif.IfdIndex, error) {

	rawExif, err := go_exif.SearchAndExtractExifWithReader(r)

	if err != nil {
		return nil, fmt.Errorf("Failed to extract EXIF data, %w", err)
	}

	im, err := go_exifcommon.NewIfdMappingWithStandard()

	if err != nil {
		return nil, fmt.Errorf("Failed to create ifd mapping, %w", err)
	}

	ti := go_exif.NewTagIndex()

	_, index, err := go_exif.Collect(im, ti, rawExif)

	if err != nil {
		return nil, fmt.Errorf("Failed to collect EXIF data, %w", err)
	}

	return &index, nil
}

func TagValue(r io.Reader, tag_name string, ifds_name string) (*go_exif.IfdTagEntry, error) {

	index, err := TagIndex(r)

	if err != nil {
		return nil, err
	}

	return TagValueWithIndex(index, tag_name, ifds_name)
}

func TagValueWithIndex(index *go_exif.IfdIndex, tag_name string, ifds_name string) (*go_exif.IfdTagEntry, error) {

	rsp, err := TagValuesWithIndex(index, tag_name, ifds_name)

	if err != nil {
		return nil, err
	}

	v, exists := rsp[ifds_name]

	if !exists {
		return nil, fmt.Errorf("Tag not found")
	}

	return v, nil
}

func TagValues(r io.Reader, tag_name string, require_ifd ...string) (map[string]*go_exif.IfdTagEntry, error) {

	index, err := TagIndex(r)

	if err != nil {
		return nil, err
	}

	return TagValuesWithIndex(index, tag_name, require_ifd...)
}

func TagValuesWithIndex(index *go_exif.IfdIndex, tag_name string, require_ifd ...string) (map[string]*go_exif.IfdTagEntry, error) {

	results := make(map[string]*go_exif.IfdTagEntry)

	for _, ifd := range index.Ifds {

		ident := ifd.IfdIdentity()

		logger := slog.Default()
		logger = logger.With("identity", ident.String())
		logger = logger.With("tag", tag_name)

		// logger.Debug("Check IDF for tag")

		if len(require_ifd) > 0 && !slices.Contains(require_ifd, ident.String()) {
			// logger.Debug("IDF not contained in required IFD list")
			continue
		}

		rsp, err := ifd.FindTagWithName(tag_name)

		if err != nil {
			logger.Debug("Failed to locate tag", "error", err)
			continue
		}

		switch len(rsp) {
		case 0:
			continue
		case 1:
			results[ident.String()] = rsp[0]
		default:
			logger.Warn("Multiple results for tag", "count", len(rsp))
			continue
		}
	}

	return results, nil
}
