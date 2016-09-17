# go-iiif

![spanking cat](misc/go-iiif-spanking-cat.png)

## What is this?

This is a fork [@greut's iiif](https://github.com/greut/iiif) package that moves most of the processing logic in to discrete Go packages and defines source, derivative and graphics details in a [JSON config file](README.md#config-files). There is an additional caching layer for both source images and derivatives.

I did this to better understand the architecture behind (and to address my own concerns about) version 2 of the [IIIF Image API](http://iiif.io/api/image/2.1/index.html).

For the time being this package will probably not support the other IIIF Metadata or Publication APIs. Honestly, as of this writing it may still be lacking some parts of Image API but it's a start and it does all the basics.

_And by "forked" I mean that [@greut](https://github.com/greut) and I decided that [it was best](https://github.com/greut/iiif/pull/2) for this code and his code to wave at each other across the divide but not necessarily to hold hands._

## Setup

Currently all the image processing is handled by the [bimg](https://github.com/h2non/bimg/) Go package which requires the [libvips](http://www.vips.ecs.soton.ac.uk/index.php?title=VIPS) C library be installed. There is a detailed [setup script](ubuntu/setup.sh) available for Ubuntu. Eventually there will be pure-Go alternatives for wrangling images. Otherwise all other depedencies are included with this repository in the [vendor](vendor) directory.

Once you have things like`Go` and `libvips` installed just type:

```
$> make bin
```

## Usage

`go-iiif` was designed to expose all of its functionality outside of the included tools although that hasn't been documented yet. The [source code for the iiif-tile-seed tool](cmd/iiif-tile-seed.go) is a good place to start poking around if you're curious.

## Tools

### iiif-server

```
$> bin/iiif-server -config config.json
2016/09/01 15:45:07 Serving 127.0.0.1:8080 with pid 12075

curl -s localhost:8080/184512_5f7f47e5b3c66207_x.jpg/full/full/0/default.jpg
curl -s localhost:8080/184512_5f7f47e5b3c66207_x.jpg/125,15,200,200/full/0/default.jpg
curl -s localhost:8080/184512_5f7f47e5b3c66207_x.jpg/pct:41.6,7.5,40,70/full/0/default.jpg
curl -s localhost:8080/184512_5f7f47e5b3c66207_x.jpg/full/full/270/default.png
```

`iiif-server` is a HTTP server that supports version 2.1 of the [IIIF Image API](http://iiif.io/api/image/2.1/).

#### Endpoints

##### GET /{ID}/info.json

```
$> curl -s http://localhost:8082/184512_5f7f47e5b3c66207_x.jpg/info.json | python -mjson.tool
{
    "@context": "http://iiif.io/api/image/2/context.json",
    "@id": "http://localhost:8082/184512_5f7f47e5b3c66207_x.jpg",
    "@type": "iiif:Image",
    "height": 4096,
    "profile": [
        "http://iiif.io/api/image/2/level2.json",
        {
            "formats": [
                "tif",
                "webp",
                "jpg",
                "png"
            ],
            "qualities": [
                "default",
                "color"
            ],
            "supports": [
                "full",
                "regionByPx",
                "regionByPct",
                "sizeByWh",
                "full",
                "max",
                "sizeByW",
                "sizeByH",
                "sizeByPct",
                "sizeByConfinedWh",
                "none",
                "rotationBy90s",
                "mirroring",
                "baseUriRedirect",
                "cors",
                "jsonldMediaType"
            ]
        }
    ],
    "protocol": "http://iiif.io/api/image",
    "width": 3897
}
```

Return the [profile description](http://iiif.io/api/image/2.1/#profile-description) for an identifier.

##### GET /{ID}/{REGION}/{SIZE}/{ROTATION}/{QUALITY}.{FORMAT}

```
$> curl -s http://localhost:8082/184512_5f7f47e5b3c66207_x.jpg/pct:41,7,40,70/,250/0/default.jpg
```

Return an image derived from an identifier and one or more [IIIF parameters](http://iiif.io/api/image/2.1/#image-request-parameters). For example:

![spanking cat, cropped](misc/go-iiif-crop.jpg)

##### GET /debug/vars

```
$> curl -s 127.0.0.1:8080/debug/vars | python -mjson.tool | grep Cache
    "CacheHit": 4,
    "CacheMiss": 16,
    "CacheSet": 16,

$> curl -s 127.0.0.1:8080/debug/vars | python -mjson.tool | grep Transforms
    "TransformsAvgTimeMS": 1833.875,
    "TransformsCount": 16,
```

This exposes all the usual Go [expvar](https://golang.org/pkg/expvar/) debugging output along with the following additional properies:

* CacheHit - _the total number of (derivative) images successfully returned from cache_
* CacheMiss - _the total number of (derivative) images not found in the cache_
* CacheSet - _the total number of (derivative) images added to the cache_
* TransformsAvgTimeMS - _the average amount of time in milliseconds to transforms a source image in to a derivative_
* TransformsCount - _the total number of source images transformed in to a derivative_

_Note: This endpoint is only available from the machine the server is running on._

### iiif-tile-seed

```
$> ./bin/iiif-tile-seed -options /path/to/source/image.jpg

Usage of ./bin/iiif-tile-seed:
  -config string
	Path to a valid go-iiif config file
  -endpoint string
	The endpoint (scheme, host and optionally port) that will serving these tiles, used for generating an 'info.json' for each source image (default "http://localhost:8080")
  -refresh
	Refresh a tile even if already exists (default false)
  -scale-factors string
	A comma-separated list of scale factors to seed tiles with (default "4")
```

Generate (seed) all the tiled derivatives for a source image for use with the [Leaflet-IIIF]() plugin.

## Config files

There is a [sample config file](config.json.example) included with this repo. Proper documentation is being written but right now the easiest way to understand config files is that consist of five top-level groupings, with nested section-specific details. They are:

### level

```
	"level": {
		"compliance": "2"
	}
```

Indicates which level of IIIF Image API compliance the server (or associated tools) should support. Basically, there is no reason to ever change this right now.

### graphics

```
	"graphics": {
		"source": { "name": "VIPS" }
	}
```

Details about how images should be processed. Because only [libvips]() is supported for image processing right now there is no reason to change this.

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

Finally, maybe you've got an IIIF implementation that [knows how to do things not defined in the spec](https://github.com/thisisaaronland/go-iiif/issues/1). This is also where you would add them.

#### compliance 

Here's how that dynamic plays out in reality. The table below lists all the IIIF parameters and their associate features. Each feature lists its syntax and whether or not it is required and supported [according to the official spec](compliance/level2.go) but then also according to the [example `go-iiif` config file](config.json.example), included with this repo.

_This table was generated using the [iiif-dump-config](cmd/iiif-dump-config.go) tool and if anyone can tell me how to make Markdown tables (in GitHub) render colours I would be grateful._

##### [region](http://iiif.io/api/image/2.1/index.html#region)
| feature | syntax | required (spec) | supported (spec) | required (config) | supported (config) |
|---|---|---|---|---|---|
| **full** | full | <span style="color:green;">true</span> | <span style="color:green;">true</span> | <span style="color:green;">**true**</span> | <span style="color:green;">**true**</span> |
| **regionByPct** | pct:x,y,w,h | <span style="color:green;">true</span> | <span style="color:green;">true</span> | <span style="color:green;">**true**</span> | <span style="color:green;">**true**</span> |
| **regionByPx** | x,y,w,h | <span style="color:green;">true</span> | <span style="color:green;">true</span> | <span style="color:green;">**true**</span> | <span style="color:green;">**true**</span> |
| **regionSquare** | square | <span style="color:red;">false</span> | <span style="color:red;">false</span> | <span style="color:red;">false</span> | <span style="color:red;">false</span> |

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

##### [quality](http://iiif.io/api/image/2.1/index.html#quality)
| feature | syntax | required (spec) | supported (spec) | required (config) | supported (config) |
|---|---|---|---|---|---|
| **bitonal** | bitonal | <span style="color:green;">true</span> | <span style="color:green;">true</span> | <span style="color:green;">**true**</span> | <span style="color:red;">false</span> |
| **color** | color | <span style="color:red;">false</span> | <span style="color:green;">true</span> | <span style="color:red;">false</span> | <span style="color:green;">**true**</span> |
| **default** | default | <span style="color:green;">true</span> | <span style="color:green;">true</span> | <span style="color:green;">**true**</span> | <span style="color:green;">**true**</span> |
| **gray** | gray | <span style="color:red;">false</span> | <span style="color:red;">false</span> | <span style="color:red;">false</span> | <span style="color:red;">false</span> |

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
		"append": { "format": {
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

##### Disk

```
	"images": {
		"source": { "name": "Disk", "path": "example/images" }
	}
```

Fetch source images from a locally available filesystem.

##### URI

```
	"images": {
		"source": { "name": "URI", "path": "https://images.collection.cooperhewitt.org/{id}" }
	}
```

Fetch source images from a remote URI. The `path` parameter must be a valid (Level 4) [URI Template](http://tools.ietf.org/html/rfc6570) with an `{id}` placeholder.

#### images.cache

Caching options for source images.

##### Disk

```
	"images": {
		"cache": { "name": "Disk", "path": "example/cache" }
	}
```

Cache images to a locally available filesystem.

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

##### Disk

```
	"derivatives": {
		"cache": { "name": "Disk", "path": "example/cache" }
	}
```

Cache images to a locally available filesystem.

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

## Example

The `iiif-server` tool comes with a canned example (consisting of exactly one image) so you can see things in the context of a slippy map. Here's what you need to do to get it set up:

First, make sure have a valid `go-iiif` config file. If you don't then you can copy the example config included in this repo:

```
$> cp config.json.example config.json
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

If you've pre-seed your tiles then if you open up the network console in your browser then you should see something like this, namely that the individual tiles are returned speedy and fast:

![spanking cat](misc/go-iiif-example-cached.png)
 
## Performance and load testing

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

In the end, I ran out of inodes. 

The current strategy for seeding tiles may also be directly responsible for some of the bottlenecks. Specifically, when processing large volumes of images (defined in a CSV file) the `ifff-tile-seed` will spawn and queue as many concurrent Go routines as there are CPUs. For each of those processes then another (n) CPUs * 2 subprocesses will be spawned to generate tiles. Maybe this is just too image concurrent image processing routines to have? I mean it works but still... Or maybe it's just that every one is waiting for bytes to be written to disk. Or all of the above. I'm not sure yet.

## Notes

* The `iiif-server` does [not support TLS](https://github.com/thisisaaronland/go-iiif/issues/5) yet.
* There is no way to change the default `quality` parameter yet. It is `color`.

## See also

### IIIF stuff

* http://iiif.io/api/image/2.1/

### Go stuff

* https://github.com/greut/iiif/
* https://github.com/h2non/bimg/
* http://www.vips.ecs.soton.ac.uk/index.php?title=VIPS

### Other stuff

* [Spanking Cat](https://collection.cooperhewitt.org/objects/18382391/)
