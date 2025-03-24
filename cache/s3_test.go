package cache

import (
	"testing"

	_ "github.com/aaronland/gocloud-blob/s3"	
	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
)

func TestS3Cache(t *testing.T) {

	t.Skip()
	
	cfg := iiifconfig.CacheConfig{
		URI: "s3blob://config/config.json?region=us-west-2&credentials=env:",
	}

	_, err := NewS3Cache(cfg)

	if err != nil {
		t.Fatalf("Failed to create new S3 cache, %v", err)
	}
}
