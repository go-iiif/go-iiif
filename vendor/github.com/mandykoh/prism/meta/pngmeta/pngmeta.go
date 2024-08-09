package pngmeta

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"github.com/mandykoh/prism/meta"
	"github.com/mandykoh/prism/meta/binary"
	"io"
	"strings"
)

// Format specifies the image format handled by this package
var Format = meta.ImageFormat("PNG")

var pngSignature = [8]byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A}

// Load loads the metadata for a PNG image stream.
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

	defer func() {
		if r := recover(); r != nil {
			if !metadataExtracted {
				md = nil
			}
			err = fmt.Errorf("panic while extracting image metadata: %v", r)
		}
	}()

	allMetadataExtracted := func() bool {
		iccData, iccErr := md.ICCProfileData()
		return metadataExtracted && (iccData != nil || iccErr != nil)
	}

	pngSig := [8]byte{}
	bytesRead, err := r.Read(pngSig[:])
	if err != nil {
		return nil, err
	}
	if bytesRead != len(pngSig) {
		return nil, fmt.Errorf("unexpected EOF reading PNG header")
	}
	if pngSig != pngSignature {
		return nil, fmt.Errorf("invalid PNG signature")
	}

parseChunks:
	for {
		ch, err := readChunkHeader(r)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		switch ch.ChunkType {

		case chunkTypeIHDR:
			md.PixelWidth, err = binary.ReadU32Big(r)
			if err != nil {
				return nil, err
			}

			md.PixelHeight, err = binary.ReadU32Big(r)
			if err != nil {
				return nil, err
			}

			bitDepth, err := r.ReadByte()
			if err != nil {
				return nil, err
			}
			md.BitsPerComponent = uint32(bitDepth)

			// Skip remainder of header
			for i := uint32(0); i < ch.Length-9; i++ {
				_, err := r.ReadByte()
				if err != nil {
					return nil, err
				}
			}

			// Skip chunk CRC
			_, err = binary.ReadU32Big(r)
			if err != nil {
				return nil, err
			}

			metadataExtracted = true

			if allMetadataExtracted() {
				break parseChunks
			}

		case chunkTypeiCCP:

			profileName := strings.Builder{}
			for i := 0; i < 80; i++ {
				b, err := r.ReadByte()
				if err != nil {
					return nil, err
				}
				if b == 0x00 {
					break
				}
				profileName.WriteByte(b)
			}
			if profileName.Len() > 79 {
				return nil, fmt.Errorf("null terminator not found reading ICC profile name")
			}

			compressionMethod, err := r.ReadByte()
			if err != nil {
				return nil, err
			}
			if compressionMethod != 0x00 {
				return nil, fmt.Errorf("unknown compression method (%d)", compressionMethod)
			}

			offset := uint32(profileName.Len() + 2)
			if offset >= ch.Length {
				return nil, fmt.Errorf("invalid ICC profile chunk length")
			}

			chunkData := make([]byte, ch.Length-offset)
			bytesRead, err := r.Read(chunkData)
			if err != nil {
				return nil, err
			}
			if bytesRead != len(chunkData) {
				return nil, fmt.Errorf("unexpected EOF reading ICC profile chunk")
			}

			// Skip chunk CRC
			_, err = binary.ReadU32Big(r)
			if err != nil {
				return nil, err
			}

			// Decompress ICC profile data
			zReader, err := zlib.NewReader(bytes.NewReader(chunkData))
			if err != nil {
				md.SetICCProfileError(err)
				break
			}
			profileData := &bytes.Buffer{}
			_, err = io.Copy(profileData, zReader)
			_ = zReader.Close()
			if err == nil {
				md.SetICCProfileData(profileData.Bytes())

				if allMetadataExtracted() {
					break parseChunks
				}
			} else {
				md.SetICCProfileError(err)
			}

		case chunkTypeIDAT, chunkTypeIEND:
			break parseChunks

		default:
			// Skip chunk data bytes
			for i := uint32(0); i < ch.Length; i++ {
				_, err := r.ReadByte()
				if err != nil {
					return nil, err
				}
			}

			// Skip chunk CRC
			_, err := binary.ReadU32Big(r)
			if err != nil {
				return nil, err
			}
		}
	}

	if !metadataExtracted {
		return nil, fmt.Errorf("no metadata found")
	}

	return md, nil
}
