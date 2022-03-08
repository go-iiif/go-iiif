package info

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/pretty"
)

// MarshalJSON encodes and formats 'i' as JSON.
func MarshalJSON(i *Info) ([]byte, error) {

	body, err := json.Marshal(i)

	if err != nil {
		return nil, fmt.Errorf("Failed to marshal info, %w", err)
	}

	body = pretty.Pretty(body)
	return body, nil
}
