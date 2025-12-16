package text

import "testing"

// ═══════════════════════════════════════════════════════════════
//  XOffsetToPosition Tests
// ═══════════════════════════════════════════════════════════════

func TestXOffsetToPosition_Basic(t *testing.T) {
	txt := NewTerminal()

	line := Line{
		Content: "Hello",
		Width:   5.0,
		Start:   0,
		End:     5,
	}

	tests := []struct {
		name           string
		xOffset        float64
		wantPosition   int
		wantIsTrailing bool
		wantWithinLine bool
	}{
		{
			name:           "Before first character",
			xOffset:        -1.0,
			wantPosition:   0,
			wantIsTrailing: false,
			wantWithinLine: false,
		},
		{
			name:           "First character left half",
			xOffset:        0.3,
			wantPosition:   0,
			wantIsTrailing: false,
			wantWithinLine: true,
		},
		{
			name:           "First character right half",
			xOffset:        0.7,
			wantPosition:   0,
			wantIsTrailing: true,
			wantWithinLine: true,
		},
		{
			name:           "Middle character",
			xOffset:        2.3,
			wantPosition:   2,
			wantIsTrailing: false,
			wantWithinLine: true,
		},
		{
			name:           "After last character",
			xOffset:        6.0,
			wantPosition:   5,
			wantIsTrailing: true,
			wantWithinLine: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := txt.XOffsetToPosition(line, tt.xOffset)

			if info.Position != tt.wantPosition {
				t.Errorf("Position = %d, want %d", info.Position, tt.wantPosition)
			}

			if info.IsTrailing != tt.wantIsTrailing {
				t.Errorf("IsTrailing = %v, want %v", info.IsTrailing, tt.wantIsTrailing)
			}

			if info.IsWithinLine != tt.wantWithinLine {
				t.Errorf("IsWithinLine = %v, want %v", info.IsWithinLine, tt.wantWithinLine)
			}
		})
	}
}

func TestXOffsetToPosition_CJK(t *testing.T) {
	txt := NewTerminal()

	// "世界" = 2 characters, each 2 cells wide = 4 total width
	line := Line{
		Content: "世界",
		Width:   4.0,
		Start:   0,
		End:     2,
	}

	tests := []struct {
		name           string
		xOffset        float64
		wantPosition   int
		wantIsTrailing bool
	}{
		{
			name:           "First CJK character left half",
			xOffset:        0.5,
			wantPosition:   0,
			wantIsTrailing: false,
		},
		{
			name:           "First CJK character right half",
			xOffset:        1.5,
			wantPosition:   0,
			wantIsTrailing: true,
		},
		{
			name:           "Second CJK character",
			xOffset:        2.5,
			wantPosition:   1,
			wantIsTrailing: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := txt.XOffsetToPosition(line, tt.xOffset)

			if info.Position != tt.wantPosition {
				t.Errorf("Position = %d, want %d", info.Position, tt.wantPosition)
			}

			if info.IsTrailing != tt.wantIsTrailing {
				t.Errorf("IsTrailing = %v, want %v", info.IsTrailing, tt.wantIsTrailing)
			}
		})
	}
}

func TestXOffsetToPosition_Empty(t *testing.T) {
	txt := NewTerminal()

	line := Line{
		Content: "",
		Width:   0.0,
		Start:   0,
		End:     0,
	}

	info := txt.XOffsetToPosition(line, 0.0)

	if info.Position != 0 {
		t.Errorf("Position = %d, want 0", info.Position)
	}

	if info.CharWidth != 0 {
		t.Errorf("CharWidth = %.1f, want 0", info.CharWidth)
	}
}

// ═══════════════════════════════════════════════════════════════
//  PositionToXOffset Tests
// ═══════════════════════════════════════════════════════════════

func TestPositionToXOffset(t *testing.T) {
	txt := NewTerminal()

	line := Line{
		Content: "Hello",
		Width:   5.0,
		Start:   10,
		End:     15,
	}

	tests := []struct {
		name        string
		position    int
		wantXOffset float64
	}{
		{
			name:        "Before line start",
			position:    5,
			wantXOffset: 0.0,
		},
		{
			name:        "At line start",
			position:    10,
			wantXOffset: 0.0,
		},
		{
			name:        "Middle of line",
			position:    12,
			wantXOffset: 2.0,
		},
		{
			name:        "At line end",
			position:    15,
			wantXOffset: 5.0,
		},
		{
			name:        "After line end",
			position:    20,
			wantXOffset: 5.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xOffset := txt.PositionToXOffset(line, tt.position)

			if xOffset != tt.wantXOffset {
				t.Errorf("PositionToXOffset() = %.1f, want %.1f", xOffset, tt.wantXOffset)
			}
		})
	}
}

// ═══════════════════════════════════════════════════════════════
//  LineContainingPosition Tests
// ═══════════════════════════════════════════════════════════════

func TestLineContainingPosition(t *testing.T) {
	txt := NewTerminal()

	lines := []Line{
		{Content: "Hello", Start: 0, End: 5},
		{Content: "World", Start: 5, End: 10},
		{Content: "Test", Start: 10, End: 14},
	}

	tests := []struct {
		name     string
		position int
		wantLine int
	}{
		{
			name:     "Before all lines",
			position: -1,
			wantLine: -1,
		},
		{
			name:     "First line start",
			position: 0,
			wantLine: 0,
		},
		{
			name:     "First line middle",
			position: 2,
			wantLine: 0,
		},
		{
			name:     "Second line start",
			position: 5,
			wantLine: 1,
		},
		{
			name:     "Second line middle",
			position: 7,
			wantLine: 1,
		},
		{
			name:     "Third line",
			position: 12,
			wantLine: 2,
		},
		{
			name:     "After all lines",
			position: 20,
			wantLine: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lineIdx := txt.LineContainingPosition(lines, tt.position)

			if lineIdx != tt.wantLine {
				t.Errorf("LineContainingPosition() = %d, want %d", lineIdx, tt.wantLine)
			}
		})
	}
}

func TestLineContainingPosition_Empty(t *testing.T) {
	txt := NewTerminal()

	lineIdx := txt.LineContainingPosition([]Line{}, 0)

	if lineIdx != -1 {
		t.Errorf("LineContainingPosition() = %d, want -1", lineIdx)
	}
}

// ═══════════════════════════════════════════════════════════════
//  Round-Trip Tests
// ═══════════════════════════════════════════════════════════════

func TestPositionXOffset_RoundTrip(t *testing.T) {
	txt := NewTerminal()

	line := Line{
		Content: "Hello world",
		Width:   11.0,
		Start:   0,
		End:     11,
	}

	// Test that PositionToXOffset and XOffsetToPosition are (mostly) inverses
	for position := line.Start; position <= line.End; position++ {
		// Position -> XOffset
		xOffset := txt.PositionToXOffset(line, position)

		// XOffset -> Position (add small offset to get center of character)
		info := txt.XOffsetToPosition(line, xOffset+0.1)

		// Should get the same position back (or next one if we're at boundary)
		if info.Position != position && info.Position != position+1 {
			t.Errorf("Round trip failed for position %d: xOffset=%.2f, got position %d",
				position, xOffset, info.Position)
		}
	}
}

// ═══════════════════════════════════════════════════════════════
//  Benchmarks
// ═══════════════════════════════════════════════════════════════

func BenchmarkXOffsetToPosition(b *testing.B) {
	txt := NewTerminal()
	line := Line{
		Content: "The quick brown fox jumps over the lazy dog",
		Width:   44.0,
		Start:   0,
		End:     44,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txt.XOffsetToPosition(line, 22.5)
	}
}

func BenchmarkPositionToXOffset(b *testing.B) {
	txt := NewTerminal()
	line := Line{
		Content: "The quick brown fox jumps over the lazy dog",
		Width:   44.0,
		Start:   0,
		End:     44,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txt.PositionToXOffset(line, 22)
	}
}

func BenchmarkLineContainingPosition(b *testing.B) {
	txt := NewTerminal()
	text := "Hello world! This is a longer text that will wrap across multiple lines."
	lines := txt.Wrap(text, WrapOptions{MaxWidth: 40})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		txt.LineContainingPosition(lines, 35)
	}
}
