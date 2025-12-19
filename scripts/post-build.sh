#!/bin/bash
set -e

APP_PATH="build/bin/arc-scanner.app"
BUILD_DIR="build/darwin"

if [ ! -d "$APP_PATH" ]; then
    echo "Error: App bundle not found at $APP_PATH"
    echo "Run 'wails build' first"
    exit 1
fi

echo "Copying bundled resources to app bundle..."

# Clean and copy binaries
echo "  Copying bin/..."
rm -rf "$APP_PATH/Contents/Resources/bin"
mkdir -p "$APP_PATH/Contents/Resources/bin"
cp -r "$BUILD_DIR/bin/"* "$APP_PATH/Contents/Resources/bin/"

# Clean and copy libraries
echo "  Copying lib/..."
rm -rf "$APP_PATH/Contents/Resources/lib"
mkdir -p "$APP_PATH/Contents/Resources/lib"
cp -r "$BUILD_DIR/lib/"* "$APP_PATH/Contents/Resources/lib/"

# Clean and copy training data
echo "  Copying tessdata/..."
rm -rf "$APP_PATH/Contents/Resources/tessdata"
mkdir -p "$APP_PATH/Contents/Resources/tessdata"
cp -r "$BUILD_DIR/tessdata/"* "$APP_PATH/Contents/Resources/tessdata/"

echo "Fixing permissions and signing..."

# Remove quarantine attributes (ignore errors if attributes don't exist)
echo "  Removing quarantine attributes..."
xattr -cr "$APP_PATH/Contents/Resources/bin" 2>/dev/null || true
xattr -cr "$APP_PATH/Contents/Resources/lib" 2>/dev/null || true

# Make binary executable
chmod +x "$APP_PATH/Contents/Resources/bin/tesseract"

# Sign all dylibs first
echo "  Signing libraries..."
for dylib in "$APP_PATH/Contents/Resources/lib"/*.dylib; do
    if [ -f "$dylib" ]; then
        codesign --force --sign - "$dylib" 2>/dev/null || true
    fi
done

# Then sign the tesseract binary
echo "  Signing Tesseract binary..."
codesign --force --deep --sign - "$APP_PATH/Contents/Resources/bin/tesseract" 2>/dev/null || true

echo "Done! Bundled resources copied and signed in $APP_PATH"
