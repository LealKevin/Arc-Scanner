# Development Guide

This guide covers building and developing Arc Scanner.

## Requirements

- Go 1.24+
- Node.js 18+
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)
- Tesseract OCR (for development only)

### macOS
```bash
brew install tesseract
```

### Windows
Download from [UB-Mannheim/tesseract](https://github.com/UB-Mannheim/tesseract/wiki) and install to `C:\Program Files\Tesseract-OCR`.

## Development Mode

```bash
wails dev
```

Runs with hot reload. Frontend dev server: http://localhost:34115

## Building

### Create a Release (macOS + Windows)

```bash
./scripts/create-release.sh 0.1.0
```

Creates distribution packages in `releases/v0.1.0/`:
- `arc-scanner-macos-v0.1.0.zip`
- `arc-scanner-windows-v0.1.0.zip`

### macOS Only

**For distribution (self-contained):**
```bash
./scripts/build-bundled.sh
```
Creates `.app` at `build/bin/arc-scanner.app` with bundled Tesseract. No user setup required.

**For quick local testing:**
```bash
wails build
```
Requires system Tesseract. Not for distribution.

### Windows Only

**Cross-compile from macOS/Linux:**
```bash
./scripts/build-bundled-windows-cross.sh
```
Creates Windows executable with bundled Tesseract. Requires `p7zip` or `innoextract`:
```bash
brew install p7zip
```

**Build on Windows:**
```batch
scripts\build-bundled-windows.bat
```

**Quick build (development):**
```bash
wails build -platform windows/amd64
```

See [WINDOWS_BUILD.md](WINDOWS_BUILD.md) for detailed Windows instructions.

## Project Structure

```
arc-scanner/
├── app.go                 # Main app logic, OCR orchestration
├── main.go                # Entry point
├── window_darwin.go       # macOS window management
├── window_windows.go      # Windows window management
├── internal/
│   ├── config/            # Configuration constants
│   ├── items/             # Item database and matching
│   ├── keyboard/          # Global keyboard hooks
│   ├── scanner/           # Screenshot + Tesseract OCR
│   └── updater/           # Auto-update logic
├── frontend/
│   ├── src/
│   │   ├── App.tsx        # Main React component
│   │   └── components/    # UI components
│   └── wailsjs/           # Auto-generated Wails bindings
└── scripts/               # Build scripts
```

## Frontend Commands

```bash
cd frontend
npm install    # Install dependencies
npm run dev    # Run Vite dev server (standalone)
npm run build  # Build frontend
```

## Testing

```bash
go test ./...
```

## Architecture

### Event System

Go emits events to the frontend:
- `item-found` - Item successfully identified
- `scan-failed` - No item found
- `toggle-visibility` - Toggle overlay
- `update-available` - New version available

### OCR Pipeline

1. `robotgo` detects mouse position
2. `screenshot` captures 450x380px area around cursor
3. Image preprocessed (grayscale, invert, contrast, sharpen)
4. Tesseract extracts text (PSM 3, OEM 1)
5. Text matched against item database
6. Result emitted to frontend

### Platform-Specific Code

Build tags control compilation:
- `//go:build darwin` - macOS only
- `//go:build windows` - Windows only

## Release Workflow

See [RELEASING.md](RELEASING.md) for the release process.
