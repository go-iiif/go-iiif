package csv

import (
	gocsv "encoding/csv"
	"io"
	"os"
)

type DictReader struct {
	Reader     *gocsv.Reader
	Fieldnames []string
}

type DictWriter struct {
	Writer     *gocsv.Writer
	Fieldnames []string
}

func NewDictReader(fh io.Reader) (*DictReader, error) {

	reader := gocsv.NewReader(fh)

	row, read_err := reader.Read()

	if read_err != nil {
		return nil, read_err
	}

	dr := DictReader{Reader: reader, Fieldnames: row}
	return &dr, nil
}

func NewDictReaderFromPath(path string) (*DictReader, error) {

	fh, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	return NewDictReader(fh)
}

func (dr DictReader) Read() (map[string]string, error) {

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

func NewDictWriter(fh io.Writer, fieldnames []string) (*DictWriter, error) {

	writer := gocsv.NewWriter(fh)

	dw := DictWriter{Writer: writer, Fieldnames: fieldnames}
	return &dw, nil
}

func NewDictWriterFromPath(path string, fieldnames []string) (*DictWriter, error) {

	fh, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	if err != nil {
		return nil, err
	}

	return NewDictWriter(fh, fieldnames)
}

func (dw DictWriter) WriteHeader() {
	dw.Writer.Write(dw.Fieldnames)
	dw.Writer.Flush()
}

// to do - check flags for whether or not to be liberal when missing keys
// (20160516/thisisaaronland)

func (dw DictWriter) WriteRow(row map[string]string) {

	out := make([]string, 0)

	for _, k := range dw.Fieldnames {

		v, ok := row[k]

		if !ok {
			v = ""
		}

		out = append(out, v)
	}

	dw.Writer.Write(out)
	dw.Writer.Flush() // move me somewhere more sensible
}
