package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

type Config struct {
	Level       LevelConfig       `json:"level"`
	Profile     ProfileConfig     `json:"profile"`
	Graphics    GraphicsConfig    `json:"graphics"`
	Features    FeaturesConfig    `json:"features"`
	Images      ImagesConfig      `json:"images"`
	Derivatives DerivativesConfig `json:"derivatives"`
	Flickr      FlickrConfig      `json:"flickr,omitempty"`
	Primitive   PrimitiveConfig   `json:"primitive,omitempty"`
	Palette     PaletteConfig     `json:"palette,omitempty"`
	Custom      interface{}       `json:"custom,omitempty"`
}

type ProfileConfig struct {
	Services ServicesConfig `json:"services"`
}

type ServicesConfig struct {
	Enable ServicesToggle `json:"enable"`
}

type ServicesToggle []string

type PaletteConfig struct {
	Extruder SourceConfig   `json:"extruder"`
	Grid     SourceConfig   `json:"grid"`
	Palettes []SourceConfig `json:"palettes"`
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

type FlickrConfig struct {
	ApiKey    string `json:"apikey"`
	ApiSecret string `json:"apisecret,omitempty"`
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

func NewConfigFromFile(file string) (*Config, error) {

	body, err := ioutil.ReadFile(file)

	if err != nil {
		return nil, err
	}

	return NewConfigFromBytes(body)
}

func NewConfigFromEnv(name string) (*Config, error) {

	env, ok := os.LookupEnv(name)

	if !ok {
		return nil, errors.New("Missing environment variable by that name")
	}

	return NewConfigFromBytes([]byte(env))
}

func NewConfigFromBytes(body []byte) (*Config, error) {

	c := Config{}

	err := json.Unmarshal(body, &c)

	if err != nil {
		return nil, err
	}

	return &c, nil
}
