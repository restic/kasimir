package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	GOOS   []string `yaml:"goos"`
	GOARCH []string `yaml:"goarch"`
	GOARM  []string `yaml:"goarm"`
}

var DefaultConfig = Config{
	GOOS:   []string{"linux", "darwin"},
	GOARCH: []string{"386", "amd64", "arm64"},
	GOARM:  []string{"6", "7"},
}

// LoadConfig loads the config from filename.
func LoadConfig(filename string) (Config, error) {
	errwrap := func(err error) error {
		return fmt.Errorf("error loading config file %v: %w", filename, err)
	}

	var cfg Config

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return Config{}, errwrap(err)
	}

	err = yaml.UnmarshalStrict(buf, &cfg)
	if err != nil {
		return Config{}, errwrap(err)
	}

	// set default values
	if cfg.GOARCH == nil {
		cfg.GOARCH = DefaultConfig.GOARCH
	}
	if cfg.GOOS == nil {
		cfg.GOOS = DefaultConfig.GOOS
	}
	if cfg.GOARM == nil {
		cfg.GOARM = DefaultConfig.GOARM
	}

	return cfg, nil
}
