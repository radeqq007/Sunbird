package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type PackageManager struct {
	config   *Config
	cacheDir string
}

func NewPackageManager() (*PackageManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	cacheDir := filepath.Join(homeDir, ".sunbird", "packages")

	// Create cache directory if it doesn't exist
	err = os.MkdirAll(cacheDir, 0755)
	if err != nil {
		return nil, err
	}

	return &PackageManager{
		cacheDir: cacheDir,
	}, nil
}

func (pm *PackageManager) Add(url string) error {
	pkgURL, err := ParsePackageURL(url)
	if err != nil {
		return err
	}

	fmt.Printf("Downloading %s/%s...\n", pkgURL.Owner, pkgURL.Repo)
	err = DownloadPackage(pkgURL, pm.cacheDir)
	if err != nil {
		return err
	}

	if _, err := os.Stat("sunbird.toml"); err == nil {
		err = pm.addDependencyToConfig(pkgURL)
		if err != nil {
			fmt.Printf("Warning: failed to add to sunbird.toml: %s\n", err)
		}
	}

	return nil
}

func (pm *PackageManager) addDependencyToConfig(url *PackageURL) error {
	config, err := LoadConfig("sunbird.toml")
	if err != nil {
		return err
	}

	depName := url.Repo

	// Initialize dependencies map if needed
	if config.Dependencies == nil {
		config.Dependencies = make(map[string]DependencyInfo)
	}

	config.Dependencies[depName] = DependencyInfo{
		Git:     url.GetGitURL(),
		Version: url.Version,
	}

	// Save the config back
	return pm.saveConfig(config)
}

func (pm *PackageManager) saveConfig(config *Config) error {
	// TODO: Implement proper TOML modification cause rn it just appends it at the end of the file
	file, err := os.OpenFile("sunbird.toml", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil
	}
	defer file.Close()

	data, err := os.ReadFile("sunbird.toml")
	if err != nil {
		return err
	}

	content := string(data)
	hasDependencies := strings.Contains(content, "[dependencies]")

	// If no [dependencies] section exists, add it
	if !hasDependencies {
		_, err = file.WriteString("\n[dependencies]\n")
		if err != nil {
			return err
		}
	}

	for name, dep := range config.Dependencies {
		// Check if this dependency already exists
		if strings.Contains(content, name+" =") {
			continue // Skip if already exists
		}

		var depLine string
		if dep.Version != "" {
			depLine = fmt.Sprintf("%s = { git = \"%s\", version = \"%s\" }\n",
				name, dep.Git, dep.Version)
		} else {
			depLine = fmt.Sprintf("%s = { git = \"%s\" }\n", name, dep.Git)
		}

		_, err = file.WriteString(depLine)
		if err != nil {
			return err
		}

		fmt.Printf("  ✓ Added %s to sunbird.toml\n", name)
	}

	return nil
}
