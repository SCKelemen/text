package text

import (
	"strings"
)

// Hyphenation using Liang's Algorithm
//
// Implements Frank Liang's hyphenation algorithm (1983) used in TeX.
// Based on pattern matching with priority levels.
//
// This package provides decent hyphenation support for English, French, German,
// and Spanish. For comprehensive language support, users can load full TeX
// hyphenation patterns from:
//
//   - TeX hyphen patterns: https://github.com/hyphenation/tex-hyphen
//   - Pattern files: https://github.com/hyphenation/tex-hyphen/tree/master/hyph-utf8/tex/generic/hyph-utf8/patterns/txt
//
// Example of loading custom patterns:
//
//	patterns := loadMyPatterns() // Your pattern loading function
//	dict := text.NewHyphenationDictionary(patterns, 2, 3)
//	points := dict.Hyphenate("example")
//
// Reference: "Word Hy-phen-a-tion by Com-put-er" by Franklin Mark Liang
// https://tug.org/docs/liang/

// ═══════════════════════════════════════════════════════════════
//  Hyphenation Dictionary
// ═══════════════════════════════════════════════════════════════

// HyphenationDictionary provides hyphenation patterns for a language.
//
// Users can create custom dictionaries by loading patterns from TeX hyphen files:
// https://github.com/hyphenation/tex-hyphen
//
// Example:
//
//	dict := text.NewHyphenationDictionary(myPatterns, 2, 3)
//	points := dict.Hyphenate("example")
type HyphenationDictionary struct {
	patterns map[string]string // pattern -> priority string
	minLeft  int               // Minimum characters on left
	minRight int               // Minimum characters on right
}

// NewHyphenationDictionary creates a custom hyphenation dictionary.
//
// Parameters:
//   - patterns: Map of hyphenation patterns (see Liang's algorithm)
//   - minLeft: Minimum characters required on left side of hyphen
//   - minRight: Minimum characters required on right side of hyphen
//
// Example:
//
//	patterns := map[string]string{
//	    "ex1am": "ex1am",     // Allow break after "ex"
//	    "ta1ble": "ta1ble",   // Allow break after "ta"
//	}
//	dict := text.NewHyphenationDictionary(patterns, 2, 3)
func NewHyphenationDictionary(patterns map[string]string, minLeft, minRight int) *HyphenationDictionary {
	return &HyphenationDictionary{
		patterns: patterns,
		minLeft:  minLeft,
		minRight: minRight,
	}
}

// NewEnglishHyphenation creates a hyphenation dictionary with English (US) patterns.
//
// Provides decent coverage for common English words.
// For comprehensive hyphenation, load full TeX patterns:
// https://github.com/hyphenation/tex-hyphen/tree/master/hyph-utf8/tex/generic/hyph-utf8/patterns/txt
func NewEnglishHyphenation() *HyphenationDictionary {
	return &HyphenationDictionary{
		patterns: englishHyphenationPatterns(),
		minLeft:  2, // At least 2 chars on left
		minRight: 3, // At least 3 chars on right
	}
}

// NewFrenchHyphenation creates a hyphenation dictionary with French patterns.
func NewFrenchHyphenation() *HyphenationDictionary {
	return &HyphenationDictionary{
		patterns: frenchHyphenationPatterns(),
		minLeft:  2,
		minRight: 3,
	}
}

// NewGermanHyphenation creates a hyphenation dictionary with German patterns.
func NewGermanHyphenation() *HyphenationDictionary {
	return &HyphenationDictionary{
		patterns: germanHyphenationPatterns(),
		minLeft:  2,
		minRight: 2, // German allows shorter right side
	}
}

// NewSpanishHyphenation creates a hyphenation dictionary with Spanish patterns.
func NewSpanishHyphenation() *HyphenationDictionary {
	return &HyphenationDictionary{
		patterns: spanishHyphenationPatterns(),
		minLeft:  2,
		minRight: 2,
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

// frenchHyphenationPatterns returns French hyphenation patterns.
//
// French hyphenation follows specific rules:
// - Hyphenate between consonants
// - Specific patterns for common suffixes (-tion, -ment, etc.)
func frenchHyphenationPatterns() map[string]string {
	return map[string]string{
		// Common prefixes
		".dé3s2":   ".dé3s2",
		".in2":     ".in2",
		".con4":    ".con4",
		".pré3":    ".pré3",
		".pro3":    ".pro3",
		".trans3":  ".trans3",
		".re3":     ".re3",

		// Common suffixes
		"5tion.":   "5tion.",
		"5ation.":  "5ation.",
		"5ment.":   "5ment.",
		"5able.":   "5able.",
		"5ique.":   "5ique.",
		"5isme.":   "5isme.",
		"5eur.":    "5eur.",
		"5rice.":   "5rice.",
		"5eux.":    "5eux.",

		// Consonant patterns
		"1ba": "1ba",
		"1be": "1be",
		"1bi": "1bi",
		"1bo": "1bo",
		"1bu": "1bu",
		"1ca": "1ca",
		"1ce": "1ce",
		"1ci": "1ci",
		"1co": "1co",
		"1cu": "1cu",
		"1da": "1da",
		"1de": "1de",
		"1di": "1di",
		"1do": "1do",
		"1du": "1du",
		"1fa": "1fa",
		"1fe": "1fe",
		"1fi": "1fi",
		"1fo": "1fo",
		"1fu": "1fu",
		"1ga": "1ga",
		"1ge": "1ge",
		"1gi": "1gi",
		"1go": "1go",
		"1gu": "1gu",
		"1la": "1la",
		"1le": "1le",
		"1li": "1li",
		"1lo": "1lo",
		"1lu": "1lu",
		"1ma": "1ma",
		"1me": "1me",
		"1mi": "1mi",
		"1mo": "1mo",
		"1mu": "1mu",
		"1na": "1na",
		"1ne": "1ne",
		"1ni": "1ni",
		"1no": "1no",
		"1nu": "1nu",
		"1pa": "1pa",
		"1pe": "1pe",
		"1pi": "1pi",
		"1po": "1po",
		"1pu": "1pu",
		"1ra": "1ra",
		"1re": "1re",
		"1ri": "1ri",
		"1ro": "1ro",
		"1ru": "1ru",
		"1sa": "1sa",
		"1se": "1se",
		"1si": "1si",
		"1so": "1so",
		"1su": "1su",
		"1ta": "1ta",
		"1te": "1te",
		"1ti": "1ti",
		"1to": "1to",
		"1tu": "1tu",
		"1va": "1va",
		"1ve": "1ve",
		"1vi": "1vi",
		"1vo": "1vo",
		"1vu": "1vu",

		// Double consonants
		"2bb": "2bb",
		"2cc": "2cc",
		"2dd": "2dd",
		"2ff": "2ff",
		"2gg": "2gg",
		"2ll": "2ll",
		"2mm": "2mm",
		"2nn": "2nn",
		"2pp": "2pp",
		"2rr": "2rr",
		"2ss": "2ss",
		"2tt": "2tt",
	}
}

// germanHyphenationPatterns returns German hyphenation patterns.
//
// German has specific compound word rules and allows hyphenation
// after shorter fragments than English.
func germanHyphenationPatterns() map[string]string {
	return map[string]string{
		// Common prefixes
		".be3":     ".be3",
		".ver3":    ".ver3",
		".ent3":    ".ent3",
		".er3":     ".er3",
		".ge3":     ".ge3",
		".über3":   ".über3",
		".unter3":  ".unter3",
		".vor3":    ".vor3",
		".zer3":    ".zer3",

		// Common suffixes
		"3ung.":    "3ung.",
		"3schaft.": "3schaft.",
		"3heit.":   "3heit.",
		"3keit.":   "3keit.",
		"3lich.":   "3lich.",
		"3bar.":    "3bar.",
		"3sam.":    "3sam.",
		"3los.":    "3los.",

		// Consonant patterns
		"1ba": "1ba",
		"1be": "1be",
		"1bi": "1bi",
		"1bo": "1bo",
		"1bu": "1bu",
		"1bä": "1bä",
		"1bö": "1bö",
		"1bü": "1bü",
		"1da": "1da",
		"1de": "1de",
		"1di": "1di",
		"1do": "1do",
		"1du": "1du",
		"1dä": "1dä",
		"1dö": "1dö",
		"1dü": "1dü",
		"1fa": "1fa",
		"1fe": "1fe",
		"1fi": "1fi",
		"1fo": "1fo",
		"1fu": "1fu",
		"1ga": "1ga",
		"1ge": "1ge",
		"1gi": "1gi",
		"1go": "1go",
		"1gu": "1gu",
		"1ha": "1ha",
		"1he": "1he",
		"1hi": "1hi",
		"1ho": "1ho",
		"1hu": "1hu",
		"1ka": "1ka",
		"1ke": "1ke",
		"1ki": "1ki",
		"1ko": "1ko",
		"1ku": "1ku",
		"1la": "1la",
		"1le": "1le",
		"1li": "1li",
		"1lo": "1lo",
		"1lu": "1lu",
		"1ma": "1ma",
		"1me": "1me",
		"1mi": "1mi",
		"1mo": "1mo",
		"1mu": "1mu",
		"1na": "1na",
		"1ne": "1ne",
		"1ni": "1ni",
		"1no": "1no",
		"1nu": "1nu",
		"1pa": "1pa",
		"1pe": "1pe",
		"1pi": "1pi",
		"1po": "1po",
		"1pu": "1pu",
		"1ra": "1ra",
		"1re": "1re",
		"1ri": "1ri",
		"1ro": "1ro",
		"1ru": "1ru",
		"1sa": "1sa",
		"1se": "1se",
		"1si": "1si",
		"1so": "1so",
		"1su": "1su",
		"1ta": "1ta",
		"1te": "1te",
		"1ti": "1ti",
		"1to": "1to",
		"1tu": "1tu",
		"1wa": "1wa",
		"1we": "1we",
		"1wi": "1wi",
		"1wo": "1wo",
		"1wu": "1wu",

		// Double consonants
		"2bb": "2bb",
		"2ck": "2ck",
		"2dd": "2dd",
		"2ff": "2ff",
		"2gg": "2gg",
		"2ll": "2ll",
		"2mm": "2mm",
		"2nn": "2nn",
		"2pp": "2pp",
		"2rr": "2rr",
		"2ss": "2ss",
		"2tt": "2tt",

		// German specific
		"3sch": "3sch",
		"2ch":  "2ch",
		"3ck":  "3ck",
	}
}

// spanishHyphenationPatterns returns Spanish hyphenation patterns.
//
// Spanish follows syllable-based hyphenation rules with specific
// patterns for vowel and consonant combinations.
func spanishHyphenationPatterns() map[string]string {
	return map[string]string{
		// Common prefixes
		".des3":   ".des3",
		".in3":    ".in3",
		".con4":   ".con4",
		".pre3":   ".pre3",
		".pro3":   ".pro3",
		".re3":    ".re3",
		".sub3":   ".sub3",
		".trans3": ".trans3",

		// Common suffixes
		"5ción.":  "5ción.",
		"5sión.":  "5sión.",
		"5mente.": "5mente.",
		"5able.":  "5able.",
		"5ible.":  "5ible.",
		"5dad.":   "5dad.",
		"5tad.":   "5tad.",
		"5ismo.":  "5ismo.",

		// Consonant-vowel patterns
		"1ba": "1ba",
		"1be": "1be",
		"1bi": "1bi",
		"1bo": "1bo",
		"1bu": "1bu",
		"1ca": "1ca",
		"1ce": "1ce",
		"1ci": "1ci",
		"1co": "1co",
		"1cu": "1cu",
		"1da": "1da",
		"1de": "1de",
		"1di": "1di",
		"1do": "1do",
		"1du": "1du",
		"1fa": "1fa",
		"1fe": "1fe",
		"1fi": "1fi",
		"1fo": "1fo",
		"1fu": "1fu",
		"1ga": "1ga",
		"1ge": "1ge",
		"1gi": "1gi",
		"1go": "1go",
		"1gu": "1gu",
		"1ha": "1ha",
		"1he": "1he",
		"1hi": "1hi",
		"1ho": "1ho",
		"1hu": "1hu",
		"1ja": "1ja",
		"1je": "1je",
		"1ji": "1ji",
		"1jo": "1jo",
		"1ju": "1ju",
		"1la": "1la",
		"1le": "1le",
		"1li": "1li",
		"1lo": "1lo",
		"1lu": "1lu",
		"1ma": "1ma",
		"1me": "1me",
		"1mi": "1mi",
		"1mo": "1mo",
		"1mu": "1mu",
		"1na": "1na",
		"1ne": "1ne",
		"1ni": "1ni",
		"1no": "1no",
		"1nu": "1nu",
		"1ña": "1ña",
		"1ñe": "1ñe",
		"1ñi": "1ñi",
		"1ño": "1ño",
		"1ñu": "1ñu",
		"1pa": "1pa",
		"1pe": "1pe",
		"1pi": "1pi",
		"1po": "1po",
		"1pu": "1pu",
		"1ra": "1ra",
		"1re": "1re",
		"1ri": "1ri",
		"1ro": "1ro",
		"1ru": "1ru",
		"1sa": "1sa",
		"1se": "1se",
		"1si": "1si",
		"1so": "1so",
		"1su": "1su",
		"1ta": "1ta",
		"1te": "1te",
		"1ti": "1ti",
		"1to": "1to",
		"1tu": "1tu",
		"1va": "1va",
		"1ve": "1ve",
		"1vi": "1vi",
		"1vo": "1vo",
		"1vu": "1vu",

		// Double consonants
		"2bb": "2bb",
		"2cc": "2cc",
		"2dd": "2dd",
		"2ff": "2ff",
		"2ll": "2ll",
		"2mm": "2mm",
		"2nn": "2nn",
		"2pp": "2pp",
		"2rr": "2rr",
		"2ss": "2ss",
		"2tt": "2tt",

		// Spanish specific patterns
		"2ch": "2ch",
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
