package items

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"arc-scanner/internal/config"
)

type Repository struct {
	cachePath string
}

func NewRepository(cachePath string) *Repository {
	return &Repository{
		cachePath: cachePath,
	}
}

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

func (r *Repository) CachePath() string {
	return r.cachePath
}

func (r *Repository) CacheExists() bool {
	_, err := os.Stat(r.cachePath)
	return err == nil
}
