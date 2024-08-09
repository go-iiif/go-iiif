package icc

import (
	"bytes"
	"fmt"
	"github.com/mandykoh/prism/meta/binary"
)

type TagTable struct {
	entries map[Signature][]byte
}

func (t *TagTable) add(sig Signature, data []byte) {
	t.entries[sig] = data
}

func (t *TagTable) getProfileDescription() (string, error) {
	data := t.entries[DescSignature]

	sig, err := binary.ReadU32Big(bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	switch Signature(sig) {

	case DescSignature:
		desc, err := parseTextDescription(data)
		if err != nil {
			return "", err
		}
		return desc.ASCII, nil

	case MultiLocalisedUnicodeSignature:
		mluc, err := parseMultiLocalisedUnicode(data)
		if err != nil {
			return "", err
		}
		if enUS := mluc.getStringForLanguage([2]byte{'e', 'n'}); enUS != "" {
			return enUS, nil
		}
		return mluc.getAnyString(), nil

	default:
		return "", fmt.Errorf("unknown profile description type (%v)", Signature(sig))
	}
}

func emptyTagTable() TagTable {
	return TagTable{
		entries: make(map[Signature][]byte),
	}
}
