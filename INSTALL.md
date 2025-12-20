# Installation Instructions

## Download

Download `arc-scanner.app` from the latest release.

## macOS Security Setup

Since this app is not signed with an Apple Developer ID certificate, macOS will block it on first launch. Follow these steps:

### Method 1: Remove Quarantine (Recommended)

Open Terminal and run:
```bash
xattr -cr /path/to/arc-scanner.app
```
(Drag the .app into Terminal after `xattr -cr ` to get the path automatically)

Then double-click the app to open.

### Method 2: Right-Click to Open

1. Right-click `arc-scanner.app`
2. Click "Open"
3. Click "Open" again in the security dialog
4. This creates a permanent exception

## Permissions Setup

On first launch, the app will request two permissions:

### 1. Accessibility Permission
- Needed to detect keyboard input (pressing 'y' to scan)
- Go to: **System Settings** → **Privacy & Security** → **Accessibility**
- Click the **+** button
- Select `arc-scanner.app`
- Enable the checkbox

### 2. Screen Recording Permission
- Needed to capture screenshots for OCR
- Go to: **System Settings** → **Privacy & Security** → **Screen Recording**
- Click the **+** button
- Select `arc-scanner.app`
- Enable the checkbox

## Usage

1. Launch the app - you'll see a small green dot at the top-right corner
2. In game, hover your mouse over an item
3. Press the **'y'** key to scan
4. The item value will appear in the overlay
5. Press **'u'** to hide/show the overlay

## Troubleshooting

**App closes immediately:**
- Make sure you completed the quarantine removal steps above

**Pressing 'y' does nothing:**
- Check that Accessibility permission is granted and enabled
- Restart the app after granting permissions

**No screenshot captured:**
- Check that Screen Recording permission is granted and enabled
- Restart the app after granting permissions
