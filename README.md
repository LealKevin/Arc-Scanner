<p align="center">
  <img src="assets/logo.png" alt="Arc Scanner" width="100"/>
</p>

<h1 align="center">Arc Scanner</h1>

<p align="center">
  <em>Instantly scan Arc Raiders items and see their value</em>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/platform-macOS%20%7C%20Windows-blue" alt="Platform"/>
  <img src="https://img.shields.io/github/v/release/LealKevin/Arc-Scanner" alt="Version"/>
  <img src="https://img.shields.io/github/license/LealKevin/Arc-Scanner" alt="License"/>
</p>

<p align="center">
  <img src="assets/demo.gif" alt="Demo" width="600"/>
</p>

<!-- TODO: Add demo.gif showing: hover item → press Y → overlay appears → fades -->

## What is Arc Scanner?

Arc Scanner is a desktop overlay for Arc Raiders that identifies items using OCR and instantly displays their value. Just hover over any item in-game and press a key to scan.

## Download

Download the latest release for your platform:

**[Download for macOS](https://github.com/LealKevin/Arc-Scanner/releases/latest)** | **[Download for Windows](https://github.com/LealKevin/Arc-Scanner/releases/latest)**

## Features

- **Instant Scanning** — Press `Y` to scan any item under your cursor
- **Always-On-Top Overlay** — Works over fullscreen games without interrupting gameplay
- **Cross-Platform** — Native support for macOS and Windows
- **Self-Contained** — No external dependencies, just download and run
- **Auto-Updates** — Automatically checks for and installs new versions
- **Recycle Info** — Shows component breakdown for recyclable items
- **Minimal & Non-Intrusive** — Tiny overlay that appears only when needed

## Usage

| Key | Action |
|-----|--------|
| `Y` | Scan item under cursor |
| `U` | Toggle overlay visibility |

1. Launch Arc Scanner — a small indicator appears in the top-right corner
2. In game, hover your mouse over an item
3. Press `Y` to scan
4. The item value appears briefly in the overlay

## Installation

### macOS

1. Download `arc-scanner-macos-vX.X.X.zip` from [Releases](https://github.com/LealKevin/Arc-Scanner/releases)
2. Unzip and drag `arc-scanner.app` to Applications
3. **First launch:** Right-click → Open → Open (required for unsigned apps)
4. Grant permissions when prompted:
   - **Accessibility** — to detect keyboard input
   - **Screen Recording** — to capture screenshots for OCR

### Windows

1. Download `arc-scanner-windows-vX.X.X.zip` from [Releases](https://github.com/LealKevin/Arc-Scanner/releases)
2. Extract to a folder
3. Run `arc-scanner.exe`
4. If SmartScreen appears: Click "More info" → "Run anyway"

## Built With

<p>
  <img src="https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white" alt="Go"/>
  <img src="https://img.shields.io/badge/React-61DAFB?style=flat&logo=react&logoColor=black" alt="React"/>
  <img src="https://img.shields.io/badge/TypeScript-3178C6?style=flat&logo=typescript&logoColor=white" alt="TypeScript"/>
  <img src="https://img.shields.io/badge/Wails-DF0000?style=flat&logo=wails&logoColor=white" alt="Wails"/>
  <img src="https://img.shields.io/badge/Tesseract-5A5A5A?style=flat&logo=tesseract&logoColor=white" alt="Tesseract"/>
</p>

## Development

See [DEVELOPMENT.md](DEVELOPMENT.md) for build instructions and project architecture.

```bash
# Quick start
wails dev
```

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.

## License

[MIT License](LICENSE) — see LICENSE file for details.

## Credits

Item data provided by [metaforge.app](https://metaforge.app)
