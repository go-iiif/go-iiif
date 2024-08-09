package icc

import (
	"fmt"
	"github.com/mandykoh/prism/meta/binary"
	"time"
)

type ProfileReader struct {
	reader binary.Reader
}

func (pr *ProfileReader) ReadProfile() (p *Profile, err error) {
	defer func() {
		if r := recover(); r != nil {
			p = nil
			err = fmt.Errorf("panic while parsing ICC profile: %v", r)
		}
	}()

	profile := newProfile()

	err = pr.readHeader(&profile.Header)
	if err != nil {
		return nil, err
	}

	err = pr.readTagTable(&profile.TagTable)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (pr *ProfileReader) readDateTimeNumber() (result time.Time, err error) {
	year, err := binary.ReadU16Big(pr.reader)
	if err != nil {
		return
	}

	month, err := binary.ReadU16Big(pr.reader)
	if err != nil {
		return
	}

	day, err := binary.ReadU16Big(pr.reader)
	if err != nil {
		return
	}

	hour, err := binary.ReadU16Big(pr.reader)
	if err != nil {
		return
	}

	minute, err := binary.ReadU16Big(pr.reader)
	if err != nil {
		return
	}

	second, err := binary.ReadU16Big(pr.reader)
	if err != nil {
		return
	}

	return time.Date(int(year), time.Month(month), int(day), int(hour), int(minute), int(second), 0, time.UTC), nil
}

func (pr *ProfileReader) readHeader(header *Header) error {
	var err error

	header.ProfileSize, err = binary.ReadU32Big(pr.reader)
	if err != nil {
		return err
	}

	value, err := binary.ReadU32Big(pr.reader)
	if err != nil {
		return err
	}
	header.PreferredCMM = Signature(value)

	header.Version.Major, err = pr.reader.ReadByte()
	if err != nil {
		return err
	}
	header.Version.MinorAndRev, err = pr.reader.ReadByte()
	if err != nil {
		return err
	}

	// Reserved bytes in version field
	_, err = pr.reader.ReadByte()
	if err != nil {
		return err
	}
	_, err = pr.reader.ReadByte()
	if err != nil {
		return err
	}

	value, err = binary.ReadU32Big(pr.reader)
	if err != nil {
		return err
	}
	header.DeviceClass = DeviceClass(value)

	value, err = binary.ReadU32Big(pr.reader)
	if err != nil {
		return err
	}
	header.DataColorSpace = ColorSpace(value)

	value, err = binary.ReadU32Big(pr.reader)
	if err != nil {
		return err
	}
	header.ProfileConnectionSpace = ColorSpace(value)

	header.CreatedAt, err = pr.readDateTimeNumber()
	if err != nil {
		return err
	}

	value, err = binary.ReadU32Big(pr.reader)
	if err != nil {
		return err
	}
	if s := Signature(value); s != ProfileFileSignature {
		return fmt.Errorf("invalid profile file signature %s", s)
	}

	value, err = binary.ReadU32Big(pr.reader)
	if err != nil {
		return err
	}
	header.PrimaryPlatform = PrimaryPlatform(value)

	value, err = binary.ReadU32Big(pr.reader)
	if err != nil {
		return err
	}
	header.Embedded = (value >> 31) != 0
	header.DependsOnEmbeddedData = (value>>30)&1 != 0

	value, err = binary.ReadU32Big(pr.reader)
	if err != nil {
		return err
	}
	header.DeviceManufacturer = Signature(value)

	value, err = binary.ReadU32Big(pr.reader)
	if err != nil {
		return err
	}
	header.DeviceModel = Signature(value)

	header.DeviceAttributes, err = binary.ReadU64Big(pr.reader)
	if err != nil {
		return err
	}

	value, err = binary.ReadU32Big(pr.reader)
	if err != nil {
		return err
	}
	header.RenderingIntent = RenderingIntent(value)

	for i := 0; i < 3; i++ {
		header.PCSIlluminant[i], err = binary.ReadU32Big(pr.reader)
		if err != nil {
			return err
		}
	}

	value, err = binary.ReadU32Big(pr.reader)
	if err != nil {
		return err
	}
	header.ProfileCreator = Signature(value)

	bytesRead, err := pr.reader.Read(header.ProfileID[:])
	if err != nil {
		return err
	}
	if bytesRead < len(header.ProfileID) {
		return fmt.Errorf("unexpected EOF when reading profile ID")
	}

	// 28 reserved bytes
	for i := 0; i < 28/4; i++ {
		if _, err = binary.ReadU32Big(pr.reader); err != nil {
			return err
		}
	}

	return nil
}

func (pr *ProfileReader) readTagTable(tagTable *TagTable) error {
	const tagTableOffset = 128

	tagCount, err := binary.ReadU32Big(pr.reader)
	if err != nil {
		return err
	}

	type tagIndexEntry struct {
		offset uint32
		size   uint32
	}
	tagIndex := make(map[Signature]tagIndexEntry)

	endOfTagData := uint32(0)
	for i := uint32(0); i < tagCount; i++ {
		sig, err := binary.ReadU32Big(pr.reader)
		if err != nil {
			return err
		}

		offset, err := binary.ReadU32Big(pr.reader)
		if err != nil {
			return err
		}

		size, err := binary.ReadU32Big(pr.reader)
		if err != nil {
			return err
		}

		if offset+size > endOfTagData {
			endOfTagData = offset + size
		}

		tagIndex[Signature(sig)] = tagIndexEntry{
			offset: offset,
			size:   size,
		}
	}

	tagDataOffset := tagTableOffset + 4 + (tagCount * 12)
	tagData := make([]byte, endOfTagData-tagDataOffset)
	bytesRead, err := pr.reader.Read(tagData)
	if err != nil {
		return err
	}
	if bytesRead < len(tagData) {
		return fmt.Errorf("expected %d bytes of tag data but only got %d", len(tagData), bytesRead)
	}

	for sig, entry := range tagIndex {
		startOffset := entry.offset - tagDataOffset
		endOffset := startOffset + entry.size
		tagTable.add(sig, tagData[startOffset:endOffset])
	}

	return nil
}

func NewProfileReader(r binary.Reader) *ProfileReader {
	return &ProfileReader{
		reader: r,
	}
}
