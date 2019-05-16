# go-sanitize

A Go port of the Flamework lib_sanitize library.

## Install

You will need to have both `Go` (specifically [version 1.12](https://golang.org/dl/) or higher because we're using [Go modules](https://github.com/golang/go/wiki/Modules)) and the `make` programs installed on your computer. 

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Example

```
package main

import (
       "fmt"
       "sanitize"
       )

func main() {

     input := "foo bar\nbaz ok:'\u2318' bad:'\u0007' bad:'\u206E' woop woop"
     
     opts := sanitize.DefaultOptions()
     output, _ := sanitize.SanitizeString(input, opts)

     fmt.Println(output)

}

```     

## See also

* https://github.com/exflickr/flamework/blob/master/www/include/lib_sanitize.php
