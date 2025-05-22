package exif

import (
	"strconv"

	"github.com/dsoprea/go-exif/v3"
)

func NewIfdBuilderWithOrientation(ifd *exif.Ifd, orientation string) (*exif.IfdBuilder, error) {

	var ib *exif.IfdBuilder

	if ifd != nil {

		ib = exif.NewIfdBuilderFromExistingChain(ifd)

		ifdPath := "IFD0"

		ifd_ib, err := exif.GetOrCreateIbFromRootIb(ib, ifdPath)

		if err != nil {
			return nil, err
		}

		oint, _ := strconv.Atoi(orientation) // top left
		oint16 := uint16(oint)

		err = ifd_ib.SetStandardWithName("Orientation", []uint16{oint16})

		if err != nil {
			return nil, err
		}
	}

	return ib, nil
}
