package source

import (
	"encoding/json"
	"errors"
	iiifconfig "github.com/thisisaaronland/go-iiif/config"
	"io/ioutil"
	"log"
	"net/http"
)

type FlickrSource struct {
	Source
	apikey    string
	apisecret string
	client    *http.Client
}

type PhotoRsp struct {
	Sizes PhotoSizes `json:"sizes"`
}

type PhotoSizes struct {
	CanBlog     int         `json:"canblog"`
	CanPrint    int         `json:"canprint"`
	CanDownload int         `json:"candownload"`
	Size        []PhotoSize `json:"size"`
}

type PhotoSize struct {
	Label string `json:"label"`
	// it turns out these get returned as both strings and ints and
	// that makes Go sad but we don't really care either way so...
	// (20160920/thisisaaronland)
	// Width  int    `json:"width"`
	// Height int    `json:"height"`
	Source string `json:"source"`
	Url    string `json:"url"`
	Media  string `json:"media"`
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

	source, err := fs.GetSource(id)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", source, nil)

	if err != nil {
		return nil, err
	}

	rsp, err := fs.client.Do(req)

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

func (fs *FlickrSource) GetSource(id string) (string, error) {

	// please cache me (20160920/thisisaaronland)

	url := "https://api.flickr.com/services/rest/"

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return "", err
	}

	values := req.URL.Query()
	values.Add("photo_id", id)
	values.Add("api_key", fs.apikey)
	values.Add("format", "json")
	values.Add("method", "flickr.photos.getSizes")
	values.Add("nojsoncallback", "1")

	req.URL.RawQuery = values.Encode()

	rsp, err := fs.client.Do(req)

	if err != nil {
		return "", err
	}

	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)

	if err != nil {
		return "", err
	}

	// log.Println(string(body))

	var data PhotoRsp
	err = json.Unmarshal(body, &data)

	if err != nil {
		return "", err
	}

	by_label := make(map[string]PhotoSize)

	for _, sz := range data.Sizes.Size {
		by_label[sz.Label] = sz
	}

	possible := []string{"Original", "Large 2048", "Large 1600", "Large"}

	var source string

	for _, k := range possible {

		sz, ok := by_label[k]

		if ok {
			source = sz.Source
			break
		}
	}

	if source == "" {
		return source, errors.New("Unable to determine photo source")
	}

	return source, nil
}
