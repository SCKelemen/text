# text - Unicode-Aware Text Operations for Go

A practical, high-level text manipulation library that provides Unicode-correct operations for terminal UIs, text editors, and layout engines.

## Features

- **Unicode-Correct Width Calculation** - Properly handles CJK characters, emoji, combining marks
- **Smart Text Wrapping** - Uses UAX #14 line breaking algorithm
- **Intelligent Truncation** - Respects grapheme cluster boundaries (won't break emoji!)
- **Text Alignment** - Left, right, center, justify
- **Bidirectional Text** - Supports Arabic, Hebrew mixed with Latin (UAX #9)
- **Grapheme Awareness** - Treats emoji sequences, combining marks as single units
- **Renderer-Agnostic** - Works with terminals (cells) and canvas (pixels)

## Installation

```bash
go get github.com/SCKelemen/text
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/SCKelemen/text"
)

func main() {
    // Create terminal text handler
    txt := text.NewTerminal()

    // Measure width (accounts for wide characters)
    width := txt.Width("Hello 世界")  // 9.0 (not 7!)

    // Truncate with ellipsis
    short := txt.Truncate("Very long text here", text.TruncateOptions{
        MaxWidth: 10,
        Strategy: text.TruncateEnd,
    })
    fmt.Println(short)  // "Very lo..."

    // Wrap text
    lines := txt.Wrap("Hello world! This is a test.", text.WrapOptions{
        MaxWidth: 15,
    })
    for _, line := range lines {
        fmt.Println(line.Content)
    }

    // Align text
    aligned := txt.Align("Hello", 20, text.AlignCenter)
    fmt.Printf("|%s|\n", aligned)  // "|       Hello        |"
}
```

## Why This Library?

Most Go text libraries get Unicode wrong:

❌ `len("世界")` = 6 bytes, not 2 characters
❌ `utf8.RuneCountInString("世界")` = 2 runes, but **4 terminal cells wide**
❌ `utf8.RuneCountInString("👋🏻")` = 2 runes, but **1 grapheme cluster**

✅ `text.Width("世界")` = 4.0 cells (correct!)
✅ `text.GraphemeCount("👋🏻")` = 1 (correct!)

This library **gets it right** by using proper Unicode algorithms.

## Core Operations

### Width Measurement

Correctly measures display width accounting for:
- CJK ideographs (2 cells)
- Emoji (2 cells)
- Emoji modifiers (0 width)
- Combining marks (0 width)
- Ambiguous width characters (configurable)

```go
txt := text.NewTerminal()

txt.Width("Hello")      // 5.0
txt.Width("世界")        // 4.0 (2 + 2)
txt.Width("😀")         // 2.0
txt.Width("👋🏻")        // 2.0 (emoji + modifier)
txt.Width("é")          // 1.0 (e + combining accent)
```

### Text Wrapping

Smart wrapping using UAX #14 line breaking algorithm:

```go
txt := text.NewTerminal()

lines := txt.Wrap("Hello 世界! This is a test.", text.WrapOptions{
    MaxWidth: 15,
})

for _, line := range lines {
    fmt.Printf("%.1f: %s\n", line.Width, line.Content)
}
// Output:
// 12.0: Hello 世界!
// 15.0: This is a test.
```

### Truncation with Ellipsis

Three truncation strategies, all grapheme-aware:

```go
txt := text.NewTerminal()

// Truncate at end
txt.Truncate("Hello world", text.TruncateOptions{
    MaxWidth: 10,
    Strategy: text.TruncateEnd,
})
// Output: "Hello w..."

// Truncate in middle
txt.Truncate("Hello world", text.TruncateOptions{
    MaxWidth: 10,
    Strategy: text.TruncateMiddle,
})
// Output: "Hel...rld"

// Truncate at start
txt.Truncate("Hello world", text.TruncateOptions{
    MaxWidth: 10,
    Strategy: text.TruncateStart,
})
// Output: "...o world"
```

### Text Alignment

```go
txt := text.NewTerminal()

txt.Align("Hello", 20, text.AlignLeft)      // "Hello               "
txt.Align("Hello", 20, text.AlignRight)     // "               Hello"
txt.Align("Hello", 20, text.AlignCenter)    // "       Hello        "
txt.Align("Hello world", 20, text.AlignJustify)  // Distributes padding
```

### Grapheme Cluster Operations

Properly handles user-perceived characters:

```go
txt := text.NewTerminal()

// Emoji with skin tone modifier = 1 grapheme
graphemes := txt.Graphemes("👋🏻")
fmt.Println(len(graphemes))  // 1

// Family emoji with ZWJ = 1 grapheme
graphemes = txt.Graphemes("👨‍👩‍👧‍👦")
fmt.Println(len(graphemes))  // 1

// Combining mark = 1 grapheme
graphemes = txt.Graphemes("é")  // e + ́
fmt.Println(len(graphemes))  // 1
```

### Bidirectional Text

Supports Arabic, Hebrew mixed with Latin:

```go
txt := text.NewTerminal()

// Properly reorders mixed LTR/RTL text
display := txt.Reorder("Hello שלום world")

// Auto-detect direction
dir := txt.DetectDirection("שלום עולם")  // DirectionRTL
```

## Configuration

The library is renderer-agnostic through the `MeasureFunc`:

### Terminal (Cell-Based)

```go
txt := text.NewTerminal()
// Or explicit configuration:
txt := text.New(text.Config{
    MeasureFunc: text.TerminalMeasure,
    AmbiguousAsWide: false,  // For non-East Asian terminals
})
```

### East Asian Terminals

```go
txt := text.NewTerminalEastAsian()
// Treats ambiguous characters as wide (2 cells)
```

### Custom (e.g., Canvas/Pixels)

```go
txt := text.New(text.Config{
    MeasureFunc: func(r rune) float64 {
        return fontFace.GlyphAdvance(r)
    },
})
```

## Integration with Layout Engines

Implements the `layout.TextMetricsProvider` interface:

```go
package main

import (
    "github.com/SCKelemen/layout"
    "github.com/SCKelemen/text"
)

func main() {
    // Create metrics provider
    metrics := text.NewTerminalMetrics()

    // Inject into layout engine
    layout.SetTextMetricsProvider(metrics)

    // Now your layout engine uses proper Unicode text measurement!
    root := layout.NewNode().
        WithWidth(layout.Value{Value: 80, Unit: layout.UnitCh}).
        SetText("Hello 世界! Unicode is handled correctly.")

    root.Layout(layout.Constraint{MaxWidth: 80})
}
```

## Unicode Standards

This library coordinates multiple Unicode standards:

- **UAX #11** (East Asian Width) - Character width classification
- **UAX #14** (Line Breaking) - Line break opportunities
- **UAX #29** (Text Segmentation) - Grapheme cluster boundaries
- **UAX #9** (Bidirectional) - RTL text reordering
- **UTS #51** (Emoji) - Emoji properties and width

All provided by [`github.com/SCKelemen/unicode`](https://github.com/SCKelemen/unicode).

## Use Cases

### Terminal UI Libraries

```go
// Table cell truncation
func renderCell(content string, width int) {
    txt := text.NewTerminal()
    display := txt.Truncate(content, text.TruncateOptions{
        MaxWidth: float64(width),
        Strategy: text.TruncateEnd,
    })
    fmt.Print(display)
}
```

### Progress Bars

```go
// Fit status message in terminal width
func updateStatus(message string, termWidth int) {
    txt := text.NewTerminal()
    display := txt.Truncate(message, text.TruncateOptions{
        MaxWidth: float64(termWidth - 10),
        Strategy: text.TruncateEnd,
    })
    fmt.Printf("\r%s", display)
}
```

### Text Editors

```go
// Cursor movement by grapheme clusters
func moveCursorRight(text string, pos int) int {
    txt := text.NewTerminal()
    graphemes := txt.Graphemes(text)
    if pos < len(graphemes)-1 {
        return pos + 1
    }
    return pos
}
```

### Log Formatting

```go
// Truncate long log lines
func formatLog(message string, maxWidth int) string {
    txt := text.NewTerminal()
    return txt.Truncate(message, text.TruncateOptions{
        MaxWidth: float64(maxWidth),
        Strategy: text.TruncateEnd,
    })
}
```

## API Reference

See [GoDoc](https://pkg.go.dev/github.com/SCKelemen/text) for complete API documentation.

### Main Types

- `Text` - Main text operations handler
- `Config` - Configuration for text measurement
- `MeasureFunc` - Function to measure rune width
- `Line` - Wrapped line of text
- `Metrics` - Layout engine integration

### Functions

**Measurement:**
- `Width(s string) float64`
- `WidthRange(s string, start, end int) float64`

**Wrapping:**
- `Wrap(text string, opts WrapOptions) []Line`

**Truncation:**
- `Truncate(text string, opts TruncateOptions) string`

**Alignment:**
- `Align(text string, width float64, align Alignment) string`

**Bidirectional:**
- `Reorder(text string) string`
- `DetectDirection(text string) Direction`

**Graphemes:**
- `Graphemes(text string) []string`
- `GraphemeCount(text string) int`
- `GraphemeAt(text string, index int) string`

## Performance

The library is designed for practical performance:

- Width calculation: O(n) where n is number of runes
- Wrapping: O(n × log m) where m is number of break opportunities
- Truncation: O(n)
- All operations use efficient binary search for Unicode lookups

Benchmarks on a modern CPU:
```
BenchmarkWidth-10        1000000    1043 ns/op
BenchmarkTruncate-10      500000    2891 ns/op
BenchmarkWrap-10          200000    7234 ns/op
```

## License

MIT

## Related Projects

- [`github.com/SCKelemen/unicode`](https://github.com/SCKelemen/unicode) - Unicode algorithm implementations
- [`github.com/SCKelemen/layout`](https://github.com/SCKelemen/layout) - CSS layout engine for Go
