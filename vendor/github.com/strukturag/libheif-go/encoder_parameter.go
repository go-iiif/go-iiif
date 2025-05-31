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
	"runtime"
	"strings"
)

// EncoderParameter defines a parameter that can be set on an encoder.
type EncoderParameter struct {
	param *C.struct_heif_encoder_parameter
}

func newEncoderParameter(p *C.struct_heif_encoder_parameter) EncoderParameter {
	return EncoderParameter{
		param: p,
	}
}

// String returns information on the encoder parameter.
func (p EncoderParameter) String() string {
	props := []string{
		"name=" + p.Name(),
	}
	if min, max, values, err := p.IntegerValues(); err == nil {
		if min != nil {
			props = append(props,
				fmt.Sprintf("min=%d", *min),
			)
		}
		if max != nil {
			props = append(props,
				fmt.Sprintf("max=%d", *max),
			)
		}
		if len(values) > 0 {
			props = append(props,
				fmt.Sprintf("values=%q", values),
			)
		}
	}
	if values, err := p.StringValues(); err == nil && len(values) > 0 {
		props = append(props,
			fmt.Sprintf("values=%q", values),
		)
	}
	return fmt.Sprintf("EncoderParameter[%s]",
		strings.Join(props, ", "),
	)
}

// Name returns the name of the encoder parameter.
func (p EncoderParameter) Name() string {
	defer runtime.KeepAlive(p)

	return C.GoString(C.heif_encoder_parameter_get_name(p.param))
}

// Type returns the type of the encoder parameter.
func (p EncoderParameter) Type() EncoderParameterType {
	defer runtime.KeepAlive(p)

	return EncoderParameterType(C.heif_encoder_parameter_get_type(p.param))
}

// IntegerValues returns the minimum and maximum values (if present) or a list
// of allowed values (if defined).
func (p EncoderParameter) IntegerValues() (*int, *int, []int, error) {
	defer runtime.KeepAlive(p)

	var have_minimum C.int
	var have_maximum C.int
	var min C.int
	var max C.int
	var num_values C.int
	var values_array *C.int
	err := C.heif_encoder_parameter_get_valid_integer_values(p.param,
		&have_minimum, &have_maximum,
		&min, &max,
		&num_values, &values_array,
	)
	if err := convertHeifError(err); err != nil {
		return nil, nil, nil, err
	}

	var minValue *int
	if have_minimum != 0 {
		minValue = makePointer(int(min))
	}

	var maxValue *int
	if have_maximum != 0 {
		maxValue = makePointer(int(max))
	}

	var values []int
	if num_values > 0 {
		values = make([]int, num_values)
		for i := 0; i < int(num_values); i++ {
			values[i] = int(*values_array)
			values_array = nextPointer(values_array)
		}
	}

	return minValue, maxValue, values, nil
}

// StringValues returns a list of allowed values (if defined).
func (p EncoderParameter) StringValues() ([]string, error) {
	defer runtime.KeepAlive(p)

	var values_array **C.char
	err := C.heif_encoder_parameter_get_valid_string_values(p.param, &values_array)
	if err := convertHeifError(err); err != nil {
		return nil, err
	}

	var values []string
	for *values_array != nil {
		values = append(values, C.GoString(*values_array))
		values_array = nextPointer(values_array)
	}

	return values, nil
}
