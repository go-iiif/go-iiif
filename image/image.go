package image

import ()

type Image interface {
	Identifier() string
	Rename(string) error
	Transform(*Transformation) error // http://iiif.io/api/image/2.1/#order-of-implementation
	Update([]byte) error
	Body() []byte
	Format() string
	ContentType() string
	Dimensions() (Dimensions, error)
}

type Dimensions interface {
	Height() int
	Width() int
}
