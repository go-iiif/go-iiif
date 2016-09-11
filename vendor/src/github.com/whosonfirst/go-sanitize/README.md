# go-sanitize

A Go port of the Flamework lib_sanitize library.

## Important

YOU SHOULD NOT TRY TO USE THIS YET. IT HAS NOT BEEN FINISHED BEING TESTED YET. IT IS NOT SAFE. WELL IT _MIGHT_ BE SAFE BUT IT WOULD PROBABLY BE SAFER TO ASSUME IT ISN'T.

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
