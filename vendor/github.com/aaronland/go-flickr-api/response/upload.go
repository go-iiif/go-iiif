package response

import (
	"encoding/xml"
	"io"
)

type Upload struct {
	XMLName xml.Name `xml:"rsp" json:"-"`
	// The Flickr API response status.
	Status string       `xml:"stat,attr" json:"stat"`
	Error  *Error       `xml:"err,omitempty" json:"error,omitempty"`
	Photo  *UploadPhoto `xml:"photoid" json:"photo,omitempty"`
}

type UploadPhoto struct {
	Id             int64  `xml:",chardata" json:"id"`
	Secret         string `xml:"secret,attr,omitempty" json:"secret"`
	OriginalSecret string `xml:"originalsecret,attr,omitempty" json:"originalsecret"`
}

type UploadTicket struct {
	XMLName xml.Name `xml:"rsp"`
	// The Flickr API response status.
	Status   string `xml:"stat,attr"`
	Error    *Error `xml:"err,omitempty"`
	TicketId string `xml:"ticketid"`
}

func UnmarshalUploadResponse(fh io.Reader) (*Upload, error) {

	body, err := io.ReadAll(fh)

	if err != nil {
		return nil, err
	}

	var up *Upload

	err = xml.Unmarshal([]byte(body), &up)

	if err != nil {
		return nil, err
	}

	return up, nil
}

func UnmarshalUploadTicketResponse(fh io.Reader) (*UploadTicket, error) {

	body, err := io.ReadAll(fh)

	if err != nil {
		return nil, err
	}

	var up *UploadTicket

	err = xml.Unmarshal(body, &up)

	if err != nil {
		return nil, err
	}

	return up, nil
}
