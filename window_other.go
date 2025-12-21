//go:build !darwin && !windows

package main

// Stub for platforms other than macOS and Windows (e.g., Linux)
func setWindowAboveFullscreen() {
	// No-op on non-macOS platforms
}
