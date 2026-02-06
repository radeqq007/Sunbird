package pkg

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Package      PackageInfo               `toml:"package"`
	Dependencies map[string]DependencyInfo `toml:"dependencies"`
}

type PackageInfo struct {
	Name        string   `toml:"name"`
	Version     string   `toml:"version"`
	Description string   `toml:"description"`
	Authors     []string `toml:"authors"`
	Main        string   `toml:"main"`
}

type DependencyInfo struct {
	Git     string `toml:"git"`
	Version string `toml:"version"`
	Tag     string `toml:"tag"`
	Path    string `toml:"path"`
}

func LoadConfig(path string) (*Config, error) {
	var config Config

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = toml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
