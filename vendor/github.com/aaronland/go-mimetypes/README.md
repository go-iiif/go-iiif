# go-mimetypes

There are many mime-type lookup tables. This one is ours.

## Usage

```
package main

import (
	"flag"
	"fmt"
	"github.com/aaronland/go-mimetypes"
	"log"
	"os"
	"strings"
)

func main() {

	var extension = flag.Bool("extension", false, "Lookup mimetypes by extension")
	var mimetype = flag.Bool("mimetype", false, "Lookup extensions by mimetype")

	flag.Parse()

	for _, input := range flag.Args() {

		if *mimetype {
			t := mimetypes.TypesByExtension(input)
			fmt.Printf("%s\t%s\n", input, strings.Join(t, "\t"))
		} else if *extension {
			e := mimetypes.ExtensionsByType(input)
			fmt.Printf("%s\t%s\n", input, strings.Join(e, "\t"))
		} else {
			log.Fatal("Invalid lookup type")
		}
	}

	os.Exit(0)
}
```

Both the `TypesByExtension` and `ExtensionsByType` functions return a list of strings. Although it is unlikely that any given extension will return multiple mime-types it's certainly possible so that's what we solve. Handy helper functions to test for and return singletons are in the works.

When calling `TypesByExtension` you can pass extensions with or without a leading `.` character.

## Tools

### lookup

A simple command line tool to lookup mimetypes by extension or vice versa.

```
./bin/lookup -h
Usage of ./bin/wof-mimetype-lookup:
  -extension
    	Lookup mimetypes by extension
  -mimetype
    	Lookup extensions by mimetype
```

The output is a line-separated list containing a tab-separated list of input (extension or mimetype) followed by one or more matched. For example:

```
./bin/wof-mimetype-lookup -mimetype '.json' jpg
.json	application/json
jpg	image/jpeg

./bin/wof-mimetype-lookup -extension 'image/jpeg' 'image/gif'
image/jpeg	jpeg	jpg	jpe
image/gif	gif
```

## Lookup tables

Lookup tables are generated from two sources: The Apache `httpd` web servers mime.types definition and a custom file to add things that aren't included in the first source.

Both sources are bundled with this repo in the [tables](tables) subdirectory. To fetch a fresh version of the Apache `mime.typs` definitions run the `apache-lookup-table` Make target, like this:

```
make apache-lookup-table
curl -s -o tables/apache.txt 'https://svn.apache.org/viewvc/httpd/httpd/branches/2.2.x/docs/conf/mime.types?view=co'
```

To rebuild the (Go) lookup tables definitions run the `lookup-tables` Make target, like this:

```
make lookup-tables
if test -d pkg; then rm -rf pkg; fi
if test -d src; then rm -rf src; fi 
if test -d src; then rm -rf src; fi
mkdir -p src/github.com/whosonfirst/go-whosonfirst-mimetypes
cp *.go src/github.com/whosonfirst/go-whosonfirst-mimetypes/
cp -r lookup src/github.com/whosonfirst/go-whosonfirst-mimetypes/
if test -d vendor; then cp -r vendor/* src/; fi
if test -d lookup; then rm -rf lookup; fi
mkdir lookup
bin/wof-mimetype-build-lookup-tables -lookup extension tables/apache.txt tables/custom.txt > lookup/extension.go
bin/wof-mimetype-build-lookup-tables -lookup mimetype tables/apache.txt tables/custom.txt > lookup/mimetype.go
go fmt lookup/*.go
lookup/extension.go
lookup/mimetype.go
rm bin/wof-mimetype-build-lookup-tables
```

Note that the `lookup-tables` Make target does not rebuild the various binary tools (maybe it should) so you'll need to run the `bin` Make target to do that. There is a handy `rebuild` Make target do all of these things in go, if necessary, like this:

```
make rebuild
make[1]: Entering directory '/usr/local/whosonfirst/go-whosonfirst-mimetypes'
curl -s -o tables/apache.txt 'https://svn.apache.org/viewvc/httpd/httpd/branches/2.2.x/docs/conf/mime.types?view=co'
make[1]: Leaving directory '/usr/local/whosonfirst/go-whosonfirst-mimetypes'
make[1]: Entering directory '/usr/local/whosonfirst/go-whosonfirst-mimetypes'
if test -d pkg; then rm -rf pkg; fi
if test -d src; then rm -rf src; fi 
if test -d src; then rm -rf src; fi
mkdir -p src/github.com/whosonfirst/go-whosonfirst-mimetypes
cp *.go src/github.com/whosonfirst/go-whosonfirst-mimetypes/
cp -r lookup src/github.com/whosonfirst/go-whosonfirst-mimetypes/
if test -d vendor; then cp -r vendor/* src/; fi
if test -d lookup; then rm -rf lookup; fi
mkdir lookup
bin/wof-mimetype-build-lookup-tables -lookup extension tables/apache.txt tables/custom.txt > lookup/extension.go
bin/wof-mimetype-build-lookup-tables -lookup mimetype tables/apache.txt tables/custom.txt > lookup/mimetype.go
go fmt lookup/*.go
lookup/extension.go
lookup/mimetype.go
rm bin/wof-mimetype-build-lookup-tables
make[1]: Leaving directory '/usr/local/whosonfirst/go-whosonfirst-mimetypes'
make[1]: Entering directory '/usr/local/whosonfirst/go-whosonfirst-mimetypes'
if test -d pkg; then rm -rf pkg; fi
if test -d src; then rm -rf src; fi 
if test -d src; then rm -rf src; fi
mkdir -p src/github.com/whosonfirst/go-whosonfirst-mimetypes
cp *.go src/github.com/whosonfirst/go-whosonfirst-mimetypes/
cp -r lookup src/github.com/whosonfirst/go-whosonfirst-mimetypes/
if test -d vendor; then cp -r vendor/* src/; fi
```

## See also

* https://golang.org/pkg/mime/#ExtensionsByType
* https://golang.org/pkg/net/http/#DetectContentType
* https://svn.apache.org/viewvc/httpd/httpd/branches/2.2.x/docs/conf/mime.types?view=co
* https://www.iana.org/assignments/media-types/media-types.xhtml
