# go-iiif

![spanking cat](misc/go-iiif-spanking-cat.png)

This is a fork [@greut's iiif](https://github.com/greut/iiif) package that moves most of the processing logic in to discrete Go packages and defines source, derivative and graphics details in a [JSON config file](README.md#config-files). There is an additional caching layer for both source images and derivatives.

I did this to better understand some of the decisions behind, and to address my own concerns about, version 2 of the IIIF Image API. For the time being this package will probably not support the other IIIF Metadata or Publication APIs. Honestly, as of this writing it may still be lacking some parts of Image API but it's a start and it does all the basics.

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

`iiif-server` is a HTTP server that supports version 2.1 of the [IIIF Image API](). For example:

```
$> bin/iiif-server -config config.json
2016/09/01 15:45:07 Serving 127.0.0.1:8080 with pid 12075

curl -s localhost:8080/184512_5f7f47e5b3c66207_x.jpg/full/full/0/default.jpg
curl -s localhost:8080/184512_5f7f47e5b3c66207_x.jpg/125,15,200,200/full/0/default.jpg
curl -s localhost:8080/184512_5f7f47e5b3c66207_x.jpg/pct:41.6,7.5,40,70/full/0/default.jpg
curl -s localhost:8080/184512_5f7f47e5b3c66207_x.jpg/full/full/270/default.png
```

#### Endpoints

##### GET /level2.json

_Please write me_

##### GET /{ID}/info.json

_Please write me_

##### GET /{ID}/{REGION}/{SIZE}/{ROTATION}/{QUALITY}.{FORMAT}

_Please write me_

##### GET /debug/vars

```
$> curl -s 127.0.0.1:8080/debug/vars | python -mjson.tool | grep Cache
    "CacheHit": 4,
    "CacheMiss": 16,
    "CacheSet": 16,
        "MCacheInuse": 1200,
        "MCacheSys": 16384,

$> curl -s 127.0.0.1:8080/debug/vars | python -mjson.tool | grep Transforms
    "TransformsAvgTimeMS": 1833.875,
    "TransformsCount": 16,
```

#### Notes

* TLS is [not supported yet](https://github.com/thisisaaronland/go-iiif/issues/5).

### iiif-tile-seed

```
$> ./bin/iiif-tile-seed -options /path/to/source/image.jpg

Usage of ./bin/iiif-tile-seed:
  -config string
    	  Path to a valid go-iiif config file
  -refresh
	Refresh a tile even if already exists (default false)
  -scale-factors string
    		 A comma-separated list of scale factors to seed tiles with (default "4")
```

## Config files

There is a [sample config file](config.json.example) included with this repo. For example:

```
{
	"level": {
		"compliance": "2"
	},
	"graphics": {
		"source": { "name": "VIPS" }
	},
	"features": {
		"enable": {},
		"disable": { "rotation": [ "rotationArbitrary"] },
		"append": {}
	},
	"images": {
		"source": { "name": "Disk", "path": "example/images" },
		"cache": { "name": "Memory", "ttl": 300, "limit": 100 }
	},
	"derivatives": {
		"cache": { "name": "Disk", "path": "example/cache" }
	}
}
```

_More better documentation to follow but in the meantime..._

### level

```
	"level": {
		"compliance": "2"
	}
```

### graphics

```
	"graphics": {
		"source": { "name": "VIPS" }
	}
```

Details about how images should be processed.

#### source

Where to find source images.

##### name

Details about what will do the actual image processing.

* VIPS - currently there is only VIPS which is which is an interface to the [libvips](http://www.vips.ecs.soton.ac.uk/index.php?title=VIPS) library.

### features

```
	"features": {
		"enable": {},
		"disable": { "rotation": [ "rotationArbitrary"] },
		"append": {}
	}
```

### images

```
	"images": {
		"source": { "name": "Disk", "path": "example/images" },
		"cache": { "name": "Memory", "ttl": 300, "limit": 100 }
	}
```

Details about source images.

#### source

Where to find source images.

##### name

Possible cache sources for source images are:

* Disk - A locally available filesystem.

_Planned future source providers include reading images via S3 or an OEmbed endpoint._ 

#####  path

The path to a valid directory to find source images.

#### cache

Caching options for source images.

##### name

Possible cache sources for source images are:

* Disk  - Cache images to a locally available filesystem. Until it is possible to read source images from a remote location it's not clear why you would ever do this but I guess that's your business...

* Memory - Cache images in memory.

* Null – Because you must define a caching layer this is here to satify the requirements without actually caching anything, anywhere.

##### ttl

This is only valid for `Memory` caches and indicates the maximum number of seconds an image should live in cache.

##### limit

This is only valid for `Memory` caches and indicates the maximum number of megabytes the cache should hold at any one time.

### derivatives

```
	"derivatives": {
		"cache": { "name": "Disk", "path": "example/cache" }
	}
```

Details about derivative images.

#### cache

Caching options for derivative images.

##### name

Possible cache sources for derivative images are:

* Disk - Tache derivatives to a locally available filesystem.

* Memory - Cache images in memory.

* Null – Because you must define a caching layer this is here to satify the requirements without actually caching anything, anywhere.

_Planned future caching providers include writing derivatives to S3._

##### ttl

This is only valid for `Memory` caches and indicates the maximum number of seconds an image should live in cache.

##### limit

This is only valid for `Memory` caches and indicates the maximum number of megabytes the cache should hold at any one time.

## IIIF image API 2.1

The API specifications can be found on [iiif.io](http://iiif.io/api/image/2.1/index.html).

### [Identifier](http://iiif.io/api/image/2.1/#identifier)

* `filename`: the name of the file **(all the images are in one folder)**

### [Region](http://iiif.io/api/image/2.1/index.html#region)

* `full`: the full image
* `square`: a square area in the picture (centered)
* `x,y,w,h`: extract the specified region (as pixels)
* `pct:x,y,w,h`: extract the specified region (as percentages)

### [Size](http://iiif.io/api/image/2.1/index.html#size)

* `full`: the full image **(deprecated)**
* `max`: the full image
* `w,h`: a potentially deformed image of `w x h` **(not supported)**
* `!w,h`: a non-deformed image of maximum `w x h`
* `w,`: a non-deformed image with `w` as the width
* `,h`: a non-deformed image with `h` as the height
* `pct:n`: a non-deformed image scaled by `n` percent

### [Rotate](http://iiif.io/api/image/2.1/index.html#rotation)

* `n` a clockwise rotation of `n` degrees
* `!n` a flip is done before the rotation

__limitations__ bimg only supports rotations that are multiples of 90.

### [Quality](http://iiif.io/api/image/2.1/index.html#quality)

* `color` image in full colour
* `gray` image in grayscale
* `bitonal` image in either black or white pixels **(not supported)**
* `default` image returned in the server default quality

### [Format](http://iiif.io/api/image/2.1/index.html#format)

* `jpg`
* `png`
* `webp`
* `tiff`

__limitations__ : bimg (libvips) doesn't support writing to `jp2`, `gif` or `pdf`.

### [Profile](http://iiif.io/api/image/2.1/#image-information)

It provides all informations but the available `sizes` and `tiles`. The `sizes`
information would be much better linked with a Cache system.

### [Level2 profile](http://iiif.io/api/image/2.1/#profile-description)

It provides meta-informations about the service. **(incomplete)**

## See also

* http://iiif.io/api/image/2.1/
* https://github.com/greut/iiif/
* https://github.com/h2non/bimg/
* http://www.vips.ecs.soton.ac.uk/index.php?title=VIPS