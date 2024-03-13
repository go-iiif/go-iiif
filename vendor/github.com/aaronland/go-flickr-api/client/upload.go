package client

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"path/filepath"
)

// Most of the code in this file has been copypasted with minor
// updates from https://github.com/masci/flickr/blob/v2/upload.go

// Generate a random multipart boundary string, shamelessly copypasted from the std library
func randomBoundary() (string, error) {

	var buf [30]byte

	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		return "", err
	}

	boundary := fmt.Sprintf("%x", buf[:])
	return boundary, nil
}

// Encode the file and request parameters in a multipart body.
// File contents are streamed into the request using an io.Pipe in a separated goroutine
func streamUploadBody(ctx context.Context, body *io.PipeWriter, file_name string, boundary string, fh io.Reader, args *url.Values) error {

	// multipart writer to fill the body
	defer body.Close()

	writer := multipart.NewWriter(body)
	writer.SetBoundary(boundary)

	part, err := writer.CreateFormFile("photo", filepath.Base(file_name))

	if err != nil {
		return err
	}

	// fill the photo field
	_, err = io.Copy(part, fh)

	if err != nil {
		return err
	}

	// dump other params
	for key, val := range *args {
		_ = writer.WriteField(key, val[0])
	}

	// close the form writer
	return writer.Close()
}
