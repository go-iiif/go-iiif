package heicexif

import (
	"bytes"
	"io"
	"os"
	"sort"

	"github.com/dsoprea/go-exif/v3"
	"github.com/dsoprea/go-exif/v3/common"
	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/v2/data"
	"github.com/dsoprea/go-utility/v2/filesystem"
	"github.com/dsoprea/go-utility/v2/image"
	"go4.org/media/heif"
	"go4.org/media/heif/bmff"
)

const (
	isoFtypBox = "ftyp"

	// MinimumHeicStreamLengthForDetection describes the minimum amount of data
	// required to detect HEIC content.
	MinimumHeicStreamLengthForDetection = 1000

	// HeicMimeType describes the theoretical, official HEIC mime-type
	// (theoretical because it's not ubiquitously supported/known in tooling,
	// yet).
	HeicMimeType = "image/heif"
)

var (
	// HeicBrands is the list of MP4 brands that are associated with HEIC
	// streams.
	HeicBrands = sort.StringSlice{
		"mif1",
		"msf1",
		"heic",
	}
)

// HeicExifContext is just a simple type-alias for a byte-array.
type HeicExifContext []byte

// Exif parses the raw EXIF data and returns the root IFD.
func (hec HeicExifContext) Exif() (rootIfd *exif.Ifd, data []byte, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	im, err := exifcommon.NewIfdMappingWithStandard()
	log.PanicIf(err)

	ti := exif.NewTagIndex()

	_, index, err := exif.Collect(im, ti, hec)
	log.PanicIf(err)

	return index.RootIfd, hec, nil
}

// HeicExifMediaParser satisfies the `riimage.MediaParser` interface.
type HeicExifMediaParser struct {
}

// NewHeicExifMediaParser returns a new HeicExifMediaParser instance.
func NewHeicExifMediaParser() *HeicExifMediaParser {
	return new(HeicExifMediaParser)
}

// Parse parses a stream using an `io.ReadSeeker`.
func (hemp *HeicExifMediaParser) Parse(rs io.ReadSeeker, size int) (mc riimage.MediaContext, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	ra := rifs.NewReadSeekerToReaderAt(rs)
	hf := heif.Open(ra)

	rawExif, err := hf.EXIF()
	if err != nil {
		if err == heif.ErrNoEXIF {
			return nil, exif.ErrNoExif
		}

		log.Panic(err)
	}

	// NOTE(dustin): Kludge for broken third-party functionality. The EXIF data
	// doesn't actually start on the first byte.
	//
	// https://github.com/go4org/go4/issues/51
	rawExif, err = exif.SearchAndExtractExif(rawExif)
	log.PanicIf(err)

	return HeicExifContext(rawExif), nil
}

// ParseFile parses a stream using a file.
func (hemp *HeicExifMediaParser) ParseFile(filepath string) (mc riimage.MediaContext, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	f, err := os.Open(filepath)
	log.PanicIf(err)

	defer f.Close()

	mc, err = hemp.Parse(f, 0)
	if err != nil {
		if err == exif.ErrNoExif {
			return nil, exif.ErrNoExif
		}

		log.Panic(err)
	}

	return mc, nil
}

// ParseBytes parses a stream direct from bytes.
func (hemp *HeicExifMediaParser) ParseBytes(data []byte) (mc riimage.MediaContext, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	br := bytes.NewReader(data)

	mc, err = hemp.Parse(br, 0)
	if err != nil {
		if err == exif.ErrNoExif {
			return nil, exif.ErrNoExif
		}

		log.Panic(err)
	}

	return mc, nil
}

// LooksLikeFormat parses the data to determine if it's a compatible format.
func (hemp *HeicExifMediaParser) LooksLikeFormat(data []byte) bool {
	// We need to assert that there's enough data to parse. Since the box that
	// we care about is just about the first thing that is parsed, we shouldn't
	// require very much, and this should be universally the case.
	if len(data) < MinimumHeicStreamLengthForDetection {
		return false
	}

	br := bytes.NewReader(data)

	// Do a low-cost check for a mime-type, first. Currently, the mime-type for
	// HEIC is still largely unsupported.

	mimetype, err := ridata.GetMimetypeFromContent(br, MinimumHeicStreamLengthForDetection)
	log.PanicIf(err)

	if mimetype == HeicMimeType {
		return true
	}

	// We were unable to detect HEIC based on the mime-type. Try to parse
	// the first couple of boxes (we may not have enough data for more than
	// that, but it's not necessary).

	_, err = br.Seek(0, io.SeekStart)
	log.PanicIf(err)

	// We just need this to get us to the first box.
	heif.Open(br)

	r := bmff.NewReader(br)

	boxHeader, err := r.ReadBox()
	if err != nil {
		return false
	}

	box, err := boxHeader.Parse()
	if err != nil {
		return false
	}

	// It's an ISO 14496-12 stream.

	// Check if the first box is a "ftyp" type.
	if boxHeader.Type().String() != isoFtypBox {
		return false
	}

	ftyp, ok := box.(*bmff.FileTypeBox)
	if ok != true {
		// Box identifies as 'ftyp' but it parses-out to the wrong type.
		return false
	}

	i := HeicBrands.Search(ftyp.MajorBrand)
	if i < len(HeicBrands) && HeicBrands[i] == ftyp.MajorBrand {
		return true
	}

	// If we get here, We got a stream that's eerily similar to an HEIC stream
	// without being one.

	return false
}

var (
	// Enforce that `HeicExifMediaParser` looks like a `riimage.MediaParser`.
	_ riimage.MediaParser = new(HeicExifMediaParser)
)

func init() {
	HeicBrands.Sort()
}
