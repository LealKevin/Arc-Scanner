package items

import (
	"testing"
)

func TestCleanOCRText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "basic uppercase words",
			input:    "CRAFTING MANUAL",
			expected: []string{"CRAFTING", "MANUAL"},
		},
		{
			name:     "filters lowercase words",
			input:    "ITEM name Value",
			expected: []string{"ITEM"},
		},
		{
			name:     "handles pipe as I",
			input:    "P|PE WRENCH",
			expected: []string{"PIPE", "WRENCH"},
		},
		{
			name:     "removes periods",
			input:    "A.R.C. SCANNER",
			expected: []string{"ARC", "SCANNER"},
		},
		{
			name:     "handles multiline",
			input:    "LINE ONE\nLINE TWO",
			expected: []string{"LINE", "ONE", "LINE", "TWO"},
		},
		{
			name:     "handles quantity format",
			input:    "ITEM 5/10",
			expected: []string{"ITEM", "5/10"},
		},
		{
			name:     "empty input",
			input:    "",
			expected: nil,
		},
		{
			name:     "only lowercase",
			input:    "all lowercase words",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CleanOCRText(tt.input)
			if !slicesEqual(result, tt.expected) {
				t.Errorf("CleanOCRText(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseQuantity(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "simple quantity",
			input:    "5/10",
			expected: 5,
		},
		{
			name:     "full stack",
			input:    "15/15",
			expected: 15,
		},
		{
			name:     "quantity with text",
			input:    "ITEM NAME\n5/10",
			expected: 5,
		},
		{
			name:     "no quantity returns 1",
			input:    "ITEM NAME",
			expected: 1,
		},
		{
			name:     "invalid format",
			input:    "abc/def",
			expected: 1,
		},
		{
			name:     "empty string",
			input:    "",
			expected: 1,
		},
		{
			name:     "quantity with spaces",
			input:    " 10 /20",
			expected: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseQuantity(tt.input)
			if result != tt.expected {
				t.Errorf("ParseQuantity(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMatcher_FindItem(t *testing.T) {
	testItems := []Item{
		{ID: "crafting-manual", Name: "Crafting Manual", Value: 100},
		{ID: "pipe-wrench", Name: "Pipe Wrench", Value: 50},
		{ID: "combat-knife-i", Name: "Combat Knife I", Value: 200},
		{ID: "combat-knife-ii", Name: "Combat Knife II", Value: 300},
	}

	matcher := NewMatcher(testItems)

	tests := []struct {
		name        string
		tokens      []string
		expectedID  string
		expectError bool
	}{
		{
			name:       "exact match",
			tokens:     []string{"CRAFTING", "MANUAL"},
			expectedID: "crafting-manual",
		},
		{
			name:       "partial match in longer text",
			tokens:     []string{"SOME", "CRAFTING", "MANUAL", "TEXT"},
			expectedID: "crafting-manual",
		},
		{
			name:       "item with Roman numeral",
			tokens:     []string{"COMBAT", "KNIFE", "II"},
			expectedID: "combat-knife-ii",
		},
		{
			name:        "no match",
			tokens:      []string{"UNKNOWN", "ITEM"},
			expectError: true,
		},
		{
			name:        "empty tokens",
			tokens:      []string{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, err := matcher.FindItem(tt.tokens)

			if tt.expectError {
				if err == nil {
					t.Errorf("FindItem(%v) expected error, got item %v", tt.tokens, item)
				}
				return
			}

			if err != nil {
				t.Errorf("FindItem(%v) unexpected error: %v", tt.tokens, err)
				return
			}

			if item.ID != tt.expectedID {
				t.Errorf("FindItem(%v) = %s, want %s", tt.tokens, item.ID, tt.expectedID)
			}
		})
	}
}

func TestBuildIndex(t *testing.T) {
	items := []Item{
		{ID: "item-1", Name: "Item One", Value: 100},
		{ID: "item-2", Name: "Item Two", Value: 200},
	}

	index := BuildIndex(items)

	if len(index) != 2 {
		t.Errorf("BuildIndex created %d entries, want 2", len(index))
	}

	item, ok := index.Get("item-1")
	if !ok {
		t.Error("BuildIndex: item-1 not found in index")
	}
	if item.Value != 100 {
		t.Errorf("BuildIndex: item-1 value = %d, want 100", item.Value)
	}

	_, ok = index.Get("nonexistent")
	if ok {
		t.Error("BuildIndex: nonexistent item should not be found")
	}
}

// slicesEqual compares two string slices for equality.
func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
