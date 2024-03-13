package response

import (
	"encoding/json"
	"io"
)

type Login struct {
	Status    string `json:"stat"`
	Error     *Error `json:"error,omitempty"`
	User      *User  `json:"user,omitempty"`
	PathAlias string `json:"path_alias"`
}

func UnmarshalCheckLoginJSONResponse(fh io.Reader) (*Login, error) {

	var l *Login

	dec := json.NewDecoder(fh)
	err := dec.Decode(&l)

	if err != nil {
		return nil, err
	}

	return l, nil
}
