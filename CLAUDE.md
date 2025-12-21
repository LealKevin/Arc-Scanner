# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Arc Scanner is a desktop overlay application built with Wails (Go + React/TypeScript) that scans in-game items from Arc Raiders and displays their value. It uses OCR (Tesseract) to identify items from screenshots and shows an overlay with item information.

## Technology Stack

- **Backend**: Go 1.24 with Wails v2.11.0
- **Frontend**: React 18 + TypeScript + Vite
- **OCR**: Tesseract (via command-line execution)
- **Platform Support**: macOS and Windows (with platform-specific window management)
- **Additional Libraries**:
  - `robotgo` - for mouse position tracking (cross-platform)
  - `gohook` - for keyboard event hooks (cross-platform)
  - `screenshot` - for screen capture (cross-platform)
  - `imaging` - for image processing (cross-platform)
  - `lxn/win` - for Windows API calls (Windows only)

## Commands

### Development
```bash
wails dev
```
Runs live development mode with hot reload. Frontend dev server available at http://localhost:34115.

### Building

#### macOS

**For distribution (self-contained):**
```bash
./scripts/build-bundled.sh
```
Creates a fully self-contained .app with bundled Tesseract. No user setup required - just download and run. Output: `build/bin/arc-scanner.app` (~20MB)

**For quick local testing:**
```bash
wails build
```
Standard Wails build without bundling. Requires system Tesseract. Faster to build but NOT suitable for distribution. Output: `build/bin/arc-scanner.app` (~10MB)

**What's the difference?**
- `./scripts/build-bundled.sh` = Bundles Tesseract → Runs `wails build` → Copies resources into .app
- `wails build` = Just builds the .app (no Tesseract bundling)

#### Windows

**For distribution (self-contained) - FROM macOS/Linux:**
```bash
./scripts/build-bundled-windows-cross.sh
```
Cross-compiles Windows .exe on macOS/Linux with bundled Tesseract. Downloads Windows Tesseract binaries automatically. No Windows machine needed! Output: `build/bin/arc-scanner.exe` (~15-20MB)

**Requirements:** One of these extraction tools:
- `brew install p7zip` (provides `7z` command)
- `brew install innoextract`

**For distribution (self-contained) - FROM Windows:**
```batch
scripts\build-bundled-windows.bat
```
Builds on Windows with bundled Tesseract. Requires local Tesseract installation. Output: `build/bin/arc-scanner.exe` (~15-20MB)

**For quick local testing:**
```bash
wails build -platform windows/amd64
```
Standard Wails build without bundling. Requires system Tesseract. Faster to build but NOT suitable for distribution. Output: `build/bin/arc-scanner.exe` (~10MB)

**See WINDOWS_BUILD.md for detailed Windows build instructions.**

### Frontend Commands
```bash
cd frontend
npm install          # Install dependencies
npm run dev          # Run Vite dev server (standalone)
npm run build        # Build frontend (TypeScript + Vite)
```

## Architecture

### Application Flow

1. **Startup** (`main.go`): Creates frameless, always-on-top window positioned at top-right corner
2. **Window Management**:
   - macOS (`window_darwin.go`): Uses Cocoa framework to set window above fullscreen apps
   - Windows (`window_windows.go`): Uses Windows API to set HWND_TOPMOST
3. **Keyboard Hook** (`app.go:105-163`): Global keyboard listener for 'y' (scan) and 'u' (toggle visibility)
4. **Scan Process**:
   - Mouse position detected via `robotgo`
   - Screenshot captured relative to cursor position (450x380px box)
   - Image preprocessed (grayscale, invert, contrast, sharpen)
   - Tesseract OCR extracts text
   - Item matched against local database (`items.json`)
   - Result emitted to frontend via Wails events

### Data Model

**Items** (`items.json`): Fetched from `metaforge.app` API on first run, cached locally. Contains:
- Item ID, name, value, icon URL
- Recycle components (what you get from recycling)
- Used in (crafting recipes)

### Key Files

- `app.go`: Main application logic, item database management, OCR orchestration, platform-agnostic app data paths
- `ocr/scanner.go`: Screenshot capture and Tesseract integration with platform-specific path detection
- `window_darwin.go`: macOS-specific window level management using CGI/Cocoa
- `window_windows.go`: Windows-specific window management using Windows API
- `window_other.go`: Stub for unsupported platforms
- `frontend/src/App.tsx`: React overlay UI with fade animations for item display
- `items.json`: Cached item database (auto-generated on first run)

### Go-Frontend Communication

Wails binds Go methods to frontend and uses event system:
- **Events emitted from Go**:
  - `item-found` - Item successfully identified
  - `scan-failed` - No item found in OCR text
  - `toggle-visibility` - Toggle overlay visibility
- **Frontend** (`App.tsx`): Listens via `EventsOn()` and manages UI state

### OCR Configuration

Tesseract runs with:
- PSM 3 (fully automatic page segmentation)
- OEM 1 (LSTM only, faster)
- Whitelist: `0123456789/ ABCDEFGHIJKLMNOPQRSTUVWXYZ`

Screenshot preprocessing is critical for accuracy (see `ocr/scanner.go:44-47`).

### Platform-Specific Code

- `window_darwin.go`: macOS window management with Cocoa (build tag: `//go:build darwin`)
- `window_windows.go`: Windows window management with Windows API (build tag: `//go:build windows`)
- `window_other.go`: Stub for unsupported platforms (build tag: `//go:build !darwin`)
- `app.go`: Platform-agnostic app data paths (detects Windows vs macOS at runtime)
- `ocr/scanner.go`: Platform-specific Tesseract path detection (detects Windows vs macOS paths)

Build tags control compilation - only the appropriate platform files are compiled for each target.

## Dependencies

### External Requirements (Development Only)

**macOS:**
Tesseract OCR must be installed for development:
- `/opt/homebrew/bin/tesseract` (Homebrew ARM)
- `/usr/local/bin/tesseract` (Homebrew Intel)
- Or available in PATH

**Windows:**
Tesseract OCR must be installed for development and bundling:
- Download from https://github.com/UB-Mannheim/tesseract/wiki
- Install to `C:\Program Files\Tesseract-OCR` (default)
- Or set `TESSERACT_PATH` environment variable

### Distribution (Self-Contained)

**macOS:**
When built with `./scripts/build-bundled.sh`, the app bundles:
- Tesseract binary (102KB)
- Required dylibs (libleptonica, libtesseract, libarchive ~5.5MB)
- Training data (eng.traineddata ~3.9MB)
- **Total overhead: ~9.5MB**

Scanner auto-detects Tesseract location in this order:
1. Bundled version (in .app/Contents/Resources/bin/)
2. Homebrew installation (for development)
3. System PATH

**Windows:**
When built with `scripts\build-bundled-windows.bat`, the app bundles:
- Tesseract binary (400KB)
- Required DLLs (libleptonica, libtesseract, runtime DLLs ~8MB)
- Training data (eng.traineddata ~3.9MB)
- **Total overhead: ~12MB**

Scanner auto-detects Tesseract location in this order:
1. Bundled version (in build/windows/bin/tesseract.exe)
2. System installation (`C:\Program Files\Tesseract-OCR\tesseract.exe`)
3. System PATH

**Both platforms:**
The bundled app requires NO external dependencies - users just download and run.

### Wails Bindings

Frontend automatically gets generated bindings in `frontend/wailsjs/` - these are auto-generated and should not be manually edited.

## Platform-Specific Notes

### macOS Permissions (For End Users)

The app requires two permissions on first launch:

1. **Accessibility** - To detect keyboard input (pressing 'y' to scan)
2. **Screen Recording** - To capture screenshots for OCR

**If scanning doesn't work:**
- Open System Settings → Privacy & Security
- Go to Accessibility → Add arc-scanner.app → Enable checkbox
- Go to Screen Recording → Add arc-scanner.app → Enable checkbox
- Restart the app

### Windows Requirements (For End Users)

**WebView2 Runtime:**
- Usually pre-installed on Windows 10/11
- Download from https://developer.microsoft.com/microsoft-edge/webview2/ if needed

**No special permissions required** (unlike macOS):
- Keyboard hooks work without admin rights
- Screen capture works without special permissions
- Windows Defender may flag unsigned .exe (expected behavior)

**App Data Location:**
- macOS: `~/Library/Application Support/arc-scanner/items.json`
- Windows: `%APPDATA%\arc-scanner\items.json` (C:\Users\<username>\AppData\Roaming\arc-scanner)

## Git Commit Guidelines

- **Always ask for approval** before creating commits
- **Simple commit messages** - clear and concise, no verbose explanations
- **No AI signatures** - do not add "Generated with Claude" or similar footers to commit messages
- Commit messages should describe what changed, not how or why
