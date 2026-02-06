package pkg

import (
	"fmt"
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

func DownloadPackage(pkgURL *PackageURL, cacheDir string) error {
	cachePath := pkgURL.GetCachePath(cacheDir)

	// Check if already downloaded
	if _, err := os.Stat(cachePath); err == nil {
		fmt.Printf("  Package already cached at: %s\n", cachePath)
		return nil
	}

	// Create parent directories
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

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

	// Clone the repository
	_, err := git.PlainClone(cachePath, false, cloneOptions)
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	fmt.Printf("  ✓ Downloaded to: %s\n", cachePath)
	return nil
}
