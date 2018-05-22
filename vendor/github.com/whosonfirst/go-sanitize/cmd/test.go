package main

// This file will be deleted when the dust settles (20160909/thisisaaronland)

import (
	"fmt"
	"os"
	"sanitize"
	"strconv"
	"strings"
	_ "unicode/utf8"
)

func main() {

	opts := sanitize.DebugOptions()
	// opts.StripReserved = true
	
	total := 0
	ok := 0
	fail := 0

	// see below...

	codepoints := [][]string{
		{"00000000", "00001000"},
		{"00001110", "00011111"},
		{"01111111", "10000100"},
		{"10000110", "10011111"},
		{"1111111011111111", "1111111011111111"},
		{"1111111111111001", "1111111111111010"},
		{"11100000000000000000", "11100000000001111111"},
		{"1101100000000000", "1101111111111111"},
		{"100010000000000000000", "100111111111111111111"},
	}

	for _, pair := range codepoints {

		// I studied painting... leave me alone or just be
		// gentle when you show me how to do this the right
		// way... (20160909/thisisaaronland)

		lo, _ := strconv.ParseUint(pair[0], 2, 32)
		hi, _ := strconv.ParseUint(pair[1], 2, 32)

		// fmt.Println("#", lo, hi)

		for i := lo; i < hi; i++ {

			total += 1

			r := rune(i)
			c, _ := sanitize.SanitizeString(string(r), opts)

			if c == " { SANITIZED } " {
				ok += 1
			} else {

				fail += 1

				h := make([]string, 0)
				s := string(r)

				for j := 0; j < len(s); j++ {
					h = append(h, fmt.Sprintf("% x", s[j]))
				}

				fmt.Printf("[%s-%s] %b %U %+q %s\n", pair[0], pair[1], r, r, r, strings.Join(h, " "))
			}
		}
	}

	/*
		s := "\xf4\x8f"
		r, _ := utf8.DecodeRuneInString(s)
		c, _ := sanitize.SanitizeString(string(r), opts)
		fmt.Println("WUH", c)
	*/

	fmt.Printf("total: %d ok: %d fail: %d\n", total, ok, fail)

	input := "foo bar\nbaz ok: '\u2318' BAD:'\u0007' BAD:'\uFEFF' BAD:'\u2029' BAD:'\u0007' BAD:'\u007F' BAD:'\x0B' woop wwoop"

	/*
		for index, runeValue := range input {

			s, _ := sanitize.SanitizeString(string(runeValue), opts)
			fmt.Printf("%#U starts at byte position %d becomes '%s'\n", runeValue, index, s)
		}
	*/

	output, _ := sanitize.SanitizeString(input, opts)
	fmt.Println(output)

	os.Exit(0)
}
