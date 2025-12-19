#!/bin/bash
set -e

echo "Building self-contained Arc Scanner..."
echo ""

# Step 1: Bundle Tesseract
echo "Step 1/3: Bundling Tesseract..."
./scripts/bundle-tesseract.sh
echo ""

# Step 2: Build with Wails
echo "Step 2/3: Building app with Wails..."
wails build
echo ""

# Step 3: Copy bundled resources
echo "Step 3/3: Copying bundled resources to app..."
./scripts/post-build.sh
echo ""

echo "✓ Build complete!"
echo "  App location: build/bin/arc-scanner.app"
echo "  App is self-contained with bundled Tesseract"
