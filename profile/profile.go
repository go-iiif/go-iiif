package profile

import (
	"fmt"
	"github.com/thisisaaronland/go-iiif/image"
)

type Profile struct {
	Context  string   `json:"@profile"`
	Id       string   `json:"@id"`
	Type     string   `json:"@type"` // Optional or iiif:Image
	Protocol string   `json:"protocol"`
	Width    int      `json:"width"`
	Height   int      `json:"height"`
	Profile  []string `json:"profile"`
	//	Sizes    []string `json:"sizes"` // Optional, existing/supported sizes.
	//	Tiles    []string `json:"tiles"` // Optional
}

func NewProfile(endpoint string, im image.Image) (*Profile, error) {

	dims, err := im.Dimensions()

	if err != nil {
		return nil, err
	}

	p := Profile{
		Context:  "http://iiif.io/api/image/2/context.json",
		Id:       fmt.Sprintf("%s/%s", endpoint, im.Identifier()),
		Type:     "iiif:Image",
		Protocol: "http://iiif.io/api/image",
		Width:    dims.Width(),
		Height:   dims.Height(),
		Profile: []string{
			fmt.Sprintf("%s/level2.json", endpoint),
		},
	}

	return &p, nil
}
