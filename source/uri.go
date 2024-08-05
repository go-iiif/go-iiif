package source

// URI as in "URI Template" – this is a badly named package

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
	"github.com/jtacoma/uritemplates"
)

type URISource struct {
	Source
	template *uritemplates.UriTemplate
	client   *http.Client
}

func init() {
	ctx := context.Background()
	err := RegisterSource(ctx, "rfc6570", NewURISourceFromURI)

	if err != nil {
		panic(err)
	}
}

func NewURISourceURIFromConfig(cfg *iiifconfig.Config) (string, error) {

	uri := cfg.Images.Source.URI

	if uri == "" {
		uri = fmt.Sprintf("rfc6570://?template=%s", cfg.Images.Source.Path)
	}

	return uri, nil
}

func NewURISource(cfg *iiifconfig.Config) (Source, error) {

	uri, err := NewURISourceURIFromConfig(cfg)

	if err != nil {
		return nil, err
	}

	return NewURISourceFromURI(uri)
}

func NewURISourceFromURI(uri string) (Source, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	q := u.Query()
	t := q.Get("template")

	client := &http.Client{}
	template, err := uritemplates.Parse(t)

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
