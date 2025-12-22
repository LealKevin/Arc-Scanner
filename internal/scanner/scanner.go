// Package scanner provides OCR-based item scanning using Tesseract.
package scanner

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"arc-scanner/internal/config"

	"github.com/disintegration/imaging"
	"github.com/kbinani/screenshot"
)

// Scanner captures screenshots and performs OCR using Tesseract.
type Scanner interface {
	TakeScreenshot(x, y int) (image.Image, error)
	ProcessImage(img image.Image) (string, error)
}

// TesseractScanner implements Scanner using the Tesseract OCR engine.
type TesseractScanner struct {
	tesseractPath string
	tessdataPath  string
}

// New creates a new TesseractScanner with auto-detected paths.
func New() *TesseractScanner {
	tesseractPath := findTesseractPath()
	tessdataPath := findTessdataPath()

	slog.Info("scanner initialized",
		"tesseract", tesseractPath,
		"tessdata", tessdataPath)

	return &TesseractScanner{
		tesseractPath: tesseractPath,
		tessdataPath:  tessdataPath,
	}
}

// TakeScreenshot captures a region of the screen relative to the cursor position.
// The capture region is defined by config constants.
func (s *TesseractScanner) TakeScreenshot(x, y int) (image.Image, error) {
	captureRect := image.Rect(
		x+config.OcrBoxXOffset,
		y-config.OcrBoxYOffset,
		x+config.OcrBoxXOffset+config.OcrBoxWidth,
		y-config.OcrBoxYOffset+config.OcrBoxHeight,
	)

	img, err := screenshot.CaptureRect(captureRect)
	if err != nil {
		return nil, fmt.Errorf("failed to capture screenshot: %w", err)
	}

	// Preprocess image for better OCR accuracy
	processed := imaging.Grayscale(img)
	processed = imaging.Invert(processed)
	processed = imaging.AdjustContrast(processed, config.ContrastLevel)
	processed = imaging.Sharpen(processed, config.SharpenLevel)

	return processed, nil
}

// ProcessImage performs OCR on the given image and returns the extracted text.
func (s *TesseractScanner) ProcessImage(img image.Image) (string, error) {
	// Encode image to PNG in memory
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", fmt.Errorf("failed to encode image: %w", err)
	}

	// Build Tesseract command
	cmd := exec.Command(
		s.tesseractPath,
		"stdin",  // Read from stdin
		"stdout", // Output to stdout
		"--psm", config.TesseractPSM,
		"--oem", config.TesseractOEM,
		"-c", "tessedit_char_whitelist="+config.TesseractWhitelist,
	)

	// Hide console window on Windows
	hideConsoleWindow(cmd)

	// Set tessdata location if using bundled version
	if s.tessdataPath != "" {
		cmd.Env = append(os.Environ(), "TESSDATA_PREFIX="+s.tessdataPath)
	}

	// Pipe the image to stdin
	cmd.Stdin = &buf

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("OCR failed: %s (stderr: %s)", err, string(exitErr.Stderr))
		}
		return "", fmt.Errorf("OCR failed: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// findTesseractPath locates the Tesseract executable.
// It checks bundled version first, then system installations.
func findTesseractPath() string {
	exe, err := os.Executable()
	if err == nil {
		bundledPath := getBundledPath(exe, "bin", tesseractBinaryName())
		if _, err := os.Stat(bundledPath); err == nil {
			return bundledPath
		}
	}

	// Fallback to system installations
	for _, path := range systemTesseractPaths() {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// Final fallback to PATH
	return tesseractBinaryName()
}

// findTessdataPath locates the Tesseract training data directory.
func findTessdataPath() string {
	exe, err := os.Executable()
	if err == nil {
		tessdataPath := getBundledPath(exe, "tessdata", "")
		if _, err := os.Stat(tessdataPath); err == nil {
			return tessdataPath
		}
	}
	return "" // Let Tesseract use default
}

// getBundledPath constructs a path to bundled resources.
func getBundledPath(exe, subdir, filename string) string {
	var basePath string
	if isWindows() {
		basePath = filepath.Join(filepath.Dir(exe), "windows")
	} else {
		basePath = filepath.Join(filepath.Dir(exe), "..", "Resources")
	}

	if filename != "" {
		return filepath.Join(basePath, subdir, filename)
	}
	return filepath.Join(basePath, subdir)
}

// isWindows returns true if running on Windows.
func isWindows() bool {
	return filepath.Separator == '\\'
}

// tesseractBinaryName returns the platform-specific Tesseract binary name.
func tesseractBinaryName() string {
	if isWindows() {
		return "tesseract.exe"
	}
	return "tesseract"
}

// systemTesseractPaths returns platform-specific system installation paths.
func systemTesseractPaths() []string {
	if isWindows() {
		return []string{
			"C:\\Program Files\\Tesseract-OCR\\tesseract.exe",
			"C:\\Program Files (x86)\\Tesseract-OCR\\tesseract.exe",
		}
	}
	return []string{
		"/opt/homebrew/bin/tesseract",
		"/usr/local/bin/tesseract",
	}
}
