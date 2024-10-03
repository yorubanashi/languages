package server

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DBPaths struct {
		Artists string `yaml:"artists"`
		Songs   struct {
			Base    string `yaml:"base"`
			Indexed string `yaml:"indexed"`
		} `yaml:"songs"`
	} `yaml:"dbPaths"`

	Server struct {
		Address  string `yaml:"address"`
		Timeouts struct {
			Shutdown int `yaml:"shutdown"`
		} `yaml:"timeouts"`
	} `yaml:"server"`

	StartOptions struct {
		Index bool `yaml:"index"`
	} `yaml:"startOptions"`

	Svelte struct {
		Pages string `yaml:"pages"`
	} `yaml:"svelte"`
}

func LoadConfig(path string) (*Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = yaml.Unmarshal(bytes, cfg)
	return cfg, err
}
