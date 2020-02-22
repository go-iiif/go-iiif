package service

// https://groups.google.com/forum/#!topic/iiif-discuss/sPU5BvSWEOo
// http://palette.davidnewbury.com/

import (
	"context"
	iiifconfig "github.com/go-iiif/go-iiif/config"
	iiifimage "github.com/go-iiif/go-iiif/image"
	_ "log"
	"github.com/buckket/go-blurhash"
	"github.com/aaronland/go-image-resize"
)

type BlurHashService struct {
	Service        `json:",omitempty"`
	BlurHashContext string           `json:"@context"`
	BlurHashProfile string           `json:"profile"`
	BlurHashLabel   string           `json:"label"`
	BlurHash        string `json:"hash,omitempty"`
}

func (s *BlurHashService) Context() string {
	return s.BlurHashContext
}

func (s *BlurHashService) Profile() string {
	return s.BlurHashProfile
}

func (s *BlurHashService) Label() string {
	return s.BlurHashLabel
}

func (s *BlurHashService) Value() interface{} {
	return s.BlurHash
}

func NewBlurHashService(cfg iiifconfig.BlurHashConfig, image iiifimage.Image) (Service, error) {

	im, err := iiifimage.IIIFImageToGolangImage(image)

	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	im, err = resize.ResizeImageMax(ctx, im, cfg.Size)

	if err != nil {
		return nil, err
	}
	
	hash, err := blurhash.Encode(cfg.X, cfg.Y, &im)

	if err != nil {
		return nil, err
	}
	
	s := BlurHashService{
		BlurHashContext: "x-urn:service:go-iiif#palette",
		BlurHashProfile: "x-urn:service:go-iiif#palette",
		BlurHashLabel:   "x-urn:service:go-iiif#palette",
		BlurHash:        hash,
	}

	return &s, nil
}
