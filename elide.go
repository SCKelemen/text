package text

import "strings"

// Elision Convenience Functions
//
// Provides intuitive APIs for common text elision patterns.
// These are convenience wrappers around Truncate() with sensible defaults.

// ═══════════════════════════════════════════════════════════════
//  Simple Elision
// ═══════════════════════════════════════════════════════════════

// Elide shortens text by removing the middle portion if it exceeds maxWidth.
//
// This is a convenience function that uses TruncateMiddle strategy.
// Common for displaying file paths, URLs, and identifiers.
//
// Example:
//
//	txt := text.NewTerminal()
//	short := txt.Elide("/very/long/path/to/some/file.txt", 20)
//	// Returns: "/very/.../file.txt"
func (t *Text) Elide(text string, maxWidth float64) string {
	return t.Truncate(text, TruncateOptions{
		MaxWidth: maxWidth,
		Strategy: TruncateMiddle,
		Ellipsis: "...",
	})
}

// ElideEnd shortens text by truncating at the end.
//
// Common for displaying snippets, descriptions, and body text.
//
// Example:
//
//	txt := text.NewTerminal()
//	short := txt.ElideEnd("This is a very long description", 20)
//	// Returns: "This is a very lo..."
func (t *Text) ElideEnd(text string, maxWidth float64) string {
	return t.Truncate(text, TruncateOptions{
		MaxWidth: maxWidth,
		Strategy: TruncateEnd,
		Ellipsis: "...",
	})
}

// ElideStart shortens text by truncating at the start.
//
// Common for displaying file paths where the end is most important.
//
// Example:
//
//	txt := text.NewTerminal()
//	short := txt.ElideStart("/path/to/myfile.txt", 15)
//	// Returns: "...myfile.txt"
func (t *Text) ElideStart(text string, maxWidth float64) string {
	return t.Truncate(text, TruncateOptions{
		MaxWidth: maxWidth,
		Strategy: TruncateStart,
		Ellipsis: "...",
	})
}

// ═══════════════════════════════════════════════════════════════
//  Path-Specific Elision
// ═══════════════════════════════════════════════════════════════

// ElidePath intelligently shortens file paths.
//
// Preserves the filename and important directory information.
// Uses middle truncation but tries to break at path separators.
//
// Example:
//
//	txt := text.NewTerminal()
//	short := txt.ElidePath("/usr/local/share/applications/myapp.desktop", 30)
//	// Returns: "/usr/.../applications/myapp.desktop"
func (t *Text) ElidePath(path string, maxWidth float64) string {
	if t.Width(path) <= maxWidth {
		return path
	}

	sep := "/"
	if strings.Contains(path, "\\") {
		sep = "\\"
	}

	parts := strings.Split(path, sep)
	if len(parts) < 2 {
		return t.Truncate(path, TruncateOptions{
			MaxWidth: maxWidth,
			Strategy: TruncateMiddle,
			Ellipsis: "...",
		})
	}

	filename := parts[len(parts)-1]
	leadSlash := strings.HasPrefix(path, sep)
	firstDir := ""
	for _, p := range parts {
		if p != "" {
			firstDir = p
			break
		}
	}
	if firstDir == "" || filename == "" {
		return t.Truncate(path, TruncateOptions{
			MaxWidth: maxWidth,
			Strategy: TruncateMiddle,
			Ellipsis: "...",
		})
	}

	prefix := firstDir + sep
	if leadSlash {
		prefix = sep + prefix
	}
	candidate := prefix + "..." + sep + filename
	if t.Width(candidate) <= maxWidth {
		return candidate
	}

	return t.Truncate(path, TruncateOptions{
		MaxWidth: maxWidth,
		Strategy: TruncateMiddle,
		Ellipsis: "...",
	})
}

// ElideURL intelligently shortens URLs.
//
// Preserves the domain and important parts of the path.
//
// Example:
//
//	txt := text.NewTerminal()
//	short := txt.ElideURL("https://example.com/very/long/path/to/resource", 35)
//	// Returns: "https://example.com/.../resource"
func (t *Text) ElideURL(url string, maxWidth float64) string {
	// For now, just use middle elision
	// TODO: Implement URL-aware elision that preserves domain
	return t.Elide(url, maxWidth)
}

// ═══════════════════════════════════════════════════════════════
//  Custom Ellipsis
// ═══════════════════════════════════════════════════════════════

// ElideWith shortens text using a custom ellipsis string.
//
// Useful for different UI contexts or localization.
//
// Example:
//
//	txt := text.NewTerminal()
//	short := txt.ElideWith("Long text", 10, "…")     // Single character ellipsis
//	short = txt.ElideWith("Long text", 10, " [...] ") // Bracketed ellipsis
func (t *Text) ElideWith(text string, maxWidth float64, ellipsis string) string {
	return t.Truncate(text, TruncateOptions{
		MaxWidth: maxWidth,
		Strategy: TruncateMiddle,
		Ellipsis: ellipsis,
	})
}

// ElideEndWith shortens text at the end with custom ellipsis.
func (t *Text) ElideEndWith(text string, maxWidth float64, ellipsis string) string {
	return t.Truncate(text, TruncateOptions{
		MaxWidth: maxWidth,
		Strategy: TruncateEnd,
		Ellipsis: ellipsis,
	})
}

// ElideStartWith shortens text at the start with custom ellipsis.
func (t *Text) ElideStartWith(text string, maxWidth float64, ellipsis string) string {
	return t.Truncate(text, TruncateOptions{
		MaxWidth: maxWidth,
		Strategy: TruncateStart,
		Ellipsis: ellipsis,
	})
}

// ═══════════════════════════════════════════════════════════════
//  Unicode-Aware Ellipsis
// ═══════════════════════════════════════════════════════════════

// ElideUnicode uses the proper Unicode horizontal ellipsis character (…).
//
// This is more typographically correct than three dots (...).
//
// Example:
//
//	txt := text.NewTerminal()
//	short := txt.ElideUnicode("Long text here", 10)
//	// Returns: "Long t…re" (using U+2026)
func (t *Text) ElideUnicode(text string, maxWidth float64) string {
	return t.ElideWith(text, maxWidth, "…")
}

// ElideEndUnicode truncates at end with Unicode ellipsis (…).
func (t *Text) ElideEndUnicode(text string, maxWidth float64) string {
	return t.ElideEndWith(text, maxWidth, "…")
}

// ElideStartUnicode truncates at start with Unicode ellipsis (…).
func (t *Text) ElideStartUnicode(text string, maxWidth float64) string {
	return t.ElideStartWith(text, maxWidth, "…")
}

// ═══════════════════════════════════════════════════════════════
//  Context-Aware Elision
// ═══════════════════════════════════════════════════════════════

// ElideContext represents different elision contexts.
type ElideContext int

const (
	// ElideContextGeneral is for general text (middle elision).
	ElideContextGeneral ElideContext = iota

	// ElideContextPath is for file paths (preserve filename).
	ElideContextPath

	// ElideContextURL is for URLs (preserve domain).
	ElideContextURL

	// ElideContextEmail is for email addresses (preserve domain).
	ElideContextEmail

	// ElideContextDescription is for descriptions (end elision).
	ElideContextDescription

	// ElideContextCode is for code identifiers (middle elision).
	ElideContextCode
)

// ElideAuto automatically chooses the best elision strategy based on content.
//
// Detects if text looks like a path, URL, email, etc. and applies
// the most appropriate elision strategy.
//
// Example:
//
//	txt := text.NewTerminal()
//	txt.ElideAuto("/path/to/file.txt", 15)          // Uses path elision
//	txt.ElideAuto("https://example.com/path", 20)  // Uses URL elision
//	txt.ElideAuto("user@domain.com", 15)           // Uses email elision
func (t *Text) ElideAuto(text string, maxWidth float64) string {
	context := t.detectContext(text)
	return t.ElideForContext(text, maxWidth, context)
}

// ElideForContext elides text using strategy appropriate for the context.
func (t *Text) ElideForContext(text string, maxWidth float64, context ElideContext) string {
	switch context {
	case ElideContextPath:
		return t.ElidePath(text, maxWidth)

	case ElideContextURL:
		return t.ElideURL(text, maxWidth)

	case ElideContextEmail:
		return t.Elide(text, maxWidth)

	case ElideContextDescription:
		return t.ElideEnd(text, maxWidth)

	case ElideContextCode, ElideContextGeneral:
		return t.Elide(text, maxWidth)

	default:
		return t.Elide(text, maxWidth)
	}
}

// detectContext attempts to detect what kind of text this is.
func (t *Text) detectContext(text string) ElideContext {
	// Simple heuristics for content detection
	if len(text) == 0 {
		return ElideContextGeneral
	}

	// Check for URL
	if len(text) > 7 && (text[:7] == "http://" || text[:7] == "https:/") {
		return ElideContextURL
	}
	if len(text) > 6 && text[:6] == "ftp://" {
		return ElideContextURL
	}

	// Check for file path
	if text[0] == '/' || text[0] == '~' {
		return ElideContextPath
	}
	// Windows path
	if len(text) > 2 && text[1] == ':' && (text[2] == '\\' || text[2] == '/') {
		return ElideContextPath
	}

	// Check for email
	for _, r := range text {
		if r == '@' {
			return ElideContextEmail
		}
	}

	// Default to general
	return ElideContextGeneral
}
