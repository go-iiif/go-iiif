/*
 * Go interface to libheif
 *
 * Copyright (c) 2018-2024 struktur AG, Joachim Bauch <bauch@struktur.de>
 *
 * libheif is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as
 * published by the Free Software Foundation, either version 3 of
 * the License, or (at your option) any later version.
 *
 * libheif is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with libheif.  If not, see <http://www.gnu.org/licenses/>.
 */

package libheif

// #cgo pkg-config: libheif
// #include <stdlib.h>
// #include <string.h>
// #include <libheif/heif.h>
import "C"

import (
	"fmt"
	"image"
	"runtime"
)

func imageFromRGBA(i *image.RGBA) (*Image, error) {
	min := i.Bounds().Min
	max := i.Bounds().Max
	w := max.X - min.X
	h := max.Y - min.Y

	out, err := NewImage(w, h, ColorspaceRGB, ChromaInterleavedRGBA)
	if err != nil {
		return nil, fmt.Errorf("failed to create image: %v", err)
	}

	p, err := out.NewPlane(ChannelInterleaved, w, h, 8)
	if err != nil {
		return nil, fmt.Errorf("failed to add plane: %v", err)
	}
	p.setData([]byte(i.Pix), w*4)

	return out, nil
}

func imageFromNRGBA(i *image.NRGBA) (*Image, error) {
	min := i.Bounds().Min
	max := i.Bounds().Max
	w := max.X - min.X
	h := max.Y - min.Y

	out, err := NewImage(w, h, ColorspaceRGB, ChromaInterleavedRGBA)
	if err != nil {
		return nil, fmt.Errorf("failed to create image: %v", err)
	}

	p, err := out.NewPlane(ChannelInterleaved, w, h, 8)
	if err != nil {
		return nil, fmt.Errorf("failed to add plane: %v", err)
	}
	p.setData([]byte(i.Pix), w*4)

	return out, nil
}

func imageFromRGBA64(i *image.RGBA64, compression CompressionFormat) (*Image, error) {
	min := i.Bounds().Min
	max := i.Bounds().Max
	w := max.X - min.X
	h := max.Y - min.Y

	out, err := NewImage(w, h, ColorspaceRGB, ChromaInterleavedRRGGBBAA_BE)
	if err != nil {
		return nil, fmt.Errorf("failed to create image: %v", err)
	}

	var depth int
	switch compression {
	case CompressionAV1:
		depth = 12
	case CompressionHEVC:
		depth = 10
	default:
		depth = 16
	}
	p, err := out.NewPlane(ChannelInterleaved, w, h, depth)
	if err != nil {
		return nil, fmt.Errorf("failed to add plane: %v", err)
	}

	if depth == 16 {
		p.setData(i.Pix, w*8)
	} else {
		shift := 16 - depth
		pix := make([]byte, w*h*8)
		read_pos := 0
		write_pos := 0
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				r := (uint16(i.Pix[read_pos]) << 8) | uint16(i.Pix[read_pos+1])
				r = r >> shift
				pix[write_pos] = byte(r >> 8)
				pix[write_pos+1] = byte(r & 0xff)
				read_pos += 2
				g := (uint16(i.Pix[read_pos]) << 8) | uint16(i.Pix[read_pos+1])
				g = g >> shift
				pix[write_pos+2] = byte(g >> 8)
				pix[write_pos+3] = byte(g & 0xff)
				read_pos += 2
				b := (uint16(i.Pix[read_pos]) << 8) | uint16(i.Pix[read_pos+1])
				b = b >> shift
				pix[write_pos+4] = byte(b >> 8)
				pix[write_pos+5] = byte(b & 0xff)
				read_pos += 2
				a := (uint16(i.Pix[read_pos]) << 8) | uint16(i.Pix[read_pos+1])
				a = a >> shift
				pix[write_pos+6] = byte(a >> 8)
				pix[write_pos+7] = byte(a & 0xff)
				read_pos += 2
				write_pos += 8
			}
		}
		p.setData(pix, w*8)
	}

	return out, nil
}

func imageFromNRGBA64(i *image.NRGBA64, compression CompressionFormat) (*Image, error) {
	min := i.Bounds().Min
	max := i.Bounds().Max
	w := max.X - min.X
	h := max.Y - min.Y

	out, err := NewImage(w, h, ColorspaceRGB, ChromaInterleavedRRGGBBAA_BE)
	if err != nil {
		return nil, fmt.Errorf("failed to create image: %v", err)
	}

	var depth int
	switch compression {
	case CompressionAV1:
		depth = 12
	case CompressionHEVC:
		depth = 10
	default:
		depth = 16
	}
	p, err := out.NewPlane(ChannelInterleaved, w, h, depth)
	if err != nil {
		return nil, fmt.Errorf("failed to add plane: %v", err)
	}

	if depth == 16 {
		p.setData(i.Pix, w*8)
	} else {
		shift := 16 - depth
		pix := make([]byte, w*h*8)
		read_pos := 0
		write_pos := 0
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				r := (uint16(i.Pix[read_pos]) << 8) | uint16(i.Pix[read_pos+1])
				r = r >> shift
				pix[write_pos] = byte(r >> 8)
				pix[write_pos+1] = byte(r & 0xff)
				read_pos += 2
				g := (uint16(i.Pix[read_pos]) << 8) | uint16(i.Pix[read_pos+1])
				g = g >> shift
				pix[write_pos+2] = byte(g >> 8)
				pix[write_pos+3] = byte(g & 0xff)
				read_pos += 2
				b := (uint16(i.Pix[read_pos]) << 8) | uint16(i.Pix[read_pos+1])
				b = b >> shift
				pix[write_pos+4] = byte(b >> 8)
				pix[write_pos+5] = byte(b & 0xff)
				read_pos += 2
				a := (uint16(i.Pix[read_pos]) << 8) | uint16(i.Pix[read_pos+1])
				a = a >> shift
				pix[write_pos+6] = byte(a >> 8)
				pix[write_pos+7] = byte(a & 0xff)
				read_pos += 2
				write_pos += 8
			}
		}
		p.setData(pix, w*8)
	}

	return out, nil
}

func imageFromGray(i *image.Gray) (*Image, error) {
	min := i.Bounds().Min
	max := i.Bounds().Max
	w := max.X - min.X
	h := max.Y - min.Y

	out, err := NewImage(w, h, ColorspaceYCbCr, ChromaMonochrome)
	if err != nil {
		return nil, fmt.Errorf("failed to create image: %v", err)
	}

	const depth = 8
	pY, err := out.NewPlane(ChannelY, w, h, depth)
	if err != nil {
		return nil, fmt.Errorf("failed to add Y plane: %v", err)
	}
	pY.setData([]byte(i.Pix), i.Stride)

	return out, nil
}

func imageFromYCbCr(i *image.YCbCr) (*Image, error) {
	min := i.Bounds().Min
	max := i.Bounds().Max
	w := max.X - min.X
	h := max.Y - min.Y

	var cm Chroma
	switch sr := i.SubsampleRatio; sr {
	case image.YCbCrSubsampleRatio420:
		cm = Chroma420
	case image.YCbCrSubsampleRatio422:
		cm = Chroma422
	case image.YCbCrSubsampleRatio444:
		cm = Chroma444
	default:
		return nil, fmt.Errorf("unsupported subsample ratio: %s", sr.String())
	}

	out, err := NewImage(w, h, ColorspaceYCbCr, cm)
	if err != nil {
		return nil, fmt.Errorf("failed to create image: %v", err)
	}

	const depth = 8
	pY, err := out.NewPlane(ChannelY, w, h, depth)
	if err != nil {
		return nil, fmt.Errorf("failed to add Y plane: %v", err)
	}
	pY.setData([]byte(i.Y), i.YStride)

	switch cm {
	case Chroma420:
		halfW, halfH := (w+1)/2, (h+1)/2
		pCb, err := out.NewPlane(ChannelCb, halfW, halfH, depth)
		if err != nil {
			return nil, fmt.Errorf("failed to add Cb plane: %v", err)
		}
		pCb.setData([]byte(i.Cb), i.CStride)
		pCr, err := out.NewPlane(ChannelCr, halfW, halfH, depth)
		if err != nil {
			return nil, fmt.Errorf("failed to add Cr plane: %v", err)
		}
		pCr.setData([]byte(i.Cr), i.CStride)
	case Chroma422:
		halfW := (w + 1) / 2
		pCb, err := out.NewPlane(ChannelCb, halfW, h, depth)
		if err != nil {
			return nil, fmt.Errorf("failed to add Cb plane: %v", err)
		}
		pCb.setData([]byte(i.Cb), i.CStride)
		pCr, err := out.NewPlane(ChannelCr, halfW, h, depth)
		if err != nil {
			return nil, fmt.Errorf("failed to add Cr plane: %v", err)
		}
		pCr.setData([]byte(i.Cr), i.CStride)
	case Chroma444:
		pCb, err := out.NewPlane(ChannelCb, w, h, depth)
		if err != nil {
			return nil, fmt.Errorf("failed to add Cb plane: %v", err)
		}
		pCb.setData([]byte(i.Cb), i.CStride)
		pCr, err := out.NewPlane(ChannelCr, w, h, depth)
		if err != nil {
			return nil, fmt.Errorf("failed to add Cr plane: %v", err)
		}
		pCr.setData([]byte(i.Cr), i.CStride)
	}

	return out, nil
}

// EncoderParameterSetter is a function that can configure an encoder.
type EncoderParameterSetter func(encoder *Encoder) error

// SetEncoderQuality returns a function that sets the quality of an encoder.
func SetEncoderQuality(quality int) EncoderParameterSetter {
	return func(encoder *Encoder) error {
		return encoder.SetQuality(quality)
	}
}

// SetEncoderLossless returns a function that sets the lossless mode of an encoder.
func SetEncoderLossless(lossless LosslessMode) EncoderParameterSetter {
	return func(encoder *Encoder) error {
		return encoder.SetLossless(lossless)
	}
}

// SetEncoderLoggingLevel returns a function that sets the logging level of an encoder.
func SetEncoderLoggingLevel(level LoggingLevel) EncoderParameterSetter {
	return func(encoder *Encoder) error {
		return encoder.SetLoggingLevel(level)
	}
}

// SetEncoderParameterBool returns a function that sets a boolean parameter in an encoder.
func SetEncoderParameterBool(name string, value bool) EncoderParameterSetter {
	return func(encoder *Encoder) error {
		return encoder.SetParameterBool(name, value)
	}
}

// SetEncoderParameterInteger returns a function that sets an integer parameter in an encoder.
func SetEncoderParameterInteger(name string, value int) EncoderParameterSetter {
	return func(encoder *Encoder) error {
		return encoder.SetParameterInteger(name, value)
	}
}

// SetEncoderParameterString returns a function that sets a string parameter in an encoder.
func SetEncoderParameterString(name string, value string) EncoderParameterSetter {
	return func(encoder *Encoder) error {
		return encoder.SetParameterString(name, value)
	}
}

// SetEncoderParameter returns a function that sets an arbitrary parameter in an encoder.
func SetEncoderParameter(name string, value string) EncoderParameterSetter {
	return func(encoder *Encoder) error {
		return encoder.SetParameter(name, value)
	}
}

// EncodeFromImage is a high-level function to encode a Go Image to a new Context.
func EncodeFromImage(img image.Image, compression CompressionFormat, params ...EncoderParameterSetter) (*Context, *ImageHandle, error) {
	if err := checkLibraryVersion(); err != nil {
		return nil, nil, err
	}

	var out *Image

	switch i := img.(type) {
	default:
		return nil, nil, fmt.Errorf("unsupported image type: %T", i)
	case *image.RGBA:
		tmp, err := imageFromRGBA(i)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create image: %v", err)
		}
		out = tmp
	case *image.NRGBA:
		tmp, err := imageFromNRGBA(i)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create image: %v", err)
		}
		out = tmp
	case *image.RGBA64:
		tmp, err := imageFromRGBA64(i, compression)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create image: %v", err)
		}
		out = tmp
	case *image.NRGBA64:
		tmp, err := imageFromNRGBA64(i, compression)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create image: %v", err)
		}
		out = tmp
	case *image.Gray:
		tmp, err := imageFromGray(i)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create image: %v", err)
		}
		out = tmp
	case *image.YCbCr:
		tmp, err := imageFromYCbCr(i)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create image: %v", err)
		}
		out = tmp
	}

	ctx, err := NewContext()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create HEIF context: %v", err)
	}

	enc, err := ctx.NewEncoder(compression)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create encoder: %v", err)
	}

	for _, param := range params {
		if err := param(enc); err != nil {
			return nil, nil, fmt.Errorf("error setting parameter: %w", err)
		}
	}

	encOpts, err := NewEncodingOptions()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get encoding options: %v", err)
	}

	defer runtime.KeepAlive(ctx)
	defer runtime.KeepAlive(out)
	defer runtime.KeepAlive(enc)
	defer runtime.KeepAlive(encOpts)

	var handle ImageHandle
	err2 := C.heif_context_encode_image(ctx.context, out.image, enc.encoder, encOpts.options, &handle.handle)
	if err := convertHeifError(err2); err != nil {
		return nil, nil, fmt.Errorf("failed to encode image: %v", err)
	}

	runtime.SetFinalizer(&handle, freeHeifImageHandle)
	return ctx, &handle, nil
}
