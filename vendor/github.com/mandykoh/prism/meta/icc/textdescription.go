package icc

import (
	"bytes"
	"fmt"
	"github.com/mandykoh/prism/meta/binary"
)

type TextDescription struct {
	ASCII string
}

func parseTextDescription(data []byte) (TextDescription, error) {
	desc := TextDescription{}

	reader := bytes.NewReader(data)

	sig, err := binary.ReadU32Big(reader)
	if err != nil {
		return desc, err
	}
	if s := Signature(sig); s != DescSignature {
		return desc, fmt.Errorf("expected %v but got %v", DescSignature, s)
	}

	// Reserved field
	_, err = binary.ReadU32Big(reader)
	if err != nil {
		return desc, err
	}

	asciiCount, err := binary.ReadU32Big(reader)
	if err != nil {
		return desc, err
	}

	asciiBytes := make([]byte, asciiCount-1)
	for i := 0; i < len(asciiBytes); i++ {
		asciiBytes[i], err = reader.ReadByte()
		if err != nil {
			return desc, err
		}
	}

	// Skip terminating null
	_, err = reader.ReadByte()
	if err != nil {
		return desc, err
	}

	desc.ASCII = string(asciiBytes)

	return desc, nil
}
