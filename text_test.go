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
		{"CJK wide", "世界", 4.0}, // 2 + 2
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

func TestGraphemes(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name  string
		text  string
		want  int
	}{
		{"ASCII", "Hello", 5},
		{"CJK", "世界", 2},
		{"Emoji", "😀", 1},
		{"Emoji with modifier", "👋🏻", 1}, // Should be 1 grapheme cluster
		{"Complex emoji", "👨‍👩‍👧‍👦", 1},      // Family emoji with ZWJ
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
