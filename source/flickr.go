package source

import (
	"encoding/json"
	"errors"
	iiifcache "github.com/thisisaaronland/go-iiif/cache"
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
	cache     iiifcache.Cache
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

	cache_config := iiifconfig.CacheConfig{
		TTL:   3600,
		Limit: 1,
	}

	cache, err := iiifcache.NewMemoryCache(cache_config)

	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	apikey := config.Flickr.ApiKey
	apisecret := config.Flickr.ApiSecret

	fs := FlickrSource{
		apikey:    apikey,
		apisecret: apisecret,
		client:    client,
		cache:     cache,
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

	cached, err := fs.cache.Get(id)

	if err == nil {
		return string(cached), nil
	}

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
	log.Println(req.URL.RawQuery)

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

	go func() {
		fs.cache.Set(id, []byte(source))
	}()

	return source, nil
}
