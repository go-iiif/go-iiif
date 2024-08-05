package random

import (
	"encoding/base32"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode"
)

const min_length int = 32

var runes []rune

var r *rand.Rand

func init() {

	r = rand.New(rand.NewSource(time.Now().UnixNano()))

	runes = make([]rune, 0)

	codepoints := [][]int{
		[]int{1, 255},         // ascii
		[]int{127744, 128317}, // emoji
	}

	for _, r := range codepoints {

		first := r[0]
		last := r[1]

		for i := first; i < last; i++ {

			r := rune(i)

			if unicode.IsControl(r) {
				continue
			}

			if unicode.IsSpace(r) {
				continue
			}

			if unicode.IsMark(r) {
				continue
			}

			runes = append(runes, r)
		}
	}

}

type Options struct {
	Length       int
	Chars        int
	ASCII        bool
	AlphaNumeric bool
	Base32       bool
}

func DefaultOptions() *Options {

	opts := Options{
		Length:       min_length,
		Chars:        0,
		ASCII:        false,
		AlphaNumeric: false,
		Base32:       false,
	}

	return &opts
}

func String(opts *Options) (string, error) {

	count := len(runes)

	result := make([]string, 0)

	var last string

	// chars := 0
	b := 0

	alpha_numeric := [][]int{
		[]int{48, 57},  // (0-9)
		[]int{65, 90},  // (A-Z)
		[]int{97, 122}, // (a-z)
	}

	for b < opts.Length {

		j := r.Intn(count)
		r := runes[j]

		if opts.ASCII && r > 127 {
			continue
		}

		if opts.AlphaNumeric {

			is_alpha_numeric := false

			for _, bookends := range alpha_numeric {

				r_int := int(r)

				if r_int >= bookends[0] && r_int <= bookends[1] {
					is_alpha_numeric = true
					break
				}
			}

			if !is_alpha_numeric {
				continue
			}

		}

		c := fmt.Sprintf("%c", r)

		if c == last {
			continue
		}

		last = c

		b += len(c)

		if b <= opts.Length {
			result = append(result, c)
		} else {

			if len(result) > 2 {
				result = result[0 : len(result)-2]
			} else {
				result = make([]string, 0)
			}
			b = len(strings.Join(result, ""))
		}
	}

	s := strings.Join(result, "")

	if opts.Base32 {
		s = base32.StdEncoding.EncodeToString([]byte(s))
	}

	return s, nil
}
