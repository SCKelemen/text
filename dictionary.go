package text

import (
	"strings"

	"github.com/SCKelemen/unicode/uax29"
)

// Dictionary Support for Text Segmentation
//
// This provides extensible dictionary support for:
// - Abbreviation detection (Dr., Mrs., etc.)
// - Language-specific word breaking
// - Hyphenation dictionaries
// - Domain-specific terminology

// ═══════════════════════════════════════════════════════════════
//  Dictionary Provider Interface
// ═══════════════════════════════════════════════════════════════

// DictionaryProvider provides domain-specific dictionaries for text segmentation.
//
// This allows extending UAX #29 sentence and word boundary detection with:
// - Common abbreviations (Dr., Mrs., Ph.D., etc.)
// - Language-specific rules
// - Domain-specific terminology
// - Hyphenation patterns
type DictionaryProvider interface {
	// IsAbbreviation returns true if the word is a known abbreviation.
	// This is used to prevent sentence breaks after abbreviated words.
	//
	// Example: "Dr." -> true, "Mr." -> true, "hello" -> false
	IsAbbreviation(word string) bool

	// GetHyphenationPoints returns hyphenation points for a word.
	// Returns slice of byte indices where hyphenation is allowed.
	//
	// Example: "example" -> []int{2, 4} (ex-am-ple)
	GetHyphenationPoints(word string) []int

	// IsCompoundWord returns true if the word is a compound that shouldn't be broken.
	// Used for language-specific compound word handling.
	//
	// Example: "JavaScript" -> true
	IsCompoundWord(word string) bool
}

// ═══════════════════════════════════════════════════════════════
//  Built-in English Dictionary
// ═══════════════════════════════════════════════════════════════

// EnglishDictionary provides common English abbreviations and rules.
type EnglishDictionary struct {
	abbreviations map[string]bool
	customWords   map[string]bool
}

// NewEnglishDictionary creates a dictionary with common English abbreviations.
func NewEnglishDictionary() *EnglishDictionary {
	return &EnglishDictionary{
		abbreviations: defaultEnglishAbbreviations(),
		customWords:   make(map[string]bool),
	}
}

// defaultEnglishAbbreviations returns common English abbreviations.
func defaultEnglishAbbreviations() map[string]bool {
	abbrevs := map[string]bool{
		// Titles
		"mr":   true,
		"mrs":  true,
		"ms":   true,
		"dr":   true,
		"prof": true,
		"rev":  true,
		"hon":  true,
		"st":   true, // Saint

		// Academic
		"phd": true,
		"ba":  true,
		"bs":  true,
		"ma":  true,
		"mba": true,
		"jr":   true,
		"sr":   true,
		"esq":  true,

		// Common abbreviations
		"etc":  true,
		"ie":   true,
		"eg":   true,
		"vs":   true,
		"inc":  true,
		"ltd":  true,
		"corp": true,
		"co":   true,

		// Units
		"ft":  true,
		"in":  true,
		"lb":  true,
		"oz":  true,
		"km":  true,
		"cm":  true,
		"mm":  true,
		"kg":  true,
		"mg":  true,
		"ml":  true,

		// Time
		"am": true,
		"pm": true,
		"ad": true,
		"bc": true,
		"ce": true,

		// Other
		"no":   true, // Number
		"vol":  true,
		"ed":   true,
		"fig":  true,
		"ref":  true,
		"seq":  true,
		"ave":  true,
		"blvd": true,
		"rd":   true,
		"apt":  true,
		"dept": true,
		"min":  true,
		"max":  true,
		"approx": true,
	}

	return abbrevs
}

// IsAbbreviation implements DictionaryProvider.
func (d *EnglishDictionary) IsAbbreviation(word string) bool {
	// Normalize: remove all periods and lowercase
	normalized := strings.ToLower(strings.ReplaceAll(word, ".", ""))
	return d.abbreviations[normalized] || d.customWords[normalized]
}

// GetHyphenationPoints implements DictionaryProvider.
// Returns empty slice for now - hyphenation requires more complex rules.
//
// Future enhancement: Implement Liang's TeX hyphenation algorithm with
// language-specific pattern dictionaries.
func (d *EnglishDictionary) GetHyphenationPoints(word string) []int {
	// Placeholder - returns no hyphenation points
	// Hyphenation currently only works with manually-inserted soft hyphens (U+00AD)
	return nil
}

// IsCompoundWord implements DictionaryProvider.
func (d *EnglishDictionary) IsCompoundWord(word string) bool {
	// Common compound words that shouldn't be broken
	compounds := map[string]bool{
		"javascript": true,
		"typescript": true,
		"database":   true,
		"anybody":    true,
		"someone":    true,
		"everyone":   true,
	}

	return compounds[strings.ToLower(word)]
}

// AddAbbreviation adds a custom abbreviation to the dictionary.
func (d *EnglishDictionary) AddAbbreviation(abbrev string) {
	normalized := strings.ToLower(strings.TrimSuffix(abbrev, "."))
	d.customWords[normalized] = true
}

// AddAbbreviations adds multiple custom abbreviations.
func (d *EnglishDictionary) AddAbbreviations(abbrevs []string) {
	for _, abbrev := range abbrevs {
		d.AddAbbreviation(abbrev)
	}
}

// ═══════════════════════════════════════════════════════════════
//  Empty Dictionary (no filtering)
// ═══════════════════════════════════════════════════════════════

// EmptyDictionary provides no dictionary support (pure UAX #29).
type EmptyDictionary struct{}

// IsAbbreviation always returns false.
func (d *EmptyDictionary) IsAbbreviation(word string) bool {
	return false
}

// GetHyphenationPoints always returns nil.
func (d *EmptyDictionary) GetHyphenationPoints(word string) []int {
	return nil
}

// IsCompoundWord always returns false.
func (d *EmptyDictionary) IsCompoundWord(word string) bool {
	return false
}

// ═══════════════════════════════════════════════════════════════
//  Dictionary-Aware Sentence Segmentation
// ═══════════════════════════════════════════════════════════════

// SentencesWithDictionary splits text into sentences using dictionary support.
//
// This provides more accurate sentence boundary detection by:
// - Not breaking after known abbreviations (Dr., Mrs., etc.)
// - Handling language-specific rules
// - Supporting domain-specific terminology
//
// Example:
//
//	dict := text.NewEnglishDictionary()
//	sentences := txt.SentencesWithDictionary("Dr. Smith is here.", dict)
//	// Returns ["Dr. Smith is here."] instead of ["Dr. ", "Smith is here."]
func (t *Text) SentencesWithDictionary(text string, dict DictionaryProvider) []string {
	// Get raw UAX #29 sentence boundaries
	rawSentences := uax29.Sentences(text)

	if dict == nil {
		return rawSentences
	}

	// Post-process to merge sentences split on abbreviations
	var filtered []string
	var current string

	for i, sent := range rawSentences {
		current += sent

		// Check if this sentence ends with an abbreviation
		endsWithAbbrev := false
		trimmed := strings.TrimSpace(sent)
		if strings.HasSuffix(trimmed, ".") {
			// Get the last word before the period
			words := strings.Fields(trimmed)
			if len(words) > 0 {
				lastWord := words[len(words)-1]
				if dict.IsAbbreviation(lastWord) {
					endsWithAbbrev = true
				}
			}
		}

		// If doesn't end with abbreviation, or is last sentence, finalize it
		if !endsWithAbbrev || i == len(rawSentences)-1 {
			filtered = append(filtered, current)
			current = ""
		}
	}

	return filtered
}

// SentenceCountWithDictionary returns the number of sentences using dictionary support.
func (t *Text) SentenceCountWithDictionary(text string, dict DictionaryProvider) int {
	return len(t.SentencesWithDictionary(text, dict))
}

// ═══════════════════════════════════════════════════════════════
//  Dictionary-Aware Text Configuration
// ═══════════════════════════════════════════════════════════════

// TextConfig extends Text with dictionary support.
type TextConfig struct {
	*Text
	Dictionary DictionaryProvider
}

// NewTextWithDictionary creates a Text instance with dictionary support.
func NewTextWithDictionary(config Config, dict DictionaryProvider) *TextConfig {
	return &TextConfig{
		Text:       New(config),
		Dictionary: dict,
	}
}

// NewTerminalWithEnglishDictionary creates a terminal text handler with English dictionary.
func NewTerminalWithEnglishDictionary() *TextConfig {
	return &TextConfig{
		Text:       NewTerminal(),
		Dictionary: NewEnglishDictionary(),
	}
}

// Sentences uses the configured dictionary for sentence segmentation.
func (tc *TextConfig) Sentences(text string) []string {
	if tc.Dictionary == nil {
		return tc.Text.Sentences(text)
	}
	return tc.Text.SentencesWithDictionary(text, tc.Dictionary)
}

// SentenceCount returns sentence count using the configured dictionary.
func (tc *TextConfig) SentenceCount(text string) int {
	return len(tc.Sentences(text))
}
