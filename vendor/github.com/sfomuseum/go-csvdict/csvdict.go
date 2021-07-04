package csvdict

import (
	"encoding/csv"
	"io"
	"os"
)

type Reader struct {
	Reader     *csv.Reader
	Fieldnames []string
}

type Writer struct {
	Writer     *csv.Writer
	Fieldnames []string
}

func NewReader(fh io.Reader) (*Reader, error) {

	reader := csv.NewReader(fh)

	row, read_err := reader.Read()

	if read_err != nil {
		return nil, read_err
	}

	dr := Reader{Reader: reader, Fieldnames: row}
	return &dr, nil
}

func NewReaderFromPath(path string) (*Reader, error) {

	fh, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	return NewReader(fh)
}

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

func NewWriter(fh io.Writer, fieldnames []string) (*Writer, error) {

	writer := csv.NewWriter(fh)

	dw := Writer{Writer: writer, Fieldnames: fieldnames}
	return &dw, nil
}

func NewWriterFromPath(path string, fieldnames []string) (*Writer, error) {

	fh, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	if err != nil {
		return nil, err
	}

	return NewWriter(fh, fieldnames)
}

func (dw Writer) WriteHeader() error {

	return dw.Writer.Write(dw.Fieldnames)
}

// to do - check flags for whether or not to be liberal when missing keys
// (20160516/thisisaaronland)

func (dw Writer) WriteRow(row map[string]string) error {

	out := make([]string, 0)

	for _, k := range dw.Fieldnames {

		v, ok := row[k]

		if !ok {
			v = ""
		}

		out = append(out, v)
	}

	return dw.Writer.Write(out)
}

// Flush writes any buffered data to the underlying writer. To check if an error occurred during the Flush, call Error.
func (dw Writer) Flush() error {
	dw.Writer.Flush()
	return nil
}

// Error reports any error that has occurred during a previous Write or Flush.
func (dw Writer) Error() error {
	return dw.Writer.Error()
}
