# go-iiif-uri

Go package for working with URIs in the `go-iiif` package (and friends).

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/go-iiif/go-iiif-uri.svg)](https://pkg.go.dev/github.com/go-iiif/go-iiif-uri)

Documentation is incomplete.

## Why does this package exist?

This package exists because the `go-iiif:process` code demands a `URI` interface because sometimes URIs are more complicated than simple strings. By extension the code in the `go-iiif-aws` package, which handles invoking the `go-iiif:process` code as ECS or Lambda task _also_ needs to know about said `URI` interfaces.

If we put all the `URI` interfaces in `go-iiif:process` (which we did at first) we run in to the problem where the `go-iiif-aws` package needs to import `go-iiif`. Then we can no longer compile Linux binaries (in `go-iiif-aws`) on anything but a Linux machine because the `go-iiif` code is being invoked which in turn invokes `bimg` which invokes `libvips` which is not set up for cross-compiling.

That's it, really. It is a tiny bit sad-making but also not really a big deal.

## URIs

[go-iiif-uri](https://github.com/go-iiif/go-iiif-uri) URI strings are still a work in progress. While they may still change a bit around the edges efforts will be made to ensure backwards compatibility going forward.

Support for different [go-iiif-uri](https://github.com/go-iiif/go-iiif-uri) URIs is supported using `database/sql` -like "drivers" for packages that support the `URI` interface:

```
type URI interface {
	Scheme() string     
	String() string
	Origin() string
	Target(*url.Values) (string, error)
}
```

## Schemes

`go-iiif-uri` URI strings are defined by a named scheme which indicates how an URI should be processed, a path which is a reference to an image and zero or more query parameters which are the specific instructions for processing the URI.

The following schemes are registed by default when you import `go-iiif-uri`.

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
| id | string  | yes |
| label | string | yes |
| format | string | yes |
| original | string | no |
| secret | string | no |
| secret_o | string | no |
| prefix | string | no |
| ensure-int | bool | no |

If either the `secret` or `secret_o` parameters are absent they will be auto-generated.

The default prefix for `idsecret` URIs is to derive nested folders from the ID property. For example the default prefix for ID `34436377369` would be `344/363/773/69`. You can override this by specifying a `?prefix={CUSTOM_PREFIX}` in your `idsecret://` URI string.

### rewrite

```
rewrite:///path/to/source/image.jpg?target=/path/to/target/picture.jpg
```

The `rewrite://` URI scheme is a variant of the `file://` URI scheme except that the `target` query parameter is required and it will be used to redefine the final URI, rather than just its directory tree, of the processed image.

| Name | Type | Required |
| --- | --- | --- |
| target | string | yes |

## Example

Here's a excerpted example taken from the `go-iiif` [process/parallel.go](https://github.com/go-iiif/go-iiif/blob/master/process/parallel.go) package that processes a single source image, defined as an `idsecret://` URI, in to multiple derivatives defined in an "instructions" file.

The `idsecret://` URI is output as a string using the instructions set to define the `label` and other query parameters. That string is then used to create a new `rewrite://` URI where source is derived from the original `idsecret://` URI and the target is newly generate URI string.

```
go func(u iiifuri.URI, label Label, i IIIFInstructions) {

	var process_uri iiifuri.URI

	switch u.Scheme() {
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

## See also

* https://github.com/go-iiif/go-iiif
* https://github.com/go-iiif/go-iiif-aws