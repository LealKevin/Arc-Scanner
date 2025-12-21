# Arc Scanner

A desktop overlay app that scans in-game items from Arc Raiders and displays their value using OCR.

## Requirements

- macOS (Apple Silicon or Intel)
- Go 1.24+
- Node.js 18+
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)
- Tesseract OCR (`brew install tesseract`)

## Development

```bash
wails dev
```

Runs live development mode with hot reload.

## Building

### For Distribution (Recommended)

```bash
./scripts/create-installer.sh
```

Creates a `.pkg` installer at `build/arc-scanner-installer.pkg` that:
- Bundles Tesseract and all dependencies (no user setup required)
- Installs to `/Applications`
- Removes quarantine automatically after install

### For Local Testing

```bash
./scripts/build-bundled.sh
```

Creates a self-contained `.app` at `build/bin/arc-scanner.app` with bundled Tesseract.

### Quick Build (Development Only)

```bash
wails build
```

Standard Wails build. Requires system Tesseract. Not suitable for distribution.

## Usage

1. Launch the app - a small green dot appears at the top-right corner
2. In game, hover your mouse over an item
3. Press **'y'** to scan
4. The item value appears in the overlay
5. Press **'u'** to hide/show the overlay

## macOS Permissions

The app requires:
- **Accessibility** - to detect keyboard input
- **Screen Recording** - to capture screenshots for OCR
