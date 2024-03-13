# go-ioutil

Go package for creating instances conforming to the Go 1.16 `io.ReadSeekCloser` interface from a variety of io.Read* instances that implement some but not all of the `io.Reader`, `io.Seeker` and `io.Closer` interfaces.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/whosonfirst/go-ioutil.svg)](https://pkg.go.dev/github.com/whosonfirst/go-ioutil)

## Example

```
import (
	"bytes"
	"github.com/whosonfirst/go-ioutil"
	"io"
	"log"
)

func main(){

	fh, _ := os.Open("README.md")
	
	rsc, _ := NewReadSeekCloser(fh)

	body, _ := io.ReadAll(rsc)

	rsc.Seek(0, 0)

	body2, _ := io.ReadAll(rsc)

	same := bytes.Equal(body, body2)
	log.Printf("Same %t\n", same)

	rsc.Close()
}	
```

_Error handling removed for brevity._

## See also

* https://golang.org/pkg/io/#ReadSeekCloser