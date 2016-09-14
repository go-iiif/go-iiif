package source

import (
	"github.com/jtacoma/uritemplates"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	"io/ioutil"
	_ "log"
	"net/http"
	_ "time"
)

type URISource struct {
	Source
	template *uritemplates.UriTemplate
	client   *http.Client
}

func NewURISource(config *iiifconfig.Config) (*URISource, error) {

	cfg := config.Images

	client := &http.Client{}
	template, err := uritemplates.Parse(cfg.Source.Path)

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
	body, err := ioutil.ReadAll(rsp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}
