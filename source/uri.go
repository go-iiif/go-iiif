package source

import (
	"github.com/jtacoma/uritemplates"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
)

type URISource struct {
	Source
	template *uritemplates.UriTemplate
}

func NewURISource(config *iiifconfig.Config) (*URISource, error) {

	cfg := config.Images

	template, err := uritemplates.Parse(cfg.Source.Path)

	if err != nil {
		return nil, err
	}

	us := URISource{
		template: template,
	}

	return &us, nil
}

func (us *URISource) Read(id string) ([]byte, error) {

	values := make(map[string]interface{})
	values["id"] = id

	uri, err := us.template.Expand(values)

	if err != nil {
		return nil, err
	}

	// FETCH URI HERE...

	return []byte(uri), nil
}
