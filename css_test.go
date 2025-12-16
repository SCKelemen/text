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
