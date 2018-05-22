package main

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-csv"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func main() {

	var str_cols = flag.String("columns", "-", "Columns to filter on. A value of \"-\" means the set of unique columns for all CSV files being filtered.")
	var out = flag.String("out", "-", "Where to write the data. A value of \"-\" means write the data to STDOUT.")

	flag.Usage = func() {
		whoami := os.Args[0]
		fname := filepath.Base(whoami)

		fmt.Fprintf(os.Stderr, "Usage of %s:\n", fname)
		fmt.Printf("  $> %s -options <files>\n\n", fname)

		fmt.Printf("Valid options are:\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	files := flag.Args()

	// see also: https://github.com/whosonfirst/go-whosonfirst-csv/issues/2

	cols := make([]string, 0)

	if *str_cols == "-" {

		tmp := make(map[string]int)

		for _, path := range files {

			reader, err := csv.NewDictReaderFromPath(path)

			if err != nil {
				continue
			}

			for _, k := range reader.Fieldnames {
				tmp[k] = 1
			}
		}

		for k, _ := range tmp {
			cols = append(cols, k)
		}
	} else {

		cols = strings.Split(*str_cols, ",")
	}

	if len(cols) == 0 {
		panic("NO COLUMNS")
	}

	var writer *csv.DictWriter

	if *out == "-" {

		w, err := csv.NewDictWriter(os.Stdout, cols)

		if err != nil {
			panic(err)
		}

		writer = w

	} else {

		w, err := csv.NewDictWriterFromPath(*out, cols)

		if err != nil {
			panic(err)
		}

		writer = w
	}

	writer.WriteHeader()

	wg := new(sync.WaitGroup)
	mu := new(sync.Mutex)

	for _, path := range files {

		wg.Add(1)

		go func(path string, cols []string) {

			defer wg.Done()

			reader, err := csv.NewDictReaderFromPath(path)

			if err != nil {
				return
			}

			for {

				row, err := reader.Read()

				if err == io.EOF {
					break
				}

				if err != nil {
					continue
				}

				mu.Lock()

				writer.WriteRow(row)
				mu.Unlock()

			}

		}(path, cols)

	}

	wg.Wait()

}
