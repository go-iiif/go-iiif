package binary

import "io"

type Reader interface {
	io.Reader
	io.ByteReader
}
