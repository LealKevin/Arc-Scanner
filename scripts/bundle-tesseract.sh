#!/bin/bash
set -e

ARCH=$(uname -m)
BUILD_DIR="build/darwin"

echo "Bundling Tesseract for macOS ($ARCH)..."

# Clean and create directories
rm -rf "$BUILD_DIR/bin" "$BUILD_DIR/lib" "$BUILD_DIR/tessdata"
mkdir -p "$BUILD_DIR/bin"
mkdir -p "$BUILD_DIR/lib"
mkdir -p "$BUILD_DIR/tessdata"

# Copy Tesseract binary
if [ "$ARCH" = "arm64" ]; then
    BREW_PREFIX="/opt/homebrew"
else
    BREW_PREFIX="/usr/local"
fi

if [ ! -f "$BREW_PREFIX/bin/tesseract" ]; then
    echo "Error: Tesseract not found at $BREW_PREFIX/bin/tesseract"
    echo "Please install Tesseract: brew install tesseract"
    exit 1
fi

echo "Copying Tesseract binary from $BREW_PREFIX/bin/tesseract..."
cp "$BREW_PREFIX/bin/tesseract" "$BUILD_DIR/bin/"

# Find and copy dynamic libraries
echo "Copying dependencies..."
otool -L "$BUILD_DIR/bin/tesseract" | grep -v "/usr/lib" | grep -v "/System" | awk 'NR>1 {print $1}' | while read lib; do
    if [ -f "$lib" ]; then
        echo "  Copying $(basename $lib)..."
        cp "$lib" "$BUILD_DIR/lib/"
    fi
done

# Relink to use @executable_path
echo "Relinking libraries..."
for dylib in "$BUILD_DIR/lib"/*.dylib; do
    if [ -f "$dylib" ]; then
        libname=$(basename "$dylib")
        echo "  Relinking $libname..."

        # Get the original path from otool
        original_path=$(otool -L "$BUILD_DIR/bin/tesseract" | grep "$libname" | awk '{print $1}')

        if [ ! -z "$original_path" ]; then
            install_name_tool -change "$original_path" \
                "@executable_path/../lib/$libname" \
                "$BUILD_DIR/bin/tesseract"
        fi

        # Also fix inter-library dependencies
        otool -L "$dylib" | grep -v "/usr/lib" | grep -v "/System" | grep "\.dylib" | awk 'NR>1 {print $1}' | while read dep; do
            depname=$(basename "$dep")
            if [ -f "$BUILD_DIR/lib/$depname" ]; then
                install_name_tool -change "$dep" \
                    "@executable_path/../lib/$depname" \
                    "$dylib"
            fi
        done
    fi
done

# Download training data
echo "Downloading training data..."
if [ ! -f "$BUILD_DIR/tessdata/eng.traineddata" ]; then
    curl -L "https://github.com/tesseract-ocr/tessdata_fast/raw/main/eng.traineddata" \
        -o "$BUILD_DIR/tessdata/eng.traineddata"
    echo "  Downloaded eng.traineddata"
else
    echo "  Training data already exists"
fi

echo ""
echo "Tesseract bundling complete!"
echo "Binary: $(ls -lh $BUILD_DIR/bin/tesseract | awk '{print $5}')"
echo "Libraries: $(ls $BUILD_DIR/lib/*.dylib | wc -l | xargs) files"
echo "Training data: $(ls -lh $BUILD_DIR/tessdata/eng.traineddata | awk '{print $5}')"
