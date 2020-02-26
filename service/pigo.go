package service

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	pigo "github.com/esimov/pigo/core"
	iiifconfig "github.com/go-iiif/go-iiif/config"
	iiifimage "github.com/go-iiif/go-iiif/image"
	"gocloud.dev/blob"
	"image"
	"io/ioutil"
	"log"
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

/*

Why do I end up with different hashes for the same image...
(20200226/thisisaaronland)

go run cmd/pigo/main.go -in ~/sfomuseum/iiif/source/commissioners.jpg -out test.jpg -json -cf cascade/facefinder
Processing... ⠋
2020/02/26 11:58:59 PIGO 762 600 457200
2020/02/26 11:58:59 HASH 9f539418465730e6b7ab46e189df4bc9fbc39c97
Processing... ⠹
2020/02/26 11:58:59 {196 561 93 33.16809}
2020/02/26 11:58:59 {181 80 75 57.482105}
2020/02/26 11:58:59 {152 251 73 59.0851}
2020/02/26 11:58:59 {179 81 81 64.015015}


go run -mod vendor cmd/iiif-server/main.go -config-source file:///Users/asc/aaronland/go-iiif
2020/02/26 11:59:31 Listening on http://localhost:8080
2020/02/26 11:59:34 IIIF 762 600 457200
2020/02/26 11:59:34 HASH 105af3c4d74273310a7510dc2ab3bee9e64e98c1
2020/02/26 11:59:34 {180 81 76 65.160675}
2020/02/26 11:59:34 {196 561 93 34.383827}
2020/02/26 11:59:34 {152 251 73 60.488777}

*/

func NewPigoService(cfg iiifconfig.PigoConfig, iiif_im iiifimage.Image) (Service, error) {

	im, err := iiifimage.IIIFImageToGolangImage(iiif_im)

	if err != nil {
		return nil, err
	}

	im = pigo.ImgToNRGBA(im)

	pixels := pigo.RgbToGrayscale(im)

	bounds := im.Bounds()
	cols := bounds.Max.X
	rows := bounds.Max.Y

	log.Println("IIIF", cols, rows, len(pixels))

	hash := sha1.Sum(pixels)
	log.Printf("HASH %s", hex.EncodeToString(hash[:]))

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

	// basically all cribbed from here
	// https://github.com/esimov/pigo/blob/master/cmd/pigo/main.go

	ctx := context.Background()
	pigo := pigo.NewPigo()

	cascade_bucket, err := blob.OpenBucket(ctx, cfg.CascadeSource)

	if err != nil {
		return nil, err
	}

	cascade_fh, err := cascade_bucket.NewReader(ctx, cfg.CascadeFile, nil)

	if err != nil {
		return nil, err
	}

	defer cascade_fh.Close()

	cascade_body, err := ioutil.ReadAll(cascade_fh)

	if err != nil {
		return nil, err
	}

	classifier, err := pigo.Unpack(cascade_body)

	if err != nil {
		return nil, err
	}

	rects := make([]image.Rectangle, 0)

	angle := 0.0
	iou_threshold := 0.2
	q_threshold := float32(5.0)

	log.Println("RUN", angle, iou_threshold)

	faces := classifier.RunCascade(cParams, angle)
	faces = classifier.ClusterDetections(faces, iou_threshold)

	for _, face := range faces {
		log.Println(face)

		if face.Q > q_threshold {

			rects = append(rects, image.Rect(
				face.Col-face.Scale/2,
				face.Row-face.Scale/2,
				face.Scale,
				face.Scale,
			))
		}
	}

	s := PigoService{
		PigoContext: "x-urn:service:go-iiif#pigo",
		PigoProfile: "x-urn:service:go-iiif#pigo",
		PigoLabel:   "x-urn:service:go-iiif#pigo",
		PigoResults: rects,
	}

	return &s, nil
}
