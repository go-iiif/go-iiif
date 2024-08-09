package autometa

import (
	"fmt"
	"io"

	"github.com/mandykoh/prism/meta"
	"github.com/mandykoh/prism/meta/jpegmeta"
	"github.com/mandykoh/prism/meta/pngmeta"
	"github.com/mandykoh/prism/meta/webpmeta"
)

// Load loads the metadata for an image stream, which may be one of the
// supported image formats.
//
// Only as much of the stream is consumed as necessary to extract the metadata;
// the returned stream contains a buffered copy of the consumed data such that
// reading from it will produce the same results as fully reading the input
// stream. This provides a convenient way to load the full image after loading
// the metadata.
//
// An error is returned if basic metadata could not be extracted. The returned
// stream still provides the full image data.
func Load(r io.Reader) (md *meta.Data, imgStream io.Reader, err error) {

	loaders := []func(r io.Reader) (*meta.Data, io.Reader, error){
		pngmeta.Load,
		jpegmeta.Load,
		webpmeta.Load,
	}

	inputStream := r

	for _, loader := range loaders {
		md, nextStream, err := loader(inputStream)
		if err == nil {
			return md, nextStream, nil
		}

		inputStream = nextStream
	}

	return nil, inputStream, fmt.Errorf("unrecognised image format")
}
