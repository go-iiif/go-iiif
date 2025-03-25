# Config files

There is a [sample config file](../docs/config.json.example) included with this repo. The easiest way to understand config files is that they consist of at least five top-level groupings, with nested section-specific details, followed by zero or more implementation specific configuration blocks. The five core blocks are:

### level

```
	"level": {
		"compliance": "2"
	}
```

Indicates which level of IIIF Image API compliance the server (or associated tools) should support. Basically, there is no reason to ever change this right now.

### profile

```
    "profile": {
    	"services": {
		    ...
	} 
    }
```	

Additional configurations for a IIIF profile (aka `info.json`). Currently this is limited to defining one or more addtional `services` to append to a profile.

#### services

```
    "profile": {
    	"services": {
		    "enable": [ "palette" ]
	} 
    }
```

Services configurations are currently limited to enabling a fixed set of named services, where that fixed set numbers exactly three:

* `blurhash` for generateing a compact base-83 encoded representation of an image using the [BlurHash](https://github.com/woltapp/blurhash/blob/master/Algorithm.md) algorithm.
* `imagehash` for generating average and difference perceptual hashes of an image (as defined by the `imagehash` configuration below).
* `palette` for extracting a colour palette for an image (as defined by the `palette` configuration below).

As of this writing adding custom services is a nuisance. [There is an open issue](https://github.com/go-iiif/go-iiif/issues/71) to address this problem, but no ETA yet for its completion.

##### blurhash

```
    "blurhash": {
    	"x": 4,
	"y": 3,
	"size": 32
    }
```

`go-iiif` uses the [go-blurhash](https://github.com/buckket/go-blurhash) to generate a compact base-83 encoded representation of an image using the [BlurHash](https://github.com/woltapp/blurhash/blob/master/Algorithm.md) algorithm.

The blurhash service configuration has no specific properties as of this writing.

* **x** is the number of BlurHash components along the `x` axis.
* **y** is the number of BlurHash components along the `y` axis.
* **size** is the maximum dimension to resize an image to before attempting to generate a BlurHash.

Sample out for the `blurhash` service is included [below](#non-standard-services).

##### imagehash

```
    "imagehash": {}
```

`go-iiif` uses the [goimagehash](https://github.com/corona10/goimagehash) to extract [average](http://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html) and [difference](http://www.hackerfactor.com/blog/index.php?/archives/529-Kind-of-Like-That.html) perceptual hashes.

The imagehash service configuration has no specific properties as of this writing.

Sample out for the `imagehash` service is included [below](#non-standard-services).

##### palette

```
    "palette": {
    	"extruder": { "name": "vibrant", "count": 5 },
    	"grid": { "name": "euclidian" },
	"palettes": [
		    { "name": "crayola" },
		    { "name": "css4" }
        ]
    }
```

`go-iiif` uses the [go-colours](https://github.com/aaronland/go-colours) package to extract colours. `go-colours` itself is a work in progress so you should approach colours extraction as a service accordingly.

A palette service configuration has the following properties:

* **extruder** is a simple dictionary with a `name` and a `count` property. Since there is currently only one extruder (defined by `go-colours`) there is no need to change this.
* **grid** is a simple dictionary with a `name` property. Since there is currently only one grid (defined by `go-colours`) there is no need to change this.
* **palettes**  is a list of simple dictionaries, each of which has a `name` property. Valid names are: `crayola`, `css3` or `css4`.

Sample out for the `palette` service is included [below](#non-standard-services).

### graphics

```
	"graphics": {
		"driver": "native://"
	}
```

`native://` is the default graphics driver. Others include:

* https://github.com/go-iiif/go-iiif-vips which uses [libvips](https://github.com/jcupitt/libvips) to perform all image-processing functions.

### features

```
	"features": {
		"enable": {},
		"disable": { "rotation": [ "rotationArbitrary"] },
		"append": {}
	}
```

The `features` block allows you to enable or disable specific IIIF features. _Currently only image related features may be manipulated._

For example the level 2 spec does not say GIF outputs is required so the level 2 compliance definition in `go-iiif` disables it by default. If you are using a graphics engine (not `libvips` though) that can produce GIF files you would enable it here.

Likewise you may need to disable a feature that is supported by not required or features that are required but can't be used for one reason or another. For example `libvips` does not allow support for the following features: `sizeByDistortedWh (size), rotationArbitrary (rotation), bitonal (quality)`.

Finally, maybe you've got an IIIF implementation that [knows how to do things not defined in the spec](https://github.com/go-iiif/go-iiif/issues/1). This is also where you would add them.

#### compliance

Here's how that dynamic plays out in reality. The table below lists all the IIIF parameters and their associate features. Each feature lists its syntax and whether or not it is required and supported [according to the official spec](compliance/level2.go) but then also according to the [example `go-iiif` config file](../docs/config.json.example), included with this repo.

_This table was generated using the [iiif-dump-config](cmd/iiif-dump-config.go) tool and if anyone can tell me how to make Markdown tables (in GitHub) render colours I would be grateful._

##### [region](http://iiif.io/api/image/2.1/index.html#region)
| feature | syntax | required (spec) | supported (spec) | required (config) | supported (config) |
|---|---|---|---|---|---|
| **full** | full | <span style="color:green;">true</span> | <span style="color:green;">true</span> | <span style="color:green;">**true**</span> | <span style="color:green;">**true**</span> |
| **regionByPct** | pct:x,y,w,h | <span style="color:green;">true</span> | <span style="color:green;">true</span> | <span style="color:green;">**true**</span> | <span style="color:green;">**true**</span> |
| **regionByPx** | x,y,w,h | <span style="color:green;">true</span> | <span style="color:green;">true</span> | <span style="color:green;">**true**</span> | <span style="color:green;">**true**</span> |
| **regionSquare** | square | <span style="color:red;">false</span> | <span style="color:red;">false</span> | <span style="color:green;">**true**</span> | <span style="color:red;">false</span> |

##### [size](http://iiif.io/api/image/2.1/index.html#size)
| feature | syntax | required (spec) | supported (spec) | required (config) | supported (config) |
|---|---|---|---|---|---|
| **full** | full | <span style="color:green;">true</span> | <span style="color:green;">true</span> | <span style="color:green;">**true**</span> | <span style="color:green;">**true**</span> |
| **max** | max | <span style="color:red;">false</span> | <span style="color:green;">true</span> | <span style="color:red;">false</span> | <span style="color:green;">**true**</span> |
| **sizeByConfinedWh** | !w,h | <span style="color:green;">true</span> | <span style="color:green;">true</span> | <span style="color:green;">**true**</span> | <span style="color:green;">**true**</span> |
| **sizeByDistortedWh** | w,h | <span style="color:green;">true</span> | <span style="color:green;">true</span> | <span style="color:green;">**true**</span> | <span style="color:red;">false</span> |
| **sizeByH** | ,h | <span style="color:green;">true</span> | <span style="color:green;">true</span> | <span style="color:green;">**true**</span> | <span style="color:green;">**true**</span> |
| **sizeByPct** | pct:n | <span style="color:green;">true</span> | <span style="color:green;">true</span> | <span style="color:green;">**true**</span> | <span style="color:green;">**true**</span> |
| **sizeByW** | w, | <span style="color:green;">true</span> | <span style="color:green;">true</span> | <span style="color:green;">**true**</span> | <span style="color:green;">**true**</span> |
| **sizeByWh** | w,h | <span style="color:green;">true</span> | <span style="color:green;">true</span> | <span style="color:green;">**true**</span> | <span style="color:green;">**true**</span> |

##### [rotation](http://iiif.io/api/image/2.1/index.html#rotation)
| feature | syntax | required (spec) | supported (spec) | required (config) | supported (config) |
|---|---|---|---|---|---|
| **mirroring** | !n | <span style="color:green;">true</span> | <span style="color:green;">true</span> | <span style="color:green;">**true**</span> | <span style="color:green;">**true**</span> |
| **none** | 0 | <span style="color:green;">true</span> | <span style="color:green;">true</span> | <span style="color:green;">**true**</span> | <span style="color:green;">**true**</span> |
| **rotationArbitrary** |  | <span style="color:red;">false</span> | <span style="color:green;">true</span> | <span style="color:red;">false</span> | <span style="color:red;">false</span> |
| **rotationBy90s** | 90,180,270 | <span style="color:green;">true</span> | <span style="color:green;">true</span> | <span style="color:green;">**true**</span> | <span style="color:green;">**true**</span> |
| **noAutoRotate** | -1 | <span style="color:red;">false</span> | <span style="color:red;">false</span> | <span style="color:red;">**false**</span> | <span style="color:green;">**true**</span> |

##### [quality](http://iiif.io/api/image/2.1/index.html#quality)
| feature | syntax | required (spec) | supported (spec) | required (config) | supported (config) |
|---|---|---|---|---|---|
| **bitonal** | bitonal | <span style="color:green;">true</span> | <span style="color:green;">true</span> | <span style="color:green;">**true**</span> | <span style="color:red;">false</span> |
| **color** | color | <span style="color:red;">false</span> | <span style="color:green;">true</span> | <span style="color:red;">false</span> | <span style="color:green;">**true**</span> |
| **default** | default | <span style="color:green;">true</span> | <span style="color:green;">true</span> | <span style="color:green;">**true**</span> | <span style="color:green;">**true**</span> |
| **dither** | dither | <span style="color:red;">false</span> | <span style="color:green;">false</span> | <span style="color:red;">false</span> | <span style="color:green;">**true**</span> |
| **gray** | gray | <span style="color:red;">false</span> | <span style="color:red;">false</span> | <span style="color:red;">false</span> | <span style="color:red;">false</span> |

_Careful readers may notice the presence of an undefined (by the IIIF spec) feature named `dither`. This is a `go-iiif` -ism and discussed in detail below in the [features.append](#featuresappend) and [non-standard features](#non-standard-features) sections._

##### [format](http://iiif.io/api/image/2.1/index.html#format)
| feature | syntax | required (spec) | supported (spec) | required (config) | supported (config) |
|---|---|---|---|---|---|
| **gif** | gif | <span style="color:red;">false</span> | <span style="color:red;">false</span> | <span style="color:red;">false</span> | <span style="color:red;">false</span> |
| **jp2** | jp2 | <span style="color:red;">false</span> | <span style="color:red;">false</span> | <span style="color:red;">false</span> | <span style="color:red;">false</span> |
| **jpg** | jpg | <span style="color:green;">true</span> | <span style="color:green;">true</span> | <span style="color:green;">**true**</span> | <span style="color:green;">**true**</span> |
| **pdf** | pdf | <span style="color:red;">false</span> | <span style="color:red;">false</span> | <span style="color:red;">false</span> | <span style="color:red;">false</span> |
| **png** | png | <span style="color:green;">true</span> | <span style="color:green;">true</span> | <span style="color:green;">**true**</span> | <span style="color:green;">**true**</span> |
| **tif** | tif | <span style="color:red;">false</span> | <span style="color:red;">false</span> | <span style="color:red;">false</span> | <span style="color:green;">**true**</span> |
| **webp** | webp | <span style="color:red;">false</span> | <span style="color:red;">false</span> | <span style="color:red;">false</span> | <span style="color:green;">**true**</span> |


_Support for GIF output is not enabled by default because it is not currently supported by `bimg` (the Go library on top of `lipvips`). There is however native support for converting final images to be GIFs but you will need to [enable that by hand](https://github.com/go-iiif/go-iiif/tree/primitive#featuresenable), below._

#### features.enable

```
	"features": {
		"enable": {
			"size": [ "max" ],
			"format": [ "webp", "tif" ]
		}
	}
```

Individual features for a given parameter are enabled by including the parameter name as a key to the `features.enabled` dictionary whose value is a list of specific feature names to enable.

#### features.disable

```
	"features": {
		"disable": {
			"size": [ "sizeByDistortedWh" ] ,
			"rotation": [ "rotationArbitrary" ],
			"quality": [ "bitonal" ]
		}
	}
```

Individual features for a given parameter are disabled by including the parameter name as a key to the `features.disabled` dictionary whose value is a list of specific feature names to disabled.

#### features.append

```
	"features": {
		"append": { "quality": {
			"dither": { "syntax": "dither", "required": false, "supported": true, "match": "^dither$" }
		}}
	}
```

New features are added by including their corresponding parameter name as a key to the `features.append` dictionary whose value is a model for that feature. The data model for new features to append looks like this:

```
	NAME (STRING): {
		"syntax": SYNTAX (STRING),
		"required": BOOLEAN,
		"supported": BOOLEAN,
		"match": REGULAR_EXPRESSION (STRING)
	}

```

All keys are required.

The `supported` key is used to determine whether a given feature is enabled or not. The `match` key is used to validate user input and should be a valid regular expression that will match that value. For example here is the compliance definition for images returned in the JPEG format:

```
		"format": {
	     	       "jpg": { "syntax": "jpg",  "required": true, "supported": true, "match": "^jpe?g$" }
		}
```

_Important: It is left to you to actually implement support for new features in the code for whichever graphics engine you are using. If you don't then any new features will be ignored at best or cause fatal errors at worst._

### images

```
	"images": {
		"source": { "uri", "file:///example/images" },
		"cache": { "uri": "memory://" }
	}
```

Image configuration details consist of a "source" and a "cache" target. Each target is defined by a URI which identifies code implementing the `go-iiif/source.Source` and `go-iiif/cache.Cache` interfaecs respectively.

#### images.source

Where to find source images.



```
	"images": {
		"source": { "uri": "file:///example/images" }
	}
```

Fetch sources images from any supported [Go Cloud storage service](https://gocloud.dev/howto/blob/#services). By default support for [local filesystem (`file://`)](https://gocloud.dev/howto/blob/#local), [in-memory (`mem://`)](https://gocloud.dev/howto/blob/#local) and [S3 (`s3://`](https://gocloud.dev/howto/blob/#s3) "bucket" sources is enabled.

There is also support for the [github.com/aaronland/gocloud-blob/s3](https://github.com/aaronland/gocloud-blob/s3) package which allows for AWS credentials to be defined using [string labels](https://github.com/aaronland/go-aws-auth?tab=readme-ov-file#credentials).

```
	"images": {
		"source": { "uri": "s3blob:///bucket-name?region=us-east-1&credentials=iam:" }
	}
```

##### Flickr

```
	"images": {
		"source": { "uri": "flickr://?client-uri=oauth1://?consumer_key={KEY}&consumer_secret={SECRET}" },
		"cache": { "uri": "memory://" }
	}
```

Fetch source images from Flickr. You will need to provide a valid [Flickr API key](https://www.flickr.com/services/api/). A few caveats:

* Under the hood the code is using the [aaronland/go-flickr-api](https://github.com/aaronland/go-flickr-api) package which uses a URI-style syntax for defining client instances. Please consult [the `go-flickr-api` documentation](https://github.com/aaronland/go-flickr-api?tab=readme-ov-file#oauth1) for details on how to construct those URIs.
* The code calls the [flickr.photos.getSizes](https://www.flickr.com/services/api/flickr.photos.getSizes.html) API method and looks for the first of the following photo sizes in this order: `Original, Large 2048, Large 1600, Large`. If none are available then an error is triggered.
* Photo size lookups are not cached yet.

Here's an example [with this photo](https://www.flickr.com/photos/straup/4136870023/in/album-72157622883263698/):

![](misc/go-iiif-flickr.png)

![](misc/go-iiif-flickr-detail.png)

##### URI Templates (RFC 6570)

```
	"images": {
		"source": { "uri": "rfc6570://template=https://images.collection.cooperhewitt.org/{id}" }
	}
```

Fetch source images from a remote URI. The `template` parameter must be a valid (Level 4) [URI Template](http://tools.ietf.org/html/rfc6570) with an `{id}` placeholder.

#### images.cache

Caching options for source images.

##### Blob (or `gocloud.dev/blob`)

```
	"images": {
		"cache": { "uri": "file:///example/images" }
	}
```

Read and write cached images from any supported [Go Cloud storage service](https://gocloud.dev/howto/blob/#services). By default support for [local filesystem (`file://`)](https://gocloud.dev/howto/blob/#local), [in-memory (`mem://`)](https://gocloud.dev/howto/blob/#local) and [S3 (`s3://`](https://gocloud.dev/howto/blob/#s3) "bucket" sources is enabled.

There is also support for the [github.com/aaronland/gocloud-blob/s3](https://github.com/aaronland/gocloud-blob/s3) package which allows for AWS credentials to be defined using [string labels](https://github.com/aaronland/go-aws-auth?tab=readme-ov-file#credentials). This package is	identified using the `s3blob://` URI scheme.

```
	"images": {
		"source": { "uri": "s3blob:///bucket-name?region=us-east-1&credentials=iam:" }
	}
```

##### Memory

```
	"images": {
		"cache": { "uri": "memory://" }
	}
```

##### Null

```
	"images": {
		"cache": { "uri": "null://" }
	}
```

Because you must define a caching layer this is here to satify the requirements without actually caching anything, anywhere.

### derivatives

```
	"derivatives": {
		"cache": { "uri": "file://example/cache" }
	}
```

Details about derivative images.

#### derivatives.cache

Caching options for derivative images.

##### Blob (or `gocloud.dev/blob`)

```
	"derivatives": {
		"cache": { "uri": "file:///example/images" }
	}
```

Read and write derivative images from any supported [Go Cloud storage service](https://gocloud.dev/howto/blob/#services). By default support for [local filesystem (`file://`)](https://gocloud.dev/howto/blob/#local), [in-memory (`mem://`)](https://gocloud.dev/howto/blob/#local) and [S3 (`s3://`](https://gocloud.dev/howto/blob/#s3) "bucket" sources is enabled.

There is also support for the [github.com/aaronland/gocloud-blob/s3](https://github.com/aaronland/gocloud-blob/s3) package which allows for AWS credentials to be defined using [string labels](https://github.com/aaronland/go-aws-auth?tab=readme-ov-file#credentials). This package is identified using the `s3blob://` URI scheme. This package is	identified using the `s3blob://` URI scheme.

```
	"derivatives": {
		"cache": { "uri": "s3blob:///bucket-name?region=us-east-1&credentials=iam:" }
	}
```

##### Memory

```
	"derivatives": {
		"cache": { "uri": "memory://" }
	}
```

##### Null

```
	"derivatives": {
		"cache": { "uri": "null://" }
	}
```

Because you must define a caching layer this is here to satify the requirements without actually caching anything, anywhere.

## Non-standard features

### Non-standard region features

#### regionByPx (and "smart" cropping)

If you are using `VIPS` as a [graphics engine](#graphics) and pass a `regionByPx` instruction whose X and Y values are `-1` then the code will ask libvips to crop the image (to the dimensions defined in the W and H values) centered on whatever libvips thinks it the most interesting or relevant part of the image.

See also: https://github.com/jcupitt/libvips/issues/317

### Non-standard rotation features

`go-iiif` supports the following non-standard IIIF `rotation` features:

#### noAutoRotate

```
	"enable": {
	    "rotation": [ "noAutoRotate" ]
	}
```

If the `noAutoRotate` feature is enabled this will act as a signal to the underlying image processing library to _not_ auto-rotate images according to the EXIF `Orientation` property (assuming it is present).

This feature exists because both the `libvips` library and the `bimg` wrapper code enable auto-rotation by default but neither updates the EXIF `Orientation` property to reflect the change so every time the newly created image is read by a piece of software that supports auto-rotation (including this one) that image will be doubly-rotated (and then triply-rotated and so on...)

If the `noAutoRotate` feature is enabled is can be triggered by setting the rotation element of your request URI to be `-1`, for example:

```
https://example.com/example.jpg/{REGION}/{SIZE}/-1/{QUALITY}.{FORMAT}
```

_As of this writing the `noAutoRotate` feature does not work in combination with other rotation commands (for example `-1,180` or equivalent, meaning "do not auto-rotate but please still rotate 180 degrees") but it probably should._

### Non-standard quality features

`go-iiif` supports the following non-standard IIIF `quality` features:

#### "Crisp"-ing

```
	"append": {
	    "quality": {
			"crisp": { "syntax": "crisp", "required": false, "supported": true, "match": "^crisp(?:\\:(\\d+\\.\\d+),(\\d+\\.\\d+),(\\d+\\.\\d+))?$"
	    }
	}
```

`crisp` will apply an "UnsharpMask" filter followed by a "Median" filter on an image using the [bild/effect](https://github.com/anthonynsimon/bild/#effect) package.

The `crisp` filter takes three positional parameters:

| Position | Name | Default |
| --- | --- | --- |
| 1 | Unsharp Mask Radius | 2.0 |
| 2 | Unsharp Mask Amount | 0.5 |
| 3 | Mediam Radius | 0.025 |

For example, this:

```
http://localhost:8080/spanking-cat.jpg/-1,-1,320,320/full/0/crisp:10.0,2.0,0.05.png
```

Would produce the following image:

![spanking cat](misc/go-iiif-crisp.png)

#### Dithering

```
	"append": {
		"quality": {
			"dither": { "syntax": "dither", "required": false, "supported": true, "match": "^dither$" }
		}
	}
```

`dither` will create a black and white [halftone](https://en.wikipedia.org/wiki/Halftone) derivative of an image using the [Atkinson dithering algorithm](https://en.wikipedia.org/wiki/Dither#Algorithms). Dithering is enabled in the [example config file](../docs/config.json.example) and you can invoke it like this:

```
http://localhost:8082/184512_5f7f47e5b3c66207_x.jpg/pct:41,7,40,70/,5000/0/dither.png
```

And here's what you should see, keeping in mind that this screenshot shows only a section of the image at full size:

![spanking cat](misc/go-iiif-dither-detail.png)

There are a few caveats about dithering images:

* The first thing to know is that the dithering is a [pure Go implementation](https://github.com/koyachi/go-atkinson) so it's not handled by `lipvips`.
* The second is that the dithering happens _after_ the `libvips` processing.
* This is relevant because there are some image formats where Go does not support native encoding. For example [webp](https://godoc.org/golang.org/x/image/webp) (which is weird since it's a Google thing...)
* It is possible to track all of this stuff in code and juggle output formats and reprocessing (in `libvips`) but that code has not been written yet.
* So you will need to track the sometimes still-rocky relationship between features and output formats yourself.

#### Primitive-ing

```
	"features": {
		"append": {
			"quality": {
				"primitive": { "syntax": "primitive:mode,iterations,alpha", "required": false, "supported": true, "match": "^primitive\\:[0-5]\\,\\d+\\,\\d+$" }
			}
		}
	},
	"primitive": { "max_iterations": 100 }
```

_Note the way the `primitive` block is a top-level element in your config file._

`primitive` use [@fogleman's primitive library](https://github.com/fogleman/primitive) to reproduce the final image using geometric primitives. Like this:

![](misc/go-iiif-primitive-circles.png)

The syntax for invoking this feature is `primitive:{MODE},{ITERATIONS},{ALPHA}` where:

* **MODE** is a number between 0-5 representing which of the [primitive shapes](https://github.com/fogleman/primitive#primitives) to use. They are:
 * 0: combo
 * 1: triangle
 * 2: rectangle
 * 3: ellipse
 * 4: circle
 * 5: rotated rectangle
* **ITERATIONS** is a number between 1 and infinity (a bad idea) or 1 and the number defined in the `primitive.max_iterations` section in your config file
* **ALPHA** is a number between 0-255

For example:

```
http://localhost:8082/184512_5f7f47e5b3c66207_x.jpg/full/500,/0/primitive:5,200,255.jpg
```

Be aware that it's not exactly "fast". It's [getting faster](https://github.com/fogleman/primitive/commit/ccd349008eb7c611d690c4dd1fd9bca74b86ceb1) but it still takes a while. Also, _this code_ should probably have a flag to downsize the input image for processing (and then resizing it back up to the requested size) but that doesn't happen yet. Basically you should not enable this feature as a public-facing web service because it will take seconds (not microseconds) or sometimes even minutes to render a single 256x256 tile. For example:

```
./bin/iiif-server -host 0.0.0.0 -config config.json
2016/09/21 15:43:08 Serving [::]:8080 with pid 5877
2016/09/21 15:43:13 starting model at 2016-09-21 15:43:13.626117993 +0000 UTC
2016/09/21 15:43:13 finished step 1 in 8.229683ms
2016/09/21 15:43:16 finished step 2 in 3.019413861s
…
2016/09/21 15:45:38 finished step 100 in 2m24.626232387s
2016/09/21 15:45:39 finished model in 2m25.611790848s
```

But it is pretty darn cool!

![](misc/go-iiif-primitive-triangles.png)

![](misc/go-iiif-primitive-triangles-detail.png)

If you specify a `gif` format parameter then `go-iiif` will return an animated GIF for the requested image consisting of each intermediate stage that the `primitive` library generated the final image. For example:

```
http://localhost:8082/184512_5f7f47e5b3c66207_x.jpg/full/500,/0/primitive:5,100,255.gif
```

Which would produce this:

![](../misc/go-iiif-primitive-animated-rect.gif)

Here are examples where each of the tiles in an slippy image are animated GIFs:

* https://go-iiif.github.io/go-iiif/animated/
* https://go-iiif.github.io/go-iiif/animated/?mode=circles
* https://go-iiif.github.io/go-iiif/animated/?mode=triangles

_Note: You will need to [manually enable support for GIF images](https://github.com/go-iiif/go-iiif/tree/primitive#featuresenable) in your config file for animated GIFs to work._

## Non-standard services

Documentation for non-standard services has been moved in to [../service/README.md](../service/README.md].
