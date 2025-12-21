# Windows Build Guide

This guide explains how to build Arc Scanner for Windows as a self-contained executable.

## Prerequisites

### Required Software

1. **Go 1.24+**
   - Download from https://golang.org/dl/
   - Add to PATH

2. **Node.js 18+**
   - Download from https://nodejs.org/
   - Required for frontend build

3. **Wails CLI**
   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```

4. **Tesseract OCR** (for development and bundling)
   - Download installer from: https://github.com/UB-Mannheim/tesseract/wiki
   - Install to default location: `C:\Program Files\Tesseract-OCR`
   - Or set `TESSERACT_PATH` environment variable to custom location

5. **WebView2 Runtime** (usually pre-installed on Windows 10/11)
   - Download from: https://developer.microsoft.com/microsoft-edge/webview2/

## Building

### Option 1: Cross-Compile from macOS/Linux (Recommended!)

**This is the easiest way to build Windows releases if you're on macOS or Linux!**

Build a fully self-contained Windows .exe without needing a Windows machine:

```bash
./scripts/build-bundled-windows-cross.sh
```

**Prerequisites:**
- Install an extraction tool (one of these):
  ```bash
  brew install p7zip        # Provides 7z command
  # OR
  brew install innoextract  # Alternative extractor
  ```

**What it does:**
1. Downloads official Windows Tesseract binaries from UB-Mannheim
2. Extracts tesseract.exe, DLLs, and training data
3. Cross-compiles the Windows .exe using Wails
4. Bundles everything into a portable package

**Output:**
- `build/bin/arc-scanner.exe` - Main executable
- `build/bin/windows/` - Bundled Tesseract and dependencies
- Total size: ~15-20MB

**To distribute:**
- Copy the entire `build/bin` folder
- Users can run `arc-scanner.exe` directly (no installation needed)

---

### Option 2: Self-Contained Build on Windows

If you're building directly on a Windows machine:

```batch
scripts\build-bundled-windows.bat
```

**Prerequisites:**
- Tesseract OCR installed locally (see Prerequisites section)

**Output:**
- `build\bin\arc-scanner.exe` - Main executable
- `build\bin\windows\` - Bundled Tesseract and dependencies

**To distribute:**
- Copy the entire `build\bin` folder
- Users can run `arc-scanner.exe` directly (no installation needed)
- No external dependencies required

### Option 3: Quick Build (Development Only)

This builds without bundling Tesseract. Requires system Tesseract installation.

```bash
wails build -platform windows/amd64
```

**Output:**
- `build\bin\arc-scanner.exe` (~10MB)
- Requires Tesseract installed on the system

## Build Process Details

The self-contained build process has 3 steps:

### Step 1: Bundle Tesseract
```batch
scripts\bundle-tesseract-windows.bat
```
- Locates system Tesseract installation
- Copies `tesseract.exe` binary
- Copies all required DLLs (libleptonica, libtesseract, etc.)
- Downloads training data (`eng.traineddata`)
- Output: `build\windows\` directory

### Step 2: Build with Wails
```bash
wails build -platform windows/amd64
```
- Compiles Go backend
- Builds React frontend
- Creates `arc-scanner.exe`
- Output: `build\bin\arc-scanner.exe`

### Step 3: Copy Bundled Resources
```batch
scripts\post-build-windows.bat
```
- Copies bundled Tesseract to `build\bin\windows\`
- Creates proper directory structure for distribution
- Output: Complete portable package in `build\bin\`

## Cross-Compilation

You can build Windows executables from macOS or Linux:

```bash
# From macOS/Linux
wails build -platform windows/amd64
```

**Note:** You'll need to bundle Tesseract separately on Windows or manually download Windows binaries.

## Directory Structure

After a successful self-contained build:

```
build/
  bin/
    arc-scanner.exe          # Main executable (10-12MB)
    windows/                 # Bundled resources
      bin/
        tesseract.exe        # Tesseract binary (400KB)
        libleptonica-5.dll   # ~2MB
        libtesseract-5.dll   # ~4MB
        [other DLLs]         # ~2MB (all DLLs in same dir as tesseract.exe)
      tessdata/
        eng.traineddata      # ~3.9MB
```

## How It Works

### Tesseract Detection (ocr/scanner.go)

The app automatically detects Tesseract in this order:

1. **Bundled version**: `build\bin\windows\bin\tesseract.exe`
2. **System installation**:
   - `C:\Program Files\Tesseract-OCR\tesseract.exe`
   - `C:\Program Files (x86)\Tesseract-OCR\tesseract.exe`
3. **PATH**: `tesseract.exe` from environment PATH

### App Data Location (app.go)

User data is stored at:
```
C:\Users\<username>\AppData\Roaming\arc-scanner\items.json
```

### Window Management (window_windows.go)

Windows-specific window management using Windows API:
- Sets window to `HWND_TOPMOST` to stay above all windows
- Supports frameless, transparent overlay
- Works with fullscreen games

## Troubleshooting

### Build Fails: "Tesseract not found"

**Solution:**
1. Install Tesseract from https://github.com/UB-Mannheim/tesseract/wiki
2. Or set `TESSERACT_PATH` environment variable:
   ```batch
   set TESSERACT_PATH=C:\path\to\your\tesseract
   ```

### Build Fails: "wails command not found"

**Solution:**
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

Ensure `%GOPATH%\bin` is in your PATH.

### App Doesn't Start: "WebView2 Runtime missing"

**Solution:**
Download and install WebView2 Runtime:
https://developer.microsoft.com/microsoft-edge/webview2/

### Keyboard Hook Not Working

**Possible causes:**
1. Another app is blocking global keyboard hooks
2. Antivirus software is blocking the hook
3. App needs to run as Administrator (rare)

**Solution:**
Try running as Administrator or check antivirus settings.

### Windows Defender Flags the .exe

This is common for unsigned executables.

**For personal use:**
- Add exception in Windows Defender

**For distribution:**
- Code sign the executable with a valid certificate
- Consider building an installer package

## Performance

### Build Time
- Quick build: ~30 seconds
- Self-contained build: ~1-2 minutes (includes downloading training data)

### App Size
- Quick build: ~10MB
- Self-contained build: ~15-20MB (with bundled Tesseract)

### Runtime Performance
- OCR: ~50-100ms per scan (same as macOS)
- Memory: ~50-80MB
- CPU: Minimal when idle

## Security Notes

### Code Signing

For production distribution, code sign the executable:

```batch
signtool sign /f certificate.pfx /p password /tr http://timestamp.digicert.com /td SHA256 /fd SHA256 build\bin\arc-scanner.exe
```

### Antivirus

Some antivirus software may flag:
- Global keyboard hooks (used for 'y' key scanning)
- Screen capture (used for OCR)

These are legitimate features of the app.

## Distribution Checklist

- [ ] Run self-contained build: `scripts\build-bundled-windows.bat`
- [ ] Test on clean Windows machine (without development tools)
- [ ] Verify Tesseract bundling works (no system Tesseract needed)
- [ ] Test keyboard shortcuts ('y' for scan, 'u' for toggle)
- [ ] Test OCR accuracy with game screenshots
- [ ] (Optional) Code sign the executable
- [ ] Package `build\bin` folder for distribution
- [ ] Include README with instructions

## Support

For issues or questions:
- Check CLAUDE.md for project overview
- Review troubleshooting section above
- Check Wails documentation: https://wails.io/docs/
