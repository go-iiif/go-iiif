package config

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/aaronland/gocloud-blob/bucket"
	iiifdefaults "github.com/go-iiif/go-iiif/v6/defaults"
	"github.com/sfomuseum/go-flags/lookup"
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

// GraphicsConfig
type GraphicsConfig struct {
	Source SourceConfig `json:"source"`
}

type SourceConfig struct {
	// A valid go-iiif/cache.Cache URI. If empty this value will be derived from the other values in CacheConfig.
	URI string `json:"uri"`
	// The name of the source provider. Deprecated; Please use the URI property instead.
	Name string `json:"name"`
	// Path information for the source provider. Deprecated; Please use the URI property instead.
	Path string `json:"path,omitempty"`
	// Path information for the S3 (blob) source provider. Deprecated; Please use the URI property instead.
	Prefix string `json:"prefix,omitempty"`
	// Region information for the S3 (blob) source provider. Deprecated; Please use the URI property instead.
	Region string `json:"region,omitempty"`
	// Credentials information for the S3 (blob) source provider. Deprecated; Please use the URI property instead.
	Credentials string `json:"credentials,omitempty"`
	Tmpdir      string `json:"tmpdir,omitempty"`
	Count       int    `json:"count,omitempty"` // used by PaletteConfig.Extruder
}

// FlickrConfig defines configuration details for interacting with the Flickr photo-sharing service.
type FlickrConfig struct {
	// A valid `aaronland/go-flickr-api.Client` URI.
	ClientURI string `json:"client_uri"`
}

// PrimitiveConfig defines configuration details for using the `fogleman/primitive` package.
type PrimitiveConfig struct {
	// MaxIterations is the maximum number of iterations for perform when generating `fogleman/primitive` images.
	MaxIterations int `json:"max_iterations"`
}

// CacheConfig defines configuration details for caching objects.
type CacheConfig struct {
	// A valid go-iiif/cache.Cache URI. If empty this value will be derived from the other values in CacheConfig.
	URI string `json:"uri"`
	// The name of the caching source. Deprecated; Please use the URI property instead.
	Name string `json:"name"`
	// Path information for the caching source. Deprecated; Please use the URI property instead.
	Path string `json:"path,omitempty"`
	// Time-to-live information for the memory caching source. Deprecated; Please use the URI property instead.
	TTL int `json:"ttl,omitempty"`
	// Cache limit information for the memory caching source. Deprecated; Please use the URI property instead.
	Limit int `json:"limit,omitempty"`
	// Path information for the S3 (blob) caching source. Deprecated; Please use the URI property instead.
	Prefix string `json:"prefix,omitempty"`
	// Region information for the S3 (blob) caching source. Deprecated; Please use the URI property instead.
	Region string `json:"region,omitempty"`
	// Credentials information for the S3 (blob) caching source. Deprecated; Please use the URI property instead.
	Credentials string `json:"credentials,omitempty"`
}

// NewConfigFromFlag is DEPRECATED and will simply hand off to the `NewConfigFromString` method.
func NewConfigFromFlag(flag string) (*Config, error) {
	slog.Warn("NewConfigFromFlag has been DEPRECATED. Please use NewConfigFromString instead.")
	return NewConfigFromString(flag)
}

// NewConfigFromString returns a new `Config` instance derived from 'str'. If 'str' starts with "env:" then the remainder
// of the string will be used as the environment variable to derive config data from. If 'str' starts with "{" then the entire
// string will be used to derive config data from. Otherwise 'str' is assumed to be a local file on disk config configuration
// data.
func NewConfigFromString(str string) (*Config, error) {

	if strings.HasPrefix(str, "env:") {

		env := strings.Replace(str, "env:", "", 1)
		env = strings.Trim(env, " ")

		if env == "" {
			return nil, errors.New("Invalid environment variable")
		}

		return NewConfigFromEnv(env)
	}

	if strings.HasPrefix(str, "{") {
		return NewConfigFromBytes([]byte(str))
	}

	return NewConfigFromFile(str)
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
		return nil, fmt.Errorf("Failed to read config body, %w", err)
	}

	return NewConfigFromBytes(body)
}

// NewConfigFromFile returns a new `Config` instance derived the key 'key' in the `gocloud.dev/blob.Bucket` identified by 'bucket'.
func NewConfigFromBucket(ctx context.Context, bucket *blob.Bucket, key string) (*Config, error) {

	r, err := bucket.NewReader(ctx, key, nil)

	if err != nil {
		return nil, fmt.Errorf("Failed to create new reader for '%s', %w", key, err)
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
		return nil, fmt.Errorf("Failed to unmarshal config, %w", err)
	}

	return &c, nil
}

func LoadConfig(ctx context.Context, bucket_uri string, key string) (*Config, error) {

	if bucket_uri == iiifdefaults.URI {

		key = "config.json"

		r, err := iiifdefaults.FS.Open(key)

		if err != nil {
			return nil, fmt.Errorf("Failed to load config (%s) from defaults, %w", key, err)
		}

		return NewConfigFromReader(r)
	}

	config_bucket, err := bucket.OpenBucket(ctx, bucket_uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to open config bucket, %w", err)
	}

	defer config_bucket.Close()

	return NewConfigFromBucket(ctx, config_bucket, key)
}

func LoadConfigWithFlagSet(ctx context.Context, fs *flag.FlagSet) (*Config, error) {

	config_source, err := lookup.StringVar(fs, "config-source")

	if err != nil {
		return nil, fmt.Errorf("Failed to lookup -config-source flag, %w", err)
	}

	config_name, err := lookup.StringVar(fs, "config-name")

	if err != nil {
		return nil, fmt.Errorf("Failed to lookup -config-name flag, %w", err)
	}

	cfg, err := LoadConfig(ctx, config_source, config_name)

	if err != nil {
		return nil, err
	}

	images_source_uri, _ := lookup.StringVar(fs, "config-images-source-uri")

	if images_source_uri != "" {
		slog.Debug("Reassign images source", "uri", images_source_uri)
		cfg.Images.Source.URI = images_source_uri
	}

	derivatives_cache_uri, _ := lookup.StringVar(fs, "config-derivatives-cache-uri")

	if derivatives_cache_uri != "" {
		slog.Debug("Reassign derivatives cache", "uri", derivatives_cache_uri)
		cfg.Derivatives.Cache.URI = derivatives_cache_uri
	}

	return cfg, nil
}
