# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Arc Scanner is a desktop overlay application built with Wails (Go + React/TypeScript) that scans in-game items from Arc Raiders and displays their value. It uses OCR (Tesseract) to identify items from screenshots and shows an overlay with item information.

## Technology Stack

- **Backend**: Go 1.24 with Wails v2.11.0
- **Frontend**: React 18 + TypeScript + Vite
- **OCR**: Tesseract (via command-line execution)
- **Platform Support**: macOS (with Darwin-specific window management)
- **Additional Libraries**:
  - `robotgo` - for mouse position tracking
  - `gohook` - for keyboard event hooks
  - `screenshot` - for screen capture
  - `imaging` - for image processing

## Commands

### Development
```bash
wails dev
```
Runs live development mode with hot reload. Frontend dev server available at http://localhost:34115.

### Building
```bash
wails build
```
Creates a production build. Output binary will be named `arc-scanner`.

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
2. **Window Management** (`window_darwin.go`): Uses Cocoa framework to set window above fullscreen apps
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

- `app.go`: Main application logic, item database management, OCR orchestration
- `ocr/scanner.go`: Screenshot capture and Tesseract integration
- `window_darwin.go`: macOS-specific window level management using CGI/Cocoa
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

- `window_darwin.go`: macOS window management with Cocoa
- `window_other.go`: No-op for non-Darwin platforms
- Build tags (`//go:build darwin`) control compilation

## Dependencies

### External Requirements

- **Tesseract OCR** must be installed:
  - `/opt/homebrew/bin/tesseract` (Homebrew ARM)
  - `/usr/local/bin/tesseract` (Homebrew Intel)
  - Or available in PATH

Scanner will auto-detect Tesseract location on initialization.

### Wails Bindings

Frontend automatically gets generated bindings in `frontend/wailsjs/` - these are auto-generated and should not be manually edited.

## Git Commit Guidelines

- **Always ask for approval** before creating commits
- **Simple commit messages** - clear and concise, no verbose explanations
- **No AI signatures** - do not add "Generated with Claude" or similar footers to commit messages
- Commit messages should describe what changed, not how or why
