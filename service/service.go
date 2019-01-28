package service

// http://iiif.io/api/image/2.1/#related-services
// http://iiif.io/api/annex/services/

type Service interface {
	Context() string
	Profile() string
	Label() string
	Value() interface{}
}
