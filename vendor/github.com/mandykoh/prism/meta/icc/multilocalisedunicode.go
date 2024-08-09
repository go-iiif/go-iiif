package icc

import (
	"bytes"
	"fmt"
	"github.com/mandykoh/prism/meta/binary"
	"unicode/utf16"
)

type MultiLocalisedUnicode struct {
	entriesByLanguageCountry map[[2]byte]map[[2]byte]string
}

func (mluc *MultiLocalisedUnicode) getAnyString() string {
	for _, country := range mluc.entriesByLanguageCountry {
		for _, s := range country {
			return s
		}
	}
	return ""
}

func (mluc *MultiLocalisedUnicode) getString(language [2]byte, country [2]byte) string {
	countries, ok := mluc.entriesByLanguageCountry[language]
	if !ok {
		return ""
	}

	return countries[country]
}

func (mluc *MultiLocalisedUnicode) getStringForLanguage(language [2]byte) string {
	for _, s := range mluc.entriesByLanguageCountry[language] {
		return s
	}
	return ""
}

func (mluc *MultiLocalisedUnicode) setString(language [2]byte, country [2]byte, text string) {
	countries, ok := mluc.entriesByLanguageCountry[language]
	if !ok {
		countries = map[[2]byte]string{
			country: text,
		}
		mluc.entriesByLanguageCountry[language] = countries

	} else {
		countries[country] = text
	}
}

func parseMultiLocalisedUnicode(data []byte) (MultiLocalisedUnicode, error) {
	result := MultiLocalisedUnicode{
		entriesByLanguageCountry: make(map[[2]byte]map[[2]byte]string),
	}

	reader := bytes.NewReader(data)

	sig, err := binary.ReadU32Big(reader)
	if err != nil {
		return result, err
	}
	if s := Signature(sig); s != MultiLocalisedUnicodeSignature {
		return result, fmt.Errorf("expected %v but got %v", MultiLocalisedUnicodeSignature, s)
	}

	// Reserved field
	_, err = binary.ReadU32Big(reader)
	if err != nil {
		return result, err
	}

	recordCount, err := binary.ReadU32Big(reader)
	if err != nil {
		return result, err
	}

	recordSize, err := binary.ReadU32Big(reader)
	if err != nil {
		return result, err
	}

	for i := uint32(0); i < recordCount; i++ {

		language := [2]byte{}
		n, err := reader.Read(language[:])
		if err != nil {
			return result, err
		}
		if n < len(language) {
			return result, fmt.Errorf("unexpected eof when reading language code")
		}

		country := [2]byte{}
		n, err = reader.Read(country[:])
		if err != nil {
			return result, err
		}
		if n < len(country) {
			return result, fmt.Errorf("unexpected eof when reading country code")
		}

		stringLength, err := binary.ReadU32Big(reader)
		if err != nil {
			return result, err
		}

		stringOffset, err := binary.ReadU32Big(reader)
		if err != nil {
			return result, err
		}

		if uint64(stringOffset+stringLength) > uint64(len(data)) {
			return result, fmt.Errorf("record exceeds tag data length")
		}

		recordStringBytes := data[stringOffset : stringOffset+stringLength]
		recordStringUTF16 := make([]uint16, len(recordStringBytes)/2)
		for j := 0; j < len(recordStringUTF16); j++ {
			recordStringUTF16[j], err = binary.ReadU16Big(reader)
			if err != nil {
				return result, err
			}
		}
		result.setString(language, country, string(utf16.Decode(recordStringUTF16)))

		// Skip to next record
		for j := uint32(12); j < recordSize; j++ {
			_, err := reader.ReadByte()
			if err != nil {
				return result, err
			}
		}
	}

	return result, nil
}

type languageCountry struct {
	language [2]byte
	country  [2]byte
}

func (lc languageCountry) String() string {
	return fmt.Sprintf("%c%c_%c%c", lc.language[0], lc.language[1], lc.country[0], lc.country[1])
}
