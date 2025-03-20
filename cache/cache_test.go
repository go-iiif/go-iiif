package cache

import (
	"testing"

	iiifconfig "github.com/go-iiif/go-iiif/v6/config"
	iiifdefaults "github.com/go-iiif/go-iiif/v6/defaults"
	_ "gocloud.dev/blob/memblob"
)

func TestNewImagesCacheFromConfig(t *testing.T) {

	r, err := iiifdefaults.FS.Open("config.json")

	if err != nil {
		t.Fatalf("Failed to load config from defaults, %v", err)
	}

	defer r.Close()

	cfg, err := iiifconfig.NewConfigFromReader(r)

	if err != nil {
		t.Fatalf("Failed to derive config from defaults, %v", err)
	}

	_, err = NewImagesCacheFromConfig(cfg)

	if err != nil {
		t.Fatalf("Failed to derive images cache from defaults, %v", err)
	}
}

func TestNewDerivativesCacheFromConfig(t *testing.T) {

	r, err := iiifdefaults.FS.Open("config.json")

	if err != nil {
		t.Fatalf("Failed to load config from defaults, %v", err)
	}

	defer r.Close()

	cfg, err := iiifconfig.NewConfigFromReader(r)

	if err != nil {
		t.Fatalf("Failed to derive config fromd defaults, %v", err)
	}

	_, err = NewDerivativesCacheFromConfig(cfg)

	if err != nil {
		t.Fatalf("Failed to derive derivatives cache from defaults, %v", err)
	}
}
