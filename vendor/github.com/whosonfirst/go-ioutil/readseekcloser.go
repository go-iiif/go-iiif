package ioutil

// This is only here until there is an equivalent package/construct in the core Go language
// (20210217/thisisaaronland)

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

// Type ReadSeekCloser implements the io.Reader, io.Seeker and io.Closer interfaces.
type ReadSeekCloser struct {
	io.Reader
	io.Seeker
	io.Closer
	reader bool
	closer bool
	seeker bool
	fh     interface{}
	br     *bytes.Reader
	mu     *sync.RWMutex
}

// Create a new NewReadSeekCloser instance conforming to the Go 1.16 `io.ReadSeekCloser` interface. This method accepts the following types: io.ReadSeekCloser, io.Reader, io.ReadCloser and io.ReadSeeker.
func NewReadSeekCloser(fh interface{}) (io.ReadSeekCloser, error) {

	reader := false
	seeker := false
	closer := false

	switch fh.(type) {
	case io.ReadSeekCloser:
		return fh.(io.ReadSeekCloser), nil
	case io.Closer:
		closer = true
	case io.ReadCloser:
		reader = true
		closer = true
	case io.ReadSeeker:
		reader = true
		seeker = true
	case io.Reader:
		reader = true
	default:
		return nil, fmt.Errorf("Invalid or unsupported type: %T", fh)
	}

	mu := new(sync.RWMutex)

	rsc := &ReadSeekCloser{
		reader: reader,
		seeker: seeker,
		closer: closer,
		fh:     fh,
		mu:     mu,
	}

	return rsc, nil
}

// Read implements the standard Read interface: it reads data from the pipe, blocking until a writer arrives or the write end is closed. If the write end is closed with an error, that error is returned as err; otherwise err is `io.EOF`.
func (rsc *ReadSeekCloser) Read(p []byte) (n int, err error) {

	if rsc.seeker {
		return rsc.fh.(io.Reader).Read(p)
	}

	br, err := rsc.bytesReader()

	if err != nil {
		return 0, err
	}

	return br.Read(p)
}

// Close closes the reader; subsequent writes to the write half of the pipe will return the error `io.ErrClosedPipe`.
func (rsc *ReadSeekCloser) Close() error {

	if rsc.closer {
		return rsc.fh.(io.ReadCloser).Close()
	}

	return nil
}

// Seek implements the `io.Seeker` interface.
func (rsc *ReadSeekCloser) Seek(offset int64, whence int) (int64, error) {

	if rsc.seeker {
		return rsc.fh.(io.Seeker).Seek(offset, whence)
	}

	br, err := rsc.bytesReader()

	if err != nil {
		return 0, err
	}

	return br.Seek(offset, whence)
}

func (rsc *ReadSeekCloser) bytesReader() (*bytes.Reader, error) {

	rsc.mu.Lock()
	defer rsc.mu.Unlock()

	if rsc.br != nil {
		return rsc.br, nil
	}

	body, err := io.ReadAll(rsc.fh.(io.Reader))

	if err != nil {
		return nil, err
	}

	br := bytes.NewReader(body)
	rsc.br = br

	return br, nil
}
