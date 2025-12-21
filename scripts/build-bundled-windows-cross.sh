#!/bin/bash
set -e

echo "==============================================="
echo "Building self-contained Arc Scanner for Windows"
echo "(Cross-compilation from macOS/Linux)"
echo "==============================================="
echo ""

BUILD_DIR="build/windows"
TEMP_DIR="build/temp-tesseract"
TESSERACT_VERSION="5.3.3.20231005"
TESSERACT_URL="https://digi.bib.uni-mannheim.de/tesseract/tesseract-ocr-w64-setup-${TESSERACT_VERSION}.exe"

# Function to find mingw-w64 gcc
find_mingw_gcc() {
    local HOMEBREW_ARM="/opt/homebrew/bin/x86_64-w64-mingw32-gcc"
    local HOMEBREW_INTEL="/usr/local/bin/x86_64-w64-mingw32-gcc"

    if [ -f "$HOMEBREW_ARM" ]; then
        echo "$HOMEBREW_ARM"
    elif [ -f "$HOMEBREW_INTEL" ]; then
        echo "$HOMEBREW_INTEL"
    elif command -v x86_64-w64-mingw32-gcc &> /dev/null; then
        echo "x86_64-w64-mingw32-gcc"
    fi
}

# Function to find mingw-w64 g++
find_mingw_gxx() {
    local HOMEBREW_ARM="/opt/homebrew/bin/x86_64-w64-mingw32-g++"
    local HOMEBREW_INTEL="/usr/local/bin/x86_64-w64-mingw32-g++"

    if [ -f "$HOMEBREW_ARM" ]; then
        echo "$HOMEBREW_ARM"
    elif [ -f "$HOMEBREW_INTEL" ]; then
        echo "$HOMEBREW_INTEL"
    elif command -v x86_64-w64-mingw32-g++ &> /dev/null; then
        echo "x86_64-w64-mingw32-g++"
    fi
}

# Check for mingw-w64 (required for CGO cross-compilation)
echo "Checking for mingw-w64..."
MINGW_GCC=$(find_mingw_gcc)
MINGW_GXX=$(find_mingw_gxx)

if [ -z "$MINGW_GCC" ] || [ -z "$MINGW_GXX" ]; then
    echo ""
    echo "========================================="
    echo "ERROR: mingw-w64 is required!"
    echo "========================================="
    echo ""
    echo "This project uses robotgo and gohook which require CGO."
    echo "To cross-compile for Windows, install mingw-w64:"
    echo ""
    echo "  brew install mingw-w64"
    echo ""
    exit 1
fi

echo "  ✓ Found: $MINGW_GCC"
echo ""

# Step 1: Download and extract Windows Tesseract
echo "Step 1/4: Downloading Windows Tesseract binaries..."
echo "-----------------------------------------------"

# Clean and create directories
rm -rf "$BUILD_DIR/bin" "$BUILD_DIR/lib" "$BUILD_DIR/tessdata" "$TEMP_DIR"
mkdir -p "$BUILD_DIR/bin"
mkdir -p "$BUILD_DIR/lib"
mkdir -p "$BUILD_DIR/tessdata"
mkdir -p "$TEMP_DIR"

# Download Tesseract installer
echo "Downloading Tesseract ${TESSERACT_VERSION}..."
if ! curl -L -o "$TEMP_DIR/tesseract-installer.exe" "$TESSERACT_URL"; then
    echo "Error: Failed to download Tesseract installer"
    echo "Trying alternative method with 7z extraction..."

    # Alternative: Download portable zip if available
    echo "Please ensure you have 7z or unzip installed for extraction"
    exit 1
fi

# Check if 7z is available for extraction
if command -v 7z &> /dev/null; then
    echo "Extracting Tesseract installer with 7z..."
    cd "$TEMP_DIR"
    7z x tesseract-installer.exe -o"extracted" > /dev/null 2>&1 || {
        echo "Warning: 7z extraction had errors, continuing..."
    }
    cd - > /dev/null

    # Find and copy tesseract.exe
    if [ -f "$TEMP_DIR/extracted/tesseract.exe" ]; then
        cp "$TEMP_DIR/extracted/tesseract.exe" "$BUILD_DIR/bin/"
        echo "  ✓ tesseract.exe copied"
    else
        echo "Error: Could not find tesseract.exe in installer"
        exit 1
    fi

    # Copy DLLs
    echo "Copying DLL dependencies..."
    find "$TEMP_DIR/extracted" -name "*.dll" -exec cp {} "$BUILD_DIR/lib/" \; 2>/dev/null || true
    echo "  ✓ DLLs copied"

    # Copy traineddata if available
    if [ -f "$TEMP_DIR/extracted/tessdata/eng.traineddata" ]; then
        cp "$TEMP_DIR/extracted/tessdata/eng.traineddata" "$BUILD_DIR/tessdata/"
        echo "  ✓ eng.traineddata copied from installer"
    fi

elif command -v innoextract &> /dev/null; then
    echo "Extracting Tesseract installer with innoextract..."
    cd "$TEMP_DIR"
    innoextract tesseract-installer.exe > /dev/null 2>&1
    cd - > /dev/null

    # Copy files from innoextract output
    if [ -d "$TEMP_DIR/app" ]; then
        cp "$TEMP_DIR/app/tesseract.exe" "$BUILD_DIR/bin/" 2>/dev/null || true
        cp "$TEMP_DIR/app"/*.dll "$BUILD_DIR/lib/" 2>/dev/null || true
        cp "$TEMP_DIR/app/tessdata/eng.traineddata" "$BUILD_DIR/tessdata/" 2>/dev/null || true
        echo "  ✓ Files extracted with innoextract"
    fi
else
    echo ""
    echo "========================================="
    echo "WARNING: No extraction tool found!"
    echo "========================================="
    echo "Please install one of:"
    echo "  - 7z:          brew install p7zip"
    echo "  - innoextract: brew install innoextract"
    echo ""
    echo "Or manually download and extract Windows Tesseract:"
    echo "  1. Download from: https://github.com/UB-Mannheim/tesseract/wiki"
    echo "  2. Extract and copy files to:"
    echo "     - tesseract.exe -> $BUILD_DIR/bin/"
    echo "     - *.dll -> $BUILD_DIR/lib/"
    echo "     - tessdata/eng.traineddata -> $BUILD_DIR/tessdata/"
    echo ""
    read -p "Press Enter if you've manually extracted files, or Ctrl+C to exit..."
fi

# Download traineddata if not already present
if [ ! -f "$BUILD_DIR/tessdata/eng.traineddata" ]; then
    echo "Downloading eng.traineddata..."
    curl -L "https://github.com/tesseract-ocr/tessdata_fast/raw/main/eng.traineddata" \
        -o "$BUILD_DIR/tessdata/eng.traineddata"
    echo "  ✓ eng.traineddata downloaded"
fi

echo ""

# Step 2: Verify extracted files
echo "Step 2/4: Verifying extracted files..."
echo "-----------------------------------------------"

if [ ! -f "$BUILD_DIR/bin/tesseract.exe" ]; then
    echo "Error: tesseract.exe not found!"
    echo "Expected at: $BUILD_DIR/bin/tesseract.exe"
    exit 1
fi

DLL_COUNT=$(find "$BUILD_DIR/lib" -name "*.dll" 2>/dev/null | wc -l | tr -d ' ')
if [ "$DLL_COUNT" -lt 2 ]; then
    echo "Warning: Only $DLL_COUNT DLL files found. Expected at least 2 (libleptonica, libtesseract)"
    echo "The build may not work without required DLLs"
fi

if [ ! -f "$BUILD_DIR/tessdata/eng.traineddata" ]; then
    echo "Error: eng.traineddata not found!"
    exit 1
fi

echo "  ✓ tesseract.exe: $(ls -lh $BUILD_DIR/bin/tesseract.exe | awk '{print $5}')"
echo "  ✓ DLL files: $DLL_COUNT"
echo "  ✓ traineddata: $(ls -lh $BUILD_DIR/tessdata/eng.traineddata | awk '{print $5}')"
echo ""

# Step 3: Build with Wails
echo "Step 3/4: Building app with Wails for Windows..."
echo "-----------------------------------------------"
echo "Cross-compiling with CGO enabled..."
echo "  CC=$MINGW_GCC"
echo "  CXX=$MINGW_GXX"
echo ""

CGO_ENABLED=1 \
CC="$MINGW_GCC" \
CXX="$MINGW_GXX" \
GOOS=windows \
GOARCH=amd64 \
wails build -platform windows/amd64
echo ""

# Step 4: Copy bundled resources
echo "Step 4/4: Copying bundled resources to app..."
echo "-----------------------------------------------"

BIN_DIR="build/bin"

if [ ! -f "$BIN_DIR/arc-scanner.exe" ]; then
    echo "Error: arc-scanner.exe not found at $BIN_DIR/arc-scanner.exe"
    echo "Wails build may have failed"
    exit 1
fi

# Create directory structure next to the .exe
echo "Creating directory structure..."
mkdir -p "$BIN_DIR/windows/bin"
mkdir -p "$BIN_DIR/windows/tessdata"

# Copy bundled resources
# Note: DLLs must be in the same directory as tesseract.exe for Windows to find them
echo "Copying resources..."
cp "$BUILD_DIR/bin/tesseract.exe" "$BIN_DIR/windows/bin/" 2>/dev/null && echo "  ✓ Copied tesseract.exe" || echo "  ✗ Failed to copy tesseract.exe"
cp "$BUILD_DIR/lib"/* "$BIN_DIR/windows/bin/" 2>/dev/null && echo "  ✓ Copied DLLs to bin/" || echo "  ✗ Failed to copy DLLs"
cp "$BUILD_DIR/tessdata/eng.traineddata" "$BIN_DIR/windows/tessdata/" 2>/dev/null && echo "  ✓ Copied tessdata/" || echo "  ✗ Failed to copy tessdata/"

# Cleanup temp directory
echo "Cleaning up..."
rm -rf "$TEMP_DIR"

echo ""
echo "==============================================="
echo "Build complete!"
echo "==============================================="
echo "  App location: $BIN_DIR/arc-scanner.exe"
echo "  App is self-contained with bundled Tesseract"
echo ""
echo "Directory structure:"
echo "  build/bin/"
echo "  ├── arc-scanner.exe              (~10-12MB)"
echo "  └── windows/"
echo "      ├── bin/"
echo "      │   ├── tesseract.exe        (~400KB)"
echo "      │   └── *.dll                (~8MB)"
echo "      └── tessdata/eng.traineddata (~3.9MB)"
echo ""
echo "To distribute:"
echo "  - Copy the entire 'build/bin' folder"
echo "  - Users can run arc-scanner.exe directly"
echo "  - No installation or dependencies needed"
echo ""
echo "Total package size: ~15-20MB"
