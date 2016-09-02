package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
        Graphics    GraphicsConfig    `json:"graphics"`
	Images      ImagesConfig      `json:"images"`
	Derivatives DerivativesConfig `json:"derivatives"`
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
	Name string `json:"name"`
	Path string `json:"path,omitempty"`
}

type CacheConfig struct {
	Name string `json:"name"`
	Path string `json:"path,omitempty"`
	TTL int `json:"ttl,omitempty"`
	Limit int `json:"limit,omitempty"`
}

func NewConfigFromFile(file string) (*Config, error) {

	body, err := ioutil.ReadFile(file)

	if err != nil {
		return nil, err
	}

	c := Config{}
	err = json.Unmarshal(body, &c)

	if err != nil {
		return nil, err
	}

	return &c, nil
}
