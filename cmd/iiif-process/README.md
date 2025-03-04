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
