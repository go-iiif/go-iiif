# go-iiif-uri

Go package for working with URIs in the `go-iiif` package (and friends).

## Install

You will need to have both `Go` (specifically a version of Go more recent than 1.7 so let's just assume you need [Go 1.11](https://golang.org/dl/) or higher) and the `make` programs installed on your computer. Assuming you do just type:

```
make bin
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Why does this package exist?

This package exists because the `go-iiif:process` code demands a `URI` interface because sometimes URIs are more complicated than simple strings. By extension the code in the `go-iiif-aws` package, which handles invoking the `go-iiif:process` code as ECS or Lambda task _also_ needs to know about said `URI` interfaces.

If we put all the `URI` interfaces in `go-iiif:process` (which we did at first) we run in to the problem where the `go-iiif-aws` package needs to import `go-iiif`. Then we can no longer compile Linux binaries (in `go-iiif-aws`) on anything but a Linux machine because the `go-iiif` code is being invoked which in turn invokes `bimg` which invokes `libvips` which is not set up for cross-compiling.

That's it, really. It is a tiny bit sad-making but also not really a big deal.

## See also

* https://github.com/go-iiif/go-iiif
* https://github.com/go-iiif/go-iiif-aws