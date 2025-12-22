package items

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRepository_SaveAndLoadCache(t *testing.T) {
	// Create temporary directory for test
	tmpDir, err := os.MkdirTemp("", "arc-scanner-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cachePath := filepath.Join(tmpDir, "items.json")
	repo := NewRepository(cachePath)

	// Test items
	items := []Item{
		{ID: "test-item-1", Name: "Test Item 1", Value: 100},
		{ID: "test-item-2", Name: "Test Item 2", Value: 200},
	}

	// Test SaveToCache
	if err := repo.SaveToCache(items); err != nil {
		t.Fatalf("SaveToCache failed: %v", err)
	}

	// Verify file was created
	if !repo.CacheExists() {
		t.Error("CacheExists returned false after SaveToCache")
	}

	// Test LoadFromCache
	loaded, err := repo.LoadFromCache()
	if err != nil {
		t.Fatalf("LoadFromCache failed: %v", err)
	}

	if len(loaded) != len(items) {
		t.Errorf("LoadFromCache returned %d items, want %d", len(loaded), len(items))
	}

	// Verify item content
	for i, item := range loaded {
		if item.ID != items[i].ID {
			t.Errorf("item[%d].ID = %s, want %s", i, item.ID, items[i].ID)
		}
		if item.Value != items[i].Value {
			t.Errorf("item[%d].Value = %d, want %d", i, item.Value, items[i].Value)
		}
	}
}

func TestRepository_CacheNotExists(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "arc-scanner-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cachePath := filepath.Join(tmpDir, "nonexistent.json")
	repo := NewRepository(cachePath)

	if repo.CacheExists() {
		t.Error("CacheExists returned true for nonexistent file")
	}

	_, err = repo.LoadFromCache()
	if err == nil {
		t.Error("LoadFromCache should return error for nonexistent file")
	}
}

func TestRepository_CorruptedCache(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "arc-scanner-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cachePath := filepath.Join(tmpDir, "corrupted.json")

	// Write invalid JSON
	if err := os.WriteFile(cachePath, []byte("not valid json"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	repo := NewRepository(cachePath)

	_, err = repo.LoadFromCache()
	if err == nil {
		t.Error("LoadFromCache should return error for corrupted file")
	}
}

func TestRepository_CachePath(t *testing.T) {
	expectedPath := "/some/path/items.json"
	repo := NewRepository(expectedPath)

	if repo.CachePath() != expectedPath {
		t.Errorf("CachePath() = %s, want %s", repo.CachePath(), expectedPath)
	}
}

func TestRepository_SaveToCache_WithComplexItems(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "arc-scanner-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	cachePath := filepath.Join(tmpDir, "items.json")
	repo := NewRepository(cachePath)

	// Create items with nested structures
	recycleComponents := []RecycleEntry{
		{Quantity: 2, Component: Component{ID: "component-1", Name: "Component 1"}},
		{Quantity: 5, Component: Component{ID: "component-2", Name: "Component 2"}},
	}

	items := []Item{
		{
			ID:                "complex-item",
			Name:              "Complex Item",
			Value:             500,
			Icon:              "https://example.com/icon.png",
			RecycleComponents: &recycleComponents,
		},
	}

	// Save and reload
	if err := repo.SaveToCache(items); err != nil {
		t.Fatalf("SaveToCache failed: %v", err)
	}

	loaded, err := repo.LoadFromCache()
	if err != nil {
		t.Fatalf("LoadFromCache failed: %v", err)
	}

	if len(loaded) != 1 {
		t.Fatalf("LoadFromCache returned %d items, want 1", len(loaded))
	}

	// Verify nested structure
	if loaded[0].RecycleComponents == nil {
		t.Fatal("RecycleComponents is nil after reload")
	}

	components := *loaded[0].RecycleComponents
	if len(components) != 2 {
		t.Errorf("RecycleComponents has %d entries, want 2", len(components))
	}

	if components[0].Component.ID != "component-1" {
		t.Errorf("First component ID = %s, want component-1", components[0].Component.ID)
	}
}
