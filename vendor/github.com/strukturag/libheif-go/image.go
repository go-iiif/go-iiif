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
	"unsafe"
)

// Image contains information on a libheif image. It is either returned when
// decoding images or can be creates to encode an image using libheif.
type Image struct {
	image *C.struct_heif_image
}

// NewImage creates a new image to be used by libheif.
func NewImage(width, height int, colorspace Colorspace, chroma Chroma) (*Image, error) {
	if err := checkLibraryVersion(); err != nil {
		return nil, err
	}

	var image Image
	err := C.heif_image_create(C.int(width), C.int(height), uint32(colorspace), uint32(chroma), &image.image)
	if err := convertHeifError(err); err != nil {
		return nil, err
	}
	runtime.SetFinalizer(&image, freeHeifImage)
	return &image, nil
}

func freeHeifImage(image *Image) {
	C.heif_image_release(image.image)
	image.image = nil
}

// GetColorspace returns the colorspace of the image.
func (img *Image) GetColorspace() Colorspace {
	defer runtime.KeepAlive(img)

	return Colorspace(C.heif_image_get_colorspace(img.image))
}

// GetChromaFormat returns the chroma format of the image.
func (img *Image) GetChromaFormat() Chroma {
	defer runtime.KeepAlive(img)

	return Chroma(C.heif_image_get_chroma_format(img.image))
}

// GetWidth returns the width of the given channel of the image.
func (img *Image) GetWidth(channel Channel) int {
	defer runtime.KeepAlive(img)

	return int(C.heif_image_get_width(img.image, uint32(channel)))
}

// GetHeight returns the width of the given channel of the image.
func (img *Image) GetHeight(channel Channel) int {
	defer runtime.KeepAlive(img)

	return int(C.heif_image_get_height(img.image, uint32(channel)))
}

// GetBitsPerPixel returns the bits per pixel of the given channel of the image.
// Note that this is the number of bits used for storage of each pixel. Especially
// for HDR images this is probably not what you want, use "GetBitsPerPixelRange"
// instead.
func (img *Image) GetBitsPerPixel(channel Channel) int {
	defer runtime.KeepAlive(img)

	return int(C.heif_image_get_bits_per_pixel(img.image, uint32(channel)))
}

// GetBitsPerPixelRange returns the bits per pixel given channel of the image.
// This is the number of bits used for representing the pixel value, i.e. it will
// return "12" for 12bit HDR images (instead of "16" which would be the amount
// of bits used for storage).
func (img *Image) GetBitsPerPixelRange(channel Channel) int {
	defer runtime.KeepAlive(img)

	return int(C.heif_image_get_bits_per_pixel_range(img.image, uint32(channel)))
}

// GetImage convers the image to a Go Image object.
func (img *Image) GetImage() (image.Image, error) {
	var i image.Image
	cf := img.GetChromaFormat()
	switch cs := img.GetColorspace(); cs {
	case ColorspaceYCbCr:
		var subsample image.YCbCrSubsampleRatio
		switch cf {
		case Chroma420:
			subsample = image.YCbCrSubsampleRatio420
		case Chroma422:
			subsample = image.YCbCrSubsampleRatio422
		case Chroma444:
			subsample = image.YCbCrSubsampleRatio444
		default:
			return nil, fmt.Errorf("Unsupported YCbCr chroma format: %v", cf)
		}
		y, err := img.GetPlane(ChannelY)
		if err != nil {
			return nil, err
		}
		cb, err := img.GetPlane(ChannelCb)
		if err != nil {
			return nil, err
		}
		cr, err := img.GetPlane(ChannelCr)
		if err != nil {
			return nil, err
		}
		i = &image.YCbCr{
			Y:              y.Plane,
			Cb:             cb.Plane,
			Cr:             cr.Plane,
			YStride:        y.Stride,
			CStride:        cb.Stride,
			SubsampleRatio: subsample,
			Rect: image.Rectangle{
				Min: image.Point{
					X: 0,
					Y: 0,
				},
				Max: image.Point{
					X: img.GetWidth(ChannelY),
					Y: img.GetHeight(ChannelY),
				},
			},
		}
	case ColorspaceRGB:
		switch cf {
		case Chroma444:
			r, err := img.GetPlane(ChannelR)
			if err != nil {
				return nil, err
			}
			g, err := img.GetPlane(ChannelG)
			if err != nil {
				return nil, err
			}
			b, err := img.GetPlane(ChannelB)
			if err != nil {
				return nil, err
			}
			width := img.GetWidth(ChannelR)
			height := img.GetHeight(ChannelR)
			read_pos_r := 0
			read_pos_g := 0
			read_pos_b := 0
			write_pos := 0
			var rgba []byte
			var stride int
			if bpp := img.GetBitsPerPixelRange(ChannelR); bpp > 8 {
				// NOTE: We only support the same bits per pixel on all components.
				stride = width * 8
				rgba = make([]byte, height*stride)
				stride_add_r := r.Stride - width*2
				stride_add_g := g.Stride - width*2
				stride_add_b := b.Stride - width*2
				if bpp == 16 {
					for y := 0; y < height; y++ {
						for x := 0; x < width; x++ {
							rgba[write_pos] = r.Plane[read_pos_r]
							rgba[write_pos+1] = r.Plane[read_pos_r+1]
							rgba[write_pos+2] = g.Plane[read_pos_g]
							rgba[write_pos+3] = g.Plane[read_pos_g+1]
							rgba[write_pos+4] = b.Plane[read_pos_b]
							rgba[write_pos+5] = b.Plane[read_pos_b+1]
							rgba[write_pos+6] = 0xff
							rgba[write_pos+7] = 0xff
							read_pos_r += 2
							read_pos_g += 2
							read_pos_b += 2
							write_pos += 8
						}
						read_pos_r += stride_add_r
						read_pos_g += stride_add_g
						read_pos_b += stride_add_b
					}
				} else {
					for y := 0; y < height; y++ {
						for x := 0; x < width; x++ {
							r_value := (int16(r.Plane[read_pos_r+1]) << 8) | int16(r.Plane[read_pos_r])
							r_value = (r_value << (16 - uint(bpp))) | (r_value >> (2*uint(bpp) - 16))
							rgba[write_pos] = byte(r_value >> 8)
							rgba[write_pos+1] = byte(r_value & 0xff)
							g_value := (int16(g.Plane[read_pos_g+1]) << 8) | int16(g.Plane[read_pos_g])
							g_value = (g_value << (16 - uint(bpp))) | (g_value >> (2*uint(bpp) - 16))
							rgba[write_pos+2] = byte(g_value >> 8)
							rgba[write_pos+3] = byte(g_value & 0xff)
							b_value := (int16(b.Plane[read_pos_b+1]) << 8) | int16(b.Plane[read_pos_b])
							b_value = (b_value << (16 - uint(bpp))) | (b_value >> (2*uint(bpp) - 16))
							rgba[write_pos+4] = byte(b_value >> 8)
							rgba[write_pos+5] = byte(b_value & 0xff)
							rgba[write_pos+6] = 0xff
							rgba[write_pos+7] = 0xff
							read_pos_r += 2
							read_pos_g += 2
							read_pos_b += 2
							write_pos += 8
						}
						read_pos_r += stride_add_r
						read_pos_g += stride_add_g
						read_pos_b += stride_add_b
					}
				}

				i = &image.RGBA64{
					Pix:    rgba,
					Stride: stride,
					Rect: image.Rectangle{
						Min: image.Point{
							X: 0,
							Y: 0,
						},
						Max: image.Point{
							X: width,
							Y: height,
						},
					},
				}
			} else {
				stride = width * 4
				rgba = make([]byte, height*stride)
				stride_add_r := r.Stride - width
				stride_add_g := g.Stride - width
				stride_add_b := b.Stride - width
				for y := 0; y < height; y++ {
					for x := 0; x < width; x++ {
						rgba[write_pos] = r.Plane[read_pos_r]
						rgba[write_pos+1] = g.Plane[read_pos_g]
						rgba[write_pos+2] = b.Plane[read_pos_b]
						rgba[write_pos+3] = 0xff
						read_pos_r++
						read_pos_g++
						read_pos_b++
						write_pos += 4
					}
					read_pos_r += stride_add_r
					read_pos_g += stride_add_g
					read_pos_b += stride_add_b
				}

				i = &image.RGBA{
					Pix:    rgba,
					Stride: stride,
					Rect: image.Rectangle{
						Min: image.Point{
							X: 0,
							Y: 0,
						},
						Max: image.Point{
							X: width,
							Y: height,
						},
					},
				}
			}
		case ChromaInterleavedRGB:
			rgb, err := img.GetPlane(ChannelInterleaved)
			if err != nil {
				return nil, err
			}
			width := img.GetWidth(ChannelInterleaved)
			height := img.GetHeight(ChannelInterleaved)
			rgba := make([]byte, width*height*4)
			read_pos := 0
			write_pos := 0
			stride_add := rgb.Stride - width*3
			for y := 0; y < height; y++ {
				for x := 0; x < width; x++ {
					rgba[write_pos] = rgb.Plane[read_pos]
					rgba[write_pos+1] = rgb.Plane[read_pos+1]
					rgba[write_pos+2] = rgb.Plane[read_pos+2]
					rgba[write_pos+3] = 0xff
					read_pos += 3
					write_pos += 4
				}
				read_pos += stride_add
			}
			i = &image.RGBA{
				Pix:    rgba,
				Stride: width * 4,
				Rect: image.Rectangle{
					Min: image.Point{
						X: 0,
						Y: 0,
					},
					Max: image.Point{
						X: width,
						Y: height,
					},
				},
			}
		case ChromaInterleavedRGBA:
			rgba, err := img.GetPlane(ChannelInterleaved)
			if err != nil {
				return nil, err
			}
			i = &image.RGBA{
				Pix:    rgba.Plane,
				Stride: rgba.Stride,
				Rect: image.Rectangle{
					Min: image.Point{
						X: 0,
						Y: 0,
					},
					Max: image.Point{
						X: img.GetWidth(ChannelInterleaved),
						Y: img.GetHeight(ChannelInterleaved),
					},
				},
			}
		case ChromaInterleavedRRGGBB_BE:
			rgb, err := img.GetPlane(ChannelInterleaved)
			if err != nil {
				return nil, err
			}
			width := img.GetWidth(ChannelInterleaved)
			height := img.GetHeight(ChannelInterleaved)
			rgba := make([]byte, width*height*8)
			read_pos := 0
			write_pos := 0
			stride_add := rgb.Stride - width*6
			if bpp := img.GetBitsPerPixelRange(ChannelInterleaved); bpp != 16 {
				for y := 0; y < height; y++ {
					for x := 0; x < width; x++ {
						r_value := (int16(rgb.Plane[read_pos]) << 8) | int16(rgb.Plane[read_pos+1])
						r_value = (r_value << (16 - uint(bpp))) | (r_value >> (2*uint(bpp) - 16))
						rgba[write_pos] = byte(r_value >> 8)
						rgba[write_pos+1] = byte(r_value & 0xff)
						g_value := (int16(rgb.Plane[read_pos+2]) << 8) | int16(rgb.Plane[read_pos+3])
						g_value = (g_value << (16 - uint(bpp))) | (g_value >> (2*uint(bpp) - 16))
						rgba[write_pos+2] = byte(g_value >> 8)
						rgba[write_pos+3] = byte(g_value & 0xff)
						b_value := (int16(rgb.Plane[read_pos+4]) << 8) | int16(rgb.Plane[read_pos+5])
						b_value = (b_value << (16 - uint(bpp))) | (b_value >> (2*uint(bpp) - 16))
						rgba[write_pos+4] = byte(b_value >> 8)
						rgba[write_pos+5] = byte(b_value & 0xff)
						rgba[write_pos+6] = 0xff
						rgba[write_pos+7] = 0xff
						read_pos += 6
						write_pos += 8
					}
					read_pos += stride_add
				}
			} else {
				for y := 0; y < height; y++ {
					for x := 0; x < width; x++ {
						rgba[write_pos] = rgb.Plane[read_pos]
						rgba[write_pos+1] = rgb.Plane[read_pos+1]
						rgba[write_pos+2] = rgb.Plane[read_pos+2]
						rgba[write_pos+3] = rgb.Plane[read_pos+3]
						rgba[write_pos+4] = rgb.Plane[read_pos+4]
						rgba[write_pos+5] = rgb.Plane[read_pos+5]
						rgba[write_pos+6] = 0xff
						rgba[write_pos+7] = 0xff
						read_pos += 6
						write_pos += 8
					}
					read_pos += stride_add
				}
			}
			i = &image.RGBA64{
				Pix:    rgba,
				Stride: width * 4,
				Rect: image.Rectangle{
					Min: image.Point{
						X: 0,
						Y: 0,
					},
					Max: image.Point{
						X: width,
						Y: height,
					},
				},
			}
		case ChromaInterleavedRRGGBBAA_BE:
			rgba, err := img.GetPlane(ChannelInterleaved)
			if err != nil {
				return nil, err
			}
			width := img.GetWidth(ChannelInterleaved)
			height := img.GetHeight(ChannelInterleaved)
			var plane []byte
			if bpp := img.GetBitsPerPixelRange(ChannelInterleaved); bpp != 16 {
				read_pos := 0
				write_pos := 0
				stride_add := rgba.Stride - width*8
				plane = make([]byte, width*height*8)
				for y := 0; y < height; y++ {
					for x := 0; x < width; x++ {
						r_value := (int16(rgba.Plane[read_pos]) << 8) | int16(rgba.Plane[read_pos+1])
						r_value = (r_value << (16 - uint(bpp))) | (r_value >> (2*uint(bpp) - 16))
						plane[write_pos] = byte(r_value >> 8)
						plane[write_pos+1] = byte(r_value & 0xff)
						g_value := (int16(rgba.Plane[read_pos+2]) << 8) | int16(rgba.Plane[read_pos+3])
						g_value = (g_value << (16 - uint(bpp))) | (g_value >> (2*uint(bpp) - 16))
						plane[write_pos+2] = byte(g_value >> 8)
						plane[write_pos+3] = byte(g_value & 0xff)
						b_value := (int16(rgba.Plane[read_pos+4]) << 8) | int16(rgba.Plane[read_pos+5])
						b_value = (b_value << (16 - uint(bpp))) | (b_value >> (2*uint(bpp) - 16))
						plane[write_pos+4] = byte(b_value >> 8)
						plane[write_pos+5] = byte(b_value & 0xff)
						a_value := (int16(rgba.Plane[read_pos+6]) << 8) | int16(rgba.Plane[read_pos+7])
						a_value = (a_value << (16 - uint(bpp))) | (a_value >> (2*uint(bpp) - 16))
						plane[write_pos+6] = byte(a_value >> 8)
						plane[write_pos+7] = byte(a_value & 0xff)
						read_pos += 8
						write_pos += 8
					}
					read_pos += stride_add
				}
			} else {
				plane = rgba.Plane
			}
			i = &image.RGBA64{
				Pix:    plane,
				Stride: rgba.Stride,
				Rect: image.Rectangle{
					Min: image.Point{
						X: 0,
						Y: 0,
					},
					Max: image.Point{
						X: width,
						Y: height,
					},
				},
			}
		default:
			return nil, fmt.Errorf("Unsupported RGB chroma format: %v", cf)
		}
	default:
		return nil, fmt.Errorf("Unsupported colorspace: %v", cs)
	}

	return i, nil
}

// GetPlane returns an ImageAccess object that can be used to access the raw
// pixel values of the given channel.
func (img *Image) GetPlane(channel Channel) (*ImageAccess, error) {
	defer runtime.KeepAlive(img)

	height := C.heif_image_get_height(img.image, uint32(channel))
	if height == -1 {
		return nil, fmt.Errorf("No such channel %v", channel)
	}

	var stride C.int
	plane := C.heif_image_get_plane(img.image, uint32(channel), &stride)
	if plane == nil {
		return nil, fmt.Errorf("No such channel %v", channel)
	}

	ptr := unsafe.Pointer(plane)
	size := stride * height
	access := &ImageAccess{
		Plane:    C.GoBytes(ptr, size),
		planePtr: ptr,
		Stride:   int(stride),
		height:   int(height),
		image:    img,
	}
	return access, nil
}

// NewPlane creates a new plane for the image. Use this to set the pixel values
// of the image to encode.
func (img *Image) NewPlane(channel Channel, width, height, depth int) (*ImageAccess, error) {
	defer runtime.KeepAlive(img)

	err := C.heif_image_add_plane(img.image, uint32(channel), C.int(width), C.int(height), C.int(depth))
	if err := convertHeifError(err); err != nil {
		return nil, err
	}
	return img.GetPlane(channel)
}

// ScaleImage scales the image to the given width and height.
func (img *Image) ScaleImage(width int, height int) (*Image, error) {
	defer runtime.KeepAlive(img)

	var scaled_image Image
	err := C.heif_image_scale_image(img.image, &scaled_image.image, C.int(width), C.int(height), nil)
	if err := convertHeifError(err); err != nil {
		return nil, err
	}

	runtime.SetFinalizer(&scaled_image, freeHeifImage)
	return &scaled_image, nil
}
