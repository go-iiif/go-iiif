package service

import (
	"context"
	"github.com/corona10/goimagehash"
	iiifconfig "github.com/go-iiif/go-iiif/config"
	iiifimage "github.com/go-iiif/go-iiif/image"
	_ "log"
)

func init() {

	ctx := context.Background()
	err := RegisterService(ctx, "imagehash", initImageHashService)

	if err != nil {
		panic(err)
	}
}

func initImageHashService(ctx context.Context, cfg *iiifconfig.Config, im iiifimage.Image) (Service, error) {
	return NewImageHashService(cfg.ImageHash, im)
}

type ImageHashService struct {
	Service          `json:",omitempty"`
	ImageHashContext string `json:"@context"`
	ImageHashProfile string `json:"profile"`
	ImageHashLabel   string `json:"label"`
	ImageHashAvg     string `json:"average,omitempty"`
	ImageHashDiff    string `json:"difference,omitempty"`
}

func (s *ImageHashService) Context() string {
	return s.ImageHashContext
}

func (s *ImageHashService) Profile() string {
	return s.ImageHashProfile
}

func (s *ImageHashService) Label() string {
	return s.ImageHashLabel
}

func (s *ImageHashService) Value() interface{} {
	return map[string]string{
		"average":    s.ImageHashAvg,
		"difference": s.ImageHashDiff,
	}
}

func NewImageHashService(cfg iiifconfig.ImageHashConfig, image iiifimage.Image) (Service, error) {

	im, err := iiifimage.IIIFImageToGolangImage(image)

	if err != nil {
		return nil, err
	}

	// please do these concurrently

	avg_hash, err := goimagehash.AverageHash(im)

	if err != nil {
		return nil, err
	}

	diff_hash, err := goimagehash.DifferenceHash(im)

	if err != nil {
		return nil, err
	}

	s := ImageHashService{
		ImageHashContext: "x-urn:service:go-iiif#imagehash",
		ImageHashProfile: "x-urn:service:go-iiif#imagehash",
		ImageHashLabel:   "x-urn:service:go-iiif#imagehash",
		ImageHashAvg:     avg_hash.ToString(),
		ImageHashDiff:    diff_hash.ToString(),
	}

	return &s, nil
}
