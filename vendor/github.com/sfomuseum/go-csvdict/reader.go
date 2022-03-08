package csvdict

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

// type Reader implements a `encoding/csv` style reader for CSV documents with named columns.
type Reader struct {
	Reader     *csv.Reader
	Fieldnames []string
}

// NewReader will return a Reader instance that will load data from 'r'
func NewReader(r io.Reader) (*Reader, error) {

	reader := csv.NewReader(r)

	row, read_err := reader.Read()

	if read_err != nil {
		return nil, fmt.Errorf("Failed to read first line of document, %w", read_err)
	}

	dr := Reader{Reader: reader, Fieldnames: row}
	return &dr, nil
}

// NewReader will return a Reader instance that will load data from 'path'
func NewReaderFromPath(path string) (*Reader, error) {

	fh, err := os.Open(path)

	if err != nil {
		return nil, fmt.Errorf("Failed to open %s, %w", path, err)
	}

	return NewReader(fh)
}

// Read reads one record (a slice of fields) from r and returns a map[string]string
// mapping columns to their corresponding names, as defined in the first line of r.
func (dr Reader) Read() (map[string]string, error) {

	row, err := dr.Reader.Read()

	if err != nil {
		return nil, err
	}

	dict := make(map[string]string)

	for i, value := range row {
		key := dr.Fieldnames[i]
		dict[key] = value
	}

	return dict, nil
}
