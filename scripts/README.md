# Build Scripts

## Overview

Scripts for building a self-contained Arc Scanner app with bundled Tesseract OCR.

## Scripts

### `build-bundled.sh` (Recommended)

**Main build script** - Creates a self-contained .app with all dependencies bundled.

```bash
./scripts/build-bundled.sh
```

**What it does:**
1. Runs `bundle-tesseract.sh` to prepare Tesseract binaries
2. Runs `wails build` to compile the app
3. Runs `post-build.sh` to copy bundled resources into the .app

**Output:** `build/bin/arc-scanner.app` (fully self-contained)

### `bundle-tesseract.sh`

Prepares Tesseract binary, dependencies, and training data for bundling.

```bash
./scripts/bundle-tesseract.sh
```

**What it does:**
- Copies Tesseract binary from Homebrew installation
- Copies required dylibs (libleptonica, libtesseract, libarchive)
- Relinks dylibs to use `@executable_path` (makes them portable)
- Downloads eng.traineddata from tesseract-ocr/tessdata_fast

**Output:** `build/darwin/` directory with:
- `bin/tesseract` (102KB)
- `lib/*.dylib` (~5.5MB)
- `tessdata/eng.traineddata` (~3.9MB)

**Prerequisites:** Tesseract must be installed via Homebrew

### `post-build.sh`

Copies bundled resources from `build/darwin/` into the .app bundle.

```bash
./scripts/post-build.sh
```

**What it does:**
- Copies `build/darwin/bin/` → `arc-scanner.app/Contents/Resources/bin/`
- Copies `build/darwin/lib/` → `arc-scanner.app/Contents/Resources/lib/`
- Copies `build/darwin/tessdata/` → `arc-scanner.app/Contents/Resources/tessdata/`

**Prerequisites:** Must run after `wails build`

## Development Workflow

**For distribution (self-contained app):**
```bash
./scripts/build-bundled.sh
```

**For development (faster, uses system Tesseract):**
```bash
wails dev  # Or wails build
```

## How Bundling Works

1. **Detection Order** (in `ocr/scanner.go`):
   - First checks for bundled Tesseract in `.app/Contents/Resources/bin/`
   - Falls back to Homebrew installation (`/opt/homebrew/bin/tesseract`)
   - Final fallback to system PATH

2. **Dynamic Library Linking**:
   - `install_name_tool` relinks dylibs to use `@executable_path`
   - Allows libraries to find each other relative to the binary
   - Makes the app portable without absolute paths

3. **Training Data**:
   - Sets `TESSDATA_PREFIX` environment variable when using bundled version
   - Points Tesseract to bundled training data

## Troubleshooting

**"Tesseract not found" error during bundling:**
- Install Tesseract: `brew install tesseract`

**App can't find bundled Tesseract:**
- Verify files exist: `ls -la build/bin/arc-scanner.app/Contents/Resources/`
- Check post-build.sh ran successfully

**Code signature warnings:**
- Normal when modifying dylibs with `install_name_tool`
- App will work fine unsigned
- Users can re-sign if needed with `codesign`
