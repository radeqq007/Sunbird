package pkg

import (
	"fmt"
)

const dependencyDirectory = ".sb_modules/"

type PackageManager struct{}

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

	return nil
}
