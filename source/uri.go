package source

// URI as in "URI Template" – this is a badly named package

import (
	"fmt"
	"io"
	"net/http"

	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
	"github.com/jtacoma/uritemplates"
)

type URISource struct {
	Source
	template *uritemplates.UriTemplate
	client   *http.Client
}

func NewURISourceURIFromConfig(cfg *iiifconfig.Config) (string, error) {

	uri := cfg.Images.Source.URI

	if uri == "" {
		uri = fmt.Sprintf("file://%s", cfg.Images.Source.Path)
	}

	return uri, nil
}

func NewURISource(cfg *iiifconfig.Config) (*URISource, error) {

	uri, err := NewURISourceURIFromConfig(cfg)

	if err != nil {
		return nil, err
	}

	return NewURISourceFromURI(uri)
}

func NewURISourceFromURI(uri string) (*URISource, error) {

	client := &http.Client{}
	template, err := uritemplates.Parse(uri)

	if err != nil {
		return nil, err
	}

	us := URISource{
		template: template,
		client:   client,
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

	req, err := http.NewRequest("GET", uri, nil)

	// t1 := time.Now()
	rsp, err := us.client.Do(req)

	// t2 := time.Since(t1)
	// log.Println("time to fetch", uri, t2)

	if err != nil {
		return nil, err
	}

	defer rsp.Body.Close()
	body, err := io.ReadAll(rsp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}
