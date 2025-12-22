#!/bin/bash
set -e

APP_PATH="build/bin/arc-scanner.app"
BUILD_DIR="build/darwin"

if [ ! -d "$APP_PATH" ]; then
    echo "Error: App bundle not found at $APP_PATH"
    echo "Run 'wails build' first"
    exit 1
fi

# Verify the main executable exists
EXEC_PATH="$APP_PATH/Contents/MacOS/arc-scanner"
if [ ! -f "$EXEC_PATH" ]; then
    echo "Error: Executable not found at $EXEC_PATH"
    echo "Wails build may have failed - check for errors above"
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

# Make all bundled files writable (needed for xattr and codesign)
chmod -R u+w "$APP_PATH/Contents/Resources"

# Make tesseract binary executable
chmod +x "$APP_PATH/Contents/Resources/bin/tesseract"

# Sign in correct order: innermost to outermost
# 1. Sign all dylibs first
echo "  Signing libraries..."
for dylib in "$APP_PATH/Contents/Resources/lib"/*.dylib; do
    if [ -f "$dylib" ]; then
        codesign --force --sign - "$dylib"
    fi
done

# 2. Sign tesseract binary
echo "  Signing tesseract..."
codesign --force --sign - "$APP_PATH/Contents/Resources/bin/tesseract"

# 3. Sign main executable
echo "  Signing main executable..."
codesign --force --sign - "$APP_PATH/Contents/MacOS/arc-scanner"

# 4. Sign entire app bundle
echo "  Signing app bundle..."
codesign --force --sign - "$APP_PATH"

# 5. Remove ALL extended attributes (must be after signing)
echo "  Removing extended attributes..."
xattr -cr "$APP_PATH"

echo "Done! App bundle ready at $APP_PATH"
echo ""
echo "Note: Users on other Macs will need to right-click -> Open on first launch"
