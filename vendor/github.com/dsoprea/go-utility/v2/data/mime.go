package ridata

import (
	"io"
	"os"
	"strings"

	"net/http"

	"github.com/dsoprea/go-logging"
)

const (
	// MimetypeLeadBytesCount is the number of bytes to use for detection.
	MimetypeLeadBytesCount = 512
)

// GetMimetypeFromContent uses net/http to map from magic-bytes to mime-type.
func GetMimetypeFromContent(r io.Reader, fileSize int64) (mimetype string, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	// TODO(dustin): !! Add test.

	leadCount := int64(MimetypeLeadBytesCount)
	if fileSize > 0 && fileSize < leadCount {
		leadCount = fileSize
	}

	buffer := make([]byte, leadCount)

	n, err := io.ReadFull(r, buffer)
	if err != nil {
		// We can return EOF if a) we weren't given a filesize and the file did
		// not haveat least as many bytes as we check by default, or b) the file-
		// size is actually (0).
		if err == io.EOF {
			return "", err
		}

		log.Panic(err)
	}

	buffer = buffer[:n]

	// Always returns a valid mime-type.
	contentType := http.DetectContentType(buffer)

	contentType = strings.TrimRight(contentType, ";")

	return contentType, nil
}

// DetectMimetype is a wrapper for GetMimetypeFromContent which returns the
// mime-type for the given `File`. An empty-string is returned if it is a zero-
// length file.
func DetectMimetype(f *os.File) (mimetype string, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	originalOffsetRaw, err := f.Seek(0, os.SEEK_CUR)
	log.PanicIf(err)

	fi, err := f.Stat()
	log.PanicIf(err)

	fileSize := fi.Size()

	if fileSize == 0 {
		return "", nil
	}

	mimetype, err = GetMimetypeFromContent(f, fileSize)
	log.PanicIf(err)

	_, err = f.Seek(originalOffsetRaw, os.SEEK_SET)
	log.PanicIf(err)

	return mimetype, nil
}
