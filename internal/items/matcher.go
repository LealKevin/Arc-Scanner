package items

import (
	"strconv"
	"strings"
)

// romanSuffixes are common Roman numeral suffixes for item tiers.
// Ordered from highest to lowest to match longest suffix first.
var romanSuffixes = []string{"IV", "III", "II", "I"}

// Matcher handles matching OCR text to items in the database.
type Matcher struct {
	items []Item
}

// NewMatcher creates a new Matcher with the given item database.
func NewMatcher(items []Item) *Matcher {
	return &Matcher{
		items: items,
	}
}

// CleanOCRText processes raw OCR text into normalized tokens.
// It handles common OCR errors (e.g., '|' misread as 'I') and
// filters to uppercase words only.
func CleanOCRText(text string) []string {
	var tokens []string

	// Use a replacer for efficient multi-replacement
	replacer := strings.NewReplacer(
		"|", "I", // Common OCR misread
		".", "",  // Remove periods
	)

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = replacer.Replace(line)
		words := strings.Fields(line) // Split on any whitespace
		for _, word := range words {
			// Keep only uppercase words (item names are uppercase)
			if word == strings.ToUpper(word) && word != "" {
				tokens = append(tokens, word)
			}
		}
	}

	return tokens
}

// FindItem searches for an item matching the given OCR tokens.
// Returns the matched item or ErrItemNotFound if no match is found.
func (m *Matcher) FindItem(tokens []string) (Item, error) {
	textJoined := strings.Join(tokens, " ")
	var bestMatch Item

	for _, item := range m.items {
		// Convert item ID to search pattern (e.g., "crafting-manual" -> "CRAFTING MANUAL")
		searchName := strings.ToUpper(strings.ReplaceAll(item.ID, "-", " "))
		searchName = strings.ReplaceAll(searchName, "RECIPE", "")

		if strings.Contains(textJoined, searchName) {
			bestMatch = item

			// Check for Roman numeral suffix for tiered items
			for _, suffix := range romanSuffixes {
				if strings.Contains(textJoined, searchName+suffix) {
					bestMatch = item
					break
				}
			}
		}
	}

	if bestMatch.ID != "" {
		return bestMatch, nil
	}

	return Item{}, ErrItemNotFound
}

// ParseQuantity extracts the stack quantity from OCR text.
// Looks for patterns like "5/10" and returns the first number.
// Returns 1 if no quantity is found.
func ParseQuantity(text string) int {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if idx := strings.Index(line, "/"); idx != -1 {
			numStr := strings.TrimSpace(line[:idx])
			if num, err := strconv.Atoi(numStr); err == nil && num > 0 {
				return num
			}
		}
	}
	return 1
}
