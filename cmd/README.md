# Command line tools

You will need to have both [Go](https://go.dev/dl) and the `make` programs installed on your computer. Assuming you do just type:

```
$> make cli-tools
cd ../ && make cli-tools && cd -
go build -mod vendor -ldflags="-s -w" -o bin/iiif-server cmd/iiif-server/main.go
go build -mod vendor -ldflags="-s -w" -o bin/iiif-tile-seed cmd/iiif-tile-seed/main.go
go build -mod vendor -ldflags="-s -w" -o bin/iiif-transform cmd/iiif-transform/main.go
go build -mod vendor -ldflags="-s -w" -o bin/iiif-process cmd/iiif-process/main.go
go build -mod vendor -ldflags="-s -w" -o bin/iiif-process-and-tile cmd/iiif-process-and-tile/main.go
go build -mod vendor -ldflags="-s -w" -o bin/iiif-dump-config cmd/iiif-dump-config/main.go
/usr/local/src/go-iiif/cmd
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

As of version 2 all of the logic, including defining and parsing command line arguments, for any `go-iiif` tool that performs image processing has been moved in to the `tools` package. This change allows non-core image processing packages (like [go-iiif-vips](https://github.com/go-iiif/go-iiif-vips)) to more easily re-use functionality defined in the core `go-iiif` package. For example:

```
package main

import (
	"context"
	
	_ "github.com/go-iiif/go-iiif-vips"
	"github.com/go-iiif/go-iiif/v6/tools"
)

func main() {
	tool, _ := tools.NewProcessTool()
	tool.Run(context.Background())
}
```

Under the hood, the `tool.Run()` command is doing tool-specific work to define, parse and set command line flags and eventually invoking its `RunWithFlagSet()` method. For example:

```
package main

import (
	"context"
	_ "github.com/go-iiif/go-iiif-vips"
	"github.com/go-iiif/go-iiif/v6/tools"
	"flag"
	"github.com/sfomuseum/go-flags"	
)

func main() {
	tool, _ := tools.NewProcessTool()

	fs := flag.NewFlagSet("process", flag.ExitOnError)

	tools.AppendCommonProcessToolFlags(ctx, fs)
	tools.AppendProcessToolFlags(ctx, fs)

	flags.Parse(fs)
	flags.SetFlagsFromEnvVars(fs, "IIIF_PROCESS")

	tool.RunWithFlagSet(context.Background(), fs)
}
```

For a complete example of how this all works, and how it can be used to stitch to together custom IIIF processing tools, take a look at the source code for the [cmd/iiif-process-and-tile](cmd/iiif-process-and-tile/main.go) tool.
