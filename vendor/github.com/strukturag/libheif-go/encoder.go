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

const (
	// encoderParamStringSize is the maximum length of a string value that is
	// returned by the GetParameter / GetParameterString functions.
	encoderParamStringSize = 1024
)

// Encoder contains a libheif encoder object.
type Encoder struct {
	encoder *C.struct_heif_encoder
	id      string
	name    string
}

func freeHeifEncoder(enc *Encoder) {
	C.heif_encoder_release(enc.encoder)
	enc.encoder = nil
}

// ID returns the id of the encoder.
func (e *Encoder) ID() string {
	return e.id
}

// Name returns the name of the encoder.
func (e *Encoder) Name() string {
	return e.name
}

// SetQuality sets the quality level for the encoder.
func (e *Encoder) SetQuality(q int) error {
	defer runtime.KeepAlive(e)

	err := C.heif_encoder_set_lossy_quality(e.encoder, C.int(q))
	return convertHeifError(err)
}

// SetLossless enables or disables the lossless encoding mode.
func (e *Encoder) SetLossless(l LosslessMode) error {
	defer runtime.KeepAlive(e)

	err := C.heif_encoder_set_lossless(e.encoder, C.int(l))
	return convertHeifError(err)
}

// SetLoggingLevel sets the logging level to use for encoding.
func (e *Encoder) SetLoggingLevel(l LoggingLevel) error {
	defer runtime.KeepAlive(e)

	err := C.heif_encoder_set_logging_level(e.encoder, C.int(l))
	return convertHeifError(err)
}

// ListParameters returns a list of parameters this encoder supports.
func (e *Encoder) ListParameters() []EncoderParameter {
	defer runtime.KeepAlive(e)

	parameters := C.heif_encoder_list_parameters(e.encoder)
	if parameters == nil {
		return nil
	}

	var result []EncoderParameter
	for *parameters != nil {
		result = append(result, newEncoderParameter(*parameters))
		parameters = nextPointer(parameters)
	}
	return result
}

// SetParameter sets a parameter of any type to the string value.
// Integer values are parsed from the string.
// Boolean values can be "true", "false", "1", "0".
//
// x265 encoder specific note:
// When using the x265 encoder, you may pass any of its parameters by
// prefixing the parameter name with 'x265:'. Hence, to set the 'ctu' parameter,
// you will have to set 'x265:ctu' in libheif.
// Note that there is no checking for valid parameters when using the prefix.
func (e *Encoder) SetParameter(name string, value string) error {
	defer runtime.KeepAlive(e)

	err := C.heif_encoder_set_parameter(e.encoder, C.CString(name), C.CString(value))
	return convertHeifError(err)
}

// GetParameter returns the current value of a parameter of any type as a human
// readable string.
// The returned string is compatible with SetParameter().
func (e *Encoder) GetParameter(name string) (string, error) {
	defer runtime.KeepAlive(e)

	value := (*C.char)(C.malloc(encoderParamStringSize))
	if value == nil {
		return "", errors.New("can't allocate memory for value")
	}

	defer C.free(unsafe.Pointer(value))
	err := C.heif_encoder_get_parameter(e.encoder, C.CString(name), value, encoderParamStringSize)
	if err := convertHeifError(err); err != nil {
		return "", err
	}

	return C.GoString(value), nil
}

// HasDefault returns true if a parameter has a default value.
func (e *Encoder) HasDefault(name string) bool {
	defer runtime.KeepAlive(e)

	return C.heif_encoder_has_default(e.encoder, C.CString(name)) != 0
}

// SetParameterInteger sets the integer parameter.
func (e *Encoder) SetParameterInteger(name string, value int) error {
	defer runtime.KeepAlive(e)

	err := C.heif_encoder_set_parameter_integer(e.encoder, C.CString(name), C.int(value))
	return convertHeifError(err)
}

// GetParameterInteger returns the value of the integer parameter.
func (e *Encoder) GetParameterInteger(name string) (int, error) {
	defer runtime.KeepAlive(e)

	var value C.int
	err := C.heif_encoder_get_parameter_integer(e.encoder, C.CString(name), &value)
	if err := convertHeifError(err); err != nil {
		return 0, err
	}

	return int(value), nil
}

// SetParameterBool sets the boolean parameter.
func (e *Encoder) SetParameterBool(name string, value bool) error {
	defer runtime.KeepAlive(e)

	err := C.heif_encoder_set_parameter_boolean(e.encoder, C.CString(name), convertBool[C.int](value))
	return convertHeifError(err)
}

// GetParameterBool returns the value of the boolean parameter.
func (e *Encoder) GetParameterBool(name string) (bool, error) {
	defer runtime.KeepAlive(e)

	var value C.int
	err := C.heif_encoder_get_parameter_boolean(e.encoder, C.CString(name), &value)
	if err := convertHeifError(err); err != nil {
		return false, err
	}

	return value != 0, nil
}

// SetParameterString sets the string parameter.
func (e *Encoder) SetParameterString(name string, value string) error {
	defer runtime.KeepAlive(e)

	err := C.heif_encoder_set_parameter_string(e.encoder, C.CString(name), C.CString(value))
	return convertHeifError(err)
}

// GetParameterString returns the value of the string parameter.
func (e *Encoder) GetParameterString(name string) (string, error) {
	defer runtime.KeepAlive(e)

	value := (*C.char)(C.malloc(encoderParamStringSize))
	if value == nil {
		return "", errors.New("can't allocate memory for value")
	}

	defer C.free(unsafe.Pointer(value))
	err := C.heif_encoder_get_parameter_string(e.encoder, C.CString(name), value, encoderParamStringSize)
	if err := convertHeifError(err); err != nil {
		return "", err
	}

	return C.GoString(value), nil
}
