package pkg

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Package PackageInfo `toml:"package"`
}

type PackageInfo struct {
	Name         string           `toml:"name"`
	Version      string           `toml:"version"`
	Description  string           `toml:"description"`
	Authors      []string         `toml:"authors"`
	Main         string           `toml:"main"`
	Dependencies []DependencyInfo `toml:"dependencies"`
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

func SaveConfig(path string, config *Config) error {
	data, err := toml.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, data, 0o644)
	if err != nil {
		return err
	}

	return nil
}

func AddDependency(path, name string, dependency DependencyInfo) error {
	config, err := LoadConfig(path)
	if err != nil {
		return err
	}

	config.Package.Dependencies = append(config.Package.Dependencies, dependency)

	return SaveConfig(path, config)
}
