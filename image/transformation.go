package image

import (
	"errors"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var regionError = "IIIF 2.1 `region` argument is not recognized: %#v"
var sizeError = "IIIF 2.1 `size` argument is not recognized: %#v"
var qualityError = "IIIF 2.1 `quality` and `format` arguments were expected: %#v"
var rotationError = "IIIF 2.1 `rotation` argument is not recognized: %#v"
var rotationMissing = "libvips cannot rotate angle that isn't a multiple of 90: %#v"

var re_alpha *regexp.Regexp
var re_region *regexp.Regexp
var re_size *regexp.Regexp
var re_rotation *regexp.Regexp
var re_quality *regexp.Regexp

func init() {

	var err error

	re_alpha, err = regexp.Compile(`^[a-z]+$`)

	if err != nil {
		log.Fatal(err)
	}

	re_region, err = regexp.Compile(`^(?:pct\:)?\d+\,\d+\,\d+\,\d+$`)

	if err != nil {
		log.Fatal(err)
	}

	re_rotation, err = regexp.Compile(`^\!?\d+`)

	if err != nil {
		log.Fatal(err)
	}

	re_quality, err = regexp.Compile(`^(?:color|grey|bitonal|default|dither)$`)

	if err != nil {
		log.Fatal(err)
	}

	re_size, err = regexp.Compile(`^(?:(?:max|full)|(?:\d+\,\d+)|(?:\!\d+\,\d+)|(\d+\,)|(\,\d+)|(pct\:\d+))$`)

	if err != nil {
		log.Fatal(err)
	}

}

type RegionInstruction struct {
	X      int
	Y      int
	Height int
	Width  int
}

type SizeInstruction struct {
	Height  int
	Width   int
	Force   bool
	Enlarge bool
}

type RotationInstruction struct {
	Flip  bool
	Angle int64
}

// full
// square
// x,y,w,h (in pixels)
// pct:x,y,w,h (in percents)

func IsValidRegion(region string) (bool, error) {

	if !re_region.MatchString(region) {
		return false, errors.New("Invalid region parameter")
	}

	return true, nil
}

// max, full
// w,h (deform)
// !w,h (best fit within size)
// w, (force width)
// ,h (force height)
// pct:n (resize)

func IsValidSize(size string) (bool, error) {

	if !re_size.MatchString(size) {
		return false, errors.New("Invalid size parameter")
	}
	return true, nil
}

// n angle clockwise in degrees
// !n angle clockwise in degrees with a flip (beforehand)

func IsValidRotation(rotation string) (bool, error) {

	if !re_rotation.MatchString(rotation) {
		return false, errors.New("Invalid size parameter")
	}

	parsed, err := strconv.ParseInt(strings.Trim(rotation, "!"), 10, 64)

	if err != nil {
		return false, err
	}

	if parsed > 360 {
		return false, errors.New("Invalid rotation parameter")
	}

	return true, nil
}

func IsValidQuality(quality string) (bool, error) {

	if !re_quality.MatchString(quality) {
		return false, errors.New("Invalid quality parameter")
	}

	return true, nil
}

func IsValidFormat(format string) (bool, error) {
	return true, nil
}

type Transformation struct {
	Region   string
	Size     string
	Rotation string
	Quality  string
	Format   string
}

func NewTransformation(region string, size string, rotation string, quality string, format string) (*Transformation, error) {

	var ok bool
	var err error

	ok, err = IsValidRegion(region)

	if !ok {
		return nil, err
	}

	ok, err = IsValidSize(size)

	if !ok {
		return nil, err
	}

	ok, err = IsValidRotation(rotation)

	if !ok {
		return nil, err
	}

	ok, err = IsValidQuality(quality)

	if !ok {
		return nil, err
	}

	ok, err = IsValidFormat(format)

	if !ok {
		return nil, err
	}

	t := Transformation{
		Region:   region,
		Size:     size,
		Rotation: rotation,
		Quality:  quality,
		Format:   format,
	}

	return &t, nil
}

func (t *Transformation) ToURI(id string) string {

	nodes := []string{
		id,
		t.Region,
		t.Size,
		t.Rotation,
		t.Quality,
	}

	return fmt.Sprintf("%s.%s", strings.Join(nodes, "/"), t.Format)
}

func (t *Transformation) HasTransformation() bool {

	if t.Region != "full" {
		return true
	}

	if t.Size != "full" {
		return true
	}

	if t.Rotation != "0" {
		return true
	}

	if t.Quality != "default" {
		return true
	}

	return false
}

func (t *Transformation) RegionInstructions(im Image) (*RegionInstruction, error) {

	dims, err := im.Dimensions()

	if err != nil {
		return nil, err
	}

	width := dims.Width()
	height := dims.Height()

	if t.Region == "square" {

		var x int
		var y int

		if width < height {
			y = (height - width) / 2.
			x = width
		} else {
			x = (width - height) / 2.
			y = height
		}

		y = x

		instruction := RegionInstruction{
			X:      x,
			Y:      y,
			Width:  width,
			Height: height,
		}

		return &instruction, nil
	}

	arr := strings.Split(t.Region, ":")

	if len(arr) == 1 {

		sizes := strings.Split(arr[0], ",")

		if len(sizes) != 4 {
			message := fmt.Sprintf("Invalid region")
			return nil, errors.New(message)
		}

		x, err := strconv.ParseInt(sizes[0], 10, 64)

		if err != nil {
			return nil, err
		}

		y, err := strconv.ParseInt(sizes[1], 10, 64)

		if err != nil {
			return nil, err
		}

		w, err := strconv.ParseInt(sizes[2], 10, 64)

		if err != nil {
			return nil, err
		}

		h, err := strconv.ParseInt(sizes[3], 10, 64)

		if err != nil {
			return nil, err
		}

		instruction := RegionInstruction{
			Width:  int(w),
			Height: int(h),
			X:      int(x),
			Y:      int(y),
		}

		return &instruction, nil

	}

	if arr[0] == "pct" {

		sizes := strings.Split(arr[1], ",")

		if len(sizes) != 4 {
			message := fmt.Sprintf("Invalid region", t.Region)
			return nil, errors.New(message)
		}

		px, err := strconv.ParseFloat(sizes[0], 64)

		if err != nil {
			return nil, err
		}

		py, err := strconv.ParseFloat(sizes[1], 64)

		if err != nil {
			return nil, err
		}

		pw, err := strconv.ParseFloat(sizes[2], 64)

		if err != nil {
			return nil, err
		}

		ph, err := strconv.ParseFloat(sizes[3], 64)

		if err != nil {
			return nil, err
		}

		w := int(math.Ceil(float64(width) * pw / 100.))
		h := int(math.Ceil(float64(height) * ph / 100.))
		x := int(math.Ceil(float64(width) * px / 100.))
		y := int(math.Ceil(float64(height) * py / 100.))

		instruction := RegionInstruction{
			Width:  w,
			Height: h,
			X:      x,
			Y:      y,
		}

		return &instruction, nil

	} else {
	}

	message := fmt.Sprintf("Unrecognized region")
	return nil, errors.New(message)

}

func (t *Transformation) SizeInstructions(im Image) (*SizeInstruction, error) {

	w := 0
	h := 0
	force := false
	enlarge := false

	arr := strings.Split(t.Size, ":")

	if len(arr) == 1 {

		best := strings.HasPrefix(t.Size, "!")
		sizes := strings.Split(strings.Trim(arr[0], "!"), ",")

		if len(sizes) != 2 {
			message := fmt.Sprintf(sizeError, t.Size)
			return nil, errors.New(message)
		}

		wi, err_w := strconv.ParseInt(sizes[0], 10, 64)
		hi, err_h := strconv.ParseInt(sizes[1], 10, 64)

		if err_w != nil && err_h != nil {
			message := fmt.Sprintf(sizeError, t.Size)
			return nil, errors.New(message)

		} else if err_w == nil && err_h == nil {

			w = int(wi)
			h = int(hi)

			if best {
				enlarge = true
			} else {
				force = true
			}

		} else if err_h != nil {
			w = int(wi)
			h = 0
		} else {
			w = 0
			h = int(hi)
		}

		instruction := SizeInstruction{
			Height:  h,
			Width:   w,
			Enlarge: enlarge,
			Force:   force,
		}

		return &instruction, nil

	} else if arr[0] == "pct" {

		pct, err := strconv.ParseFloat(arr[1], 64)

		if err != nil {
			err := errors.New("invalid size")
			return nil, err
		}

		dims, err := im.Dimensions()

		if err != nil {
			return nil, err
		}

		width := dims.Width()
		height := dims.Height()

		w = int(math.Ceil(pct / 100 * float64(width)))
		h = int(math.Ceil(pct / 100 * float64(height)))

	} else {

		message := fmt.Sprintf(sizeError, t.Size)
		return nil, errors.New(message)
	}

	instruction := SizeInstruction{
		Height:  h,
		Width:   w,
		Enlarge: enlarge,
		Force:   force,
	}

	return &instruction, nil

}

func (t *Transformation) RotationInstructions(im Image) (*RotationInstruction, error) {

	flip := strings.HasPrefix(t.Rotation, "!")
	angle, err := strconv.ParseInt(strings.Trim(t.Rotation, "!"), 10, 64)

	if err != nil {
		message := fmt.Sprintf(rotationError, t.Rotation)
		return nil, errors.New(message)

	}

	instruction := RotationInstruction{
		Flip:  flip,
		Angle: angle,
	}

	return &instruction, nil
}
