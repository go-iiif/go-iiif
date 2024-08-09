# prism

[![PkgGoDev](https://pkg.go.dev/badge/github.com/mandykoh/prism)](https://pkg.go.dev/github.com/mandykoh/prism)
[![Go Report Card](https://goreportcard.com/badge/github.com/mandykoh/prism)](https://goreportcard.com/report/github.com/mandykoh/prism)
[![Build Status](https://travis-ci.org/mandykoh/prism.svg?branch=main)](https://travis-ci.org/mandykoh/prism)

`prism` aims to become a set of utilities for practical colour management and conversion in pure Go.

`prism` currently implements:

* Encoding/decoding linear colour from sRGB, Adobe RGB, Pro Photo RGB, and Display P3 encodings
* Fast LUT-based tonal response encoding/decoding
* Conversion to and from CIE xyY, CIE XYZ, and CIE Lab
* Chromatic adaptation in XYZ space between different white points
* Extracting metadata (including ICC profile) from PNG, JPEG, and WebP files

Still missing:

* Embedding of tagged colour profiles in image
* Exposing colour data from ICC profiles (to enable conversions between arbitrary profiles)
* Rendering intent support
* CMYK support

See the [API documentation](https://pkg.go.dev/github.com/mandykoh/prism) for more details.

This software is made available under an [MIT license](LICENSE).

Much of this implementation is based on information provided by [Bruce Lindbloom](http://www.brucelindbloom.com) and [Charles Poynton](http://poynton.ca), among many others who generously contribute to public edification on the esoteric science of colour.


## Rationale

Using the analogy of working with strings, a _colour space_ is to _colour encoding_ what a character set is to character encoding.

Image data provided by the standard [`image`](https://golang.org/pkg/image/) package doesn’t come with colour space information, yet colour values from the images are invariably encoded in some scheme specific to the colour space that the image is targeting. And as with character encodings, programmers then routinely take these RGB (for example) values and make incorrect assumptions about how they’ve been encoded, leading to software that sort of mostly works in common cases but fails in important ways (like supporting non-English languages or Emoji).

`prism` can be used to decode these colour values into a normalised, linear representation more suitable for image processing, and subsequently converting those back to encoded colour values in (potentially) other colour spaces.


### sRGB is the Web standard; why do I still need to worry about this?

This is like asking why we need to worry about UTF-8 when ASCII is the standard, or why we need to worry about other fonts when Times New Roman is the standard. However, there are two prominent reasons:


#### 1. Wide gamut imaging is becoming commonplace

sRGB is a narrow gamut colour space. Smartphones, computing devices, and other displays (including everything marketed under “HDR” consumer labels) increasingly use wider gamuts, and are capable of reproducing much more saturation than sRGB can represent. Images produced on these displays, taking advantage of the wider gamuts, will look incorrect when naively interpreted as sRGB.

The following example shows an image targeting Adobe RGB (a wide gamut colour space commonly used by artists and photographers; left) and what happens when the same image is incorrectly assumed to be sRGB (right). Note the loss of saturation—a common complaint with images uploaded to social media or other sites. The bright, saturated topping has become much more dull and unappetising, and the whole image has gained an unpleasant greenish cast:

![Example of incorrectly interpreting an Adobe RGB image as sRGB](doc-images/example-bad-conversion.png)

_This is not a deficiency of sRGB._ This image is well within the sRGB gamut, and a correct interpretation will look just like the version on the left (indeed, this example figure itself is actually sRGB).

Another way of stating this problem is that sRGB being the Web standard makes it the default, but doesn’t preclude other, significantly different colour spaces from common use.


#### 2. sRGB uses a non-linear tonal response curve

For efficiency and fidelity, nearly all colour encoding schemes (sRGB or otherwise) are set up to be _perceptually_ linear. But because our eyes don’t perceive brightness linearly, this means the [colour values are not linear in intensity](https://blog.johnnovak.net/2016/09/21/what-every-coder-should-know-about-gamma/), so sRGB(127, 127, 127) is not actually half as bright as sRGB(255, 255, 255).

The following example shows an image being resampled to half size, without and with corrected linear colour. The corrected example shows that the checkerboard patches are consistent in brightness with the solid patches after resizing, while the resized patches in the incorrect example end up as different colours:

![Example of correct and incorrect resampling](doc-images/example-resampling.png)

Since many image manipulation operations (such as scaling, sharpening, or blending) rely on colour values having linear intensity, applying them to non-linear colour data produces visual artefacts and generally incorrect results.

Another way of stating this problem is that colour values in images are _encoded_ (sometimes referred to as “gamma encoding” or “gamma correction”), and need to be _decoded_ rather than used directly.


## Example usage


### Metadata extraction

Image metadata can be extracted from images without needing to consume the entire image stream. Currently, this is supported for JPEG and PNG data. The following example demonstrates this using [`autometa.Load`](https://pkg.go.dev/github.com/mandykoh/prism/meta/autometa?tab=doc#Load):

```go
// Get a meta.Data instance containing image details and a stream to the full image
md, imgStream, err := autometa.Load(inFile)
if err != nil {
    panic(err)
}

// The metadata specifies the image format
imgFormat := md.Format // eg. jpeg.Format

// Load the full image after extracting metadata
img, err = jpeg.Decode(imgStream)
if err != nil {
    panic(err)
}
```

Included in the metadata are basic details about the image such as the pixel dimensions and colour depth:

```go
fmt.Printf("Image format: %s\n", md.Format)
fmt.Printf("Image height in pixels: %d\n", md.PixelHeight)
fmt.Printf("Image width in pixels: %d\n", md.PixelWidth)
fmt.Printf("Bits per component: %d\n", md.BitsPerComponent)
```

The stream returned by `autometa.Load` reproduces the full image stream, so that it can be later passed to (for example) `jpeg.Decode` to load the rest of the image. This allows information like the size of the image to be known before having to load an extremely large image.

If the image contained an ICC profile, it can be retrieved from the metadata:

```go
iccProfile, err := md.ICCProfile()
description, err := iccProfile.Description()  // eg. "sRGB IEC61966-2.1"
```

If no profile exists, `nil` is returned without an error.

`autometa.Load` delegates to format-specific loaders like `jpegmeta.Load` and `pngmeta.Load`; these can be used instead if you know the format of image.


### Colour linearisation

An image can be easily converted from its colour space encoding (eg sRGB) to a linear encoding. Because this operation can be lossy in 8-bit colour depths, it’s a good idea to first convert images to 16-bit colour (eg instances of `image.NRGBA64` or `image.RGBA64`). `prism` provides utility functions for such conversions:

```go
img = prism.ConvertImageToRGBA64(img, parallelism)
```

Then the image can be linearised (here, using itself as both source and destination):

```go
srgb.LineariseImage(img, img, parallelism)
```

Alternatively a new blank 16-bit image can be created and the original image linearised into it:

```go
linearisedImg := image.NewRGBA64(img.Bounds())
srgb.LineariseImage(linearisedImg, img, parallelism)
```

The image can then be passed to operations that expect an `image.Image` but assume linear colour. Here we pass it to the `BiLinear` rescaler to reduce the image to half its original size, which will now produce a correct result in linear space:

```go
resampled := image.NewRGBA64(image.Rect(0, 0, img.Rect.Dx()/2, img.Rect.Dy()/2))
draw.BiLinear.Scale(resampled, resampled.Rect, img, img.Bounds(), draw.Src, nil)
```

Note that the output is still linearised, so before writing the image to an output file (eg in PNG or JPEG format), we need to re-encode it back to sRGB space, and probably also want to convert it back to 8-bit colour:

```go
encodedImg := image.NewRGBA(resampled.Bounds())
srgb.EncodeImage(encodedImg, resampled, parallelism)
```


### Colour conversion

Conversions between RGB colour spaces are performed via the CIE XYZ intermediate colour space (using the `ToXYZ` and `ColorFromXYZ` functions).

The following example converts Adobe RGB (1998) pixel data to sRGB. It retrieves a pixel from an [NRGBA image](https://golang.org/pkg/image/#NRGBA), decodes it to an Adobe RGB (1998) linearised colour value, then converts that to an sRGB colour value via CIE XYZ, before finally encoding the result as an 8-bit sRGB value suitable for writing back to an `image.NRGBA`:

```go
c := inputImg.NRGBAAt(x, y)                 // Take input colour value
ac, alpha := adobergb.ColorFromNRGBA(c)     // Interpret image pixel as Adobe RGB and convert to linear representation
sc := srgb.ColorFromXYZ(ac.ToXYZ())         // Convert to XYZ, then from XYZ to sRGB linear representation
outputImg.SetNRGBA(x, y, sc.ToNRGBA(alpha)) // Write sRGB-encoded value to output image
```


### Chromatic adaptation

Adobe RGB (1998) and sRGB are both specified referring to a standard D65 white point. However, Pro Photo RGB references a D50 white point. When converting between white points, a chromatic adaptation is required to compensate for a shift in warmness/coolness that would otherwise occur.

The following example prepares such a chromatic adaptation (using the [`AdaptBetweenXYYWhitePoints`](https://pkg.go.dev/github.com/mandykoh/prism/ciexyz?tab=doc#AdaptBetweenXYYWhitePoints) function), then uses it in converting from Pro Photo RGB to sRGB:

```go
adaptation := ciexyz.AdaptBetweenXYYWhitePoints(
    prophotorgb.StandardWhitePoint,         // From D50
    srgb.StandardWhitePoint,                // To D65
)

c := inputImg.NRGBAAt(x, y)                 // Take input colour value
pc, alpha := prophotorgb.ColorFromNRGBA(c)  // Interpret image pixel as Pro Photo RGB and convert to linear representation

xyz := pc.ToXYZ()                           // Convert from Pro Photo RGB to CIE XYZ
xyz = adaptation.Apply(xyz)                 // Apply chromatic adaptation from D50 to D65

sc := srgb.ColorFromXYZ(xyz)                // Convert from CIE XYZ to sRGB linear representation
outputImg.SetNRGBA(x, y, sc.ToNRGBA(alpha)) // Write sRGB-encoded value to output image
```
