package source

import (
	"context"
	"net/url"
	"testing"
)

func TestNewFlickrSourceFromURI(t *testing.T) {

	ctx := context.Background()

	client_uri := "oauth1://?consumer_key={KEY}&consumer_secret={SECRET}&oauth_token={TOKEN}&oauth_token_secret={SECRET}"

	source_q := url.Values{}
	source_q.Set("client-uri", client_uri)

	source_u := url.URL{}
	source_u.Scheme = "flickr"
	source_u.RawQuery = source_q.Encode()
	
	uri := source_u.String()

	_, err := NewFlickrSourceFromURI(uri)

	if err != nil {
		t.Fatalf("Failed to create Flickr source '%s', %v", uri, err)
	}

	_, err = NewSource(ctx, uri)

	if err != nil {
		t.Fatalf("Failed to create Flickr source calling NewSource, %v", err)
	}

	
}
