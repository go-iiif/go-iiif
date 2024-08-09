package tools

import (
	"testing"
)

func TestDefaultURIFunc(t *testing.T) {

	tests := []string{
		"idsecret:///IMG_9998.jpg?id=9998&secret=abc&secret_o=def&format=jpg&label=x",
	}

	fn := DefaultURIFunc()

	for _, u := range tests {

		_, err := fn(u)

		if err != nil {
			t.Fatalf("URI func failed for '%s', %v", u, err)
		}
	}
}
