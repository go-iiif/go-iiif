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
        "jpg",
        "png",
        "tif",
        "gif",
        "webp"
      ],
      "qualities": [
        "dither",
        "crisp",
        "default",
        "color"
      ],
      "supports": [
        "full",
        "regionByPx",
        "regionByPct",
        "regionSquare",
        "full",
        "max",
        "sizeByPct",
        "sizeByWh",
        "sizeByW",
        "sizeByH",
        "sizeByConfinedWh",
        "sizeByDistortedWh",
        "none",
        "rotationBy90s",
        "mirroring",
        "noAutoRotate",
        "baseUriRedirect",
        "cors",
        "jsonldMediaType"
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
          "name": "marekm4",
          "hex": "#dddad7",
          "reference": "dddad7",
          "closest": [
            {
              "name": "Timberwolf",
              "hex": "#dbd7d2",
              "reference": "crayola"
            },
            {
              "name": "gainsboro",
              "hex": "#dcdcdc",
              "reference": "css4"
            }
          ]
        },
        {
          "name": "marekm4",
          "hex": "#594a36",
          "reference": "594a36",
          "closest": [
            {
              "name": "Raw Umber",
              "hex": "#714b23",
              "reference": "crayola"
            },
            {
              "name": "darkolivegreen",
              "hex": "#556b2f",
              "reference": "css4"
            }
          ]
        },
        {
          "name": "marekm4",
          "hex": "#9b8e6f",
          "reference": "9b8e6f",
          "closest": [
            {
              "name": "Beaver",
              "hex": "#9f8170",
              "reference": "crayola"
            },
            {
              "name": "grey",
              "hex": "#808080",
              "reference": "css4"
            }
          ]
        },
        {
          "name": "marekm4",
          "hex": "#88775a",
          "reference": "88775a",
          "closest": [
            {
              "name": "Shadow",
              "hex": "#8a795d",
              "reference": "crayola"
            },
            {
              "name": "dimgrey",
              "hex": "#696969",
              "reference": "css4"
            }
          ]
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