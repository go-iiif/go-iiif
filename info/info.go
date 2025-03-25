// package info provides methods for producing IIIF Image API "info.json" records.
package info

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	iiifimage "github.com/go-iiif/go-iiif/v7/image"
	iiiflevel "github.com/go-iiif/go-iiif/v7/level"
	iiifservice "github.com/go-iiif/go-iiif/v7/service"
)

// The URI for the IIIF Image API protocol.
const IMAGE_PROTOCOL string = "http://iiif.io/api/image"

// The URI for the IIIF Image API version 2.0 context.
const IMAGE_V2_CONTEXT string = "http://iiif.io/api/image/2/context.json"

// type Tile provides a struct for describing scaled (or resized) derivatives of an image.
type Size struct {
	// The height in pixels of each derivative.
	Height int `json:"height"`
	// The width in pixels of each derivative.
	Width int `json:"width"`
}

// type Tile provides a struct for describing tiled representations of an image.
type Tile struct {
	// The height in pixels of each tile.
	Height int `json:"height"`
	// The width in pixels of each tile.
	Width int `json:"width"`
	// The list of scale factors applied to the source image to produce tiles.
	ScaleFactors []int `json:"scaleFactors"`
}

// type Info provides a struct representing a IIIF Image API "info.json" record.
type Info struct {
	// The height in pixels of the image being described.
	Height int `json:"height"`
	// The width in pixels of the image being described.
	Width int `json:"width"`
	// The IIIF API context for the image being described.
	Context string `json:"@context"`
	// The ID (URI) of the image being described.
	Id string `json:"@id"`
	// The IIIF API protocol used to describe this image.
	Protocol string `json:"protocol"`
	// The IIIF profile (level) used to described this image.
	Profile []interface{} `json:"profile"`
	// Zero or more `Tile` entries describing zoomable tiles for this image.
	Tiles []*Tile `json:"tiles,omitempty"`
	// Zero or more `Size` entries describing deriative sizes for this image.
	Sizes []*Size `json:"sizes,omitempty"`
	// Zero or more `iiifservice.Service` entries describing additional or non-standard representations of this image.
	Services []iiifservice.Service `json:"service,omitempty"`
}

// New returns a new `Info` describing 'im' using IIIF level 'l' and the Image API represented by 'iiif_context'.
func New(iiif_context string, l iiiflevel.Level, im iiifimage.Image) (*Info, error) {

	dims, err := im.Dimensions()

	if err != nil {
		return nil, fmt.Errorf("Failed to derive dimensions for image, %w", err)
	}

	endpoint := strings.TrimRight(l.Endpoint(), "/")
	identifier := strings.TrimLeft(im.Identifier(), "/")
	id := fmt.Sprintf("%s/%s", endpoint, identifier)

	i := &Info{
		Context:  iiif_context,
		Protocol: IMAGE_PROTOCOL,
		Id:       id,
		Profile: []interface{}{
			l.Profile(),
			l,
		},
		Height: dims.Height(),
		Width:  dims.Width(),
	}

	return i, nil
}

func (i *Info) MarshalJSON(wr io.Writer) error {
	enc := json.NewEncoder(wr)
	return enc.Encode(i)
}
