# go-exif-update

Go package to for updating EXIF data in JPEG files.

This is a thin wrapper around the dsoprea's [go-exif](https://github.com/dsoprea/go-exif) and [go-jpeg-image-structure](https://github.com/dsoprea/go-jpeg-image-structure) packages and includes command-line tools for updating the EXIF data JPEG files using key-value parameters as well as a WebAssembly (wasm) binary for updating EXIF data in JavaScript (or other languages that support wasm binaries).

The blog post [Updating EXIF metadata in JavaScript (and WebAssembly)](https://millsfield.sfomuseum.org/blog/2021/04/14/exif/) includes an example of the WebAssembly binary in action.

## Important

As of this writing the majority of EXIF tags are _not_ supported. Currently only EXIF tags of with the following EXIF types are supported: `ASCII`, `BYTE`, `RATIONAL` and `SRATIONAL` are supported. This is not ideal but I am still trying to get familiar with the requirements of the `go-exif` package. Contributions and patches for the other remaining EXIF tag types is welcomed.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/sfomuseum/go-exif-update.svg)](https://pkg.go.dev/github.com/sfomuseum/go-exif-update)

### Example

_Error handling removed for the sake of brevity._

```
package main

import (
	"flag"
	"github.com/sfomuseum/go-exif-update"
	"log"
	"os"
)

func main() {

	exif_props := map[string]interface{}{
		"Artist": "Bob",
		"Copyright": "SFO Museum",
	}
	
	for _, path := range paths {

		fh, _ := os.Open(path)
		defer fh.Close()

		update.UpdateExif(fh, os.Stdout, exif_props)
	}
}
```

The properties passed to the [UpdateExif](https://pkg.go.dev/github.com/sfomuseum/go-exif-update#UpdateExif) method are assumed to be suitably typed to be used as input for the [exif.SetStandardWithName](https://pkg.go.dev/github.com/dsoprea/go-exif/v3#IfdBuilder.SetStandardWithName) method. For anything other than simple strings this can become a fairly involved process. I am still trying the understand the `go-exif` documentation and requirements so any guidance I have to offer remains limited.

This package also provides a convenience [PrepareAndUpdateExif](https://pkg.go.dev/github.com/sfomuseum/go-exif-update#PrepareUpdateExif) method which attempts to translate EXIF property values in to their corresponding `go-exif` type before invoking the `UpdateExif` method. For example:

```
	exif_props := map[string]interface{}{
		"FNumber": "11/1",	// FNumber is stored as a `RATIONAL` value rather than as a string (or `ASCI`)
	}
	
	for _, path := range paths {

		fh, _ := os.Open(path)
		defer fh.Close()

		update.PrepareAndUpdateExif(fh, os.Stdout, exif_props)
	}
```

### Non-standard tags

The `PrepareAndUpdateExif` method also supports a limited set of custom non-standard tags. These are provided to simplify the process of assigning multiple or complex EXIF tags.

The list of currently supported non-standard tags is:

| Name | Value | Description | Notes |
| --- | --- | --- | --- |
| `X-Latitude` | float64 | Assign a decimal latitude value in to its corresponding `GPSLatitude` and `GPSLatitudeRef` properties. | |
| `X-Longitude` | float64 | Assign a decimal longitude value in to its corresponding `GPSLongitude` and `GPSLongitudeRef` properties. | |

## Tools

```
$> make cli
GOOS=js GOARCH=wasm go build -mod vendor -o www/wasm/update_exif.wasm cmd/update-exif-wasm/main.go
GOOS=js GOARCH=wasm go build -mod vendor -o www/wasm/supported_tags.wasm cmd/tags-supported-wasm/main.go
go build -mod vendor -o bin/tags-is-supported cmd/tags-is-supported/main.go
go build -mod vendor -o bin/tags-supported cmd/tags-supported/main.go
go build -mod vendor -o bin/update-exif cmd/update-exif/main.go
go build -mod vendor -o bin/server cmd/update-exif-server/main.go
```

As part of the build process for tools the two WebAssembly (wasm) binaries that are used by the `update-exif-server` tool are compiled and placee. You can also build the wasm binaries separately using the `wasm` Makefile target:

```
$> make wasm
GOOS=js GOARCH=wasm go build -mod vendor -o www/wasm/update_exif.wasm cmd/update-exif-wasm/main.go
GOOS=js GOARCH=wasm go build -mod vendor -o www/wasm/supported_tags.wasm cmd/tags-supported-wasm/main.go
```

### tags-is-supported

Command-line tool for indicating whether a named EXIF tag is supported by the sfomuseum/go-exif-update package.

```
$> ./bin/tags-is-supported -h
Command-line tool for indicating whether a named EXIF tag is supported by the sfomuseum/go-exif-update package.

Usage:
	./bin/tags-is-supported tag(N) tag(N) tag(N)
```

For example:

```
$> ./bin/tags-is-supported Copyright ImageWidth
Copyright true
ImageWidth false
```

_As of this writing non-standard tags, described above, are not considered when processing tags._

### tags-supported

Command-line tool that prints a list of EXIF tag names, sorted alphabetically, that are supported by the sfomuseum/go-exif-update package.

```
$> ./bin/tags-supported -h
Command-line tool that prints a list of EXIF tag names, sorted alphabetically, that are supported by the sfomuseum/go-exif-update package.

Usage:
	./bin/tags-supported
```

_As of this writing non-standard tags, described above, are not considered when listing supported tags._

### update-exif

Command-line tool for updating the EXIF properties in one or more JPEG images.

```
$> ./bin/update-exif -h
Command-line tool for updating the EXIF properties in one or more JPEG images. Images are not updated in place but written to STDOUT.

Usage:
	./bin/update-exif [options] image(N) image(N) image(N)

  -property value
    	One or more {TAG}={VALUE} EXIF properties. {TAG} must be a recognized EXIF tag.
```

For example:

```
$> ./bin/update-exif \
	-property Artist=Bob \
	-property Copyright=Alice \
	test1.jpg > test2.jpg
```

And then verifying the data using the `exiv2` tool:

```
$> exiv2 -pa test2.jpg 
Exif.Image.Artist                            Ascii       4  Bob
Exif.Image.Copyright                         Ascii       6  Alice
```

This tool uses the `PrepareAndUpdateExif` method and supports non-standard tags, both described above. For example:

```
$> ./bin/update-exif \
	-property X-Latitude=37.61799 \
	-property X-Longitude=-122.384864 \
	-property Artist=Walrus \
	fixtures/walrus.jpg \

> walrus-exif.jpg

$> exiv2 -pa walrus-exif.jpg 
Exif.Image.GPSTag                            Long        1  45
Exif.GPSInfo.GPSLongitude                    Rational    3  122deg 23' 5"
Exif.GPSInfo.GPSLongitudeRef                 Ascii       2  West
Exif.GPSInfo.GPSLatitude                     Rational    3  37deg 37' 4"
Exif.GPSInfo.GPSLatitudeRef                  Ascii       2  North
Exif.Image.Artist                            Ascii       7  Walrus
```

### update-exif-wasm-demo

HTTP server for demonstrating the use of the update_exif WebAssembly binary.

```
$> ./bin/update-exif-wasm-demo -h
HTTP server for demonstrating the use of the update_exif WebAssembly binary.

Usage:
	 ./bin/update-exif-wasm-demo [options]
  -bootstrap-prefix string
    	A relative path to append to all Bootstrap-related paths the server will listen for requests on.
  -server-uri string
    	A valid aaronland/go-http-server URI. (default "http://localhost:8080")
```	

For example:

```
$> ./bin/update-exif-wasm-demo 
2021/04/09 17:32:21 Listening on http://localhost:8080
```

This will start a web server on `localhost:8080` hosting a small web application for updating the EXIF tags in a static image. The application is bundled with two WebAssembly binaries:

* `supported_tags.wasm` that will return a JSON-encoded list of supported EXIF tags.
* `update_exif.wasm` that will update an image encoded in a base64 data URL from a JSON-encoded dictionary of EXIF tags and values.

You should see something like this:

![](docs/images/update-exif-wasm-demo.png)

_It's early days so this web application lacks any kind of polish._

Once the image has been uploaded it will trigger a download prompt (using eligrey's [FileSaver.js](https://github.com/eligrey/FileSaver.js). PNG and GIF images will be re-encoded as JPEG images dynamically.

![](docs/images/update-exif-wasm-demo-save.png)

An abbreviated version of the code to use the `update_exif.wasm` binary in JavaScript looks like this:

```
    var update = { "Artist": "Bob" };	
    var enc_update = JSON.stringify(update);
    
    var img = document.getElementById("image");
    
    var canvas = document.createElement("canvas");
    canvas.width = img.width;
    canvas.height = img.height;
    var ctx = canvas.getContext("2d");
    ctx.drawImage(img, 0, 0);
    var b64_img = canvas.toDataURL("image/jpeg", 1.0);
    
    update_exif(b64_img, enc_update).then(data_url => {

        var blob = dataURLToBlob(data_url);
    	saveAs(blob, "example.jpg");

    ).catch(err => {

    	console.log("Failed to update EXIF data, ", err);
    });
    
```

For a complete example consult the [www/javascript/index.js](www/javascript/index.js) file.

The WASM binary included with tool uses the `PrepareAndUpdateExif` method and supports non-standard tags, both described above.

## Supported tags

The following EXIF tags are supported by this package. This list was generated by the `tags-supported` tool.

```
AnalogBalance
AntiAliasStrength
ApertureValue
ApertureValue
Artist
AsShotPreProfileMatrix
AsShotProfileName
AsShotWhiteXY
BaselineExposure
BaselineNoise
BaselineSharpness
BatteryLevel
BestQualityScale
BlackLevel
BlackLevelDeltaH
BlackLevelDeltaV
BodySerialNumber
BrightnessValue
BrightnessValue
CFAPattern
CFAPlaneColor
CameraCalibration1
CameraCalibration2
CameraCalibrationSignature
CameraOwnerName
CameraSerialNumber
ChromaBlurRadius
ClipPath
ColorMatrix1
ColorMatrix2
CompressedBitsPerPixel
CompressedBitsPerPixel
Copyright
CurrentPreProfileMatrix
DNGBackwardVersion
DNGPrivateData
DNGVersion
DateTime
DateTimeDigitized
DateTimeOriginal
DateTimeOriginal
DefaultScale
DigitalZoomRatio
DocumentName
DotRange
ExposureBiasValue
ExposureBiasValue
ExposureIndex
ExposureIndex
ExposureTime
ExposureTime
FNumber
FNumber
FlashEnergy
FlashEnergy
FocalLength
FocalLength
FocalPlaneXResolution
FocalPlaneXResolution
FocalPlaneYResolution
FocalPlaneYResolution
ForwardMatrix1
ForwardMatrix2
GPSAltitude
GPSAltitudeRef
GPSDOP
GPSDateStamp
GPSDestBearing
GPSDestBearingRef
GPSDestDistance
GPSDestDistanceRef
GPSDestLatitude
GPSDestLatitudeRef
GPSDestLongitude
GPSDestLongitudeRef
GPSImgDirection
GPSImgDirectionRef
GPSLatitude
GPSLatitudeRef
GPSLongitude
GPSLongitudeRef
GPSMapDatum
GPSMeasureMode
GPSSatellites
GPSSpeed
GPSSpeedRef
GPSStatus
GPSTimeStamp
GPSTrack
GPSTrackRef
GPSVersionID
HostComputer
ImageDescription
ImageHistory
ImageID
ImageResources
ImageUniqueID
InkNames
InteroperabilityIndex
LensInfo
LensMake
LensModel
LensSerialNumber
LensSpecification
LinearResponseLimit
LocalizedCameraModel
Make
MaxApertureValue
MaxApertureValue
Model
NoiseProfile
NoiseReductionApplied
OriginalRawFileName
PreviewApplicationName
PreviewApplicationVersion
PreviewDateTime
PreviewSettingsDigest
PreviewSettingsName
PrimaryChromaticities
ProcessingSoftware
ProfileCalibrationSignature
ProfileCopyright
ProfileHueSatMapData1
ProfileHueSatMapData2
ProfileLookTableData
ProfileName
ProfileToneCurve
RawDataUniqueID
ReductionMatrix1
ReductionMatrix2
ReferenceBlackWhite
RelatedImageFileFormat
RelatedSoundFile
SecurityClassification
ShadowScale
ShutterSpeedValue
ShutterSpeedValue
Software
SpectralSensitivity
SubSecTime
SubSecTimeDigitized
SubSecTimeOriginal
SubjectDistance
TIFFEPStandardID
TargetPrinter
TimeZoneOffset
UniqueCameraModel
WhitePoint
XClipPathUnits
XMLPacket
XPAuthor
XPComment
XPKeywords
XPSubject
XPTitle
XResolution
YCbCrCoefficients
YClipPathUnits
YResolution
```

## See also

* https://github.com/dsoprea/go-exif
* https://github.com/dsoprea/go-jpeg-image-structure
* https://exiftool.org/TagNames/EXIF.html
* https://exiv2.org/tags.html
* https://www.media.mit.edu/pia/Research/deepview/exif.html
