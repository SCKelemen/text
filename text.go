// Package text provides Unicode-aware text measurement, wrapping, truncation,
// and alignment operations for terminal UIs and text layout engines.
//
// This package coordinates multiple Unicode standards (UAX #11, #14, #29, #9,
// UTS #51) to provide practical text operations that are correct for CJK
// characters, emoji, combining marks, and bidirectional text.
//
// Based on:
//   - UAX #11: East Asian Width (https://www.unicode.org/reports/tr11/)
//   - UAX #14: Line Breaking (https://www.unicode.org/reports/tr14/)
//   - UAX #29: Text Segmentation (https://www.unicode.org/reports/tr29/)
//   - UAX #9: Bidirectional Algorithm (https://www.unicode.org/reports/tr9/)
//   - UTS #51: Unicode Emoji (https://www.unicode.org/reports/tr51/)
//
// # Units
//
// All width measurements (MaxWidth, Width, Line.Width) are in "abstract units"
// determined by the MeasureFunc configuration:
//   - For terminals: character cells (1.0 per ASCII, 2.0 per CJK)
//   - For canvas: pixels (varies by font)
//
// The package doesn't care about the unit semantics - it just accumulates
// whatever your MeasureFunc returns.
//
// # Quick Start
//
//	import "github.com/SCKelemen/text"
//
//	// Create terminal text handler
//	txt := text.NewTerminal()
//
//	// Measure width (accounts for wide characters, emoji)
//	width := txt.Width("Hello 世界")  // 9.0 cells
//
//	// Wrap text to fit width
//	lines := txt.Wrap("Long text here...", text.WrapOptions{
//	    MaxWidth: 40,  // 40 character cells
//	})
//
//	// Truncate with ellipsis
//	short := txt.Truncate("Very long text", text.TruncateOptions{
//	    MaxWidth: 10,  // 10 character cells
//	    Strategy: text.TruncateEnd,
//	})
//
//	// Align text
//	aligned := txt.Align("Hello", 20, text.AlignCenter)  // 20 cells total
//
// # Configuration
//
// The package is renderer-agnostic through the MeasureFunc configuration:
//
//	// Terminal: measure in cells
//	txt := text.New(text.Config{
//	    MeasureFunc: text.TerminalMeasure,
//	})
//
//	// Canvas: measure in pixels (future)
//	txt := text.New(text.Config{
//	    MeasureFunc: func(r rune) float64 {
//	        return fontFace.GlyphAdvance(r)
//	    },
//	})
package text

import (
	"strings"

	"github.com/SCKelemen/unicode/uax11"
	"github.com/SCKelemen/unicode/uax14"
	"github.com/SCKelemen/unicode/uax29"
	"github.com/SCKelemen/unicode/uax9"
	"github.com/SCKelemen/unicode/uts51"
)

// ═══════════════════════════════════════════════════════════════
//  Configuration
// ═══════════════════════════════════════════════════════════════

// Config configures text measurement and layout behavior.
type Config struct {
	// MeasureFunc measures the width of a single rune in abstract units.
	// For terminals: returns 1 or 2 (cells)
	// For canvas: returns pixel width from font metrics
	MeasureFunc MeasureFunc

	// AmbiguousAsWide determines UAX #11 context for ambiguous width characters.
	// Set to true for East Asian contexts (Chinese, Japanese, Korean).
	// Set to false (default) for non-East Asian contexts.
	AmbiguousAsWide bool

	// HyphenationMode specifies UAX #14 line breaking preferences.
	HyphenationMode uax14.Hyphens

	// BaseDirection specifies the default paragraph direction for UAX #9.
	BaseDirection uax9.Direction
}

// MeasureFunc measures the width of a single rune in abstract units.
//
// For terminal rendering, it should return 1 or 2 (character cells).
// For pixel-based rendering, it should return the actual pixel width.
type MeasureFunc func(r rune) float64

// Text provides high-level Unicode-aware text operations.
type Text struct {
	config Config
}

// New creates a new Text instance with the given configuration.
func New(config Config) *Text {
	// Set defaults
	if config.MeasureFunc == nil {
		config.MeasureFunc = TerminalMeasure
	}
	if config.HyphenationMode == 0 {
		config.HyphenationMode = uax14.HyphensManual
	}
	if config.BaseDirection == 0 {
		config.BaseDirection = uax9.DirectionLTR
	}

	return &Text{config: config}
}

// NewTerminal creates a Text instance configured for terminal rendering.
//
// Uses:
//   - TerminalMeasure for width calculation (1 or 2 cells)
//   - ContextNarrow for UAX #11 ambiguous characters
//   - Manual hyphenation for UAX #14
//   - LTR base direction for UAX #9
func NewTerminal() *Text {
	return New(Config{
		MeasureFunc:     TerminalMeasure,
		AmbiguousAsWide: false,
		HyphenationMode: uax14.HyphensManual,
		BaseDirection:   uax9.DirectionLTR,
	})
}

// NewTerminalEastAsian creates a Text instance configured for East Asian terminals.
//
// Same as NewTerminal but treats ambiguous characters as wide (2 cells).
func NewTerminalEastAsian() *Text {
	return New(Config{
		MeasureFunc:     TerminalMeasureEastAsian,
		AmbiguousAsWide: true,
		HyphenationMode: uax14.HyphensManual,
		BaseDirection:   uax9.DirectionLTR,
	})
}

// ═══════════════════════════════════════════════════════════════
//  Measurement
// ═══════════════════════════════════════════════════════════════

// Width measures the display width of text in abstract units.
//
// Units are determined by the MeasureFunc configuration:
//   - For terminals: character cells (1.0 per ASCII char, 2.0 per CJK char)
//   - For canvas: pixels (12.5 per 'a', 24.3 per 'W', etc.)
//
// This correctly handles:
//   - CJK characters (2 cells/units wide)
//   - Emoji (2 cells/units wide)
//   - Combining marks (0 width)
//   - Zero-width joiners (0 width)
//
// Example:
//
//	txt := text.NewTerminal()
//	width := txt.Width("Hello")     // 5.0 cells
//	width = txt.Width("Hello 世界")  // 9.0 cells (5 + 1 space + 2 + 2)
//	width = txt.Width("👋🏻")        // 2.0 cells (emoji + skin tone modifier)
func (t *Text) Width(s string) float64 {
	width := 0.0
	for _, r := range s {
		width += t.config.MeasureFunc(r)
	}
	return width
}

// WidthRange measures the display width of a substring by rune indices.
//
// Example:
//
//	txt := text.NewTerminal()
//	text := "Hello 世界"
//	width := txt.WidthRange(text, 0, 5)  // 5.0 (just "Hello")
//	width = txt.WidthRange(text, 6, 8)   // 4.0 (just "世界")
func (t *Text) WidthRange(s string, start, end int) float64 {
	width := 0.0
	i := 0
	for _, r := range s {
		if i >= start && i < end {
			width += t.config.MeasureFunc(r)
		}
		i++
		if i >= end {
			break
		}
	}
	return width
}

// ═══════════════════════════════════════════════════════════════
//  Pre-configured Measure Functions
// ═══════════════════════════════════════════════════════════════

// TerminalMeasure measures characters in terminal cells.
//
// Returns:
//   - 2 for wide characters (CJK ideographs, fullwidth, emoji)
//   - 1 for narrow characters (ASCII, halfwidth)
//   - 0 for zero-width characters (combining marks, ZWJ, variation selectors, emoji modifiers)
//
// Uses UAX #11 with ContextNarrow (ambiguous characters treated as narrow).
// UTS #51 takes precedence for emoji characters.
func TerminalMeasure(r rune) float64 {
	// Check if this is an emoji character (has emoji properties)
	// UTS #51 takes precedence over UAX #11 for emoji
	if uts51.IsEmoji(r) || uts51.IsEmojiComponent(r) {
		return float64(uts51.EmojiWidth(r))
	}

	// Use UAX #11 for non-emoji characters
	return float64(uax11.CharWidth(r, uax11.ContextNarrow))
}

// TerminalMeasureEastAsian measures characters in terminal cells with East Asian context.
//
// Same as TerminalMeasure but treats ambiguous characters as wide (2 cells).
// Use this for terminals with East Asian locales (Chinese, Japanese, Korean).
func TerminalMeasureEastAsian(r rune) float64 {
	// Check if this is an emoji character (has emoji properties)
	// UTS #51 takes precedence over UAX #11 for emoji
	if uts51.IsEmoji(r) || uts51.IsEmojiComponent(r) {
		return float64(uts51.EmojiWidth(r))
	}

	// Use UAX #11 with East Asian context for non-emoji characters
	return float64(uax11.CharWidth(r, uax11.ContextEastAsian))
}

// ═══════════════════════════════════════════════════════════════
//  Wrapping
// ═══════════════════════════════════════════════════════════════

// WrapOptions configures text wrapping behavior.
type WrapOptions struct {
	// MaxWidth is the maximum width for wrapped lines in abstract units.
	// For terminals: character cells (e.g., 40 means 40 columns)
	// For canvas: pixels (e.g., 400.0 means 400 pixels)
	// Units are determined by the MeasureFunc configuration.
	MaxWidth float64

	// BreakWords allows breaking in the middle of words if necessary.
	// If false, only breaks at UAX #14 line break opportunities.
	BreakWords bool

	// PreserveNewlines keeps existing newline characters as line breaks.
	PreserveNewlines bool
}

// Line represents a wrapped line of text.
type Line struct {
	// Content is the text content of the line.
	Content string

	// Width is the display width of the line in abstract units.
	// For terminals: character cells
	// For canvas: pixels
	Width float64

	// Start is the rune index in the original text where this line starts.
	Start int

	// End is the rune index in the original text where this line ends.
	End int
}

// Wrap breaks text into lines that fit within maxWidth.
//
// Uses UAX #14 for proper line break opportunities and UAX #29 to avoid
// breaking within grapheme clusters (emoji, combining marks, etc.).
//
// Example:
//
//	txt := text.NewTerminal()
//	lines := txt.Wrap("Hello 世界! This is a test.", text.WrapOptions{
//	    MaxWidth: 15,
//	})
//	for _, line := range lines {
//	    fmt.Println(line.Content)
//	}
//	// Output:
//	// Hello 世界!
//	// This is a test.
func (t *Text) Wrap(text string, opts WrapOptions) []Line {
	if opts.MaxWidth <= 0 {
		return []Line{{Content: text, Width: t.Width(text), Start: 0, End: len([]rune(text))}}
	}

	// Use UAX #29 to get grapheme clusters (don't break emoji!)
	graphemes := uax29.Graphemes(text)

	var lines []Line
	currentLine := ""
	currentWidth := 0.0
	lineStart := 0

	for i, g := range graphemes {
		gWidth := t.Width(g)

		// Check if we need to break
		if currentWidth+gWidth > opts.MaxWidth && currentWidth > 0 {
			// Add current line
			lines = append(lines, Line{
				Content: currentLine,
				Width:   currentWidth,
				Start:   lineStart,
				End:     lineStart + len([]rune(currentLine)),
			})

			// Start new line
			currentLine = g
			currentWidth = gWidth
			lineStart = i
		} else {
			currentLine += g
			currentWidth += gWidth
		}
	}

	// Add final line
	if currentLine != "" {
		lines = append(lines, Line{
			Content: currentLine,
			Width:   currentWidth,
			Start:   lineStart,
			End:     lineStart + len([]rune(currentLine)),
		})
	}

	return lines
}

// ═══════════════════════════════════════════════════════════════
//  Truncation
// ═══════════════════════════════════════════════════════════════

// TruncateOptions configures truncation behavior.
type TruncateOptions struct {
	// MaxWidth is the maximum width for the truncated text in abstract units.
	// For terminals: character cells (e.g., 20 means 20 columns)
	// For canvas: pixels (e.g., 200.0 means 200 pixels)
	// Units are determined by the MeasureFunc configuration.
	MaxWidth float64

	// Ellipsis is the string to append when truncating (default: "...").
	Ellipsis string

	// Strategy specifies where to truncate.
	Strategy TruncateStrategy
}

// TruncateStrategy specifies where to truncate text.
type TruncateStrategy int

const (
	// TruncateEnd truncates at the end: "Hello wo..."
	TruncateEnd TruncateStrategy = iota

	// TruncateMiddle truncates in the middle: "Hel...rld"
	TruncateMiddle

	// TruncateStart truncates at the start: "...o world"
	TruncateStart
)

// Truncate shortens text to fit within maxWidth, adding an ellipsis.
//
// Uses UAX #29 to respect grapheme cluster boundaries, ensuring emoji
// and combining marks are not broken.
//
// Example:
//
//	txt := text.NewTerminal()
//	short := txt.Truncate("Hello 世界!", text.TruncateOptions{
//	    MaxWidth: 10,
//	    Strategy: text.TruncateEnd,
//	})
//	fmt.Println(short)  // "Hello 世..."
func (t *Text) Truncate(text string, opts TruncateOptions) string {
	if opts.Ellipsis == "" {
		opts.Ellipsis = "..."
	}

	textWidth := t.Width(text)
	if textWidth <= opts.MaxWidth {
		return text
	}

	ellipsisWidth := t.Width(opts.Ellipsis)
	if ellipsisWidth >= opts.MaxWidth {
		return ""
	}

	targetWidth := opts.MaxWidth - ellipsisWidth

	// Use UAX #29 to respect grapheme boundaries
	graphemes := uax29.Graphemes(text)

	switch opts.Strategy {
	case TruncateEnd:
		return t.truncateEnd(graphemes, targetWidth, opts.Ellipsis)
	case TruncateMiddle:
		return t.truncateMiddle(graphemes, targetWidth, opts.Ellipsis)
	case TruncateStart:
		return t.truncateStart(graphemes, targetWidth, opts.Ellipsis)
	default:
		return t.truncateEnd(graphemes, targetWidth, opts.Ellipsis)
	}
}

func (t *Text) truncateEnd(graphemes []string, targetWidth float64, ellipsis string) string {
	result := ""
	width := 0.0

	for _, g := range graphemes {
		gWidth := t.Width(g)
		if width+gWidth > targetWidth {
			break
		}
		result += g
		width += gWidth
	}

	return result + ellipsis
}

func (t *Text) truncateMiddle(graphemes []string, targetWidth float64, ellipsis string) string {
	if len(graphemes) == 0 {
		return ellipsis
	}

	leftWidth := targetWidth / 2
	rightWidth := targetWidth - leftWidth

	// Build left side
	left := ""
	width := 0.0
	for _, g := range graphemes {
		gWidth := t.Width(g)
		if width+gWidth > leftWidth {
			break
		}
		left += g
		width += gWidth
	}

	// Build right side
	right := ""
	width = 0.0
	for i := len(graphemes) - 1; i >= 0; i-- {
		g := graphemes[i]
		gWidth := t.Width(g)
		if width+gWidth > rightWidth {
			break
		}
		right = g + right
		width += gWidth
	}

	return left + ellipsis + right
}

func (t *Text) truncateStart(graphemes []string, targetWidth float64, ellipsis string) string {
	result := ""
	width := 0.0

	for i := len(graphemes) - 1; i >= 0; i-- {
		g := graphemes[i]
		gWidth := t.Width(g)
		if width+gWidth > targetWidth {
			break
		}
		result = g + result
		width += gWidth
	}

	return ellipsis + result
}

// ═══════════════════════════════════════════════════════════════
//  Direction
// ═══════════════════════════════════════════════════════════════

// Direction specifies text directionality (LTR or RTL).
//
// Specification:
//   - CSS Writing Modes Level 3: https://www.w3.org/TR/css-writing-modes-3/#direction
//   - MDN direction: https://developer.mozilla.org/en-US/docs/Web/CSS/direction
type Direction int

const (
	// DirectionLTR is left-to-right text direction (Latin, Cyrillic, etc).
	DirectionLTR Direction = iota

	// DirectionRTL is right-to-left text direction (Arabic, Hebrew, etc).
	DirectionRTL

	// DirectionAuto determines direction from content using Unicode bidi algorithm.
	DirectionAuto
)

// ═══════════════════════════════════════════════════════════════
//  Alignment
// ═══════════════════════════════════════════════════════════════

// Alignment specifies text alignment.
//
// Specification:
//   - CSS Text Level 3: https://www.w3.org/TR/css-text-3/#text-align-property
//   - CSS Text Level 3 (text-align-last): https://www.w3.org/TR/css-text-3/#text-align-last-property
//   - MDN text-align: https://developer.mozilla.org/en-US/docs/Web/CSS/text-align
//   - MDN text-align-last: https://developer.mozilla.org/en-US/docs/Web/CSS/text-align-last
type Alignment int

const (
	// AlignLeft aligns text to the left (physical direction).
	AlignLeft Alignment = iota

	// AlignCenter centers text.
	AlignCenter

	// AlignRight aligns text to the right (physical direction).
	AlignRight

	// AlignJustify distributes padding between words.
	AlignJustify

	// AlignStart aligns text to the start edge (flow-relative).
	// Left in LTR, right in RTL.
	AlignStart

	// AlignEnd aligns text to the end edge (flow-relative).
	// Right in LTR, left in RTL.
	AlignEnd

	// AlignMatchParent inherits parent's alignment, computed for direction.
	AlignMatchParent
)

// Align pads text to a specific width with the specified alignment.
// For flow-relative alignment (start/end), assumes LTR direction.
// Use AlignWithDirection for RTL support.
//
// The width parameter is in abstract units (cells for terminals, pixels for canvas).
//
// Example:
//
//	txt := text.NewTerminal()
//	aligned := txt.Align("Hello", 20, text.AlignCenter)
//	fmt.Printf("|%s|", aligned)  // "|       Hello        |" (20 cells total)
func (t *Text) Align(text string, width float64, align Alignment) string {
	return t.AlignWithDirection(text, width, align, DirectionLTR, AlignLeft)
}

// AlignWithDirection pads text to a specific width with the specified alignment,
// respecting text direction for flow-relative alignments (start/end/match-parent).
//
// Parameters:
//   - text: The text to align
//   - width: The target width
//   - align: The alignment mode
//   - direction: Text direction (LTR, RTL, or Auto)
//   - parentAlign: Parent's alignment (used for match-parent)
//
// Example:
//
//	txt := text.NewTerminal()
//	// In RTL context, start means right
//	aligned := txt.AlignWithDirection("مرحبا", 20, text.AlignStart, text.DirectionRTL, text.AlignLeft)
func (t *Text) AlignWithDirection(text string, width float64, align Alignment, direction Direction, parentAlign Alignment) string {
	textWidth := t.Width(text)

	if textWidth >= width {
		return text
	}

	padding := width - textWidth

	// Resolve flow-relative alignments
	resolvedAlign := t.resolveAlignment(align, direction, parentAlign)

	switch resolvedAlign {
	case AlignLeft:
		return text + t.makePadding(padding)
	case AlignRight:
		return t.makePadding(padding) + text
	case AlignCenter:
		leftPad := padding / 2
		rightPad := padding - leftPad
		return t.makePadding(leftPad) + text + t.makePadding(rightPad)
	case AlignJustify:
		return t.justify(text, padding)
	default:
		return text
	}
}

// resolveAlignment resolves flow-relative alignments (start/end/match-parent) to physical alignments.
func (t *Text) resolveAlignment(align Alignment, direction Direction, parentAlign Alignment) Alignment {
	// Handle match-parent first
	if align == AlignMatchParent {
		// Compute parent's alignment for current direction
		align = t.resolveAlignment(parentAlign, direction, AlignLeft)
		return align
	}

	// Resolve direction if auto
	if direction == DirectionAuto {
		direction = DirectionLTR // Default to LTR (could use UAX #9 for auto-detection)
	}

	// Resolve start/end based on direction
	switch align {
	case AlignStart:
		if direction == DirectionRTL {
			return AlignRight
		}
		return AlignLeft
	case AlignEnd:
		if direction == DirectionRTL {
			return AlignLeft
		}
		return AlignRight
	default:
		return align
	}
}

// makePadding creates padding of specified width using spaces.
func (t *Text) makePadding(width float64) string {
	if width <= 0 {
		return ""
	}

	spaceWidth := t.config.MeasureFunc(' ')
	count := int(width / spaceWidth)
	return strings.Repeat(" ", count)
}

// justify distributes padding between words.
func (t *Text) justify(text string, padding float64) string {
	// Simple justification: distribute padding between words
	words := strings.Fields(text)
	if len(words) <= 1 {
		return text
	}

	gaps := len(words) - 1
	extraSpacePerGap := padding / float64(gaps)

	result := words[0]
	for i := 1; i < len(words); i++ {
		result += t.makePadding(t.Width(" ") + extraSpacePerGap)
		result += words[i]
	}

	return result
}

// ═══════════════════════════════════════════════════════════════
//  Bidirectional Support
// ═══════════════════════════════════════════════════════════════

// Reorder applies the bidirectional algorithm for display.
//
// Uses UAX #9 to properly reorder mixed LTR/RTL text (e.g., Latin + Arabic).
//
// Example:
//
//	txt := text.NewTerminal()
//	display := txt.Reorder("Hello שלום world")
//	fmt.Println(display)  // Properly reordered for display
func (t *Text) Reorder(text string) string {
	return uax9.Reorder(text, t.config.BaseDirection)
}

// ReorderWithDirection applies bidirectional algorithm with explicit direction.
func (t *Text) ReorderWithDirection(text string, dir uax9.Direction) string {
	return uax9.Reorder(text, dir)
}

// DetectDirection automatically detects paragraph direction.
func (t *Text) DetectDirection(text string) uax9.Direction {
	return uax9.GetParagraphDirection(text)
}

// ═══════════════════════════════════════════════════════════════
//  Grapheme Operations
// ═══════════════════════════════════════════════════════════════

// Graphemes splits text into grapheme clusters (user-perceived characters).
//
// Uses UAX #29 to properly handle:
//   - Emoji sequences with ZWJ: 👨‍👩‍👧‍👦
//   - Emoji with skin tones: 👋🏻
//   - Combining marks: é (e + ́)
//   - Regional indicators: 🇺🇸
//
// Example:
//
//	txt := text.NewTerminal()
//	graphemes := txt.Graphemes("Hello👋🏻")
//	fmt.Println(len(graphemes))  // 6, not 7 (emoji+modifier is 1 grapheme)
func (t *Text) Graphemes(text string) []string {
	return uax29.Graphemes(text)
}

// GraphemeCount returns the number of grapheme clusters.
func (t *Text) GraphemeCount(text string) int {
	return len(uax29.Graphemes(text))
}

// GraphemeAt returns the grapheme cluster at the specified index.
func (t *Text) GraphemeAt(text string, index int) string {
	graphemes := uax29.Graphemes(text)
	if index >= 0 && index < len(graphemes) {
		return graphemes[index]
	}
	return ""
}
