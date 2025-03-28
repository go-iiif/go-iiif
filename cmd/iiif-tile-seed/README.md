### iiif-tile-seed

```
$> bin/iiif-tile-seed -h
Generate IIIF Level-0 image tiles for one or images.

Usage:
	 bin/iiif-tile-seed[options] uri(N) uri(N)

Valid options are:
  -config-derivatives-cache-uri string
    	If present this value will be used to assign the 'derivatives.cache.uri' property in the config file. Note: The 'derivatives.cache.uri' property takes precedence over other properties in 'derivatives.cache' block.
  -config-images-source-uri string
    	If present this value will be used to assign the 'images.source.uri' property in the config file. Note: The 'images.source.uri' property takes precedence over other properties in 'images.source' block.
  -config-name string
    	The name of your go-iiif config file. This value will be ignored if -config-source is 'defaults://'. (default "config.json")
  -config-source string
    	A valid Go Cloud bucket URI where your go-iiif config file is located. Optionally, if 'defaults://' is specified then the default config bundled with this package will be used. (default "defaults://")
  -endpoint string
    	The endpoint (scheme, host and optionally port) that will serving these tiles, used for generating an 'info.json' for each source image (default "http://localhost:8080")
  -format string
    	A valid IIIF format parameter (default "jpg")
  -generate-html
    	If true then the tiles directory will be updated to include HTML/JavaScript/CSS assets to display tiles as a "slippy" map (using the leaflet-iiif.js library.
  -mode string
    	Valid options are: cli, csv, fsnotify, lambda (default "cli")
  -no-extension
    	Remove any extension from destination folder name. For example the target (destination) folder for tiles produced from a source file called 'example.jpg' would be 'example'.
  -processes int
    	The number of concurrent processes to use when tiling images (default 10)
  -quality string
    	A valid IIIF quality parameter - if "default" then the code will try to determine which format you've set as the default (default "default")
  -refresh
    	Refresh a tile even if already exists (default false)
  -scale-factors string
    	A comma-separated list of scale factors to seed tiles with (default "8,4,2,1")
  -verbose
    	Enable verbose (debug) logging.
```

## Example

```
$> bin/iiif-tile-seed \
	-config-images-source-uri file:///usr/local/src/go-iiif/static/example/images \
	-config-derivatives-cache-uri file:///usr/local/src/go-iiif/work \
	-scale-factors '8,4,2,1' \
	-verbose \
	'rewrite:///spanking-cat.jpg?target=spank'
```
## URIs and identifiers

Identifiers for source images can be passed to `iiif-tiles-seed` in of two way:

1. A space-separated list of identifiers
2. A space-separated list of _comma-separated_ identifiers indicating the identifier for the source image followed by the identifier for the newly generated tiles

For example:

```
$> ./bin/iiif-tile-seed [options] 191733_5755a1309e4d66a7_k.jpg
```

Or:

```
$> ./bin/iiif-tile-seed [options] 191733_5755a1309e4d66a7_k.jpg,191/733/191733_5755a1309e4d66a7
```

In many cases the first option will suffice but sometimes you might need to create new identifiers or structure existing identifiers according to their output, for example avoiding the need to store lots of file in a single directory. It's up to you.

## CSV input

You can also run `iiif-tile-seed` pass a list of identifiers as a CSV file. To do so include the `-mode csv` argument, like this:

```
$> ./bin/iiif-tile-seed [options] -mode csv CSVFILE
```

CSV should have the following columns:

| Name | Required | Notes |
| --- | --- | --- |
| `source_filename` | yes | The name of the file to generate tiles for. |
| `source_root` | no | A custom source URI to read `source_filename` from. If empty then the default value from the config (`Images.Source.URI`) file will be used. |
| `target_filename` | no | The final name of the directory to store tiles in. If empty then the value of `source_filename` is used. |
| `target_root` | no | A custom target URI to write tiles to. If empty then the default value from the config file (`Derivatives.Cache.URI`) will be used. |

