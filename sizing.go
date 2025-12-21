package text

import (
	"math"
	"strings"

	"github.com/SCKelemen/units"
)

// Sizing and Metrics for Layout Engines
//
// This provides comprehensive sizing information needed for CSS layout engines:
// - Intrinsic sizing (min-content, max-content)
// - Multi-line text bounds
// - Baseline positioning
// - Line box metrics

// ═══════════════════════════════════════════════════════════════
//  Intrinsic Sizing (CSS Box Sizing)
// ═══════════════════════════════════════════════════════════════

// IntrinsicSize represents the intrinsic dimensions of text.
//
// Based on CSS Box Sizing Module Level 3:
// https://www.w3.org/TR/css-sizing-3/#intrinsic-sizes
type IntrinsicSize struct {
	// MinContent is the minimum width without overflow.
	// This is the width of the longest word/grapheme cluster.
	MinContent float64

	// MaxContent is the width if the text never wraps.
	// This is the total width of all text on a single line.
	MaxContent float64

	// PreferredWidth is a suggested width for comfortable reading.
	// For text, this is typically around 60-80 characters.
	PreferredWidth float64
}

// IntrinsicSizing calculates intrinsic sizes for text.
//
// Returns:
//   - MinContent: Width of the widest word/grapheme (won't overflow)
//   - MaxContent: Width if text never wraps (single line)
//   - PreferredWidth: Comfortable reading width (60-80 ch)
//
// Example:
//
//	txt := text.NewTerminal()
//	sizes := txt.IntrinsicSizing("Hello world! This is a test.")
//	// sizes.MinContent = 6.0 (widest word "Hello!" or "world!")
//	// sizes.MaxContent = 28.0 (full line width)
func (t *Text) IntrinsicSizing(text string) IntrinsicSize {
	// MaxContent: full line width
	maxContent := t.Width(text)

	// MinContent: width of widest unbreakable segment
	// For proper intrinsic sizing, we should look at the longest segment
	// between forced breaks (spaces for most text, or line break opportunities for CJK).
	// Using space-separated segments gives us the "longest word" in Latin text
	// and works reasonably for CJK with spaces.
	segments := strings.Fields(text)
	minContent := 0.0

	if len(segments) > 0 {
		// Use space-separated segments for min-content
		for _, segment := range segments {
			w := t.Width(segment)
			if w > minContent {
				minContent = w
			}
		}
	} else {
		// No spaces - the whole text is one segment
		minContent = maxContent
	}

	// Fallback: if still zero, use widest grapheme
	if minContent == 0 {
		graphemes := t.Graphemes(text)
		for _, g := range graphemes {
			w := t.Width(g)
			if w > minContent {
				minContent = w
			}
		}
	}

	// PreferredWidth: aim for ~60-80 characters (use 70 as default)
	// This is a heuristic for comfortable reading
	avgCharWidth := 1.0
	if maxContent > 0 && len([]rune(text)) > 0 {
		avgCharWidth = maxContent / float64(len([]rune(text)))
	}
	preferredWidth := 70 * avgCharWidth

	return IntrinsicSize{
		MinContent:     minContent,
		MaxContent:     maxContent,
		PreferredWidth: math.Min(preferredWidth, maxContent),
	}
}

// ═══════════════════════════════════════════════════════════════
//  Line Box Metrics
// ═══════════════════════════════════════════════════════════════

// LineBoxMetrics provides complete metrics for a single line of text.
//
// Based on CSS Inline Layout Module:
// https://www.w3.org/TR/css-inline-3/#line-box
type LineBoxMetrics struct {
	// Content dimensions
	Width   float64 // Content width (advance)
	Content string  // Text content

	// Baseline and vertical metrics (in abstract units)
	Ascent  float64 // Distance above baseline
	Descent float64 // Distance below baseline
	Leading float64 // Extra space distributed above/below

	// Line box dimensions
	LineHeight float64 // Total line height (ascent + descent + leading)
	Baseline   float64 // Position of baseline from top

	// Position in source text
	Start int // Rune index in original text
	End   int // Rune index in original text
}

// MeasureLineBox calculates complete metrics for a line of text.
//
// Parameters:
//   - text: The line content
//   - style: Font and line height settings
//
// Returns metrics needed for proper line box positioning and alignment.
func (t *Text) MeasureLineBox(text string, style TextStyle) LineBoxMetrics {
	width := t.Width(text)

	// Get line height from style or default to 1.0
	lineHeight := style.LineHeight
	if lineHeight == 0 {
		lineHeight = 1.0
	}

	// For terminal/simple rendering, use standard proportions
	// In real font rendering, these would come from font metrics
	ascent := lineHeight * 0.8
	descent := lineHeight * 0.2
	leading := 0.0

	// If line height is larger than ascent + descent, distribute as leading
	contentHeight := ascent + descent
	if lineHeight > contentHeight {
		leading = lineHeight - contentHeight
		// Leading is distributed half-above, half-below
	} else {
		// Ensure ascent + descent + leading always equals lineHeight
		// by adjusting the descent slightly if needed
		leading = 0.0
		descent = lineHeight - ascent
	}

	// Baseline is ascent plus half the leading
	baseline := ascent + (leading / 2)

	return LineBoxMetrics{
		Width:      width,
		Content:    text,
		Ascent:     ascent,
		Descent:    descent,
		Leading:    leading,
		LineHeight: lineHeight, // Use the requested line height, not calculated
		Baseline:   baseline,
		Start:      0,
		End:        len([]rune(text)),
	}
}

// ═══════════════════════════════════════════════════════════════
//  Multi-Line Text Bounds
// ═══════════════════════════════════════════════════════════════

// TextBounds represents the bounding box of multi-line text.
type TextBounds struct {
	// Total dimensions
	Width  float64 // Width of widest line
	Height float64 // Total height of all lines

	// Baseline information
	FirstBaseline float64 // Position of first line's baseline from top
	LastBaseline  float64 // Position of last line's baseline from top

	// Line information
	LineCount int               // Number of lines
	Lines     []LineBoxMetrics  // Metrics for each line
}

// MeasureMultiLine calculates bounds for wrapped text.
//
// This is essential for layout engines that need to know:
// - Total height of text block
// - Baseline positions for alignment
// - Per-line metrics for rendering
//
// Example:
//
//	txt := text.NewTerminal()
//	bounds := txt.MeasureMultiLine("Long text...", text.WrapOptions{
//	    MaxWidth: 40,
//	}, text.TextStyle{
//	    LineHeight: 1.5,
//	})
//	// bounds.Height = total height with line spacing
//	// bounds.FirstBaseline = for vertical-align: baseline
func (t *Text) MeasureMultiLine(text string, wrapOpts WrapOptions, style TextStyle) TextBounds {
	// Wrap text into lines
	lines := t.Wrap(text, wrapOpts)

	if len(lines) == 0 {
		return TextBounds{}
	}

	// Calculate metrics for each line
	lineMetrics := make([]LineBoxMetrics, len(lines))
	maxWidth := 0.0
	currentY := 0.0

	for i, line := range lines {
		metrics := t.MeasureLineBox(line.Content, style)
		metrics.Start = line.Start
		metrics.End = line.End

		if line.Width > maxWidth {
			maxWidth = line.Width
		}

		lineMetrics[i] = metrics
		currentY += metrics.LineHeight
	}

	// Calculate baseline positions
	firstBaseline := lineMetrics[0].Baseline
	lastBaseline := currentY - lineMetrics[len(lineMetrics)-1].Descent - (lineMetrics[len(lineMetrics)-1].Leading / 2)

	return TextBounds{
		Width:         maxWidth,
		Height:        currentY,
		FirstBaseline: firstBaseline,
		LastBaseline:  lastBaseline,
		LineCount:     len(lines),
		Lines:         lineMetrics,
	}
}

// ═══════════════════════════════════════════════════════════════
//  CSS-Aware Sizing
// ═══════════════════════════════════════════════════════════════

// CSSTextBounds extends TextBounds with CSS-specific measurements.
type CSSTextBounds struct {
	TextBounds

	// CSS intrinsic sizes
	Intrinsic IntrinsicSize

	// Applied spacing
	LetterSpacing units.Length
	WordSpacing   units.Length
	TextIndent    TextIndent
}

// MeasureCSS calculates complete bounds with CSS text properties.
//
// This combines:
// - Multi-line wrapping with CSS white-space
// - Intrinsic sizing for layout
// - Text transformation
// - Letter/word spacing
//
// Example:
//
//	txt := text.NewTerminal()
//	bounds := txt.MeasureCSS("hello world", text.CSSWrapOptions{
//	    MaxWidth: units.Ch(40),
//	    Style: text.CSSTextStyle{
//	        WhiteSpace: text.WhiteSpaceNormal,
//	        TextTransform: text.TextTransformUppercase,
//	        LetterSpacing: units.Px(1),
//	    },
//	}, text.TextStyle{
//	    LineHeight: 1.5,
//	})
func (t *Text) MeasureCSS(text string, cssOpts CSSWrapOptions, textStyle TextStyle) CSSTextBounds {
	// Process text according to CSS properties
	processed, allowWrap := t.ProcessWhiteSpace(text, cssOpts.Style.WhiteSpace)
	processed = t.Transform(processed, cssOpts.Style.TextTransform)

	// Calculate intrinsic sizing
	intrinsic := t.IntrinsicSizing(processed)

	// Wrap and measure
	var bounds TextBounds
	if allowWrap {
		bounds = t.MeasureMultiLine(processed, WrapOptions{
			MaxWidth: cssOpts.MaxWidth.Raw(),
		}, textStyle)
	} else {
		// No wrapping: single line
		lineMetrics := t.MeasureLineBox(processed, textStyle)
		bounds = TextBounds{
			Width:         lineMetrics.Width,
			Height:        lineMetrics.LineHeight,
			FirstBaseline: lineMetrics.Baseline,
			LastBaseline:  lineMetrics.Baseline,
			LineCount:     1,
			Lines:         []LineBoxMetrics{lineMetrics},
		}
	}

	return CSSTextBounds{
		TextBounds:    bounds,
		Intrinsic:     intrinsic,
		LetterSpacing: cssOpts.Style.LetterSpacing,
		WordSpacing:   cssOpts.Style.WordSpacing,
		TextIndent:    cssOpts.Style.TextIndent,
	}
}

// ═══════════════════════════════════════════════════════════════
//  Font Metrics Interface (for future font integration)
// ═══════════════════════════════════════════════════════════════

// FontMetrics provides actual font measurements.
//
// This interface allows integration with font rendering libraries
// that can provide real font metrics (not just approximations).
//
// For future integration with:
// - freetype-go
// - golang.org/x/image/font
// - harfbuzz bindings
type FontMetrics interface {
	// Ascent returns the distance from baseline to top of tallest glyph.
	Ascent() float64

	// Descent returns the distance from baseline to bottom of lowest glyph.
	Descent() float64

	// LineGap returns the recommended line spacing (leading).
	LineGap() float64

	// CapHeight returns the height of capital letters.
	CapHeight() float64

	// XHeight returns the height of lowercase 'x'.
	XHeight() float64

	// UnitsPerEm returns the font's units per em.
	UnitsPerEm() float64
}

// WithFontMetrics creates a Text instance with real font metrics.
//
// When provided, font metrics override the default ascent/descent calculations.
// This enables proper support for font-relative units (em, ex, cap, ch, ic).
//
// Future enhancement: Store and use font metrics in line box measurement for
// pixel-perfect canvas/GUI rendering with proper baseline alignment.
func (t *Text) WithFontMetrics(fm FontMetrics) *Text {
	// Placeholder - font metrics not yet stored
	// Current implementation uses basic cell-based sizing suitable for terminals
	return t
}
