# iiif-transform

Transform one or more images using the IIIF API.

```
$> ./bin/iiif-transform -h
Transform one or more images using the IIIF API.

Usage:
	 ./bin/iiif-transform[options] uri(N) uri(N)

Valid options are:
  -config-derivatives-cache-uri string
    	If present this value will be used to assign the 'derivatives.cache.uri' property in the config file. Note: The 'derivatives.cache.uri' property takes precedence over other properties in 'derivatives.cache' block.
  -config-images-source-uri string
    	If present this value will be used to assign the 'images.source.uri' property in the config file. Note: The 'images.source.uri' property takes precedence over other properties in 'images.source' block.
  -config-name string
    	The name of your go-iiif config file. This value will be ignored if -config-source is 'defaults://'. (default "config.json")
  -config-source string
    	A valid Go Cloud bucket URI where your go-iiif config file is located. Optionally, if 'defaults://' is specified then the default config bundled with this package will be used. (default "defaults://")
  -format string
    	A valid IIIF 2.0 format value. (default "jpg")
  -mode string
    	Valid options are: cli, fsnotify, lambda (default "cli")
  -quality string
    	A valid IIIF 2.0 quality value. (default "default")
  -region string
    	A valid IIIF 2.0 region value. (default "full")
  -rotation string
    	A valid IIIF 2.0 rotation value. (default "0")
  -size string
    	A valid IIIF 2.0 size value. (default "full")
  -verbose
    	Enabled verbose (debug) loggging.
```	
