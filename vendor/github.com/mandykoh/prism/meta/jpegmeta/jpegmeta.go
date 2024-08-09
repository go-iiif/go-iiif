package jpegmeta

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/mandykoh/prism/meta"
	"github.com/mandykoh/prism/meta/binary"
	"io"
)

// Format specifies the image format handled by this package
var Format = meta.ImageFormat("JPEG")

var iccProfileIdentifier = []byte("ICC_PROFILE\x00")

// Load loads the metadata for a JPEG image stream.
//
// Only as much of the stream is consumed as necessary to extract the metadata;
// the returned stream contains a buffered copy of the consumed data such that
// reading from it will produce the same results as fully reading the input
// stream. This provides a convenient way to load the full image after loading
// the metadata.
//
// An error is returned if basic metadata could not be extracted. The returned
// stream still provides the full image data.
func Load(r io.Reader) (md *meta.Data, imgStream io.Reader, err error) {
	rewindBuffer := &bytes.Buffer{}
	tee := io.TeeReader(r, rewindBuffer)
	md, err = extractMetadata(bufio.NewReader(tee))
	return md, io.MultiReader(rewindBuffer, r), err
}

func extractMetadata(r binary.Reader) (md *meta.Data, err error) {
	metadataExtracted := false
	md = &meta.Data{Format: Format}
	segReader := NewSegmentReader(r)

	defer func() {
		if r := recover(); r != nil {
			if !metadataExtracted {
				md = nil
			}
			err = fmt.Errorf("panic while extracting image metadata: %v", r)
		}
	}()

	var iccProfileChunks [][]byte
	var iccProfileChunksExtracted int

	allMetadataExtracted := func() bool {
		return metadataExtracted &&
			iccProfileChunks != nil &&
			iccProfileChunksExtracted == len(iccProfileChunks)
	}

	soiSegment, err := segReader.ReadSegment()
	if err != nil {
		return nil, err
	}
	if soiSegment.Marker.Type != markerTypeStartOfImage {
		return nil, fmt.Errorf("stream does not begin with start-of-image")
	}

parseSegments:
	for {
		segment, err := segReader.ReadSegment()
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("unexpected EOF")
			}
			return nil, err
		}

		switch segment.Marker.Type {

		case markerTypeStartOfFrameBaseline,
			markerTypeStartOfFrameProgressive:
			md.BitsPerComponent = uint32(segment.Data[0])
			md.PixelHeight = uint32(segment.Data[1])<<8 | uint32(segment.Data[2])
			md.PixelWidth = uint32(segment.Data[3])<<8 | uint32(segment.Data[4])
			metadataExtracted = true

			if allMetadataExtracted() {
				break parseSegments
			}

		case markerTypeStartOfScan,
			markerTypeEndOfImage:
			break parseSegments

		case markerTypeApp2:
			if len(segment.Data) < len(iccProfileIdentifier)+2 {
				continue
			}

			for i := range iccProfileIdentifier {
				if segment.Data[i] != iccProfileIdentifier[i] {
					continue parseSegments
				}
			}

			iccData, iccErr := md.ICCProfileData()
			if iccData != nil || iccErr != nil {
				continue
			}

			chunkTotal := segment.Data[len(iccProfileIdentifier)+1]
			if iccProfileChunks == nil {
				iccProfileChunks = make([][]byte, chunkTotal)
			} else if int(chunkTotal) != len(iccProfileChunks) {
				md.SetICCProfileError(fmt.Errorf("inconsistent ICC profile chunk count"))
				continue
			}

			chunkNum := segment.Data[len(iccProfileIdentifier)]
			if chunkNum == 0 || int(chunkNum) > len(iccProfileChunks) {
				md.SetICCProfileError(fmt.Errorf("invalid ICC profile chunk number"))
				continue
			}
			if iccProfileChunks[chunkNum-1] != nil {
				md.SetICCProfileError(fmt.Errorf("duplicated ICC profile chunk"))
				continue
			}
			iccProfileChunksExtracted++
			iccProfileChunks[chunkNum-1] = segment.Data[len(iccProfileIdentifier)+2:]

			if allMetadataExtracted() {
				break parseSegments
			}
		}
	}

	if !metadataExtracted {
		return nil, fmt.Errorf("no metadata found")
	}

	// Incomplete or missing ICC profile
	if len(iccProfileChunks) != iccProfileChunksExtracted {
		_, iccErr := md.ICCProfileData()
		if iccErr == nil {
			md.SetICCProfileError(fmt.Errorf("incomplete ICC profile data"))
		}
		return md, nil
	}

	iccProfileData := bytes.Buffer{}
	for i := range iccProfileChunks {
		iccProfileData.Write(iccProfileChunks[i])
	}
	md.SetICCProfileData(iccProfileData.Bytes())

	return md, nil
}
