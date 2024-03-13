package auth

import (
	"math/rand"
	"time"
)

// Common interface for authentication request tokens
type RequestToken interface {
	// A temporary authentication request token. This is used to create an authorization token request.
	Token() string
	// A temporary authentication request secret. This is used to exchange an authorization token for an access token.
	Secret() string
}

// Common interface for authentication authorization tokens
type AuthorizationToken interface {
	// A temporary authentication authorization token
	Token() string
	// A temporary authentication authorization verification string. This is used to exchange an authorization	token for an access token.
	Verifier() string
}

// Common interface for authentication access tokens
type AccessToken interface {
	// A permanent access token associatin an application to a user account.
	Token() string
	// A permanent secret key for an access token.
	Secret() string
}

// Generate a random string of 8 chars, needed for OAuth1 signatures.
func GenerateNonce() string {

	rand.Seed(time.Now().UTC().UnixNano())

	// For convenience, use a set of chars we don't need to url-escape
	var letters = []rune("123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ")

	b := make([]rune, 8)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
