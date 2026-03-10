package text

import (
	"testing"
)

func TestWidth(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected float64
	}{
		{"ASCII", "Hello", 5.0},
		{"CJK wide", "世界", 4.0},   // 2 + 2
		{"Mixed", "Hello世界", 9.0}, // 5 + 4
		{"Emoji", "😀", 2.0},
		{"Emoji with modifier", "👋🏻", 2.0}, // emoji + skin tone = still 2
		{"Space", " ", 1.0},
	}

	txt := NewTerminal()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := txt.Width(tt.text)
			if got != tt.expected {
				t.Errorf("Width(%q) = %.1f, want %.1f", tt.text, got, tt.expected)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name     string
		text     string
		maxWidth float64
		strategy TruncateStrategy
		want     string
	}{
		{
			name:     "No truncation needed",
			text:     "Hello",
			maxWidth: 10,
			strategy: TruncateEnd,
			want:     "Hello",
		},
		{
			name:     "Truncate end ASCII",
			text:     "Hello world",
			maxWidth: 8,
			strategy: TruncateEnd,
			want:     "Hello...",
		},
		{
			name:     "Truncate end with CJK",
			text:     "Hello世界",
			maxWidth: 8,
			strategy: TruncateEnd,
			want:     "Hello...",
		},
		{
			name:     "Truncate middle",
			text:     "Hello world",
			maxWidth: 8,
			strategy: TruncateMiddle,
			want:     "He...ld",
		},
		{
			name:     "Truncate start",
			text:     "Hello world",
			maxWidth: 8,
			strategy: TruncateStart,
			want:     "...world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := txt.Truncate(tt.text, TruncateOptions{
				MaxWidth: tt.maxWidth,
				Strategy: tt.strategy,
			})
			if got != tt.want {
				t.Errorf("Truncate() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAlign(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name  string
		text  string
		width float64
		align Alignment
		want  string
	}{
		{
			name:  "Left align",
			text:  "Hello",
			width: 10,
			align: AlignLeft,
			want:  "Hello     ",
		},
		{
			name:  "Right align",
			text:  "Hello",
			width: 10,
			align: AlignRight,
			want:  "     Hello",
		},
		{
			name:  "Center align",
			text:  "Hello",
			width: 11,
			align: AlignCenter,
			want:  "   Hello   ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := txt.Align(tt.text, tt.width, tt.align)
			if got != tt.want {
				t.Errorf("Align() = %q, want %q", got, tt.want)
			}
			// Verify width is correct
			if txt.Width(got) != tt.width {
				t.Errorf("Aligned text width = %.1f, want %.1f", txt.Width(got), tt.width)
			}
		})
	}
}

func TestWrap(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name      string
		text      string
		maxWidth  float64
		wantLines int
	}{
		{
			name:      "No wrap needed",
			text:      "Hello",
			maxWidth:  10,
			wantLines: 1,
		},
		{
			name:      "Simple wrap",
			text:      "Hello world test",
			maxWidth:  10,
			wantLines: 2,
		},
		{
			name:      "CJK wrap",
			text:      "Hello世界test",
			maxWidth:  10,
			wantLines: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := txt.Wrap(tt.text, WrapOptions{
				MaxWidth: tt.maxWidth,
			})
			if len(lines) != tt.wantLines {
				t.Errorf("Wrap() returned %d lines, want %d", len(lines), tt.wantLines)
			}
			// Verify no line exceeds maxWidth
			for i, line := range lines {
				if line.Width > tt.maxWidth {
					t.Errorf("Line %d width %.1f exceeds maxWidth %.1f: %q",
						i, line.Width, tt.maxWidth, line.Content)
				}
			}
		})
	}
}

func TestWrap_BreakWordsOption(t *testing.T) {
	txt := NewTerminal()
	text := "supercalifragilistic"

	noBreak := txt.Wrap(text, WrapOptions{
		MaxWidth:   8,
		BreakWords: false,
	})
	if len(noBreak) != 1 {
		t.Fatalf("BreakWords=false should keep unbreakable word on one line, got %d lines", len(noBreak))
	}
	if noBreak[0].Content != text {
		t.Fatalf("BreakWords=false line content = %q, want %q", noBreak[0].Content, text)
	}
	if noBreak[0].Width <= 8 {
		t.Fatalf("BreakWords=false expected overflow line width > 8, got %.1f", noBreak[0].Width)
	}

	breakWords := txt.Wrap(text, WrapOptions{
		MaxWidth:   8,
		BreakWords: true,
	})
	if len(breakWords) <= 1 {
		t.Fatalf("BreakWords=true should break long word, got %d lines", len(breakWords))
	}
	for i, line := range breakWords {
		if line.Width > 8 {
			t.Fatalf("BreakWords=true line %d width %.1f exceeds maxWidth 8", i, line.Width)
		}
	}
}

func TestWrap_PreserveNewlines(t *testing.T) {
	txt := NewTerminal()
	text := "a\n👨‍👩‍👧‍👦b"

	lines := txt.Wrap(text, WrapOptions{
		MaxWidth:         20,
		PreserveNewlines: true,
	})
	if len(lines) != 2 {
		t.Fatalf("Wrap() with PreserveNewlines returned %d lines, want 2", len(lines))
	}

	if lines[0].Content != "a" || lines[0].Start != 0 || lines[0].End != 1 {
		t.Fatalf("line 0 = %+v, want content=%q start=0 end=1", lines[0], "a")
	}
	if lines[1].Content != "👨‍👩‍👧‍👦b" || lines[1].Start != 2 || lines[1].End != 10 {
		t.Fatalf("line 1 = %+v, want content=%q start=2 end=10", lines[1], "👨‍👩‍👧‍👦b")
	}
}

func TestWrap_RuneIndicesWithGrapheme(t *testing.T) {
	txt := NewTerminal()
	text := "👨‍👩‍👧‍👦a"

	lines := txt.Wrap(text, WrapOptions{
		MaxWidth:   2,
		BreakWords: true,
	})
	if len(lines) != 2 {
		t.Fatalf("Wrap() returned %d lines, want 2", len(lines))
	}

	if lines[0].Content != "👨‍👩‍👧‍👦" || lines[0].Start != 0 || lines[0].End != 7 {
		t.Fatalf("line 0 = %+v, want content=%q start=0 end=7", lines[0], "👨‍👩‍👧‍👦")
	}
	if lines[1].Content != "a" || lines[1].Start != 7 || lines[1].End != 8 {
		t.Fatalf("line 1 = %+v, want content=%q start=7 end=8", lines[1], "a")
	}
}

func TestGraphemes(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name string
		text string
		want int
	}{
		{"ASCII", "Hello", 5},
		{"CJK", "世界", 2},
		{"Emoji", "😀", 1},
		{"Emoji with modifier", "👋🏻", 1}, // Should be 1 grapheme cluster
		{"Complex emoji", "👨‍👩‍👧‍👦", 1},  // Family emoji with ZWJ
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graphemes := txt.Graphemes(tt.text)
			if len(graphemes) != tt.want {
				t.Errorf("Graphemes(%q) = %d clusters, want %d",
					tt.text, len(graphemes), tt.want)
			}
		})
	}
}

func BenchmarkWidth(b *testing.B) {
	txt := NewTerminal()
	text := "Hello 世界! This is a test with emoji 😀 and CJK."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txt.Width(text)
	}
}

func BenchmarkTruncate(b *testing.B) {
	txt := NewTerminal()
	text := "Hello 世界! This is a long text that needs truncation with emoji 😀."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txt.Truncate(text, TruncateOptions{
			MaxWidth: 20,
			Strategy: TruncateEnd,
		})
	}
}

func BenchmarkWrap(b *testing.B) {
	txt := NewTerminal()
	text := "Hello 世界! This is a long text that needs wrapping across multiple lines with proper Unicode support including emoji 😀 and combining marks."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txt.Wrap(text, WrapOptions{
			MaxWidth: 40,
		})
	}
}

func TestWidthRange(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name     string
		text     string
		start    int
		end      int
		expected float64
	}{
		{"ASCII range", "Hello world", 0, 5, 5.0},
		{"Middle range", "Hello world", 6, 11, 5.0},
		{"Empty range", "Hello world", 5, 5, 0.0},
		{"CJK range", "Hello世界!", 5, 7, 4.0}, // 世界 = 4 cells
		{"Emoji ZWJ cluster", "👨‍👩‍👧‍👦", 0, 7, 2.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width := txt.WidthRange(tt.text, tt.start, tt.end)
			if width != tt.expected {
				t.Errorf("WidthRange(%q, %d, %d) = %.1f, want %.1f",
					tt.text, tt.start, tt.end, width, tt.expected)
			}
		})
	}
}

func TestTerminalMeasureEastAsian(t *testing.T) {
	// Test East Asian context where ambiguous characters are wide
	tests := []struct {
		name     string
		char     rune
		expected float64
	}{
		{"ASCII", 'A', 1.0},
		{"CJK", '世', 2.0},
		{"Emoji", '😀', 2.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width := TerminalMeasureEastAsian(tt.char)
			if width != tt.expected {
				t.Errorf("TerminalMeasureEastAsian(%q) = %.1f, want %.1f",
					tt.char, width, tt.expected)
			}
		})
	}
}

func TestGraphemeCount(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name     string
		text     string
		expected int
	}{
		{"ASCII", "Hello", 5},
		{"CJK", "世界", 2},
		{"Emoji with modifier", "👋🏻", 1},
		{"Mixed", "Hello👋", 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := txt.GraphemeCount(tt.text)
			if count != tt.expected {
				t.Errorf("GraphemeCount(%q) = %d, want %d",
					tt.text, count, tt.expected)
			}
		})
	}
}

func TestGraphemeAt(t *testing.T) {
	txt := NewTerminal()
	text := "Hello👋🏻"

	tests := []struct {
		name     string
		index    int
		expected string
	}{
		{"First char", 0, "H"},
		{"Middle char", 2, "l"},
		{"Emoji", 5, "👋🏻"},
		{"Out of range", 10, ""},
		{"Negative index", -1, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := txt.GraphemeAt(text, tt.index)
			if result != tt.expected {
				t.Errorf("GraphemeAt(%q, %d) = %q, want %q",
					text, tt.index, result, tt.expected)
			}
		})
	}
}

func TestReorderWithDirection(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name      string
		text      string
		direction Direction
	}{
		{"LTR explicit", "Hello world", DirectionLTR},
		{"RTL explicit", "שלום עולם", DirectionRTL},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify it doesn't panic
			result := txt.ReorderWithDirection(tt.text, toUAX9Direction(tt.direction))
			if result == "" && tt.text != "" {
				t.Errorf("ReorderWithDirection returned empty string")
			}
		})
	}
}

func TestResolveAlignment(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name        string
		align       Alignment
		direction   Direction
		parentAlign Alignment
		expected    Alignment
	}{
		{"Start LTR", AlignStart, DirectionLTR, AlignLeft, AlignLeft},
		{"Start RTL", AlignStart, DirectionRTL, AlignLeft, AlignRight},
		{"End LTR", AlignEnd, DirectionLTR, AlignLeft, AlignRight},
		{"End RTL", AlignEnd, DirectionRTL, AlignLeft, AlignLeft},
		{"Match parent", AlignMatchParent, DirectionLTR, AlignCenter, AlignCenter},
		{"Left unchanged", AlignLeft, DirectionLTR, AlignLeft, AlignLeft},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := txt.resolveAlignment(tt.align, tt.direction, tt.parentAlign)
			if result != tt.expected {
				t.Errorf("resolveAlignment() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestWidthBytes(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name     string
		bytes    []byte
		expected float64
	}{
		{"ASCII bytes", []byte("Hello"), 5.0},
		{"CJK bytes", []byte("世界"), 4.0},
		{"Mixed bytes", []byte("Hello世界"), 9.0},
		{"Emoji bytes", []byte("😀"), 2.0},
		{"Empty bytes", []byte(""), 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := txt.WidthBytes(tt.bytes)
			if got != tt.expected {
				t.Errorf("WidthBytes(%q) = %.1f, want %.1f", tt.bytes, got, tt.expected)
			}
		})
	}
}

func TestWidthUpTo(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name           string
		text           string
		maxWidth       float64
		expectedWidth  float64
		expectedExceed bool
	}{
		{
			name:           "Within limit",
			text:           "Hello",
			maxWidth:       10.0,
			expectedWidth:  5.0,
			expectedExceed: false,
		},
		{
			name:           "Exactly at limit",
			text:           "Hello",
			maxWidth:       5.0,
			expectedWidth:  5.0,
			expectedExceed: false,
		},
		{
			name:           "Exceeds limit - stops early",
			text:           "Hello world",
			maxWidth:       8.0,
			expectedWidth:  9.0, // Stops at "Hello wor" (9 chars) - returns width when exceeded
			expectedExceed: true,
		},
		{
			name:           "Empty text",
			text:           "",
			maxWidth:       10.0,
			expectedWidth:  0.0,
			expectedExceed: false,
		},
		{
			name:           "CJK exceeds - stops early",
			text:           "世界世界",
			maxWidth:       5.0,
			expectedWidth:  6.0, // Stops at "世界世" (6 cells) - returns width when exceeded
			expectedExceed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width, exceeded := txt.WidthUpTo(tt.text, tt.maxWidth)
			if width != tt.expectedWidth {
				t.Errorf("WidthUpTo(%q, %.1f) width = %.1f, want %.1f",
					tt.text, tt.maxWidth, width, tt.expectedWidth)
			}
			if exceeded != tt.expectedExceed {
				t.Errorf("WidthUpTo(%q, %.1f) exceeded = %v, want %v",
					tt.text, tt.maxWidth, exceeded, tt.expectedExceed)
			}
		})
	}
}

func TestWidthMany(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name     string
		strings  []string
		expected []float64
	}{
		{
			name:     "Mixed strings",
			strings:  []string{"Hello", "世界", "Test"},
			expected: []float64{5.0, 4.0, 4.0},
		},
		{
			name:     "Empty slice",
			strings:  []string{},
			expected: []float64{},
		},
		{
			name:     "Single string",
			strings:  []string{"Hello"},
			expected: []float64{5.0},
		},
		{
			name:     "With emoji",
			strings:  []string{"😀", "Hi😀", "世界"},
			expected: []float64{2.0, 4.0, 4.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := txt.WidthMany(tt.strings)
			if len(got) != len(tt.expected) {
				t.Errorf("WidthMany() returned %d widths, want %d", len(got), len(tt.expected))
				return
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("WidthMany()[%d] = %.1f, want %.1f", i, got[i], tt.expected[i])
				}
			}
		})
	}
}
