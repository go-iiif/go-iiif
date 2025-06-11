# Services

## Non-standard services

### palette

`go-iiif` supports the following additional services for profiles:

* `blurhash` for generateing a compact base-83 encoded representation of an image.
* `imagehash` for generating average and difference perceptual hashes of an image.
* `palette` for extracting one or more colour palettes for an image using the [aaronland/go-colours](https://github.com/aaronland/go-colours) package.

Details for configuring these service are discussed [above](#services) but here is the output for a service with the default settings:

```
$> curl -s localhost:8080/spanking-cat.jpg/info.json | jq '.service'
[
  {
    "@context": "x-urn:service:go-iiif#palette",
    "profile": "x-urn:service:go-iiif#palette",
    "label": "x-urn:service:go-iiif#palette",
    "palette": [
      {
        "name": "#dddad7",
        "hex": "#dddad7",
        "reference": "marekm4",
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
        "name": "#594a36",
        "hex": "#594a36",
        "reference": "marekm4",
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
        "name": "#9b8e6f",
        "hex": "#9b8e6f",
        "reference": "marekm4",
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
        "name": "#88775a",
        "hex": "#88775a",
        "reference": "marekm4",
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
```

_Please remember that the [aaronland/go-colours](https://github.com/aaronland/go-colours) package, used "under the hood", is itself is a work in progress so you should approach the `palette` service accordingly._

### Writing your own non-standard services

Services are invoked by the `go-iiif` codebase using URI-style identifiers. For example, assuming an "example" service you would invoke it like this:

```
    	service_name := "example"	
	service_uri := fmt.Sprintf("%s://", service_name)
	service, _ := iiifservice.NewService(ctx, service_uri, cfg, im)
```

In addition to implementing the `service.Service` interface custom services need to also "register" themselves on initialization with a (golang) context, a (go-iiif), a unique scheme used to identify the service and a `service.ServiceInitializationFunc` callback function. The callback function implements the following interface:

```
type ServiceInitializationFunc func(ctx context.Context, config *iiifconfig.Config, im iiifimage.Image) (Service, error)
```

Here is an abbreviated example, with error handling removed for the sake of brevity. For real working examples, take a look at any of the built-in services in the [services](services) directory.

```
package example	// for example "github.com/example/go-iiif-example"

import (
	"context"
	
	iiifconfig "github.com/go-iiif/go-iiif/v8/config"
	iiifimage "github.com/go-iiif/go-iiif/v8/image"	
	iiifservice "github.com/go-iiif/go-iiif/v8/service"	
)

func init() {
	ctx := context.Background()
	iiifservice.RegisterService(ctx, "example", initExampleService)
}

func initExampleService(ctx context.Context, cfg *iiifconfig.Config, im iiifimage.Image) (iiifservice.Service, error) {
	return NewExampleService(cfg, im)
}

type ExampleService struct {
	iiifservice.Service        `json:",omitempty"`
	// your properties here...
}

// your implementation of the iiifservice.Service interface here...

func NewExampleService(cfg *iiifconfig.Config, im iiifimage.Image) (iiifservice.Service, error){

     // presumably you'd do something with im here...
     
     s := &ExampleService{
       // your properties here...
     }

     return s, nil
}
```

Finally, you will need to create custom versions of any `go-iiif` tools you want to you use your new service. For example, here's a modified version of the [cmd/iiif-server/main.go](cmd/iiif-server/main.go) server implementation.

```
package main

import (
)

import (
	"context"
	"log"
	
        _ "github.com/example/go-iiif-example"
	_ "github.com/aaronland/go-cloud-s3blob"
	_ "github.com/go-iiif/go-iiif/v8/native"
	_ "gocloud.dev/blob/fileblob"       

	"github.com/go-iiif/go-iiif/v8/tools"
)

func main() {

	tool, err := tools.NewIIIFServerTool()

	if err != nil {
		log.Fatal(err)
	}

	err = tool.Run(context.Background())

	if err != nil {
		log.Fatal(err)
	}
}
```

 The only change from the default server tool is the addition of the `_ "github.com/example/go-iiif-example"` import statement. That will allow the core `go-iiif` software to find and use your custom service.

It's unfortunate that using custom and bespoke services requires compiling your own version of the `go-iiif` tools but such is life when you are using a language like Go.
