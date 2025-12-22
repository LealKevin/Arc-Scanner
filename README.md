# Arc Scanner

A desktop overlay app that scans in-game items from Arc Raiders and displays their value using OCR.

## Download

Download the latest release from the [Releases page](https://github.com/Mizuba74/arc-scanner/releases).

### macOS Installation

1. Download `arc-scanner-macos-vX.X.X.zip`
2. Unzip and drag `arc-scanner.app` to your Applications folder
3. **First launch:** Right-click the app → Click "Open" → Click "Open" in the dialog
   - This is required because the app is not notarized with Apple
4. Grant permissions when prompted:
   - **Accessibility** - to detect keyboard input (press 'y' to scan)
   - **Screen Recording** - to capture screenshots for OCR

### Windows Installation

1. Download `arc-scanner-windows-vX.X.X.zip`
2. Extract to a folder (keep the folder structure intact)
3. Run `arc-scanner.exe`
4. If Windows SmartScreen appears: Click "More info" → "Run anyway"
   - This warning appears because the app is not signed with a certificate

---

## Development

### Requirements

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

### Create a Release (macOS + Windows)

```bash
./scripts/create-release.sh 0.1.0
```

Creates distribution packages in `releases/v0.1.0/`:
- `arc-scanner-macos-v0.1.0.zip` - macOS app bundle
- `arc-scanner-windows-v0.1.0.zip` - Windows executable with dependencies

### Build macOS Only

```bash
./scripts/build-bundled.sh
```

Creates a self-contained `.app` at `build/bin/arc-scanner.app` with bundled Tesseract.

### Build Windows Only (Cross-Compile from macOS)

```bash
./scripts/build-bundled-windows-cross.sh
```

Creates Windows executable at `build/bin/arc-scanner.exe` with bundled Tesseract.

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
