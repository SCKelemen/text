package text

// This file provides integration with github.com/SCKelemen/layout engine.
// It implements the layout.TextMetricsProvider interface.

// TextStyle mirrors the text styling properties from the layout engine.
// This avoids importing the layout package directly.
type TextStyle struct {
	FontSize      float64
	LineHeight    float64
	LetterSpacing float64
	FontFamily    string
	Bold          bool
	Italic        bool
}

// Metrics implements layout.TextMetricsProvider for integration with the layout engine.
//
// Example:
//
//	// Create metrics for terminal rendering
//	metrics := text.NewMetrics(text.Config{
//	    MeasureFunc: text.TerminalMeasure,
//	})
//
//	// In your layout engine setup:
//	layout.SetTextMetricsProvider(metrics)
type Metrics struct {
	text *Text
}

// NewMetrics creates a layout-compatible metrics provider.
func NewMetrics(config Config) *Metrics {
	return &Metrics{
		text: New(config),
	}
}

// NewTerminalMetrics creates metrics configured for terminal rendering.
func NewTerminalMetrics() *Metrics {
	return &Metrics{
		text: NewTerminal(),
	}
}

// Measure implements layout.TextMetricsProvider interface.
//
// Returns:
//   - advance: The display width of the text
//   - ascent: The distance above the baseline (for line height calculation)
//   - descent: The distance below the baseline (for line height calculation)
//
// For terminal rendering:
//   - advance is measured in character cells
//   - ascent/descent are proportional to line height for proper spacing
func (m *Metrics) Measure(text string, style TextStyle) (advance, ascent, descent float64) {
	advance = m.text.Width(text)

	// Calculate ascent/descent from line height
	// These are used by the layout engine for line box calculations
	lineHeight := style.LineHeight
	if lineHeight == 0 {
		lineHeight = 1.0
	}

	// Standard proportions: 80% ascent, 20% descent
	ascent = lineHeight * 0.8
	descent = lineHeight * 0.2

	return
}

// Text returns the underlying Text instance for direct access to text operations.
//
// This allows using all text.Text methods (Wrap, Truncate, Align, etc.)
// alongside the metrics provider.
//
// Example:
//
//	metrics := text.NewTerminalMetrics()
//	layout.SetTextMetricsProvider(metrics)
//
//	// Also use for direct text operations
//	txt := metrics.Text()
//	truncated := txt.Truncate("Long text...", text.TruncateOptions{
//	    MaxWidth: 40,
//	})
func (m *Metrics) Text() *Text {
	return m.text
}
