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
	"unsafe"
)

// formatVersion formats a numeric version of libheif encoded as BCD
// 0xHHMMLL00 = HH.MM.LL.
// For example: 0x02143000 is version 2.14.30
func formatVersion(value int) string {
	return fmt.Sprintf("%d.%d.%d",
		value>>24,
		(value>>16)&0xff,
		(value>>8)&0xff,
	)
}

// init initializes the libheif library.
func init() {
	C.heif_init(nil)
}

// checkLibraryVersion checks if the loaded libheif library has at least the
// version that was used while compiling.
func checkLibraryVersion() error {
	runtime_version := C.heif_get_version_number()
	if runtime_version >= build_version {
		return nil
	}

	return fmt.Errorf("expected at least libheif version %s, got %s",
		formatVersion(build_version),
		formatVersion(int(runtime_version)),
	)
}

// GetVersion returns the libheif version string.
func GetVersion() string {
	return C.GoString(C.heif_get_version())
}

// HaveDecoderForFormat checks if a decoder is available for the given format.
// Note that the decoder still may not be able to decode all variants of that format.
// You will have to query that further or just try to decode and check the returned error.
func HaveDecoderForFormat(format CompressionFormat) bool {
	return C.heif_have_decoder_for_format(C.enum_heif_compression_format(format)) != 0
}

// HaveEncoderForFormat checks if an encoder is available for the given format.
// Note that the encoder may be limited to a certain subset of features (e.g. only 8 bit, only lossy).
// You will have to query the specific capabilities further.
func HaveEncoderForFormat(format CompressionFormat) bool {
	return C.heif_have_encoder_for_format(C.enum_heif_compression_format(format)) != 0
}

func convertItemIDs(ids []C.heif_item_id, count int) []int {
	result := make([]int, count)
	for i := 0; i < count; i++ {
		result[i] = int(ids[i])
	}
	return result
}

func convertBool[T C.uchar | C.int](value bool) T { // nolint
	if value {
		return 1
	}

	return 0
}

func nextPointer[T any](value *T) *T {
	return (*T)(unsafe.Pointer(uintptr(unsafe.Pointer(value)) + unsafe.Sizeof(value)))
}

func makePointer[T any](value T) *T {
	return &value
}
