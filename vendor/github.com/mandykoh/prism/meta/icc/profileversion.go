package icc

import "fmt"

type Version struct {
	Major       byte
	MinorAndRev byte
}

func (pv Version) String() string {
	return fmt.Sprintf("%d.%d.%d", pv.Major, pv.MinorAndRev>>4, pv.MinorAndRev&3)
}
