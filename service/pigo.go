package service

import (
	"context"
	pigo "github.com/esimov/pigo/core"
	iiifconfig "github.com/go-iiif/go-iiif/config"
	iiifimage "github.com/go-iiif/go-iiif/image"
	_ "log"
)

func init() {

	ctx := context.Background()
	err := RegisterService(ctx, "pigo", initPigoService)

	if err != nil {
		panic(err)
	}
}

func initPigoService(ctx context.Context, cfg *iiifconfig.Config, im iiifimage.Image) (Service, error) {
	return NewPigoService(cfg.Pigo, im)
}

type PigoService struct {
	Service     `json:",omitempty"`
	PigoContext string      `json:"@context"`
	PigoProfile string      `json:"profile"`
	PigoLabel   string      `json:"label"`
	PigoResults interface{} `json:"results"` // FIX ME
}

func (s *PigoService) Context() string {
	return s.PigoContext
}

func (s *PigoService) Profile() string {
	return s.PigoProfile
}

func (s *PigoService) Label() string {
	return s.PigoLabel
}

func (s *PigoService) Value() interface{} {
	return s.PigoResults
}

func NewPigoService(cfg iiifconfig.PigoConfig, image iiifimage.Image) (Service, error) {

	im, err := iiifimage.IIIFImageToGolangImage(image)

	if err != nil {
		return nil, err
	}

	pixels := pigo.RgbToGrayscale(im)

	bounds := im.Bounds()
	cols := bounds.Max.X
	rows := bounds.Max.Y

	cParams := pigo.CascadeParams{
		MinSize:     cfg.MinSize,
		MaxSize:     cfg.MaxSize,
		ShiftFactor: cfg.ShiftFactor,
		ScaleFactor: cfg.ScaleFactor,

		ImageParams: pigo.ImageParams{
			Pixels: pixels,
			Rows:   rows,
			Cols:   cols,
			Dim:    cols,
		},
	}

	pigo := pigo.NewPigo()

	var cascadeFile []byte // FIX ME...
	classifier, err := pigo.Unpack(cascadeFile)

	if err != nil {
		return nil, err
	}

	angle := 0.0 // cascade rotation angle. 0.0 is 0 radians and 1.0 is 2*pi radians

	dets := classifier.RunCascade(cParams, angle)

	s := PigoService{
		PigoContext: "x-urn:service:go-iiif#pigo",
		PigoProfile: "x-urn:service:go-iiif#pigo",
		PigoLabel:   "x-urn:service:go-iiif#pigo",
		PigoResults: dets,
	}

	return &s, nil
}
