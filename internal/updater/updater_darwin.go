//go:build darwin

package updater

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// ApplyUpdate applies the downloaded update on macOS
// It creates a helper script that waits for the app to exit, then copies the new version
func (u *Updater) ApplyUpdate() error {
	if u.downloadedPath == "" {
		return fmt.Errorf("no update downloaded")
	}

	// Find the .app bundle in the extracted directory
	newAppPath := ""
	entries, err := os.ReadDir(u.downloadedPath)
	if err != nil {
		return fmt.Errorf("failed to read extracted directory: %w", err)
	}

	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".app" {
			newAppPath = filepath.Join(u.downloadedPath, entry.Name())
			break
		}
	}

	if newAppPath == "" {
		return fmt.Errorf("no .app bundle found in update")
	}

	// Get the current app path
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Navigate up from MacOS/arc-scanner to the .app bundle
	// execPath = /path/to/arc-scanner.app/Contents/MacOS/arc-scanner
	currentAppPath := filepath.Dir(filepath.Dir(filepath.Dir(execPath)))

	// Create a shell script to perform the update
	scriptContent := fmt.Sprintf(`#!/bin/bash
# Wait for the app to exit
sleep 2

# Remove old app
rm -rf "%s"

# Copy new app
cp -R "%s" "%s"

# Remove quarantine attribute
xattr -cr "%s"

# Relaunch
open "%s"

# Clean up temp files
rm -rf "%s"
rm -- "$0"
`, currentAppPath, newAppPath, currentAppPath, currentAppPath, currentAppPath, filepath.Dir(u.downloadedPath))

	// Write the script to a temp file
	scriptPath := filepath.Join(os.TempDir(), "arc-scanner-update.sh")
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		return fmt.Errorf("failed to write update script: %w", err)
	}

	// Execute the script in background
	cmd := exec.Command("/bin/bash", scriptPath)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start update script: %w", err)
	}

	return nil
}
