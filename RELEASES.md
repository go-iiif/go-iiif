# Releases

The current release is `github.com/go-iiif/go-iiif/v8`.

### v8

Version 8.0.0 does NOT change anything in the core IIIF image processing code but has been updated to reflect changes the [aaronland/go-colours](https://github.com/aaronland/go-colours) package which, by extension, introduce a number of backwards incompatible changes to [the "palette" (colour extraction) service](https://github.com/go-iiif/go-iiif/tree/main/service) and the `go-iiif` config file definitions.

#### Config changes

Specific changes to the config file are:

1. The top-level `palette`, `blurhash` and `imagehash` (service) config keys have been replaced with `palette_service`, `blurhash_service` and `imagehash_service` respectively.
2. The `name` property for `aaronland/go-colours` operators has been replaced with a `uri` property. The URIs for these operators map their corresponding types in the `aaronland/go-colours` package.
3. The `palette_service` config now expects extruders to be defined as a list (or dictionaries) in the `extruders` key.
3. The `palette_service` config now expects (colour) palettes to be defined as a list (or dictionaries) in the `palettes` key.

For example, this:

```
    "profile": {
    	"services": {
		    "enable": [
		    	"palette",
			"blurhash",
		    	"imagehash"
		    ]
	}
    },
    "palette": {
    	"extruder": { "name": "vibrant", "count": 5 },
    	"grid": { "name": "euclidian" },
	"palette": [
		    { "name": "crayola" },
		    { "name": "css4" }
        ]
    },
    "blurhash": { "x": 8, "y": 8, "size": 200 },
    "imagehash": {},
```

Becomes:

```
    "profile": {
    	"services": {
		    "enable": [
		    	"palette://",
			"blurhash://",
		    	"imagehash://"
		    ]
	}
    },
    "palette_service": {
    	"extruders": [
		{ "uri": "vibrant://", "count": 5 }
	],
    	"grid": { "uri": "euclidian://" },
	"palettes": [
		    { "uri": "crayola://" },
		    { "uri": "css4://" }
        ]
    },
    "blurhash_service": { "x": 8, "y": 8, "size": 200 },
    "imagehash_service": {},
```

The built-in config file defined in the `defaults` package has been updated accordingly, as have the examples in the [docs/config](docs/config) folder.

### v7

Version 7.0.0 is introduces a number of backwards incompatible changes to how the package is structured and the interface and method signatures it exposes. The entire `tools` subpackage has been replaced by `app`. Additionally:

* The `cmd/iiif-process-and-tile` tool has been removed. Other command line tools should not have any user-facing changes save for the removal of deprecated flags.
* Some properties in `config.Config` blocks have been removed in favour of a single URI-style declarative syntax.
* The `disk`, `memory` and `s3` source and cache providers have been removed and been replaced by the [gocloud.dev/blob](#) equivalents (`file://`, `mem://` and `s3://` or `s3blob://`).

#### Config changes

##### Graphics

Prior to version 7 the syntax for graphics drivers was:

```
    "graphics": {
        "source": { "name": "native" }
    },
```

Which now becomes:

```
    "graphics": {
        "driver": "native://"
    },
```

##### Images and derivatives

Prior to version 7 the syntax for image and derivative sources was:

```
    "images": {
	"source": { "name": "Blob", "path": "s3blob://{BUCKET}?prefix={PREFIX}/&region={REGION}&credentials={CREDENTIALS}" },
	"cache": { "name": "Memory", "ttl": 300, "limit": 100 }
    },
```

Which now becomes:

```
    "images": {
	"source": { "uri": "s3blob://{BUCKET}?prefix={PREFIX}/&region={REGION}&credentials={CREDENTIALS}" },
	"cache": { "uri": "memory://" }
    },
```   

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


