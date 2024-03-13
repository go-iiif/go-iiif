package response

import (
	"encoding/json"
	"io"
)

// CheckTicket is a struct that maps to the Flickr API flickr.photos.upload.checkTickets
// method response.
type CheckTicket struct {
	// An Uploader instance that maps to the "uploader" element in the API response.
	Uploader *Uploader `json:"uploader"`
}

// Uploader is a struct that maps to the "uploader" element in a Flickr API
// flickr.photos.upload.checkTickets method response.
type Uploader struct {
	// An array of UploaderTicket instances that map to the "ticket" array in the API response.
	Tickets []*UploaderTicket `json:"ticket"`
}

// UploaderTicket is a struct that maps to the "ticket" array in a Flickr API
// flickr.photos.upload.checkTickets method response.
type UploaderTicket struct {
	// A Flickr API upload ticket ID.
	TicketId string `json:"id"`
	// A numeric flag (1 or 0) indicating whether an upload ticket has been completed.
	Complete int `json:"complete"`
	// The Flickr photo ID for a successful upload. Note that this is encoded (by the Flickr API) a string.
	PhotoId string `json:"photoid"`
	// The creation time (Unix timestamp) for a successful upload. Note that this is encoded (by the Flickr API) a string.
	Imported string `json:"imported"`
}

// Unmarshal Flickr API flickr.photos.upload method response, for asynchronous uploads, in to a UploadTicket instance.
func UmarshalUploadTicketResponse(fh io.Reader) (*UploadTicket, error) {

	var ut *UploadTicket

	dec := json.NewDecoder(fh)
	err := dec.Decode(&ut)

	if err != nil {
		return nil, err
	}

	return ut, nil
}

// Unmarshal Flickr API flickr.photos.upload.checkTickets method response into a CheckTicket instance.
func UnmarshalCheckTicketResponse(fh io.Reader) (*CheckTicket, error) {

	var ct *CheckTicket

	dec := json.NewDecoder(fh)
	err := dec.Decode(&ct)

	if err != nil {
		return nil, err
	}

	return ct, nil
}
