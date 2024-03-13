// package ioutil provides methods for creating a new instance conforming to the Go 1.16 io.ReadSeekCloser interface from a variety of io.Read* instances that implement some but not all of the io.Reader, io.Seeker and io.Closer interfaces.
//
// Example
//
//	import (
//		"bytes"
//		"github.com/whosonfirst/go-ioutil"
//		"io"
//		"log"
//	)
//
//	func main(){
//
//		fh, _ := os.Open("README.md")
//
//		rsc, _ := NewReadSeekCloser(fh)
//
//		body, _ := io.ReadAll(rsc)
//
//		rsc.Seek(0, 0)
//
//		body2, _ := io.ReadAll(rsc)
//
//		same := bytes.Equal(body, body2)
//		log.Printf("Same %t\n", same)
//
//		rsc.Close()
//	}
package ioutil
