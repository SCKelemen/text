package text

import (
	"testing"
)

// Comprehensive International Text Support Tests
//
// Tests for:
// - LTR (Left-to-Right) text
// - RTL (Right-to-Left) text (Arabic, Hebrew)
// - Bidirectional text (mixed LTR/RTL)
// - Vertical text layout
// - CJK characters (Chinese, Japanese, Korean)
// - Emoji (including modifiers and sequences)
// - Various scripts

// ═══════════════════════════════════════════════════════════════
//  LTR Text Tests
// ═══════════════════════════════════════════════════════════════

func TestLTR_Basic(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name     string
		text     string
		expected float64
	}{
		{"English", "Hello World", 11.0},
		{"Spanish", "Hola Mundo", 10.0},
		{"French", "Bonjour Monde", 13.0},
		{"German", "Hallo Welt", 10.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width := txt.Width(tt.text)
			if width != tt.expected {
				t.Errorf("Width(%q) = %.1f, want %.1f", tt.text, width, tt.expected)
			}
		})
	}
}

// ═══════════════════════════════════════════════════════════════
//  RTL Text Tests (Arabic, Hebrew)
// ═══════════════════════════════════════════════════════════════

func TestRTL_Arabic(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name string
		text string
	}{
		{"Arabic hello", "مرحبا"},
		{"Arabic sentence", "مرحبا بك في العالم"},
		{"Arabic with numbers", "العدد 123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width := txt.Width(tt.text)
			if width <= 0 {
				t.Errorf("Width(%q) = %.1f, want > 0", tt.text, width)
			}

			// Test reordering
			reordered := txt.Reorder(tt.text)
			if reordered == "" {
				t.Errorf("Reorder(%q) returned empty string", tt.text)
			}
		})
	}
}

func TestRTL_Hebrew(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name string
		text string
	}{
		{"Hebrew hello", "שלום"},
		{"Hebrew sentence", "שלום עולם"},
		{"Hebrew with numbers", "מספר 456"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width := txt.Width(tt.text)
			if width <= 0 {
				t.Errorf("Width(%q) = %.1f, want > 0", tt.text, width)
			}

			// Test direction detection
			dir := txt.DetectDirection(tt.text)
			if dir == 0 {
				t.Errorf("DetectDirection(%q) failed", tt.text)
			}
		})
	}
}

// ═══════════════════════════════════════════════════════════════
//  Bidirectional Text Tests
// ═══════════════════════════════════════════════════════════════

func TestBidirectional_Mixed(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name string
		text string
	}{
		{"English + Arabic", "Hello مرحبا"},
		{"English + Hebrew", "World שלום"},
		{"Numbers in RTL", "The number is 123 in مرحبا"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width := txt.Width(tt.text)
			if width <= 0 {
				t.Errorf("Width(%q) = %.1f, want > 0", tt.text, width)
			}

			// Test reordering doesn't crash
			reordered := txt.Reorder(tt.text)
			if len(reordered) == 0 && len(tt.text) > 0 {
				t.Errorf("Reorder(%q) returned empty for non-empty input", tt.text)
			}
		})
	}
}

// ═══════════════════════════════════════════════════════════════
//  CJK Character Tests
// ═══════════════════════════════════════════════════════════════

func TestCJK_Chinese(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name     string
		text     string
		expected float64
	}{
		{"Simple Chinese", "世界", 4.0},    // 2 + 2
		{"Chinese sentence", "你好世界", 8.0}, // 2 + 2 + 2 + 2
		{"Mixed Chinese-English", "Hello世界", 9.0}, // 5 + 4
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width := txt.Width(tt.text)
			if width != tt.expected {
				t.Errorf("Width(%q) = %.1f, want %.1f", tt.text, width, tt.expected)
			}
		})
	}
}

func TestCJK_Japanese(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name string
		text string
	}{
		{"Hiragana", "こんにちは"},
		{"Katakana", "コンニチハ"},
		{"Kanji", "日本語"},
		{"Mixed Japanese", "こんにちは世界"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width := txt.Width(tt.text)
			if width <= 0 {
				t.Errorf("Width(%q) = %.1f, want > 0", tt.text, width)
			}

			// Test wrapping
			lines := txt.Wrap(tt.text, WrapOptions{MaxWidth: 6})
			if len(lines) == 0 {
				t.Errorf("Wrap(%q) returned no lines", tt.text)
			}
		})
	}
}

func TestCJK_Korean(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name string
		text string
	}{
		{"Korean hello", "안녕하세요"},
		{"Korean world", "세계"},
		{"Mixed Korean", "Hello 세계"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width := txt.Width(tt.text)
			if width <= 0 {
				t.Errorf("Width(%q) = %.1f, want > 0", tt.text, width)
			}
		})
	}
}

// ═══════════════════════════════════════════════════════════════
//  Emoji Tests
// ═══════════════════════════════════════════════════════════════

func TestEmoji_Basic(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name     string
		text     string
		expected float64
	}{
		{"Smiley", "😀", 2.0},
		{"Heart", "❤️", 2.0},
		{"Thumbs up", "👍", 2.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width := txt.Width(tt.text)
			if width != tt.expected {
				t.Errorf("Width(%q) = %.1f, want %.1f", tt.text, width, tt.expected)
			}
		})
	}
}

func TestEmoji_Modifiers(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name     string
		text     string
		expected float64
	}{
		{"Wave with skin tone", "👋🏻", 2.0}, // Base + modifier = 2
		{"Thumbs up with skin tone", "👍🏽", 2.0},
		{"Person with skin tone", "👨🏾", 2.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width := txt.Width(tt.text)
			if width != tt.expected {
				t.Errorf("Width(%q) = %.1f, want %.1f", tt.text, width, tt.expected)
			}

			// Should be single grapheme
			graphemes := txt.Graphemes(tt.text)
			if len(graphemes) != 1 {
				t.Errorf("Graphemes(%q) = %d clusters, want 1", tt.text, len(graphemes))
			}
		})
	}
}

func TestEmoji_ZWJSequences(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name string
		text string
	}{
		{"Family", "👨‍👩‍👧‍👦"},
		{"Woman technologist", "👩‍💻"},
		{"Rainbow flag", "🏳️‍🌈"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width := txt.Width(tt.text)
			if width <= 0 {
				t.Errorf("Width(%q) = %.1f, want > 0", tt.text, width)
			}

			// Should be single grapheme
			graphemes := txt.Graphemes(tt.text)
			if len(graphemes) != 1 {
				t.Errorf("Graphemes(%q) = %d clusters, want 1", tt.text, len(graphemes))
			}
		})
	}
}

func TestEmoji_Flags(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name string
		text string
	}{
		{"US flag", "🇺🇸"},
		{"UK flag", "🇬🇧"},
		{"Japan flag", "🇯🇵"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width := txt.Width(tt.text)
			if width <= 0 {
				t.Errorf("Width(%q) = %.1f, want > 0", tt.text, width)
			}

			// Regional indicator pairs should be single grapheme
			graphemes := txt.Graphemes(tt.text)
			if len(graphemes) != 1 {
				t.Errorf("Graphemes(%q) = %d clusters, want 1", tt.text, len(graphemes))
			}
		})
	}
}

// ═══════════════════════════════════════════════════════════════
//  Vertical Text Tests
// ═══════════════════════════════════════════════════════════════

func TestVertical_CJK(t *testing.T) {
	txt := NewTerminal()

	style := VerticalTextStyle{
		WritingMode:     WritingModeVerticalRL,
		TextOrientation: TextOrientationMixed,
	}

	tests := []struct {
		name string
		text string
	}{
		{"Chinese vertical", "世界"},
		{"Japanese vertical", "日本語"},
		{"Mixed vertical", "Hello世界"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics := txt.MeasureVertical(tt.text, style)

			if metrics.Advance <= 0 {
				t.Errorf("MeasureVertical(%q).Advance = %.1f, want > 0", tt.text, metrics.Advance)
			}

			if metrics.InlineSize <= 0 {
				t.Errorf("MeasureVertical(%q).InlineSize = %.1f, want > 0", tt.text, metrics.InlineSize)
			}
		})
	}
}

func TestVertical_Wrapping(t *testing.T) {
	txt := NewTerminal()

	style := VerticalTextStyle{
		WritingMode:     WritingModeVerticalRL,
		TextOrientation: TextOrientationMixed,
	}

	text := "世界こんにちは日本語"
	columns := txt.WrapVertical(text, VerticalWrapOptions{
		MaxBlockSize: 5.0,
		Style:        style,
	})

	if len(columns) == 0 {
		t.Error("WrapVertical returned no columns")
	}

	for i, col := range columns {
		if col.Advance > 5.0 {
			t.Errorf("Column %d advance %.1f exceeds max 5.0", i, col.Advance)
		}
	}
}

// ═══════════════════════════════════════════════════════════════
//  Various Scripts Tests
// ═══════════════════════════════════════════════════════════════

func TestScripts_Cyrillic(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name string
		text string
	}{
		{"Russian", "Привет мир"},
		{"Ukrainian", "Привіт світ"},
		{"Serbian", "Здраво свете"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width := txt.Width(tt.text)
			if width <= 0 {
				t.Errorf("Width(%q) = %.1f, want > 0", tt.text, width)
			}
		})
	}
}

func TestScripts_Greek(t *testing.T) {
	txt := NewTerminal()

	text := "Γειά σου κόσμε"
	width := txt.Width(text)
	if width <= 0 {
		t.Errorf("Width(%q) = %.1f, want > 0", text, width)
	}
}

func TestScripts_Thai(t *testing.T) {
	txt := NewTerminal()

	text := "สวัสดีชาวโลก"
	width := txt.Width(text)
	if width <= 0 {
		t.Errorf("Width(%q) = %.1f, want > 0", text, width)
	}

	// Test wrapping doesn't crash
	lines := txt.Wrap(text, WrapOptions{MaxWidth: 10})
	if len(lines) == 0 {
		t.Error("Wrap returned no lines for Thai text")
	}
}

func TestScripts_Devanagari(t *testing.T) {
	txt := NewTerminal()

	text := "नमस्ते दुनिया" // Hindi
	width := txt.Width(text)
	if width <= 0 {
		t.Errorf("Width(%q) = %.1f, want > 0", text, width)
	}
}

// ═══════════════════════════════════════════════════════════════
//  Complex Mixed Script Tests
// ═══════════════════════════════════════════════════════════════

func TestMixed_AllScripts(t *testing.T) {
	txt := NewTerminal()

	// Mix of English, CJK, Arabic, emoji
	text := "Hello 世界 مرحبا 😀"

	width := txt.Width(text)
	if width <= 0 {
		t.Errorf("Width(%q) = %.1f, want > 0", text, width)
	}

	// Test wrapping doesn't crash
	lines := txt.Wrap(text, WrapOptions{MaxWidth: 15})
	if len(lines) == 0 {
		t.Error("Wrap returned no lines for mixed script text")
	}

	// Test truncation doesn't crash
	truncated := txt.Truncate(text, TruncateOptions{
		MaxWidth: 10,
		Strategy: TruncateEnd,
	})
	if len(truncated) == 0 && len(text) > 0 {
		t.Error("Truncate returned empty for non-empty mixed script text")
	}
}

// ═══════════════════════════════════════════════════════════════
//  Combining Marks Tests
// ═══════════════════════════════════════════════════════════════

func TestCombiningMarks(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name           string
		text           string
		expectedGraphemes int
	}{
		{"e with acute", "é", 1},  // e + combining acute
		{"a with tilde", "ã", 1}, // a + combining tilde
		{"Complex diacritics", "ṽ", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphemes := txt.Graphemes(tt.text)
			if len(graphemes) != tt.expectedGraphemes {
				t.Errorf("Graphemes(%q) = %d, want %d", tt.text, len(graphemes), tt.expectedGraphemes)
			}
		})
	}
}
