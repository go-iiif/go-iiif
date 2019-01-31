package process

import (
	iiifuri "github.com/aaronland/go-iiif-uri"
	iiifimage "github.com/thisisaaronland/go-iiif/image"
)

type Label string

type Processor interface {
	ProcessURIWithInstructions(iiifuri.URI, Label, IIIFInstructions) (iiifuri.URI, iiifimage.Image, error)
}
