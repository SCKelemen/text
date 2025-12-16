# Comparison: github.com/SCKelemen/text vs golang.org/x/text

## Overview

**golang.org/x/text** is a low-level Unicode and internationalization library.
**github.com/SCKelemen/text** is a high-level text layout library for UI and layout engines.

## Scope Differences

### golang.org/x/text
- **Focus**: Low-level i18n primitives
- **Use case**: Building i18n applications, text processing pipelines
- **Level**: Building blocks for Unicode operations

### github.com/SCKelemen/text
- **Focus**: High-level text layout and rendering
- **Use case**: Terminal UIs, layout engines, text editors
- **Level**: Complete text layout solution

## Feature Comparison

| Feature | golang.org/x/text | github.com/SCKelemen/text |
|---------|-------------------|---------------------------|
| **Unicode Normalization** | ✅ (unicode/norm) | ❌ Not needed for layout |
| **Collation/Sorting** | ✅ (collate) | ❌ Not a layout concern |
| **Language Detection** | ✅ (language) | ❌ Not a layout concern |
| **Message Translation** | ✅ (message) | ❌ Not a layout concern |
| **Text Encoding** | ✅ (encoding) | ❌ Not a layout concern |
| **Width Measurement** | ⚠️ Basic (width) | ✅ Complete (UAX #11 + UTS #51) |
| **Line Breaking** | ❌ | ✅ Full UAX #14 |
| **Grapheme Clusters** | ❌ | ✅ Full UAX #29 |
| **Text Wrapping** | ❌ | ✅ Unicode-aware |
| **Text Truncation** | ❌ | ✅ Grapheme-aware |
| **Text Alignment** | ❌ | ✅ Left/Right/Center/Justify |
| **Bidirectional Text** | ⚠️ Basic (bidi) | ✅ Full UAX #9 integration |
| **Vertical Text** | ❌ | ✅ CSS Writing Modes + UAX #50 |
| **CSS Text Module** | ❌ | ✅ Level 3/4 |
| **Intrinsic Sizing** | ❌ | ✅ min-content, max-content |
| **Line Box Metrics** | ❌ | ✅ Baseline, ascent, descent |
| **Dictionary Support** | ❌ | ✅ Abbreviation-aware segmentation |

## golang.org/x/text Packages

### What They Provide

1. **unicode/norm** - Unicode normalization (NFC, NFD, NFKC, NFKD)
   - Canonical/compatibility equivalence
   - Not needed for text layout

2. **collate** - Language-sensitive string comparison
   - Sorting, searching
   - Not a layout concern

3. **language** - BCP 47 language tags
   - Language identification
   - Not directly used in layout (could be for dictionary selection)

4. **message** - i18n message translation
   - Printf-style internationalization
   - Not a layout concern

5. **encoding** - Character encoding conversion
   - Legacy encodings to/from UTF-8
   - Modern apps use UTF-8, not needed

6. **width** - East Asian Width **[Overlaps with our library]**
   - Basic width calculation
   - **Our library provides**: More complete UAX #11 + emoji support

7. **bidi** - Bidirectional algorithm **[Overlaps with our library]**
   - Basic UAX #9 implementation
   - **Our library provides**: Full UAX #9 with convenient APIs

## What We Add Beyond golang.org/x/text

### 1. Complete Unicode Text Layout Stack
- UAX #11 (East Asian Width) - with emoji
- UAX #14 (Line Breaking) - with hyphenation modes
- UAX #29 (Graphemes, Words, Sentences) - complete segmentation
- UAX #50 (Vertical Orientation) - for vertical text
- UTS #51 (Emoji) - proper emoji width

### 2. CSS Text Module Implementation
- White space processing (normal, pre, nowrap, etc.)
- Text transformations (uppercase, capitalize, fullwidth, etc.)
- Word/line breaking control
- Letter and word spacing
- Text alignment

### 3. Layout Engine Features
- Intrinsic sizing (min-content, max-content)
- Line box metrics (baseline, ascent, descent)
- Multi-line text bounds
- First/last baseline for alignment
- Vertical text layout

### 4. High-Level Operations
- Text wrapping with width constraints
- Grapheme-aware truncation with ellipsis
- Text alignment (left, right, center, justify)
- Dictionary-aware sentence segmentation

### 5. Type-Safe CSS Units
- Integration with github.com/SCKelemen/units
- Proper Length types (px, em, ch, etc.)
- Clean API boundaries

## When to Use Each

### Use golang.org/x/text when:
- Building i18n applications
- Need message translation
- Working with legacy text encodings
- Implementing collation/sorting
- Need Unicode normalization
- Language detection

### Use github.com/SCKelemen/text when:
- Building terminal UIs (TUI libraries)
- Implementing text editors
- Creating layout engines
- Rendering text in games
- Any application that needs proper text measurement and layout
- Need CSS-compatible text handling

## Can They Work Together?

**Yes!** They serve different purposes:

```go
import (
    xtext "golang.org/x/text"
    "golang.org/x/text/language"
    "golang.org/x/text/message"

    "github.com/SCKelemen/text"
)

// Use x/text for i18n
p := message.NewPrinter(language.English)
translated := p.Sprintf("Hello, %s!", name)

// Use our library for layout
txt := text.NewTerminal()
lines := txt.Wrap(translated, text.WrapOptions{
    MaxWidth: 40,
})
```

## golang.org/x/text/width Issues

The `width` package in golang.org/x/text has limitations:

1. **No Emoji Support** - Doesn't handle emoji modifiers correctly
   - `width.LookupRune('🏻')` gives wrong width
   - Our library integrates UTS #51

2. **Limited Context** - Basic wide/narrow classification
   - Our library: Full UAX #11 with ambiguous character handling

3. **No Layout Operations** - Just width lookup
   - Our library: Wrapping, truncation, alignment

## Summary

**golang.org/x/text**: Low-level i18n building blocks
**github.com/SCKelemen/text**: High-level text layout for UIs

They complement each other:
- Use **x/text** for internationalization concerns
- Use **our library** for text measurement, layout, and rendering

Our library fills the gap between Unicode primitives and actual text rendering needs.
