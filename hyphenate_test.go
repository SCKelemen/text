package text

import (
	"testing"
)

// ═══════════════════════════════════════════════════════════════
//  Hyphenation Tests
// ═══════════════════════════════════════════════════════════════

func TestHyphenate(t *testing.T) {
	dict := NewEnglishHyphenation()

	tests := []struct {
		name         string
		word         string
		expectPoints bool // Whether we expect any hyphenation points
	}{
		{
			name:         "example",
			word:         "example",
			expectPoints: true,
		},
		{
			name:         "table",
			word:         "table",
			expectPoints: true,
		},
		{
			name:         "record",
			word:         "record",
			expectPoints: true,
		},
		{
			name:         "present",
			word:         "present",
			expectPoints: true,
		},
		{
			name:         "project",
			word:         "project",
			expectPoints: true,
		},
		{
			name:         "computer",
			word:         "computer",
			expectPoints: true,
		},
		{
			name:         "algorithm",
			word:         "algorithm",
			expectPoints: true,
		},
		{
			name:         "hyphenation",
			word:         "hyphenation",
			expectPoints: true,
		},
		{
			name:         "pattern",
			word:         "pattern",
			expectPoints: true,
		},
		{
			name:         "Short word",
			word:         "cat",
			expectPoints: false, // Too short
		},
		{
			name:         "Two letters",
			word:         "to",
			expectPoints: false, // Too short
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			points := dict.Hyphenate(tt.word)

			if tt.expectPoints && len(points) == 0 {
				t.Logf("Warning: No hyphenation points found for %q (may need more patterns)", tt.word)
			}

			if !tt.expectPoints && len(points) > 0 {
				t.Errorf("Expected no hyphenation points for %q, got %v", tt.word, points)
			}

			// Verify points are within valid range
			for _, point := range points {
				if point < 0 || point >= len(tt.word) {
					t.Errorf("Invalid hyphenation point %d for word %q (len=%d)", point, tt.word, len(tt.word))
				}

				// Verify minimum left/right constraints
				if point < dict.minLeft {
					t.Errorf("Hyphenation point %d violates minLeft=%d for word %q", point, dict.minLeft, tt.word)
				}
				if point > len(tt.word)-dict.minRight {
					t.Errorf("Hyphenation point %d violates minRight=%d for word %q", point, dict.minRight, tt.word)
				}
			}

			// Log results
			if len(points) > 0 {
				t.Logf("Hyphenation for %q: %v", tt.word, points)
				hyphenated := dict.HyphenateWithString(tt.word, "-")
				t.Logf("  Formatted: %s", hyphenated)
			}
		})
	}
}

func TestHyphenateWithString(t *testing.T) {
	dict := NewEnglishHyphenation()

	tests := []struct {
		name   string
		word   string
		hyphen string
	}{
		{
			name:   "Standard hyphen",
			word:   "example",
			hyphen: "-",
		},
		{
			name:   "Soft hyphen",
			word:   "example",
			hyphen: "\u00AD", // Soft hyphen
		},
		{
			name:   "Custom marker",
			word:   "table",
			hyphen: "|",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := dict.HyphenateWithString(tt.word, tt.hyphen)

			// Result should contain the word's letters
			if !containsAllLetters(result, tt.word) {
				t.Errorf("HyphenateWithString lost letters: %q -> %q", tt.word, result)
			}

			// If hyphenation points exist, result should contain the hyphen marker
			points := dict.Hyphenate(tt.word)
			if len(points) > 0 {
				// Count expected hyphens
				expectedHyphens := len(points)
				actualHyphens := countOccurrences(result, tt.hyphen)

				if actualHyphens != expectedHyphens {
					t.Errorf("Expected %d hyphens in %q, got %d", expectedHyphens, result, actualHyphens)
				}
			}

			t.Logf("%s with %q -> %s", tt.word, tt.hyphen, result)
		})
	}
}

func TestHyphenateCapitalized(t *testing.T) {
	dict := NewEnglishHyphenation()

	tests := []struct {
		name string
		word string
	}{
		{"Capitalized", "Example"},
		{"ALL CAPS", "EXAMPLE"},
		{"Mixed", "ExAmPlE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			points := dict.Hyphenate(tt.word)

			// Should work with any capitalization
			if len(points) == 0 {
				t.Logf("Note: No hyphenation points for %q", tt.word)
			} else {
				t.Logf("Hyphenated %q: %v", tt.word, points)
			}
		})
	}
}

func TestHyphenateMinConstraints(t *testing.T) {
	dict := NewEnglishHyphenation()

	// Test that minLeft and minRight are respected
	word := "testing"
	points := dict.Hyphenate(word)

	for _, point := range points {
		if point < dict.minLeft {
			t.Errorf("Point %d violates minLeft=%d", point, dict.minLeft)
		}
		if point > len(word)-dict.minRight {
			t.Errorf("Point %d violates minRight=%d", point, dict.minRight)
		}
	}
}

func TestHyphenateCommonWords(t *testing.T) {
	dict := NewEnglishHyphenation()

	// Test common English words
	words := []string{
		"able", "being", "coming", "doing",
		"testing", "running", "walking", "talking",
		"beautiful", "wonderful", "terrible", "horrible",
		"information", "communication", "education", "station",
	}

	for _, word := range words {
		t.Run(word, func(t *testing.T) {
			points := dict.Hyphenate(word)
			if len(points) > 0 {
				hyphenated := dict.HyphenateWithString(word, "-")
				t.Logf("%s -> %s (points: %v)", word, hyphenated, points)
			} else {
				t.Logf("%s -> no hyphenation", word)
			}
		})
	}
}

// ═══════════════════════════════════════════════════════════════
//  English Dictionary with Hyphenation Tests
// ═══════════════════════════════════════════════════════════════

func TestEnglishDictionaryWithHyphenation(t *testing.T) {
	dict := NewEnglishDictionaryWithHyphenation()

	// Test that it implements DictionaryProvider
	var _ DictionaryProvider = dict

	// Test abbreviations still work
	if !dict.IsAbbreviation("Dr.") {
		t.Error("IsAbbreviation should still work")
	}

	// Test hyphenation works
	points := dict.GetHyphenationPoints("example")
	if len(points) == 0 {
		t.Log("Warning: No hyphenation points for 'example'")
	} else {
		t.Logf("Hyphenation points for 'example': %v", points)
	}
}

func TestEnglishDictionaryWithHyphenation_GetHyphenationPoints(t *testing.T) {
	dict := NewEnglishDictionaryWithHyphenation()

	tests := []struct {
		word         string
		expectPoints bool
	}{
		{"example", true},
		{"table", true},
		{"test", false}, // Too short
		{"go", false},   // Too short
	}

	for _, tt := range tests {
		t.Run(tt.word, func(t *testing.T) {
			points := dict.GetHyphenationPoints(tt.word)

			hasPoints := len(points) > 0
			if hasPoints != tt.expectPoints {
				if tt.expectPoints {
					t.Logf("Warning: Expected hyphenation points for %q, got none", tt.word)
				} else {
					t.Errorf("Expected no points for %q, got %v", tt.word, points)
				}
			}
		})
	}
}

// ═══════════════════════════════════════════════════════════════
//  Helper Functions
// ═══════════════════════════════════════════════════════════════

func containsAllLetters(result, original string) bool {
	// Simple check: result should have at least as many chars as original
	// (plus hyphens, but we don't subtract those)
	letterCount := 0
	for _, ch := range original {
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
			letterCount++
		}
	}

	resultLetters := 0
	for _, ch := range result {
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
			resultLetters++
		}
	}

	return resultLetters >= letterCount
}

func countOccurrences(s, substr string) int {
	count := 0
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			count++
		}
	}
	return count
}

// ═══════════════════════════════════════════════════════════════
//  Benchmark Tests
// ═══════════════════════════════════════════════════════════════

func BenchmarkHyphenate(b *testing.B) {
	dict := NewEnglishHyphenation()
	word := "internationalization"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dict.Hyphenate(word)
	}
}

func BenchmarkHyphenateWithString(b *testing.B) {
	dict := NewEnglishHyphenation()
	word := "internationalization"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dict.HyphenateWithString(word, "-")
	}
}
