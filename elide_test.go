package text

import (
	"testing"
)

// ═══════════════════════════════════════════════════════════════
//  Basic Elision Tests
// ═══════════════════════════════════════════════════════════════

func TestElide(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name     string
		text     string
		maxWidth float64
		want     string
	}{
		{
			name:     "Middle elision",
			text:     "Hello world",
			maxWidth: 8,
			want:     "He...ld",
		},
		{
			name:     "No elision needed",
			text:     "Short",
			maxWidth: 10,
			want:     "Short",
		},
		{
			name:     "Path-like text",
			text:     "/very/long/path/to/file.txt",
			maxWidth: 20,
			want:     "/very/.../file.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := txt.Elide(tt.text, tt.maxWidth)
			if got != tt.want {
				t.Errorf("Elide() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestElideEnd(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name     string
		text     string
		maxWidth float64
		want     string
	}{
		{
			name:     "End elision",
			text:     "Hello world",
			maxWidth: 8,
			want:     "Hello...",
		},
		{
			name:     "Description text",
			text:     "This is a very long description",
			maxWidth: 15,
			want:     "This is a ve...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := txt.ElideEnd(tt.text, tt.maxWidth)
			if got != tt.want {
				t.Errorf("ElideEnd() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestElideStart(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name     string
		text     string
		maxWidth float64
		want     string
	}{
		{
			name:     "Start elision",
			text:     "Hello world",
			maxWidth: 8,
			want:     "...world",
		},
		{
			name:     "Path end important",
			text:     "/path/to/myfile.txt",
			maxWidth: 15,
			want:     "...myfile.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := txt.ElideStart(tt.text, tt.maxWidth)
			if got != tt.want {
				t.Errorf("ElideStart() = %q, want %q", got, tt.want)
			}
		})
	}
}

// ═══════════════════════════════════════════════════════════════
//  Custom Ellipsis Tests
// ═══════════════════════════════════════════════════════════════

func TestElideWith(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name     string
		text     string
		maxWidth float64
		ellipsis string
		want     string
	}{
		{
			name:     "Single character ellipsis",
			text:     "Hello world",
			maxWidth: 8,
			ellipsis: "…",
			want:     "Hel…rld",
		},
		{
			name:     "Bracketed ellipsis",
			text:     "Hello world",
			maxWidth: 12,
			ellipsis: "[...]",
			want:     "He[...]ld",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := txt.ElideWith(tt.text, tt.maxWidth, tt.ellipsis)
			if got != tt.want {
				t.Errorf("ElideWith() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestElideUnicode(t *testing.T) {
	txt := NewTerminal()

	got := txt.ElideUnicode("Hello world", 8)
	// Should use Unicode horizontal ellipsis (…)
	if len(got) == 0 {
		t.Error("ElideUnicode returned empty string")
	}

	// Check that it contains the Unicode ellipsis
	hasUnicodeEllipsis := false
	for _, r := range got {
		if r == '…' {
			hasUnicodeEllipsis = true
			break
		}
	}

	if !hasUnicodeEllipsis {
		t.Errorf("ElideUnicode() = %q, expected to contain '…'", got)
	}
}

// ═══════════════════════════════════════════════════════════════
//  Context Detection Tests
// ═══════════════════════════════════════════════════════════════

func TestDetectContext(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name string
		text string
		want ElideContext
	}{
		{
			name: "HTTP URL",
			text: "http://example.com/path",
			want: ElideContextURL,
		},
		{
			name: "HTTPS URL",
			text: "https://example.com/path",
			want: ElideContextURL,
		},
		{
			name: "Unix path",
			text: "/usr/local/bin/app",
			want: ElideContextPath,
		},
		{
			name: "Home path",
			text: "~/Documents/file.txt",
			want: ElideContextPath,
		},
		{
			name: "Windows path",
			text: "C:\\Users\\Name\\file.txt",
			want: ElideContextPath,
		},
		{
			name: "Email",
			text: "user@example.com",
			want: ElideContextEmail,
		},
		{
			name: "General text",
			text: "Just some text",
			want: ElideContextGeneral,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := txt.detectContext(tt.text)
			if got != tt.want {
				t.Errorf("detectContext(%q) = %v, want %v", tt.text, got, tt.want)
			}
		})
	}
}

func TestElideAuto(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name     string
		text     string
		maxWidth float64
	}{
		{
			name:     "Auto detect URL",
			text:     "https://example.com/very/long/path/to/resource",
			maxWidth: 25,
		},
		{
			name:     "Auto detect path",
			text:     "/usr/local/share/applications/myapp.desktop",
			maxWidth: 25,
		},
		{
			name:     "Auto detect email",
			text:     "verylongemailaddress@example.com",
			maxWidth: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := txt.ElideAuto(tt.text, tt.maxWidth)

			// Just verify it returns something shorter
			gotWidth := txt.Width(got)
			if gotWidth > tt.maxWidth {
				t.Errorf("ElideAuto() width %.1f exceeds maxWidth %.1f", gotWidth, tt.maxWidth)
			}

			// Should contain ellipsis
			hasEllipsis := false
			for _, r := range got {
				if r == '.' {
					hasEllipsis = true
					break
				}
			}

			if !hasEllipsis && got != tt.text {
				t.Errorf("ElideAuto() = %q, expected ellipsis", got)
			}
		})
	}
}

// ═══════════════════════════════════════════════════════════════
//  Unicode Text Elision Tests
// ═══════════════════════════════════════════════════════════════

func TestElide_CJK(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name     string
		text     string
		maxWidth float64
	}{
		{
			name:     "Chinese text",
			text:     "这是一个非常长的中文句子需要省略",
			maxWidth: 15,
		},
		{
			name:     "Japanese text",
			text:     "これは非常に長い日本語の文です",
			maxWidth: 15,
		},
		{
			name:     "Korean text",
			text:     "이것은 매우 긴 한국어 문장입니다",
			maxWidth: 15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := txt.Elide(tt.text, tt.maxWidth)

			gotWidth := txt.Width(got)
			if gotWidth > tt.maxWidth {
				t.Errorf("Elide() width %.1f exceeds maxWidth %.1f", gotWidth, tt.maxWidth)
			}

			// Should be shorter than original
			if txt.Width(got) >= txt.Width(tt.text) {
				t.Error("Elided text should be shorter than original")
			}
		})
	}
}

func TestElide_Emoji(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name     string
		text     string
		maxWidth float64
	}{
		{
			name:     "Text with emoji",
			text:     "Hello 😀 world 🎉",
			maxWidth: 12,
		},
		{
			name:     "Emoji with modifiers",
			text:     "👋🏻 Hello 👨‍👩‍👧‍👦",
			maxWidth: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := txt.Elide(tt.text, tt.maxWidth)

			gotWidth := txt.Width(got)
			if gotWidth > tt.maxWidth {
				t.Errorf("Elide() width %.1f exceeds maxWidth %.1f", gotWidth, tt.maxWidth)
			}
		})
	}
}

// ═══════════════════════════════════════════════════════════════
//  Benchmark Tests
// ═══════════════════════════════════════════════════════════════

func BenchmarkElide(b *testing.B) {
	txt := NewTerminal()
	text := "/very/long/path/to/some/directory/with/file.txt"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txt.Elide(text, 25)
	}
}

func BenchmarkElideAuto(b *testing.B) {
	txt := NewTerminal()
	text := "https://example.com/very/long/path/to/resource"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txt.ElideAuto(text, 25)
	}
}
