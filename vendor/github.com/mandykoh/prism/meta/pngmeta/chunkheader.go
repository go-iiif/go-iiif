package pngmeta

import (
	"fmt"
	"github.com/mandykoh/prism/meta/binary"
)

type chunkHeader struct {
	Length    uint32
	ChunkType [4]byte
}

func (ch chunkHeader) String() string {
	return fmt.Sprintf("%c%c%c%c(%d)", ch.ChunkType[0], ch.ChunkType[1], ch.ChunkType[2], ch.ChunkType[3], ch.Length)
}

func readChunkHeader(r binary.Reader) (chunkHeader, error) {
	ch := chunkHeader{}

	var err error
	ch.Length, err = binary.ReadU32Big(r)
	if err != nil {
		return ch, err
	}

	bytesRead, err := r.Read(ch.ChunkType[:])
	if err != nil {
		return ch, err
	}
	if bytesRead != len(ch.ChunkType) {
		return ch, fmt.Errorf("unexpected EOF reading chunk type")
	}

	return ch, nil
}
