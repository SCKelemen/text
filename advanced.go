package text

import (
	"sort"
	"strings"

	"github.com/SCKelemen/unicode/uax14"
)

// Advanced CSS Text Module Features
//
// Implements CSS Text Level 3/4 advanced features:
// - text-justify
// - hanging-punctuation
// - tab-size
// - text-wrap
// - text-spacing-trim
// - wrap-before/wrap-after

// ═══════════════════════════════════════════════════════════════
//  Text Justify (CSS Text §8)
// ═══════════════════════════════════════════════════════════════

// TextJustify specifies the justification method.
// Based on CSS Text Module Level 3 §8.2:
// https://www.w3.org/TR/css-text-3/#text-justify-property
type TextJustify int

const (
	// TextJustifyAuto allows the browser to choose (typically inter-word).
	TextJustifyAuto TextJustify = iota

	// TextJustifyNone disables justification.
	TextJustifyNone

	// TextJustifyInterWord adds space between words.
	// Best for scripts that use word separators (Latin, Greek, Cyrillic).
	TextJustifyInterWord

	// TextJustifyInterCharacter adds space between characters.
	// Best for CJK scripts without word separators.
	TextJustifyInterCharacter

	// TextJustifyDistribute distributes space evenly (similar to inter-character).
	// Legacy value, maps to inter-character.
	TextJustifyDistribute
)

// JustifyText applies justification to text to fit a specific width.
//
// Uses the specified justification method to distribute extra space.
//
// Example:
//
//	txt := text.NewTerminal()
//	justified := txt.JustifyText("Hello world", 20, text.TextJustifyInterWord)
func (t *Text) JustifyText(text string, targetWidth float64, method TextJustify) string {
	currentWidth := t.Width(text)

	// No justification needed
	if currentWidth >= targetWidth {
		return text
	}

	extraSpace := targetWidth - currentWidth

	switch method {
	case TextJustifyNone:
		return text

	case TextJustifyInterWord, TextJustifyAuto:
		return t.justifyInterWord(text, extraSpace)

	case TextJustifyInterCharacter, TextJustifyDistribute:
		return t.justifyInterCharacter(text, extraSpace)

	default:
		return t.justifyInterWord(text, extraSpace)
	}
}

// justifyInterWord distributes space between words.
func (t *Text) justifyInterWord(text string, extraSpace float64) string {
	words := strings.Fields(text)
	if len(words) <= 1 {
		return text // Can't justify single word
	}

	gaps := len(words) - 1
	spacePerGap := extraSpace / float64(gaps)

	var result strings.Builder
	for i, word := range words {
		result.WriteString(word)
		if i < len(words)-1 {
			// Add original space plus extra
			result.WriteString(t.makePadding(1.0 + spacePerGap))
		}
	}

	return result.String()
}

// justifyInterCharacter distributes space between characters.
func (t *Text) justifyInterCharacter(text string, extraSpace float64) string {
	graphemes := t.Graphemes(text)
	if len(graphemes) <= 1 {
		return text
	}

	gaps := len(graphemes) - 1
	spacePerGap := extraSpace / float64(gaps)

	var result strings.Builder
	for i, g := range graphemes {
		result.WriteString(g)
		if i < len(graphemes)-1 {
			result.WriteString(t.makePadding(spacePerGap))
		}
	}

	return result.String()
}

// ═══════════════════════════════════════════════════════════════
//  Hanging Punctuation (CSS Text §6)
// ═══════════════════════════════════════════════════════════════

// HangingPunctuation controls whether punctuation can hang outside the line box.
// Based on CSS Text Module Level 3 §6:
// https://www.w3.org/TR/css-text-3/#hanging-punctuation-property
type HangingPunctuation int

const (
	// HangingPunctuationNone disables hanging punctuation.
	HangingPunctuationNone HangingPunctuation = 0

	// HangingPunctuationFirst allows opening brackets/quotes to hang at start.
	HangingPunctuationFirst HangingPunctuation = 1 << 0

	// HangingPunctuationLast allows closing brackets/quotes to hang at end.
	HangingPunctuationLast HangingPunctuation = 1 << 1

	// HangingPunctuationForceEnd allows stops/commas to hang at end if needed.
	HangingPunctuationForceEnd HangingPunctuation = 1 << 2

	// HangingPunctuationAllowEnd allows stops/commas to hang at end if they don't fit.
	HangingPunctuationAllowEnd HangingPunctuation = 1 << 3
)

// Common hanging punctuation characters
var (
	openingPunctuation = map[rune]bool{
		'(':      true, '[': true, '{': true,
		'"':      true, '\'': true, // Straight quotes
		'\u201C': true, '\u2018': true, // " " ' '
		'«':      true, '‹': true, '「': true, '『': true,
	}

	closingPunctuation = map[rune]bool{
		')':      true, ']': true, '}': true,
		'"':      true, '\'': true, // Straight quotes
		'\u201D': true, '\u2019': true, // " " ' '
		'»':      true, '›': true, '」': true, '』': true,
	}

	stopsPunctuation = map[rune]bool{
		'.': true, ',': true, ';': true, ':': true,
		'!': true, '?': true,
		'。': true, '、': true, '，': true, '．': true,
	}
)

// IsOpeningPunctuation returns true if the rune is opening punctuation.
func IsOpeningPunctuation(r rune) bool {
	return openingPunctuation[r]
}

// IsClosingPunctuation returns true if the rune is closing punctuation.
func IsClosingPunctuation(r rune) bool {
	return closingPunctuation[r]
}

// IsStopPunctuation returns true if the rune is a stop/comma.
func IsStopPunctuation(r rune) bool {
	return stopsPunctuation[r]
}

// ShouldHang determines if punctuation should hang outside the line box.
func (t *Text) ShouldHang(text string, position int, mode HangingPunctuation) (shouldHang bool, hangWidth float64) {
	if mode == HangingPunctuationNone {
		return false, 0
	}

	runes := []rune(text)
	if position < 0 || position >= len(runes) {
		return false, 0
	}

	r := runes[position]

	// Check first position
	if position == 0 && (mode&HangingPunctuationFirst) != 0 {
		if IsOpeningPunctuation(r) {
			return true, t.config.MeasureFunc(r)
		}
	}

	// Check last position
	if position == len(runes)-1 {
		if (mode&HangingPunctuationLast) != 0 && IsClosingPunctuation(r) {
			return true, t.config.MeasureFunc(r)
		}
		if (mode&HangingPunctuationForceEnd) != 0 && IsStopPunctuation(r) {
			return true, t.config.MeasureFunc(r)
		}
		if (mode&HangingPunctuationAllowEnd) != 0 && IsStopPunctuation(r) {
			return true, t.config.MeasureFunc(r)
		}
	}

	return false, 0
}

// ═══════════════════════════════════════════════════════════════
//  Tab Size (CSS Text §7)
// ═══════════════════════════════════════════════════════════════

// TabSize controls the width of tab characters.
// Based on CSS Text Module Level 3 §7:
// https://www.w3.org/TR/css-text-3/#tab-size-property
type TabSize struct {
	// Value is the tab width
	Value float64

	// Unit determines if Value is spaces or a length
	Unit TabSizeUnit
}

// TabSizeUnit specifies how tab size is measured.
type TabSizeUnit int

const (
	// TabSizeSpaces measures tabs in number of space characters.
	TabSizeSpaces TabSizeUnit = iota

	// TabSizeLength measures tabs in CSS length units.
	TabSizeLength
)

// DefaultTabSize returns the default tab size (8 spaces).
func DefaultTabSize() TabSize {
	return TabSize{
		Value: 8,
		Unit:  TabSizeSpaces,
	}
}

// ExpandTabs expands tab characters according to tab-size.
func (t *Text) ExpandTabs(text string, tabSize TabSize) string {
	if !strings.Contains(text, "\t") {
		return text
	}

	var result strings.Builder
	column := 0.0

	for _, r := range text {
		if r == '\t' {
			// Calculate tab stop position
			var tabWidth float64
			if tabSize.Unit == TabSizeSpaces {
				spaceWidth := t.config.MeasureFunc(' ')
				tabStop := tabSize.Value * spaceWidth
				// Advance to next tab stop
				tabWidth = tabStop - (column - (float64(int(column/tabStop)) * tabStop))
			} else {
				// Use length directly
				tabWidth = tabSize.Value
			}

			// Insert spaces to reach tab stop
			numSpaces := int(tabWidth / t.config.MeasureFunc(' '))
			result.WriteString(strings.Repeat(" ", numSpaces))
			column += tabWidth
		} else {
			result.WriteRune(r)
			if r == '\n' {
				column = 0
			} else {
				column += t.config.MeasureFunc(r)
			}
		}
	}

	return result.String()
}

// ═══════════════════════════════════════════════════════════════
//  Text Wrap (CSS Text Level 4)
// ═══════════════════════════════════════════════════════════════

// TextWrap controls advanced text wrapping strategies.
// Based on CSS Text Module Level 4:
// https://www.w3.org/TR/css-text-4/#text-wrap-property
type TextWrap int

const (
	// TextWrapWrap allows wrapping at any break point (default).
	TextWrapWrap TextWrap = iota

	// TextWrapNowrap disables wrapping.
	TextWrapNowrap

	// TextWrapBalance tries to balance line lengths.
	// Good for headings and short blocks of text.
	TextWrapBalance

	// TextWrapStable minimizes reflow when editing.
	// Only wraps when next word doesn't fit.
	TextWrapStable

	// TextWrapPretty optimizes for readability.
	// Avoids orphans and short last lines.
	TextWrapPretty
)

// WrapBalanced wraps text with balanced line lengths.
func (t *Text) WrapBalanced(text string, maxWidth float64) []Line {
	words := t.Words(text)
	if len(words) == 0 {
		return nil
	}

	// Calculate total width
	totalWidth := t.Width(text)

	// Estimate optimal line count
	optimalLines := int((totalWidth / maxWidth) + 0.5)
	if optimalLines < 1 {
		optimalLines = 1
	}

	// Try to distribute evenly
	targetWidth := totalWidth / float64(optimalLines)

	// Build lines
	var lines []Line
	currentLine := ""
	currentWidth := 0.0
	lineStart := 0

	for _, word := range words {
		wordWidth := t.Width(word)
		spaceWidth := t.Width(" ")

		// Check if adding this word exceeds target or max width
		nextWidth := currentWidth
		if currentWidth > 0 {
			nextWidth += spaceWidth
		}
		nextWidth += wordWidth

		if currentWidth > 0 && (nextWidth > maxWidth || nextWidth > targetWidth) {
			// Start new line
			lines = append(lines, Line{
				Content: strings.TrimSpace(currentLine),
				Width:   currentWidth,
				Start:   lineStart,
				End:     lineStart + len([]rune(currentLine)),
			})
			currentLine = word
			currentWidth = wordWidth
			lineStart += len([]rune(currentLine))
		} else {
			if currentWidth > 0 {
				currentLine += " "
				currentWidth += spaceWidth
			}
			currentLine += word
			currentWidth += wordWidth
		}
	}

	// Add last line
	if currentLine != "" {
		lines = append(lines, Line{
			Content: strings.TrimSpace(currentLine),
			Width:   currentWidth,
			Start:   lineStart,
			End:     lineStart + len([]rune(currentLine)),
		})
	}

	return lines
}

// WrapPretty wraps text optimizing for readability.
func (t *Text) WrapPretty(text string, maxWidth float64) []Line {
	// Pretty wrapping avoids:
	// - Orphans (single word on last line)
	// - Very short last lines
	// - Hyphenated words at end of paragraphs

	lines := t.Wrap(text, WrapOptions{MaxWidth: maxWidth})

	// If last line is too short (< 40% of max), try to pull word from previous
	if len(lines) >= 2 {
		lastLine := lines[len(lines)-1]
		if lastLine.Width < maxWidth*0.4 {
			// Pull last word from second-to-last line
			secondLast := lines[len(lines)-2]
			words := strings.Fields(secondLast.Content)

			if len(words) > 1 {
				// Move last word to last line
				pullWord := words[len(words)-1]
				newSecondLast := strings.Join(words[:len(words)-1], " ")
				newLastLine := pullWord + " " + lastLine.Content

				lines[len(lines)-2].Content = newSecondLast
				lines[len(lines)-2].Width = t.Width(newSecondLast)
				lines[len(lines)-1].Content = newLastLine
				lines[len(lines)-1].Width = t.Width(newLastLine)
			}
		}
	}

	return lines
}

// ═══════════════════════════════════════════════════════════════
//  Text Spacing Trim (CSS Text Level 4)
// ═══════════════════════════════════════════════════════════════

// TextSpacingTrim controls trimming of CJK spacing.
// Based on CSS Text Module Level 4:
// https://www.w3.org/TR/css-text-4/#text-spacing-trim-property
type TextSpacingTrim int

const (
	// TextSpacingTrimNone disables spacing trim.
	TextSpacingTrimNone TextSpacingTrim = iota

	// TextSpacingTrimSpaceAll trims space before/after ideographic characters.
	TextSpacingTrimSpaceAll

	// TextSpacingTrimSpaceFirst trims space at line start.
	TextSpacingTrimSpaceFirst

	// TextSpacingTrimAuto uses language-specific rules.
	TextSpacingTrimAuto
)

// IsCJKIdeograph returns true if the rune is a CJK ideograph.
func IsCJKIdeograph(r rune) bool {
	return (r >= 0x4E00 && r <= 0x9FFF) || // CJK Unified Ideographs
		(r >= 0x3400 && r <= 0x4DBF) || // CJK Extension A
		(r >= 0x20000 && r <= 0x2A6DF) || // CJK Extension B
		(r >= 0x2A700 && r <= 0x2B73F) || // CJK Extension C
		(r >= 0x2B740 && r <= 0x2B81F) || // CJK Extension D
		(r >= 0x2B820 && r <= 0x2CEAF) || // CJK Extension E
		(r >= 0xF900 && r <= 0xFAFF) || // CJK Compatibility Ideographs
		(r >= 0x2F800 && r <= 0x2FA1F) // CJK Compatibility Supplement
}

// TrimCJKSpacing trims spacing around CJK characters according to the trim mode.
func (t *Text) TrimCJKSpacing(text string, mode TextSpacingTrim) string {
	if mode == TextSpacingTrimNone {
		return text
	}

	runes := []rune(text)
	var result []rune

	for i, r := range runes {
		// Skip spaces between CJK characters
		if r == ' ' && i > 0 && i < len(runes)-1 {
			prevCJK := IsCJKIdeograph(runes[i-1])
			nextCJK := IsCJKIdeograph(runes[i+1])

			if mode == TextSpacingTrimSpaceAll && prevCJK && nextCJK {
				continue // Skip this space
			}
		}

		// Skip leading space before CJK at line start
		if mode == TextSpacingTrimSpaceFirst && r == ' ' && i < len(runes)-1 {
			if i == 0 || runes[i-1] == '\n' {
				if IsCJKIdeograph(runes[i+1]) {
					continue
				}
			}
		}

		result = append(result, r)
	}

	return string(result)
}

// ═══════════════════════════════════════════════════════════════
//  Wrap Before/After Controls
// ═══════════════════════════════════════════════════════════════

// WrapControl specifies wrap behavior before/after an element.
type WrapControl int

const (
	// WrapControlAuto uses normal line breaking rules.
	WrapControlAuto WrapControl = iota

	// WrapControlAvoid avoids breaking before/after if possible.
	WrapControlAvoid

	// WrapControlAlways forces a break before/after.
	WrapControlAlways
)

// WrapPoint represents a position where wrapping behavior is controlled.
type WrapPoint struct {
	Position int         // Position in text (rune index)
	Before   WrapControl // Wrap behavior before this position
	After    WrapControl // Wrap behavior after this position
}

// WrapWithControls wraps text respecting wrap-before/wrap-after controls.
//
// This function modifies line break opportunities based on CSS wrap-before
// and wrap-after properties:
// - WrapControlAvoid: Avoids breaking at this position if possible
// - WrapControlAlways: Forces a break at this position
// - WrapControlAuto: Uses normal UAX #14 break opportunities
//
// Example:
//
//	controls := []WrapPoint{
//	    {Position: 5, After: WrapControlAlways}, // Force break after position 5
//	    {Position: 10, Before: WrapControlAvoid}, // Avoid breaking before position 10
//	}
//	lines := txt.WrapWithControls("Hello world test", 20, controls)
func (t *Text) WrapWithControls(text string, maxWidth float64, controls []WrapPoint) []Line {
	if len(controls) == 0 {
		// No controls, use regular wrapping
		return t.Wrap(text, WrapOptions{MaxWidth: maxWidth})
	}

	// Build a map of positions to controls for quick lookup
	controlMap := make(map[int]WrapPoint)
	for _, wp := range controls {
		controlMap[wp.Position] = wp
	}

	// Get UAX #14 break opportunities
	breakOpportunities := uax14.FindLineBreakOpportunities(text, uax14.HyphensManual)

	// Filter and modify break opportunities based on controls
	var allowedBreaks []int
	allowedBreaks = append(allowedBreaks, 0) // Always include start

	for _, breakPos := range breakOpportunities {
		if breakPos == 0 {
			continue // Skip start, already added
		}

		// Check for controls at this position
		allow := true

		// Check wrap-after control for previous position
		if ctrl, ok := controlMap[breakPos-1]; ok {
			if ctrl.After == WrapControlAvoid {
				allow = false
			} else if ctrl.After == WrapControlAlways {
				allow = true // Force this break
			}
		}

		// Check wrap-before control for current position
		if ctrl, ok := controlMap[breakPos]; ok {
			if ctrl.Before == WrapControlAvoid {
				allow = false
			} else if ctrl.Before == WrapControlAlways {
				allow = true // Force this break
			}
		}

		if allow {
			allowedBreaks = append(allowedBreaks, breakPos)
		}
	}

	// Add forced breaks from WrapControlAlways
	for pos, ctrl := range controlMap {
		if ctrl.After == WrapControlAlways {
			// Force break after this position
			found := false
			for _, bp := range allowedBreaks {
				if bp == pos+1 {
					found = true
					break
				}
			}
			if !found && pos+1 <= len(text) {
				allowedBreaks = append(allowedBreaks, pos+1)
			}
		}
		if ctrl.Before == WrapControlAlways {
			// Force break before this position
			found := false
			for _, bp := range allowedBreaks {
				if bp == pos {
					found = true
					break
				}
			}
			if !found && pos > 0 {
				allowedBreaks = append(allowedBreaks, pos)
			}
		}
	}

	// Sort break opportunities
	sort.Ints(allowedBreaks)

	// Build lines from allowed break points
	return t.buildLinesFromBreaks(text, allowedBreaks, maxWidth)
}

// buildLinesFromBreaks creates lines from break opportunities.
func (t *Text) buildLinesFromBreaks(text string, breakPoints []int, maxWidth float64) []Line {
	if len(breakPoints) == 0 {
		return []Line{{
			Content: text,
			Width:   t.Width(text),
			Start:   0,
			End:     len([]rune(text)),
		}}
	}

	var lines []Line
	currentLine := ""
	currentWidth := 0.0
	lineStartIdx := 0

	runes := []rune(text)

	for i := 1; i < len(breakPoints); i++ {
		start := breakPoints[i-1]
		end := breakPoints[i]

		if start >= len(runes) {
			break
		}
		if end > len(runes) {
			end = len(runes)
		}

		segment := string(runes[start:end])
		testLine := currentLine + segment
		testWidth := t.Width(testLine)

		if testWidth > maxWidth && currentLine != "" {
			// Line is full, commit current line
			lines = append(lines, Line{
				Content: currentLine,
				Width:   currentWidth,
				Start:   lineStartIdx,
				End:     lineStartIdx + len([]rune(currentLine)),
			})

			// Start new line
			currentLine = segment
			currentWidth = t.Width(segment)
			lineStartIdx += len([]rune(currentLine))
		} else {
			// Add to current line
			currentLine = testLine
			currentWidth = testWidth
		}
	}

	// Add final line
	if currentLine != "" {
		lines = append(lines, Line{
			Content: currentLine,
			Width:   currentWidth,
			Start:   lineStartIdx,
			End:     lineStartIdx + len([]rune(currentLine)),
		})
	}

	return lines
}

// ═══════════════════════════════════════════════════════════════
//  Text Autospace (CSS Text Level 4)
// ═══════════════════════════════════════════════════════════════

// TextAutospace controls automatic spacing around ideographic characters.
// Based on CSS Text Module Level 4:
// https://www.w3.org/TR/css-text-4/#text-autospace-property
type TextAutospace int

const (
	// TextAutospaceNormal creates extra spacing as specified.
	TextAutospaceNormal TextAutospace = iota

	// TextAutospaceNoAutospace disables automatic spacing.
	TextAutospaceNoAutospace

	// TextAutospaceAuto uses language-specific spacing rules.
	TextAutospaceAuto
)

// AutospaceFlags specifies which autospace features to apply.
type AutospaceFlags int

const (
	// AutospaceNone disables all autospace features.
	AutospaceNone AutospaceFlags = 0

	// AutospaceIdeographAlpha adds spacing between ideographic and non-ideographic characters.
	AutospaceIdeographAlpha AutospaceFlags = 1 << 0

	// AutospaceIdeographNumeric adds spacing between ideographic and numeric characters.
	AutospaceIdeographNumeric AutospaceFlags = 1 << 1

	// AutospacePunctuation adjusts spacing around fullwidth punctuation.
	AutospacePunctuation AutospaceFlags = 1 << 2

	// AutospaceAll enables all autospace features.
	AutospaceAll = AutospaceIdeographAlpha | AutospaceIdeographNumeric | AutospacePunctuation
)

// IsIdeographic returns true if the rune is an ideographic character.
// Includes CJK ideographs, Hiragana, Katakana, and Hangul.
func IsIdeographic(r rune) bool {
	// CJK Ideographs (already in IsCJKIdeograph)
	if IsCJKIdeograph(r) {
		return true
	}

	// Hiragana
	if r >= 0x3040 && r <= 0x309F {
		return true
	}

	// Katakana
	if r >= 0x30A0 && r <= 0x30FF {
		return true
	}

	// Hangul Syllables
	if r >= 0xAC00 && r <= 0xD7AF {
		return true
	}

	// Hangul Jamo
	if r >= 0x1100 && r <= 0x11FF {
		return true
	}

	return false
}

// IsFullwidthPunctuation returns true if the rune is fullwidth punctuation.
func IsFullwidthPunctuation(r rune) bool {
	fullwidthPunct := map[rune]bool{
		'、': true, '。': true, '，': true, '．': true,
		'：': true, '；': true, '！': true, '？': true,
		'「': true, '」': true, '『': true, '』': true,
		'（': true, '）': true, '【': true, '】': true,
		'《': true, '》': true, '〈': true, '〉': true,
		'〔': true, '〕': true, '｛': true, '｝': true,
	}
	return fullwidthPunct[r]
}

// IsOpeningFullwidthPunctuation returns true if the rune is opening fullwidth punctuation.
func IsOpeningFullwidthPunctuation(r rune) bool {
	opening := map[rune]bool{
		'「': true, '『': true, '（': true, '【': true,
		'《': true, '〈': true, '〔': true, '｛': true,
	}
	return opening[r]
}

// IsClosingFullwidthPunctuation returns true if the rune is closing fullwidth punctuation.
func IsClosingFullwidthPunctuation(r rune) bool {
	closing := map[rune]bool{
		'」': true, '』': true, '）': true, '】': true,
		'》': true, '〉': true, '〕': true, '｝': true,
	}
	return closing[r]
}

// ApplyAutospace applies automatic spacing according to text-autospace rules.
//
// Example:
//
//	txt := text.NewTerminal()
//	result := txt.ApplyAutospace("Hello世界123", text.AutospaceAll)
//	// Returns "Hello 世界 123" (with spacing between scripts)
func (t *Text) ApplyAutospace(text string, flags AutospaceFlags) string {
	if flags == AutospaceNone {
		return text
	}

	runes := []rune(text)
	if len(runes) == 0 {
		return text
	}

	var result []rune
	result = append(result, runes[0])

	for i := 1; i < len(runes); i++ {
		prev := runes[i-1]
		curr := runes[i]

		// Check if we need to insert space
		needSpace := false

		// Ideograph-Alpha spacing
		if (flags & AutospaceIdeographAlpha) != 0 {
			prevIdeo := IsIdeographic(prev)
			currIdeo := IsIdeographic(curr)
			prevAlpha := (prev >= 'A' && prev <= 'Z') || (prev >= 'a' && prev <= 'z')
			currAlpha := (curr >= 'A' && curr <= 'Z') || (curr >= 'a' && curr <= 'z')

			// Add space between ideographic and alphabetic
			if (prevIdeo && currAlpha) || (prevAlpha && currIdeo) {
				needSpace = true
			}
		}

		// Ideograph-Numeric spacing
		if (flags & AutospaceIdeographNumeric) != 0 {
			prevIdeo := IsIdeographic(prev)
			currIdeo := IsIdeographic(curr)
			prevNum := (prev >= '0' && prev <= '9')
			currNum := (curr >= '0' && curr <= '9')

			// Add space between ideographic and numeric
			if (prevIdeo && currNum) || (prevNum && currIdeo) {
				// Don't add space if already has space or punctuation
				if prev != ' ' && !IsFullwidthPunctuation(prev) {
					needSpace = true
				}
			}
		}

		// Punctuation spacing
		if (flags & AutospacePunctuation) != 0 {
			// Reduce space after opening punctuation
			if IsOpeningFullwidthPunctuation(prev) && curr == ' ' {
				// Skip the space (don't add it)
				continue
			}

			// Reduce space before closing punctuation
			if prev == ' ' && IsClosingFullwidthPunctuation(curr) {
				// Remove the previous space we added
				if len(result) > 0 && result[len(result)-1] == ' ' {
					result = result[:len(result)-1]
				}
			}
		}

		// Insert spacing if needed
		if needSpace && prev != ' ' {
			result = append(result, ' ')
		}

		result = append(result, curr)
	}

	return string(result)
}

// ApplyAutospaceMode applies text-autospace with predefined mode.
func (t *Text) ApplyAutospaceMode(text string, mode TextAutospace) string {
	switch mode {
	case TextAutospaceNoAutospace:
		return text
	case TextAutospaceAuto, TextAutospaceNormal:
		return t.ApplyAutospace(text, AutospaceAll)
	default:
		return text
	}
}

// ═══════════════════════════════════════════════════════════════
//  Word Break: Auto-Phrase (CSS Text Level 4)
// ═══════════════════════════════════════════════════════════════

// PhraseBreaker is an interface for language-specific phrase segmentation.
//
// This enables word-break: auto-phrase for CJK languages by allowing users
// to provide their own dictionary/ML-based segmentation.
//
// Users should integrate external libraries for phrase breaking:
//   - Chinese: jieba, pkuseg, HanLP
//   - Japanese: MeCab, kuromoji, Sudachi
//   - Korean: KoNLPy, Mecab-ko
//   - Thai: ICU, libthai
//
// Example:
//
//	type MyJiebaBreaker struct {
//	    segmenter *jieba.Segmenter
//	}
//
//	func (j *MyJiebaBreaker) FindPhrases(text string) []int {
//	    segments := j.segmenter.Cut(text, true)
//	    positions := []int{0}
//	    offset := 0
//	    for _, seg := range segments {
//	        offset += len([]rune(seg))
//	        positions = append(positions, offset)
//	    }
//	    return positions
//	}
//
//	breaker := &MyJiebaBreaker{segmenter: jieba.NewSegmenter()}
//	lines := txt.WrapWithPhrases("你好世界", 20, breaker)
type PhraseBreaker interface {
	// FindPhrases returns phrase boundary positions (rune indices).
	// Positions should be sorted and include 0 as the first position.
	//
	// Example: "你好世界" might return [0, 2, 4] for "你好" and "世界"
	FindPhrases(text string) []int
}

// WrapWithPhrases wraps text using phrase boundaries from a PhraseBreaker.
//
// This implements CSS word-break: auto-phrase, which breaks at natural
// phrase boundaries in languages that don't use spaces between words.
//
// The PhraseBreaker interface allows users to integrate language-specific
// dictionaries without requiring this library to ship with large dictionary
// files or ML models.
//
// Example:
//
//	// User provides their own phrase breaker
//	breaker := &MyChineseBreaker{}
//	lines := txt.WrapWithPhrases("你好世界，这是一个测试。", 20, breaker)
func (t *Text) WrapWithPhrases(text string, maxWidth float64, breaker PhraseBreaker) []Line {
	if breaker == nil {
		// Fallback to regular wrapping
		return t.Wrap(text, WrapOptions{MaxWidth: maxWidth})
	}

	// Get phrase boundaries from user-provided breaker
	phraseBoundaries := breaker.FindPhrases(text)

	// Use these as line break opportunities
	return t.buildLinesFromBreaks(text, phraseBoundaries, maxWidth)
}

// WrapWithPhrasesAndControls combines phrase breaking with wrap controls.
//
// This allows fine-grained control over phrase-based line breaking.
//
// Example:
//
//	controls := []WrapPoint{
//	    {Position: 5, After: WrapControlAvoid}, // Don't break after phrase at position 5
//	}
//	lines := txt.WrapWithPhrasesAndControls(text, 20, breaker, controls)
func (t *Text) WrapWithPhrasesAndControls(text string, maxWidth float64, breaker PhraseBreaker, controls []WrapPoint) []Line {
	if breaker == nil {
		// Fallback to wrap with controls only
		return t.WrapWithControls(text, maxWidth, controls)
	}

	// Get phrase boundaries
	phraseBoundaries := breaker.FindPhrases(text)

	// Build control map
	controlMap := make(map[int]WrapPoint)
	for _, wp := range controls {
		controlMap[wp.Position] = wp
	}

	// Filter phrase boundaries based on controls
	var allowedBreaks []int
	allowedBreaks = append(allowedBreaks, 0) // Always include start

	for _, breakPos := range phraseBoundaries {
		if breakPos == 0 {
			continue
		}

		allow := true

		// Check wrap-after control for previous position
		if ctrl, ok := controlMap[breakPos-1]; ok {
			if ctrl.After == WrapControlAvoid {
				allow = false
			} else if ctrl.After == WrapControlAlways {
				allow = true
			}
		}

		// Check wrap-before control for current position
		if ctrl, ok := controlMap[breakPos]; ok {
			if ctrl.Before == WrapControlAvoid {
				allow = false
			} else if ctrl.Before == WrapControlAlways {
				allow = true
			}
		}

		if allow {
			allowedBreaks = append(allowedBreaks, breakPos)
		}
	}

	// Add forced breaks
	for pos, ctrl := range controlMap {
		if ctrl.After == WrapControlAlways || ctrl.Before == WrapControlAlways {
			found := false
			checkPos := pos
			if ctrl.After == WrapControlAlways {
				checkPos = pos + 1
			}
			for _, bp := range allowedBreaks {
				if bp == checkPos {
					found = true
					break
				}
			}
			if !found && checkPos > 0 && checkPos <= len(text) {
				allowedBreaks = append(allowedBreaks, checkPos)
			}
		}
	}

	// Sort breaks
	sort.Ints(allowedBreaks)

	// Build lines
	return t.buildLinesFromBreaks(text, allowedBreaks, maxWidth)
}
