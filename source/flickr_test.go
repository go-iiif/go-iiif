package source

import (
	"context"
	"fmt"
	"net/url"
	"testing"
)

func TestNewFlickrSourceFromURI(t *testing.T) {

	ctx := context.Background()

	client_uri := "oauth1://?consumer_key=abcdef&consumer_secret=124567&oauth_token=omgwtf&oauth_token_secret=bbq"

	source_q := url.Values{}
	source_q.Set("client-uri", client_uri)

	source_u := url.URL{}
	source_u.Scheme = "flickr"
	source_u.RawQuery = source_q.Encode()

	uri := source_u.String()

	s, err := NewSource(ctx, uri)

	if err != nil {
		t.Fatalf("Failed to create Flickr source calling NewSource, %v", err)
	}

	// Because Source interface doesn't have a String() method yet
	// This is scheduled for v7
	str_u := fmt.Sprintf("%s", s)

	// I don't understand why net/url.URL.String() strips the "//" from unknown
	// schemes. Anyway, the point is to strip OAuth secrets from the Flickr client
	// URI.
	expected_str := "flickr:?client-uri=oauth1%3A%3Fconsumer_key%3D%257BKEY%257D%26consumer_secret%3D%257BSECRET%257D%26oauth_token%3D%257BTOKEN%257D%26oauth_token_secret%3D%257BSECRET%257D"

	if str_u != expected_str {
		t.Fatalf("String value for source is not expected safe URI")
	}

}
