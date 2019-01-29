package process

import (
	iiifimage "github.com/thisisaaronland/go-iiif/image"
)

type Processor interface {
	ProcessURIWithInstructionSet(string, IIIFInstructionSet) (map[string]interface{}, error)
	ProcessURIWithInstructions(string, IIIFInstructions) (string, iiifimage.Image, error)
}
