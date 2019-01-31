package uri

type IIIFURI struct {
	URI
	raw string
}

func NewIIIFURI(raw string) (URI, error) {

	u := IIIFURI{
		raw: raw,
	}

	return &u, nil
}

func (u *IIIFURI) URL() string {
	return u.raw
}

func (u *IIIFURI) String() string {
	return u.raw
}
