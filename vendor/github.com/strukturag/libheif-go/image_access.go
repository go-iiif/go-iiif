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
	"unsafe"
)

// ImageAccess contains information on how to access raw image data.
type ImageAccess struct {
	Plane    []byte
	planePtr unsafe.Pointer
	Stride   int
	height   int

	image *Image // need this reference to make sure the image is not GC'ed while we access it
}

func (i *ImageAccess) setData(data []byte, stride int) {
	dataPtr := unsafe.Pointer(&data[0])
	// Handle common case directly
	if stride == i.Stride {
		dstP := i.planePtr
		srcP := dataPtr
		C.memcpy(dstP, srcP, C.size_t(i.height*stride))
	} else {
		for y := 0; y < i.height; y++ {
			dstP := unsafe.Add(i.planePtr, y*i.Stride)
			srcP := unsafe.Add(dataPtr, y*stride)
			C.memcpy(dstP, srcP, C.size_t(stride))
		}
	}
	i.Plane = C.GoBytes(i.planePtr, C.int(i.height*i.Stride))
}
