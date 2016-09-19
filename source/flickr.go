package source

import (
	"github.com/jtacoma/uritemplates"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	"io/ioutil"
	"net/http"
)

type FlickrSource struct {
	Source
	apikey    string
	apisecret string
	client    *http.Client
}

func NewFlickrSource(config *iiifconfig.Config) (*FlickrSource, error) {

	client := &http.Client{}

	apikey := "fixme"
	apisecret := "fixme"

	fs := FlickrSource{
		apikey:    apikey,
		apisecret: apisecret,
		client:    client,
	}

	return &fs, nil
}

func (fs *FlickrSource) Read(id string) ([]byte, error) {

	return nil, errors.New("Please write me")

	/*

	   call https://www.flickr.com/services/api/flickr.photos.getSizes.html

	   check for @label=Original
	   check for @label=Large
	   grab @source

	   cache results ?

	*/

	url := "http://farm2.staticflickr.com/1103/567229075_6dc09dc6da_o.jpg"

	req, err := http.NewRequest("GET", url, nil)

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
