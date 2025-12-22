// Package config provides centralized configuration constants for arc-scanner.
package config

const (
	// OCR capture area dimensions (relative to mouse cursor position)
	OcrBoxWidth   = 450
	OcrBoxHeight  = 380
	OcrBoxYOffset = 400
	OcrBoxXOffset = 0

	// Keyboard shortcuts
	ScanKey   = 'y'
	ToggleKey = 'u'

	// API configuration
	MetaForgeAPIBase = "https://metaforge.app/api/arc-raiders/items"
	APIPageSize      = 100
	APIPageCount     = 6

	// Tesseract configuration
	TesseractPSM       = "3" // Fully automatic page segmentation
	TesseractOEM       = "1" // LSTM only (faster)
	TesseractWhitelist = "0123456789/ ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Image processing parameters
	ContrastLevel = 20
	SharpenLevel  = 20
)
