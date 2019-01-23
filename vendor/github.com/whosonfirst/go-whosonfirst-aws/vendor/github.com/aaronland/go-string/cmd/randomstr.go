package main

import (
	"flag"
	"fmt"
	"github.com/aaronland/go-string/random"
	"log"
)

func main() {

	var ascii = flag.Bool("ascii", false, "Only include ASCII characters")
	var alpha = flag.Bool("alphanumeric", false, "Only include alpha-numeric characters (this causes the -ascii flag to be set to true)")
	var length = flag.Int("length", 32, "Minimum length of the random string, in bytes")
	var chars = flag.Int("chars", 0, "Minimum length of the random string, in characters")
	var base32 = flag.Bool("base32", false, "Encode as a base32 string")

	flag.Parse()

	if *alpha {
		*ascii = true
	}

	opts := random.DefaultOptions()
	opts.ASCII = *ascii
	opts.AlphaNumeric = *alpha
	opts.Length = *length
	opts.Chars = *chars
	opts.Base32 = *base32

	s, err := random.String(opts)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(s)
}
