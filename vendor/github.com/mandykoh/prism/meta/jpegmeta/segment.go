package jpegmeta

import (
	"fmt"
	"github.com/mandykoh/prism/meta/binary"
	"io"
)

var invalidSegment = segment{Marker: invalidMarker}

type segment struct {
	Marker marker
	Data   []byte
}

func makeSegment(markerType byte, r io.ByteReader) (segment, error) {
	m, err := makeMarker(markerType, r)
	return segment{Marker: m}, err
}

func readSegment(r binary.Reader) (segment, error) {
	m, err := readMarker(r)
	if err != nil {
		return invalidSegment, err
	}

	seg := segment{
		Marker: m,
	}
	if m.DataLength > 0 {
		seg.Data = make([]byte, m.DataLength)
	}

	n, err := io.ReadFull(r, seg.Data)
	if err != nil {
		return invalidSegment, err
	}
	if n < len(seg.Data) {
		return invalidSegment, fmt.Errorf("expected %d bytes of segment data but read %d", m.DataLength, n)
	}

	return seg, nil
}
