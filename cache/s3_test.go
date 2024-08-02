package cache

import (
	"testing"

	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
)

func TestS3Cache(t testing.T) {

	cfg := &iiifconfig.CacheConfig{
		Path:        "config.json",
		Prefix:      "",
		Region:      "us-east-1",
		Credentials: "session",
	}

	_, err := NewS3Cache(cfg)

	if err != nil {
		t.Fatalf("Failed to create new S3 cache, %v", err)
	}
}
