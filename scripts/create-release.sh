#!/bin/bash
set -e

VERSION="${1:-0.1.0}"
RELEASE_DIR="releases/v$VERSION"

echo "Creating release v$VERSION..."
echo ""

# Create release directory
rm -rf "$RELEASE_DIR"
mkdir -p "$RELEASE_DIR"

# Build macOS
echo "=== Building macOS ==="
./scripts/build-bundled.sh
echo ""

# Create macOS zip
echo "Packaging macOS..."
cd build/bin
zip -r "../../$RELEASE_DIR/arc-scanner-macos-v$VERSION.zip" arc-scanner.app
cd ../..
echo "  Created: $RELEASE_DIR/arc-scanner-macos-v$VERSION.zip"
echo ""

# Build Windows (cross-compile)
echo "=== Building Windows ==="
if [ -f "./scripts/build-bundled-windows-cross.sh" ]; then
    ./scripts/build-bundled-windows-cross.sh

    # Create Windows zip with proper structure
    echo "Packaging Windows..."
    cd build/bin
    zip -r "../../$RELEASE_DIR/arc-scanner-windows-v$VERSION.zip" arc-scanner.exe windows/
    cd ../..
    echo "  Created: $RELEASE_DIR/arc-scanner-windows-v$VERSION.zip"
else
    echo "Warning: Windows cross-compile script not found, skipping Windows build"
fi

echo ""
echo "=== Release v$VERSION Complete ==="
echo ""
echo "Files created in $RELEASE_DIR:"
ls -lh "$RELEASE_DIR"
echo ""
echo "Next steps:"
echo "  1. Create a GitHub release at: https://github.com/YOUR_USERNAME/arc-scanner/releases/new"
echo "  2. Tag: v$VERSION"
echo "  3. Upload the zip files from $RELEASE_DIR"
echo "  4. Add release notes with install instructions"
