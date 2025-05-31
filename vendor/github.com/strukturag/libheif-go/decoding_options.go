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
	"errors"
	"runtime"
	"unsafe"
)

// DecodingOptions contain options that are used for decoding.
type DecodingOptions struct {
	options *C.struct_heif_decoding_options
}

func freeHeifDecodingOptions(options *DecodingOptions) {
	if options.options.decoder_id != nil {
		C.free(unsafe.Pointer(options.options.decoder_id))
	}
	C.heif_decoding_options_free(options.options)
	options.options = nil
}

// NewDecodingOptions creates new decoding options.
func NewDecodingOptions() (*DecodingOptions, error) {
	if err := checkLibraryVersion(); err != nil {
		return nil, err
	}

	options := &DecodingOptions{
		options: C.heif_decoding_options_alloc(),
	}
	if options.options == nil {
		return nil, errors.New("Could not allocate decoding options")
	}

	runtime.SetFinalizer(options, freeHeifDecodingOptions)
	options.options.version = 5
	options.options.color_conversion_options.version = 1
	return options, nil
}

// SetIgnoreTransformations sets whether geometric transformations like
// cropping, rotation, mirroring should be ignored.
func (o *DecodingOptions) SetIgnoreTransformations(ignore bool) {
	o.options.ignore_transformations = convertBool[C.uchar](ignore)
}

// GetIgnoreTransformations returns true if geometric transformations like
// cropping, rotation, mirroring should be ignored.
func (o *DecodingOptions) GetIgnoreTransformations() bool {
	return o.options.ignore_transformations != 0
}

// SetConvertHDRTo8Bit defines whether HDR images should be converted to 8bit
// during decoding.
func (o *DecodingOptions) SetConvertHDRTo8Bit(convert bool) {
	o.options.convert_hdr_to_8bit = convertBool[C.uchar](convert)
}

// GetConvertHDRTo8Bit returns true if HDR images will be converted to 8bit
// during decoding.
func (o *DecodingOptions) GetConvertHDRTo8Bit() bool {
	return o.options.convert_hdr_to_8bit != 0
}

// SetStrictDecoding enabled strict decoding and an error is returned for
// invalid input. Otherwise, it will try its best and add decoding warnings
// to the decoded heif_image. Default is non-strict.
func (o *DecodingOptions) SetStrictDecoding(strict bool) {
	o.options.strict_decoding = convertBool[C.uchar](strict)
}

// GetStrictDecoding returns true if strict decoding is enabled.
func (o *DecodingOptions) GetStrictDecoding() bool {
	return o.options.strict_decoding != 0
}

// SetDecoderId sets the id of the decoder to use. If an empty id is specified
// (the default), the highest priority decoder is chosen.
// The priority is defined in the plugin.
func (o *DecodingOptions) SetDecoderId(decoder string) {
	if o.options.decoder_id != nil {
		C.free(unsafe.Pointer(o.options.decoder_id))
	}

	if decoder == "" {
		o.options.decoder_id = nil
	} else {
		o.options.decoder_id = C.CString(decoder)
	}
}

// GetDecoderId returns the decoder id that should be used.
func (o *DecodingOptions) GetDecoderId() string {
	if o.options.decoder_id == nil {
		return ""
	}

	return C.GoString(o.options.decoder_id)
}

// SetChromaDownsamplingAlgorithm sets the chroma downsampling algorithm to use.
func (o *DecodingOptions) SetChromaDownsamplingAlgorithm(algorithm ChromaDownsamplingAlgorithm) {
	o.options.color_conversion_options.preferred_chroma_downsampling_algorithm = uint32(algorithm)
}

// GetChromaDownsamplingAlgorithm returns the chroma downsampling algorithm to use.
func (o *DecodingOptions) GetChromaDownsamplingAlgorithm() ChromaDownsamplingAlgorithm {
	return ChromaDownsamplingAlgorithm(o.options.color_conversion_options.preferred_chroma_downsampling_algorithm)
}

// SetChromaUpsamplingAlgorithm sets the chroma upsampling algorithm to use.
func (o *DecodingOptions) SetChromaUpsamplingAlgorithm(algorithm ChromaUpsamplingAlgorithm) {
	o.options.color_conversion_options.preferred_chroma_upsampling_algorithm = uint32(algorithm)
}

// GetChromaUpsamplingAlgorithm returns the chroma upsampling algorithm to use.
func (o *DecodingOptions) GetChromaUpsamplingAlgorithm() ChromaUpsamplingAlgorithm {
	return ChromaUpsamplingAlgorithm(o.options.color_conversion_options.preferred_chroma_upsampling_algorithm)
}

// SetOnlyUsePreferredChromaAlgorithm enforces to use the preferred algorithm.
// If set to false, libheif may also use a different algorithm if the preferred
// one is not available.
func (o *DecodingOptions) SetOnlyUsePreferredChromaAlgorithm(preferred bool) {
	o.options.color_conversion_options.only_use_preferred_chroma_algorithm = convertBool[C.uchar](preferred)
}

// GetOnlyUsePreferredChromaAlgorithm returns true if only the preferred chroma algorithm should be used
func (o *DecodingOptions) GetOnlyUsePreferredChromaAlgorithm() bool {
	return o.options.color_conversion_options.only_use_preferred_chroma_algorithm != 0
}
