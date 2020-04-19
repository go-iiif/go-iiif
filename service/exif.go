package service

import (
	"bytes"
	"context"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	iiifconfig "github.com/go-iiif/go-iiif/config"
	iiifimage "github.com/go-iiif/go-iiif/image"
	_ "log"
)

func init() {

	ctx := context.Background()
	err := RegisterService(ctx, "exif", initExifService)

	if err != nil {
		panic(err)
	}
}

func initExifService(ctx context.Context, cfg *iiifconfig.Config, im iiifimage.Image) (Service, error) {
	return NewExifService(cfg.Exif, im)
}

type ExifService struct {
	Service          `json:",omitempty"`
	ExifContext string `json:"@context"`
	ExifProfile string `json:"profile"`
	ExifLabel   string `json:"label"`
	ExifData    *exif.Exif `json:"data"`
}

func (s *ExifService) Context() string {
	return s.ExifContext
}

func (s *ExifService) Profile() string {
	return s.ExifProfile
}

func (s *ExifService) Label() string {
	return s.ExifLabel
}

func (s *ExifService) Value() interface{} {
	return s.ExifData
}

func NewExifService(cfg iiifconfig.ExifConfig, image iiifimage.Image) (Service, error) {

	br := bytes.NewReader(image.Body())

	exif.RegisterParsers(mknote.All...)

	data, err := exif.Decode(br)

	if err != nil {
		return nil, err
	}

	s := ExifService{
		ExifContext: "x-urn:service:go-iiif#exif",
		ExifProfile: "x-urn:service:go-iiif#exif",
		ExifLabel:   "x-urn:service:go-iiif#exif",
		ExifData:     data,
	}

	return &s, nil
}
