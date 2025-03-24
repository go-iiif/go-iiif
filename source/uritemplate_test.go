package source

import (
	"context"
	"testing"
)

func TestNewURISourceFromURI(t *testing.T) {

	ctx := context.Background()

	uri := "rfc6570://?template=http://127.0.0.1/{name}"

	_, err := NewSource(ctx, uri)

	if err != nil {
		t.Fatalf("Failed to create URI template source calling NewSource, %v", err)
	}

}
