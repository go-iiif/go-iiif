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

// HeifError contains information about an error in libheif.
type HeifError struct {
	Code    ErrorCode
	Subcode SuberrorCode
	Message string
}

// Error returns the human readable error message.
func (e *HeifError) Error() string {
	return e.Message
}

func convertHeifError(cerror C.struct_heif_error) error {
	if ErrorCode(cerror.code) == ErrorOK {
		return nil
	}

	return &HeifError{
		Code:    ErrorCode(cerror.code),
		Subcode: SuberrorCode(cerror.subcode),
		Message: C.GoString(cerror.message),
	}
}
