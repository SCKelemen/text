package text

import (
	"github.com/SCKelemen/unicode/uax50"
)

// Vertical Text Layout Support (UAX #50)
// Based on: https://www.w3.org/TR/css-writing-modes-4/

// ═══════════════════════════════════════════════════════════════
//  Writing Mode (CSS Writing Modes Level 4)
// ═══════════════════════════════════════════════════════════════

// WritingMode specifies the block flow direction.
// Based on CSS Writing Modes Module Level 4: https://www.w3.org/TR/css-writing-modes-4/#block-flow
type WritingMode int

const (
	// WritingModeHorizontalTB flows top to bottom, inline left to right.
	// This is the default for Latin scripts.
	WritingModeHorizontalTB WritingMode = iota

	// WritingModeVerticalRL flows right to left, inline top to bottom.
	// Used in traditional Chinese, Japanese, Korean typography.
	WritingModeVerticalRL

	// WritingModeVerticalLR flows left to right, inline top to bottom.
	// Used in some Mongolian and certain historical scripts.
	WritingModeVerticalLR

	// WritingModeSidewaysRL flows right to left with characters rotated 90° clockwise.
	// Used for mixed scripts in vertical layout.
	WritingModeSidewaysRL

	// WritingModeSidewaysLR flows left to right with characters rotated 90° counter-clockwise.
	WritingModeSidewaysLR
)

// ═══════════════════════════════════════════════════════════════
//  Text Orientation (CSS Writing Modes Level 4)
// ═══════════════════════════════════════════════════════════════

// TextOrientation controls glyph orientation in vertical text.
// Based on CSS Writing Modes Module Level 4 §7.1: https://www.w3.org/TR/css-writing-modes-4/#text-orientation
type TextOrientation int

const (
	// TextOrientationMixed uses upright orientation for CJK and rotates other characters.
	// This is the default for vertical writing modes.
	TextOrientationMixed TextOrientation = iota

	// TextOrientationUpright forces all characters to upright orientation.
	TextOrientationUpright

	// TextOrientationSideways rotates all characters 90° clockwise.
	TextOrientationSideways

	// TextOrientationSidewaysRight is deprecated, use Sideways instead.
	TextOrientationSidewaysRight
)

// ═══════════════════════════════════════════════════════════════
//  Text Combine Upright (Tate-chu-yoko)
// ═══════════════════════════════════════════════════════════════

// TextCombineUpright controls horizontal text runs in vertical layout.
// Based on CSS Writing Modes Module Level 4 §9: https://www.w3.org/TR/css-writing-modes-4/#text-combine-upright
type TextCombineUpright int

const (
	// TextCombineUprightNone disables text combining.
	TextCombineUprightNone TextCombineUpright = iota

	// TextCombineUprightAll enables combining for all characters.
	// Used for numbers, acronyms, etc. in vertical text (tate-chu-yoko in Japanese).
	TextCombineUprightAll

	// TextCombineUprightDigits enables combining for digit sequences only.
	TextCombineUprightDigits
)

// ═══════════════════════════════════════════════════════════════
//  Vertical Text Configuration
// ═══════════════════════════════════════════════════════════════

// VerticalTextStyle configures vertical text layout properties.
type VerticalTextStyle struct {
	WritingMode        WritingMode
	TextOrientation    TextOrientation
	TextCombineUpright TextCombineUpright

	// GlyphOrientationVertical controls glyph rotation in vertical text.
	// Deprecated in favor of TextOrientation, but kept for compatibility.
	// Value in degrees: 0, 90, -90, or auto.
	GlyphOrientationVertical float64
}

// DefaultVerticalTextStyle returns default vertical text style (horizontal layout).
func DefaultVerticalTextStyle() VerticalTextStyle {
	return VerticalTextStyle{
		WritingMode:              WritingModeHorizontalTB,
		TextOrientation:          TextOrientationMixed,
		TextCombineUpright:       TextCombineUprightNone,
		GlyphOrientationVertical: 0,
	}
}

// ═══════════════════════════════════════════════════════════════
//  Character Orientation Detection
// ═══════════════════════════════════════════════════════════════

// CharOrientation returns the orientation for a character in vertical text.
//
// This implements the UAX #50 algorithm combined with CSS text-orientation property.
// For TextOrientationMixed, it respects UAX #50 defaults.
// For TextOrientationUpright, all characters are upright.
// For TextOrientationSideways, all characters are rotated.
func (t *Text) CharOrientation(r rune, orientation TextOrientation) uax50.Orientation {
	switch orientation {
	case TextOrientationUpright:
		return uax50.Upright

	case TextOrientationSideways, TextOrientationSidewaysRight:
		return uax50.Rotated

	case TextOrientationMixed:
		// Use UAX #50 defaults
		return uax50.LookupOrientation(r)

	default:
		return uax50.LookupOrientation(r)
	}
}

// IsUpright returns true if the character should be displayed upright in vertical text.
func (t *Text) IsUpright(r rune, style VerticalTextStyle) bool {
	orientation := t.CharOrientation(r, style.TextOrientation)
	return orientation == uax50.Upright || orientation == uax50.TransformedUpright
}

// IsRotated returns true if the character should be rotated in vertical text.
func (t *Text) IsRotated(r rune, style VerticalTextStyle) bool {
	orientation := t.CharOrientation(r, style.TextOrientation)
	return orientation == uax50.Rotated || orientation == uax50.TransformedRotated
}

// ═══════════════════════════════════════════════════════════════
//  Vertical Metrics
// ═══════════════════════════════════════════════════════════════

// VerticalMetrics provides measurements for vertical text layout.
type VerticalMetrics struct {
	// Advance is the vertical advance (height) for this text segment.
	Advance float64

	// InlineSize is the inline dimension (width in vertical layout).
	InlineSize float64

	// BlockSize is the block dimension (height in vertical layout).
	BlockSize float64

	// BaselineOffset is the offset from the baseline.
	BaselineOffset float64
}

// MeasureVertical measures text dimensions for vertical layout.
//
// In vertical layout:
//   - Advance is the vertical distance (top to bottom or bottom to top)
//   - InlineSize is the width (perpendicular to flow)
//   - BlockSize is the height (parallel to flow)
func (t *Text) MeasureVertical(text string, style VerticalTextStyle) VerticalMetrics {
	var metrics VerticalMetrics

	switch style.WritingMode {
	case WritingModeHorizontalTB:
		// Horizontal layout: use regular width as advance
		metrics.Advance = t.Width(text)
		metrics.InlineSize = t.Width(text)
		metrics.BlockSize = 1.0 // Assume 1 line height

	case WritingModeVerticalRL, WritingModeVerticalLR:
		// Vertical layout: each character takes vertical space
		graphemes := t.Graphemes(text)
		metrics.Advance = float64(len(graphemes))

		// Calculate maximum inline size (widest character)
		maxWidth := 0.0
		for _, g := range graphemes {
			w := t.Width(g)
			if w > maxWidth {
				maxWidth = w
			}
		}
		metrics.InlineSize = maxWidth
		metrics.BlockSize = float64(len(graphemes))

	case WritingModeSidewaysRL, WritingModeSidewaysLR:
		// Sideways: rotated horizontal text
		metrics.Advance = t.Width(text)
		metrics.InlineSize = 1.0 // Rotated height
		metrics.BlockSize = t.Width(text)
	}

	return metrics
}

// ═══════════════════════════════════════════════════════════════
//  Vertical Line Breaking
// ═══════════════════════════════════════════════════════════════

// VerticalWrapOptions extends WrapOptions for vertical text.
type VerticalWrapOptions struct {
	// MaxBlockSize is the maximum block dimension before wrapping.
	// For vertical text, this is the maximum column height.
	MaxBlockSize float64

	// Style configures vertical text properties.
	Style VerticalTextStyle

	// BaseOptions provides standard wrapping options.
	BaseOptions WrapOptions
}

// VerticalLine represents a line in vertical layout (a column).
type VerticalLine struct {
	Content    string
	Advance    float64 // Vertical advance (column height)
	InlineSize float64 // Horizontal size (column width)
	Start      int
	End        int
}

// WrapVertical wraps text for vertical layout.
//
// In vertical layout, "lines" are vertical columns that flow from top to bottom.
// When a column reaches MaxBlockSize, text wraps to the next column.
func (t *Text) WrapVertical(text string, opts VerticalWrapOptions) []VerticalLine {
	if opts.MaxBlockSize <= 0 {
		// No wrapping
		metrics := t.MeasureVertical(text, opts.Style)
		return []VerticalLine{{
			Content:    text,
			Advance:    metrics.Advance,
			InlineSize: metrics.InlineSize,
			Start:      0,
			End:        len([]rune(text)),
		}}
	}

	// For vertical text, wrap by grapheme clusters
	graphemes := t.Graphemes(text)
	var lines []VerticalLine

	currentColumn := ""
	currentHeight := 0.0
	maxWidth := 0.0
	columnStart := 0

	for i, g := range graphemes {
		gHeight := 1.0 // Each grapheme takes 1 unit of vertical space
		gWidth := t.Width(g)

		// Check if adding this grapheme exceeds the column height
		if currentHeight+gHeight > opts.MaxBlockSize && currentHeight > 0 {
			// Start new column
			lines = append(lines, VerticalLine{
				Content:    currentColumn,
				Advance:    currentHeight,
				InlineSize: maxWidth,
				Start:      columnStart,
				End:        columnStart + len([]rune(currentColumn)),
			})

			currentColumn = g
			currentHeight = gHeight
			maxWidth = gWidth
			columnStart = i
		} else {
			currentColumn += g
			currentHeight += gHeight
			if gWidth > maxWidth {
				maxWidth = gWidth
			}
		}
	}

	// Add final column
	if currentColumn != "" {
		lines = append(lines, VerticalLine{
			Content:    currentColumn,
			Advance:    currentHeight,
			InlineSize: maxWidth,
			Start:      columnStart,
			End:        columnStart + len([]rune(currentColumn)),
		})
	}

	return lines
}

// ═══════════════════════════════════════════════════════════════
//  Utility Functions
// ═══════════════════════════════════════════════════════════════

// IsVerticalWritingMode returns true if the writing mode is vertical.
func IsVerticalWritingMode(mode WritingMode) bool {
	return mode == WritingModeVerticalRL ||
		mode == WritingModeVerticalLR ||
		mode == WritingModeSidewaysRL ||
		mode == WritingModeSidewaysLR
}

// IsHorizontalWritingMode returns true if the writing mode is horizontal.
func IsHorizontalWritingMode(mode WritingMode) bool {
	return mode == WritingModeHorizontalTB
}
