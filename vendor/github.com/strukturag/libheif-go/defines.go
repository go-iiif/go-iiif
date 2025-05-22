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

//go:generate scripts/generate-defines.py defines.gen.go

const (
	// ChromaInterleaved24Bit is an alias for ChromaInterleavedRGB.
	//
	// Deprecated: use ChromaInterleavedRGB instead
	ChromaInterleaved24Bit = ChromaInterleavedRGB
	// ChromaInterleaved32Bit is an alias for ChromaInterleavedRGBA
	//
	// Deprecated: use ChromaInterleavedRGBA instead
	ChromaInterleaved32Bit = ChromaInterleavedRGBA
)

type LosslessMode int

const (
	LosslessModeDisabled LosslessMode = iota
	LosslessModeEnabled
)

type LoggingLevel int

const (
	LoggingLevelNone LoggingLevel = iota
	LoggingLevelBasic
	LoggingLevelAdvanced
	LoggingLevelFull
)
