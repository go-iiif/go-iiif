# drivers

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
	NewImageFromConfigWithSource(context.Context, *iiifconfig.Config, iiifsource.Source, string) (iiifimage.Image, error)
	NewImageFromConfigWithCache(context.Context, *iiifconfig.Config, iiifcache.Cache, string) (iiifimage.Image, error)
	NewImageFromConfig(context.Context, *iiifconfig.Config, string) (iiifimage.Image, error)
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
	iiifdriver.RegisterDriver(context.Background(), "native", dr)
}
```

And then in your code you might do something like this:

```
import (
	"context"

	_ "github.com/go-iiif/go-iiif/v6/native"
	
	"github.com/aaronland/gocloud-blob/bucket"	
	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
	iiifdriver "github.com/go-iiif/go-iiif/v6/driver"	
)

ctx := context.Background()
	
config_bucket, _ := bucket.OpenBucket(ctx, "file:///etc/go-iiif")

cfg, _ := config.NewConfigFromBucket(ctx, config_bucket, "config.json")

driver, _ := iiifdriver.NewDriver(ctx, cfg.Graphics.Driver)
```

That's really the only change to existing code. Careful readers may note the calls to `bucket.OpenBucket` and `config.NewConfigFromBucket` to load `go-iiif` configuration files. This is discussed below. In the meantime the only other change is to update the previously default `graphics.source` property in the configuration file from `VIPS` (or `vips`) to `native`. For example:

```
    "graphics": {
	"driver": { "name": "vips://" }
    }
```

Becomes:

```
    "graphics": {
	"driver": { "name": "native://" }    
    }
```

The value of the `graphics.driver` property should match the URI that driver used to register itself with `go-iiif`. For working examples you should consult either of the following packages:

* https://github.com/go-iiif/go-iiif/v6/tree/main/native
* https://github.com/go-iiif/go-iiif-vips
