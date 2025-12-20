#!/usr/bin/env bash
set -e

ARCH=$(uname -m)
BUILD_DIR="build/darwin"

echo "Bundling Tesseract for macOS ($ARCH) with ALL dependencies..."

# Clean and create directories
rm -rf "$BUILD_DIR/bin" "$BUILD_DIR/lib" "$BUILD_DIR/tessdata"
mkdir -p "$BUILD_DIR/bin"
mkdir -p "$BUILD_DIR/lib"
mkdir -p "$BUILD_DIR/tessdata"

# Detect Homebrew prefix
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

# Copy Tesseract binary
echo "Copying Tesseract binary..."
cp "$BREW_PREFIX/bin/tesseract" "$BUILD_DIR/bin/"

# Track bundled libraries to avoid duplicates (using file instead of associative array for compatibility)
BUNDLED_FILE=$(mktemp)
trap "rm -f $BUNDLED_FILE" EXIT

is_bundled() {
    grep -q "^$1$" "$BUNDLED_FILE" 2>/dev/null
}

mark_bundled() {
    echo "$1" >> "$BUNDLED_FILE"
}

# Function to get non-system dependencies (absolute paths only)
get_deps() {
    local binary="$1"
    otool -L "$binary" 2>/dev/null | awk 'NR>1 {print $1}' | grep -v "^/usr/lib" | grep -v "^/System" | grep -v "^@"
}

# Function to get @rpath dependencies
get_rpath_deps() {
    local binary="$1"
    otool -L "$binary" 2>/dev/null | awk 'NR>1 {print $1}' | grep "^@rpath" | sed 's|^@rpath/||'
}

# Function to resolve library path (handles symlinks and opt paths)
resolve_lib() {
    local lib="$1"
    if [ -f "$lib" ]; then
        echo "$lib"
    elif [ -f "$BREW_PREFIX/lib/$(basename $lib)" ]; then
        echo "$BREW_PREFIX/lib/$(basename $lib)"
    else
        # Try to find in opt subdirs
        local found=$(find "$BREW_PREFIX/opt" -name "$(basename $lib)" -type f 2>/dev/null | head -1)
        if [ -n "$found" ]; then
            echo "$found"
        fi
    fi
}

# Recursive function to bundle dependencies
bundle_deps() {
    local binary="$1"

    # Get absolute path dependencies
    local deps=$(get_deps "$binary")

    for dep in $deps; do
        local depname=$(basename "$dep")

        # Skip if already bundled
        if is_bundled "$depname"; then
            continue
        fi

        # Resolve the actual path
        local resolved=$(resolve_lib "$dep")

        if [ -n "$resolved" ] && [ -f "$resolved" ]; then
            echo "  Bundling: $depname"
            cp "$resolved" "$BUILD_DIR/lib/"
            chmod 755 "$BUILD_DIR/lib/$depname"
            mark_bundled "$depname"

            # Recursively process this library's dependencies
            bundle_deps "$BUILD_DIR/lib/$depname"
        else
            echo "  Warning: Could not find $dep"
        fi
    done

    # Get @rpath dependencies (like libsharpyuv.0.dylib)
    local rpath_deps=$(get_rpath_deps "$binary")

    for depname in $rpath_deps; do
        # Skip if already bundled
        if is_bundled "$depname"; then
            continue
        fi

        # Search for the library in Homebrew
        local resolved=$(resolve_lib "$depname")

        if [ -n "$resolved" ] && [ -f "$resolved" ]; then
            echo "  Bundling (@rpath): $depname"
            cp "$resolved" "$BUILD_DIR/lib/"
            chmod 755 "$BUILD_DIR/lib/$depname"
            mark_bundled "$depname"

            # Recursively process this library's dependencies
            bundle_deps "$BUILD_DIR/lib/$depname"
        else
            echo "  Warning: Could not find @rpath dependency $depname"
        fi
    done
}

# Bundle all dependencies recursively
echo "Bundling dependencies recursively..."
bundle_deps "$BUILD_DIR/bin/tesseract"

echo ""
echo "Relinking libraries to use @executable_path..."

# Fix the tesseract binary
echo "  Relinking tesseract binary..."
for dep in $(get_deps "$BUILD_DIR/bin/tesseract"); do
    depname=$(basename "$dep")
    if [ -f "$BUILD_DIR/lib/$depname" ]; then
        install_name_tool -change "$dep" "@executable_path/../lib/$depname" "$BUILD_DIR/bin/tesseract" 2>/dev/null || true
    fi
done

# Fix each library's install name and dependencies
for lib in "$BUILD_DIR/lib"/*.dylib; do
    if [ -f "$lib" ]; then
        libname=$(basename "$lib")
        echo "  Relinking $libname..."

        # Change the library's own install name
        install_name_tool -id "@executable_path/../lib/$libname" "$lib" 2>/dev/null || true

        # Change all its absolute path dependencies to use @executable_path
        for dep in $(get_deps "$lib"); do
            depname=$(basename "$dep")
            if [ -f "$BUILD_DIR/lib/$depname" ]; then
                install_name_tool -change "$dep" "@executable_path/../lib/$depname" "$lib" 2>/dev/null || true
            fi
        done

        # Change all @rpath dependencies to use @executable_path
        for depname in $(get_rpath_deps "$lib"); do
            if [ -f "$BUILD_DIR/lib/$depname" ]; then
                install_name_tool -change "@rpath/$depname" "@executable_path/../lib/$depname" "$lib" 2>/dev/null || true
            fi
        done
    fi
done

# Download training data
echo ""
echo "Downloading training data..."
if [ ! -f "$BUILD_DIR/tessdata/eng.traineddata" ]; then
    curl -L "https://github.com/tesseract-ocr/tessdata_fast/raw/main/eng.traineddata" \
        -o "$BUILD_DIR/tessdata/eng.traineddata"
    echo "  Downloaded eng.traineddata"
else
    echo "  Training data already exists"
fi

# Verify no homebrew paths or unbundled @rpath remain
echo ""
echo "Verifying all dependencies are bundled..."
MISSING=0
for lib in "$BUILD_DIR/lib"/*.dylib "$BUILD_DIR/bin/tesseract"; do
    if [ -f "$lib" ]; then
        # Check for Homebrew paths
        BAD_DEPS=$(otool -L "$lib" 2>/dev/null | grep "$BREW_PREFIX" | grep -v "^$lib" || true)
        if [ -n "$BAD_DEPS" ]; then
            echo "  WARNING: $(basename $lib) still has unbundled Homebrew deps:"
            echo "$BAD_DEPS"
            MISSING=1
        fi

        # Check for @rpath dependencies that weren't bundled
        RPATH_DEPS=$(get_rpath_deps "$lib")
        for depname in $RPATH_DEPS; do
            if [ ! -f "$BUILD_DIR/lib/$depname" ]; then
                echo "  WARNING: $(basename $lib) has unbundled @rpath dep: $depname"
                MISSING=1
            fi
        done
    fi
done

if [ $MISSING -eq 0 ]; then
    echo "  All dependencies properly bundled!"
fi

echo ""
echo "Tesseract bundling complete!"
echo "Binary: $(ls -lh $BUILD_DIR/bin/tesseract | awk '{print $5}')"
echo "Libraries: $(ls $BUILD_DIR/lib/*.dylib 2>/dev/null | wc -l | xargs) files"
echo "Training data: $(ls -lh $BUILD_DIR/tessdata/eng.traineddata | awk '{print $5}')"
ls -la "$BUILD_DIR/lib/"
