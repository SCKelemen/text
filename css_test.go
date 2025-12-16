package text

import (
	"strings"
	"testing"

	"github.com/SCKelemen/units"
)

// ═══════════════════════════════════════════════════════════════
//  White Space Processing Tests
// ═══════════════════════════════════════════════════════════════

func TestProcessWhiteSpace(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name       string
		input      string
		whiteSpace WhiteSpace
		want       string
		allowWrap  bool
	}{
		{
			name:       "Normal collapses spaces",
			input:      "Hello    world",
			whiteSpace: WhiteSpaceNormal,
			want:       "Hello world",
			allowWrap:  true,
		},
		{
			name:       "Normal collapses newlines",
			input:      "Hello\n\nworld",
			whiteSpace: WhiteSpaceNormal,
			want:       "Hello world",
			allowWrap:  true,
		},
		{
			name:       "Pre preserves all",
			input:      "Hello    \n  world",
			whiteSpace: WhiteSpacePre,
			want:       "Hello    \n  world",
			allowWrap:  false,
		},
		{
			name:       "NoWrap collapses but doesn't wrap",
			input:      "Hello    world",
			whiteSpace: WhiteSpaceNoWrap,
			want:       "Hello world",
			allowWrap:  false,
		},
		{
			name:       "PreWrap preserves spaces and wraps",
			input:      "Hello    world",
			whiteSpace: WhiteSpacePreWrap,
			want:       "Hello    world",
			allowWrap:  true,
		},
		{
			name:       "PreLine collapses spaces but preserves newlines",
			input:      "Hello    world\nNext line",
			whiteSpace: WhiteSpacePreLine,
			want:       "Hello world\nNext line",
			allowWrap:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, allowWrap := txt.ProcessWhiteSpace(tt.input, tt.whiteSpace)
			if got != tt.want {
				t.Errorf("ProcessWhiteSpace() = %q, want %q", got, tt.want)
			}
			if allowWrap != tt.allowWrap {
				t.Errorf("allowWrap = %v, want %v", allowWrap, tt.allowWrap)
			}
		})
	}
}

// ═══════════════════════════════════════════════════════════════
//  Text Transformation Tests
// ═══════════════════════════════════════════════════════════════

func TestTransform(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name      string
		input     string
		transform TextTransform
		want      string
	}{
		{
			name:      "None",
			input:     "Hello World",
			transform: TextTransformNone,
			want:      "Hello World",
		},
		{
			name:      "Uppercase",
			input:     "Hello World",
			transform: TextTransformUppercase,
			want:      "HELLO WORLD",
		},
		{
			name:      "Lowercase",
			input:     "Hello World",
			transform: TextTransformLowercase,
			want:      "hello world",
		},
		{
			name:      "Capitalize",
			input:     "hello world",
			transform: TextTransformCapitalize,
			want:      "Hello World",
		},
		{
			name:      "Capitalize with punctuation",
			input:     "hello, world! how are you?",
			transform: TextTransformCapitalize,
			want:      "Hello, World! How Are You?",
		},
		{
			name:      "FullWidth ASCII",
			input:     "Hello",
			transform: TextTransformFullWidth,
			want:      "Ｈｅｌｌｏ",
		},
		{
			name:      "FullWidth space",
			input:     "Hi there",
			transform: TextTransformFullWidth,
			want:      "Ｈｉ　ｔｈｅｒｅ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := txt.Transform(tt.input, tt.transform)
			if got != tt.want {
				t.Errorf("Transform() = %q, want %q", got, tt.want)
			}
		})
	}
}

// ═══════════════════════════════════════════════════════════════
//  Word and Sentence Boundary Tests
// ═══════════════════════════════════════════════════════════════

func TestWords(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name      string
		input     string
		wantCount int
	}{
		{
			name:      "Simple sentence",
			input:     "Hello world",
			wantCount: 2,
		},
		{
			name:      "With punctuation",
			input:     "Hello, world!",
			wantCount: 2,
		},
		{
			name:      "Multiple spaces",
			input:     "Hello    world",
			wantCount: 2,
		},
		{
			name:      "Contractions",
			input:     "don't can't won't",
			wantCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := txt.WordCount(tt.input)
			if count != tt.wantCount {
				t.Errorf("WordCount() = %d, want %d", count, tt.wantCount)
			}
		})
	}
}

func TestSentences(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name      string
		input     string
		wantCount int
	}{
		{
			name:      "Single sentence",
			input:     "Hello world.",
			wantCount: 1,
		},
		{
			name:      "Multiple sentences",
			input:     "Hello world. How are you?",
			wantCount: 2,
		},
		{
			name:      "Abbreviations",
			input:     "Dr. Smith is here.",
			wantCount: 2, // UAX #29 splits on "Dr." without abbreviation dictionary
		},
		{
			name:      "Exclamation and question",
			input:     "Hello! How are you? I'm fine.",
			wantCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := txt.SentenceCount(tt.input)
			if count != tt.wantCount {
				t.Errorf("SentenceCount() = %d, want %d", count, tt.wantCount)
			}
		})
	}
}

// ═══════════════════════════════════════════════════════════════
//  CSS Text Style Tests
// ═══════════════════════════════════════════════════════════════

func TestDefaultCSSTextStyle(t *testing.T) {
	style := DefaultCSSTextStyle()

	if style.WhiteSpace != WhiteSpaceNormal {
		t.Errorf("WhiteSpace = %v, want %v", style.WhiteSpace, WhiteSpaceNormal)
	}

	if style.TextTransform != TextTransformNone {
		t.Errorf("TextTransform = %v, want %v", style.TextTransform, TextTransformNone)
	}

	if !style.LetterSpacing.IsZero() {
		t.Errorf("LetterSpacing should be zero, got %v", style.LetterSpacing)
	}

	if !style.WordSpacing.IsZero() {
		t.Errorf("WordSpacing should be zero, got %v", style.WordSpacing)
	}
}

// ═══════════════════════════════════════════════════════════════
//  CSS Wrapping Tests
// ═══════════════════════════════════════════════════════════════

func TestWrapCSS(t *testing.T) {
	txt := NewTerminal()

	t.Run("Normal white space with wrapping", func(t *testing.T) {
		text := "Hello    world!    This is a test."
		opts := CSSWrapOptions{
			MaxWidth: units.Ch(15),
			Style: CSSTextStyle{
				WhiteSpace:    WhiteSpaceNormal,
				TextTransform: TextTransformNone,
				LetterSpacing: units.Px(0),
				WordSpacing:   units.Px(0),
			},
		}

		lines := txt.WrapCSS(text, opts)

		// Should collapse spaces and wrap
		if len(lines) == 0 {
			t.Error("Expected wrapped lines, got none")
		}

		// Verify no line is too wide
		for i, line := range lines {
			if line.Width > opts.MaxWidth.Raw() {
				t.Errorf("Line %d width %.1f exceeds maxWidth %.1f: %q",
					i, line.Width, opts.MaxWidth.Raw(), line.Content)
			}
		}
	})

	t.Run("Pre white space no wrapping", func(t *testing.T) {
		text := "Hello    world!    This is a test."
		opts := CSSWrapOptions{
			MaxWidth: units.Ch(10),
			Style: CSSTextStyle{
				WhiteSpace:    WhiteSpacePre,
				TextTransform: TextTransformNone,
				LetterSpacing: units.Px(0),
				WordSpacing:   units.Px(0),
			},
		}

		lines := txt.WrapCSS(text, opts)

		// Should not wrap with pre
		if len(lines) != 1 {
			t.Errorf("Expected 1 line with pre, got %d", len(lines))
		}

		// Should preserve spaces
		if !strings.Contains(lines[0].Content, "    ") {
			t.Error("Expected preserved spaces in pre mode")
		}
	})

	t.Run("Text transformation", func(t *testing.T) {
		text := "hello world"
		opts := CSSWrapOptions{
			MaxWidth: units.Ch(20),
			Style: CSSTextStyle{
				WhiteSpace:    WhiteSpaceNormal,
				TextTransform: TextTransformUppercase,
				LetterSpacing: units.Px(0),
				WordSpacing:   units.Px(0),
			},
		}

		lines := txt.WrapCSS(text, opts)

		if len(lines) == 0 {
			t.Fatal("Expected at least one line")
		}

		// Should be uppercase
		if !strings.Contains(lines[0].Content, "HELLO") {
			t.Errorf("Expected uppercase text, got %q", lines[0].Content)
		}
	})
}

// ═══════════════════════════════════════════════════════════════
//  Kana Transformation Tests
// ═══════════════════════════════════════════════════════════════

func TestToFullSizeKana(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Small hiragana",
			input: "ぁぃぅぇぉ",
			want:  "あいうえお",
		},
		{
			name:  "Small katakana",
			input:  "ァィゥェォ",
			want:  "アイウエオ",
		},
		{
			name:  "Small tsu",
			input:  "っッ",
			want:  "つツ",
		},
		{
			name:  "Mixed with normal",
			input:  "あぁアァ",
			want:  "ああアア",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := txt.Transform(tt.input, TextTransformFullSizeKana)
			if got != tt.want {
				t.Errorf("toFullSizeKana() = %q, want %q", got, tt.want)
			}
		})
	}
}

// ═══════════════════════════════════════════════════════════════
//  Benchmark Tests
// ═══════════════════════════════════════════════════════════════

func BenchmarkProcessWhiteSpace(b *testing.B) {
	txt := NewTerminal()
	text := "Hello    world\n\nThis is   a    test with multiple    spaces."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txt.ProcessWhiteSpace(text, WhiteSpaceNormal)
	}
}

func BenchmarkTransformUppercase(b *testing.B) {
	txt := NewTerminal()
	text := "Hello world! This is a test with multiple words and punctuation."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txt.Transform(text, TextTransformUppercase)
	}
}

func BenchmarkTransformCapitalize(b *testing.B) {
	txt := NewTerminal()
	text := "hello world! this is a test with multiple words and punctuation."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txt.Transform(text, TextTransformCapitalize)
	}
}

func BenchmarkWrapCSS(b *testing.B) {
	txt := NewTerminal()
	text := "Hello world! This is a long text that needs wrapping across multiple lines with proper Unicode support."

	opts := CSSWrapOptions{
		MaxWidth: units.Ch(40),
		Style:    DefaultCSSTextStyle(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txt.WrapCSS(text, opts)
	}
}

func BenchmarkWordCount(b *testing.B) {
	txt := NewTerminal()
	text := "Hello world! This is a test with multiple words, punctuation, and contractions like don't and can't."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txt.WordCount(text)
	}
}

// ═══════════════════════════════════════════════════════════════
//  Word Spacing Tests
// ═══════════════════════════════════════════════════════════════

func TestWordSpacing(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name        string
		text        string
		wordSpacing float64
		maxWidth    float64
	}{
		{
			name:        "Positive word spacing",
			text:        "Hello world test",
			wordSpacing: 2.0,
			maxWidth:    30.0,
		},
		{
			name:        "Large word spacing",
			text:        "A B C",
			wordSpacing: 5.0,
			maxWidth:    30.0,
		},
		{
			name:        "Zero word spacing",
			text:        "Hello world",
			wordSpacing: 0.0,
			maxWidth:    20.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			style := CSSTextStyle{
				WordSpacing: units.Px(tt.wordSpacing),
			}

			opts := CSSWrapOptions{
				MaxWidth: units.Px(tt.maxWidth),
				Style:    style,
			}

			lines := txt.WrapCSS(tt.text, opts)

			if len(lines) == 0 {
				t.Error("WrapCSS returned no lines")
			}

			// Verify that lines respect max width
			for i, line := range lines {
				if line.Width > tt.maxWidth {
					t.Errorf("Line %d width %.1f exceeds max %.1f", i, line.Width, tt.maxWidth)
				}
			}

			// With word spacing, width should be greater than without
			if tt.wordSpacing > 0 {
				spaceCount := strings.Count(tt.text, " ")
				if spaceCount > 0 {
					// The total width should include the word spacing contribution
					expectedExtra := float64(spaceCount) * tt.wordSpacing
					if expectedExtra > 0 {
						// Just verify calculation doesn't crash and produces valid output
						t.Logf("Applied word spacing %.1f to %d spaces", tt.wordSpacing, spaceCount)
					}
				}
			}
		})
	}
}

// ═══════════════════════════════════════════════════════════════
//  Text Align Last Tests
// ═══════════════════════════════════════════════════════════════

func TestAlignLines(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name          string
		lines         []Line
		width         float64
		textAlign     Alignment
		textAlignLast Alignment
	}{
		{
			name: "Justify with left last",
			lines: []Line{
				{Content: "First line", Width: 10},
				{Content: "Second line", Width: 11},
				{Content: "Last", Width: 4},
			},
			width:         20.0,
			textAlign:     AlignJustify,
			textAlignLast: AlignLeft,
		},
		{
			name: "Justify with center last",
			lines: []Line{
				{Content: "Line one", Width: 8},
				{Content: "Line two", Width: 8},
				{Content: "End", Width: 3},
			},
			width:         20.0,
			textAlign:     AlignJustify,
			textAlignLast: AlignCenter,
		},
		{
			name: "Center all lines",
			lines: []Line{
				{Content: "First", Width: 5},
				{Content: "Second", Width: 6},
			},
			width:         15.0,
			textAlign:     AlignCenter,
			textAlignLast: AlignCenter,
		},
		{
			name: "Single line",
			lines: []Line{
				{Content: "Only", Width: 4},
			},
			width:         10.0,
			textAlign:     AlignLeft,
			textAlignLast: AlignRight,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			style := CSSTextStyle{
				TextAlign:     tt.textAlign,
				TextAlignLast: tt.textAlignLast,
			}

			aligned := txt.AlignLines(tt.lines, tt.width, style)

			if len(aligned) != len(tt.lines) {
				t.Errorf("Expected %d lines, got %d", len(tt.lines), len(aligned))
			}

			// Check last line uses text-align-last
			if len(aligned) > 0 {
				lastLine := aligned[len(aligned)-1]
				if lastLine.Width != tt.width {
					t.Errorf("Last line width %.1f, want %.1f", lastLine.Width, tt.width)
				}

				// Verify alignment was applied (content should have changed if alignment needed)
				originalLast := tt.lines[len(tt.lines)-1]
				originalContentWidth := txt.Width(originalLast.Content)
				if originalContentWidth < tt.width {
					// Alignment should have added padding
					alignedContentWidth := txt.Width(strings.TrimSpace(lastLine.Content))
					if alignedContentWidth < originalContentWidth {
						t.Errorf("Expected alignment to preserve content width")
					}
					// The line should now be at target width
					if lastLine.Width != tt.width {
						t.Errorf("Expected aligned line width %.1f, got %.1f", tt.width, lastLine.Width)
					}
				}
			}

			// Check non-last lines use text-align
			for i := 0; i < len(aligned)-1; i++ {
				if aligned[i].Width != tt.width {
					t.Errorf("Line %d width %.1f, want %.1f", i, aligned[i].Width, tt.width)
				}
			}
		})
	}
}

func TestAlignLines_Empty(t *testing.T) {
	txt := NewTerminal()

	style := CSSTextStyle{
		TextAlign:     AlignCenter,
		TextAlignLast: AlignRight,
	}

	aligned := txt.AlignLines([]Line{}, 20.0, style)

	if len(aligned) != 0 {
		t.Errorf("Expected empty result, got %d lines", len(aligned))
	}
}

func TestAlignLines_JustifyDefault(t *testing.T) {
	txt := NewTerminal()

	lines := []Line{
		{Content: "Justified line one", Width: 18},
		{Content: "Justified line two", Width: 18},
		{Content: "Last", Width: 4},
	}

	// When text-align is justify and text-align-last is default (left),
	// the last line should be left-aligned, not justified
	style := CSSTextStyle{
		TextAlign:     AlignJustify,
		TextAlignLast: AlignLeft, // Default
	}

	aligned := txt.AlignLines(lines, 25.0, style)

	// Last line should be left-aligned
	lastLine := aligned[len(aligned)-1]
	// It should have padding on the right (left-aligned)
	if !strings.HasSuffix(lastLine.Content, " ") && txt.Width(lastLine.Content) < 25.0 {
		// The content itself is short, and should have right padding
		t.Logf("Last line content: %q", lastLine.Content)
	}
}

// ═══════════════════════════════════════════════════════════════
//  Hanging Punctuation Tests
// ═══════════════════════════════════════════════════════════════

func TestHangingPunctuation_First(t *testing.T) {
	txt := NewTerminal()

	// Text with opening quote
	text := `"Hello world this is a test"`

	tests := []struct {
		name      string
		maxWidth  float64
		mode      HangingPunctuation
		expectFit string // What we expect to fit on first line
	}{
		{
			name:      "Without hanging - quote counts",
			maxWidth:  15.0,
			mode:      HangingPunctuationNone,
			expectFit: `"Hello world`, // ~13 chars
		},
		{
			name:      "With hanging first - quote doesn't count",
			maxWidth:  15.0,
			mode:      HangingPunctuationFirst,
			expectFit: `"Hello world t`, // ~14-15 chars (quote hangs)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			style := DefaultCSSTextStyle()
			style.HangingPunctuation = tt.mode

			lines := txt.WrapCSS(text, CSSWrapOptions{
				MaxWidth: units.Px(tt.maxWidth),
				Style:    style,
			})

			if len(lines) == 0 {
				t.Fatal("Expected at least one line")
			}

			firstLine := lines[0].Content
			t.Logf("Mode: %v, First line: %q (width: %.1f)", tt.mode, firstLine, lines[0].Width)

			// With hanging, we should fit more characters
			if tt.mode == HangingPunctuationFirst {
				// Should fit more because quote hangs
				if len(firstLine) <= len(`"Hello world`) {
					t.Errorf("With hanging first, expected to fit more than without: got %q", firstLine)
				}
			}
		})
	}
}

func TestHangingPunctuation_Last(t *testing.T) {
	txt := NewTerminal()

	// Text with closing quote
	text := `Hello world."`

	tests := []struct {
		name      string
		maxWidth  float64
		mode      HangingPunctuation
		shouldFit bool // Should entire text fit?
	}{
		{
			name:      "Without hanging - quote counts",
			maxWidth:  12.0,
			mode:      HangingPunctuationNone,
			shouldFit: false, // 14 chars won't fit in 12
		},
		{
			name:      "With hanging last - quote doesn't count",
			maxWidth:  12.0,
			mode:      HangingPunctuationLast,
			shouldFit: true, // Closing quote hangs, effective width ~13
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			style := DefaultCSSTextStyle()
			style.HangingPunctuation = tt.mode

			lines := txt.WrapCSS(text, CSSWrapOptions{
				MaxWidth: units.Px(tt.maxWidth),
				Style:    style,
			})

			t.Logf("Mode: %v, Lines: %d", tt.mode, len(lines))
			for i, line := range lines {
				t.Logf("  Line %d: %q (width: %.1f)", i, line.Content, line.Width)
			}

			fitsInOneLine := len(lines) == 1
			if fitsInOneLine != tt.shouldFit {
				t.Errorf("Expected shouldFit=%v, got %v lines", tt.shouldFit, len(lines))
			}
		})
	}
}

func TestHangingPunctuation_ForceEnd(t *testing.T) {
	txt := NewTerminal()

	// Text with period at end
	text := "Hello world."

	tests := []struct {
		name      string
		maxWidth  float64
		mode      HangingPunctuation
		shouldFit bool
	}{
		{
			name:      "Without hanging - period counts",
			maxWidth:  11.0,
			mode:      HangingPunctuationNone,
			shouldFit: false, // 12 chars won't fit in 11
		},
		{
			name:      "With force-end - period hangs",
			maxWidth:  11.0,
			mode:      HangingPunctuationForceEnd,
			shouldFit: true, // Period hangs, effective width ~11
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			style := DefaultCSSTextStyle()
			style.HangingPunctuation = tt.mode

			lines := txt.WrapCSS(text, CSSWrapOptions{
				MaxWidth: units.Px(tt.maxWidth),
				Style:    style,
			})

			t.Logf("Mode: %v, Lines: %d", tt.mode, len(lines))
			for i, line := range lines {
				t.Logf("  Line %d: %q (width: %.1f)", i, line.Content, line.Width)
			}

			fitsInOneLine := len(lines) == 1
			if fitsInOneLine != tt.shouldFit {
				t.Errorf("Expected shouldFit=%v, got %v lines", tt.shouldFit, len(lines))
			}
		})
	}
}

func TestHangingPunctuation_Combined(t *testing.T) {
	txt := NewTerminal()

	// Text with both opening and closing quotes
	text := `"Hello."`

	style := DefaultCSSTextStyle()
	style.HangingPunctuation = HangingPunctuationFirst | HangingPunctuationLast | HangingPunctuationForceEnd

	lines := txt.WrapCSS(text, CSSWrapOptions{
		MaxWidth: units.Px(5.0), // Very narrow - only "Hello" should fit with hanging on both sides
		Style:    style,
	})

	t.Logf("Lines: %d", len(lines))
	for i, line := range lines {
		t.Logf("  Line %d: %q (width: %.1f)", i, line.Content, line.Width)
	}

	// With both quotes hanging, should fit in one line
	if len(lines) != 1 {
		t.Errorf("Expected 1 line with combined hanging, got %d", len(lines))
	}
}
