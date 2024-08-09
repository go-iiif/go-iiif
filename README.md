# go-iiif

![spanking cat](misc/go-iiif-spanking-cat.png)

## What is this?

This began as a fork of [@greut's iiif](https://github.com/greut/iiif) package that moves all of the processing logic for the [IIIF Image API](http://iiif.io/api/image/) in to discrete Go packages and defines source, derivative and graphics details in a [JSON config file](README.md#config-files). There is an additional caching layer for both source images and derivatives.

I did this to better understand the architecture behind (and to address my own concerns about) the [IIIF Image API](http://iiif.io/api/image/2.1/index.html).

For the time being this package will probably not support the other IIIF Metadata or Publication APIs.

_And by "forked" I mean that [@greut](https://github.com/greut) and I decided that [it was best](https://github.com/greut/iiif/pull/2) for this code and his code to wave at each other across the divide but not necessarily to hold hands._

## Releases

The current release is `github.com/go-iiif/go-iiif/v6`.

### v7

Version 7.0.0 is in the early stages of development and it is expected that it will introduce a number of backwards incompatible changes to how the `tools` package is structured and the interface and method signatures it exposes. 

Additionally, many properties in `config.Config` blocks will be removed in favour of a single URI-style declarative syntax.

The goal is to leave the command-line tools unchanged with the exception of removing any flags that have previously been flagged as deprecated.

There is no expected release date for the `v7` branch yet.

### v6

Version 6.0.0 was updated to use the [aaronland/go-flickr-api](https://github.com/aaronland/go-flickr-api) package to retrieve Flickr photos which introduced backwards incompatible changes in the `config.FlickrConfig` block.

### v5

Version 5.0.0 and higher of the `go-iiif` package introduces three backwards incompatible changes from previous versions. They are:

* The `tile/seed.go` package and `cmd/iiif-tile-seed` tool assume IIIF Level 0 profiles rather than Level 2 to account for [issue #92](https://github.com/go-iiif/go-iiif/issues/92).

* The `profile` package and types have been removed. The code to generate `info.json` files has been moved in to the `info` package.

* The interface for the `level` package has been changed. Specifically the `Profile` method has been changed to return a URI string and there is a new `Endpoint` method.

### v2

Version 2.0.0 and higher of the `go-iiif` package introduces three backwards incompatible changes from previous versions. They are:

* The removal of the `libvips` and `bimg` package for default image processing and the introduction of "drivers" for defining image processing functionality.
* The use of the [Go Cloud](https://gocloud.dev/) `Bucket` and `Blob` interfaces for reading and writing files.
* The introduction of [go-iiif-uri](https://github.com/go-iiif/go-iiif-uri) URI strings rather than paths or filenames to define images for processing.

All three changes are discussed in detail below.

## Drivers

`go-iiif` was first written with the [libvips](https://github.com/jcupitt/libvips) library and [bimg](https://github.com/h2non/bimg/) Go wrapper for image processing. `libvips` is pretty great but it introduces non-trivial build and setup requirements. As of version 2.0 `go-iiif` no longer uses `libvips` by default but instead does all its image processing using native (Go) code. This allows `go-iiif` to run on any platform supported by Go without the need for external dependencies.

A longer discussion about drivers and how they work follows but if you want or need to use `libvips` for image processing you should use the [go-iiif-vips](https://github.com/go-iiif/go-iiif-vips) package.

Support for alternative image processing libraries, like `libvips` is supported through the use of "drivers" (similar to the way the Go `database/sql` package works). A driver needs to support the `driver.Driver` interface which looks like this:

```
import (
	iiifcache "github.com/go-iiif/go-iiif/v6/cache"
	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
	iiifsource "github.com/go-iiif/go-iiif/v6/source"
)

type Driver interface {
	NewImageFromConfigWithSource(*iiifconfig.Config, iiifsource.Source, string) (iiifimage.Image, error)
	NewImageFromConfigWithCache(*iiifconfig.Config, iiifcache.Cache, string) (iiifimage.Image, error)
	NewImageFromConfig(*iiifconfig.Config, string) (iiifimage.Image, error)
}
```

The idea here is that the bulk of the `go-iiif` code isn't aware of who or how images are _actually_ being processed only that it can reliably pass around things that implement the `image.Image` interface (the `go-iiif` image interface, not the Go language interface).

Drivers are expected to "register" themselves through the `driver.RegisterDriver` method at runtime. For example:

```
package native

import (
	iiifdriver "github.com/go-iiif/go-iiif/v6/driver"
)

func init() {

	dr, err := NewNativeDriver()

	if err != nil {
		panic(err)
	}

	iiifdriver.RegisterDriver("native", dr)
}
```

And then in your code you might do something like this:

```
import (
	"context"
	"github.com/aaronland/gocloud-blob-bucket"	
	_ "github.com/go-iiif/go-iiif/v6/native"
	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
	iiifdriver "github.com/go-iiif/go-iiif/v6/driver"	
)

ctx := context.Background()
	
config_bucket, _ := bucket.OpenBucket(ctx, "file:///etc/go-iiif")

cfg, _ := config.NewConfigFromBucket(ctx, config_bucket, "config.json")

driver, _ := iiifdriver.NewDriverFromConfig(cfg)
```

That's really the only change to existing code. Careful readers may note the calls to `bucket.OpenBucket` and `config.NewConfigFromBucket` to load `go-iiif` configuration files. This is discussed below. In the meantime the only other change is to update the previously default `graphics.source` property in the configuration file from `VIPS` (or `vips`) to `native`. For example:

```
    "graphics": {
	"source": { "name": "VIPS" }
    }
```

Becomes:

```
    "graphics": {
	"source": { "name": "native" }
    }
```

The value of the `graphics.source` property should match the name that driver uses to register itself with `go-iiif`.

The rest of the code in `go-iiif` has been updated to expect a `driver.Driver` object and to invoke the relevant `NewImageFrom...` method as needed. It is assumed that the driver package in question will also implement it's own implementation of the `go-iiif` `image.Image` interface. For working examples you should consult either of the following packages:

* https://github.com/go-iiif/go-iiif/v6/tree/master/native
* https://github.com/go-iiif/go-iiif-vips

## Buckets

Starting with version 2 the `go-iiif` package uses the [Go Cloud](https://gocloud.dev/) `Bucket` and `Blob` interfaces for reading and writing all files. For example, instead of doing this:

```
cfg, _ := config.NewConfigFromFile("/etc/go-iiif/config.json")
```

It is now necessary to do this:

```
config_bucket, _ := bucket.OpenBucket(ctx, "file:///etc/go-iiif")
cfg, _ := config.NewConfigFromBucket(ctx, config_bucket, "config.json")
```
This allows for configuration files, and others, to be stored and retrieved from [any "bucket" source that is supported by the Go Cloud package](https://gocloud.dev/howto/blob/#services), notably remote storage services like AWS S3.

The `source` and `caching` layers have also been updated accordingly but support for the older `Disk`, `S3` and `Memory` sources has been updated to use the `Go Cloud` packages so there is no need to update any existing `go-iiif` configuration files.

## URIs

[go-iiif-uri](https://github.com/go-iiif/go-iiif-uri) URI strings are still a work in progress. While they may still change a bit around the edges efforts will be made to ensure backwards compatibility going forward.

`go-iiif-uri` URI strings are defined by a named scheme which indicates how an URI should be processed, a path which is a reference to an image and zero or more query parameters which are the specific instructions for processing the URI.

### file

```
file:///path/to/source/image.jpg
```

```
file:///path/to/source/image.jpg?target=/path/to/target/image.jpg
```

The `file://` URI scheme is basically just a path or filename. It has an option `target` property which allows the name of the source image to be changed. These filenames are _not_ the final name of the image as processed by `go-iiif` but the name of the directory structure that files will be written to, as in the weird IIIF instructions-based URIs. 

Valid parameters for the `file://` URI scheme are:

| Name | Type | Required |
| --- | --- | --- |
| target | string | no |

### idsecret

```
idsecret:///path/to/source/image.jpg?id=1234&secret=s33kret&secret_o=seekr3t&label
```

The `idsecret://` URI scheme is designed to rewrite a source image URI to {UNIQUE_ID} + {SECRET} + {LABEL} style filenames. For example `cat.jpg` becomes `1234_s33kret_b.jpg` and specifically `123/4/1234_s33kret_b.jpg` where the unique ID is used to generate a nested directory tree in which the final image lives.

The `idsecret://` URI scheme was developed for use with `go-iiif` "instructions" files where a single image produced multiple derivatives that need to share commonalities in their final URIs.

Valid parameters for the `idsecret://` URI scheme are:

| Name | Type | Required |
| --- | --- | --- |
| id | int64 | yes |
| label | string | yes |
| format | string | yes |
| original | string | no |
| secret | string | no |
| secret_o | string | no |

If either the `secret` or `secret_o` parameters are absent they will be auto-generated.

### rewrite

```
rewrite:///path/to/source/image.jpg?target=/path/to/target/picture.jpg
```

The `rewrite://` URI scheme is a variant of the `file://` URI scheme except that the `target` query parameter is required and it will be used to redefine the final URI, rather than just its directory tree, of the processed image.

| Name | Type | Required |
| --- | --- | --- |
| target | string | yes |

## Example

Here's a excerpted example taken from the [process/parallel.go](process/parallel.go) package that processes a single source image, defined as an `idsecret://` URI, in to multiple derivatives defined in an "instructions" file.

The `idsecret://` URI is output as a string using the instructions set to define the `label` and other query parameters. That string is then used to create a new `rewrite://` URI where source is derived from the original `idsecret://` URI and the target is newly generate URI string.

```
go func(u iiifuri.URI, label Label, i IIIFInstructions) {

	var process_uri iiifuri.URI

	switch u.Driver() {
	case "idsecret":

		str_label := fmt.Sprintf("%s", label)

		opts := &url.Values{}
		opts.Set("label", str_label)
		opts.Set("format", i.Format)

		if str_label == "o" {
			opts.Set("original", "1")
		}

		target_str, _ := u.Target(opts)

		origin := u.Origin()

		rw_str := fmt.Sprintf("%s?target=%s", origin, target_str)
		rw_str = iiifuri.NewRewriteURIString(rw_str)

		rw_uri, err := iiifuri.NewURI(rw_str)

		process_uri = rw_uri

	default:
		process_uri = u
	}

	new_uri, im, _ := pr.ProcessURIWithInstructions(process_uri, label, i)
	// do something with new_uri and im here...
	
}(u, label, i)
```

## Command line tools

`go-iiif` was designed to expose all of its functionality outside of the included tools although that hasn't been documented yet. The source code for the [iiif-tile-seed](cmd/iiif-tile-seed.go), [iiif-transform](cmd/iiif-transform.go) and [iiif-process](cmd/iiif-process.go) tools is a good place to start poking around if you're curious.

For the sake of shortening this file documentation for the command line tools has been moved in to [cmd/README.md](cmd/README.md).

## Config files

There is a [sample config file](docs/config.json.example) included with this repo. The easiest way to understand config files is that they consist of at least five top-level groupings, with nested section-specific details, followed by zero or more implementation specific configuration blocks. The five core blocks are:

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
		"source": { "name": "native" }
	}
```

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

Here's how that dynamic plays out in reality. The table below lists all the IIIF parameters and their associate features. Each feature lists its syntax and whether or not it is required and supported [according to the official spec](compliance/level2.go) but then also according to the [example `go-iiif` config file](docs/config.json.example), included with this repo.

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
		"source": { "name": "Disk", "path": "example/images" },
		"cache": { "name": "Memory", "ttl": 300, "limit": 100 }
	}
```

Details about source images.

#### images.source

Where to find source images.

##### Blob

```
	"images": {
		"source": { "name": "Blob", "path": "file:///example/images" }
	}
```

Fetch sources images from any supported [Go Cloud storage service](https://gocloud.dev/howto/blob/#services).

Some notes about the `Blob` source:

* Under the hood the [github.com/aaronland/go-cloud-s3blob](https://github.com/aaronland/go-cloud-s3blob) package is used to open Go Cloud "buckets", which are the parent containers for Go Cloud "blobs". This is specifically in order to be able to specify AWS S3 credentials using string values: `env:` (read credentials from AWS defined environment variables), `iam:` (assume AWS IAM credentials), `{AWS_PROFILE_NAME}`, `{AWS_CREDENTIALS_PATH}:{AWS_PROFILE_NAME}`. For example:

```
	"images": {
		"source": { "name": "Blob", "path": "s3:///bucket-name?region=us-east-1&credentials=iam:" }
	}
```

##### Disk

```
	"images": {
		"source": { "name": "Disk", "path": "example/images" }
	}
```

Fetch source images from a locally available filesystem.

_The `Disk` source is still supported but has been replaced by the `Blob` source._

##### Flickr

```
	"images": {
		"source": { "name": "Flickr" },
		"cache": { "name": "Memory", "ttl": 60, "limit": 100 }
	},
	"flickr": {
		"client_uri": "oauth1://?consumer_key={KEY}&consumer_secret={SECRET}",
	}
```

Fetch source images from Flickr. You will need to provide a valid [Flickr API key](https://www.flickr.com/services/api/). A few caveats:

* Under the hood the code is using the [aaronland/go-flickr-api](https://github.com/aaronland/go-flickr-api) package which uses a URI-style syntax for defining client instances. Please consult [the `go-flickr-api` documentation](https://github.com/aaronland/go-flickr-api?tab=readme-ov-file#oauth1) for details on how to construct those URIs.
* The code calls the [flickr.photos.getSizes](https://www.flickr.com/services/api/flickr.photos.getSizes.html) API method and looks for the first of the following photo sizes in this order: `Original, Large 2048, Large 1600, Large`. If none are available then an error is triggered.
* Photo size lookups are not cached yet.

Here's an example [with this photo](https://www.flickr.com/photos/straup/4136870023/in/album-72157622883263698/):

![](misc/go-iiif-flickr.png)

![](misc/go-iiif-flickr-detail.png)

##### S3

```
	"images": {
		"source": { "name": "S3", "path": "your.S3.bucket", "region": "us-east-1", "credentials": "default" }
	}
```

Fetch source images from Amazon's S3 service. S3 caches assume that that the `path` key is the name of the S3 bucket you are reading from. S3 caches have three addition properties:

* **prefix** is an optional path to a sub-path inside of your S3 bucket where images are stored.
* **region** is the name of the AWS region where your S3 bucket lives. Sorry [this is an AWS-ism](https://docs.aws.amazon.com/sdk-for-go/v1/developerguide/configuring-sdk.html)
* **credentials** is a string describing [how your AWS credentials are defined](https://docs.aws.amazon.com/general/latest/gr/aws-access-keys-best-practices.html). Valid options are:
 * `env:` - Signals that you you have defined valid AWS credentials as environment variables
 * `shared:PATH_TO_SHARED_CREDENTIALS_FILE:SHARED_CREDENTIALS_PROFILE` - Signals that your AWS credentials are in a [shared credentials files](https://docs.aws.amazon.com/general/latest/gr/aws-access-keys-best-practices.html) and that `go-iiif` should use a specific profile - _this syntax is deprecated and you should just use use `PATH_TO_SHARED_CREDENTIALS_FILE:SHARED_CREDENTIALS_PROFILE` instead._
 * `iam:` - Signals that you are using `go-iiif` in an AWS environment with suitable roles and permissioning for working with S3. The details of how and where you configure IAM roles are outside the scope of this document.
 * `PATH_TO_SHARED_CREDENTIALS_FILE:SHARED_CREDENTIALS_PROFILE` - Signals that your AWS credentials are in a [shared credentials files](https://docs.aws.amazon.com/general/latest/gr/aws-access-keys-best-practices.html) and that `go-iiif` should use a specific profile
 * `SHARED_CREDENTIALS_PROFILE` - Signals that your AWS credentials are in default location(s) for the current user and that `go-iiif` should use a specific profile
 
For the sake of backwards compatibilty if the value of `credentials` is any other string then it will be assumed to be the name of the profile you wish to use for a valid [credential files](https://docs.aws.amazon.com/sdk-for-go/v1/developerguide/configuring-sdk.html) in the home directory of the current user. Likewise if the value of `credentials` is an _empty_ string (or absent) it will be assumed that valid AWS access credentials have been defined as environment variables.

It is not possible to define your AWS credentials as properties in your `go-iiif` config file.

Important: If you are both reading source files and writing cached derivatives to S3 in the same bucket make sure they have **different** prefixes. If you don't then AWS will happily overwrite your original source files with the directory (which shares the same names as the original file) containing your derivatives. Good times.

_The `S3` source is still supported but has been replaced by the `Blob` source._

![](misc/go-iiif-aws-source-cache.png)

##### URI

```
	"images": {
		"source": { "name": "URI", "path": "https://images.collection.cooperhewitt.org/{id}" }
	}
```

Fetch source images from a remote URI. The `path` parameter must be a valid (Level 4) [URI Template](http://tools.ietf.org/html/rfc6570) with an `{id}` placeholder.

#### images.cache

Caching options for source images.

##### Blob

```
	"images": {
		"cache": { "name": "Blob", "path": "file:///example/images" }
	}
```

Cache sources images to any supported [Go Cloud storage service](https://gocloud.dev/howto/blob/#services).

Some notes about the `Blob` cache:

* Under the hood the [github.com/aaronland/go-cloud-s3blob](https://github.com/aaronland/go-cloud-s3blob) package is used to open Go Cloud "buckets", which are the parent containers for Go Cloud "blobs". This is specifically in order to be able to specify AWS S3 credentials using string values: `env:` (read credentials from AWS defined environment variables), `iam:` (assume AWS IAM credentials), `{AWS_PROFILE_NAME}`, `{AWS_CREDENTIALS_PATH}:{AWS_PROFILE_NAME}`. For example:

```
	"images": {
		"cache": { "name": "Blob", "path": "s3:///bucket-name?region=us-east-1&credentials=iam:" }
	}
```

* Under the hood the `Blob` cache supports an optional `acl={ACL}` query parameter in the path property (which is equivalent to a Go Cloud URI definition). This is to account for [the inability to assign permissions](https://github.com/google/go-cloud/issues/1108) when writing Go Cloud blob objects. Currently the `acl=ACL` parameter is only honoured for `s3://` URIs but [patches for other sources would be welcome](https://github.com/go-iiif/go-iiif/blob/master/cache/blob.go). Additionally it is only possible to assign ACLs for a Go Cloud "bucket" URI and not invidiual "blobs". For example:

```
	"images": {
		"cache": { "name": "Blob", "path": "s3:///bucket-name?region=us-east-1&credentials=iam:&acl=public-read" }
	}
```

##### Disk

```
	"images": {
		"cache": { "name": "Disk", "path": "example/cache" }
	}
```

Cache images to a locally available filesystem.

_The `Disk` cache is still supported but has been replaced by the `Blob` cache._

##### Memory

```
	"images": {
		"cache": { "name": "Memory", "ttl": 300, "limit": 100 }
	}
```

Cache images in memory. Memory caches have two addition properties:

* **ttl** is the maximum number of seconds an image should live in cache.
* **limit** the maximum number of megabytes the cache should hold at any one time.

##### Null

```
	"images": {
		"cache": { "name": "Null" }
	}
```

Because you must define a caching layer this is here to satify the requirements without actually caching anything, anywhere.

### derivatives

```
	"derivatives": {
		"cache": { "name": "Disk", "path": "example/cache" }
	}
```

Details about derivative images.

#### derivatives.cache

Caching options for derivative images.

##### Blob

```
	"derivatives": {
		"cache": { "name": "Blob", "path": "file:///example/images" }
	}
```

Cache derivation images to any supported [Go Cloud storage service](https://gocloud.dev/howto/blob/#services).

Some notes about the `Blob` cache:

* Under the hood the [github.com/aaronland/go-cloud-s3blob](https://github.com/aaronland/go-cloud-s3blob) package is used to open Go Cloud "buckets", which are the parent containers for Go Cloud "blobs". This is specifically in order to be able to specify AWS S3 credentials using string values: `env:` (read credentials from AWS defined environment variables), `iam:` (assume AWS IAM credentials), `{AWS_PROFILE_NAME}`, `{AWS_CREDENTIALS_PATH}:{AWS_PROFILE_NAME}`. For example:

```
	"derivatives": {
		"cache": { "name": "Blob", "path": "s3:///bucket-name?region=us-east-1&credentials=iam:" }
	}
```

* Under the hood the `Blob` cache supports an optional `acl={ACL}` query parameter in the path property (which is equivalent to a Go Cloud URI definition). This is to account for [the inability to assign permissions](https://github.com/google/go-cloud/issues/1108) when writing Go Cloud blob objects. Currently the `acl=ACL` parameter is only honoured for `s3://` URIs but [patches for other sources would be welcome](https://github.com/go-iiif/go-iiif/blob/master/cache/blob.go). Additionally it is only possible to assign ACLs for a Go Cloud "bucket" URI and not invidiual "blobs". For example:

```
	"derivatives": {
		"cache": { "name": "Blob", "path": "s3:///bucket-name?region=us-east-1&credentials=iam:&acl=public-read" }
	}
```

##### Disk

```
	"derivatives": {
		"cache": { "name": "Disk", "path": "example/cache" }
	}
```

Cache images to a locally available filesystem.

_The `Disk` cache is still supported but has been replaced by the `Blob` cache._

##### Memory

```
	"derivatives": {
		"cache": { "name": "Memory", "ttl": 300, "limit": 100 }
	}
```

Cache images in memory. Memory caches have two addition properties:

* **ttl** is the maximum number of seconds an image should live in cache.
* **limit** the maximum number of megabytes the cache should hold at any one time.

##### Null

```
	"derivatives": {
		"cache": { "name": "Null" }
	}
```

Because you must define a caching layer this is here to satify the requirements without actually caching anything, anywhere.

##### S3

```
	"derivatives": {
		"cache": { "name": "S3", "path": "your.S3.bucket", "region": "us-east-1", "credentials": "default" }
	}
```

Cache images using Amazon's S3 service. S3 caches assume that that the `path` key is the name of the S3 bucket you are reading from. S3 caches have three addition properties:

* **prefix** is an optional path to a sub-path inside of your S3 bucket where images are stored.
* **region** is the name of the AWS region where your S3 bucket lives. Sorry [this is an AWS-ism](https://docs.aws.amazon.com/sdk-for-go/v1/developerguide/configuring-sdk.html)
* **credentials** is a string describing [how your AWS credentials are defined](https://docs.aws.amazon.com/general/latest/gr/aws-access-keys-best-practices.html). Valid options are:
 * `env:` - Signals that you you have defined valid AWS credentials as environment variables
 * `shared:PATH_TO_SHARED_CREDENTIALS_FILE:SHARED_CREDENTIALS_PROFILE` - Signals that your AWS credentials are in a [shared credentials files](https://docs.aws.amazon.com/general/latest/gr/aws-access-keys-best-practices.html) and that `go-iiif` should use a specific profile - _this syntax is deprecated and you should just use use `PATH_TO_SHARED_CREDENTIALS_FILE:SHARED_CREDENTIALS_PROFILE` instead._
 * `iam:` - Signals that you are using `go-iiif` in an AWS environment with suitable roles and permissioning for working with S3. The details of how and where you configure IAM roles are outside the scope of this document.
 * `PATH_TO_SHARED_CREDENTIALS_FILE:SHARED_CREDENTIALS_PROFILE` - Signals that your AWS credentials are in a [shared credentials files](https://docs.aws.amazon.com/general/latest/gr/aws-access-keys-best-practices.html) and that `go-iiif` should use a specific profile
 * `SHARED_CREDENTIALS_PROFILE` - Signals that your AWS credentials are in default location(s) for the current user and that `go-iiif` should use a specific profile

For the sake of backwards compatibilty if the value of `credentials` is any other string then it will be assumed to be the name of the profile you wish to use for a valid [credential files](https://docs.aws.amazon.com/sdk-for-go/v1/developerguide/configuring-sdk.html) in the home directory of the current user. Likewise if the value of `credentials` is an _empty_ string (or absent) it will be assumed that valid AWS access credentials have been defined as environment variables.

It is not possible to define your AWS credentials as properties in your `go-iiif` config file.

Important: If you are both reading source files and writing cached derivatives to S3 in the same bucket make sure they have **different** prefixes. If you don't then AWS will happily overwrite your original source files with the directory (which shares the same names as the original file) containing your derivatives. Good times.

_The `S3` cache is still supported but has been replaced by the `Blob` cache._

![](misc/go-iiif-aws-source-cache.png)

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

`dither` will create a black and white [halftone](https://en.wikipedia.org/wiki/Halftone) derivative of an image using the [Atkinson dithering algorithm](https://en.wikipedia.org/wiki/Dither#Algorithms). Dithering is enabled in the [example config file](docs/config.json.example) and you can invoke it like this:

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

![](misc/go-iiif-primitive-animated-rect.gif)

Here are examples where each of the tiles in an slippy image are animated GIFs:

* https://go-iiif.github.io/go-iiif/animated/
* https://go-iiif.github.io/go-iiif/animated/?mode=circles
* https://go-iiif.github.io/go-iiif/animated/?mode=triangles

_Note: You will need to [manually enable support for GIF images](https://github.com/go-iiif/go-iiif/tree/primitive#featuresenable) in your config file for animated GIFs to work._

## Non-standard services

### palette

`go-iiif` supports the following additional services for profiles:

* `blurhash` for generateing a compact base-83 encoded representation of an image.
* `imagehash` for generating average and difference perceptual hashes of an image.
* `palette` for extracting a colour palette for an image.

Details for configuring these service are discussed [above](#services) but here is the output for a service with the default settings:

```
$> curl -s localhost:8080/spanking.jpg/info.json | jq '.service'

[
  {
    "@context": "x-urn:service:go-iiif#palette",
    "profile": "x-urn:service:go-iiif#palette",
    "label": "x-urn:service:go-iiif#palette",
    "palette": [
      {
        "name": "#4e3c24",
        "hex": "#4e3c24",
        "reference": "vibrant"
      },
      {
        "name": "#9d8959",
        "hex": "#9d8959",
        "reference": "vibrant"
      },
      {
        "name": "#c7bca6",
        "hex": "#c7bca6",
        "reference": "vibrant"
      },
      {
        "name": "#5a4b36",
        "hex": "#5a4b36",
        "reference": "vibrant"
      }
    ]
  },
  {
    "@context": "x-urn:service:go-iiif#blurhash",
    "profile": "x-urn:service:go-iiif#blurhash",
    "label": "x-urn:service:go-iiif#blurhash",
    "hash": "LOOWsZxu_4-;~pj[Rjof-;kBIAWB"
  },
  {
    "@context": "x-urn:service:go-iiif#imagehash",
    "profile": "x-urn:service:go-iiif#imagehash",
    "label": "x-urn:service:go-iiif#imagehash",
    "average": "a:ffffc7e7c3c3c3c3",
    "difference": "d:c48c0c0e8e8f0e0f"
  }
]
```

_Please remember that `go-colours` itself is a work in progress so you should approach the `palette` service accordingly._

### Writing your own non-standard services

Services are invoked by the `go-iiif` codebase using URI-style identifiers. For example, assuming an "example" service you would invoke it like this:

```
    	service_name := "example"	
	service_uri := fmt.Sprintf("%s://", service_name)
	service, _ := iiifservice.NewService(ctx, service_uri, cfg, im)
```

In addition to implementing the `service.Service` interface custom services need to also "register" themselves on initialization with a (golang) context, a (go-iiif), a unique scheme used to identify the service and a `service.ServiceInitializationFunc` callback function. The callback function implements the following interface:

```
type ServiceInitializationFunc func(ctx context.Context, config *iiifconfig.Config, im iiifimage.Image) (Service, error)
```

Here is an abbreviated example, with error handling removed for the sake of brevity. For real working examples, take a look at any of the built-in services in the [services](services) directory.

```
package example	// for example "github.com/example/go-iiif-example"

import (
	"context"
	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
	iiifimage "github.com/go-iiif/go-iiif/v6/image"	
	iiifservice "github.com/go-iiif/go-iiif/v6/service"	
)

func init() {
	ctx := context.Background()
	iiifservice.RegisterService(ctx, "example", initExampleService)
}

func initExampleService(ctx context.Context, cfg *iiifconfig.Config, im iiifimage.Image) (iiifservice.Service, error) {
	return NewExampleService(cfg, im)
}

type ExampleService struct {
	iiifservice.Service        `json:",omitempty"`
	// your properties here...
}

// your implementation of the iiifservice.Service interface here...

func NewExampleService(cfg *iiifconfig.Config, im iiifimage.Image) (iiifservice.Service, error){

     // presumably you'd do something with im here...
     
     s := &ExampleService{
       // your properties here...
     }

     return s, nil
}
```

Finally, you will need to create custom versions of any `go-iiif` tools you want to you use your new service. For example, here's a modified version of the [cmd/iiif-server/main.go](cmd/iiif-server/main.go) server implementation.

```
package main

import (
       _ "github.com/example/go-iiif-example"
)

import (
	"context"
	_ "github.com/aaronland/go-cloud-s3blob"
	_ "github.com/go-iiif/go-iiif/v6/native"
	"github.com/go-iiif/go-iiif/v6/tools"
	_ "gocloud.dev/blob/fileblob"
	"log"
)

func main() {

	tool, err := tools.NewIIIFServerTool()

	if err != nil {
		log.Fatal(err)
	}

	err = tool.Run(context.Background())

	if err != nil {
		log.Fatal(err)
	}
}
```

 The only change from the default server tool is the addition of the `_ "github.com/example/go-iiif-example"` import statement. That will allow the core `go-iiif` software to find and use your custom service.

It's unfortunate that using custom and bespoke services requires compiling your own version of the `go-iiif` tools but such is life when you are using a language like Go.

## Example

There is a live demo of the [Leaflet-IIIF](https://github.com/mejackreed/Leaflet-IIIF) slippymap provider used in conjunction with a series of tiles images generated using the `iiif-tile-seed` utility available for viewing over here:

https://go-iiif.github.io/go-iiif/

The `iiif-server` tool also comes with a canned example (consisting of exactly one image) so you can see things in the context of a slippy map. Here's what you need to do to get it set up:

First, make sure have a valid `go-iiif` config file. If you don't then you can copy the example config included in this repo:

```
$> cp docs/config.json.example config.json
```

Next, pre-seed some tiles for an image. You don't necessarily need to do this step but it's included to show you how it's done:

```
$> ./bin/iiif-tile-seed -config config.json -endpoint http://localhost:8082 -scale-factors 8,4,2,1 184512_5f7f47e5b3c66207_x.jpg
```

_Note how we are specifying the endpoint where these tiles will be served from. That's necessary so that we can also pre-seed a [profile description](http://iiif.io/api/image/2.1/#profile-description) for each image as well as tiles._

Finally start up the `iiff-server` and be sure to pass the `-example` flag:

```
$> ./bin/iiif-server -config config.json -host localhost -port 8082 -example
```

Now if you visit `http://localhost:8082/example/` in your browser you should see this:

![spanking cat](misc/go-iiif-example.png)

Assuming you've pre-seed your tiles if you open up the network console in your browser then you should see something like this, namely that the individual tiles are returned speedy and fast:

![spanking cat](misc/go-iiif-example-cached.png)

### Generating static images

![spanking cat](misc/go-iiif-example-screenshot.png)

The example included with `go-iiif` has an added super power which is the ability to create a static image of the current state of the map/image.

Just click the handy `📷` button to the bottom right of the image and you will be prompted for where you'd like to save your new image.

_This is not a feature of `go-iiif` itself. It's [entirely client-side magic in your browser](example/index.js) but it's still pretty cool..._

## Performance and load testing

### iiif-tile-seed

Processing individual or small batches of images `go-iiif` ranges from pretty fast to very fast. For example here is a picture of Spanking Cat width a [maximum dimension of 4096 pixels](https://images.collection.cooperhewitt.org/184512_5f7f47e5b3c66207_x.jpg):

```
$> ./bin/iiif-tile-seed -config config.json -refresh -scale-factors 8,4,2,1 184512_5f7f47e5b3c66207_x.jpg
[184512_5f7f47e5b3c66207_x.jpg] time to process 340 tiles: 27.537429902s
```

So any individual tile is pretty speedy but in the aggregate it starts to add up. I will need to do some continued digging to make sure that the source image isn't being processed unnecessarily for each tile. Here is the same image but with a [maximum dimension of 2048 pixels](https://images.collection.cooperhewitt.org/184512_b812003c86c3525b_k.jpg):

```
$> ./bin/iiif-tile-seed -config config.json -refresh -scale-factors 4,2,1 184512_b812003c86c3525b_k.jpg
[184512_b812003c86c3525b_k.jpg] time to process 84 tiles: 1.894074539s
```

Note that we are only generating tiles for three scale factors instead of four. But that's not where things slow down as we can see seeding tiles for only three scale factors for the larger image:

```
$> ./bin/iiif-tile-seed -config config.json -refresh -scale-factors 4,2,1 184512_5f7f47e5b3c66207_x.jpg
[184512_5f7f47e5b3c66207_x.jpg] time to process 336 tiles: 26.925253066s
```

For processing large, or large volumes of, images the bottlenecks will be:

* CPU usage crunching pixels
* Disk I/O writing tiles to disk
* Running out of inodes

That said on a machine with 8 CPUs and 32GB RAM I was able to run the machine hot with all the CPUs pegged at 100% usage and seed 100, 000 (2048x pixel) images yielding a little over 3 million, or approximately 70GB of, tiles in 24 hours. Some meaningful but not overwhelming amount of time was spent fetching source images across the network so presumably things would be faster reading from a local filesystem.

Memory usage across all the `iiif-tile-seed` processes never went above 5GB and, in the end, I ran out of inodes.

The current strategy for seeding tiles may also be directly responsible for some of the bottlenecks. Specifically, when processing large volumes of images (defined in a CSV file) the `ifff-tile-seed` will spawn and queue as many concurrent Go routines as there are CPUs. For each of those processes then another (n) CPUs * 2 subprocesses will be spawned to generate tiles. Maybe this is just too image concurrent image processing routines to have? I mean it works but still... Or maybe it's just that every one is waiting for bytes to be written to disk. Or all of the above. I'm not sure yet.

### iiif-server

All of the notes so far have assumed that you are using `iiif-tile-seed`. If you are running `iiif-server` the principle concern will be getting overwhelmed by too many requests for too many different images, especially if they are large, and running out of memory. That is why you can define an [in-memory cache](https://github.com/go-iiif/go-iiif/blob/master/README.md#memory) for source images but that will only be of limited use if your problem is handling concurrent requests. It is probably worth adding checks and throttles around current memory usage to the various handlers...

## Docker

Yes. There is a [Dockerfile](Dockerfile) included with this distribution. It will build a container with the following tools:

* The `iiif-server` tool.
* The `iiif-process` command-line tool.
* The `iiif-tile-seed` command-line tool.

To build the container run:

```
$> docker build -f Dockerfile -t go-iiif .
```

To start the `iiif-server` tool run:

```
$> docker run -it -p 6161:8080 \
   -v /usr/local/go-iiif/docker/etc:/etc/iiif-server \
   -v /usr/local/go-iiif/docker/images:/usr/local/iiif-server \
   go-iiif \
   /bin/iiif-server -host 0.0.0.0 \
   -config-source file:///etc/iiif-server
   
2018/06/20 23:03:10 Listening for requests at 0.0.0.0:8080
```

See the way we are mapping `/etc/iiif-server` and `/usr/local/iiif-server` to local directories? By default the `iiif-server` Dockerfile does not bundle config files or images. Maybe some day, but that day is not today.

Then, in another terminal:

```
$> curl localhost:6161/test.jpg/info.json
{"@context":"http://iiif.io/api/image/2/context.json","@id":"http://localhost:6161/test.jpg","@type":"iiif:Image","protocol":"http://iiif.io/api/image","width":3897,"height":4096,"profile":["http://iiif.io/api/image/2/level2.json",{"formats":["gif","webp","jpg","png","tif"],"qualities":["default","color","dither"],"supports":["full","regionByPx","regionByPct","regionSquare","sizeByDistortedWh","sizeByWh","full","max","sizeByW","sizeByH","sizeByPct","sizeByConfinedWh","none","rotationBy90s","mirroring","noAutoRotate","baseUriRedirect","cors","jsonldMediaType"]}],"service":[{"@context":"x-urn:service:go-iiif#palette","profile":"x-urn:service:go-iiif#palette","label":"x-urn:service:go-iiif#palette","palette":[{"name":"#2f2013","hex":"#2f2013","reference":"vibrant"},{"name":"#9e8e65","hex":"#9e8e65","reference":"vibrant"},{"name":"#c6bca6","hex":"#c6bca6","reference":"vibrant"},{"name":"#5f4d32","hex":"#5f4d32","reference":"vibrant"}]}]}
```

Let's say you're using S3 as an image source and reading (S3) credentials from environment variables (something like `{"source": { "name": "S3", "path": "{BUCKET}", "region": "us-east-1", "credentials": "env:" }`) then you would start up `iiif-server` like this:

```
$> docker run -it -p 6161:8080 \
       -v /usr/local/go-iiif/docker/etc:/etc/iiif-server \
       -v /usr/local/go-iiif/docker/images:/usr/local/iiif-server \
       -e AWS_ACCESS_KEY_ID={AWS_KEY} -e AWS_SECRET_ACCESS_KEY={AWS_SECRET} \
       go-iiif \
       /bin/iiif-server -host 0.0.0.0 \
       -config-source file:///etc/iiif-server
```

The process an image using the `iiif-process` tool you would run something like:

```
$> docker run \
   -v /usr/local/go-iiif/docker/etc:/etc/go-iiif \
   go-iiif /bin/iiif-process \
   -config-source file:///etc/go-iiif \
   -instructions-source file:///etc/go-iiif \
   file:///test.jpg
```

To tile an image using the `iiif-tile-seed` tool you would run something like:

```
$> docker run -v /usr/local/go-iiif-vips/docker:/usr/local/go-iiif \
	go-iiif /bin/iiif-tile-seed \
	-config-source file:////usr/local/go-iiif/config \
	-scale-factors 1,2,4,8 \
	file:///zuber.jpg
```

Again, see the way we're mapping `/etc/go-iiif` to a local folder, like we do in the `iiif-server` Docker example? The same rules apply here.

### Amazon ECS

I still find ECS to be a world of [poorly-to-weirdly documented](https://aws.amazon.com/getting-started/tutorials/deploy-docker-containers/) strangeness .Remy Dewolf's
[AWS Fargate: First hands-on experience and review](https://medium.com/@remy.dewolf/aws-fargate-first-hands-on-experience-and-review-1b52fca2148e) is a pretty good introduction.

I have gotten IIIF-related things to work in ECS but it's always a bit nerve-wracking and I haven't completely internalized the steps in order to repeat them to someone else. What follows should be considered a "current best-attempt".

#### iiif-server

What follows are non-comprehensive notes for getting `iiif-server` to work under
ECS. The bad news is that it's fiddly (and weird, did I mention that?) The good
news is that I _did_ get it to work.

These are not detailed instructions for setting up `iiif-server` in ECS from
scratch. You should consult the [Amazon Elastic Container Service Documentation](https://aws.amazon.com/documentation/ecs/)
for that.

What follows assumes that you're using an S3 "source" for source images and
derivatives. I have not tried any of this with [EBS volumes mounted as Docker
volumes](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/using_data_volumes.html)
so if you have I'd love to hear about it.

##### Services

You will need to ensure that the service has `Auto-assign public IP(s)`
enabled.This is necessary in order to fetch the actual Docker container.

The corollary to that is that unless you are _wanting_ to expose your instances
of `iiif-server` to the public internet you will need to add a security group
(to your ECS service) with suitable restrictions.

##### Task definitions

##### Command 

`/bin/iiif-server`

###### Port mappings

`8080`

###### Environment variables

| Variable | Value |
| --- | --- |
| `IIIF_CONFIG_JSON` | _Your IIIF config file encoded as a string_ |
| `IIIF_SERVER_CONFIG` | `env:IIIF_CONFIG_JSON` | 
| `AWS_ACCESS_KEY_ID` | _A valid AWS access key ID for talking to the S3 bucket defined in your config file_ | 
| `AWS_SECRET_ACCESS_KEY` | _A valid AWS access secret for talking to the S3 bucket defined in your config file_ |

As of this writing the `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` variables
are necessary because if you specify `credential="iam:"` for an S3 source (in
your IIIF config file) the server fails to start up with a weird error I've
never seen before. Computers...

#### iiif-process

I have had an easier time setting up a Docker-ized `iiif-process` container in ECS and running it as a simple ECS task.

##### Services

You will need to ensure that the service has `Auto-assign public IP(s)` enabled.This is necessary in order to fetch the actual Docker container.

#### Task definitions

Your "task definition" will need a suitable AWS IAM role with the following properties:

* A trust definition with `ecs-tasks.amazonaws.com`

And the following policies assigned to it:

* `AmazonECSTaskExecutionRolePolicy`
* A custom policy with the necessary permissions your task will need to read-from and write-to source and derivative caches (typically S3)

The task should be run in `awsvpc` network mode and required the `FARGATE` capability.

Unlike the `iiif-server` container as of this writing it is not possible to pass in the IIIF config file (or the instructions file) as an environment variable. I've never really loved that approach and want to reconsider it for all the `ifff-` tools.

This means you have a container that can run `/bin/iiif-process` but where does it find any of it's configuration information? The short answer is you don't use the `Dockerfile.process` Dockerfile in this package. Or you create a local copy of it customizing it as necessary.

Instead you should use the `Dockerfile.process.ecs` Dockerfile defined in the [go-iiif-aws](https://github.com/go-iiif/go-iiif-aws) package.

This package will create a custom `iiif-process` container copying a custom IIIF config and instructions file into `/etc/go-iiif/config.json` and `/etc/go-iiif/instructions.json` respectively. This is the container image that you would then upload as a task to your AWS ECS account.

It will also build a `iiif-process-ecs` tool that can be:

* Used to invoke your task directly from the command-line, passing in one or more URIs to process.
* Bundled as an AWS Lambda function that can be run to invoke your task.
* Used to invoke that Lambda function (to invoke your task) from the command-line.

The Dockerfile in the `go-iiif-aws` package will build the `iiif-process` binary from _this package_ but otherwise manages all of the ECS, Lambda and other AWS-specific code in its own codebase.

## Notes

* The `iiif-server` does [not support TLS](https://github.com/go-iiif/go-iiif/issues/5) yet.
* There is no way to change the default `quality` parameter yet. It is `color`.

## Bugs?

Probably. Please consult [the currently known-known issues](https://github.com/go-iiif/go-iiif/issues) and if you don't see what ails you please feel free to add it.

## See also

### IIIF stuff

* http://iiif.io/api/image/2.1/

### go-iiig stuff

* https://github.com/go-iiif/go-iiif-vips
* https://github.com/go-iiif/go-iiif-uri
* https://github.com/go-iiif/go-iiif-www

### Go stuff

* https://github.com/greut/iiif/
* https://github.com/anthonynsimon/bild
* https://github.com/muesli/smartcrop

### Slippy map stuff

* https://github.com/mejackreed/Leaflet-IIIF
* https://github.com/mapbox/leaflet-image

### Blog posts

* http://www.aaronland.info/weblog/2016/09/18/marshmallows/#iiif
* http://www.aaronland.info/weblog/2017/03/05/record/#numbers
* https://labs.cooperhewitt.org/2017/parting-gifts/
* https://millsfield.sfomuseum.org/blog/2018/07/18/iiif/
* https://millsfield.sfomuseum.org/blog/2019/02/12/iiif-aws/
* https://millsfield.sfomuseum.org/blog/2019/11/13/iiif-v2/

### Other stuff

* [Spanking Cat](https://collection.cooperhewitt.org/objects/18382391/)
