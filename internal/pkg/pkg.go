package pkg

import (
	"fmt"
)

const dependencyDirectory = ".sb_modules/"

type PackageManager struct {
	config *Config
}

func NewPackageManager() *PackageManager {
	return &PackageManager{}
}

func (pm *PackageManager) Add(url string) error {
	pkgURL, err := ParsePackageURL(url)
	if err != nil {
		return err
	}

	fmt.Printf("Downloading %s/%s...\n", pkgURL.Owner, pkgURL.Repo)
	err = DownloadPackage(pkgURL)
	if err != nil {
		return err
	}

	// if _, err := os.Stat("sunbird.toml"); err == nil {
	// 	err = pm.addDependencyToConfig(pkgURL)
	// 	if err != nil {
	// 		fmt.Printf("Warning: failed to add to sunbird.toml: %s\n", err)
	// 	}
	// }

	return nil
}
