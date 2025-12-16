package text

// Text Position Utilities
//
// Pure text-level operations for working with positions within lines.
// These are primitives that a layout engine can use to build coordinate-based
// hit testing or cursor positioning.
//
// This library focuses on TEXT concerns:
// - Finding positions within text content
// - Converting between positions and x-offsets within a line
// - Identifying which line contains a position
//
// Layout concerns (screen coordinates, y-positions) belong in a layout engine.

// ═══════════════════════════════════════════════════════════════
//  Position to X-Offset (within a line)
// ═══════════════════════════════════════════════════════════════

// PositionToXOffset finds the x-offset of a text position within a line.
//
// This converts a character position to its visual offset from the start of the line.
// Returns 0 if position is before the line, line.Width if after.
//
// Example:
//
//	txt := text.NewTerminal()
//	line := text.Line{Content: "Hello", Start: 0, End: 5}
//	offset := txt.PositionToXOffset(line, 2) // Returns 2.0 (position of 'l')
func (t *Text) PositionToXOffset(line Line, position int) float64 {
	if position <= line.Start {
		return 0
	}

	if position >= line.End {
		return line.Width
	}

	// Calculate width up to position
	relativePos := position - line.Start
	runes := []rune(line.Content)

	if relativePos > len(runes) {
		relativePos = len(runes)
	}

	textUpToPos := string(runes[:relativePos])
	return t.Width(textUpToPos)
}

// ═══════════════════════════════════════════════════════════════
//  X-Offset to Position (within a line)
// ═══════════════════════════════════════════════════════════════

// XOffsetInfo contains information about a position found at an x-offset.
type XOffsetInfo struct {
	// Position in the text (rune index)
	Position int

	// CharXOffset is the x-offset of the character's left edge
	CharXOffset float64

	// CharWidth is the width of the character
	CharWidth float64

	// IsTrailing indicates if offset is in trailing half of character
	// Used for determining cursor position: before or after character
	IsTrailing bool

	// IsWithinLine indicates if the offset is within the line's width
	IsWithinLine bool
}

// XOffsetToPosition finds the character position at the given x-offset within a line.
//
// This is the core text operation for finding where a click or cursor should go
// within a line's content. The layout engine provides the x-offset within the line.
//
// Example:
//
//	txt := text.NewTerminal()
//	line := text.Line{Content: "Hello", Width: 5.0, Start: 0, End: 5}
//	info := txt.XOffsetToPosition(line, 2.3)
//	// info.Position = 2 (character 'l')
//	// info.IsTrailing = false (left half of 'l')
func (t *Text) XOffsetToPosition(line Line, xOffset float64) XOffsetInfo {
	content := line.Content
	if len(content) == 0 {
		return XOffsetInfo{
			Position:     line.Start,
			CharXOffset:  0,
			CharWidth:    0,
			IsTrailing:   false,
			IsWithinLine: xOffset >= 0 && xOffset <= line.Width,
		}
	}

	// Special case: before start of line
	if xOffset <= 0 {
		firstChar := string([]rune(content)[0])
		return XOffsetInfo{
			Position:     line.Start,
			CharXOffset:  0,
			CharWidth:    t.Width(firstChar),
			IsTrailing:   false,
			IsWithinLine: false,
		}
	}

	// Walk through characters to find position at xOffset
	graphemes := t.Graphemes(content)
	currentX := 0.0
	runeOffset := 0

	for _, grapheme := range graphemes {
		charWidth := t.Width(grapheme)

		// Check if xOffset is within this character
		if xOffset >= currentX && xOffset < currentX+charWidth {
			// Determine if in leading or trailing half
			isTrailing := (xOffset - currentX) > (charWidth / 2)

			return XOffsetInfo{
				Position:     line.Start + runeOffset,
				CharXOffset:  currentX,
				CharWidth:    charWidth,
				IsTrailing:   isTrailing,
				IsWithinLine: true,
			}
		}

		currentX += charWidth
		runeOffset += len([]rune(grapheme))
	}

	// xOffset is beyond the end of the line
	lastGrapheme := graphemes[len(graphemes)-1]
	lastWidth := t.Width(lastGrapheme)

	return XOffsetInfo{
		Position:     line.End,
		CharXOffset:  currentX - lastWidth,
		CharWidth:    lastWidth,
		IsTrailing:   true,
		IsWithinLine: xOffset <= line.Width,
	}
}

// ═══════════════════════════════════════════════════════════════
//  Line Identification
// ═══════════════════════════════════════════════════════════════

// LineContainingPosition finds which line contains the given text position.
//
// Returns:
//   - Line index if position is within a line
//   - -1 if position is before all lines
//   - len(lines) if position is after all lines
//
// Example:
//
//	txt := text.NewTerminal()
//	lines := []text.Line{
//	    {Start: 0, End: 5},
//	    {Start: 5, End: 10},
//	}
//	lineIdx := txt.LineContainingPosition(lines, 7) // Returns 1
func (t *Text) LineContainingPosition(lines []Line, position int) int {
	if len(lines) == 0 {
		return -1
	}

	for i, line := range lines {
		// Use exclusive end check for all lines except the last
		// This ensures position at line boundary goes to the next line
		if i == len(lines)-1 {
			// Last line: inclusive end
			if position >= line.Start && position <= line.End {
				return i
			}
		} else {
			// Other lines: exclusive end
			if position >= line.Start && position < line.End {
				return i
			}
		}
	}

	// Position not found
	if position < lines[0].Start {
		return -1
	}
	return len(lines)
}
