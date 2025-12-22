package config

const (
	// OCR capture area dimensions (relative to mouse cursor position)
	OcrBoxWidth   = 450
	OcrBoxHeight  = 480
	OcrBoxYOffset = 400
	OcrBoxXOffset = 0

	ScanKey   = 'y'
	ToggleKey = 'u'

	MetaForgeAPIBase = "https://metaforge.app/api/arc-raiders/items"
	APIPageSize      = 100
	APIPageCount     = 6

	TesseractPSM       = "3" // Fully automatic page segmentation
	TesseractOEM       = "1" // LSTM only (faster)
	TesseractWhitelist = "0123456789/' ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	ContrastLevel = 20
	SharpenLevel  = 20
)
