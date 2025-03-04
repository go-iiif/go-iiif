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

