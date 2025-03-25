# go-iiif

![spanking cat](misc/go-iiif-spanking-cat.png)

## Motivation

This began as a fork of [@greut's iiif](https://github.com/greut/iiif) package that moves all of the processing logic for the [IIIF Image API](http://iiif.io/api/image/) in to discrete Go packages and defines source, derivative and graphics details in a [JSON config file](README.md#config-files). There is an additional caching layer for both source images and derivatives.

I did this to better understand the architecture behind (and to address my own concerns about) the [IIIF Image API](http://iiif.io/api/image/2.1/index.html). For the time being this package will probably not support the other IIIF Metadata or Publication APIs.

_And by "forked" I mean that [@greut](https://github.com/greut) and I decided that [it was best](https://github.com/greut/iiif/pull/2) for this code and his code to wave at each other across the divide but not necessarily to hold hands._

## Releases

The current release is `github.com/go-iiif/go-iiif/v6`.

Documentation for releases has been moved in to [RELEASES.md](RELEASES.md).

## Drivers

`go-iiif` was first written with the [libvips](https://github.com/jcupitt/libvips) library and [bimg](https://github.com/h2non/bimg/) Go wrapper for image processing. `libvips` is pretty great but it introduces non-trivial build and setup requirements. As of version 2.0 `go-iiif` no longer uses `libvips` by default but instead does all its image processing using native (Go) code. This allows `go-iiif` to run on any platform supported by Go without the need for external dependencies.

Detailed documentation for drivers has been moved in to [driver/README.md](driver/README.md])

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

### Example

Here's a excerpted example taken from the [process/parallel.go](process/parallel.go) package that processes a single source image, defined as an `idsecret://` URI, in to multiple derivatives defined in an "instructions" file.

The `idsecret://` URI is output as a string using the instructions set to define the `label` and other query parameters. That string is then used to create a new `rewrite://` URI where source is derived from the original `idsecret://` URI and the target is newly generate URI string.

```
go func(ctx context.Context, u iiifuri.URI, label Label, i IIIFInstructions) {

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

	new_uri, im, _ := pr.ProcessURIWithInstructions(ctx, process_uri, label, i)
	// do something with new_uri and im here...
	
}(...)
```

## Command line tools

`go-iiif` was designed to expose all of its functionality outside of the included tools although that hasn't been documented yet. The source code for the [iiif-tile-seed](cmd/iiif-tile-seed.go), [iiif-transform](cmd/iiif-transform.go) and [iiif-process](cmd/iiif-process.go) tools is a good place to start poking around if you're curious.

### Building

Run the handy `cli` Makefile target to build all the tools:

```
$> make cli
go build -mod vendor -ldflags="-s -w" -o bin/iiif-server cmd/iiif-server/main.go
go build -mod vendor -ldflags="-s -w" -o bin/iiif-tile-seed cmd/iiif-tile-seed/main.go
go build -mod vendor -ldflags="-s -w" -o bin/iiif-transform cmd/iiif-transform/main.go
go build -mod vendor -ldflags="-s -w" -o bin/iiif-process cmd/iiif-process/main.go
go build -mod vendor -ldflags="-s -w" -o bin/iiif-dump-config cmd/iiif-dump-config/main.go
```

### iiif-transform

Transform one or more images using the IIIF API. For detailed usage consult [cmd/iiif-transform/README.md](cmd/iiif-transform/README.md)

### iiif-tile-seed

For detailed usage consult [cmd/iiif-tile-seed/README.md](cmd/iiif-tile-seed/README.md)

### iiif-process

Generate IIIF Level-0 image tiles for one or images. For detailed usage consult [cmd/iiif-process/README.md](cmd/iiif-proces/README.md)

### iiif-server

Expose the IIIF Image API via an HTTP endpoint. For detailed usage consult [cmd/iiif-server/README.md](cmd/iiif-server/README.md)

### iiif-dump-config

Emit a go-iiif config file as Markdown. For detailed usage consult [cmd/iiif-dump-config/README.md](cmd/iiif-dump-config/README.md)

## Config files

Documentation for config files has been moved in to [config/README.md](config/README.md].

## Examples

The easiest way to try things out is to use the handy `debug-{SOMETHING}` Makefile targets which will perform operations on files bundled with this package (in the [fixtures](fixtures) directory).

### Generating Level-0 tiles

```
$> make debug-seed
if test -d /usr/local/src/go-iiif/fixtures/cache/spank; then rm -rf /usr/local/src/go-iiif/fixtures/cache/spank; fi
go run cmd/iiif-tile-seed/main.go \
		-config-images-source-uri file:///usr/local/src/go-iiif/fixtures/images \
		-config-derivatives-cache-uri file:///usr/local/src/go-iiif/fixtures/cache \
		-verbose \
		-generate-html \
		'rewrite:///spanking-cat.jpg?target=spank'
2025/03/24 17:55:42 DEBUG Verbose logging enabled
2025/03/24 17:55:42 DEBUG New tiled image origin=spanking-cat.jpg target=spank

... time passes, with lots of debugging information

2025/03/24 17:56:08 DEBUG Tile seeding complete source=spanking-cat.jpg target=spank count=340
2025/03/24 17:56:08 INFO Generate HTML index page for tiles source=spanking-cat.jpg alt=spank
2025/03/24 17:56:08 DEBUG Successfully wrote blob "bucket uri"=file:///usr/local/src/go-iiif/fixtures/cache uri=spank/leaflet.iiif.bundle.js
2025/03/24 17:56:08 DEBUG Successfully wrote blob "bucket uri"=file:///usr/local/src/go-iiif/fixtures/cache uri=spank/leaflet.css
2025/03/24 17:56:08 DEBUG Successfully wrote blob "bucket uri"=file:///usr/local/src/go-iiif/fixtures/cache uri=spank/index.html
2025/03/24 17:56:08 DEBUG Time to seed tiles source=spanking-cat.jpg target=spank time=25.699858709s
```

And then:

```
$> open fixtures/cache/spank/index.html
```

### Generating derivatives using an "instructions" file

```
$> make debug-process
if test -d /usr/local/src/go-iiif/fixtures/cache/999; then rm -rf /usr/local/src/go-iiif/fixtures/cache/999; fi
go run cmd/iiif-process/main.go \
		-config-derivatives-cache-uri file:///usr/local/src/go-iiif/fixtures/cache \
		-config-images-source-uri file:///usr/local/src/go-iiif/fixtures/images \
		-report \
		-report-bucket-uri file:///usr/local/src/go-iiif/fixtures/reports \
		-report-html \
		-verbose \
		'idsecret:///spanking-cat.jpg?id=9998&secret=abc&secret_o=def&format=jpg&label=x'
2025/03/24 17:57:13 DEBUG Verbose logging enabled

... time passes, with lots of debugging information

2025/03/24 17:57:17 DEBUG Successfully wrote blob "bucket uri"=file:///usr/local/src/go-iiif/fixtures/cache uri=999/8/9998_abc_k.jpg
2025/03/24 17:57:17 DEBUG Return transformation uri="rewrite:///spanking-cat.jpg?target=999%2F8%2F9998_abc_k.jpg" origin=spanking-cat.jpg target=999/8/9998_abc_k.jpg "source cache"=memory:// "destination cache"=file:///usr/local/src/go-iiif/fixtures/cache "new uri"=file:///999/8/9998_abc_k.jpg
2025/03/24 17:57:17 DEBUG Successfully wrote blob "bucket uri"=file:///usr/local/src/go-iiif/fixtures/cache uri=999/8/index.html
```

And then:

```
$> open fixtures/cache/999/8/index.html
```

### Running a IIIF API endpoint

```
$> make debug-server
mkdir -p fixtures/cache
go run cmd/iiif-server/main.go \
		-config-derivatives-cache-uri file:///usr/local/src/go-iiif/fixtures/cache \
		-config-images-source-uri file:///usr/local/src/go-iiif/fixtures/images \
		-example \
		-verbose
2025/03/24 17:55:18 DEBUG Verbose logging enabled
2025/03/24 17:55:18 INFO Listening for requests address=http://localhost:8080
```

## Performance and load testing

For processing large, or large volumes of, images the bottlenecks will be:

* CPU usage crunching pixels
* Disk I/O writing tiles to disk
* Running out of inodes

That said on a machine with 8 CPUs and 32GB RAM I was able to run the machine hot with all the CPUs pegged at 100% usage and seed 100, 000 (2048x pixel) images yielding a little over 3 million, or approximately 70GB of, tiles in 24 hours. Some meaningful but not overwhelming amount of time was spent fetching source images across the network so presumably things would be faster reading from a local filesystem.

Memory usage across all the `iiif-tile-seed` processes never went above 5GB and, in the end, I ran out of inodes.

The current strategy for seeding tiles may also be directly responsible for some of the bottlenecks. Specifically, when processing large volumes of images (defined in a CSV file) the `ifff-tile-seed` will spawn and queue as many concurrent Go routines as there are CPUs. For each of those processes then another (n) CPUs * 2 subprocesses will be spawned to generate tiles. Maybe this is just too image concurrent image processing routines to have? I mean it works but still... Or maybe it's just that every one is waiting for bytes to be written to disk. Or all of the above. I'm not sure yet.

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
