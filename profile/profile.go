package profile

import (
	_ "fmt"
	iiifservice "github.com/go-iiif/go-iiif/v5/service"
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

func (p *Profile) AddService(s iiifservice.Service) {
	p.Services = append(p.Services, s)
}
