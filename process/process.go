package process

import (
	iiifimage "github.com/thisisaaronland/go-iiif/image"
)

type Processor interface {
	ProcessURIWithInstructions(string, string, IIIFInstructions) (string, iiifimage.Image, error)
}
