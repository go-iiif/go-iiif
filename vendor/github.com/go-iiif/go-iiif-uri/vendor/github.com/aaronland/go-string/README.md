# go-string

There are many string libraries. This one is mine.

## Install

You will need to have both `Go` (specifically a version of Go more recent than 1.6 so let's just assume you need [Go 1.8](https://golang.org/dl/) or higher) and the `make` programs installed on your computer. Assuming you do just type:

```
make bin
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Example

```
package main

import (
	"fmt"
	"github.com/aaronland/go-string/random"
)

func main() {

	opts := random.DefaultOptions()
	opts.Length = 40
	opts.Chars = 20

	s, _ := random.String(opts)

	fmt.Println(s)
}
```

_Error handling omitted for the sake of brevity._

## Tools

### randomstr

```
./bin/randomstr -h
Usage of ./bin/randomstr:
  -alphanumeric
    	Only include alpha-numeric characters (this causes the -ascii flag to be set to true)
  -ascii
    	Only include ASCII characters
  -chars int
    	Minimum length of the random string, in characters
  -length int
    	Minimum length of the random string, in bytes (default 32)
```	