package cache

import (
	"context"
	"testing"

	iiifconfig "github.com/go-iiif/go-iiif/v8/config"
	iiifdefaults "github.com/go-iiif/go-iiif/v8/defaults"
)

func TestNewImagesCacheFromConfig(t *testing.T) {

	ctx := context.Background()

	r, err := iiifdefaults.FS.Open("config.json")

	if err != nil {
		t.Fatalf("Failed to load config from defaults, %v", err)
	}

	defer r.Close()

	cfg, err := iiifconfig.NewConfigFromReader(r)

	if err != nil {
		t.Fatalf("Failed to derive config from defaults, %v", err)
	}

	_, err = NewCache(ctx, cfg.Images.Cache.URI)

	if err != nil {
		t.Fatalf("Failed to derive images cache from defaults, %v", err)
	}
}

func TestNewDerivativesCacheFromConfig(t *testing.T) {

	ctx := context.Background()

	r, err := iiifdefaults.FS.Open("config.json")

	if err != nil {
		t.Fatalf("Failed to load config from defaults, %v", err)
	}

	defer r.Close()

	cfg, err := iiifconfig.NewConfigFromReader(r)

	if err != nil {
		t.Fatalf("Failed to derive config fromd defaults, %v", err)
	}

	_, err = NewCache(ctx, cfg.Derivatives.Cache.URI)

	if err != nil {
		t.Fatalf("Failed to derive derivatives cache from defaults, %v", err)
	}
}
