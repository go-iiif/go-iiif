package csvdict

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

// type Writer implements a `encoding/csv` style writer for CSV documents with named columns.
type Writer struct {
	csv_writer *csv.Writer
	fieldnames []string
	mu         *sync.RWMutex
}

// NewWriter will return a new Writer that writes to 'wr'.
func NewWriter(wr io.Writer) (*Writer, error) {

	csv_writer := csv.NewWriter(wr)

	mu := new(sync.RWMutex)

	dw := Writer{
		csv_writer: csv_writer,
		mu:         mu,
	}

	return &dw, nil
}

// NewWriter will return a new Writer that writes to 'path'.
func NewWriterFromPath(path string) (*Writer, error) {

	abs_path, err := filepath.Abs(path)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive absolute path for %s, %w", path, err)
	}

	wr, err := os.OpenFile(abs_path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	if err != nil {
		return nil, fmt.Errorf("Failed to open %s for writing, %w", path, err)
	}

	return NewWriter(wr)
}

// WriteRow writes the values of row as CSV-encoded data. The order of those values is determined
// by their position defined in the list of keys in 'row' which are sorted alphabetically. WriteRow
// is thread-safe.
func (dw *Writer) WriteRow(row map[string]string) error {

	dw.mu.Lock()
	defer dw.mu.Unlock()

	if len(dw.fieldnames) == 0 {

		fieldnames := make([]string, 0)

		for k, _ := range row {
			fieldnames = append(fieldnames, k)
		}

		sort.Strings(fieldnames)
		err := dw.csv_writer.Write(fieldnames)

		if err != nil {
			return err
		}

		dw.fieldnames = fieldnames
	}

	out := make([]string, len(dw.fieldnames))

	for idx, k := range dw.fieldnames {

		v, ok := row[k]

		if !ok {
			v = ""
		}

		out[idx] = v
	}

	return dw.csv_writer.Write(out)
}

// Flush writes any buffered data to the underlying writer. To check if an error occurred during the Flush, call Error.
func (dw *Writer) Flush() error {
	dw.csv_writer.Flush()
	return nil
}

// Error reports any error that has occurred during a previous Write or Flush.
func (dw *Writer) Error() error {
	return dw.csv_writer.Error()
}
