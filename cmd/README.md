# Command line tools

You will need to have both `Go` (specifically version [1.12](https://golang.org/dl/) or higher) and the `make` programs installed on your computer. Assuming you do just type:

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
	"github.com/go-iiif/go-iiif/v7/tools"
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
	"github.com/go-iiif/go-iiif/v7/tools"
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

### iiif-process

```
$> ./bin/iiif-process -h
Usage of process:
  -config-derivatives-cache-uri string
    	If present this value will be used to assign the 'derivatives.cache.uri' property in the config file. Note: The 'derivatives.cache.uri' property takes precedence over other properties in 'derivatives.cache' block.
  -config-images-source-uri string
    	If present this value will be used to assign the 'images.source.uri' property in the config file. Note: The 'images.source.uri' property takes precedence over other properties in 'images.source' block.
  -config-name string
    	The name of your go-iiif config file. This value will be ignored if -config-source is 'defaults://'. (default "config.json")
  -config-source string
    	A valid Go Cloud bucket URI where your go-iiif config file is located. Optionally, if 'defaults://' is specified then the default config bundled with this package will be used.
  -generate-report-html
    	Generate an HTML page showing all the images listed in a process report.
  -instructions-name string
    	The name of your go-iiif instructions file. This value will be ignored if -instructions-source is 'defaults://'. (default "instructions.json")
  -instructions-source string
    	A valid Go Cloud bucket URI where your go-iiif "instructions" processing file is located. Optionally, if 'defaults://' is specified then the default instructions set bundled with this package will be used.
  -mode string
    	Valid modes are: cli, csv, fsnotify, lambda. (default "cli")
  -report
    	Store a process report (JSON) for each URI in the cache tree.
  -report-source string
    	A valid Go Cloud bucket URI where your report file will be saved. If empty reports will be stored alongside derivative (or cached) images.
  -report-template string
    	A valid URI template for generating process report filenames. (default "process_{sha256_origin}.json")
  -verbose
    	Enabled verbose (debug) loggging.
```

Perform a series of IIIF image processing tasks, defined in a JSON-based "instructions" file, on one or more (IIIF) URIs. For example:

```
$> ./bin/iiif-process -config config.json -instructions instructions.json -uri source/IMG_0084.JPG | jq

{
  "source/IMG_0084.JPG": {
    "dimensions": {
      "b": [
        2048,
        1536
      ],
      "d": [
        320,
        320
      ],
      "o": [
        4032,
        3024
      ]
    },
    "palette": [
      {
        "name": "#b87531",
        "hex": "#b87531",
        "reference": "vibrant"
      },
      {
        "name": "#805830",
        "hex": "#805830",
        "reference": "vibrant"
      },
      {
        "name": "#7a7a82",
        "hex": "#7a7a82",
        "reference": "vibrant"
      },
      {
        "name": "#c7c3b3",
        "hex": "#c7c3b3",
        "reference": "vibrant"
      },
      {
        "name": "#5c493a",
        "hex": "#5c493a",
        "reference": "vibrant"
      }
    ],
    "uris": {
      "b": "source/IMG_0084.JPG/full/!2048,1536/0/color.jpg",
      "d": "source/IMG_0084.JPG/-1,-1,320,320/full/0/dither.jpg",
      "o": "source/IMG_0084.JPG/full/full/-1/color.jpg"
    }
  }
}
```

Images are read-from and stored-to whatever source or derivatives caches defined in your `config.json` file.

#### "instructions" files

An instruction file is a JSON-encoded dictionary. Keys are user-defined and values are dictionary of IIIF one or more transformation instructions. For example:

```
{
    "o": {"size": "full", "format": "", "rotation": "-1" },
    "b": {"size": "!2048,1536", "format": "jpg" },
    "d": {"size": "full", "quality": "dither", "region": "-1,-1,320,320", "format": "jpg" }	
}

```

The complete list of possible instructions is:

```
type IIIFInstructions struct {
	Region   string `json:"region"`
	Size     string `json:"size"`
	Rotation string `json:"rotation"`
	Quality  string `json:"quality"`
	Format   string `json:"format"`
}
```

As of this writing there is no explicit response type for image beyond `map[string]interface{}`. There probably could be but it's still early days.

#### "report" files

"Report" files are JSON files that contain the list of files created, their dimensions and the output of any (IIIF) services that have been configured.

For example, if you ran the following `iiif-process` command:

```
$> go run -mod vendor cmd/iiif-process/main.go \
   -config-source file:///usr/local/go-iiif/docs \
   -instructions-source file:///usr/local/go-iiif/docs \
   -report test.jpg
```

The default `-report-template` URI template is `process_{sha256_origin}.json` so the resultant process report would be created at `test.jpg/process_0d407ee6406a1216f2366674a1a9ff71361d5bef47021f8eb8b51f95e319dd56.json`.

As in: `hex(sha256("test.jpg")) == 0d407ee6406a1216f2366674a1a9ff71361d5bef47021f8eb8b51f95e319dd56.json`.

Currently, there is only one optional suffix (`{sha256_origin}`) defined but in the future the hope is to make these customizable. The output of the report will look something like this, depending on which services are enabled or not:

```
{
  "blurhash": ":JK_E@_4?bM}?vM|.8WB~pt6RjWCRjf6jtWBx^WBNGoLRjoeWAj]ogWBj?j[ofayayofxvaeWBoeWBofRjofozfPj@a{f6j[f6j[kEaxj[a{WBt7WBj[t8j?aeayj[ayayj[",
  "dimensions": {
    "b": [
      1152,
      1536
    ],
    "d": [
      320,
      320
    ],
    "o": [
      2995,
      3993
    ]
  },
  "imagehash": {
    "average": "a:fffdf1f1e1818181",
    "difference": "d:0141050103031303"
  },
  "origin": "test.jpg",
  "origin_fingerprint": "572e4ee59493efcdc4356ba3e142b19661ff60fa",
  "origin_uri": "file:///test.jpg",
  "palette": [
    {
      "name": "#87837f",
      "hex": "#87837f",
      "reference": "vibrant"
    },
    {
      "name": "#c7c4bf",
      "hex": "#c7c4bf",
      "reference": "vibrant"
    },
    {
      "name": "#483c2c",
      "hex": "#483c2c",
      "reference": "vibrant"
    }
  ],
  "uris": {
    "b": "file:///test.jpg/full/!2048,1536/0/color.jpg",
    "d": "file:///test.jpg/-1,-1,320,320/full/0/dither.jpg",
    "o": "file:///test.jpg/full/full/-1/color.jpg"
  }
}
```

#### "lambda" mode

If you are running this tool in Lambda mode you will need to map environment variables to their command line flag equivalents. This is handled automatically so long as the environment variables you set follows these conventions:

* The name of a flag is upper-cased
* Any instances of `-` are replaced by `_`
* The final environment variable is prefixed by `IIIF_`

For example the command line flag `-mode` becomes the AWS Lambda environment variable `IIIF_MODE`.

### iiif-process-and-tile

```
$> ./bin/iiif-process-and-tile -h
Usage of iiif-process-and-tile:
  -config-derivatives-cache-uri string
    	If present this value will be used to assign the 'derivatives.cache.uri' property in the config file. Note: The 'derivatives.cache.uri' property takes precedence over other properties in 'derivatives.cache' block.
  -config-images-source-uri string
    	If present this value will be used to assign the 'images.source.uri' property in the config file. Note: The 'images.source.uri' property takes precedence over other properties in 'images.source' block.
  -config-name string
    	The name of your go-iiif config file. This value will be ignored if -config-source is 'defaults://'. (default "config.json")
  -config-source string
    	A valid Go Cloud bucket URI where your go-iiif config file is located. Optionally, if 'defaults://' is specified then the default config bundled with this package will be used.
  -csv-source string
    	 (default "A valid Go Cloud bucket URI where your CSV tileseed files are located.")
  -endpoint string
    	The endpoint (scheme, host and optionally port) that will serving these tiles, used for generating an 'info.json' for each source image (default "http://localhost:8080")
  -format string
    	A valid IIIF format parameter (default "jpg")
  -generate-report-html
    	Generate an HTML page showing all the images listed in a process report.
  -generate-tiles-html
    	If true then the tiles directory will be updated to include HTML/JavaScript/CSS assets to display tiles as a "slippy" map (using the leaflet-iiif.js library.
  -instructions-name string
    	The name of your go-iiif instructions file. This value will be ignored if -instructions-source is 'defaults://'. (default "instructions.json")
  -instructions-source string
    	A valid Go Cloud bucket URI where your go-iiif "instructions" processing file is located. Optionally, if 'defaults://' is specified then the default instructions set bundled with this package will be used.
  -logfile string
    	Write logging information to this file
  -loglevel string
    	The amount of logging information to include, valid options are: debug, info, status, warning, error, fatal (default "info")
  -mode string
    	Valid modes are: cli, csv, fsnotify, lambda. (default "cli")
  -noextension
    	Remove any extension from destination folder name.
  -processes int
    	The number of concurrent processes to use when tiling images (default 10)
  -quality string
    	A valid IIIF quality parameter - if "default" then the code will try to determine which format you've set as the default (default "default")
  -refresh
    	Refresh a tile even if already exists (default false)
  -report
    	Store a process report (JSON) for each URI in the cache tree.
  -report-source string
    	A valid Go Cloud bucket URI where your report file will be saved. If empty reports will be stored alongside derivative (or cached) images.
  -report-template string
    	A valid URI template for generating process report filenames. (default "process_{sha256_origin}.json")
  -scale-factors string
    	A comma-separated list of scale factors to seed tiles with (default "4")
  -synchronous
    	Run tools synchronously.
  -tiles-prefix string
    	A relative URL to use a prefix when storing tiles.
  -verbose
    	Enabled verbose (debug) loggging.
```

This tool wraps the functionality of the `iiif-process` and `iiif-tile-seed` tools in to a single operation to be performed on one or more URIs.

Processing and tile-seeding operations happen asynchronously by default but can be made to happen sequentially with the `-synchronous` flag.

#### "lambda" mode

If you are running this tool in Lambda mode you will need to map environment variables to their command line flag equivalents. This is handled automatically so long as the environment variables you set follows these conventions:

* The name of a flag is upper-cased
* Any instances of `-` are replaced by `_`
* The final environment variable is prefixed by `IIIF_PROCESS_AND_TILE_`

For example the command line flag `-mode` becomes the AWS Lambda environment variable `IIIF_PROCESS_AND_TILE_MODE`.

### iiif-server

```
$> ./bin/iiif-server -h
Usage of server:
  -config-derivatives-cache-uri string
    	If present this value will be used to assign the 'derivatives.cache.uri' property in the config file. Note: The 'derivatives.cache.uri' property takes precedence over other properties in 'derivatives.cache' block.
  -config-images-source-uri string
    	If present this value will be used to assign the 'images.source.uri' property in the config file. Note: The 'images.source.uri' property takes precedence over other properties in 'images.source' block.
  -config-name string
    	The name of your go-iiif config file. This value will be ignored if -config-source is 'defaults://'. (default "config.json")
  -config-source string
    	A valid Go Cloud bucket URI where your go-iiif config file is located. Optionally, if 'defaults://' is specified then the default config bundled with this package will be used.
  -example
    	Add an /example endpoint to the server for testing and demonstration purposes
  -example-root string
    	An explicit path to a folder containing example assets (default "example")
  -host string
    	Bind the server to this host. THIS FLAG IS DEPRECATED: Please use -server-uri instead.
  -port int
    	Bind the server to this port. THIS FLAG IS DEPRECATED: Please use -server-uri instead.
  -protocol string
    	The protocol for iiif-server server to listen on. Valid protocols are: http, lambda. THIS FLAG IS DEPRECATED: Please use -server-uri instead.
  -server-uri string
    	A valid aaronland/go-http-server URI (default "http://localhost:8080")
  -verbose
    	Enabled verbose (debug) loggging.
```

For example:

```
$> bin/iiif-server -config config.json
2016/09/01 15:45:07 Serving 127.0.0.1:8080 with pid 12075

curl -s localhost:8080/184512_5f7f47e5b3c66207_x.jpg/full/full/0/default.jpg
curl -s localhost:8080/184512_5f7f47e5b3c66207_x.jpg/125,15,200,200/full/0/default.jpg
curl -s localhost:8080/184512_5f7f47e5b3c66207_x.jpg/pct:41.6,7.5,40,70/full/0/default.jpg
curl -s localhost:8080/184512_5f7f47e5b3c66207_x.jpg/full/full/270/default.png
```

`iiif-server` is a HTTP server that supports version 2.1 of the [IIIF Image API](http://iiif.io/api/image/2.1/).

#### "lambda" mode

If you are running this tool in Lambda mode you will need to map environment variables to their command line flag equivalents. This is handled automatically so long as the environment variables you set follows these conventions:

* The name of a flag is upper-cased
* Any instances of `-` are replaced by `_`
* The final environment variable is prefixed by `IIIF_`

For example the command line flag `-mode` becomes the AWS Lambda environment variable `IIIF_MODE`.

#### Endpoints

Although the identifier parameter (`{ID}`) in the examples below suggests that is is only string characters up to and until a `/` character, it can in fact contain multiple `/` separated strings. For example, either of these two URLs is valid

```
http://localhost:8082/191733_5755a1309e4d66a7_k.jpg/info.json
http://localhost:8082/191/733/191733_5755a1309e4d66a7/info.json
```

Where the identified will be interpreted as `191733_5755a1309e4d66a7_k.jpg` and `191/733/191733_5755a1309e4d66a7` respectively. Identifiers containing one or more `../` strings will be made to feel bad about themselves.

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
		"dither",
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
$> ./bin/iiif-tile-seed -h
Usage of tileseed:
  -config-derivatives-cache-uri string
    	If present this value will be used to assign the 'derivatives.cache.uri' property in the config file. Note: The 'derivatives.cache.uri' property takes precedence over other properties in 'derivatives.cache' block.
  -config-images-source-uri string
    	If present this value will be used to assign the 'images.source.uri' property in the config file. Note: The 'images.source.uri' property takes precedence over other properties in 'images.source' block.
  -config-name string
    	The name of your go-iiif config file. This value will be ignored if -config-source is 'defaults://'. (default "config.json")
  -config-source string
    	A valid Go Cloud bucket URI where your go-iiif config file is located. Optionally, if 'defaults://' is specified then the default config bundled with this package will be used.
  -csv-source string
    	 (default "A valid Go Cloud bucket URI where your CSV tileseed files are located.")
  -endpoint string
    	The endpoint (scheme, host and optionally port) that will serving these tiles, used for generating an 'info.json' for each source image (default "http://localhost:8080")
  -format string
    	A valid IIIF format parameter (default "jpg")
  -generate-tiles-html
    	If true then the tiles directory will be updated to include HTML/JavaScript/CSS assets to display tiles as a "slippy" map (using the leaflet-iiif.js library.
  -logfile string
    	Write logging information to this file
  -loglevel string
    	The amount of logging information to include, valid options are: debug, info, status, warning, error, fatal (default "info")
  -mode string
    	Valid modes are: cli, csv, fsnotify, lambda. (default "cli")
  -noextension
    	Remove any extension from destination folder name.
  -processes int
    	The number of concurrent processes to use when tiling images (default 10)
  -quality string
    	A valid IIIF quality parameter - if "default" then the code will try to determine which format you've set as the default (default "default")
  -refresh
    	Refresh a tile even if already exists (default false)
  -scale-factors string
    	A comma-separated list of scale factors to seed tiles with (default "4")
  -verbose
    	Enabled verbose (debug) loggging.
```

Generate (seed) all the tiled derivatives for a source image for use with the [Leaflet-IIIF](https://github.com/mejackreed/Leaflet-IIIF) plugin.

#### iiif-tile-seed and identifiers

Identifiers for source images can be passed to `iiif-tiles-seed` in of two way:

1. A space-separated list of identifiers
2. A space-separated list of _comma-separated_ identifiers indicating the identifier for the source image followed by the identifier for the newly generated tiles

For example:

```
$> ./bin/iiif-tile-seed -options 191733_5755a1309e4d66a7_k.jpg
```

Or:

```
$> ./bin/iiif-tile-seed -options 191733_5755a1309e4d66a7_k.jpg,191/733/191733_5755a1309e4d66a7
```

In many cases the first option will suffice but sometimes you might need to create new identifiers or structure existing identifiers according to their output, for example avoiding the need to store lots of file in a single directory. It's up to you.

You can also run `iiif-tile-seed` pass a list of identifiers as a CSV file. To do so include the `-mode csv` argument, like this:

```
$> ./bin/iiif-tile-seed -options -mode csv CSVFILE
```

Your CSV file must contain a header specifying a `source_id` and `alternate_id` column, like this:

```
source_id,alternate_id
191733_5755a1309e4d66a7_k.jpg,191733_5755a1309e4d66a7
```

While all columns are required if `alternate_id` is empty the code will simply default to using `source_id` for all operations.

_Important: The use of alternate IDs is not fully supported by `iiif-server` yet. Which is to say to the logic for how to convert a source identifier to an alternate identifier is still outside the scope of `go-iiif` so unless you have pre-rendered all of your tiles or other derivatives (in which case the check for cached derivatives at the top of the imgae handler will be triggered) then the server won't know where to write new alternate files._

#### "lambda" mode

If you are running this tool in Lambda mode you will need to map environment variables to their command line flag equivalents. This is handled automatically so long as the environment variables you set follows these conventions:

* The name of a flag is upper-cased
* Any instances of `-` are replaced by `_`
* The final environment variable is prefixed by `IIIF_`

For example the command line flag `-mode` becomes the AWS Lambda environment variable `IIIF_MODE`.
