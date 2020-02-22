package service

// https://groups.google.com/forum/#!topic/iiif-discuss/sPU5BvSWEOo
// http://palette.davidnewbury.com/

import (
	_ "context"
	"github.com/corona10/goimagehash"
	iiifconfig "github.com/go-iiif/go-iiif/config"
	iiifimage "github.com/go-iiif/go-iiif/image"
	_ "log"
)

type ImageHashService struct {
	Service          `json:",omitempty"`
	ImageHashContext string `json:"@context"`
	ImageHashProfile string `json:"profile"`
	ImageHashLabel   string `json:"label"`
	ImageHashAvg     string `json:"hash,omitempty"`
	ImageHashDiff    string `json:"hash,omitempty"`
	ImageHashExt     string `json:"hash,omitempty"`
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
	return s.ImageHashAvg
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

	ext_hash, err := goimagehash.ExtAverageHash(im, 8, 8)

	if err != nil {
		return nil, err
	}

	s := ImageHashService{
		ImageHashContext: "x-urn:service:go-iiif#palette",
		ImageHashProfile: "x-urn:service:go-iiif#palette",
		ImageHashLabel:   "x-urn:service:go-iiif#palette",
		ImageHashAvg:     avg_hash.ToString(),
		ImageHashDiff:    diff_hash.ToString(),
		ImageHashExt:     ext_hash.ToString(),
	}

	return &s, nil
}
