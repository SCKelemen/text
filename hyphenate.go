package text

import (
	"strings"
)

// Hyphenation using Liang's Algorithm
//
// Implements Frank Liang's hyphenation algorithm (1983) used in TeX.
// Based on pattern matching with priority levels.
//
// Reference: "Word Hy-phen-a-tion by Com-put-er" by Franklin Mark Liang
// https://tug.org/docs/liang/

// ═══════════════════════════════════════════════════════════════
//  Hyphenation Dictionary
// ═══════════════════════════════════════════════════════════════

// HyphenationDictionary provides hyphenation patterns for a language.
type HyphenationDictionary struct {
	patterns map[string]string // pattern -> priority string
	minLeft  int               // Minimum characters on left
	minRight int               // Minimum characters on right
}

// NewEnglishHyphenation creates a hyphenation dictionary with English patterns.
//
// Uses a subset of TeX hyphenation patterns for demonstration.
// For production use, load full pattern files from:
// https://github.com/hyphenation/tex-hyphen
func NewEnglishHyphenation() *HyphenationDictionary {
	return &HyphenationDictionary{
		patterns: englishHyphenationPatterns(),
		minLeft:  2, // At least 2 chars on left
		minRight: 3, // At least 3 chars on right
	}
}

// englishHyphenationPatterns returns a subset of English hyphenation patterns.
//
// Pattern format: letters with numbers indicating break priority.
// Numbers are between positions. Odd numbers allow breaks, even prevent.
//
// Example: "ex1am" means allow break after "ex" (priority 1)
func englishHyphenationPatterns() map[string]string {
	return map[string]string{
		// Common prefixes
		".anti5":  ".anti5",
		".co4me":  ".co4me",
		".co4op":  ".co4op",
		".dis3":   ".dis3",
		".ex1":    ".ex1",
		".inter3": ".inter3",
		".multi3": ".multi3",
		".non1":   ".non1",
		".post3":  ".post3",
		".pre3":   ".pre3",
		".pro3":   ".pro3",
		".re3":    ".re3",
		".semi3":  ".semi3",
		".sub3":   ".sub3",
		".super5": ".super5",
		".trans3": ".trans3",
		".un1":    ".un1",
		".under3": ".under3",

		// Common suffixes
		"5able.":   "5able.",
		"5ible.":   "5ible.",
		"5ing.":    "5ing.",
		"5tion.":   "5tion.",
		"5sion.":   "5sion.",
		"5ness.":   "5ness.",
		"5ment.":   "5ment.",
		"5ful.":    "5ful.",
		"5less.":   "5less.",
		"5ous.":    "5ous.",
		"5ive.":    "5ive.",
		"3ence.":   "3ence.",
		"3ance.":   "3ance.",
		"3ity.":    "3ity.",
		"3ency.":   "3ency.",
		"3ancy.":   "3ancy.",
		"5er.":     "5er.",
		"5est.":    "5est.",
		"5ed.":     "5ed.",

		// Common patterns
		"1ba":     "1ba",
		"1be":     "1be",
		"1bi":     "1bi",
		"1bo":     "1bo",
		"1bu":     "1bu",
		"1ca":     "1ca",
		"1ce":     "1ce",
		"1ci":     "1ci",
		"1co":     "1co",
		"1cu":     "1cu",
		"1da":     "1da",
		"1de":     "1de",
		"1di":     "1di",
		"1do":     "1do",
		"1du":     "1du",
		"1ga":     "1ga",
		"1ge":     "1ge",
		"1gi":     "1gi",
		"1go":     "1go",
		"1gu":     "1gu",
		"1la":     "1la",
		"1le":     "1le",
		"1li":     "1li",
		"1lo":     "1lo",
		"1lu":     "1lu",
		"1ma":     "1ma",
		"1me":     "1me",
		"1mi":     "1mi",
		"1mo":     "1mo",
		"1mu":     "1mu",
		"1na":     "1na",
		"1ne":     "1ne",
		"1ni":     "1ni",
		"1no":     "1no",
		"1nu":     "1nu",
		"1pa":     "1pa",
		"1pe":     "1pe",
		"1pi":     "1pi",
		"1po":     "1po",
		"1pu":     "1pu",
		"1ra":     "1ra",
		"1re":     "1re",
		"1ri":     "1ri",
		"1ro":     "1ro",
		"1ru":     "1ru",
		"1sa":     "1sa",
		"1se":     "1se",
		"1si":     "1si",
		"1so":     "1so",
		"1su":     "1su",
		"1ta":     "1ta",
		"1te":     "1te",
		"1ti":     "1ti",
		"1to":     "1to",
		"1tu":     "1tu",
		"1va":     "1va",
		"1ve":     "1ve",
		"1vi":     "1vi",
		"1vo":     "1vo",
		"1vu":     "1vu",

		// Double consonants
		"2bb":  "2bb",
		"2cc":  "2cc",
		"2dd":  "2dd",
		"2ff":  "2ff",
		"2gg":  "2gg",
		"2ll":  "2ll",
		"2mm":  "2mm",
		"2nn":  "2nn",
		"2pp":  "2pp",
		"2rr":  "2rr",
		"2ss":  "2ss",
		"2tt":  "2tt",

		// Specific words
		"ta1ble":    "ta1ble",
		"rec1ord":   "rec1ord",
		"pre1sent":  "pre1sent",
		"ex1am":     "ex1am",
		"exam1ple":  "exam1ple",
		"con1test":  "con1test",
		"pro1ject":  "pro1ject",
		"in1for":    "in1for",
		"com1put":   "com1put",
		"al1go":     "al1go",
		"hyph1en":   "hyph1en",
		"pat1tern":  "pat1tern",
	}
}

// Hyphenate returns hyphenation points for a word using Liang's algorithm.
//
// Returns byte indices where hyphenation is allowed.
// Uses pattern matching with priority levels to determine break points.
//
// Example:
//
//	dict := text.NewEnglishHyphenation()
//	points := dict.Hyphenate("example")
//	// Returns []int{2, 4} for ex-am-ple
func (h *HyphenationDictionary) Hyphenate(word string) []int {
	if len(word) < h.minLeft+h.minRight {
		return nil // Too short to hyphenate
	}

	// Normalize word: lowercase and add delimiters
	normalized := "." + strings.ToLower(word) + "."

	// Initialize priority array (one value between each character)
	// Length is len(normalized) + 1 to account for positions
	priorities := make([]int, len(normalized)+1)

	// Apply all matching patterns
	for pattern := range h.patterns {
		h.applyPattern(normalized, pattern, priorities)
	}

	// Extract hyphenation points
	var points []int
	for i := h.minLeft; i < len(word)-h.minRight; i++ {
		// i+1 because priorities[0] is before first char
		// Odd priorities indicate allowed breaks
		if priorities[i+1]%2 == 1 {
			points = append(points, i)
		}
	}

	return points
}

// applyPattern applies a single hyphenation pattern to the word.
func (h *HyphenationDictionary) applyPattern(word, pattern string, priorities []int) {
	// Extract letters and numbers from pattern
	patternLetters := ""
	patternNumbers := make([]int, len(pattern)+1)
	pos := 0

	for i := 0; i < len(pattern); i++ {
		ch := pattern[i]
		if ch >= '0' && ch <= '9' {
			patternNumbers[pos] = int(ch - '0')
		} else {
			patternLetters += string(ch)
			pos++
		}
	}

	// Find all occurrences of the letter pattern in the word
	for i := 0; i <= len(word)-len(patternLetters); i++ {
		if word[i:i+len(patternLetters)] == patternLetters {
			// Apply priority numbers
			for j := 0; j <= len(patternLetters); j++ {
				if patternNumbers[j] > priorities[i+j] {
					priorities[i+j] = patternNumbers[j]
				}
			}
		}
	}
}

// HyphenateWithString returns the hyphenated word with hyphens inserted.
//
// Example:
//
//	dict := text.NewEnglishHyphenation()
//	result := dict.HyphenateWithString("example", "-")
//	// Returns "ex-am-ple"
func (h *HyphenationDictionary) HyphenateWithString(word, hyphen string) string {
	points := h.Hyphenate(word)
	if len(points) == 0 {
		return word
	}

	var result strings.Builder
	lastPos := 0

	for _, pos := range points {
		result.WriteString(word[lastPos:pos])
		result.WriteString(hyphen)
		lastPos = pos
	}
	result.WriteString(word[lastPos:])

	return result.String()
}

// ═══════════════════════════════════════════════════════════════
//  Integration with EnglishDictionary
// ═══════════════════════════════════════════════════════════════

// EnglishDictionaryWithHyphenation extends EnglishDictionary with hyphenation.
type EnglishDictionaryWithHyphenation struct {
	*EnglishDictionary
	hyphenation *HyphenationDictionary
}

// NewEnglishDictionaryWithHyphenation creates a full-featured English dictionary.
func NewEnglishDictionaryWithHyphenation() *EnglishDictionaryWithHyphenation {
	return &EnglishDictionaryWithHyphenation{
		EnglishDictionary: NewEnglishDictionary(),
		hyphenation:       NewEnglishHyphenation(),
	}
}

// GetHyphenationPoints implements DictionaryProvider with actual hyphenation.
func (d *EnglishDictionaryWithHyphenation) GetHyphenationPoints(word string) []int {
	if d.hyphenation == nil {
		return nil
	}
	return d.hyphenation.Hyphenate(word)
}
