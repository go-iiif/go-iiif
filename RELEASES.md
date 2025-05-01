# Releases

The current release is `github.com/go-iiif/go-iiif/v8`.

### v8

#### Config changes

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
    	"extruder": { "uri": "vibrant://", "count": 5 },
    	"grid": { "uri": "euclidian://" },
	"palettes": [
		    { "uri": "crayola://" },
		    { "uri": "css4://" }
        ]
    },
    "blurhash_service": { "x": 8, "y": 8, "size": 200 },
    "imagehash_service": {},
```

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


