package text

import (
	"math"
	"strings"
)

// Knuth-Plass Line Breaking Algorithm
//
// This implements the optimal line breaking algorithm from TeX, as described in
// "Breaking Paragraphs into Lines" by Donald E. Knuth and Michael F. Plass.
//
// Unlike greedy line breaking (which fills each line as much as possible),
// Knuth-Plass finds the globally optimal set of breakpoints that minimizes
// the "badness" of the entire paragraph.
//
// References:
//   - Knuth & Plass (1981): https://www.eprg.org/G53DOC/pdfs/knuth-plass-breaking.pdf
//   - Practical implementation: https://defoe.dev/blog/optimal-text-layout

// ═══════════════════════════════════════════════════════════════
//  Knuth-Plass Configuration
// ═══════════════════════════════════════════════════════════════

// KnuthPlassOptions configures the Knuth-Plass line breaking algorithm.
type KnuthPlassOptions struct {
	// MaxWidth is the target line width
	MaxWidth float64

	// Tolerance controls how much stretch/shrink is acceptable
	// Higher values allow more variation in line lengths
	// Default: 1.0 (TeX uses 1.0)
	Tolerance float64

	// FitnessClass controls which fitness classes are compatible
	// Prevents mixing very tight and very loose lines adjacent to each other
	// Default: true
	FitnessClass bool

	// Hyphenate enables hyphenation for better line breaks
	// Default: false
	Hyphenate bool

	// HyphenPenalty is the penalty for breaking at a hyphen
	// Default: 50
	HyphenPenalty float64

	// LinePenalty is the penalty for each line (encourages fewer lines)
	// Default: 10
	LinePenalty float64
}

// DefaultKnuthPlassOptions returns sensible defaults.
func DefaultKnuthPlassOptions(maxWidth float64) KnuthPlassOptions {
	return KnuthPlassOptions{
		MaxWidth:      maxWidth,
		Tolerance:     1.0,
		FitnessClass:  true,
		Hyphenate:     false,
		HyphenPenalty: 50,
		LinePenalty:   10,
	}
}

// ═══════════════════════════════════════════════════════════════
//  Internal Types
// ═══════════════════════════════════════════════════════════════

// breakpoint represents a potential line break position.
type breakpoint struct {
	position int     // Position in text (rune index)
	demerits float64 // Total demerits up to this point
	ratio    float64 // Adjustment ratio for the line ending here
	line     int     // Line number
	fitness  int     // Fitness class (0=tight, 1=normal, 2=loose, 3=very loose)
	prev     *breakpoint
}

// box represents an item in the text (word or glue).
type box struct {
	content  string
	width    float64
	position int    // Starting position in original text
	isGlue   bool   // True for spaces, false for words
	penalty  float64 // Penalty for breaking after this item
}

// ═══════════════════════════════════════════════════════════════
//  Knuth-Plass Algorithm
// ═══════════════════════════════════════════════════════════════

// WrapKnuthPlass wraps text using the Knuth-Plass optimal line breaking algorithm.
//
// This produces better-looking paragraphs than greedy wrapping by considering
// all possible break points and choosing the set that minimizes total "badness".
//
// Example:
//
//	txt := text.NewTerminal()
//	opts := text.DefaultKnuthPlassOptions(40.0)
//	lines := txt.WrapKnuthPlass("The quick brown fox jumps over the lazy dog", opts)
func (t *Text) WrapKnuthPlass(text string, opts KnuthPlassOptions) []Line {
	// Break text into boxes (words and glue/spaces)
	boxes := t.textToBoxes(text)
	if len(boxes) == 0 {
		return nil
	}

	// Find optimal breakpoints using dynamic programming
	breakpoints := t.findOptimalBreakpoints(boxes, opts)
	if len(breakpoints) == 0 {
		// Fallback to greedy if Knuth-Plass fails
		return t.Wrap(text, WrapOptions{MaxWidth: opts.MaxWidth})
	}

	// Convert breakpoints to lines
	return t.breakpointsToLines(text, boxes, breakpoints)
}

// textToBoxes converts text into a sequence of boxes (words) and glue (spaces).
func (t *Text) textToBoxes(text string) []box {
	var boxes []box
	position := 0

	// Split by spaces
	parts := strings.Split(text, " ")

	for i, part := range parts {
		if part == "" {
			continue
		}

		// Add word box
		boxes = append(boxes, box{
			content:  part,
			width:    t.Width(part),
			position: position,
			isGlue:   false,
			penalty:  0,
		})
		position += len([]rune(part))

		// Add glue (space) after word, except for last word
		if i < len(parts)-1 {
			boxes = append(boxes, box{
				content:  " ",
				width:    t.Width(" "),
				position: position,
				isGlue:   true,
				penalty:  0,
			})
			position++
		}
	}

	return boxes
}

// findOptimalBreakpoints uses dynamic programming to find the best set of breakpoints.
func (t *Text) findOptimalBreakpoints(boxes []box, opts KnuthPlassOptions) []int {
	if len(boxes) == 0 {
		return nil
	}

	// Initialize with a breakpoint at the start
	active := []*breakpoint{
		{
			position: 0,
			demerits: 0,
			ratio:    0,
			line:     0,
			fitness:  1, // Normal fitness
			prev:     nil,
		},
	}

	// Try to find a breakpoint after each box
	for i := 0; i < len(boxes); i++ {
		// Skip glue - we only break after words
		if boxes[i].isGlue {
			continue
		}

		// Try breaking after this box from each active breakpoint
		var newActive []*breakpoint

		for _, activeNode := range active {
			// Calculate line width from activeNode to current position
			lineWidth := t.calculateLineWidth(boxes, activeNode.position, i)

			// Calculate adjustment ratio
			ratio := (opts.MaxWidth - lineWidth) / opts.MaxWidth
			if ratio < -1 {
				// Line is too full, can't break here
				continue
			}

			// Calculate badness
			badness := t.calculateBadness(ratio, opts.Tolerance)
			if badness >= 10000 {
				continue // Too bad
			}

			// Calculate demerits
			penalty := boxes[i].penalty
			if boxes[i].content[len(boxes[i].content)-1] == '-' {
				penalty = opts.HyphenPenalty
			}

			demerits := math.Pow(opts.LinePenalty+badness+penalty, 2)
			totalDemerits := activeNode.demerits + demerits

			// Determine fitness class
			fitness := t.calculateFitness(ratio)

			// Add fitness penalty if adjacent lines have incompatible fitness
			if opts.FitnessClass && math.Abs(float64(fitness-activeNode.fitness)) > 1 {
				totalDemerits += 100
			}

			// Create new breakpoint
			newBreakpoint := &breakpoint{
				position: i + 1,
				demerits: totalDemerits,
				ratio:    ratio,
				line:     activeNode.line + 1,
				fitness:  fitness,
				prev:     activeNode,
			}

			newActive = append(newActive, newBreakpoint)
		}

		// Update active nodes (keep best ones)
		if len(newActive) > 0 {
			active = append(active, newActive...)
			// Prune: keep only best breakpoints for each line number
			active = t.pruneBreakpoints(active)
		}
	}

	// Find best final breakpoint
	var best *breakpoint
	minDemerits := math.MaxFloat64

	for _, node := range active {
		if node.demerits < minDemerits {
			minDemerits = node.demerits
			best = node
		}
	}

	if best == nil {
		return nil
	}

	// Reconstruct breakpoint positions
	var positions []int
	for node := best; node != nil; node = node.prev {
		if node.position > 0 {
			positions = append([]int{node.position}, positions...)
		}
	}

	return positions
}

// calculateLineWidth computes the width of text from boxes[start] to boxes[end].
func (t *Text) calculateLineWidth(boxes []box, start, end int) float64 {
	width := 0.0
	for i := start; i <= end && i < len(boxes); i++ {
		width += boxes[i].width
	}
	return width
}

// calculateBadness computes the badness of a line based on its adjustment ratio.
func (t *Text) calculateBadness(ratio float64, tolerance float64) float64 {
	if ratio < -1 {
		return 10000 // Infinitely bad
	}
	if ratio > tolerance {
		return 10000 // Too loose
	}
	return 100 * math.Pow(math.Abs(ratio), 3)
}

// calculateFitness determines the fitness class based on adjustment ratio.
//
// Fitness classes:
//   - 0: tight (ratio < -0.5)
//   - 1: normal (-0.5 <= ratio <= 0.5)
//   - 2: loose (0.5 < ratio <= 1.0)
//   - 3: very loose (ratio > 1.0)
func (t *Text) calculateFitness(ratio float64) int {
	if ratio < -0.5 {
		return 0 // Tight
	}
	if ratio <= 0.5 {
		return 1 // Normal
	}
	if ratio <= 1.0 {
		return 2 // Loose
	}
	return 3 // Very loose
}

// pruneBreakpoints keeps only the best breakpoints for each line number.
func (t *Text) pruneBreakpoints(breakpoints []*breakpoint) []*breakpoint {
	// Group by line number
	byLine := make(map[int][]*breakpoint)
	for _, bp := range breakpoints {
		byLine[bp.line] = append(byLine[bp.line], bp)
	}

	// Keep best from each line
	var pruned []*breakpoint
	for _, bps := range byLine {
		// Keep the best (lowest demerits)
		best := bps[0]
		for _, bp := range bps[1:] {
			if bp.demerits < best.demerits {
				best = bp
			}
		}
		pruned = append(pruned, best)
	}

	return pruned
}

// breakpointsToLines converts breakpoint positions to Line objects.
func (t *Text) breakpointsToLines(text string, boxes []box, breakpoints []int) []Line {
	var lines []Line
	runes := []rune(text)

	start := 0
	for _, breakPos := range breakpoints {
		if breakPos > len(boxes) {
			breakPos = len(boxes)
		}

		// Find text position of this breakpoint
		textPos := 0
		if breakPos > 0 && breakPos <= len(boxes) {
			textPos = boxes[breakPos-1].position + len([]rune(boxes[breakPos-1].content))
		}

		if textPos > len(runes) {
			textPos = len(runes)
		}

		// Extract line content
		content := string(runes[start:textPos])
		content = strings.TrimSpace(content)

		if content != "" {
			lines = append(lines, Line{
				Content: content,
				Width:   t.Width(content),
				Start:   start,
				End:     start + len([]rune(content)),
			})
		}

		start = textPos
		// Skip space
		if start < len(runes) && runes[start] == ' ' {
			start++
		}
	}

	// Add remaining text as final line
	if start < len(runes) {
		content := strings.TrimSpace(string(runes[start:]))
		if content != "" {
			lines = append(lines, Line{
				Content: content,
				Width:   t.Width(content),
				Start:   start,
				End:     len(runes),
			})
		}
	}

	return lines
}
