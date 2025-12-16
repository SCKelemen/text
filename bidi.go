package text

import (
	"strings"

	"github.com/SCKelemen/unicode/uax9"
)

// Bidirectional Text Support (UAX #9)
//
// This file provides additional bidirectional text utilities beyond the
// basic Reorder() and DetectDirection() methods in text.go.
//
// The bidirectional algorithm handles:
// - Mixed Latin and Arabic/Hebrew text
// - Proper number handling in RTL contexts
// - Bracket and punctuation mirroring
// - Nested embeddings and isolates
//
// References:
//   - UAX #9: https://www.unicode.org/reports/tr9/
//   - Bidirectional Algorithm: https://www.unicode.org/reports/tr9/#Basic_Display_Algorithm
//   - W3C CSS Writing Modes: https://www.w3.org/TR/css-writing-modes-3/#direction
//   - MDN direction: https://developer.mozilla.org/en-US/docs/Web/CSS/direction

// ═══════════════════════════════════════════════════════════════
//  Helper Functions for Bidirectional Text
// ═══════════════════════════════════════════════════════════════

// toUAX9Direction converts our Direction type to uax9.Direction.
func toUAX9Direction(d Direction) uax9.Direction {
	switch d {
	case DirectionLTR:
		return uax9.DirectionLTR
	case DirectionRTL:
		return uax9.DirectionRTL
	case DirectionAuto:
		return uax9.DirectionAuto
	default:
		return uax9.DirectionLTR
	}
}

// fromUAX9Direction converts uax9.Direction to our Direction type.
func fromUAX9Direction(d uax9.Direction) Direction {
	switch d {
	case uax9.DirectionLTR:
		return DirectionLTR
	case uax9.DirectionRTL:
		return DirectionRTL
	case uax9.DirectionAuto:
		return DirectionAuto
	default:
		return DirectionLTR
	}
}

// ═══════════════════════════════════════════════════════════════
//  Paragraph-Level Reordering
// ═══════════════════════════════════════════════════════════════

// ReorderParagraph reorders a paragraph, handling line breaks properly.
//
// This function:
// 1. Splits text by line breaks
// 2. Detects direction for each paragraph
// 3. Reorders each line independently
// 4. Rejoins with line breaks
//
// Example:
//
//	txt := text.NewTerminal()
//
//	input := "Hello world\nمرحبا العالم\nMixed text"
//	result := txt.ReorderParagraph(input, text.DirectionAuto)
//	// Each line reordered according to its detected direction
func (t *Text) ReorderParagraph(text string, baseDirection Direction) string {
	if text == "" {
		return text
	}

	// Split by line breaks (paragraph separators)
	lines := strings.Split(text, "\n")
	reordered := make([]string, len(lines))

	dir := toUAX9Direction(baseDirection)

	for i, line := range lines {
		// Each line is treated as a separate paragraph
		reordered[i] = uax9.Reorder(line, dir)
	}

	return strings.Join(reordered, "\n")
}

// ═══════════════════════════════════════════════════════════════
//  Integration with Line Breaking
// ═══════════════════════════════════════════════════════════════

// ReorderLine reorders a single line after wrapping.
//
// This is used internally by wrapping functions to ensure each wrapped line
// is properly reordered for display. Layout engines should call this after
// computing line breaks.
//
// Example:
//
//	txt := text.NewTerminal()
//
//	line := text.Line{
//	    Content: "Hello مرحبا world",
//	    Width:   15.0,
//	    Start:   0,
//	    End:     17,
//	}
//
//	// Reorder the line content
//	line.Content = txt.ReorderLine(line.Content, text.DirectionAuto)
func (t *Text) ReorderLine(text string, direction Direction) string {
	dir := toUAX9Direction(direction)
	return uax9.Reorder(text, dir)
}

// ═══════════════════════════════════════════════════════════════
//  Bracket Mirroring
// ═══════════════════════════════════════════════════════════════

// mirrorBrackets maps opening brackets to their mirrored closing equivalents.
var mirrorBrackets = map[rune]rune{
	'(': ')',
	')': '(',
	'[': ']',
	']': '[',
	'{': '}',
	'}': '{',
	'<': '>',
	'>': '<',
	'«': '»',
	'»': '«',
	'‹': '›',
	'›': '‹',
	'⟨': '⟩',
	'⟩': '⟨',
	'⟪': '⟫',
	'⟫': '⟪',
	'⟬': '⟭',
	'⟭': '⟬',
	'⟮': '⟯',
	'⟯': '⟮',
	// CJK brackets
	'「': '」',
	'」': '「',
	'『': '』',
	'』': '『',
	'〈': '〉',
	'〉': '〈',
	'《': '》',
	'》': '《',
	'【': '】',
	'】': '【',
	'〔': '〕',
	'〕': '〔',
	'〖': '〗',
	'〗': '〖',
	'〘': '〙',
	'〙': '〘',
	'〚': '〛',
	'〛': '〚',
}

// MirrorBrackets mirrors brackets for RTL display.
//
// In RTL text, opening brackets like '(' should become ')' and vice versa.
// This function is typically called automatically by the reordering algorithm,
// but is exposed for cases where manual bracket handling is needed.
//
// Note: The UAX #9 algorithm handles this automatically in most cases.
// This function is provided for special situations or debugging.
//
// Example:
//
//	txt := text.NewTerminal()
//
//	result := txt.MirrorBrackets("(hello)")
//	// result: ")hello(" (brackets mirrored)
func (t *Text) MirrorBrackets(text string) string {
	runes := []rune(text)
	for i, r := range runes {
		if mirrored, ok := mirrorBrackets[r]; ok {
			runes[i] = mirrored
		}
	}
	return string(runes)
}

// ═══════════════════════════════════════════════════════════════
//  Bidirectional Character Type Query
// ═══════════════════════════════════════════════════════════════

// BidiClass represents the bidirectional character type.
type BidiClass = uax9.BidiClass

// Bidirectional character classes (re-exported from uax9)
const (
	// Strong types
	ClassL  = uax9.ClassL  // Left-to-Right
	ClassR  = uax9.ClassR  // Right-to-Left
	ClassAL = uax9.ClassAL // Right-to-Left Arabic

	// Weak types
	ClassEN  = uax9.ClassEN  // European Number
	ClassES  = uax9.ClassES  // European Number Separator
	ClassET  = uax9.ClassET  // European Number Terminator
	ClassAN  = uax9.ClassAN  // Arabic Number
	ClassCS  = uax9.ClassCS  // Common Number Separator
	ClassNSM = uax9.ClassNSM // Nonspacing Mark
	ClassBN  = uax9.ClassBN  // Boundary Neutral

	// Neutral types
	ClassB  = uax9.ClassB  // Paragraph Separator
	ClassS  = uax9.ClassS  // Segment Separator
	ClassWS = uax9.ClassWS // Whitespace
	ClassON = uax9.ClassON // Other Neutrals

	// Explicit formatting types
	ClassLRE = uax9.ClassLRE // Left-to-Right Embedding
	ClassLRO = uax9.ClassLRO // Left-to-Right Override
	ClassRLE = uax9.ClassRLE // Right-to-Left Embedding
	ClassRLO = uax9.ClassRLO // Right-to-Left Override
	ClassPDF = uax9.ClassPDF // Pop Directional Format
	ClassLRI = uax9.ClassLRI // Left-to-Right Isolate
	ClassRLI = uax9.ClassRLI // Right-to-Left Isolate
	ClassFSI = uax9.ClassFSI // First Strong Isolate
	ClassPDI = uax9.ClassPDI // Pop Directional Isolate
)

// GetBidiClass returns the bidirectional character type for a rune.
//
// This is useful for understanding how a character behaves in bidirectional
// text or for implementing custom bidirectional algorithms.
//
// Example:
//
//	txt := text.NewTerminal()
//
//	class := txt.GetBidiClass('a')      // ClassL (Left-to-Right)
//	class := txt.GetBidiClass('א')      // ClassR (Right-to-Left Hebrew)
//	class := txt.GetBidiClass('ا')      // ClassAL (Right-to-Left Arabic)
//	class := txt.GetBidiClass('5')      // ClassEN (European Number)
func (t *Text) GetBidiClass(r rune) BidiClass {
	return uax9.GetBidiClass(r)
}
