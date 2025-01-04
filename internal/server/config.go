package server

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DBPaths struct {
		Base    string `yaml:"base"`
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

func (c Config) SongPath(lang, artist, title string) string {
	return fmt.Sprintf("%s/%s/%s/%s/%s.yaml", c.DBPaths.Base, lang, c.DBPaths.Songs.Base, artist, title)
}

func (c Config) SongBasePath(lang string) string {
	return fmt.Sprintf("%s/%s/%s", c.DBPaths.Base, lang, c.DBPaths.Songs.Base)
}

func (c Config) IndexedSongsPath(lang string) string {
	return fmt.Sprintf("%s/%s/%s", c.DBPaths.Base, lang, c.DBPaths.Songs.Indexed)
}
