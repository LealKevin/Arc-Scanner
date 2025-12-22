# Release Workflow

## How to Release a New Version

### Step 1: Update Version Number

Decide on the new version (e.g., `0.2.0`). Follow semantic versioning:
- **MAJOR** (1.0.0): Breaking changes
- **MINOR** (0.2.0): New features, backwards compatible
- **PATCH** (0.1.1): Bug fixes

### Step 2: Build Both Platforms

From the project root:

```bash
# Build macOS
APP_VERSION=0.2.0 ./scripts/build-bundled.sh

# Build Windows (cross-compile from Mac)
./scripts/build-bundled-windows-cross.sh
```

### Step 3: Create Release Zips

```bash
# Create release directory
mkdir -p releases/v0.2.0

# Zip macOS app
cd build/bin
zip -r ../../releases/v0.2.0/arc-scanner-macos-v0.2.0.zip arc-scanner.app

# Zip Windows app
zip -r ../../releases/v0.2.0/arc-scanner-windows-v0.2.0.zip arc-scanner.exe windows/
cd ../..
```

### Step 4: Create GitHub Release

```bash
gh release create v0.2.0 \
  releases/v0.2.0/arc-scanner-macos-v0.2.0.zip \
  releases/v0.2.0/arc-scanner-windows-v0.2.0.zip \
  --title "v0.2.0" \
  --notes "Release notes here..."
```

### Step 5: Verify

1. Check release page: https://github.com/LealKevin/Arc-Scanner/releases
2. Test auto-update: Run an older version and verify it detects the new release

---

## Quick Reference (Copy-Paste)

Replace `0.2.0` with your version number:

```bash
VERSION=0.2.0

# 1. Build
APP_VERSION=$VERSION ./scripts/build-bundled.sh
./scripts/build-bundled-windows-cross.sh

# 2. Package
mkdir -p releases/v$VERSION
cd build/bin
zip -r ../../releases/v$VERSION/arc-scanner-macos-v$VERSION.zip arc-scanner.app
zip -r ../../releases/v$VERSION/arc-scanner-windows-v$VERSION.zip arc-scanner.exe windows/
cd ../..

# 3. Release
gh release create v$VERSION \
  releases/v$VERSION/arc-scanner-macos-v$VERSION.zip \
  releases/v$VERSION/arc-scanner-windows-v$VERSION.zip \
  --title "v$VERSION" \
  --notes "What's new in this release..."
```

---

## Asset Naming Convention

The auto-updater expects these exact filenames:
- `arc-scanner-macos-vX.Y.Z.zip` - macOS bundle
- `arc-scanner-windows-vX.Y.Z.zip` - Windows bundle

---

## Prerequisites

- [GitHub CLI](https://cli.github.com/) (`gh`) installed and authenticated
- For Windows cross-compile: `brew install p7zip` or `brew install innoextract`
