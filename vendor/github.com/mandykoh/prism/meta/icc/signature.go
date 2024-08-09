package icc

import "fmt"

type Signature uint32

const (
	ProfileFileSignature           Signature = 0x61637370 // 'acsp'
	DescSignature                  Signature = 0x64657363 // 'desc'
	MultiLocalisedUnicodeSignature Signature = 0x6D6C7563 // 'mluc'
)

func (s Signature) String() string {
	v := [4]byte{
		maskNull(byte((s >> 24) & 0xff)),
		maskNull(byte((s >> 16) & 0xff)),
		maskNull(byte((s >> 8) & 0xff)),
		maskNull(byte(s & 0xff)),
	}

	return fmt.Sprintf("'%c%c%c%c'", v[0], v[1], v[2], v[3])
}

func maskNull(b byte) byte {
	if b == 0 {
		return ' '
	}
	return b
}
