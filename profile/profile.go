package profile

import (
	"fmt"
	iiifimage "github.com/go-iiif/go-iiif/v5/image"
	iiifservice "github.com/go-iiif/go-iiif/v5/service"
	"path/filepath"
)

type Profile struct {
	Context  string                `json:"@context"`
	Id       string                `json:"@id"`
	Type     string                `json:"@type,omitempty"`
	Protocol string                `json:"protocol"`
	Width    int                   `json:"width"`
	Height   int                   `json:"height"`
	Profile  []interface{}         `json:"profile"`
	Sizes    []string              `json:"sizes,omitempty"`
	Tiles    []string              `json:"tiles,omitempty"`
	Services []iiifservice.Service `json:"service,omitempty"`
}

func (p *Profile) AddImage(endpoint string, im iiifimage.Image) error {

	dims, err := im.Dimensions()

	if err != nil {
		return fmt.Errorf("Failed to derive dimensions for image, %w", err)
	}

	p.Id = filepath.Join(endpoint, im.Identifier())
	p.Height = dims.Height()
	p.Width = dims.Width()

	return nil
}

func (p *Profile) AddService(s iiifservice.Service) {
	p.Services = append(p.Services, s)
}
