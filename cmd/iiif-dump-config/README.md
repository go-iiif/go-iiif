# iiif-dump-config

Emit a go-iiif config file as Markdown.

```
$> ./bin/iiif-dump-config -h
Emit a go-iiif config file as Markdown.

Usage:
	 ./bin/iiif-dump-config[options]

Valid options are:
  -config-derivatives-cache-uri string
    	If present this value will be used to assign the 'derivatives.cache.uri' property in the config file. Note: The 'derivatives.cache.uri' property takes precedence over other properties in 'derivatives.cache' block.
  -config-images-source-uri string
    	If present this value will be used to assign the 'images.source.uri' property in the config file. Note: The 'images.source.uri' property takes precedence over other properties in 'images.source' block.
  -config-name string
    	The name of your go-iiif config file. This value will be ignored if -config-source is 'defaults://'. (default "config.json")
  -config-source string
    	A valid Go Cloud bucket URI where your go-iiif config file is located. Optionally, if 'defaults://' is specified then the default config bundled with this package will be used. (default "defaults://")
  -verbose
    	Enabled verbose (debug) loggging.
```	