package text

import (
	"testing"
)

// ═══════════════════════════════════════════════════════════════
//  English Dictionary Tests
// ═══════════════════════════════════════════════════════════════

func TestEnglishDictionary_IsAbbreviation(t *testing.T) {
	dict := NewEnglishDictionary()

	tests := []struct {
		name string
		word string
		want bool
	}{
		{"Title Mr", "Mr.", true},
		{"Title Dr", "Dr.", true},
		{"Title Mrs", "Mrs.", true},
		{"Academic PhD", "Ph.D.", true}, // Now recognized after normalizing periods
		{"Common etc", "etc.", true},
		{"Common ie", "i.e.", true},
		{"Common eg", "e.g.", true},
		{"Not abbreviation", "Hello", false},
		{"Not abbreviation with period", "Hello.", false},
		{"Uppercase", "DR.", true}, // Should be case-insensitive
		{"Lowercase", "dr", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dict.IsAbbreviation(tt.word)
			if got != tt.want {
				t.Errorf("IsAbbreviation(%q) = %v, want %v", tt.word, got, tt.want)
			}
		})
	}
}

func TestEnglishDictionary_AddAbbreviation(t *testing.T) {
	dict := NewEnglishDictionary()

	// Initially not recognized
	if dict.IsAbbreviation("Acme.") {
		t.Error("Should not recognize 'Acme.' initially")
	}

	// Add custom abbreviation
	dict.AddAbbreviation("Acme")

	// Now should be recognized
	if !dict.IsAbbreviation("Acme.") {
		t.Error("Should recognize 'Acme.' after adding")
	}

	// Should work case-insensitively
	if !dict.IsAbbreviation("acme.") {
		t.Error("Should recognize 'acme.' (lowercase)")
	}
}

func TestEnglishDictionary_AddAbbreviations(t *testing.T) {
	dict := NewEnglishDictionary()

	abbrevs := []string{"NASA", "FBI", "CIA"}
	dict.AddAbbreviations(abbrevs)

	for _, abbrev := range abbrevs {
		if !dict.IsAbbreviation(abbrev + ".") {
			t.Errorf("Should recognize %q after adding", abbrev)
		}
	}
}

// ═══════════════════════════════════════════════════════════════
//  Dictionary-Aware Sentence Segmentation Tests
// ═══════════════════════════════════════════════════════════════

func TestSentencesWithDictionary(t *testing.T) {
	txt := NewTerminal()
	dict := NewEnglishDictionary()

	tests := []struct {
		name      string
		input     string
		wantCount int
	}{
		{
			name:      "Single sentence with abbreviation",
			input:     "Dr. Smith is here.",
			wantCount: 1,
		},
		{
			name:      "Multiple abbreviations",
			input:     "Mr. Jones and Mrs. Smith went to the store.",
			wantCount: 1,
		},
		{
			name:      "Abbreviation at end",
			input:     "The meeting is at 3 p.m.",
			wantCount: 1,
		},
		{
			name:      "Multiple sentences",
			input:     "Dr. Smith is here. He is a doctor.",
			wantCount: 2,
		},
		{
			name:      "No abbreviations",
			input:     "Hello world. How are you?",
			wantCount: 2,
		},
		{
			name:      "Common abbreviations",
			input:     "See Fig. 1 for details. The result is shown in Fig. 2.",
			wantCount: 2, // Breaks after "details." which is correct
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sentences := txt.SentencesWithDictionary(tt.input, dict)
			if len(sentences) != tt.wantCount {
				t.Errorf("SentencesWithDictionary() returned %d sentences, want %d",
					len(sentences), tt.wantCount)
				for i, s := range sentences {
					t.Logf("  Sentence %d: %q", i+1, s)
				}
			}
		})
	}
}

func TestSentencesWithDictionary_CustomAbbreviations(t *testing.T) {
	txt := NewTerminal()
	dict := NewEnglishDictionary()

	// Add custom abbreviation
	dict.AddAbbreviation("Acme")

	text := "Acme. Corp is a company. They sell products."
	sentences := txt.SentencesWithDictionary(text, dict)

	// Should treat "Acme." as abbreviation and not break
	if len(sentences) != 2 {
		t.Errorf("Expected 2 sentences, got %d", len(sentences))
		for i, s := range sentences {
			t.Logf("  Sentence %d: %q", i+1, s)
		}
	}
}

func TestSentencesWithDictionary_NilDictionary(t *testing.T) {
	txt := NewTerminal()

	text := "Dr. Smith is here."
	sentences := txt.SentencesWithDictionary(text, nil)

	// Without dictionary, should use raw UAX #29 (splits on "Dr.")
	if len(sentences) < 2 {
		t.Errorf("Expected UAX #29 to split on 'Dr.', got %d sentences", len(sentences))
	}
}

// ═══════════════════════════════════════════════════════════════
//  TextConfig Tests
// ═══════════════════════════════════════════════════════════════

func TestNewTerminalWithEnglishDictionary(t *testing.T) {
	tc := NewTerminalWithEnglishDictionary()

	if tc.Text == nil {
		t.Error("Text should not be nil")
	}

	if tc.Dictionary == nil {
		t.Error("Dictionary should not be nil")
	}

	// Test that it correctly handles abbreviations
	text := "Dr. Smith is here."
	count := tc.SentenceCount(text)

	if count != 1 {
		t.Errorf("Expected 1 sentence with English dictionary, got %d", count)
	}
}

func TestTextConfig_Sentences(t *testing.T) {
	tc := NewTerminalWithEnglishDictionary()

	tests := []struct {
		name      string
		input     string
		wantCount int
	}{
		{
			name:      "With abbreviation",
			input:     "Dr. Smith is here.",
			wantCount: 1,
		},
		{
			name:      "Multiple sentences",
			input:     "Hello. World.",
			wantCount: 2,
		},
		{
			name:      "Mixed",
			input:     "Mr. Jones said hello. Then he left.",
			wantCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sentences := tc.Sentences(tt.input)
			if len(sentences) != tt.wantCount {
				t.Errorf("Sentences() returned %d sentences, want %d",
					len(sentences), tt.wantCount)
			}
		})
	}
}

// ═══════════════════════════════════════════════════════════════
//  Empty Dictionary Tests
// ═══════════════════════════════════════════════════════════════

func TestEmptyDictionary(t *testing.T) {
	dict := &EmptyDictionary{}

	if dict.IsAbbreviation("Dr.") {
		t.Error("EmptyDictionary should not recognize any abbreviations")
	}

	if dict.IsCompoundWord("JavaScript") {
		t.Error("EmptyDictionary should not recognize any compound words")
	}

	if points := dict.GetHyphenationPoints("example"); points != nil {
		t.Error("EmptyDictionary should return nil hyphenation points")
	}
}

// ═══════════════════════════════════════════════════════════════
//  Benchmark Tests
// ═══════════════════════════════════════════════════════════════

func BenchmarkSentencesWithDictionary(b *testing.B) {
	txt := NewTerminal()
	dict := NewEnglishDictionary()
	text := "Dr. Smith and Mrs. Jones went to the store. They bought milk and eggs. The total was $10.50."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txt.SentencesWithDictionary(text, dict)
	}
}

func BenchmarkSentencesWithoutDictionary(b *testing.B) {
	txt := NewTerminal()
	text := "Dr. Smith and Mrs. Jones went to the store. They bought milk and eggs. The total was $10.50."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txt.Sentences(text)
	}
}

func BenchmarkIsAbbreviation(b *testing.B) {
	dict := NewEnglishDictionary()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dict.IsAbbreviation("Dr.")
		dict.IsAbbreviation("Hello")
	}
}
