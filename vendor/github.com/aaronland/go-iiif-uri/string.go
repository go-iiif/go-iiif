package uri

type StringURI struct {
	URI
	raw string
}

func NewStringURI(raw string) (URI, error) {

	u := StringURI{
		raw: raw,
	}

	return &u, nil
}

func (u *StringURI) URL() string {
	return u.raw
}

func (u *StringURI) String() string {
	return u.raw
}
