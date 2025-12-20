#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
BUILD_DIR="$PROJECT_DIR/build"
INSTALLER_DIR="$PROJECT_DIR/installer"
STAGING_DIR="$BUILD_DIR/installer-staging"
PKG_OUTPUT="$BUILD_DIR/arc-scanner-installer.pkg"

echo "=== Arc Scanner Installer Builder ==="
echo ""

# Step 1: Build the app with bundled dependencies
echo "Step 1: Building app with bundled dependencies..."
"$SCRIPT_DIR/build-bundled.sh"

if [ ! -d "$BUILD_DIR/bin/arc-scanner.app" ]; then
    echo "Error: App not found at $BUILD_DIR/bin/arc-scanner.app"
    exit 1
fi

# Step 2: Prepare staging directory
echo ""
echo "Step 2: Preparing installer staging..."
rm -rf "$STAGING_DIR"
mkdir -p "$STAGING_DIR/Applications"

# Copy the built app to staging
cp -R "$BUILD_DIR/bin/arc-scanner.app" "$STAGING_DIR/Applications/"

# Step 3: Ensure postinstall script is executable
echo ""
echo "Step 3: Preparing installer scripts..."
chmod +x "$INSTALLER_DIR/scripts/postinstall"

# Clear any extended attributes from scripts
xattr -cr "$INSTALLER_DIR/scripts" 2>/dev/null || true

# Step 4: Build the package
echo ""
echo "Step 4: Building installer package..."
rm -f "$PKG_OUTPUT"

pkgbuild \
    --root "$STAGING_DIR" \
    --scripts "$INSTALLER_DIR/scripts" \
    --identifier "com.arc-scanner.app" \
    --version "1.0.0" \
    --install-location "/" \
    "$PKG_OUTPUT"

# Step 5: Verify
echo ""
echo "Step 5: Verifying package..."
if [ -f "$PKG_OUTPUT" ]; then
    PKG_SIZE=$(ls -lh "$PKG_OUTPUT" | awk '{print $5}')
    echo "Success! Installer created:"
    echo "  Path: $PKG_OUTPUT"
    echo "  Size: $PKG_SIZE"
    echo ""
    echo "=== Distribution Instructions ==="
    echo "1. Share arc-scanner-installer.pkg with users"
    echo "2. User double-clicks the .pkg"
    echo "3. User goes to System Settings > Privacy & Security > 'Open Anyway'"
    echo "4. After installation, arc-scanner.app works with double-click!"
else
    echo "Error: Package was not created"
    exit 1
fi
