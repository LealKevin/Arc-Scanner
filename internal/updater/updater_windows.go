//go:build windows

package updater

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// ApplyUpdate applies the downloaded update on Windows
// It creates a batch script that waits for the app to exit, then copies the new version
func (u *Updater) ApplyUpdate() error {
	if u.downloadedPath == "" {
		return fmt.Errorf("no update downloaded")
	}

	// Get the current app directory
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	appDir := filepath.Dir(execPath)
	appName := filepath.Base(execPath)

	// Create a batch script to perform the update
	// The script waits for the app to exit, copies new files, and relaunches
	scriptContent := fmt.Sprintf(`@echo off
setlocal

REM Wait for the app to exit
timeout /t 2 /nobreak >nul

REM Copy new files (xcopy handles directories)
xcopy /E /Y /Q "%s\*" "%s\"

REM Relaunch the app
start "" "%s"

REM Clean up temp files
rmdir /S /Q "%s"
del "%%~f0"
`, u.downloadedPath, appDir, filepath.Join(appDir, appName), filepath.Dir(u.downloadedPath))

	// Write the script to a temp file
	scriptPath := filepath.Join(os.TempDir(), "arc-scanner-update.bat")
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		return fmt.Errorf("failed to write update script: %w", err)
	}

	// Execute the script in background using cmd /c start
	cmd := exec.Command("cmd", "/c", "start", "/b", scriptPath)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start update script: %w", err)
	}

	return nil
}
