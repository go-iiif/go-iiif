// package response provides methods for unmarshaling Flickr API responses in to
// structs and instances corresponding to specific API methods. Not all Flickr API
// responses are supported yet and may not be until time and circumstance warrant
// their inclusion.
package response

import (
	_ "encoding/json"
	"encoding/xml"
	"fmt"
	"io"
)

// Error is a struct containing information about a failed API request.
type Error struct {
	// The numeric code for the error.
	Code int `xml:"code,attr" json:"code"`
	// The message associated with the error.
	Message string `xml:"msg,attr" json:"msg"`
}

// Return a Flickr API error as a string containing both the error code and message.
func (e *Error) Error() string {
	return fmt.Sprintf("%d %s", e.Code, e.Message)
}

// Response is a struct containing only minimal information about an API request, notably it's Status and optionally an Error associated with the request.
type Response struct {
	XMLName xml.Name `xml:"rsp" json:"-"`
	// A string label indicating whether or not an API request succeeded or failed, "ok" and "error" respectively.
	Status string `xml:"stat,attr" json:"stat"`
	// An optional Error instance containing specific details about a failed API request.
	Error *Error `xml:"err,omitempty" json:"err,omitempty"`
}

// Unmarshal a Flickr API response in to a Response struct.
func UnmarshalResponse(fh io.Reader) (*Response, error) {

	body, err := io.ReadAll(fh)

	if err != nil {
		return nil, err
	}

	var rsp *Response

	err = xml.Unmarshal([]byte(body), &rsp)

	if err != nil {
		return nil, err
	}

	return rsp, nil
}
