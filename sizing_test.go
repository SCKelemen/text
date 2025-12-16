package text

import (
	"testing"

	"github.com/SCKelemen/units"
)

// ═══════════════════════════════════════════════════════════════
//  Intrinsic Sizing Tests
// ═══════════════════════════════════════════════════════════════

func TestIntrinsicSizing(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name              string
		text              string
		wantMinContentMin float64 // Minimum expected min-content
		wantMaxContent    float64
	}{
		{
			name:              "Single word",
			text:              "Hello",
			wantMinContentMin: 5.0,
			wantMaxContent:    5.0,
		},
		{
			name:              "Multiple words",
			text:              "Hello world",
			wantMinContentMin: 5.0, // At least as wide as "Hello" or "world"
			wantMaxContent:    11.0,
		},
		{
			name:              "Long word",
			text:              "Supercalifragilisticexpialidocious",
			wantMinContentMin: 34.0,
			wantMaxContent:    34.0,
		},
		{
			name:              "CJK text",
			text:              "世界 你好",
			wantMinContentMin: 4.0, // "世界" or "你好" = 4 cells
			wantMaxContent:    9.0, // 4 + 1 space + 4
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sizes := txt.IntrinsicSizing(tt.text)

			if sizes.MinContent < tt.wantMinContentMin {
				t.Errorf("MinContent = %.1f, want >= %.1f", sizes.MinContent, tt.wantMinContentMin)
			}

			if sizes.MaxContent != tt.wantMaxContent {
				t.Errorf("MaxContent = %.1f, want %.1f", sizes.MaxContent, tt.wantMaxContent)
			}

			if sizes.PreferredWidth > sizes.MaxContent {
				t.Errorf("PreferredWidth %.1f > MaxContent %.1f", sizes.PreferredWidth, sizes.MaxContent)
			}

			if sizes.MinContent > sizes.MaxContent {
				t.Errorf("MinContent %.1f > MaxContent %.1f", sizes.MinContent, sizes.MaxContent)
			}
		})
	}
}

// ═══════════════════════════════════════════════════════════════
//  Line Box Metrics Tests
// ═══════════════════════════════════════════════════════════════

func TestMeasureLineBox(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name       string
		text       string
		lineHeight float64
		wantWidth  float64
	}{
		{
			name:       "Simple text",
			text:       "Hello",
			lineHeight: 1.0,
			wantWidth:  5.0,
		},
		{
			name:       "CJK text",
			text:       "世界",
			lineHeight: 1.0,
			wantWidth:  4.0,
		},
		{
			name:       "With line height",
			text:       "Test",
			lineHeight: 1.5,
			wantWidth:  4.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics := txt.MeasureLineBox(tt.text, TextStyle{
				LineHeight: tt.lineHeight,
			})

			if metrics.Width != tt.wantWidth {
				t.Errorf("Width = %.1f, want %.1f", metrics.Width, tt.wantWidth)
			}

			if metrics.Content != tt.text {
				t.Errorf("Content = %q, want %q", metrics.Content, tt.text)
			}

			expectedLineHeight := tt.lineHeight
			if expectedLineHeight == 0 {
				expectedLineHeight = 1.0
			}

			if metrics.LineHeight != expectedLineHeight {
				t.Errorf("LineHeight = %.1f, want %.1f", metrics.LineHeight, expectedLineHeight)
			}

			// Ascent + Descent + Leading should equal LineHeight
			totalHeight := metrics.Ascent + metrics.Descent + metrics.Leading
			if totalHeight != metrics.LineHeight {
				t.Errorf("Ascent + Descent + Leading = %.1f, want %.1f",
					totalHeight, metrics.LineHeight)
			}

			// Baseline should be within line height
			if metrics.Baseline < 0 || metrics.Baseline > metrics.LineHeight {
				t.Errorf("Baseline %.1f outside [0, %.1f]", metrics.Baseline, metrics.LineHeight)
			}
		})
	}
}

// ═══════════════════════════════════════════════════════════════
//  Multi-Line Text Bounds Tests
// ═══════════════════════════════════════════════════════════════

func TestMeasureMultiLine(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name          string
		text          string
		maxWidth      float64
		lineHeight    float64
		wantLineCount int
	}{
		{
			name:          "Single line",
			text:          "Hello",
			maxWidth:      20.0,
			lineHeight:    1.0,
			wantLineCount: 1,
		},
		{
			name:          "Two lines",
			text:          "Hello world test",
			maxWidth:      10.0,
			lineHeight:    1.0,
			wantLineCount: 2,
		},
		{
			name:          "Three lines with line height",
			text:          "Hello world test example",
			maxWidth:      10.0,
			lineHeight:    1.5,
			wantLineCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bounds := txt.MeasureMultiLine(tt.text, WrapOptions{
				MaxWidth: tt.maxWidth,
			}, TextStyle{
				LineHeight: tt.lineHeight,
			})

			if bounds.LineCount != tt.wantLineCount {
				t.Errorf("LineCount = %d, want %d", bounds.LineCount, tt.wantLineCount)
			}

			if len(bounds.Lines) != tt.wantLineCount {
				t.Errorf("len(Lines) = %d, want %d", len(bounds.Lines), tt.wantLineCount)
			}

			// Width should be at most maxWidth (or one line if forced)
			if bounds.Width > tt.maxWidth && tt.wantLineCount > 1 {
				t.Errorf("Width %.1f > maxWidth %.1f", bounds.Width, tt.maxWidth)
			}

			// Height should be positive
			if bounds.Height <= 0 {
				t.Errorf("Height = %.1f, want > 0", bounds.Height)
			}

			// First baseline should be within first line height
			if tt.wantLineCount > 0 {
				firstLineHeight := bounds.Lines[0].LineHeight
				if bounds.FirstBaseline < 0 || bounds.FirstBaseline > firstLineHeight {
					t.Errorf("FirstBaseline %.1f outside [0, %.1f]",
						bounds.FirstBaseline, firstLineHeight)
				}
			}

			// Last baseline should be within total height
			if bounds.LastBaseline < 0 || bounds.LastBaseline > bounds.Height {
				t.Errorf("LastBaseline %.1f outside [0, %.1f]",
					bounds.LastBaseline, bounds.Height)
			}
		})
	}
}

func TestMeasureMultiLine_Empty(t *testing.T) {
	txt := NewTerminal()

	bounds := txt.MeasureMultiLine("", WrapOptions{MaxWidth: 10}, TextStyle{})

	if bounds.LineCount != 0 {
		t.Errorf("Empty text should have 0 lines, got %d", bounds.LineCount)
	}

	if bounds.Height != 0 {
		t.Errorf("Empty text should have 0 height, got %.1f", bounds.Height)
	}
}

// ═══════════════════════════════════════════════════════════════
//  CSS Text Bounds Tests
// ═══════════════════════════════════════════════════════════════

func TestMeasureCSS(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name       string
		text       string
		whiteSpace WhiteSpace
		transform  TextTransform
		maxWidth   float64
		wantLines  int
	}{
		{
			name:       "Normal white space",
			text:       "Hello    world",
			whiteSpace: WhiteSpaceNormal,
			transform:  TextTransformNone,
			maxWidth:   20.0,
			wantLines:  1,
		},
		{
			name:       "Pre preserves spaces",
			text:       "Hello    world",
			whiteSpace: WhiteSpacePre,
			transform:  TextTransformNone,
			maxWidth:   10.0,
			wantLines:  1, // No wrapping with pre
		},
		{
			name:       "Uppercase transform",
			text:       "hello world",
			whiteSpace: WhiteSpaceNormal,
			transform:  TextTransformUppercase,
			maxWidth:   20.0,
			wantLines:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bounds := txt.MeasureCSS(tt.text, CSSWrapOptions{
				MaxWidth: units.Px(tt.maxWidth),
				Style: CSSTextStyle{
					WhiteSpace:    tt.whiteSpace,
					TextTransform: tt.transform,
					LetterSpacing: units.Px(0),
					WordSpacing:   units.Px(0),
				},
			}, TextStyle{
				LineHeight: 1.0,
			})

			if bounds.LineCount != tt.wantLines {
				t.Errorf("LineCount = %d, want %d", bounds.LineCount, tt.wantLines)
			}

			// Intrinsic sizing should be calculated
			if bounds.Intrinsic.MinContent <= 0 {
				t.Error("MinContent should be > 0")
			}

			if bounds.Intrinsic.MaxContent < bounds.Intrinsic.MinContent {
				t.Errorf("MaxContent %.1f < MinContent %.1f",
					bounds.Intrinsic.MaxContent, bounds.Intrinsic.MinContent)
			}

			// Check text transformation was applied
			if tt.transform == TextTransformUppercase && bounds.LineCount > 0 {
				firstLine := bounds.Lines[0].Content
				if firstLine != txt.Transform(tt.text, tt.transform) {
					// Collapsed spaces may affect this, so just verify it's not lowercase
					for _, r := range firstLine {
						if r >= 'a' && r <= 'z' {
							t.Errorf("Text should be uppercase, got %q", firstLine)
							break
						}
					}
				}
			}
		})
	}
}

// ═══════════════════════════════════════════════════════════════
//  Benchmark Tests
// ═══════════════════════════════════════════════════════════════

func BenchmarkIntrinsicSizing(b *testing.B) {
	txt := NewTerminal()
	text := "Hello world! This is a test with multiple words and some longer text."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txt.IntrinsicSizing(text)
	}
}

func BenchmarkMeasureLineBox(b *testing.B) {
	txt := NewTerminal()
	text := "Hello world! This is a test."
	style := TextStyle{LineHeight: 1.5}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txt.MeasureLineBox(text, style)
	}
}

func BenchmarkMeasureMultiLine(b *testing.B) {
	txt := NewTerminal()
	text := "Hello world! This is a longer text that will wrap across multiple lines when constrained to a specific width."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txt.MeasureMultiLine(text, WrapOptions{MaxWidth: 40}, TextStyle{LineHeight: 1.5})
	}
}

func BenchmarkMeasureCSS(b *testing.B) {
	txt := NewTerminal()
	text := "Hello world! This is a longer text with CSS properties applied."
	opts := CSSWrapOptions{
		MaxWidth: units.Ch(40),
		Style:    DefaultCSSTextStyle(),
	}
	style := TextStyle{LineHeight: 1.5}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txt.MeasureCSS(text, opts, style)
	}
}
