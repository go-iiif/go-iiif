package process

import (
	iiifuri "github.com/go-iiif/go-iiif-uri"
	iiifimage "github.com/go-iiif/go-iiif/v7/image"
)

type Label string

type Processor interface {
	ProcessURIWithInstructions(iiifuri.URI, Label, IIIFInstructions) (iiifuri.URI, iiifimage.Image, error)
}
