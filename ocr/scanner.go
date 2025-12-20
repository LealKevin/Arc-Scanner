package ocr

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/kbinani/screenshot"
)

type Scanner struct {
	tesseractPath string
}

func findTesseractPath() string {
	exe, err := os.Executable()
	if err == nil {
		// Try bundled version first (platform-specific)
		var bundledPath string

		if filepath.Separator == '\\' {
			// Windows: arc-scanner.exe is in build/bin/
			// Look for build/windows/bin/tesseract.exe
			bundledPath = filepath.Join(
				filepath.Dir(exe),
				"..", "windows", "bin", "tesseract.exe",
			)
		} else {
			// macOS: From arc-scanner.app/Contents/MacOS/arc-scanner
			// Navigate to arc-scanner.app/Contents/Resources/bin/tesseract
			bundledPath = filepath.Join(
				filepath.Dir(exe),
				"..", "Resources", "bin", "tesseract",
			)
		}

		if _, err := os.Stat(bundledPath); err == nil {
			return bundledPath
		}
	}

	// Fallback to system installations (for development)
	if filepath.Separator == '\\' {
		// Windows system paths
		windowsPaths := []string{
			"C:\\Program Files\\Tesseract-OCR\\tesseract.exe",
			"C:\\Program Files (x86)\\Tesseract-OCR\\tesseract.exe",
		}
		for _, path := range windowsPaths {
			if _, err := os.Stat(path); err == nil {
				return path
			}
		}
	} else {
		// macOS Homebrew paths
		if _, err := os.Stat("/opt/homebrew/bin/tesseract"); err == nil {
			return "/opt/homebrew/bin/tesseract"
		}
		if _, err := os.Stat("/usr/local/bin/tesseract"); err == nil {
			return "/usr/local/bin/tesseract"
		}
	}

	// Final fallback to PATH
	if filepath.Separator == '\\' {
		return "tesseract.exe"
	}
	return "tesseract"
}

func getTessdataPath() string {
	exe, err := os.Executable()
	if err == nil {
		var tessdataPath string

		if filepath.Separator == '\\' {
			// Windows: build/windows/tessdata
			tessdataPath = filepath.Join(
				filepath.Dir(exe),
				"..", "windows", "tessdata",
			)
		} else {
			// macOS: arc-scanner.app/Contents/Resources/tessdata
			tessdataPath = filepath.Join(
				filepath.Dir(exe),
				"..", "Resources", "tessdata",
			)
		}

		if _, err := os.Stat(tessdataPath); err == nil {
			return tessdataPath
		}
	}
	return "" // Let Tesseract use default
}

func NewScanner() *Scanner {
	tesseractPath := findTesseractPath()
	fmt.Printf("Using Tesseract at: %s\n", tesseractPath)
	return &Scanner{
		tesseractPath: tesseractPath,
	}
}

func (s *Scanner) Close() {
}

func (s *Scanner) TakeScreenshot(x, y int) (image.Image, error) {
	captureRect := image.Rect(x, y-400, x+450, y-200+380)

	img, err := screenshot.CaptureRect(captureRect)
	if err != nil {
		return nil, fmt.Errorf("failed to capture screenshot: %w", err)
	}

	greyImg := imaging.Grayscale(img)
	greyImg = imaging.Invert(greyImg)
	greyImg = imaging.AdjustContrast(greyImg, 20)
	greyImg = imaging.Sharpen(greyImg, 20)

	return greyImg, nil
}

func (s *Scanner) ProcessImage(img image.Image) (string, error) {
	// Encode image to PNG in memory
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", fmt.Errorf("failed to encode image: %w", err)
	}

	// Create Tesseract command with stdin input
	cmd := exec.Command(
		s.tesseractPath,
		"stdin",      // Read from stdin instead of file
		"stdout",     // Output to stdout
		"--psm", "3", // Fully automatic page segmentation
		"--oem", "1", // LSTM only (faster)
		"-c", "tessedit_char_whitelist=0123456789/ ABCDEFGHIJKLMNOPQRSTUVWXYZ",
	)

	// Set tessdata location if using bundled version
	if tessdataPath := getTessdataPath(); tessdataPath != "" {
		cmd.Env = append(os.Environ(),
			"TESSDATA_PREFIX="+tessdataPath)
	}

	// Pipe the image buffer to stdin
	cmd.Stdin = &buf

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("OCR failed: %s (stderr: %s)", err, string(exitErr.Stderr))
		}
		return "", fmt.Errorf("OCR failed: %w", err)
	}

	text := strings.TrimSpace(string(output))
	return text, nil
}
