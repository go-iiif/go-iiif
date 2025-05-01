package source

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/aaronland/go-flickr-api/client"
	iiifcache "github.com/go-iiif/go-iiif/v8/cache"
)

type FlickrSource struct {
	Source
	cache         iiifcache.Cache
	http_client   *http.Client
	flickr_client client.Client
	uri           string
	safe_uri      string
}

func init() {
	ctx := context.Background()
	err := RegisterSource(ctx, "flickr", NewFlickrSource)
	if err != nil {
		panic(err)
	}
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

func NewFlickrSource(ctx context.Context, uri string) (Source, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	client_uri := q.Get("client-uri")
	flickr_client, err := client.NewClient(ctx, client_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new Flickr client, %w", err)
	}

	cache_uri := "memory://?ttl=3600&limit=1"
	cache, err := iiifcache.NewCache(ctx, cache_uri)

	if err != nil {
		return nil, err
	}

	//

	client_u, err := url.Parse(client_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse Flickr client URI, %w", err)
	}

	client_q := url.Values{}
	client_q.Set("consumer_key", "{KEY}")
	client_q.Set("consumer_secret", "{SECRET}")
	client_q.Set("oauth_token", "{TOKEN}")
	client_q.Set("oauth_token_secret", "{SECRET}")

	client_u.RawQuery = client_q.Encode()

	safe_q := url.Values{}
	safe_q.Set("client-uri", client_u.String())

	safe_u, _ := url.Parse(uri)
	safe_u.RawQuery = safe_q.Encode()

	safe_uri := safe_u.String()

	//

	http_client := &http.Client{}

	fs := FlickrSource{
		flickr_client: flickr_client,
		http_client:   http_client,
		cache:         cache,
		uri:           uri,
		safe_uri:      safe_uri,
	}

	return &fs, nil
}

func (fs *FlickrSource) String() string {
	return fs.safe_uri
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

func (fs *FlickrSource) Close() error {
	return nil
}
