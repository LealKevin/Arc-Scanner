package items

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"arc-scanner/internal/config"
)

// Repository handles item data persistence and retrieval.
type Repository struct {
	cachePath string
}

// NewRepository creates a new Repository with the given cache path.
func NewRepository(cachePath string) *Repository {
	return &Repository{
		cachePath: cachePath,
	}
}

// FetchFromAPI retrieves items from the MetaForge API.
// Returns the full list of items or an error if the API is unavailable.
func (r *Repository) FetchFromAPI() ([]Item, error) {
	var allItems []Item

	slog.Info("fetching items from API")

	for page := 1; page <= config.APIPageCount; page++ {
		url := fmt.Sprintf("%s?minimal=true&includeComponents=true&limit=%d&page=%d",
			config.MetaForgeAPIBase, config.APIPageSize, page)

		slog.Debug("fetching page", "page", page, "url", url)

		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrAPIUnavailable, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("%w: status %d", ErrAPIUnavailable, resp.StatusCode)
		}

		var response Response
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, fmt.Errorf("failed to decode API response: %w", err)
		}

		allItems = append(allItems, response.Data...)
	}

	slog.Info("items fetched from API", "count", len(allItems))
	return allItems, nil
}

// LoadFromCache reads items from the local cache file.
// Returns the items or an error if the cache is missing or corrupted.
func (r *Repository) LoadFromCache() ([]Item, error) {
	data, err := os.ReadFile(r.cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err // Let caller handle missing file
		}
		return nil, fmt.Errorf("failed to read cache: %w", err)
	}

	var items []Item
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCacheCorrupted, err)
	}

	slog.Info("items loaded from cache", "count", len(items))
	return items, nil
}

// SaveToCache writes items to the local cache file.
func (r *Repository) SaveToCache(items []Item) error {
	file, err := os.Create(r.cachePath)
	if err != nil {
		return fmt.Errorf("failed to create cache file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(items); err != nil {
		return fmt.Errorf("failed to encode items: %w", err)
	}

	slog.Info("items saved to cache", "path", r.cachePath, "count", len(items))
	return nil
}

// CachePath returns the path to the cache file.
func (r *Repository) CachePath() string {
	return r.cachePath
}

// CacheExists returns true if the cache file exists.
func (r *Repository) CacheExists() bool {
	_, err := os.Stat(r.cachePath)
	return err == nil
}
