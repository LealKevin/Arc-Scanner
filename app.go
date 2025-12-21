package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"arc-scanner/ocr"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx     context.Context
	scanner *ocr.Scanner
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	runtime.WindowSetAlwaysOnTop(ctx, true)
	screens, err := runtime.ScreenGetAll(ctx)
	if err != nil {
		println("Error getting screens:", err.Error())
		return
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
	y := 0

	runtime.WindowSetSize(ctx, windowWidth, windowHeight)
	runtime.WindowSetPosition(ctx, x, y)

	// Set window level above fullscreen apps
	// Need a small delay to ensure window is fully created
	go func() {
		time.Sleep(100 * time.Millisecond)
		setWindowAboveFullscreen()
	}()

	// Get app data directory
	appDataDir, err := getAppDataDir()
	if err != nil {
		panic(err)
	}
	err = os.MkdirAll(appDataDir, 0755)
	if err != nil {
		panic(err)
	}
	itemsPath := filepath.Join(appDataDir, "items.json")

	_, err = os.Stat(itemsPath)
	if os.IsNotExist(err) {
		fmt.Println("items.json not found, creating...")
		items, err := fetchItems()
		if err != nil {
			panic(err)
		}

		err = saveItemsJSON(items, itemsPath)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("items.json found, loading...")
	}

	itemsJson, err := getItems(itemsPath)
	if err != nil {
		panic(err)
	}

	itemsMap := buildItemIndex(itemsJson)
	fmt.Println("Items retrieved, count:", len(itemsMap))

	a.scanner = ocr.NewScanner()

	go func() {
		evChan := hook.Start()
		defer hook.End()

		fmt.Println("Keyboard hook started. Press 'k' to scan.")

		for ev := range evChan {
			if ev.Keychar == 'y' {
				timeStart := time.Now()
				X, Y := robotgo.Location()
				fmt.Println("X:", X, "Y:", Y)

				// Time screenshot
				screenshotStart := time.Now()
				img, err := a.scanner.TakeScreenshot(X, Y)
				if err != nil {
					fmt.Println("Screenshot error:", err)
					continue
				}
				fmt.Printf("Screenshot took: %v\n", time.Since(screenshotStart))

				// Time OCR
				ocrStart := time.Now()
				text, err := a.scanner.ProcessImage(img)
				if err != nil {
					fmt.Println("OCR error:", err)
					continue
				}
				fmt.Printf("OCR took: %v\n", time.Since(ocrStart))

				cleanText := cleanText(text)
				item, err := findItemInText(cleanText, itemsJson)
				if err != nil {
					fmt.Println("Item not found in text:", err)
					// Emit empty item so frontend shows "Object not found"
					runtime.EventsEmit(a.ctx, "scan-failed", nil)
					timeElapsed := time.Since(timeStart)
					fmt.Println("Time elapsed:", timeElapsed)
					continue
				}

				quantity := getQuantity(text)
				fmt.Println(text)
				fmt.Println("Quantity:", quantity)
				fmt.Println("Item found:", item)
				fmt.Println(itemsMap.getItemInfo(item))

				runtime.EventsEmit(a.ctx, "item-found", item)
				timeElapsed := time.Since(timeStart)
				fmt.Println("Time elapsed:", timeElapsed)
			}
			if ev.Keychar == 'u' {
				fmt.Println("y pressed - toggling overlay visibility")
				runtime.EventsEmit(a.ctx, "toggle-visibility", nil)
			}
		}
	}()
}

type Response struct {
	Data []Item `json:"data"`
}

type Item struct {
	ID                string          `json:"id"`
	Name              string          `json:"name"`
	Value             int             `json:"value"`
	Icon              string          `json:"icon"`
	RecycleComponents *[]RecycleEntry `json:"recycle_components"`
	UsedIn            *[]UsedInEntry  `json:"used_in"`
}

type RecycleEntry struct {
	Quantity  int       `json:"quantity"`
	Component Component `json:"component"`
}

type UsedInEntry struct {
	Quantity int  `json:"quantity"`
	Item     Item `json:"item"`
}

type Component struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func fetchItems() (Items []Item, err error) {
	var items []Item
	fmt.Println("Getting items...")

	apiString := "https://metaforge.app/api/arc-raiders/items?minimal=true&includeComponents=true&limit=100&page="
	for i := 1; i <= 6; i++ {
		var response Response
		apiStrinfFinal := apiString + fmt.Sprintf("%d", i)
		fmt.Println(apiStrinfFinal)

		data, err := http.Get(apiStrinfFinal)
		if err != nil {
			return nil, fmt.Errorf("failed to get items: %w", err)
		}
		defer data.Body.Close()

		err = json.NewDecoder(data.Body).Decode(&response)
		if err != nil {
			return nil, fmt.Errorf("failed to decode items: %w", err)
		}

		items = append(items, response.Data...)
	}

	fmt.Println("Items retrieved, count:", len(items))

	fmt.Println("Done!")
	return items, nil
}

func getItems(path string) (Items []Item, err error) {
	var items []Item

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read items.json: %w", err)
	}

	err = json.Unmarshal(data, &items)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal items: %w", err)
	}

	return items, nil
}

func saveItemsJSON(items []Item, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create items.json: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(items)
	if err != nil {
		return fmt.Errorf("failed to encode items: %w", err)
	}
	return nil
}

type ItemMap map[string]Item

func buildItemIndex(items []Item) ItemMap {
	itemMap := make(ItemMap)
	for _, item := range items {
		itemMap[item.ID] = item
	}
	return itemMap
}

func (itemMap ItemMap) getItemByID(id string) (Item, bool) {
	item, ok := itemMap[id]
	return item, ok
}

func cleanText(text string) []string {
	var texts []string
	lines := strings.Split(text, "\n")
	for _, line := range lines {

		line = strings.ReplaceAll(line, "|", "I")
		line = strings.ReplaceAll(line, ".", "")
		words := strings.Split(line, " ")
		for _, word := range words {
			if word == strings.ToUpper(word) {
				texts = append(texts, word)
			}
		}
	}
	return texts
}

func getQuantity(text string) int {
	fmt.Println(text)
	quantity := 1
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if strings.Contains(line, "/") {
			parts := strings.Split(line, "/")
			if len(parts) >= 2 {
				numStr := strings.TrimSpace(parts[0])
				if num, err := strconv.Atoi(numStr); err == nil {
					return num
				}
			}
		}
	}
	return quantity
}

func findItemInText(text []string, items []Item) (Item, error) {
	textJoined := strings.Join(text, " ")
	romanSuffix := []string{"IV", "III", "II", "I"}
	var itemFound Item

	for _, item := range items {
		cleanItemName := strings.ToUpper(strings.ReplaceAll(item.ID, "-", " "))
		cleanItemName = strings.ReplaceAll(cleanItemName, "RECIPE", "")
		if strings.Contains(textJoined, cleanItemName) {
			itemFound = item
			for _, suffix := range romanSuffix {
				if strings.Contains(textJoined, cleanItemName+suffix) {
					itemFound = item
				}
			}
		}
	}
	fmt.Println("Here")
	fmt.Println(itemFound)
	if itemFound.ID != "" {
		return itemFound, nil
	}
	return Item{}, fmt.Errorf("item not found in text: %s", textJoined)
}

func (itemMap ItemMap) getItemInfo(item Item) string {
	var info string
	info += fmt.Sprintf("ID: %s\n", item.ID)
	info += fmt.Sprintf("Name: %s\n", item.Name)
	info += fmt.Sprintf("Value: %d\n", item.Value)

	if item.RecycleComponents != nil {
		info += "Recycle Components:\n"
		valueOfRecycleComponents := 0
		for _, entry := range *item.RecycleComponents {
			component, ok := itemMap.getItemByID(entry.Component.ID)
			if !ok {
				panic(fmt.Errorf("component not found: %s", entry.Component.ID))
			}
			valueOfRecycleComponents += entry.Quantity * component.Value
		}
		for _, entry := range *item.RecycleComponents {
			info += fmt.Sprintf("\tQuantity: %d\n", entry.Quantity)
			info += fmt.Sprintf("\tComponent: %s\n", entry.Component.Name)

		}
		info += fmt.Sprintf("\tTotal value: %d\n", valueOfRecycleComponents)
	}
	return info
}

func getAppDataDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Platform-specific app data directory
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
