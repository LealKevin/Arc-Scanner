package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"arc-scanner/internal/config"
	"arc-scanner/internal/items"
	"arc-scanner/internal/keyboard"
	"arc-scanner/internal/scanner"

	"github.com/go-vgo/robotgo"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App is the main application struct.
type App struct {
	ctx     context.Context
	scanner scanner.Scanner
	matcher *items.Matcher
	repo    *items.Repository
}

// NewApp creates a new App instance.
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. It initializes all components.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	if err := a.initWindow(ctx); err != nil {
		slog.Error("failed to initialize window", "error", err)
	}

	itemsList, err := a.initItems()
	if err != nil {
		slog.Error("failed to initialize items", "error", err)
		runtime.EventsEmit(ctx, "startup-error", err.Error())
		return
	}

	a.scanner = scanner.New()
	a.matcher = items.NewMatcher(itemsList)

	a.initKeyboardHook(ctx, itemsList)

	slog.Info("application started", "items", len(itemsList))
}

// initWindow sets up the application window.
func (a *App) initWindow(ctx context.Context) error {
	runtime.WindowSetAlwaysOnTop(ctx, true)

	screens, err := runtime.ScreenGetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get screens: %w", err)
	}

	var currentScreen runtime.Screen
	for _, screen := range screens {
		if screen.IsCurrent {
			currentScreen = screen
			break
		}
	}

	// Fallback: use first screen if IsCurrent not found (Windows compatibility)
	if currentScreen.Size.Width == 0 && len(screens) > 0 {
		currentScreen = screens[0]
	}

	windowWidth := 1
	windowHeight := 1
	x := currentScreen.Size.Width - windowWidth

	runtime.WindowSetSize(ctx, windowWidth, windowHeight)
	runtime.WindowSetPosition(ctx, x, 0)

	// Set window level above fullscreen apps (needs delay for window creation)
	go func() {
		time.Sleep(100 * time.Millisecond)
		setWindowAboveFullscreen()
	}()

	return nil
}

// initItems loads or fetches the item database.
func (a *App) initItems() ([]items.Item, error) {
	appDataDir, err := getAppDataDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get app data directory: %w", err)
	}

	if err := os.MkdirAll(appDataDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create app data directory: %w", err)
	}

	cachePath := filepath.Join(appDataDir, "items.json")
	a.repo = items.NewRepository(cachePath)

	// Try to load from cache first
	if a.repo.CacheExists() {
		slog.Info("loading items from cache")
		itemsList, err := a.repo.LoadFromCache()
		if err == nil {
			return itemsList, nil
		}
		slog.Warn("cache load failed, fetching from API", "error", err)
	}

	// Fetch from API
	slog.Info("fetching items from API")
	itemsList, err := a.repo.FetchFromAPI()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch items: %w", err)
	}

	// Save to cache
	if err := a.repo.SaveToCache(itemsList); err != nil {
		slog.Warn("failed to save cache", "error", err)
	}

	return itemsList, nil
}

// initKeyboardHook sets up the global keyboard listener.
func (a *App) initKeyboardHook(ctx context.Context, itemsList []items.Item) {
	itemsMap := items.BuildIndex(itemsList)
	hook := keyboard.New(ctx)

	// Register scan key handler
	hook.Register(config.ScanKey, func() {
		a.handleScan(itemsMap)
	})

	// Register toggle visibility handler
	hook.Register(config.ToggleKey, func() {
		slog.Debug("toggling overlay visibility")
		runtime.EventsEmit(a.ctx, "toggle-visibility", nil)
	})

	hook.Start()
}

// handleScan performs a scan at the current mouse position.
func (a *App) handleScan(itemsMap items.ItemMap) {
	startTime := time.Now()
	x, y := robotgo.Location()

	slog.Debug("scanning", "x", x, "y", y)

	// Take screenshot
	img, err := a.scanner.TakeScreenshot(x, y)
	if err != nil {
		slog.Error("screenshot failed", "error", err)
		return
	}

	// Perform OCR
	text, err := a.scanner.ProcessImage(img)
	if err != nil {
		slog.Error("OCR failed", "error", err)
		return
	}

	// Clean and match item
	tokens := items.CleanOCRText(text)
	item, err := a.matcher.FindItem(tokens)
	if err != nil {
		slog.Debug("item not found", "tokens", tokens)
		runtime.EventsEmit(a.ctx, "scan-failed", nil)
		return
	}

	quantity := items.ParseQuantity(text)
	slog.Info("item found",
		"name", item.Name,
		"value", item.Value,
		"quantity", quantity,
		"duration", time.Since(startTime))

	// Log recycle info if available
	if item.RecycleComponents != nil {
		logRecycleInfo(item, itemsMap)
	}

	runtime.EventsEmit(a.ctx, "item-found", item)
}

// logRecycleInfo logs the recycling value breakdown for an item.
func logRecycleInfo(item items.Item, itemsMap items.ItemMap) {
	totalValue := 0
	for _, entry := range *item.RecycleComponents {
		component, ok := itemsMap.Get(entry.Component.ID)
		if ok {
			totalValue += entry.Quantity * component.Value
		}
	}
	slog.Debug("recycle value", "item", item.Name, "total", totalValue)
}

// getAppDataDir returns the platform-specific application data directory.
func getAppDataDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	var appDataDir string
	if os.Getenv("OS") == "Windows_NT" || filepath.Separator == '\\' {
		// Windows: Use APPDATA environment variable
		appData := os.Getenv("APPDATA")
		if appData == "" {
			appData = filepath.Join(homeDir, "AppData", "Roaming")
		}
		appDataDir = filepath.Join(appData, "arc-scanner")
	} else {
		// macOS and other Unix-like systems
		appDataDir = filepath.Join(homeDir, "Library", "Application Support", "arc-scanner")
	}

	return appDataDir, nil
}
