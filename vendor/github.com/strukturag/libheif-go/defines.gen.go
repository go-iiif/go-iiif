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

const build_version = (1<<24) | (16<<16) | (2<<8)

type ErrorCode C.enum_heif_error_code

const (
	// Everything ok, no error occurred.
	ErrorOK                       ErrorCode = C.heif_error_Ok
	// Input file does not exist.
	ErrorInputDoesNotExist        ErrorCode = C.heif_error_Input_does_not_exist
	// Error in input file. Corrupted or invalid content.
	ErrorInvalidInput             ErrorCode = C.heif_error_Invalid_input
	// Input file type is not supported.
	ErrorUnsupportedFiletype      ErrorCode = C.heif_error_Unsupported_filetype
	// Image requires an unsupported decoder feature.
	ErrorUnsupportedFeature       ErrorCode = C.heif_error_Unsupported_feature
	// Library API has been used in an invalid way.
	ErrorUsage                    ErrorCode = C.heif_error_Usage_error
	// Could not allocate enough memory.
	ErrorMemoryAllocation         ErrorCode = C.heif_error_Memory_allocation_error
	// The decoder plugin generated an error
	ErrorDecoderPlugin            ErrorCode = C.heif_error_Decoder_plugin_error
	// The encoder plugin generated an error
	ErrorEncoderPlugin            ErrorCode = C.heif_error_Encoder_plugin_error
	// Error during encoding or when writing to the output
	ErrorEncoding                 ErrorCode = C.heif_error_Encoding_error
	// Application has asked for a color profile type that does not exist
	ErrorColorProfileDoesNotExist ErrorCode = C.heif_error_Color_profile_does_not_exist
	// Error loading a dynamic plugin
	ErrorPluginLoading            ErrorCode = C.heif_error_Plugin_loading_error
)

type SuberrorCode C.enum_heif_suberror_code

const (
	// no further information available
	SuberrorUnspecified                          SuberrorCode = C.heif_suberror_Unspecified
	// End of data reached unexpectedly.
	SuberrorEndOfData                            SuberrorCode = C.heif_suberror_End_of_data
	// Size of box (defined in header) is wrong
	SuberrorInvalidBoxSize                       SuberrorCode = C.heif_suberror_Invalid_box_size
	// Mandatory 'ftyp' box is missing
	SuberrorNoFtypBox                            SuberrorCode = C.heif_suberror_No_ftyp_box
	SuberrorNoIdatBox                            SuberrorCode = C.heif_suberror_No_idat_box
	SuberrorNoMetaBox                            SuberrorCode = C.heif_suberror_No_meta_box
	SuberrorNoHdlrBox                            SuberrorCode = C.heif_suberror_No_hdlr_box
	SuberrorNoHvcCBox                            SuberrorCode = C.heif_suberror_No_hvcC_box
	SuberrorNoPitmBox                            SuberrorCode = C.heif_suberror_No_pitm_box
	SuberrorNoIpcoBox                            SuberrorCode = C.heif_suberror_No_ipco_box
	SuberrorNoIpmaBox                            SuberrorCode = C.heif_suberror_No_ipma_box
	SuberrorNoIlocBox                            SuberrorCode = C.heif_suberror_No_iloc_box
	SuberrorNoIinfBox                            SuberrorCode = C.heif_suberror_No_iinf_box
	SuberrorNoIprpBox                            SuberrorCode = C.heif_suberror_No_iprp_box
	SuberrorNoIrefBox                            SuberrorCode = C.heif_suberror_No_iref_box
	SuberrorNoPictHandler                        SuberrorCode = C.heif_suberror_No_pict_handler
	// An item property referenced in the 'ipma' box is not existing in the 'ipco' container.
	SuberrorIpmaBoxReferencesNonexistingProperty SuberrorCode = C.heif_suberror_Ipma_box_references_nonexisting_property
	// No properties have been assigned to an item.
	SuberrorNoPropertiesAssignedToItem           SuberrorCode = C.heif_suberror_No_properties_assigned_to_item
	// Image has no (compressed) data
	SuberrorNoItemData                           SuberrorCode = C.heif_suberror_No_item_data
	// Invalid specification of image grid (tiled image)
	SuberrorInvalidGridData                      SuberrorCode = C.heif_suberror_Invalid_grid_data
	// Tile-images in a grid image are missing
	SuberrorMissingGridImages                    SuberrorCode = C.heif_suberror_Missing_grid_images
	SuberrorInvalidCleanAperture                 SuberrorCode = C.heif_suberror_Invalid_clean_aperture
	// Invalid specification of overlay image
	SuberrorInvalidOverlayData                   SuberrorCode = C.heif_suberror_Invalid_overlay_data
	// Overlay image completely outside of visible canvas area
	SuberrorOverlayImageOutsideOfCanvas          SuberrorCode = C.heif_suberror_Overlay_image_outside_of_canvas
	SuberrorAuxiliaryImageTypeUnspecified        SuberrorCode = C.heif_suberror_Auxiliary_image_type_unspecified
	SuberrorNoOrInvalidPrimaryItem               SuberrorCode = C.heif_suberror_No_or_invalid_primary_item
	SuberrorNoInfeBox                            SuberrorCode = C.heif_suberror_No_infe_box
	SuberrorUnknownColorProfileType              SuberrorCode = C.heif_suberror_Unknown_color_profile_type
	SuberrorWrongTileImageChromaFormat           SuberrorCode = C.heif_suberror_Wrong_tile_image_chroma_format
	SuberrorInvalidFractionalNumber              SuberrorCode = C.heif_suberror_Invalid_fractional_number
	SuberrorInvalidImageSize                     SuberrorCode = C.heif_suberror_Invalid_image_size
	SuberrorInvalidPixiBox                       SuberrorCode = C.heif_suberror_Invalid_pixi_box
	SuberrorNoAV1CBox                            SuberrorCode = C.heif_suberror_No_av1C_box
	SuberrorWrongTileImagePixelDepth             SuberrorCode = C.heif_suberror_Wrong_tile_image_pixel_depth
	SuberrorUnknownNCLXColorPrimaries            SuberrorCode = C.heif_suberror_Unknown_NCLX_color_primaries
	SuberrorUnknownNCLXTransferCharacteristics   SuberrorCode = C.heif_suberror_Unknown_NCLX_transfer_characteristics
	SuberrorUnknownNCLXMatrixCoefficients        SuberrorCode = C.heif_suberror_Unknown_NCLX_matrix_coefficients
	// Invalid specification of region item
	SuberrorInvalidRegionData                    SuberrorCode = C.heif_suberror_Invalid_region_data
	// A security limit preventing unreasonable memory allocations was exceeded by the input file.
	// Please check whether the file is valid. If it is, contact us so that we could increase the
	// security limits further.
	SuberrorSecurityLimitExceeded                SuberrorCode = C.heif_suberror_Security_limit_exceeded
	// also used for Invalid_input
	SuberrorNonexistingItemReferenced            SuberrorCode = C.heif_suberror_Nonexisting_item_referenced
	// An API argument was given a NULL pointer, which is not allowed for that function.
	SuberrorNullPointerArgument                  SuberrorCode = C.heif_suberror_Null_pointer_argument
	// Image channel referenced that does not exist in the image
	SuberrorNonexistingImageChannelReferenced    SuberrorCode = C.heif_suberror_Nonexisting_image_channel_referenced
	// The version of the passed plugin is not supported.
	SuberrorUnsupportedPluginVersion             SuberrorCode = C.heif_suberror_Unsupported_plugin_version
	// The version of the passed writer is not supported.
	SuberrorUnsupportedWriterVersion             SuberrorCode = C.heif_suberror_Unsupported_writer_version
	// The given (encoder) parameter name does not exist.
	SuberrorUnsupportedParameter                 SuberrorCode = C.heif_suberror_Unsupported_parameter
	// The value for the given parameter is not in the valid range.
	SuberrorInvalidParameterValue                SuberrorCode = C.heif_suberror_Invalid_parameter_value
	// Error in property specification
	SuberrorInvalidProperty                      SuberrorCode = C.heif_suberror_Invalid_property
	// Image reference cycle found in iref
	SuberrorItemReferenceCycle                   SuberrorCode = C.heif_suberror_Item_reference_cycle
	// Image was coded with an unsupported compression method.
	SuberrorUnsupportedCodec                     SuberrorCode = C.heif_suberror_Unsupported_codec
	// Image is specified in an unknown way, e.g. as tiled grid image (which is supported)
	SuberrorUnsupportedImageType                 SuberrorCode = C.heif_suberror_Unsupported_image_type
	SuberrorUnsupportedDataVersion               SuberrorCode = C.heif_suberror_Unsupported_data_version
	// The conversion of the source image to the requested chroma / colorspace is not supported.
	SuberrorUnsupportedColorConversion           SuberrorCode = C.heif_suberror_Unsupported_color_conversion
	SuberrorUnsupportedItemConstructionMethod    SuberrorCode = C.heif_suberror_Unsupported_item_construction_method
	SuberrorUnsupportedHeaderCompressionMethod   SuberrorCode = C.heif_suberror_Unsupported_header_compression_method
	SuberrorUnsupportedBitDepth                  SuberrorCode = C.heif_suberror_Unsupported_bit_depth
	SuberrorCannotWriteOutputData                SuberrorCode = C.heif_suberror_Cannot_write_output_data
	SuberrorEncoderInitialization                SuberrorCode = C.heif_suberror_Encoder_initialization
	SuberrorEncoderEncoding                      SuberrorCode = C.heif_suberror_Encoder_encoding
	SuberrorEncoderCleanup                       SuberrorCode = C.heif_suberror_Encoder_cleanup
	SuberrorTooManyRegions                       SuberrorCode = C.heif_suberror_Too_many_regions
	// a specific plugin file cannot be loaded
	SuberrorPluginLoadingError                   SuberrorCode = C.heif_suberror_Plugin_loading_error
	// trying to remove a plugin that is not loaded
	SuberrorPluginIsNotLoaded                    SuberrorCode = C.heif_suberror_Plugin_is_not_loaded
	// error while scanning the directory for plugins
	SuberrorCannotReadPluginDirectory            SuberrorCode = C.heif_suberror_Cannot_read_plugin_directory
)

type CompressionFormat C.enum_heif_compression_format

const (
	CompressionUndefined    CompressionFormat = C.heif_compression_undefined
	CompressionHEVC         CompressionFormat = C.heif_compression_HEVC
	CompressionAVC          CompressionFormat = C.heif_compression_AVC
	CompressionJPEG         CompressionFormat = C.heif_compression_JPEG
	CompressionAV1          CompressionFormat = C.heif_compression_AV1
	CompressionVVC          CompressionFormat = C.heif_compression_VVC
	CompressionEVC          CompressionFormat = C.heif_compression_EVC
	// ISO/IEC 15444-16:2021
	CompressionJPEG2000     CompressionFormat = C.heif_compression_JPEG2000
	// ISO/IEC 23001-17:2023
	CompressionUncompressed CompressionFormat = C.heif_compression_uncompressed
)

type Chroma C.enum_heif_chroma

const (
	ChromaUndefined              Chroma = C.heif_chroma_undefined
	ChromaMonochrome             Chroma = C.heif_chroma_monochrome
	Chroma420                    Chroma = C.heif_chroma_420
	Chroma422                    Chroma = C.heif_chroma_422
	Chroma444                    Chroma = C.heif_chroma_444
	ChromaInterleavedRGB         Chroma = C.heif_chroma_interleaved_RGB
	ChromaInterleavedRGBA        Chroma = C.heif_chroma_interleaved_RGBA
	// HDR, big endian.
	ChromaInterleavedRRGGBB_BE   Chroma = C.heif_chroma_interleaved_RRGGBB_BE
	// HDR, big endian.
	ChromaInterleavedRRGGBBAA_BE Chroma = C.heif_chroma_interleaved_RRGGBBAA_BE
	// HDR, little endian.
	ChromaInterleavedRRGGBB_LE   Chroma = C.heif_chroma_interleaved_RRGGBB_LE
	// HDR, little endian.
	ChromaInterleavedRRGGBBAA_LE Chroma = C.heif_chroma_interleaved_RRGGBBAA_LE
)

type Colorspace C.enum_heif_colorspace

const (
	ColorspaceUndefined  Colorspace = C.heif_colorspace_undefined
	// heif_colorspace_YCbCr should be used with one of these heif_chroma values:
	// * heif_chroma_444
	// * heif_chroma_422
	// * heif_chroma_420
	ColorspaceYCbCr      Colorspace = C.heif_colorspace_YCbCr
	// heif_colorspace_RGB should be used with one of these heif_chroma values:
	// * heif_chroma_444 (for planar RGB)
	// * heif_chroma_interleaved_RGB
	// * heif_chroma_interleaved_RGBA
	// * heif_chroma_interleaved_RRGGBB_BE
	// * heif_chroma_interleaved_RRGGBBAA_BE
	// * heif_chroma_interleaved_RRGGBB_LE
	// * heif_chroma_interleaved_RRGGBBAA_LE
	ColorspaceRGB        Colorspace = C.heif_colorspace_RGB
	// heif_colorspace_monochrome should only be used with heif_chroma = heif_chroma_monochrome
	ColorspaceMonochrome Colorspace = C.heif_colorspace_monochrome
)

type Channel C.enum_heif_channel

const (
	ChannelY           Channel = C.heif_channel_Y
	ChannelCb          Channel = C.heif_channel_Cb
	ChannelCr          Channel = C.heif_channel_Cr
	ChannelR           Channel = C.heif_channel_R
	ChannelG           Channel = C.heif_channel_G
	ChannelB           Channel = C.heif_channel_B
	ChannelAlpha       Channel = C.heif_channel_Alpha
	ChannelInterleaved Channel = C.heif_channel_interleaved
)

type ProgressStep C.enum_heif_progress_step

const (
	ProgressStepTotal    ProgressStep = C.heif_progress_step_total
	ProgressStepLoadTile ProgressStep = C.heif_progress_step_load_tile
)

type ChromaDownsamplingAlgorithm C.enum_heif_chroma_downsampling_algorithm

const (
	ChromaDownsamplingNearestNeighbor ChromaDownsamplingAlgorithm = C.heif_chroma_downsampling_nearest_neighbor
	ChromaDownsamplingAverage         ChromaDownsamplingAlgorithm = C.heif_chroma_downsampling_average
	// Combine with 'heif_chroma_upsampling_bilinear' for best quality.
	// Makes edges look sharper when using YUV 420 with bilinear chroma upsampling.
	ChromaDownsamplingSharpYuv        ChromaDownsamplingAlgorithm = C.heif_chroma_downsampling_sharp_yuv
)

type ChromaUpsamplingAlgorithm C.enum_heif_chroma_upsampling_algorithm

const (
	ChromaUpsamplingNearestNeighbor ChromaUpsamplingAlgorithm = C.heif_chroma_upsampling_nearest_neighbor
	ChromaUpsamplingBilinear        ChromaUpsamplingAlgorithm = C.heif_chroma_upsampling_bilinear
)

type EncoderParameterType C.enum_heif_encoder_parameter_type

const (
	EncoderParameterTypeInteger EncoderParameterType = C.heif_encoder_parameter_type_integer
	EncoderParameterTypeBoolean EncoderParameterType = C.heif_encoder_parameter_type_boolean
	EncoderParameterTypeString  EncoderParameterType = C.heif_encoder_parameter_type_string
)
