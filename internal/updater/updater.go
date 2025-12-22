package updater

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/blang/semver"
)

// UpdateInfo contains information about an available update
type UpdateInfo struct {
	Version     string `json:"version"`
	URL         string `json:"url"`
	DownloadURL string `json:"downloadUrl"`
	ReleaseNotes string `json:"releaseNotes"`
	PublishedAt string `json:"publishedAt"`
}

// Updater handles checking for and applying updates
type Updater struct {
	owner          string
	repo           string
	currentVersion string
	downloadedPath string
	httpClient     *http.Client
}

// GitHubRelease represents a GitHub release API response
type GitHubRelease struct {
	TagName     string        `json:"tag_name"`
	HTMLURL     string        `json:"html_url"`
	Body        string        `json:"body"`
	PublishedAt string        `json:"published_at"`
	Assets      []GitHubAsset `json:"assets"`
}

// GitHubAsset represents an asset in a GitHub release
type GitHubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

// New creates a new Updater instance
func New(owner, repo, currentVersion string) *Updater {
	return &Updater{
		owner:          owner,
		repo:           repo,
		currentVersion: currentVersion,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CheckForUpdate checks GitHub for a newer version
func (u *Updater) CheckForUpdate() (*UpdateInfo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", u.owner, u.repo)
	fmt.Printf("DEBUG: checking URL: %s\n", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "arc-scanner-updater")

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release: %w", err)
	}
	defer resp.Body.Close()

	fmt.Printf("DEBUG: response status: %d\n", resp.StatusCode)

	if resp.StatusCode == http.StatusNotFound {
		// No releases yet
		fmt.Println("DEBUG: 404 - no releases found")
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Parse versions for comparison
	latestVersion := strings.TrimPrefix(release.TagName, "v")
	currentVersion := strings.TrimPrefix(u.currentVersion, "v")

	fmt.Printf("DEBUG: latestVersion=%s, currentVersion=%s\n", latestVersion, currentVersion)

	latest, err := semver.Parse(latestVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to parse latest version '%s': %w", latestVersion, err)
	}

	current, err := semver.Parse(currentVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to parse current version '%s': %w", currentVersion, err)
	}

	fmt.Printf("DEBUG: latest=%v, current=%v, current.GTE(latest)=%v\n", latest, current, current.GTE(latest))

	// No update if current >= latest
	if current.GTE(latest) {
		return nil, nil
	}

	// Find the appropriate asset for this platform
	downloadURL := u.findAssetURL(release.Assets)
	if downloadURL == "" {
		return nil, fmt.Errorf("no compatible asset found for %s", runtime.GOOS)
	}

	return &UpdateInfo{
		Version:      latestVersion,
		URL:          release.HTMLURL,
		DownloadURL:  downloadURL,
		ReleaseNotes: release.Body,
		PublishedAt:  release.PublishedAt,
	}, nil
}

// findAssetURL finds the download URL for the current platform
func (u *Updater) findAssetURL(assets []GitHubAsset) string {
	var platformKey string
	switch runtime.GOOS {
	case "darwin":
		platformKey = "macos"
	case "windows":
		platformKey = "windows"
	default:
		return ""
	}

	for _, asset := range assets {
		name := strings.ToLower(asset.Name)
		if strings.Contains(name, platformKey) && strings.HasSuffix(name, ".zip") {
			return asset.BrowserDownloadURL
		}
	}
	return ""
}

// DownloadUpdate downloads the update to a temporary location
// progressFn is called with progress percentage (0-100)
func (u *Updater) DownloadUpdate(ctx context.Context, info *UpdateInfo, progressFn func(percent int)) error {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "arc-scanner-update-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	zipPath := filepath.Join(tempDir, "update.zip")

	// Download the file
	req, err := http.NewRequestWithContext(ctx, "GET", info.DownloadURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// Create output file
	out, err := os.Create(zipPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	// Download with progress tracking
	totalSize := resp.ContentLength
	var downloaded int64

	buf := make([]byte, 32*1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			_, writeErr := out.Write(buf[:n])
			if writeErr != nil {
				return fmt.Errorf("failed to write to file: %w", writeErr)
			}
			downloaded += int64(n)

			if totalSize > 0 && progressFn != nil {
				percent := int(float64(downloaded) / float64(totalSize) * 100)
				progressFn(percent)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}
	}

	// Extract the zip
	extractDir := filepath.Join(tempDir, "extracted")
	if err := u.extractZip(zipPath, extractDir); err != nil {
		return fmt.Errorf("failed to extract update: %w", err)
	}

	u.downloadedPath = extractDir
	return nil
}

// extractZip extracts a zip file to the destination directory
func (u *Updater) extractZip(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	for _, f := range r.File {
		fpath := filepath.Join(destDir, f.Name)

		// Security check for zip slip vulnerability
		if !strings.HasPrefix(fpath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), 0755); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

// GetDownloadedPath returns the path where the update was extracted
func (u *Updater) GetDownloadedPath() string {
	return u.downloadedPath
}

// GetCurrentVersion returns the current app version
func (u *Updater) GetCurrentVersion() string {
	return u.currentVersion
}
