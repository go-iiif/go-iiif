package csvdict

import (
	"encoding/csv"
	"fmt"
	"io"
	"iter"
	"os"
	"path/filepath"
)

// type Reader implements a `encoding/csv` style reader for CSV documents with named columns.
type Reader struct {
	csv_reader *csv.Reader
	fieldnames []string
}

// NewReader will return a Reader instance that will load data from 'path'
func NewReaderFromPath(path string) (*Reader, error) {

	abs_path, err := filepath.Abs(path)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive absolute path for %s, %w", path, err)
	}

	r, err := os.Open(abs_path)

	if err != nil {
		return nil, fmt.Errorf("Failed to open %s, %w", path, err)
	}

	return NewReader(r)
}

// NewReader will return a Reader instance that will load data from 'r'
func NewReader(io_r io.Reader) (*Reader, error) {

	csv_r := csv.NewReader(io_r)

	fieldnames, err := csv_r.Read()

	if err != nil {
		return nil, err
	}

	r := &Reader{
		csv_reader: csv_r,
		fieldnames: fieldnames,
	}

	return r, nil
}

// Read reads one record (a slice of fields) from r and returns a map[string]string
// mapping columns to their corresponding names, as defined in the first line of r.
func (r Reader) Read() (map[string]string, error) {

	row, err := r.csv_reader.Read()

	if err != nil {
		return nil, err
	}

	dict := make(map[string]string)

	for i, value := range row {
		key := r.fieldnames[i]
		dict[key] = value
	}

	return dict, nil
}

// Iterate throush all from r and yielding a map[string]string mapping columns to their corresponding
// names, as defined in the first line of r.
func (r Reader) Iterate() iter.Seq2[map[string]string, error] {

	return func(yield func(map[string]string, error) bool) {

		for {

			dict, err := r.Read()

			if err == io.EOF {
				break
			}

			if !yield(dict, nil) {
				return
			}
		}
	}
}
