package source

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/jtacoma/uritemplates"
)

type URITemplateSource struct {
	Source
	template *uritemplates.UriTemplate
	client   *http.Client
	uri      string
}

func init() {
	ctx := context.Background()
	err := RegisterSource(ctx, "rfc6570", NewURITemplateSource)

	if err != nil {
		panic(err)
	}
}

func NewURITemplateSource(ctx context.Context, uri string) (Source, error) {

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

	us := URITemplateSource{
		template: template,
		client:   client,
	}

	return &us, nil
}

func (us *URITemplateSource) String() string {
	return us.uri
}

func (us *URITemplateSource) Read(id string) ([]byte, error) {

	values := make(map[string]interface{})
	values["id"] = id

	uri, err := us.template.Expand(values)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", uri, nil)

	rsp, err := us.client.Do(req)

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

func (us *URITemplateSource) Close() error {
	return nil
}
