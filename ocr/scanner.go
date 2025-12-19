package ocr

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"os/exec"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/kbinani/screenshot"
)

type Scanner struct {
	tesseractPath string
}

func NewScanner() *Scanner {
	tesseractPath := "tesseract" // Default to PATH

	if _, err := os.Stat("/opt/homebrew/bin/tesseract"); err == nil {
		tesseractPath = "/opt/homebrew/bin/tesseract"
	} else if _, err := os.Stat("/usr/local/bin/tesseract"); err == nil {
		tesseractPath = "/usr/local/bin/tesseract"
	}

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
