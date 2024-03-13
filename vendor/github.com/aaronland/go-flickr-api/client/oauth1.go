package client

import (
	"context"
	"fmt"
	"github.com/aaronland/go-flickr-api/auth"
	"github.com/whosonfirst/go-ioutil"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// The default Flickr endpoint for OAuth1 authorization requests.
const OAUTH1_AUTHORIZE_ENDPOINT string = "https://www.flickr.com/services/oauth/authorize"

// The default Flickr endpoint for OAuth1 request token requests.
const OAUTH1_REQUEST_TOKEN_ENDPOINT string = "https://www.flickr.com/services/oauth/request_token"

// The default Flickr endpoint for OAuth1 access token requests.
const OAUTH1_ACCESS_TOKEN_ENDPOINT string = "https://www.flickr.com/services/oauth/access_token"

func init() {

	ctx := context.Background()
	err := RegisterClient(ctx, "oauth1", NewOAuth1Client)

	if err != nil {
		panic(err)
	}
}

// OAuth1Client implements the Client interface for invoking the Flickr API using the OAuth1 authentication
// and authorization mechanism.
type OAuth1Client struct {
	http_client        *http.Client
	api_endpoint       *url.URL
	consumer_key       string
	consumer_secret    string
	oauth_token        string
	oauth_token_secret string
}

// Create a new OAuth1Client instance conforming to the Client interface. OAuth1Client instances are
// create by passing in a context.Context instance and a URI string in the form of:
// oauth1://?consumer_key={KEY}&consumer_secret={SECRET} or:
// oauth1://?consumer_key={KEY}&consumer_secret={SECRET}&oauth1_token={TOKEN}&oauth1_token_secret={SECRET}
func NewOAuth1Client(ctx context.Context, uri string) (Client, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	key := q.Get("consumer_key")
	secret := q.Get("consumer_secret")

	if key == "" {
		return nil, fmt.Errorf("Missing ?consumer_key parameter")
	}

	if secret == "" {
		return nil, fmt.Errorf("Missing ?consumer_secret parameter")
	}

	http_client := &http.Client{}

	cl := &OAuth1Client{
		http_client:     http_client,
		consumer_key:    key,
		consumer_secret: secret,
	}

	oauth_token := q.Get("oauth_token")
	oauth_token_secret := q.Get("oauth_token_secret")

	if oauth_token != "" {
		cl.oauth_token = oauth_token
	}

	if oauth_token_secret != "" {
		cl.oauth_token_secret = oauth_token_secret
	}

	return cl, nil
}

// Return a new Client instance that uses the credentials included in the auth.AccessToken instance.
func (cl *OAuth1Client) WithAccessToken(ctx context.Context, access_token auth.AccessToken) (Client, error) {

	new_cl := &OAuth1Client{
		http_client:        cl.http_client,
		consumer_key:       cl.consumer_key,
		consumer_secret:    cl.consumer_secret,
		oauth_token:        access_token.Token(),
		oauth_token_secret: access_token.Secret(),
	}

	return new_cl, nil
}

// Call the Flickr API and create a new request token as part of the token authorization flow.
func (cl *OAuth1Client) GetRequestToken(ctx context.Context, cb_url string) (auth.RequestToken, error) {

	endpoint, err := url.Parse(OAUTH1_REQUEST_TOKEN_ENDPOINT)

	if err != nil {
		return nil, err
	}

	http_method := "GET"

	args := &url.Values{}
	args.Set("oauth_callback", cb_url)

	args, err = cl.signArgs(http_method, endpoint, args, "")

	if err != nil {
		return nil, err
	}

	endpoint.RawQuery = args.Encode()

	req, err := http.NewRequest(http_method, endpoint.String(), nil)

	if err != nil {
		return nil, err
	}

	fh, err := cl.call(ctx, req)

	if err != nil {
		return nil, err
	}

	defer fh.Close()

	rsp_body, err := io.ReadAll(fh)

	if err != nil {
		return nil, err
	}

	return auth.UnmarshalOAuth1RequestToken(string(rsp_body))
}

// Generate the URL using a request token and permissions string used to redirect a user to in order to authorize a token request.
func (cl *OAuth1Client) GetAuthorizationURL(ctx context.Context, req auth.RequestToken, perms string) (string, error) {

	q := url.Values{}
	q.Set("oauth_token", req.Token())

	if perms != "" {
		q.Set("perms", perms)
	}

	endpoint, err := url.Parse(OAUTH1_AUTHORIZE_ENDPOINT)

	if err != nil {
		return "", err
	}

	endpoint.RawQuery = q.Encode()

	return endpoint.String(), nil
}

// Call the Flickr API to exchange a request and authorization token for a permanent access token.
func (cl *OAuth1Client) GetAccessToken(ctx context.Context, req_token auth.RequestToken, auth_token auth.AuthorizationToken) (auth.AccessToken, error) {

	endpoint, err := url.Parse(OAUTH1_ACCESS_TOKEN_ENDPOINT)

	if err != nil {
		return nil, err
	}

	http_method := "GET"

	args := &url.Values{}

	// See what's going on here? The token is coming from the authentication
	// response but the secret is coming from the request response. It took
	// me a long time to figure that out... (20210331/thisisaaronland)

	args.Set("oauth_token", auth_token.Token())
	args.Set("oauth_verifier", auth_token.Verifier())

	args, err = cl.signArgs(http_method, endpoint, args, req_token.Secret())

	if err != nil {
		return nil, err
	}

	endpoint.RawQuery = args.Encode()

	req, err := http.NewRequest(http_method, endpoint.String(), nil)

	if err != nil {
		return nil, err
	}

	fh, err := cl.call(ctx, req)

	if err != nil {
		return nil, err
	}

	defer fh.Close()

	rsp_body, err := io.ReadAll(fh)

	if err != nil {
		return nil, err
	}

	return auth.UnmarshalOAuth1AccessToken(string(rsp_body))
}

// Execute a Flickr API method. If not "format" parameter in include in the url.Values instance passed to the method API responses will be returned as JSON (by automatically assign the 'nojsoncallback=1' and 'format=json' parameters).
func (cl *OAuth1Client) ExecuteMethod(ctx context.Context, args *url.Values) (io.ReadSeekCloser, error) {

	endpoint, err := url.Parse(API_ENDPOINT)

	if err != nil {
		return nil, err
	}

	http_method := "GET"

	if args.Get("format") == "" {
		args.Set("nojsoncallback", "1")
		args.Set("format", "json")
	}

	if cl.oauth_token != "" {
		args.Set("oauth_token", cl.oauth_token)
	}

	args, err = cl.signArgs(http_method, endpoint, args, cl.oauth_token_secret)

	if err != nil {
		return nil, err
	}

	endpoint.RawQuery = args.Encode()

	req, err := http.NewRequest(http_method, endpoint.String(), nil)

	if err != nil {
		return nil, err
	}

	return cl.call(ctx, req)
}

// Upload an image using the Flickr API.
func (cl *OAuth1Client) Upload(ctx context.Context, fh io.Reader, args *url.Values) (io.ReadSeekCloser, error) {

	endpoint, err := url.Parse(UPLOAD_ENDPOINT)

	if err != nil {
		return nil, err
	}

	return cl.upload(ctx, endpoint, fh, args)
}

// Replace an image using the Flickr API.
func (cl *OAuth1Client) Replace(ctx context.Context, fh io.Reader, args *url.Values) (io.ReadSeekCloser, error) {

	endpoint, err := url.Parse(REPLACE_ENDPOINT)

	if err != nil {
		return nil, err
	}

	return cl.upload(ctx, endpoint, fh, args)
}

func (cl *OAuth1Client) upload(ctx context.Context, endpoint *url.URL, fh io.Reader, args *url.Values) (io.ReadSeekCloser, error) {

	http_method := "POST"

	args.Set("oauth_token", cl.oauth_token)

	args, err := cl.signArgs(http_method, endpoint, args, cl.oauth_token_secret)

	if err != nil {
		return nil, err
	}

	fname := "upload"
	boundary, err := randomBoundary()

	if err != nil {
		return nil, err
	}

	r, w := io.Pipe()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {

		err := streamUploadBody(ctx, w, fname, boundary, fh, args)

		if err != nil {
			log.Printf("Failed to stream upload body for '%s', %v", fname, err)
			cancel()
		}
	}()

	req, err := http.NewRequest(http_method, endpoint.String(), r)

	if err != nil {
		return nil, err
	}

	req.Header.Set("content-type", "multipart/form-data; boundary="+boundary)
	req.ContentLength = -1 // unknown

	// This response is formatted in the REST API response style.
	// https://www.flickr.com/services/api/response.rest.html

	return cl.call(ctx, req)
}

func (cl *OAuth1Client) call(ctx context.Context, req *http.Request) (io.ReadSeekCloser, error) {

	req = req.WithContext(ctx)

	rsp, err := cl.http_client.Do(req)

	if err != nil {
		return nil, err
	}

	if rsp.StatusCode != http.StatusOK {
		rsp.Body.Close()
		return nil, fmt.Errorf("API call failed with status '%s'", rsp.Status)
	}

	return ioutil.NewReadSeekCloser(rsp.Body)
}

func (cl *OAuth1Client) signArgs(http_method string, endpoint *url.URL, args *url.Values, secret string) (*url.Values, error) {

	now := time.Now()
	ts := now.Unix()

	str_ts := strconv.FormatInt(ts, 10)

	nonce := auth.GenerateNonce()

	args.Set("oauth_version", "1.0")
	args.Set("oauth_signature_method", "HMAC-SHA1")

	args.Set("oauth_nonce", nonce)
	args.Set("oauth_timestamp", str_ts)
	args.Set("oauth_consumer_key", cl.consumer_key)

	sig := cl.getSignature(http_method, endpoint, args, secret)
	args.Set("oauth_signature", sig)

	return args, nil
}

func (cl *OAuth1Client) getSignature(http_method string, endpoint *url.URL, args *url.Values, token_secret string) string {

	key := fmt.Sprintf("%s&%s", url.QueryEscape(cl.consumer_secret), url.QueryEscape(token_secret))
	base_string := auth.GenerateOAuth1SigningBaseString(http_method, endpoint, args)

	return auth.GenerateOAuth1Signature(key, base_string)
}
