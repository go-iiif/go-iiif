package config

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func TestNewConfigFromFile(t *testing.T) {

	path := "../docs/config.json.example"

	_, err := NewConfigFromFile(path)

	if err != nil {
		t.Fatalf("Failed to derive config from %s, %v", path, err)
	}
}

func TestNewConfigFromReader(t *testing.T) {

	path := "../docs/config.json.example"

	r, err := os.Open(path)

	if err != nil {
		t.Fatalf("Failed to open %s for reading, %v", path, err)
	}

	defer r.Close()

	_, err = NewConfigFromReader(r)

	if err != nil {
		t.Fatalf("Failed to derive config from %s, %v", path, err)
	}
}

func TestNewConfigFromBytes(t *testing.T) {

	path := "../docs/config.json.example"

	r, err := os.Open(path)

	if err != nil {
		t.Fatalf("Failed to open %s for reading, %v", path, err)
	}

	defer r.Close()

	body, err := io.ReadAll(r)

	if err != nil {
		t.Fatalf("Failed to read %s, %v", path, err)
	}

	_, err = NewConfigFromBytes(body)

	if err != nil {
		t.Fatalf("Failed to derive config from %s, %v", path, err)
	}
}

func TestNewConfigFromEnv(t *testing.T) {

	path := "../docs/config.json.example"

	r, err := os.Open(path)

	if err != nil {
		t.Fatalf("Failed to open %s for reading, %v", path, err)
	}

	defer r.Close()

	body, err := io.ReadAll(r)

	if err != nil {
		t.Fatalf("Failed to read %s, %v", path, err)
	}

	env_var := "IIIF_CONFIG"

	err = os.Setenv(env_var, string(body))

	if err != nil {
		t.Fatalf("Failed to assign %s environment variable, %v", env_var, err)
	}

	str_config := os.Getenv(env_var)

	if str_config == "" {
		t.Fatalf("Environment variable %s is empty", env_var)
	}

	_, err = NewConfigFromEnv(env_var)

	if err != nil {
		t.Fatalf("Failed to derive config from environment variable %s, %v", env_var, err)
	}
}

func TestNewConfigFromString(t *testing.T) {

	path := "../docs/config.json.example"

	r, err := os.Open(path)

	if err != nil {
		t.Fatalf("Failed to open %s for reading, %v", path, err)
	}

	defer r.Close()

	body, err := io.ReadAll(r)

	if err != nil {
		t.Fatalf("Failed to read %s, %v", path, err)
	}

	env_var := "IIIF_CONFIG"

	err = os.Setenv(env_var, string(body))

	if err != nil {
		t.Fatalf("Failed to assign %s environment variable, %v", env_var, err)
	}

	str_config := os.Getenv(env_var)

	if str_config == "" {
		t.Fatalf("Environment variable %s is empty", env_var)
	}

	flags := []string{
		fmt.Sprintf("env:%s", env_var),
		str_config,
		path,
	}

	for _, fl := range flags {

		_, err = NewConfigFromString(fl)

		if err != nil {
			t.Fatalf("Failed to derive config from flag '%s', %v", fl, err)
		}
	}

}
