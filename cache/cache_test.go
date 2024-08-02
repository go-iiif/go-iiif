package cache

import (
	"testing"
)

func TestNewImagesCacheFromConfig(t testing.T) {

	path := "../docs/config.json.example"

	cfg, err := NewConfigFromFile(path)

	if err != nil {
		t.Fatalf("Failed to derive config from %s, %v", path, err)
	}

	_, err := NewImagesCacheFromConfig(cfg)

	if err != nil {
		t.Fatalf("Failed to derive images cache %s, %v", path, err)
	}
}

func TestNewDerivativesCacheFromConfig(t testing.T) {

	path := "../docs/config.json.example"

	cfg, err := NewConfigFromFile(path)

	if err != nil {
		t.Fatalf("Failed to derive config from %s, %v", path, err)
	}

	_, err := NewDerivativesCacheFromConfig(cfg)

	if err != nil {
		t.Fatalf("Failed to derive derivatives cache %s, %v", path, err)
	}
}
