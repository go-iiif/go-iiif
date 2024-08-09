package funcs

import (
	"time"
)

func FormatStringTime(strtime string, input_format string, output_format string) string {

	t, err := time.Parse(input_format, strtime)

	if err != nil {
		return strtime
	}

	return t.Format(output_format)
}

func FormatUnixTime(ts int64, output_format string) string {

	t := time.Unix(ts, 0).UTC()

	return t.Format(output_format)
}
