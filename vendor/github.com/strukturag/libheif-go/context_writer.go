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
/*
#include <stdlib.h>
#include <string.h>
#include <libheif/heif.h>
*/
import "C"

import (
	"io"
	"unsafe"
)

type writerData struct {
	w   io.Writer
	err error
}

//export writeGo
func writeGo(data *C.void, size C.size_t, userdata *C.void) C.struct_heif_error {
	writer := (*writerData)(unsafe.Pointer(userdata))
	if writer.err != nil {
		return C.struct_heif_error{
			code:    C.heif_error_Ok,
			subcode: C.heif_suberror_Unspecified,
		}
	}

	_, err := writer.w.Write(C.GoBytes(unsafe.Pointer(data), C.int(size)))
	if err != nil {
		writer.err = err
		return C.struct_heif_error{
			code:    C.heif_error_Usage_error,
			subcode: C.heif_suberror_Unspecified,
		}
	}

	return C.struct_heif_error{
		code:    C.heif_error_Ok,
		subcode: C.heif_suberror_Unspecified,
	}
}
