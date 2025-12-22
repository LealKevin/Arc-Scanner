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
	"arc-scanner/internal/updater"

	"github.com/go-vgo/robotgo"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Version is set at build time via ldflags: -ldflags "-X main.Version=1.0.0"
var Version = "dev"

type App struct {
	ctx     context.Context
	scanner scanner.Scanner
	matcher *items.Matcher
	repo    *items.Repository
	updater *updater.Updater
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Request permissions upfront (macOS only) so user only needs one restart
	requestPermissions()

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

	// Initialize updater
	a.updater = updater.New("LealKevin", "Arc-Scanner", Version)

	// Check for updates in background
	go a.checkForUpdates()

	a.initKeyboardHook(ctx, itemsList)

	slog.Info("application started", "items", len(itemsList), "version", Version)
}

func (a *App) checkForUpdates() {
	// Wait for frontend to be ready before checking
	time.Sleep(500 * time.Millisecond)

	slog.Info("checking for updates", "currentVersion", Version)

	// Skip update check in dev mode
	if Version == "dev" {
		slog.Info("skipping update check in dev mode")
		return
	}

	info, err := a.updater.CheckForUpdate()
	if err != nil {
		slog.Error("failed to check for updates", "error", err)
		return
	}

	if info == nil {
		slog.Info("no updates available")
		return
	}

	slog.Info("update available", "version", info.Version, "url", info.URL)
	runtime.EventsEmit(a.ctx, "update-available", info)
}

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
	// Windows needs more time for WebView2 transparency to initialize properly
	go func() {
		time.Sleep(500 * time.Millisecond)
		setWindowAboveFullscreen()
	}()

	return nil
}

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

	if a.repo.CacheExists() {
		slog.Info("loading items from cache")
		itemsList, err := a.repo.LoadFromCache()
		if err == nil {
			return itemsList, nil
		}
		slog.Warn("cache load failed, fetching from API", "error", err)
	}

	slog.Info("fetching items from API")
	itemsList, err := a.repo.FetchFromAPI()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch items: %w", err)
	}

	if err := a.repo.SaveToCache(itemsList); err != nil {
		slog.Warn("failed to save cache", "error", err)
	}

	return itemsList, nil
}

func (a *App) initKeyboardHook(ctx context.Context, itemsList []items.Item) {
	itemsMap := items.BuildIndex(itemsList)
	hook := keyboard.New(ctx)

	hook.Register(config.ScanKey, func() {
		a.handleScan(itemsMap)
	})

	hook.Register(config.ToggleKey, func() {
		slog.Debug("toggling overlay visibility")
		runtime.EventsEmit(a.ctx, "toggle-visibility", nil)
	})

	hook.Start()
}

func (a *App) handleScan(itemsMap items.ItemMap) {
	startTime := time.Now()
	x, y := robotgo.Location()

	slog.Debug("scanning", "x", x, "y", y)

	img, err := a.scanner.TakeScreenshot(x, y)
	if err != nil {
		slog.Error("screenshot failed", "error", err)
		return
	}

	text, err := a.scanner.ProcessImage(img)
	if err != nil {
		slog.Error("OCR failed", "error", err)
		return
	}

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

	if item.RecycleComponents != nil {
		logRecycleInfo(item, itemsMap)
	}

	runtime.EventsEmit(a.ctx, "item-found", item)
}

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

// GetVersion returns the current app version
func (a *App) GetVersion() string {
	return Version
}

// DownloadUpdate downloads the available update
// Emits "update-progress" events with percentage (0-100)
// Emits "update-ready" when download is complete
func (a *App) DownloadUpdate(info *updater.UpdateInfo) error {
	slog.Info("starting update download", "version", info.Version)

	err := a.updater.DownloadUpdate(a.ctx, info, func(percent int) {
		runtime.EventsEmit(a.ctx, "update-progress", percent)
	})

	if err != nil {
		slog.Error("failed to download update", "error", err)
		runtime.EventsEmit(a.ctx, "update-error", err.Error())
		return err
	}

	slog.Info("update downloaded successfully")
	runtime.EventsEmit(a.ctx, "update-ready", nil)
	return nil
}

// ApplyUpdateAndRestart applies the downloaded update and restarts the app
func (a *App) ApplyUpdateAndRestart() error {
	slog.Info("applying update and restarting")

	if err := a.updater.ApplyUpdate(); err != nil {
		slog.Error("failed to apply update", "error", err)
		return err
	}

	// Quit the app - the update script will relaunch it
	runtime.Quit(a.ctx)
	return nil
}
