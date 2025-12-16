package text

import (
	"strings"
	"testing"
)

// ═══════════════════════════════════════════════════════════════
//  Text Justify Tests
// ═══════════════════════════════════════════════════════════════

func TestJustifyText_InterWord(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name        string
		text        string
		targetWidth float64
		method      TextJustify
	}{
		{
			name:        "Basic inter-word justification",
			text:        "Hello world",
			targetWidth: 15.0,
			method:      TextJustifyInterWord,
		},
		{
			name:        "Multiple words",
			text:        "The quick brown fox",
			targetWidth: 25.0,
			method:      TextJustifyInterWord,
		},
		{
			name:        "Auto defaults to inter-word",
			text:        "Text here",
			targetWidth: 15.0,
			method:      TextJustifyAuto,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			justified := txt.JustifyText(tt.text, tt.targetWidth, tt.method)
			width := txt.Width(justified)

			// Justified text should be wider than original
			originalWidth := txt.Width(tt.text)
			if width <= originalWidth {
				t.Errorf("Justified width %.1f not greater than original %.1f", width, originalWidth)
			}

			// Should be approximately target width (allowing for rounding)
			if width < tt.targetWidth-0.5 || width > tt.targetWidth+0.5 {
				t.Errorf("Justified width %.1f not close to target %.1f", width, tt.targetWidth)
			}
		})
	}
}

func TestJustifyText_InterCharacter(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name        string
		text        string
		targetWidth float64
	}{
		{
			name:        "CJK text",
			text:        "世界和平",
			targetWidth: 12.0,
		},
		{
			name:        "Mixed CJK and Latin",
			text:        "Hello世界",
			targetWidth: 15.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			justified := txt.JustifyText(tt.text, tt.targetWidth, TextJustifyInterCharacter)
			width := txt.Width(justified)

			originalWidth := txt.Width(tt.text)
			if width <= originalWidth {
				t.Errorf("Justified width %.1f not greater than original %.1f", width, originalWidth)
			}
		})
	}
}

func TestJustifyText_None(t *testing.T) {
	txt := NewTerminal()

	text := "Hello world"
	justified := txt.JustifyText(text, 20.0, TextJustifyNone)

	if justified != text {
		t.Errorf("JustifyNone should return unchanged text, got %q", justified)
	}
}

func TestJustifyText_NoExtraSpace(t *testing.T) {
	txt := NewTerminal()

	text := "Hello"
	targetWidth := txt.Width(text)

	justified := txt.JustifyText(text, targetWidth, TextJustifyInterWord)

	// Should return unchanged if already at target width
	if justified != text {
		t.Errorf("Text at target width should be unchanged, got %q", justified)
	}
}

func TestJustifyText_SingleWord(t *testing.T) {
	txt := NewTerminal()

	text := "Hello"
	justified := txt.JustifyText(text, 20.0, TextJustifyInterWord)

	// Single word can't be justified (no gaps)
	if justified != text {
		t.Errorf("Single word should be unchanged, got %q", justified)
	}
}

// ═══════════════════════════════════════════════════════════════
//  Hanging Punctuation Tests
// ═══════════════════════════════════════════════════════════════

func TestShouldHang_First(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name     string
		text     string
		position int
		mode     HangingPunctuation
		want     bool
	}{
		{
			name:     "Opening quote at start",
			text:     `"Hello"`,
			position: 0,
			mode:     HangingPunctuationFirst,
			want:     true,
		},
		{
			name:     "Opening paren at start",
			text:     "(Hello)",
			position: 0,
			mode:     HangingPunctuationFirst,
			want:     true,
		},
		{
			name:     "CJK opening bracket",
			text:     "「こんにちは」",
			position: 0,
			mode:     HangingPunctuationFirst,
			want:     true,
		},
		{
			name:     "Not at start",
			text:     `Hello "world"`,
			position: 6,
			mode:     HangingPunctuationFirst,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shouldHang, width := txt.ShouldHang(tt.text, tt.position, tt.mode)

			if shouldHang != tt.want {
				t.Errorf("ShouldHang() = %v, want %v", shouldHang, tt.want)
			}

			if shouldHang && width <= 0 {
				t.Error("Hanging punctuation should have positive width")
			}
		})
	}
}

func TestShouldHang_Last(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name     string
		text     string
		mode     HangingPunctuation
		want     bool
	}{
		{
			name: "Closing quote at end",
			text: `"Hello"`,
			mode: HangingPunctuationLast,
			want: true,
		},
		{
			name: "Closing paren at end",
			text: "(Hello)",
			mode: HangingPunctuationLast,
			want: true,
		},
		{
			name: "CJK closing bracket",
			text: "「こんにちは」",
			mode: HangingPunctuationLast,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runes := []rune(tt.text)
			position := len(runes) - 1

			shouldHang, width := txt.ShouldHang(tt.text, position, tt.mode)

			if shouldHang != tt.want {
				t.Errorf("ShouldHang() = %v, want %v", shouldHang, tt.want)
			}

			if shouldHang && width <= 0 {
				t.Error("Hanging punctuation should have positive width")
			}
		})
	}
}

func TestShouldHang_ForceEnd(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name string
		text string
		want bool
	}{
		{name: "Period at end", text: "Hello.", want: true},
		{name: "Comma at end", text: "Hello,", want: true},
		{name: "Question at end", text: "Hello?", want: true},
		{name: "Exclamation at end", text: "Hello!", want: true},
		{name: "CJK stop at end", text: "こんにちは。", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runes := []rune(tt.text)
			position := len(runes) - 1

			shouldHang, _ := txt.ShouldHang(tt.text, position, HangingPunctuationForceEnd)

			if shouldHang != tt.want {
				t.Errorf("ShouldHang() = %v, want %v", shouldHang, tt.want)
			}
		})
	}
}

func TestShouldHang_None(t *testing.T) {
	txt := NewTerminal()

	text := `"Hello."`
	shouldHang, _ := txt.ShouldHang(text, 0, HangingPunctuationNone)

	if shouldHang {
		t.Error("HangingPunctuationNone should never hang")
	}
}

func TestShouldHang_Combined(t *testing.T) {
	txt := NewTerminal()

	text := `"Hello."`
	mode := HangingPunctuationFirst | HangingPunctuationForceEnd | HangingPunctuationLast

	// Check first position (opening quote)
	shouldHang, _ := txt.ShouldHang(text, 0, mode)
	if !shouldHang {
		t.Error("Combined mode should hang opening quote at start")
	}

	// Check last position (closing quote)
	runes := []rune(text)
	shouldHang, _ = txt.ShouldHang(text, len(runes)-1, mode)
	if !shouldHang {
		t.Error("Combined mode should hang closing quote at end")
	}

	// Check period position (second to last)
	text2 := "Hello."
	runes2 := []rune(text2)
	shouldHang, _ = txt.ShouldHang(text2, len(runes2)-1, mode)
	if !shouldHang {
		t.Error("Combined mode should hang period at end with ForceEnd")
	}
}

func TestIsOpeningPunctuation(t *testing.T) {
	tests := []struct {
		r    rune
		want bool
	}{
		{'(', true},
		{'[', true},
		{'{', true},
		{'"', true},
		{'「', true},
		{')', false},
		{'a', false},
	}

	for _, tt := range tests {
		if got := IsOpeningPunctuation(tt.r); got != tt.want {
			t.Errorf("IsOpeningPunctuation(%c) = %v, want %v", tt.r, got, tt.want)
		}
	}
}

func TestIsClosingPunctuation(t *testing.T) {
	tests := []struct {
		r    rune
		want bool
	}{
		{')', true},
		{']', true},
		{'}', true},
		{'"', true},
		{'」', true},
		{'(', false},
		{'a', false},
	}

	for _, tt := range tests {
		if got := IsClosingPunctuation(tt.r); got != tt.want {
			t.Errorf("IsClosingPunctuation(%c) = %v, want %v", tt.r, got, tt.want)
		}
	}
}

func TestIsStopPunctuation(t *testing.T) {
	tests := []struct {
		r    rune
		want bool
	}{
		{'.', true},
		{',', true},
		{'!', true},
		{'?', true},
		{'。', true},
		{'、', true},
		{'a', false},
		{'(', false},
	}

	for _, tt := range tests {
		if got := IsStopPunctuation(tt.r); got != tt.want {
			t.Errorf("IsStopPunctuation(%c) = %v, want %v", tt.r, got, tt.want)
		}
	}
}

// ═══════════════════════════════════════════════════════════════
//  Tab Size Tests
// ═══════════════════════════════════════════════════════════════

func TestExpandTabs_Spaces(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name    string
		text    string
		tabSize TabSize
	}{
		{
			name: "Single tab, default size",
			text: "Hello\tworld",
			tabSize: TabSize{
				Value: 8,
				Unit:  TabSizeSpaces,
			},
		},
		{
			name: "Multiple tabs",
			text: "A\tB\tC",
			tabSize: TabSize{
				Value: 4,
				Unit:  TabSizeSpaces,
			},
		},
		{
			name: "Tab at start",
			text: "\tIndented",
			tabSize: TabSize{
				Value: 4,
				Unit:  TabSizeSpaces,
			},
		},
		{
			name: "Tab at end",
			text: "Text\t",
			tabSize: TabSize{
				Value: 8,
				Unit:  TabSizeSpaces,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expanded := txt.ExpandTabs(tt.text, tt.tabSize)

			// Should not contain tabs
			if strings.Contains(expanded, "\t") {
				t.Error("Expanded text should not contain tabs")
			}

			// Should contain spaces
			if !strings.Contains(expanded, " ") {
				t.Error("Expanded text should contain spaces")
			}

			// Should be wider than original (tabs replaced with spaces)
			if len(expanded) <= len(tt.text) {
				t.Errorf("Expanded text length %d not greater than original %d", len(expanded), len(tt.text))
			}
		})
	}
}

func TestExpandTabs_NoTabs(t *testing.T) {
	txt := NewTerminal()

	text := "No tabs here"
	tabSize := DefaultTabSize()
	expanded := txt.ExpandTabs(text, tabSize)

	if expanded != text {
		t.Errorf("Text without tabs should be unchanged, got %q", expanded)
	}
}

func TestExpandTabs_Multiline(t *testing.T) {
	txt := NewTerminal()

	text := "Line1\tA\nLine2\tB"
	tabSize := TabSize{Value: 4, Unit: TabSizeSpaces}
	expanded := txt.ExpandTabs(text, tabSize)

	lines := strings.Split(expanded, "\n")
	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(lines))
	}

	for i, line := range lines {
		if strings.Contains(line, "\t") {
			t.Errorf("Line %d still contains tabs", i)
		}
	}
}

func TestDefaultTabSize(t *testing.T) {
	tabSize := DefaultTabSize()

	if tabSize.Value != 8 {
		t.Errorf("Default tab size value = %.1f, want 8", tabSize.Value)
	}

	if tabSize.Unit != TabSizeSpaces {
		t.Errorf("Default tab size unit = %v, want TabSizeSpaces", tabSize.Unit)
	}
}

// ═══════════════════════════════════════════════════════════════
//  Text Wrap Tests
// ═══════════════════════════════════════════════════════════════

func TestWrapBalanced(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name     string
		text     string
		maxWidth float64
	}{
		{
			name:     "Short heading",
			text:     "The Quick Brown Fox",
			maxWidth: 15.0,
		},
		{
			name:     "Medium text",
			text:     "This is a test of balanced wrapping for better readability",
			maxWidth: 20.0,
		},
		{
			name:     "Single word fits",
			text:     "Hello",
			maxWidth: 10.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := txt.WrapBalanced(tt.text, tt.maxWidth)

			if len(lines) == 0 {
				t.Error("WrapBalanced returned no lines")
			}

			// Check all lines respect max width
			for i, line := range lines {
				if line.Width > tt.maxWidth {
					t.Errorf("Line %d width %.1f exceeds max %.1f", i, line.Width, tt.maxWidth)
				}
			}

			// Check lines are relatively balanced (if multiple lines)
			if len(lines) > 1 {
				minWidth := lines[0].Width
				maxLineWidth := lines[0].Width

				for _, line := range lines[1:] {
					if line.Width < minWidth {
						minWidth = line.Width
					}
					if line.Width > maxLineWidth {
						maxLineWidth = line.Width
					}
				}

				// Lines should be within 50% of each other (balanced)
				ratio := minWidth / maxLineWidth
				if ratio < 0.5 {
					t.Logf("Lines not well balanced: min=%.1f, max=%.1f, ratio=%.2f", minWidth, maxLineWidth, ratio)
				}
			}
		})
	}
}

func TestWrapPretty(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name     string
		text     string
		maxWidth float64
	}{
		{
			name:     "Avoid orphans",
			text:     "This is a test of pretty wrapping mode",
			maxWidth: 20.0,
		},
		{
			name:     "Short paragraph",
			text:     "The quick brown fox jumps over the lazy dog",
			maxWidth: 25.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := txt.WrapPretty(tt.text, tt.maxWidth)

			if len(lines) == 0 {
				t.Error("WrapPretty returned no lines")
			}

			// Check all lines respect max width
			for i, line := range lines {
				if line.Width > tt.maxWidth {
					t.Errorf("Line %d width %.1f exceeds max %.1f", i, line.Width, tt.maxWidth)
				}
			}

			// If there are multiple lines, last line shouldn't be too short
			if len(lines) >= 2 {
				lastLine := lines[len(lines)-1]
				// After pretty wrapping, last line should be reasonable length
				// (at least 20% of max width if possible)
				if lastLine.Width < tt.maxWidth*0.2 {
					t.Logf("Last line might be too short: %.1f", lastLine.Width)
				}
			}
		})
	}
}

// ═══════════════════════════════════════════════════════════════
//  Text Spacing Trim Tests
// ═══════════════════════════════════════════════════════════════

func TestTrimCJKSpacing_SpaceAll(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name string
		text string
	}{
		{
			name: "Chinese with spaces",
			text: "世界 和平",
		},
		{
			name: "Japanese with spaces",
			text: "日本 語",
		},
		{
			name: "Korean with spaces",
			text: "안녕 하세요",
		},
		{
			name: "Multiple spaces",
			text: "中国 人民 共和国",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trimmed := txt.TrimCJKSpacing(tt.text, TextSpacingTrimSpaceAll)

			// Should have fewer or equal spaces
			originalSpaces := strings.Count(tt.text, " ")
			trimmedSpaces := strings.Count(trimmed, " ")

			if trimmedSpaces > originalSpaces {
				t.Errorf("Trimmed text has more spaces (%d) than original (%d)", trimmedSpaces, originalSpaces)
			}

			// Should be shorter or equal length
			if len(trimmed) > len(tt.text) {
				t.Errorf("Trimmed text longer than original")
			}
		})
	}
}

func TestTrimCJKSpacing_Mixed(t *testing.T) {
	txt := NewTerminal()

	// Mixed CJK and Latin - only trim spaces between CJK
	text := "Hello 世界 and 和平"
	trimmed := txt.TrimCJKSpacing(text, TextSpacingTrimSpaceAll)

	// Space between "Hello" and "世界" should remain
	// Space between "世界" and "and" should remain
	// Space between "and" and "和平" should remain
	// Only space between CJK characters would be trimmed (none in this case)

	if !strings.Contains(trimmed, "Hello") {
		t.Error("Should preserve Latin text")
	}
}

func TestTrimCJKSpacing_SpaceFirst(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name string
		text string
	}{
		{
			name: "Space at start before CJK",
			text: " 世界",
		},
		{
			name: "Space after newline before CJK",
			text: "Line1\n 世界",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trimmed := txt.TrimCJKSpacing(tt.text, TextSpacingTrimSpaceFirst)

			// Should have fewer or equal spaces
			if len(trimmed) > len(tt.text) {
				t.Errorf("Trimmed text longer than original")
			}
		})
	}
}

func TestTrimCJKSpacing_None(t *testing.T) {
	txt := NewTerminal()

	text := "世界 和平"
	trimmed := txt.TrimCJKSpacing(text, TextSpacingTrimNone)

	if trimmed != text {
		t.Errorf("TrimNone should return unchanged text, got %q", trimmed)
	}
}

func TestIsCJKIdeograph(t *testing.T) {
	tests := []struct {
		name string
		r    rune
		want bool
	}{
		{"Chinese common", '中', true},
		{"Chinese common 2", '国', true},
		{"Japanese kanji", '日', true},
		{"Korean hanja", '韓', true},
		{"CJK Extension A", '\u3400', true},
		{"Latin", 'a', false},
		{"Hiragana", 'あ', false},
		{"Katakana", 'ア', false},
		{"Hangul", '가', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCJKIdeograph(tt.r); got != tt.want {
				t.Errorf("IsCJKIdeograph(%c/U+%04X) = %v, want %v", tt.r, tt.r, got, tt.want)
			}
		})
	}
}

// ═══════════════════════════════════════════════════════════════
//  Wrap Controls Tests
// ═══════════════════════════════════════════════════════════════

func TestWrapWithControls(t *testing.T) {
	txt := NewTerminal()

	text := "The quick brown fox jumps over the lazy dog"
	controls := []WrapPoint{
		{Position: 10, Before: WrapControlAvoid},
		{Position: 20, After: WrapControlAlways},
	}

	lines := txt.WrapWithControls(text, 25.0, controls)

	if len(lines) == 0 {
		t.Error("WrapWithControls returned no lines")
	}

	// Basic check: all lines should respect max width
	for i, line := range lines {
		if line.Width > 25.0 {
			t.Errorf("Line %d width %.1f exceeds max 25.0", i, line.Width)
		}
	}
}

func TestWrapWithControls_Empty(t *testing.T) {
	txt := NewTerminal()

	text := "Hello world"
	lines := txt.WrapWithControls(text, 20.0, nil)

	if len(lines) == 0 {
		t.Error("WrapWithControls with no controls returned no lines")
	}
}

// ═══════════════════════════════════════════════════════════════
//  Word Break: Auto-Phrase Tests
// ═══════════════════════════════════════════════════════════════

// MockPhraseBreaker is a simple test implementation of PhraseBreaker.
type MockPhraseBreaker struct {
	boundaries []int
}

func (m *MockPhraseBreaker) FindPhrases(text string) []int {
	return m.boundaries
}

func TestWrapWithPhrases(t *testing.T) {
	txt := NewTerminal()

	// Simulate CJK text: "你好世界这是测试" (12 chars)
	// We'll use ASCII for simplicity but pretend it's segmented
	text := "HelloWorldTest"

	// Mock breaker that segments at: Hello|World|Test
	breaker := &MockPhraseBreaker{
		boundaries: []int{0, 5, 10, 14}, // After "Hello", "World", "Test"
	}

	lines := txt.WrapWithPhrases(text, 10.0, breaker)

	if len(lines) == 0 {
		t.Fatal("WrapWithPhrases returned no lines")
	}

	t.Logf("Lines: %d", len(lines))
	for i, line := range lines {
		t.Logf("  Line %d: %q (width: %.1f)", i, line.Content, line.Width)
	}

	// With phrase boundaries at 5, 10, 14 and maxWidth=10,
	// we should get phrases respecting those boundaries
	// Each phrase should fit within the width
	for i, line := range lines {
		if line.Width > 10.0 {
			t.Errorf("Line %d width %.1f exceeds max 10.0", i, line.Width)
		}
	}
}

func TestWrapWithPhrases_NilBreaker(t *testing.T) {
	txt := NewTerminal()

	text := "Hello world"
	lines := txt.WrapWithPhrases(text, 20.0, nil)

	if len(lines) == 0 {
		t.Error("WrapWithPhrases with nil breaker returned no lines")
	}
}

func TestWrapWithPhrasesAndControls(t *testing.T) {
	txt := NewTerminal()

	text := "HelloWorldTest"

	breaker := &MockPhraseBreaker{
		boundaries: []int{0, 5, 10, 14},
	}

	controls := []WrapPoint{
		{Position: 5, After: WrapControlAvoid}, // Don't break after "Hello"
	}

	lines := txt.WrapWithPhrasesAndControls(text, 15.0, breaker, controls)

	if len(lines) == 0 {
		t.Fatal("WrapWithPhrasesAndControls returned no lines")
	}

	t.Logf("Lines: %d", len(lines))
	for i, line := range lines {
		t.Logf("  Line %d: %q (width: %.1f)", i, line.Content, line.Width)
	}

	// The control should prevent breaking after position 5
	// So "HelloWorld" might stay together if it fits
	for i, line := range lines {
		if line.Width > 15.0 {
			t.Errorf("Line %d width %.1f exceeds max 15.0", i, line.Width)
		}
	}
}

// ═══════════════════════════════════════════════════════════════
//  Benchmark Tests
// ═══════════════════════════════════════════════════════════════

func BenchmarkJustifyText_InterWord(b *testing.B) {
	txt := NewTerminal()
	text := "The quick brown fox jumps over the lazy dog"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txt.JustifyText(text, 60.0, TextJustifyInterWord)
	}
}

func BenchmarkJustifyText_InterCharacter(b *testing.B) {
	txt := NewTerminal()
	text := "世界和平与发展"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txt.JustifyText(text, 20.0, TextJustifyInterCharacter)
	}
}

func BenchmarkExpandTabs(b *testing.B) {
	txt := NewTerminal()
	text := "Line1\tA\tB\tC\nLine2\tD\tE\tF"
	tabSize := DefaultTabSize()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txt.ExpandTabs(text, tabSize)
	}
}

func BenchmarkWrapBalanced(b *testing.B) {
	txt := NewTerminal()
	text := "This is a longer piece of text that will be wrapped across multiple lines to test the balanced wrapping algorithm performance"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txt.WrapBalanced(text, 30.0)
	}
}

func BenchmarkWrapPretty(b *testing.B) {
	txt := NewTerminal()
	text := "This is a longer piece of text that will be wrapped across multiple lines to test the pretty wrapping algorithm performance"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txt.WrapPretty(text, 30.0)
	}
}

func BenchmarkTrimCJKSpacing(b *testing.B) {
	txt := NewTerminal()
	text := "中国 人民 共和国 的 发展 与 和平"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txt.TrimCJKSpacing(text, TextSpacingTrimSpaceAll)
	}
}

// ═══════════════════════════════════════════════════════════════
//  Text Autospace Tests
// ═══════════════════════════════════════════════════════════════

func TestApplyAutospace_IdeographAlpha(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "Chinese with English",
			text:     "Hello世界",
			expected: "Hello 世界",
		},
		{
			name:     "English with Chinese",
			text:     "世界Hello",
			expected: "世界 Hello",
		},
		{
			name:     "Multiple transitions",
			text:     "Hello世界World中国",
			expected: "Hello 世界 World 中国",
		},
		{
			name:     "Japanese Hiragana with English",
			text:     "Helloこんにちは",
			expected: "Hello こんにちは",
		},
		{
			name:     "Japanese Katakana with English",
			text:     "テストTest",
			expected: "テスト Test",
		},
		{
			name:     "Korean with English",
			text:     "안녕Hello",
			expected: "안녕 Hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := txt.ApplyAutospace(tt.text, AutospaceIdeographAlpha)

			if result != tt.expected {
				t.Errorf("ApplyAutospace(%q) = %q, want %q", tt.text, result, tt.expected)
			}
		})
	}
}

func TestApplyAutospace_IdeographNumeric(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "Chinese with number",
			text:     "世界123",
			expected: "世界 123",
		},
		{
			name:     "Number with Chinese",
			text:     "123世界",
			expected: "123 世界",
		},
		{
			name:     "Multiple numbers",
			text:     "世界123中国456",
			expected: "世界 123 中国 456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := txt.ApplyAutospace(tt.text, AutospaceIdeographNumeric)

			if result != tt.expected {
				t.Errorf("ApplyAutospace(%q) = %q, want %q", tt.text, result, tt.expected)
			}
		})
	}
}

func TestApplyAutospace_Punctuation(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name string
		text string
	}{
		{
			name: "Opening bracket",
			text: "「世界」",
		},
		{
			name: "Parentheses",
			text: "（テスト）",
		},
		{
			name: "Mixed punctuation",
			text: "《中国》",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := txt.ApplyAutospace(tt.text, AutospacePunctuation)

			// Just verify it doesn't crash
			if len(result) == 0 && len(tt.text) > 0 {
				t.Errorf("ApplyAutospace returned empty for non-empty input")
			}
		})
	}
}

func TestApplyAutospace_All(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name string
		text string
	}{
		{
			name: "Complex mixed text",
			text: "Hello世界123テスト",
		},
		{
			name: "With punctuation",
			text: "「世界」test123",
		},
		{
			name: "Multiple scripts",
			text: "English中文한글テスト123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := txt.ApplyAutospace(tt.text, AutospaceAll)

			// Verify result is longer (due to added spaces)
			if len(result) < len(tt.text) {
				t.Errorf("ApplyAutospace should not shorten text")
			}

			t.Logf("Input:  %q", tt.text)
			t.Logf("Output: %q", result)
		})
	}
}

func TestApplyAutospace_None(t *testing.T) {
	txt := NewTerminal()

	text := "Hello世界123"
	result := txt.ApplyAutospace(text, AutospaceNone)

	if result != text {
		t.Errorf("AutospaceNone should return unchanged text, got %q", result)
	}
}

func TestApplyAutospaceMode(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name string
		text string
		mode TextAutospace
	}{
		{
			name: "Normal mode",
			text: "Hello世界",
			mode: TextAutospaceNormal,
		},
		{
			name: "Auto mode",
			text: "Hello世界",
			mode: TextAutospaceAuto,
		},
		{
			name: "No autospace",
			text: "Hello世界",
			mode: TextAutospaceNoAutospace,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := txt.ApplyAutospaceMode(tt.text, tt.mode)

			if tt.mode == TextAutospaceNoAutospace {
				if result != tt.text {
					t.Errorf("NoAutospace should return unchanged, got %q", result)
				}
			} else {
				// Normal and Auto should add spacing
				if len(result) <= len(tt.text) {
					t.Logf("Note: No spacing added for %q with mode %v", tt.text, tt.mode)
				}
			}

			t.Logf("Mode %v: %q -> %q", tt.mode, tt.text, result)
		})
	}
}

func TestIsIdeographic(t *testing.T) {
	tests := []struct {
		name string
		r    rune
		want bool
	}{
		// CJK Ideographs
		{"Chinese common", '中', true},
		{"Chinese common 2", '国', true},

		// Hiragana
		{"Hiragana あ", 'あ', true},
		{"Hiragana か", 'か', true},

		// Katakana
		{"Katakana ア", 'ア', true},
		{"Katakana カ", 'カ', true},

		// Hangul
		{"Hangul 가", '가', true},
		{"Hangul 한", '한', true},

		// Non-ideographic
		{"Latin A", 'A', false},
		{"Latin a", 'a', false},
		{"Digit 1", '1', false},
		{"Space", ' ', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIdeographic(tt.r); got != tt.want {
				t.Errorf("IsIdeographic(%c/U+%04X) = %v, want %v", tt.r, tt.r, got, tt.want)
			}
		})
	}
}

func TestIsFullwidthPunctuation(t *testing.T) {
	tests := []struct {
		r    rune
		want bool
	}{
		{'、', true},
		{'。', true},
		{'，', true},
		{'「', true},
		{'」', true},
		{'（', true},
		{'）', true},
		{'.', false},
		{',', false},
		{'a', false},
	}

	for _, tt := range tests {
		if got := IsFullwidthPunctuation(tt.r); got != tt.want {
			t.Errorf("IsFullwidthPunctuation(%c) = %v, want %v", tt.r, got, tt.want)
		}
	}
}

func BenchmarkApplyAutospace(b *testing.B) {
	txt := NewTerminal()
	text := "Hello世界123テストKorean한글"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txt.ApplyAutospace(text, AutospaceAll)
	}
}
