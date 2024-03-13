package source

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/aaronland/go-flickr-api/client"
	iiifcache "github.com/go-iiif/go-iiif/v6/cache"
	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
)

type FlickrSource struct {
	Source
	cache         iiifcache.Cache
	http_client   *http.Client
	flickr_client client.Client
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

	if config.Flickr.ClientURI == "" {
		return nil, fmt.Errorf("Missing config.Flickr.ClientURI property")
	}

	ctx := context.Background()

	flickr_client, err := client.NewClient(ctx, config.Flickr.ClientURI)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new Flickr client, %w", err)
	}

	http_client := &http.Client{}

	fs := FlickrSource{
		flickr_client: flickr_client,
		http_client:   http_client,
		cache:         cache,
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

	rsp, err := fs.http_client.Do(req)

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
	values.Add("method", "flickr.photos.getSizes")
	values.Add("nojsoncallback", "1")

	ctx := context.Background()

	rsp, err := fs.flickr_client.ExecuteMethod(ctx, &values)

	if err != nil {
		return "", err
	}

	defer rsp.Close()

	body, err := io.ReadAll(rsp)

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
