package service

import (
	"context"
	_ "log"

	"github.com/aaronland/go-image-resize"
	"github.com/buckket/go-blurhash"
	iiifconfig "github.com/go-iiif/go-iiif/v5/config"
	iiifimage "github.com/go-iiif/go-iiif/v5/image"
)

func init() {

	ctx := context.Background()
	err := RegisterService(ctx, "blurhash", initBlurHashService)

	if err != nil {
		panic(err)
	}
}

func initBlurHashService(ctx context.Context, cfg *iiifconfig.Config, im iiifimage.Image) (Service, error) {
	return NewBlurHashService(cfg.BlurHash, im)
}

type BlurHashService struct {
	Service         `json:",omitempty"`
	BlurHashContext string `json:"@context"`
	BlurHashProfile string `json:"profile"`
	BlurHashLabel   string `json:"label"`
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

	hash, err := blurhash.Encode(cfg.X, cfg.Y, im)

	if err != nil {
		return nil, err
	}

	s := BlurHashService{
		BlurHashContext: "x-urn:service:go-iiif#blurhash",
		BlurHashProfile: "x-urn:service:go-iiif#blurhash",
		BlurHashLabel:   "x-urn:service:go-iiif#blurhash",
		BlurHash:        hash,
	}

	return &s, nil
}
