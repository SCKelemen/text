package text

import (
	"strings"
	"unicode"

	"github.com/SCKelemen/unicode/uax14"
	"github.com/SCKelemen/unicode/uax29"
	"github.com/SCKelemen/units"
)

// CSS Text Module Level 3/4 Implementation
// Based on: https://www.w3.org/TR/css-text-3/ and https://www.w3.org/TR/css-text-4/

// ═══════════════════════════════════════════════════════════════
//  White Space Processing (CSS Text §4)
// ═══════════════════════════════════════════════════════════════

// WhiteSpace controls how white space is handled inside an element.
// Based on CSS Text Module Level 3 §4: https://www.w3.org/TR/css-text-3/#white-space-property
type WhiteSpace int

const (
	// WhiteSpaceNormal collapses white space sequences and wraps lines.
	// Newlines are treated as spaces.
	WhiteSpaceNormal WhiteSpace = iota

	// WhiteSpacePre preserves all white space and does not wrap.
	// Like HTML <pre> tag.
	WhiteSpacePre

	// WhiteSpaceNoWrap collapses white space but does not wrap lines.
	WhiteSpaceNoWrap

	// WhiteSpacePreWrap preserves white space and wraps lines.
	WhiteSpacePreWrap

	// WhiteSpacePreLine collapses white space sequences but preserves newlines.
	WhiteSpacePreLine

	// WhiteSpaceBreakSpaces preserves white space and allows breaking at spaces.
	WhiteSpaceBreakSpaces
)

// ═══════════════════════════════════════════════════════════════
//  Text Transformation (CSS Text §2)
// ═══════════════════════════════════════════════════════════════

// TextTransform controls case transformation of text.
// Based on CSS Text Module Level 3 §2: https://www.w3.org/TR/css-text-3/#text-transform-property
type TextTransform int

const (
	// TextTransformNone performs no transformation.
	TextTransformNone TextTransform = iota

	// TextTransformUppercase converts all characters to uppercase.
	TextTransformUppercase

	// TextTransformLowercase converts all characters to lowercase.
	TextTransformLowercase

	// TextTransformCapitalize capitalizes the first character of each word.
	TextTransformCapitalize

	// TextTransformFullWidth converts characters to their fullwidth forms.
	// Used in East Asian typography.
	TextTransformFullWidth

	// TextTransformFullSizeKana converts small kana to full-size equivalents.
	TextTransformFullSizeKana
)

// ═══════════════════════════════════════════════════════════════
//  Word Breaking (CSS Text §5)
// ═══════════════════════════════════════════════════════════════

// WordBreak controls word breaking rules for CJK and other scripts.
// Based on CSS Text Module Level 3 §5.2: https://www.w3.org/TR/css-text-3/#word-break-property
type WordBreak int

const (
	// WordBreakNormal uses default line break rules.
	WordBreakNormal WordBreak = iota

	// WordBreakBreakAll allows breaks between any characters for non-CJK scripts.
	WordBreakBreakAll

	// WordBreakKeepAll prevents breaks between CJK characters.
	// Only breaks at whitespace and punctuation.
	WordBreakKeepAll

	// WordBreakBreakWord is like normal, but allows breaking within words
	// if there are no acceptable break points in the line.
	WordBreakBreakWord
)

// ═══════════════════════════════════════════════════════════════
//  Line Breaking (CSS Text §5)
// ═══════════════════════════════════════════════════════════════

// LineBreak controls line breaking strictness for CJK text.
// Based on CSS Text Module Level 3 §5.3: https://www.w3.org/TR/css-text-3/#line-break-property
type LineBreak int

const (
	// LineBreakAuto uses the default line breaking rule.
	LineBreakAuto LineBreak = iota

	// LineBreakLoose uses the least restrictive line break rule.
	// Allows breaks at more positions.
	LineBreakLoose

	// LineBreakNormal uses the common line break rule.
	LineBreakNormal

	// LineBreakStrict uses the most restrictive line break rule.
	// Follows traditional typography rules strictly.
	LineBreakStrict

	// LineBreakAnywhere allows breaks at any character, even within words.
	LineBreakAnywhere
)

// ═══════════════════════════════════════════════════════════════
//  Overflow Wrapping (CSS Text §5)
// ═══════════════════════════════════════════════════════════════

// OverflowWrap controls whether to break within words to prevent overflow.
// Based on CSS Text Module Level 3 §5.5: https://www.w3.org/TR/css-text-3/#overflow-wrap-property
type OverflowWrap int

const (
	// OverflowWrapNormal breaks only at allowed break points.
	OverflowWrapNormal OverflowWrap = iota

	// OverflowWrapBreakWord breaks within words if necessary to prevent overflow.
	OverflowWrapBreakWord

	// OverflowWrapAnywhere breaks at any point if necessary.
	OverflowWrapAnywhere
)

// ═══════════════════════════════════════════════════════════════
//  Hyphens (CSS Text §4.3)
// ═══════════════════════════════════════════════════════════════

// Hyphens controls hyphenation behavior.
// Based on CSS Text Module Level 3 §4.3: https://www.w3.org/TR/css-text-3/#hyphenation
type Hyphens int

const (
	// HyphensNone disables all hyphenation.
	HyphensNone Hyphens = iota

	// HyphensManual only allows hyphenation at manually specified points (soft hyphens).
	HyphensManual

	// HyphensAuto allows automatic hyphenation using language-appropriate rules.
	HyphensAuto
)

// ═══════════════════════════════════════════════════════════════
//  Text Overflow (CSS UI §3.1)
// ═══════════════════════════════════════════════════════════════

// TextOverflow controls how overflowing inline content is signaled to users.
// Based on CSS Basic User Interface Module Level 4 §3.1:
// https://www.w3.org/TR/css-ui-4/#text-overflow
type TextOverflow int

const (
	// TextOverflowClip clips the text at the content edge (no ellipsis).
	TextOverflowClip TextOverflow = iota

	// TextOverflowEllipsis displays an ellipsis ('...') to represent clipped text.
	TextOverflowEllipsis

	// TextOverflowString displays a custom string to represent clipped text.
	// The string value is specified separately.
	TextOverflowString

	// TextOverflowFade fades out the end of the text (not widely supported).
	TextOverflowFade
)

// ═══════════════════════════════════════════════════════════════
//  CSS Text Style Configuration
// ═══════════════════════════════════════════════════════════════

// CSSTextStyle represents CSS text properties for text layout.
type CSSTextStyle struct {
	// White space handling
	WhiteSpace WhiteSpace

	// Text transformation
	TextTransform TextTransform

	// Word and line breaking
	WordBreak    WordBreak
	LineBreak    LineBreak
	OverflowWrap OverflowWrap
	Hyphens      Hyphens

	// Text overflow (CSS Overflow Module)
	// https://drafts.csswg.org/css-overflow/#text-overflow
	TextOverflow      TextOverflow // How to handle overflow (applies to both ends if TextOverflowEnd not set)
	TextOverflowEnd   TextOverflow // How to handle overflow at end (if different from start)
	TextOverflowClipString string  // Custom string for clip mode
	TextOverflowEllipsisString string // Custom ellipsis string (default: "...")

	// Spacing (using CSS length units)
	LetterSpacing units.Length // Additional spacing between characters
	WordSpacing   units.Length // Additional spacing between words

	// Indentation and alignment
	TextIndent    units.Length // First line indentation
	TextAlign     Alignment    // Horizontal alignment
	TextAlignLast Alignment    // Alignment of last line (justify becomes start)
	VerticalAlign Alignment    // Vertical alignment within line box
	HangingPunct  bool         // Allow punctuation to hang outside line box
}

// DefaultCSSTextStyle returns a CSSTextStyle with default values matching CSS defaults.
func DefaultCSSTextStyle() CSSTextStyle {
	return CSSTextStyle{
		WhiteSpace:             WhiteSpaceNormal,
		TextTransform:          TextTransformNone,
		WordBreak:              WordBreakNormal,
		LineBreak:              LineBreakAuto,
		OverflowWrap:           OverflowWrapNormal,
		Hyphens:                HyphensManual,
		TextOverflow:           TextOverflowClip,
		TextOverflowEllipsisString: "...",
		LetterSpacing:          units.Px(0),
		WordSpacing:            units.Px(0),
		TextIndent:             units.Px(0),
		TextAlign:              AlignLeft,
		TextAlignLast:          AlignLeft,
		VerticalAlign:          AlignLeft,
		HangingPunct:           false,
	}
}

// ═══════════════════════════════════════════════════════════════
//  White Space Processing Implementation
// ═══════════════════════════════════════════════════════════════

// ProcessWhiteSpace processes text according to CSS white-space property.
// Returns the processed text and whether line wrapping is allowed.
func (t *Text) ProcessWhiteSpace(text string, whiteSpace WhiteSpace) (processed string, allowWrap bool) {
	switch whiteSpace {
	case WhiteSpaceNormal:
		return t.collapseWhiteSpace(text, true), true

	case WhiteSpacePre:
		return text, false

	case WhiteSpaceNoWrap:
		return t.collapseWhiteSpace(text, true), false

	case WhiteSpacePreWrap:
		return text, true

	case WhiteSpacePreLine:
		return t.collapseWhiteSpace(text, false), true

	case WhiteSpaceBreakSpaces:
		return text, true

	default:
		return text, true
	}
}

// collapseWhiteSpace collapses sequences of white space into single spaces.
// If collapseNewlines is true, newlines are treated as spaces.
func (t *Text) collapseWhiteSpace(text string, collapseNewlines bool) string {
	var result strings.Builder
	result.Grow(len(text))

	inSpace := false
	for _, r := range text {
		isSpace := unicode.IsSpace(r)
		isNewline := r == '\n' || r == '\r'

		if isNewline && !collapseNewlines {
			result.WriteRune('\n')
			inSpace = false
			continue
		}

		if isSpace {
			if !inSpace {
				result.WriteRune(' ')
				inSpace = true
			}
			continue
		}

		result.WriteRune(r)
		inSpace = false
	}

	return strings.TrimSpace(result.String())
}

// ═══════════════════════════════════════════════════════════════
//  Text Transformation Implementation
// ═══════════════════════════════════════════════════════════════

// Transform applies text transformation according to CSS text-transform property.
func (t *Text) Transform(text string, transform TextTransform) string {
	switch transform {
	case TextTransformNone:
		return text

	case TextTransformUppercase:
		return strings.ToUpper(text)

	case TextTransformLowercase:
		return strings.ToLower(text)

	case TextTransformCapitalize:
		return t.capitalize(text)

	case TextTransformFullWidth:
		return t.toFullWidth(text)

	case TextTransformFullSizeKana:
		return t.toFullSizeKana(text)

	default:
		return text
	}
}

// capitalize capitalizes the first letter of each word.
func (t *Text) capitalize(text string) string {
	// Use UAX #29 word boundaries for proper capitalization
	words := uax29.Words(text)
	var result strings.Builder
	result.Grow(len(text))

	for _, word := range words {
		if len(word) == 0 {
			continue
		}

		// Check if this is a word (not punctuation or whitespace)
		firstRune := []rune(word)[0]
		if unicode.IsLetter(firstRune) {
			runes := []rune(word)
			runes[0] = unicode.ToUpper(runes[0])
			result.WriteString(string(runes))
		} else {
			result.WriteString(word)
		}
	}

	return result.String()
}

// toFullWidth converts ASCII characters to their fullwidth forms.
func (t *Text) toFullWidth(text string) string {
	var result strings.Builder
	result.Grow(len(text) * 2) // Fullwidth characters are larger in bytes

	for _, r := range text {
		// ASCII range: U+0021-U+007E -> Fullwidth: U+FF01-U+FF5E
		if r >= 0x21 && r <= 0x7E {
			result.WriteRune(r - 0x21 + 0xFF01)
		} else if r == 0x20 { // Space -> Fullwidth space
			result.WriteRune(0x3000)
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// toFullSizeKana converts small kana to full-size equivalents.
func (t *Text) toFullSizeKana(text string) string {
	// Map of small kana to full-size kana
	smallToFull := map[rune]rune{
		'ぁ': 'あ', 'ぃ': 'い', 'ぅ': 'う', 'ぇ': 'え', 'ぉ': 'お',
		'ゃ': 'や', 'ゅ': 'ゆ', 'ょ': 'よ', 'ゎ': 'わ',
		'ァ': 'ア', 'ィ': 'イ', 'ゥ': 'ウ', 'ェ': 'エ', 'ォ': 'オ',
		'ャ': 'ヤ', 'ュ': 'ユ', 'ョ': 'ヨ', 'ヮ': 'ワ',
		'っ': 'つ', 'ッ': 'ツ',
	}

	var result strings.Builder
	result.Grow(len(text))

	for _, r := range text {
		if full, ok := smallToFull[r]; ok {
			result.WriteRune(full)
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// ═══════════════════════════════════════════════════════════════
//  Word Boundary Support
// ═══════════════════════════════════════════════════════════════

// Words splits text into words using UAX #29 word boundaries.
// Returns a slice of word segments (includes spaces and punctuation).
func (t *Text) Words(text string) []string {
	return uax29.Words(text)
}

// WordCount returns the number of words in the text.
func (t *Text) WordCount(text string) int {
	words := uax29.Words(text)
	count := 0
	for _, word := range words {
		// Only count actual words (not whitespace or punctuation)
		if len(word) > 0 {
			firstRune := []rune(word)[0]
			if unicode.IsLetter(firstRune) || unicode.IsNumber(firstRune) {
				count++
			}
		}
	}
	return count
}

// ═══════════════════════════════════════════════════════════════
//  Sentence Boundary Support
// ═══════════════════════════════════════════════════════════════

// Sentences splits text into sentences using UAX #29 sentence boundaries.
func (t *Text) Sentences(text string) []string {
	return uax29.Sentences(text)
}

// SentenceCount returns the number of sentences in the text.
func (t *Text) SentenceCount(text string) int {
	return len(uax29.Sentences(text))
}

// ═══════════════════════════════════════════════════════════════
//  Advanced Wrapping with CSS Text Properties
// ═══════════════════════════════════════════════════════════════

// CSSWrapOptions extends WrapOptions with CSS text properties.
type CSSWrapOptions struct {
	MaxWidth units.Length
	Style    CSSTextStyle
}

// WrapCSS wraps text according to CSS text properties.
// This is a more sophisticated version of Wrap that handles white-space,
// word-break, line-break, and other CSS properties.
func (t *Text) WrapCSS(text string, opts CSSWrapOptions) []Line {
	// Process white space first
	processed, allowWrap := t.ProcessWhiteSpace(text, opts.Style.WhiteSpace)

	if !allowWrap {
		// No wrapping allowed
		return []Line{{
			Content: processed,
			Width:   t.Width(processed),
			Start:   0,
			End:     len([]rune(processed)),
		}}
	}

	// Apply text transformation
	processed = t.Transform(processed, opts.Style.TextTransform)

	// Convert CSS properties to UAX #14 line breaking options
	hyphenMode := uax14.HyphensManual
	switch opts.Style.Hyphens {
	case HyphensNone:
		hyphenMode = uax14.HyphensNone
	case HyphensManual:
		hyphenMode = uax14.HyphensManual
	case HyphensAuto:
		hyphenMode = uax14.HyphensAuto
	}

	// Find line break opportunities using UAX #14
	breakPoints := uax14.FindLineBreakOpportunities(processed, hyphenMode)

	// Build lines using break opportunities
	return t.buildLinesFromBreakPoints(processed, breakPoints, opts)
}

// buildLinesFromBreakPoints creates lines from UAX #14 break points.
func (t *Text) buildLinesFromBreakPoints(text string, breakPoints []int, opts CSSWrapOptions) []Line {
	if len(breakPoints) == 0 {
		return []Line{{
			Content: text,
			Width:   t.Width(text),
			Start:   0,
			End:     len([]rune(text)),
		}}
	}

	var lines []Line
	lineStart := 0
	maxWidth := opts.MaxWidth.Raw()

	for i := 1; i < len(breakPoints); i++ {
		segment := text[breakPoints[i-1]:breakPoints[i]]
		segmentWidth := t.Width(segment)

		// Apply letter spacing
		if !opts.Style.LetterSpacing.IsZero() {
			graphemes := t.Graphemes(segment)
			segmentWidth += float64(len(graphemes)-1) * opts.Style.LetterSpacing.Raw()
		}

		// Apply word spacing
		if !opts.Style.WordSpacing.IsZero() {
			spaceCount := strings.Count(segment, " ")
			segmentWidth += float64(spaceCount) * opts.Style.WordSpacing.Raw()
		}

		// Check if this segment fits
		if segmentWidth <= maxWidth || len(lines) == 0 {
			// Create a line
			lines = append(lines, Line{
				Content: segment,
				Width:   segmentWidth,
				Start:   lineStart,
				End:     lineStart + len([]rune(segment)),
			})
			lineStart += len([]rune(segment))
		}
	}

	return lines
}

// ═══════════════════════════════════════════════════════════════
//  Text Overflow Implementation (CSS Overflow Module)
// ═══════════════════════════════════════════════════════════════

// ApplyTextOverflow applies CSS text-overflow property to text that exceeds maxWidth.
//
// Based on CSS Overflow Module Level 3 §3:
// https://drafts.csswg.org/css-overflow/#text-overflow
//
// Example:
//
//	txt := text.NewTerminal()
//	result := txt.ApplyTextOverflow("Very long text here", 15, text.CSSTextStyle{
//	    TextOverflow: text.TextOverflowEllipsis,
//	    TextOverflowEllipsisString: "…",
//	})
//	// Returns: "Very long t…"
func (t *Text) ApplyTextOverflow(text string, maxWidth float64, style CSSTextStyle) string {
	currentWidth := t.Width(text)

	// No overflow, return as-is
	if currentWidth <= maxWidth {
		return text
	}

	// Determine overflow mode (use TextOverflow for both ends if TextOverflowEnd not set)
	overflowMode := style.TextOverflow
	if style.TextOverflowEnd != TextOverflowClip {
		// If end mode is set differently, use it for end truncation
		overflowMode = style.TextOverflowEnd
	}

	switch overflowMode {
	case TextOverflowClip:
		// Just clip at maxWidth (no indicator)
		return t.clipAtWidth(text, maxWidth)

	case TextOverflowEllipsis:
		// Add ellipsis
		ellipsis := style.TextOverflowEllipsisString
		if ellipsis == "" {
			ellipsis = "..."
		}
		return t.Truncate(text, TruncateOptions{
			MaxWidth: maxWidth,
			Strategy: TruncateEnd,
			Ellipsis: ellipsis,
		})

	case TextOverflowString:
		// Use custom overflow string
		clipString := style.TextOverflowClipString
		if clipString == "" {
			clipString = "..."
		}
		return t.Truncate(text, TruncateOptions{
			MaxWidth: maxWidth,
			Strategy: TruncateEnd,
			Ellipsis: clipString,
		})

	case TextOverflowFade:
		// Fade is not widely supported, fallback to ellipsis
		ellipsis := style.TextOverflowEllipsisString
		if ellipsis == "" {
			ellipsis = "..."
		}
		return t.Truncate(text, TruncateOptions{
			MaxWidth: maxWidth,
			Strategy: TruncateEnd,
			Ellipsis: ellipsis,
		})

	default:
		return t.clipAtWidth(text, maxWidth)
	}
}

// clipAtWidth clips text at the exact width without any indicator.
func (t *Text) clipAtWidth(text string, maxWidth float64) string {
	graphemes := t.Graphemes(text)
	result := ""
	width := 0.0

	for _, g := range graphemes {
		gWidth := t.Width(g)
		if width+gWidth > maxWidth {
			break
		}
		result += g
		width += gWidth
	}

	return result
}

// ═══════════════════════════════════════════════════════════════
//  Text Alignment with CSS Support
// ═══════════════════════════════════════════════════════════════

// AlignLines aligns multiple lines according to CSS text-align and text-align-last properties.
//
// Based on CSS Text Module Level 3:
// - text-align: https://www.w3.org/TR/css-text-3/#text-align-property
// - text-align-last: https://www.w3.org/TR/css-text-3/#text-align-last-property
//
// Example:
//
//	lines := []Line{
//	    {Content: "First line", Width: 10},
//	    {Content: "Second line", Width: 11},
//	    {Content: "Last", Width: 4},
//	}
//	aligned := txt.AlignLines(lines, 20.0, text.CSSTextStyle{
//	    TextAlign:     text.AlignJustify,
//	    TextAlignLast: text.AlignLeft,
//	})
func (t *Text) AlignLines(lines []Line, width float64, style CSSTextStyle) []Line {
	if len(lines) == 0 {
		return lines
	}

	result := make([]Line, len(lines))
	copy(result, lines)

	for i := range result {
		isLastLine := (i == len(result)-1)

		// Determine which alignment to use
		align := style.TextAlign
		if isLastLine && style.TextAlignLast != AlignLeft {
			// Use text-align-last for the last line
			// If text-align-last is not set (AlignLeft is default), use text-align
			align = style.TextAlignLast

			// Special handling: if text-align is justify but last line uses default,
			// the last line should use start alignment (left in LTR)
			if style.TextAlign == AlignJustify && style.TextAlignLast == AlignLeft {
				align = AlignLeft
			}
		}

		// Apply alignment
		result[i].Content = t.Align(result[i].Content, width, align)
		result[i].Width = width
	}

	return result
}
