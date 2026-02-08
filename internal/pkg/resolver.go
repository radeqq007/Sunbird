package pkg

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type PackageURL struct {
	Host    string // eg. github.com
	Owner   string
	Repo    string
	Version string
	IsTag   bool
}

func ParsePackageURL(url string) (*PackageURL, error) {
	parts := strings.Split(url, "@")

	var path, version string
	if len(parts) == 2 {
		path = parts[0]
		version = parts[1]
	} else if len(parts) == 1 {
		path = parts[0]
		version = "" // latest
	} else {
		return nil, fmt.Errorf("invalid package URL format: %s", url)
	}

	// Parse the path: host/owner/repo
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 3 {
		return nil, fmt.Errorf("invalid package path: %s (expected host/owner/repo)", path)
	}

	return &PackageURL{
		Host:    pathParts[0],
		Owner:   pathParts[1],
		Repo:    pathParts[2],
		Version: version,
		IsTag:   strings.HasPrefix(version, "v"),
	}, nil
}

func (p *PackageURL) GetGitURL() string {
	return fmt.Sprintf("https://%s/%s/%s.git", p.Host, p.Owner, p.Repo)
}

func (p *PackageURL) GetCachePath(cacheDir string) string {
	versionSuffix := ""
	if p.Version != "" {
		versionSuffix = "@" + p.Version
	} else {
		versionSuffix = "@latest"
	}

	return filepath.Join(cacheDir, p.Host, p.Owner, p.Repo+versionSuffix)
}

func DownloadPackage(pkgURL *PackageURL) error {
	gitURL := pkgURL.GetGitURL()
	fmt.Printf("  Cloning from: %s\n", gitURL)

	cloneOptions := &git.CloneOptions{
		URL:      gitURL,
		Progress: os.Stdout,
	}

	// If a specific version is specified
	if pkgURL.Version != "" {
		if pkgURL.IsTag {
			// Clone specific tag
			cloneOptions.ReferenceName = plumbing.NewTagReferenceName(pkgURL.Version)
			cloneOptions.SingleBranch = true
		} else {
			// Clone specific branch
			cloneOptions.ReferenceName = plumbing.NewBranchReferenceName(pkgURL.Version)
			cloneOptions.SingleBranch = true
		}
	}


	tempDir := filepath.Join(os.TempDir(), "sunbird_cache", pkgURL.Owner, pkgURL.Repo)

	if _, err := os.Stat(tempDir); err == nil {
		err = os.RemoveAll(tempDir)
		if err != nil {
			return fmt.Errorf("failed to remove existing temp directory: %w", err)
		}
	}

	// Clone the repository
	_, err := git.PlainClone(tempDir, false, cloneOptions)
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	// Read config from the cloned package, get the name and move the package to .sunbird/host/owner/repo/version
	configPath := filepath.Join(tempDir, "sunbird.toml")
	config, err := LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("Failed to load config from downloaded package: %w", err)
	}

	packageName := config.Package.Name
	if packageName == "" {
		return fmt.Errorf("Package name is missing in sunbird.toml")
	}

	finalPath := filepath.Join(dependencyDirectory, packageName)

	// Remove existing directory if it exists
	if _, err :=
		os.Stat(finalPath); err == nil {
		err = os.RemoveAll(finalPath)
		if err != nil {
			return fmt.Errorf("failed to remove existing directory: %w", err)
		}
	}

	err = os.MkdirAll(filepath.Dir(finalPath), 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory for package: %w", err)
	}

	err = moveDirectory(tempDir, finalPath)
	if err != nil {
		return fmt.Errorf("failed to move package to final location: %w", err)
	}

	fmt.Printf("  ✓ Downloaded\n")
	return nil
}

func moveDirectory(src, dst string) error {
	// Try renaming first (cross-device moves will fail)
	err := os.Rename(src, dst)
	if err == nil {
		return nil
	}

	// Fallback to copying the directory
	return moveViaCopy(src, dst)
}

func moveViaCopy(src, dst string) error {
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Determine the destination path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		// It's a file, copy it
		return copyFile(path, targetPath)
	})

	if err != nil {
		return err
	}

	// Clean up the source
	return os.RemoveAll(src)
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, info.Mode())
}