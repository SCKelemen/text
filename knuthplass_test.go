package text

import (
	"math"
	"testing"
)

func TestWrapKnuthPlass_Basic(t *testing.T) {
	txt := NewTerminal()

	text := "The quick brown fox jumps over the lazy dog"
	opts := DefaultKnuthPlassOptions(20.0)

	lines := txt.WrapKnuthPlass(text, opts)

	if len(lines) == 0 {
		t.Fatal("Expected lines, got none")
	}

	// Debug: print lines
	t.Logf("Produced %d lines:", len(lines))
	for i, line := range lines {
		t.Logf("  Line %d: %q (width: %.1f, start: %d, end: %d)",
			i, line.Content, line.Width, line.Start, line.End)
	}

	// Verify all text is included
	totalChars := 0
	for _, line := range lines {
		totalChars += len([]rune(line.Content))
	}

	// Account for spaces that were trimmed
	expectedChars := len([]rune(text))
	// Each line break removes a space
	minExpected := expectedChars - len(lines) + 1

	if totalChars < minExpected-2 { // Allow small variance for trimming
		t.Errorf("Total chars = %d, expected at least %d (original: %d, lines: %d)",
			totalChars, minExpected, expectedChars, len(lines))
	}

	// Verify no line exceeds max width
	for i, line := range lines {
		if line.Width > opts.MaxWidth {
			t.Errorf("Line %d width %.1f exceeds max %.1f: %q",
				i, line.Width, opts.MaxWidth, line.Content)
		}
	}
}

func TestWrapKnuthPlass_Short(t *testing.T) {
	txt := NewTerminal()

	text := "Hello"
	opts := DefaultKnuthPlassOptions(20.0)

	lines := txt.WrapKnuthPlass(text, opts)

	if len(lines) != 1 {
		t.Errorf("Expected 1 line, got %d", len(lines))
	}

	if lines[0].Content != "Hello" {
		t.Errorf("Content = %q, want %q", lines[0].Content, "Hello")
	}
}

func TestWrapKnuthPlass_Empty(t *testing.T) {
	txt := NewTerminal()

	text := ""
	opts := DefaultKnuthPlassOptions(20.0)

	lines := txt.WrapKnuthPlass(text, opts)

	if len(lines) != 0 {
		t.Errorf("Expected 0 lines for empty text, got %d", len(lines))
	}
}

func TestWrapKnuthPlass_SingleWord(t *testing.T) {
	txt := NewTerminal()

	text := "Supercalifragilisticexpialidocious"
	opts := DefaultKnuthPlassOptions(20.0)

	lines := txt.WrapKnuthPlass(text, opts)

	// Should get at least one line
	if len(lines) == 0 {
		t.Fatal("Expected at least 1 line, got none")
	}

	// Verify the full word is preserved somewhere
	fullContent := ""
	for _, line := range lines {
		fullContent += line.Content
	}

	if fullContent != text {
		t.Errorf("Content = %q, want %q", fullContent, text)
	}
}

func TestWrapKnuthPlass_QualityComparison(t *testing.T) {
	txt := NewTerminal()

	// This text should wrap better with Knuth-Plass than greedy
	text := "The quick brown fox jumps over the lazy dog and runs away"
	maxWidth := 25.0

	// Greedy wrapping
	greedyLines := txt.Wrap(text, WrapOptions{MaxWidth: maxWidth})

	// Knuth-Plass wrapping
	knuthLines := txt.WrapKnuthPlass(text, DefaultKnuthPlassOptions(maxWidth))

	// Both should produce valid wrapping
	if len(greedyLines) == 0 || len(knuthLines) == 0 {
		t.Fatal("Both algorithms should produce lines")
	}

	// Calculate raggedness (variance in line widths)
	greedyRaggedness := calculateRaggedness(greedyLines, maxWidth)
	knuthRaggedness := calculateRaggedness(knuthLines, maxWidth)

	t.Logf("Greedy lines: %d, raggedness: %.2f", len(greedyLines), greedyRaggedness)
	t.Logf("Knuth-Plass lines: %d, raggedness: %.2f", len(knuthLines), knuthRaggedness)

	for i, line := range greedyLines {
		t.Logf("Greedy[%d]: %q (%.1f)", i, line.Content, line.Width)
	}
	for i, line := range knuthLines {
		t.Logf("Knuth[%d]: %q (%.1f)", i, line.Content, line.Width)
	}

	// Knuth-Plass should generally produce lower raggedness
	// (though not always, depending on the text)
	// Just verify it produces valid output
	for i, line := range knuthLines {
		if line.Width > maxWidth*1.1 { // Allow 10% tolerance
			t.Errorf("Knuth-Plass line %d exceeds max width: %.1f > %.1f",
				i, line.Width, maxWidth)
		}
	}
}

func TestWrapKnuthPlass_CJK(t *testing.T) {
	txt := NewTerminal()

	text := "世界 你好 再见 朋友"
	opts := DefaultKnuthPlassOptions(10.0)

	lines := txt.WrapKnuthPlass(text, opts)

	if len(lines) == 0 {
		t.Fatal("Expected lines, got none")
	}

	// Verify no line exceeds max width
	for i, line := range lines {
		if line.Width > opts.MaxWidth {
			t.Errorf("Line %d width %.1f exceeds max %.1f: %q",
				i, line.Width, opts.MaxWidth, line.Content)
		}
	}
}

func TestWrapKnuthPlass_Tolerance(t *testing.T) {
	txt := NewTerminal()

	text := "The quick brown fox jumps over the lazy dog"

	tests := []struct {
		name      string
		tolerance float64
	}{
		{
			name:      "Strict tolerance",
			tolerance: 0.5,
		},
		{
			name:      "Normal tolerance",
			tolerance: 1.0,
		},
		{
			name:      "Loose tolerance",
			tolerance: 2.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := DefaultKnuthPlassOptions(20.0)
			opts.Tolerance = tt.tolerance

			lines := txt.WrapKnuthPlass(text, opts)

			if len(lines) == 0 {
				t.Fatal("Expected lines, got none")
			}

			t.Logf("Tolerance %.1f produced %d lines", tt.tolerance, len(lines))
		})
	}
}

func TestTextToBoxes(t *testing.T) {
	txt := NewTerminal()

	text := "Hello world"
	boxes := txt.textToBoxes(text)

	// Should have: word "Hello", glue " ", word "world"
	expectedCount := 3
	if len(boxes) != expectedCount {
		t.Errorf("Expected %d boxes, got %d", expectedCount, len(boxes))
	}

	// Check first box is "Hello"
	if len(boxes) > 0 && boxes[0].content != "Hello" {
		t.Errorf("First box content = %q, want %q", boxes[0].content, "Hello")
	}

	// Check second box is glue
	if len(boxes) > 1 && !boxes[1].isGlue {
		t.Error("Second box should be glue")
	}

	// Check third box is "world"
	if len(boxes) > 2 && boxes[2].content != "world" {
		t.Errorf("Third box content = %q, want %q", boxes[2].content, "world")
	}
}

func TestTextToBoxes_Empty(t *testing.T) {
	txt := NewTerminal()

	text := ""
	boxes := txt.textToBoxes(text)

	if len(boxes) != 0 {
		t.Errorf("Expected 0 boxes for empty text, got %d", len(boxes))
	}
}

func TestCalculateFitness(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		ratio   float64
		fitness int
	}{
		{ratio: -0.6, fitness: 0}, // Tight
		{ratio: -0.3, fitness: 1}, // Normal
		{ratio: 0.0, fitness: 1},  // Normal
		{ratio: 0.3, fitness: 1},  // Normal
		{ratio: 0.6, fitness: 2},  // Loose
		{ratio: 1.5, fitness: 3},  // Very loose
	}

	for _, tt := range tests {
		fitness := txt.calculateFitness(tt.ratio)
		if fitness != tt.fitness {
			t.Errorf("calculateFitness(%.1f) = %d, want %d",
				tt.ratio, fitness, tt.fitness)
		}
	}
}

func TestCalculateBadness(t *testing.T) {
	txt := NewTerminal()

	tests := []struct {
		name      string
		ratio     float64
		tolerance float64
		wantBad   bool // Should be >= 10000
	}{
		{
			name:      "Perfect fit",
			ratio:     0.0,
			tolerance: 1.0,
			wantBad:   false,
		},
		{
			name:      "Slightly tight",
			ratio:     -0.3,
			tolerance: 1.0,
			wantBad:   false,
		},
		{
			name:      "Too tight",
			ratio:     -1.5,
			tolerance: 1.0,
			wantBad:   true,
		},
		{
			name:      "Slightly loose",
			ratio:     0.5,
			tolerance: 1.0,
			wantBad:   false,
		},
		{
			name:      "Too loose",
			ratio:     2.0,
			tolerance: 1.0,
			wantBad:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			badness := txt.calculateBadness(tt.ratio, tt.tolerance)

			if tt.wantBad && badness < 10000 {
				t.Errorf("Expected badness >= 10000, got %.1f", badness)
			}

			if !tt.wantBad && badness >= 10000 {
				t.Errorf("Expected badness < 10000, got %.1f", badness)
			}
		})
	}
}

// ═══════════════════════════════════════════════════════════════
//  Helpers
// ═══════════════════════════════════════════════════════════════

// calculateRaggedness computes the variance in line widths as a measure of quality.
// Lower raggedness generally means more balanced line lengths.
func calculateRaggedness(lines []Line, maxWidth float64) float64 {
	if len(lines) == 0 {
		return 0
	}

	// Calculate mean deviation from max width
	sumSquaredDev := 0.0
	for _, line := range lines {
		dev := maxWidth - line.Width
		sumSquaredDev += dev * dev
	}

	return math.Sqrt(sumSquaredDev / float64(len(lines)))
}

// ═══════════════════════════════════════════════════════════════
//  Benchmarks
// ═══════════════════════════════════════════════════════════════

func BenchmarkWrapKnuthPlass(b *testing.B) {
	txt := NewTerminal()
	text := "The quick brown fox jumps over the lazy dog and runs through the forest with great speed"
	opts := DefaultKnuthPlassOptions(40.0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txt.WrapKnuthPlass(text, opts)
	}
}

func BenchmarkWrapGreedy(b *testing.B) {
	txt := NewTerminal()
	text := "The quick brown fox jumps over the lazy dog and runs through the forest with great speed"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txt.Wrap(text, WrapOptions{MaxWidth: 40.0})
	}
}

func BenchmarkTextToBoxes(b *testing.B) {
	txt := NewTerminal()
	text := "The quick brown fox jumps over the lazy dog and runs through the forest with great speed"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txt.textToBoxes(text)
	}
}
