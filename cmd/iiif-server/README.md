# iiif-server

Expose the IIIF Image API via an HTTP endpoint.

```
$> ./bin/iiif-server -h
Expose the IIIF Image API via an HTTP endpoint.

Usage:
	 ./bin/iiif-server[options]

Valid options are:
  -config-derivatives-cache-uri string
    	If present this value will be used to assign the 'derivatives.cache.uri' property in the config file. Note: The 'derivatives.cache.uri' property takes precedence over other properties in 'derivatives.cache' block.
  -config-images-source-uri string
    	If present this value will be used to assign the 'images.source.uri' property in the config file. Note: The 'images.source.uri' property takes precedence over other properties in 'images.source' block.
  -config-name string
    	The name of your go-iiif config file. This value will be ignored if -config-source is 'defaults://'. (default "config.json")
  -config-source string
    	A valid Go Cloud bucket URI where your go-iiif config file is located. Optionally, if 'defaults://' is specified then the default config bundled with this package will be used. (default "defaults://")
  -example
    	Add an example endpoint (at /) to the server for testing and demonstration purposes.
  -server-uri string
    	A registered aaronland/go-http-server.Server URI. (default "http://localhost:8080")
  -verbose
    	Enabled verbose (debug) loggging.
```

## Example

```
$> bin/iiif-server \
	-config-derivatives-cache-uri file:///usr/local/src/go-iiif/fixtures/cache \
	-config-images-source-uri file:///usr/local/src/go-iiif/fixtures/images \
	-verbose
2025/03/25 14:15:08 DEBUG Verbose logging enabled
2025/03/25 14:15:08 INFO Listening for requests address=http://localhost:8080
```

## Endpoints

Although the identifier parameter (`{ID}`) in the examples below suggests that is is only string characters up to and until a `/` character, it can in fact contain multiple `/` separated strings. For example, either of these two URLs is valid

```
http://localhost:8080/191733_5755a1309e4d66a7_k.jpg/info.json
http://localhost:8080/191/733/191733_5755a1309e4d66a7/info.json
```

Where the identified will be interpreted as `191733_5755a1309e4d66a7_k.jpg` and `191/733/191733_5755a1309e4d66a7` respectively. Identifiers containing one or more `../` strings will be made to feel bad about themselves.

### GET /{ID}/info.json

```
$> curl -s http://localhost:8080/spanking-cat.jpg/info.json | python -mjson.tool
{
  "height": 4096,
  "width": 3897,
  "@context": "http://iiif.io/api/image/2/context.json",
  "@id": "http://localhost:8080/spanking-cat.jpg",
  "protocol": "http://iiif.io/api/image",
  "profile": [
    "http://iiif.io/api/image/2/level2.json",
    {
      "formats": [
        "webp",
        "jpg",
        "png",
        "tif",
        "gif"
      ],
      "qualities": [
        "dither",
        "default",
        "color",
        "crisp"
      ],
      "supports": [
        "regionSquare",
        "full",
        "regionByPx",
        "regionByPct",
        "sizeByPct",
        "sizeByDistortedWh",
        "full",
        "max",
        "sizeByConfinedWh",
        "sizeByWh",
        "sizeByW",
        "sizeByH",
        "mirroring",
        "noAutoRotate",
        "none",
        "rotationBy90s",
        "cors",
        "jsonldMediaType",
        "baseUriRedirect"
      ]
    }
  ],
  "service": [
    {
      "@context": "x-urn:service:go-iiif#palette",
      "profile": "x-urn:service:go-iiif#palette",
      "label": "x-urn:service:go-iiif#palette",
      "palette": [
        {
          "name": "6b5c43",
          "hex": "#6b5c43",
          "reference": "vibrant"
        },
        {
          "name": "cfcbc9",
          "hex": "#cfcbc9",
          "reference": "vibrant"
        },
        {
          "name": "a69f92",
          "hex": "#a69f92",
          "reference": "vibrant"
        },
        {
          "name": "e7e5e0",
          "hex": "#e7e5e0",
          "reference": "vibrant"
        },
        {
          "name": "dad8d7",
          "hex": "#dad8d7",
          "reference": "vibrant"
        }
      ]
    },
    {
      "@context": "x-urn:service:go-iiif#blurhash",
      "profile": "x-urn:service:go-iiif#blurhash",
      "label": "x-urn:service:go-iiif#blurhash",
      "hash": ":NOWsZa{_Nt69Fofxuof_2a{M{ofofaej[ay%Nj[IAj[%MayRjj[jZoft7V@WCofaxj[xuoMRjV@j[t7ofWBoeayfRkCayaxofj[WBaxofozayRjaxofRjayt7ofWBayj?ay"
    },
    {
      "@context": "x-urn:service:go-iiif#imagehash",
      "profile": "x-urn:service:go-iiif#imagehash",
      "label": "x-urn:service:go-iiif#imagehash",
      "average": "a:ffffc7e7c3c3c3c3",
      "difference": "d:c48c0c0e8e8f0e0f"
    }
  ]
}
```

Return the [profile description](http://iiif.io/api/image/2.1/#profile-description) for an identifier.

### GET /{ID}/{REGION}/{SIZE}/{ROTATION}/{QUALITY}.{FORMAT}

```
$> curl -s http://localhost:8080/spanking-cat.jpg/pct:41,7,40,70/,250/0/default.jpg
```

Return an image derived from an identifier and one or more [IIIF parameters](http://iiif.io/api/image/2.1/#image-request-parameters). For example:

![spanking cat, cropped](../../misc/go-iiif-crop.jpg)