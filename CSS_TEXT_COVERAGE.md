# CSS Text Module Coverage

This document tracks the implementation status of CSS Text Module Level 3 and Level 4 specifications in this library.

## Specification References

- [CSS Text Module Level 3](https://www.w3.org/TR/css-text-3/) (W3C Recommendation)
- [CSS Text Module Level 4](https://drafts.csswg.org/css-text-4/) (Editor's Draft)
- [MDN CSS Text](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_text)
- [web.dev Typography](https://web.dev/learn/css/typography)

## CSS Text Module Level 3 - Coverage

**Status: 100% Complete (15/15 properties)**

| Property | Status | Implementation | Spec Section |
|----------|--------|----------------|--------------|
| `text-transform` | ✅ | `css.go:TextTransform` | [§2.1](https://www.w3.org/TR/css-text-3/#text-transform-property) |
| `white-space` | ✅ | `css.go:WhiteSpace` | [§4](https://www.w3.org/TR/css-text-3/#white-space-property) |
| `tab-size` | ✅ | `advanced.go:TabSize` | [§4.2](https://www.w3.org/TR/css-text-3/#tab-size-property) |
| `word-break` | ✅ | `css.go:WordBreak` | [§5.2](https://www.w3.org/TR/css-text-3/#word-break-property) |
| `line-break` | ✅ | `css.go:LineBreak` | [§5.3](https://www.w3.org/TR/css-text-3/#line-break-property) |
| `hyphens` | ✅ | `css.go:Hyphens` | [§5.4](https://www.w3.org/TR/css-text-3/#hyphens-property) |
| `overflow-wrap` | ✅ | `css.go:OverflowWrap` | [§5.5](https://www.w3.org/TR/css-text-3/#overflow-wrap-property) |
| `text-align` | ✅ | `text.go:Alignment` | [§6.1](https://www.w3.org/TR/css-text-3/#text-align-property) |
| `text-align-all` | ✅ | `css.go:CSSTextStyle` | [§6.2](https://www.w3.org/TR/css-text-3/#text-align-all-property) |
| `text-align-last` | ✅ | `css.go:CSSTextStyle.TextAlignLast` | [§6.3](https://www.w3.org/TR/css-text-3/#text-align-last-property) |
| `text-justify` | ✅ | `advanced.go:TextJustify` | [§6.4](https://www.w3.org/TR/css-text-3/#text-justify-property) |
| `word-spacing` | ✅ | `css.go:CSSTextStyle.WordSpacing` | [§7.1](https://www.w3.org/TR/css-text-3/#word-spacing-property) |
| `letter-spacing` | ✅ | `css.go:CSSTextStyle.LetterSpacing` | [§7.2](https://www.w3.org/TR/css-text-3/#letter-spacing-property) |
| `text-indent` | ✅ | `css.go:TextIndent` | [§8.1](https://www.w3.org/TR/css-text-3/#text-indent-property) |
| `hanging-punctuation` | ✅ | `advanced.go:HangingPunctuation` | [§8.2](https://www.w3.org/TR/css-text-3/#hanging-punctuation-property) |

## CSS Text Module Level 4 - Coverage

**Status: 100% Complete for Text Library Scope (18/18 in-scope properties)**

### ✅ Implemented Properties

| Property | Status | Implementation | Spec Section |
|----------|--------|----------------|--------------|
| `white-space-collapse` | ✅ | `advanced.go:WhiteSpaceCollapse` | [§3.1](https://drafts.csswg.org/css-text-4/#white-space-collapsing) |
| `white-space-trim` | ✅ | `advanced.go:WhiteSpaceTrim` | [§3.4](https://drafts.csswg.org/css-text-4/#white-space-trim-property) |
| `text-wrap-mode` | ✅ | `advanced.go:TextWrapMode` | [§5.1](https://drafts.csswg.org/css-text-4/#text-wrap-mode-property) |
| `text-wrap-style` | ✅ | `advanced.go:TextWrapStyle` | [§5.4](https://drafts.csswg.org/css-text-4/#text-wrap-style-property) |
| `text-wrap` | ✅ | `advanced.go:TextWrap` | [§5.7](https://drafts.csswg.org/css-text-4/#text-wrap-property) |
| `word-space-transform` | ✅ | `advanced.go:WordSpaceTransform` | [§3.3](https://drafts.csswg.org/css-text-4/#word-space-transform-property) |
| `wrap-inside` | ✅ | `advanced.go:WrapInside` | [§5.6](https://drafts.csswg.org/css-text-4/#wrap-inside-property) |
| `wrap-before` | ✅ | `advanced.go:WrapBefore` | [§5.5](https://drafts.csswg.org/css-text-4/#wrap-before-property) |
| `wrap-after` | ✅ | `advanced.go:WrapAfter` | [§5.5](https://drafts.csswg.org/css-text-4/#wrap-after-property) |
| `text-autospace` | ✅ | `advanced.go:TextAutospace` | [§8.3](https://drafts.csswg.org/css-text-4/#text-autospace-property) |
| `text-spacing-trim` | ✅ | `advanced.go:TextSpacingTrim` | [§8.4](https://drafts.csswg.org/css-text-4/#text-spacing-trim-property) |
| `text-spacing` | ✅ | `advanced.go:TextSpacing` | [§8.5](https://drafts.csswg.org/css-text-4/#text-spacing-property) |
| `text-group-align` | ✅ | `advanced.go:TextGroupAlign` | [§7.5](https://drafts.csswg.org/css-text-4/#text-group-align-property) |
| `line-padding` | ✅ | `advanced.go:LinePadding` | [§8.6](https://drafts.csswg.org/css-text-4/#line-padding-property) |
| `text-align: start` | ✅ | `text.go:AlignStart` | [§7.1](https://drafts.csswg.org/css-text-4/#text-align-property) |
| `text-align: end` | ✅ | `text.go:AlignEnd` | [§7.1](https://drafts.csswg.org/css-text-4/#text-align-property) |
| `text-align: match-parent` | ✅ | `text.go:AlignMatchParent` | [§7.1](https://drafts.csswg.org/css-text-4/#text-align-property) |
| `direction` | ✅ | `text.go:Direction` | [§2.4](https://drafts.csswg.org/css-text-4/#direction) |

### ⚠️ Out of Scope Properties

These properties require layout engine or rendering capabilities beyond text measurement:

| Property | Reason | Belongs In |
|----------|--------|------------|
| `hyphenate-character` | Rendering/display concern | Render Engine |
| `hyphenate-limit-zone` | Requires layout box calculations | Layout Engine |
| `hyphenate-limit-chars` | Could implement but minimal value without full hyphenation | Layout Engine |
| `hyphenate-limit-lines` | Requires multi-line layout context | Layout Engine |
| `hyphenate-limit-last` | Requires multi-line layout context | Layout Engine |

## Unicode Standards Coverage

Beyond CSS specifications, this library implements several Unicode standards:

| Standard | Status | Implementation | Purpose |
|----------|--------|----------------|---------|
| [UAX #9: Bidirectional Algorithm](https://unicode.org/reports/tr9/) | ✅ | `bidi.go` + `unicode/uax9` | RTL/LTR text reordering |
| [UAX #11: East Asian Width](https://unicode.org/reports/tr11/) | ✅ | `unicode/uax11` | Character width calculation |
| [UAX #14: Line Breaking](https://unicode.org/reports/tr14/) | ✅ | `unicode/uax14` | Line break opportunities |
| [UAX #29: Text Segmentation](https://unicode.org/reports/tr29/) | ✅ | `unicode/uax29` | Grapheme/word/sentence boundaries |
| [UTS #51: Unicode Emoji](https://unicode.org/reports/tr51/) | ✅ | `unicode/uts51` | Emoji properties and sequences |

## Advanced Text Layout Algorithms

| Algorithm | Status | Implementation | Reference |
|-----------|--------|----------------|-----------|
| **Knuth-Plass Line Breaking** | ✅ | `knuthplass.go` | [Breaking Paragraphs into Lines (1981)](https://www.eprg.org/G53DOC/pdfs/knuth-plass-breaking.pdf) |
| **Greedy Line Breaking** | ✅ | `text.go:Wrap()` | Standard first-fit algorithm |
| **Balanced Text Wrapping** | ✅ | `advanced.go:WrapBalanced()` | CSS Text Level 4 §5.4 |
| **Pretty Text Wrapping** | ✅ | `advanced.go:WrapPretty()` | CSS Text Level 4 §5.4 |

### Knuth-Plass Algorithm

The optimal line breaking algorithm from TeX:
- **Paper**: [Breaking Paragraphs into Lines](https://www.eprg.org/G53DOC/pdfs/knuth-plass-breaking.pdf) by Donald E. Knuth and Michael F. Plass (1981)
- **Implementation**: Uses dynamic programming to minimize total "badness" across a paragraph
- **Features**: Configurable tolerance, fitness classes, hyphenation penalties
- **Usage**: Opt-in via `WrapKnuthPlass()` - does not affect standard CSS wrapping

## Text Position Helpers

While not CSS specifications, these utilities support building layout engines:

| Function | Purpose | Use Case |
|----------|---------|----------|
| `XOffsetToPosition()` | Find text position at x-offset | Click/cursor positioning |
| `PositionToXOffset()` | Find x-offset of text position | Cursor rendering |
| `LineContainingPosition()` | Find line containing position | Multi-line cursor movement |

## Library Scope

This library focuses on **text concerns** only:

### ✅ In Scope (Text Library)
- Text measurement (width, height, intrinsic sizing)
- Line breaking and wrapping
- Text transformation (case, spacing)
- Bidirectional reordering (logical → visual)
- Text-position helpers (offset ↔ position within lines)

### ❌ Out of Scope (Layout Engine)
- Positioning lines at (x, y) coordinates
- Box model (padding, margins, borders)
- Vertical layout direction
- Coordinate-based hit testing
- Scroll management

### ❌ Out of Scope (Render Engine)
- ANSI codes and terminal control
- Colors and styling
- Actual drawing/rendering
- Font rasterization

## Conformance Testing

All CSS-related functionality has comprehensive test coverage:
- Unit tests for each CSS property
- Integration tests for property combinations
- Edge case testing (empty text, CJK, emoji, RTL)
- Performance benchmarks

Run tests:
```bash
go test ./...
```

## Future Considerations

Potential additions that remain in text library scope:

1. **Hyphenation Implementation**
   - Currently have hooks (`Hyphens` type, `HyphenPenalty`)
   - Could implement Liang's TeX hyphenation algorithm
   - Requires user-provided dictionaries (no built-in dictionaries)
   - Reference: [Liang's dissertation (1983)](https://tug.org/docs/liang/)

2. **Ruby Annotation Metrics**
   - Calculate line height adjustments for ruby text
   - Reference: [CSS Ruby Annotation Layout Module](https://drafts.csswg.org/css-ruby/)

3. **Vertical Text Metrics**
   - Character rotation affects measurement
   - Already have `TextOrientation` type
   - Could expand for better vertical text support

However, the current implementation is **feature-complete** for the core text layout use case.

## Contributing

When adding new features, ensure they:
1. Are text measurement/breaking concerns (not layout/rendering)
2. Include spec links in documentation
3. Have comprehensive test coverage
4. Follow existing code organization patterns

## License

This implementation follows the specifications as defined by the W3C and Unicode Consortium. See LICENSE file for code licensing.
