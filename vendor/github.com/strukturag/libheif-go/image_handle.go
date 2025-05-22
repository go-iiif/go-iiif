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
	"runtime"
	"unsafe"
)

// ImageHandle contains information about an image in a libheif Context.
type ImageHandle struct {
	handle *C.struct_heif_image_handle
}

func freeHeifImageHandle(c *ImageHandle) {
	C.heif_image_handle_release(c.handle)
	c.handle = nil
}

// IsPrimaryImage checks if the image handle is for a primary image.
func (h *ImageHandle) IsPrimaryImage() bool {
	defer runtime.KeepAlive(h)

	return C.heif_image_handle_is_primary_image(h.handle) != 0
}

// GetWidth returns the width of the image handle.
func (h *ImageHandle) GetWidth() int {
	defer runtime.KeepAlive(h)

	return int(C.heif_image_handle_get_width(h.handle))
}

// GetHeight returns the height of the image handle.
func (h *ImageHandle) GetHeight() int {
	defer runtime.KeepAlive(h)

	return int(C.heif_image_handle_get_height(h.handle))
}

// HasAlphaChannel checks if the image handle has an alpha channel.
func (h *ImageHandle) HasAlphaChannel() bool {
	defer runtime.KeepAlive(h)

	return C.heif_image_handle_has_alpha_channel(h.handle) != 0
}

// HasDepthImage checks if the image handle has a depth channel.
func (h *ImageHandle) HasDepthImage() bool {
	defer runtime.KeepAlive(h)

	return C.heif_image_handle_has_depth_image(h.handle) != 0
}

// GetNumberOfDepthImages returns the number of depth images in the image handle.
func (h *ImageHandle) GetNumberOfDepthImages() int {
	defer runtime.KeepAlive(h)

	return int(C.heif_image_handle_get_number_of_depth_images(h.handle))
}

// GetListOfDepthImageIDs returns the list of depth image ids in the image handle.
func (h *ImageHandle) GetListOfDepthImageIDs() []int {
	defer runtime.KeepAlive(h)

	num := int(C.heif_image_handle_get_number_of_depth_images(h.handle))
	if num == 0 {
		return []int{}
	}

	origIDs := make([]C.heif_item_id, num)
	C.heif_image_handle_get_list_of_depth_image_IDs(h.handle, &origIDs[0], C.int(num))
	return convertItemIDs(origIDs, num)
}

// GetDepthImageHandle returns the image handle for the given depth image id.
func (h *ImageHandle) GetDepthImageHandle(depth_image_id int) (*ImageHandle, error) {
	defer runtime.KeepAlive(h)

	var handle ImageHandle
	err := C.heif_image_handle_get_depth_image_handle(h.handle, C.heif_item_id(depth_image_id), &handle.handle)
	if err := convertHeifError(err); err != nil {
		return nil, err
	}

	runtime.SetFinalizer(&handle, freeHeifImageHandle)
	return &handle, nil
}

// GetNumberOfThumbnails returns the number of thumbnails in the image handle.
func (h *ImageHandle) GetNumberOfThumbnails() int {
	defer runtime.KeepAlive(h)

	return int(C.heif_image_handle_get_number_of_thumbnails(h.handle))
}

// GetListOfThumbnailIDs returns the list of thumbnail ids in the image handle.
func (h *ImageHandle) GetListOfThumbnailIDs() []int {
	defer runtime.KeepAlive(h)

	num := int(C.heif_image_handle_get_number_of_thumbnails(h.handle))
	if num == 0 {
		return []int{}
	}

	origIDs := make([]C.heif_item_id, num)
	C.heif_image_handle_get_list_of_thumbnail_IDs(h.handle, &origIDs[0], C.int(num))
	return convertItemIDs(origIDs, num)
}

// GetThumbnail returns the image handle for the given thumbnail id.
func (h *ImageHandle) GetThumbnail(thumbnail_id int) (*ImageHandle, error) {
	defer runtime.KeepAlive(h)

	var handle ImageHandle
	err := C.heif_image_handle_get_thumbnail(h.handle, C.heif_item_id(thumbnail_id), &handle.handle)
	runtime.SetFinalizer(&handle, freeHeifImageHandle)
	return &handle, convertHeifError(err)
}

// DecodeImage decodes the image to the provided colorspace and chroma.
func (h *ImageHandle) DecodeImage(colorspace Colorspace, chroma Chroma, options *DecodingOptions) (*Image, error) {
	defer runtime.KeepAlive(h)

	var image Image

	var opt *C.struct_heif_decoding_options
	if options != nil {
		opt = options.options
	}

	err := C.heif_decode_image(h.handle, &image.image, uint32(colorspace), uint32(chroma), opt)
	if err := convertHeifError(err); err != nil {
		return nil, err
	}

	runtime.SetFinalizer(&image, freeHeifImage)
	return &image, nil
}

func (h *ImageHandle) GetMetadataBlockIDs(filter string) []int {
	defer runtime.KeepAlive(h)

	var f *C.char
	if filter != "" {
		f = C.CString(filter)
		defer C.free(unsafe.Pointer(f))
	}
	num := int(C.heif_image_handle_get_number_of_metadata_blocks(h.handle, f))
	if num == 0 {
		return nil
	}

	ids := make([]C.heif_item_id, num)
	C.heif_image_handle_get_list_of_metadata_block_IDs(h.handle, f, &ids[0], C.int(num))
	return convertItemIDs(ids, num)
}

func (h *ImageHandle) GetMetadataContentType(block_id int) string {
	defer runtime.KeepAlive(h)

	ct := C.heif_image_handle_get_metadata_content_type(h.handle, C.heif_item_id(block_id))
	if ct == nil {
		return ""
	}

	return C.GoString(ct)
}

func (h *ImageHandle) GetMetadata(block_id int) ([]byte, error) {
	defer runtime.KeepAlive(h)

	var err C.struct_heif_error
	var result []byte
	if size := C.heif_image_handle_get_metadata_size(h.handle, C.heif_item_id(block_id)); size > 0 {
		result = make([]byte, size)
		err = C.heif_image_handle_get_metadata(h.handle, C.heif_item_id(block_id), unsafe.Pointer(&result[0]))
	} else {
		err = C.heif_image_handle_get_metadata(h.handle, C.heif_item_id(block_id), nil)
	}
	if err := convertHeifError(err); err != nil {
		return nil, err
	}

	return result, nil
}
