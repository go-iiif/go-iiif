package http

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/whosonfirst/go-sanitize"
	gohttp "net/http"
	gourl "net/url"
	"strings"
)

type IIIFParameters struct {
	Identifier string
	Region     string
	Size       string
	Rotation   string
	Quality    string
	Format     string
}

type IIIFQueryParser struct {
	opts *sanitize.Options
	vars map[string]string
}

func NewIIIFQueryParser(r *gohttp.Request) (*IIIFQueryParser, error) {

	opts := sanitize.DefaultOptions()
	vars := mux.Vars(r)

	p := IIIFQueryParser{
		opts: opts,
		vars: vars,
	}

	return &p, nil
}

func (p *IIIFQueryParser) GetIIIFParameter(key string) (string, error) {

	var err error

	value := p.vars[key]

	value, err = sanitize.SanitizeString(value, p.opts)

	if err != nil {
		return "", err
	}

	value, err = gourl.QueryUnescape(value)

	if err != nil {
		return "", err
	}

	// This should be already be stripped out by the time we get here but just
	// in case... (20160926/thisisaaronland)

	if strings.Contains(value, "../") {
		msg := fmt.Sprintf("Invalid key %s", key)
		err := errors.New(msg)
		return "", err
	}

	return value, nil
}

func (p *IIIFQueryParser) GetIIIFParameters() (*IIIFParameters, error) {

	id, err := p.GetIIIFParameter("identifier")

	if err != nil {
		return nil, err
	}

	region, err := p.GetIIIFParameter("region")

	if err != nil {
		return nil, err
	}

	size, err := p.GetIIIFParameter("size")

	if err != nil {
		return nil, err
	}

	rotation, err := p.GetIIIFParameter("rotation")

	if err != nil {
		return nil, err
	}

	quality, err := p.GetIIIFParameter("quality")

	if err != nil {
		return nil, err
	}

	format, err := p.GetIIIFParameter("format")

	if err != nil {
		return nil, err
	}

	params := IIIFParameters{
		Identifier: id,
		Region:     region,
		Size:       size,
		Rotation:   rotation,
		Quality:    quality,
		Format:     format,
	}

	return &params, nil
}
