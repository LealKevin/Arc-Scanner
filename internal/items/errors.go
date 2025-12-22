package items

import "errors"

var (
	ErrItemNotFound   = errors.New("item not found in OCR text")
	ErrAPIUnavailable = errors.New("failed to fetch items from API")
	ErrCacheCorrupted = errors.New("failed to parse cached items")
)
