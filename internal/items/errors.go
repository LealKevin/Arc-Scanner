package items

import "errors"

// Sentinel errors for item operations.
var (
	// ErrItemNotFound is returned when an item cannot be matched from OCR text.
	ErrItemNotFound = errors.New("item not found in OCR text")

	// ErrAPIUnavailable is returned when the items API cannot be reached.
	ErrAPIUnavailable = errors.New("failed to fetch items from API")

	// ErrCacheCorrupted is returned when the cached items file cannot be parsed.
	ErrCacheCorrupted = errors.New("failed to parse cached items")
)
