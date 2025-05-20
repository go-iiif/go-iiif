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

import (
	"image"
	"image/color"
	"io"
)

// --- High-level decoding API, always decodes primary image (if present).

func decodePrimaryImageFromReader(r io.Reader) (*ImageHandle, error) {
	ctx, err := NewContext()
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	if err := ctx.ReadFromMemory(data); err != nil {
		return nil, err
	}

	handle, err := ctx.GetPrimaryImageHandle()
	if err != nil {
		return nil, err
	}

	return handle, nil
}

func decodeImage(r io.Reader) (image.Image, error) {
	handle, err := decodePrimaryImageFromReader(r)
	if err != nil {
		return nil, err
	}

	img, err := handle.DecodeImage(ColorspaceUndefined, ChromaUndefined, nil)
	if err != nil {
		return nil, err
	}

	return img.GetImage()
}

func decodeConfig(r io.Reader) (image.Config, error) {
	var config image.Config
	handle, err := decodePrimaryImageFromReader(r)
	if err != nil {
		return config, err
	}

	config = image.Config{
		ColorModel: color.YCbCrModel,
		Width:      handle.GetWidth(),
		Height:     handle.GetHeight(),
	}
	return config, nil
}

func init() {
	image.RegisterFormat("heif", "????ftypheic", decodeImage, decodeConfig)
	image.RegisterFormat("heif", "????ftypheim", decodeImage, decodeConfig)
	image.RegisterFormat("heif", "????ftypheis", decodeImage, decodeConfig)
	image.RegisterFormat("heif", "????ftypheix", decodeImage, decodeConfig)
	image.RegisterFormat("heif", "????ftyphevc", decodeImage, decodeConfig)
	image.RegisterFormat("heif", "????ftyphevm", decodeImage, decodeConfig)
	image.RegisterFormat("heif", "????ftyphevs", decodeImage, decodeConfig)
	image.RegisterFormat("heif", "????ftypmif1", decodeImage, decodeConfig)
	image.RegisterFormat("avif", "????ftypavif", decodeImage, decodeConfig)
	image.RegisterFormat("avif", "????ftypavis", decodeImage, decodeConfig)
}
