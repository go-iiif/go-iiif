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

extern struct heif_error writeGo(void* data, size_t size, void* userdata);

struct heif_error writeCgo(struct heif_context* ctx, const void* data, size_t size, void* userdata) {
	struct heif_error err = writeGo((char*)data, size, userdata);
	if (!err.message) {
		switch (err.code) {
			case heif_error_Ok:
				err.message = "Success";
				break;
			default:
				err.message = "Error writing";
				break;
		}
	}
	return err;
}
*/
import "C"

import (
	"errors"
	"fmt"
	"io"
	"runtime"
	"unsafe"
)

// Context is a libheif context.
type Context struct {
	context *C.struct_heif_context
}

// NewContext creates a new libheif context that can be used for decoding and
// encoding of images.
func NewContext() (*Context, error) {
	if err := checkLibraryVersion(); err != nil {
		return nil, err
	}

	ctx := &Context{
		context: C.heif_context_alloc(),
	}
	if ctx.context == nil {
		return nil, errors.New("Could not allocate context")
	}

	runtime.SetFinalizer(ctx, freeHeifContext)
	return ctx, nil
}

func freeHeifContext(c *Context) {
	C.heif_context_free(c.context)
	c.context = nil
}

// ReadFromFile loads the image from the given filename in the current context.
func (c *Context) ReadFromFile(filename string) error {
	defer runtime.KeepAlive(c)

	c_filename := C.CString(filename)
	defer C.free(unsafe.Pointer(c_filename))

	err := C.heif_context_read_from_file(c.context, c_filename, nil)
	return convertHeifError(err)
}

// ReadFromMemory loads the image from the given data in the current context.
func (c *Context) ReadFromMemory(data []byte) error {
	defer runtime.KeepAlive(c)

	// TODO: Use reader API internally.
	err := C.heif_context_read_from_memory(c.context, unsafe.Pointer(&data[0]), C.size_t(len(data)), nil)
	return convertHeifError(err)
}

func (c *Context) convertEncoderDescriptor(d *C.struct_heif_encoder_descriptor) (*Encoder, error) {
	defer runtime.KeepAlive(c)

	cid := C.heif_encoder_descriptor_get_id_name(d)
	cname := C.heif_encoder_descriptor_get_name(d)
	enc := &Encoder{
		id:   C.GoString(cid),
		name: C.GoString(cname),
	}
	err := C.heif_context_get_encoder(c.context, d, &enc.encoder)
	if err := convertHeifError(err); err != nil {
		return nil, err
	}

	runtime.SetFinalizer(enc, freeHeifEncoder)
	return enc, nil
}

// NewEncoder creates a new encoder with the given compression format.
func (c *Context) NewEncoder(compression CompressionFormat) (*Encoder, error) {
	defer runtime.KeepAlive(c)

	const max = 1
	descriptors := make([]*C.struct_heif_encoder_descriptor, max)
	num := int(C.heif_context_get_encoder_descriptors(c.context, uint32(compression), nil, &descriptors[0], C.int(max)))
	if num == 0 {
		return nil, fmt.Errorf("no encoder for compression %v", compression)
	}

	return c.convertEncoderDescriptor(descriptors[0])
}

// Write saves the current image.
func (c *Context) Write(w io.Writer) error {
	defer runtime.KeepAlive(c)

	writer := &C.struct_heif_writer{
		writer_api_version: 1,

		write: (*[0]byte)(C.writeCgo),
	}
	writerData := &writerData{
		w: w,
	}

	var p runtime.Pinner
	p.Pin(w)
	defer p.Unpin()

	err := C.heif_context_write(c.context, writer, unsafe.Pointer(writerData))
	if writerData.err != nil {
		// Bubble up error returned by passed io.Writer
		return writerData.err
	}

	return convertHeifError(err)
}

// WriteToFile saves the current image to the given file.
func (c *Context) WriteToFile(filename string) error {
	defer runtime.KeepAlive(c)

	err := C.heif_context_write_to_file(c.context, C.CString(filename))
	return convertHeifError(err)
}

// GetNumberOfTopLevelImages returns the number of top-level images.
func (c *Context) GetNumberOfTopLevelImages() int {
	defer runtime.KeepAlive(c)

	i := int(C.heif_context_get_number_of_top_level_images(c.context))
	return i
}

// IsTopLevelImageID checks if a given id is a top-level image.
func (c *Context) IsTopLevelImageID(ID int) bool {
	defer runtime.KeepAlive(c)

	ok := C.heif_context_is_top_level_image_ID(c.context, C.heif_item_id(ID)) != 0
	return ok
}

// GetListOfTopLevelImageIDs returns a list of top-level image ids.
func (c *Context) GetListOfTopLevelImageIDs() []int {
	defer runtime.KeepAlive(c)

	num := int(C.heif_context_get_number_of_top_level_images(c.context))
	if num == 0 {
		return []int{}
	}

	origIDs := make([]C.heif_item_id, num)
	C.heif_context_get_list_of_top_level_image_IDs(c.context, &origIDs[0], C.int(num))
	return convertItemIDs(origIDs, num)
}

// GetPrimaryImageID returns the id of the primary image.
func (c *Context) GetPrimaryImageID() (int, error) {
	defer runtime.KeepAlive(c)

	var id C.heif_item_id
	err := C.heif_context_get_primary_image_ID(c.context, &id)
	if err := convertHeifError(err); err != nil {
		return 0, err
	}

	return int(id), nil
}

// GetPrimaryImageHandle returns the image handle of the primary image.
func (c *Context) GetPrimaryImageHandle() (*ImageHandle, error) {
	defer runtime.KeepAlive(c)

	var handle ImageHandle
	err := C.heif_context_get_primary_image_handle(c.context, &handle.handle)
	if err := convertHeifError(err); err != nil {
		return nil, err
	}

	runtime.SetFinalizer(&handle, freeHeifImageHandle)
	return &handle, convertHeifError(err)
}

// GetImageHandle returns the image handle of the given image id.
func (c *Context) GetImageHandle(id int) (*ImageHandle, error) {
	defer runtime.KeepAlive(c)

	var handle ImageHandle
	err := C.heif_context_get_image_handle(c.context, C.heif_item_id(id), &handle.handle)
	if err := convertHeifError(err); err != nil {
		return nil, err
	}

	runtime.SetFinalizer(&handle, freeHeifImageHandle)
	return &handle, nil
}

func (c *Context) AddExifMetadata(handle *ImageHandle, data []byte) error {
	runtime.KeepAlive(c)
	runtime.KeepAlive(handle)

	dataPtr := unsafe.Pointer(&data[0])
	err := C.heif_context_add_exif_metadata(c.context, handle.handle, dataPtr, C.int(len(data)))
	return convertHeifError(err)
}

func (c *Context) AddXmpMetadata(handle *ImageHandle, data []byte) error {
	runtime.KeepAlive(c)
	runtime.KeepAlive(handle)

	dataPtr := unsafe.Pointer(&data[0])
	err := C.heif_context_add_XMP_metadata(c.context, handle.handle, dataPtr, C.int(len(data)))
	return convertHeifError(err)
}

func (c *Context) AddGenericMetadata(handle *ImageHandle, data []byte, item_type string, content_type string) error {
	runtime.KeepAlive(c)
	runtime.KeepAlive(handle)

	dataPtr := unsafe.Pointer(&data[0])
	var it *C.char
	if item_type != "" {
		it = C.CString(item_type)
		defer C.free(unsafe.Pointer(it))
	}
	var ct *C.char
	if content_type != "" {
		ct = C.CString(content_type)
		defer C.free(unsafe.Pointer(ct))
	}
	err := C.heif_context_add_generic_metadata(c.context, handle.handle, dataPtr, C.int(len(data)), it, ct)
	return convertHeifError(err)
}
