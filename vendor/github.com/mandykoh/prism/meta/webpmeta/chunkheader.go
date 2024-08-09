package webpmeta

import (
	"fmt"

	"github.com/mandykoh/prism/meta/binary"
)

type chunkHeader struct {
	ChunkType [4]byte
	Length    uint32
}

func (ch chunkHeader) String() string {
	return fmt.Sprintf("%c%c%c%c(%d)", ch.ChunkType[0], ch.ChunkType[1], ch.ChunkType[2], ch.ChunkType[3], ch.Length)
}

func readChunkHeader(r binary.Reader) (ch chunkHeader, err error) {
	bytesRead, err := r.Read(ch.ChunkType[:])
	if err != nil {
		return ch, err
	}
	if bytesRead != len(ch.ChunkType) {
		return ch, fmt.Errorf("unexpected EOF reading chunk type")
	}

	ch.Length, err = binary.ReadU32Little(r)
	if err != nil {
		return ch, err
	}

	return ch, nil
}
