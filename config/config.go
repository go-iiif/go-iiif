package config

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"

	"gocloud.dev/blob"
)

// type Config is a struct containing configuration details for IIIF processes and services.
type Config struct {
	// Level is a `LevelConfig` instance detailing the IIIF level in use.
	Level LevelConfig `json:"level"`
	// Profile is a `ProfileConfig` instance detailing the IIIF profile in use.
	Profile ProfileConfig `json:"profile"`
	// Graphics is a `GraphicsConfig` instance detailing the graphics processor used for IIIF processes.
	Graphics GraphicsConfig `json:"graphics"`
	// Features is a `ProfileConfig` instance detailing the IIIF features in use.
	Features FeaturesConfig `json:"features"`
	// Images is a `ImagesConfig` detailing how and where IIIF source images are stored.
	Images ImagesConfig `json:"images"`
	// Derivatives is a `DerivativesConfig` detailing how and where IIIF derivative images are stored.
	Derivatives DerivativesConfig `json:"derivatives"`
	Flickr      FlickrConfig      `json:"flickr,omitempty"`
	Primitive   PrimitiveConfig   `json:"primitive,omitempty"`
	Palette     PaletteConfig     `json:"palette,omitempty"`
	BlurHash    BlurHashConfig    `json:"blurhash,omitempty"`
	ImageHash   ImageHashConfig   `json:"imagehash,omitempty"`
	Custom      interface{}       `json:"custom,omitempty"`
}

// ProfileConfig defines configuration details for the IIIF profile in use.
type ProfileConfig struct {
	// Services is a `ServicesConfig` instance detailing IIIF services in use.
	Services ServicesConfig `json:"services"`
}

// ServicesConfig defines configuration details for the IIIF services in use.
type ServicesConfig struct {
	// Enable is a list of `ServiceToggle` instance to enable for IIIF processing.
	Enable ServicesToggle `json:"enable"`
}

type ServicesToggle []string

// PaletteConfig details configuration details for colour palette extraction services.
type PaletteConfig struct {
	Extruder SourceConfig   `json:"extruder"`
	Grid     SourceConfig   `json:"grid"`
	Palettes []SourceConfig `json:"palettes"`
}

// BlurHashConfig defines configuration details for blurhash generation services.
type BlurHashConfig struct {
	X    int `json:"x"`
	Y    int `json:"y"`
	Size int `json:"size"`
}

type ImageHashConfig struct {
}

type LevelConfig struct {
	Compliance string `json:"compliance"`
}

type FeaturesConfig struct {
	Enable  FeaturesToggle `json:"enable"`
	Disable FeaturesToggle `json:"disable"`
	Append  FeaturesAppend `json:"append"`
}

type FeaturesToggle map[string][]string

type FeaturesAppend map[string]map[string]FeaturesDetails

type FeaturesDetails struct {
	Syntax    string `json:"syntax"`
	Required  bool   `json:"required"`
	Supported bool   `json:"supported"`
	Match     string `json:"match,omitempty"`
}

type ImagesConfig struct {
	Source SourceConfig `json:"source"`
	Cache  CacheConfig  `json:"cache"`
}

type DerivativesConfig struct {
	Cache CacheConfig `json:"cache"`
}

type GraphicsConfig struct {
	Source SourceConfig `json:"source"`
}

type SourceConfig struct {
	Name        string `json:"name"`
	Path        string `json:"path,omitempty"`
	Prefix      string `json:"prefix,omitempty"`
	Region      string `json:"region,omitempty"`
	Credentials string `json:"credentials,omitempty"`
	Tmpdir      string `json:"tmpdir,omitempty"`
	Count       int    `json:"count,omitempty"` // used by PaletteConfig.Extruder
}

// FlickrConfig defines confiruation
type FlickrConfig struct {
	// A valid `aaronland/go-flickr-api.Client` URI.
	ClientURI string `json:"client_uri"`
}

type PrimitiveConfig struct {
	MaxIterations int `json:"max_iterations"`
}

type CacheConfig struct {
	Name        string `json:"name"`
	Path        string `json:"path,omitempty"`
	TTL         int    `json:"ttl,omitempty"`
	Limit       int    `json:"limit,omitempty"`
	Prefix      string `json:"prefix,omitempty"`
	Region      string `json:"region,omitempty"`
	Credentials string `json:"credentials,omitempty"`
}

func NewConfigFromFlag(flag string) (*Config, error) {

	if strings.HasPrefix(flag, "env:") {

		env := strings.Replace(flag, "env:", "", 1)
		env = strings.Trim(env, " ")

		if env == "" {
			return nil, errors.New("Invalid environment variable")
		}

		return NewConfigFromEnv(env)
	}

	if strings.HasPrefix(flag, "{") {
		return NewConfigFromBytes([]byte(flag))
	}

	return NewConfigFromFile(flag)
}

// NewConfigFromFile returns a new `Config` instance derived from 'file' which is assumed to be a local file on disk.
func NewConfigFromFile(file string) (*Config, error) {

	body, err := os.ReadFile(file)

	if err != nil {
		return nil, err
	}

	return NewConfigFromBytes(body)
}

// NewConfigFromReader returns a new `Config` instance derived from 'r'.
func NewConfigFromReader(r io.Reader) (*Config, error) {

	body, err := io.ReadAll(r)

	if err != nil {
		return nil, err
	}

	return NewConfigFromBytes(body)
}

// NewConfigFromFile returns a new `Config` instance derived the key 'key' in the `gocloud.dev/blob.Bucket` identified by 'bucket'.
func NewConfigFromBucket(ctx context.Context, bucket *blob.Bucket, key string) (*Config, error) {

	r, err := bucket.NewReader(ctx, key, nil)

	if err != nil {
		return nil, err
	}

	defer r.Close()

	return NewConfigFromReader(r)
}

// NewConfigFromFile returns a new `Config` instance derived from the environment variable 'name'.
func NewConfigFromEnv(name string) (*Config, error) {

	env, ok := os.LookupEnv(name)

	if !ok {
		return nil, errors.New("Missing environment variable by that name")
	}

	return NewConfigFromBytes([]byte(env))
}

// NewConfigFromFile returns a new `Config` instance derived from 'body'.
func NewConfigFromBytes(body []byte) (*Config, error) {

	c := Config{}

	err := json.Unmarshal(body, &c)

	if err != nil {
		return nil, err
	}

	return &c, nil
}
