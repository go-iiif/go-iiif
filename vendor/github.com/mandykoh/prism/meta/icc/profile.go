package icc

type Profile struct {
	Header   Header
	TagTable TagTable
}

func (p *Profile) Description() (string, error) {
	return p.TagTable.getProfileDescription()
}

func newProfile() *Profile {
	return &Profile{
		TagTable: emptyTagTable(),
	}
}
