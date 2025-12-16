package text

import (
	"testing"
)

// ═══════════════════════════════════════════════════════════════
//  Reorder (Existing Method) Tests
// ═══════════════════════════════════════════════════════════════

func TestReorder_Basic(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "English text",
			input: "Hello world",
		},
		{
			name:  "Hebrew text",
			input: "שלום",
		},
		{
			name:  "Arabic text",
			input: "مرحبا",
		},
		{
			name:  "Mixed text",
			input: "Hello שלום world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := txt.Reorder(tt.input)

			// Just verify it doesn't panic and returns something
			if result == "" && tt.input != "" {
				t.Error("Reorder() returned empty string")
			}

			t.Logf("Input:  %q", tt.input)
			t.Logf("Result: %q", result)
		})
	}
}

// ═══════════════════════════════════════════════════════════════
//  ReorderParagraph Tests
// ═══════════════════════════════════════════════════════════════

func TestReorderParagraph(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name  string
		input string
		dir   Direction
	}{
		{
			name:  "Single line LTR",
			input: "Hello world",
			dir:   DirectionLTR,
		},
		{
			name:  "Multiple LTR lines",
			input: "Hello world\nHow are you\nGoodbye",
			dir:   DirectionLTR,
		},
		{
			name:  "Single line RTL",
			input: "שלום",
			dir:   DirectionRTL,
		},
		{
			name:  "Mixed directions with auto",
			input: "Hello\nשלום\nWorld",
			dir:   DirectionAuto,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := txt.ReorderParagraph(tt.input, tt.dir)

			if result == "" && tt.input != "" {
				t.Error("ReorderParagraph() returned empty string")
			}

			t.Logf("Input:\n%s", tt.input)
			t.Logf("Result:\n%s", result)
		})
	}
}

func TestReorderParagraph_Empty(t *testing.T) {
	txt := NewTerminal()

	result := txt.ReorderParagraph("", DirectionLTR)

	if result != "" {
		t.Errorf("ReorderParagraph() = %q, want empty string", result)
	}
}

// ═══════════════════════════════════════════════════════════════
//  ReorderLine Tests
// ═══════════════════════════════════════════════════════════════

func TestReorderLine(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name  string
		input string
		dir   Direction
	}{
		{
			name:  "LTR line",
			input: "Hello world",
			dir:   DirectionLTR,
		},
		{
			name:  "RTL line",
			input: "שלום עולם",
			dir:   DirectionRTL,
		},
		{
			name:  "Mixed line",
			input: "Hello שלום",
			dir:   DirectionAuto,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := txt.ReorderLine(tt.input, tt.dir)

			if result == "" && tt.input != "" {
				t.Error("ReorderLine() returned empty string")
			}

			t.Logf("Input:  %q", tt.input)
			t.Logf("Result: %q", result)
		})
	}
}

// ═══════════════════════════════════════════════════════════════
//  MirrorBrackets Tests
// ═══════════════════════════════════════════════════════════════

func TestMirrorBrackets(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Round brackets",
			input: "(hello)",
			want:  ")hello(",
		},
		{
			name:  "Square brackets",
			input: "[world]",
			want:  "]world[",
		},
		{
			name:  "Curly brackets",
			input: "{test}",
			want:  "}test{",
		},
		{
			name:  "Angle brackets",
			input: "<tag>",
			want:  ">tag<",
		},
		{
			name:  "Mixed brackets",
			input: "([{<>}])",
			want:  ")]}><{[(",
		},
		{
			name:  "CJK brackets",
			input: "「こんにちは」",
			want:  "」こんにちは「",
		},
		{
			name:  "No brackets",
			input: "hello",
			want:  "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := txt.MirrorBrackets(tt.input)

			if result != tt.want {
				t.Errorf("MirrorBrackets(%q) = %q, want %q", tt.input, result, tt.want)
			}
		})
	}
}

// ═══════════════════════════════════════════════════════════════
//  GetBidiClass Tests
// ═══════════════════════════════════════════════════════════════

func TestGetBidiClass(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name  string
		char  rune
		want  BidiClass
	}{
		{
			name: "Latin letter",
			char: 'a',
			want: ClassL,
		},
		{
			name: "Hebrew letter",
			char: 'א',
			want: ClassR,
		},
		{
			name: "Arabic letter",
			char: 'ا',
			want: ClassAL,
		},
		{
			name: "European number",
			char: '5',
			want: ClassEN,
		},
		{
			name: "Space",
			char: ' ',
			want: ClassWS,
		},
		{
			name: "Comma",
			char: ',',
			want: ClassCS,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := txt.GetBidiClass(tt.char)

			if result != tt.want {
				t.Errorf("GetBidiClass(%q) = %v, want %v", tt.char, result, tt.want)
			}
		})
	}
}

// ═══════════════════════════════════════════════════════════════
//  Integration Tests
// ═══════════════════════════════════════════════════════════════

func TestBidi_Integration(t *testing.T) {
	txt := NewTerminal()

	t.Run("Wrap and reorder RTL text", func(t *testing.T) {
		text := "שלום עולם זהו טקסט ארוך מאוד שצריך להיות עטוף"

		// Wrap the text first
		lines := txt.Wrap(text, WrapOptions{MaxWidth: 20})

		// Then reorder each line
		for i := range lines {
			lines[i].Content = txt.ReorderLine(lines[i].Content, DirectionRTL)
		}

		if len(lines) == 0 {
			t.Error("Expected wrapped lines")
		}

		t.Logf("Wrapped and reordered %d lines", len(lines))
		for i, line := range lines {
			t.Logf("Line %d: %q", i, line.Content)
		}
	})

	t.Run("Reorder mixed text paragraph", func(t *testing.T) {
		text := "Hello world\nשלום עולם\nMixed text"

		result := txt.ReorderParagraph(text, DirectionAuto)

		t.Logf("Input:\n%s", text)
		t.Logf("Result:\n%s", result)
	})
}

// ═══════════════════════════════════════════════════════════════
//  Benchmarks
// ═══════════════════════════════════════════════════════════════

func BenchmarkReorder_LTR(b *testing.B) {
	txt := NewTerminal()
	text := "The quick brown fox jumps over the lazy dog"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txt.Reorder(text)
	}
}

func BenchmarkReorderParagraph(b *testing.B) {
	txt := NewTerminal()
	text := "Hello world\nHow are you\nGoodbye"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txt.ReorderParagraph(text, DirectionLTR)
	}
}

func BenchmarkReorderLine(b *testing.B) {
	txt := NewTerminal()
	text := "Hello שלום world"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txt.ReorderLine(text, DirectionAuto)
	}
}

func BenchmarkMirrorBrackets(b *testing.B) {
	txt := NewTerminal()
	text := "(hello [world] {test})"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txt.MirrorBrackets(text)
	}
}
