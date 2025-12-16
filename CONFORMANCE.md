# Conformance Testing

This document describes the official conformance tests available for the standards implemented in this library and how to run them.

## Overview

This library implements several Unicode standards and CSS specifications. Where official conformance test data exists, we reference it and provide guidance for running tests.

## Unicode Standards Conformance

### UAX #9: Bidirectional Algorithm

**Status**: ✅ **100% Conformance** (513,494/513,494 tests passing)

The `unicode/uax9` package has already achieved 100% conformance with official Unicode test data:
- `BidiTest.txt` - Full bidirectional algorithm test suite
- `BidiCharacterTest.txt` - Character-level bidi tests

**Test Files**: Already included in `unicode/uax9/` directory
**Source**: [Unicode Bidi Test Data](https://www.unicode.org/Public/UCD/latest/ucd/BidiTest.txt)

### UAX #14: Line Breaking Algorithm

**Status**: ⚠️ **Partial Implementation** (Simplified for practical use)

Official test data available but not included by default:
- `LineBreakTest.txt` - Complete line breaking test suite

**Running Tests**:
```bash
# Download the official test file
cd unicode/uax14
curl -O https://www.unicode.org/Public/UCD/latest/ucd/auxiliary/LineBreakTest.txt

# Run conformance tests
go test -v -run TestOfficialUnicodeVectors
```

**Source**: [LineBreakTest.txt](https://www.unicode.org/Public/UCD/latest/ucd/auxiliary/LineBreakTest.txt)

**Note**: Our UAX #14 implementation is simplified for practical text wrapping and may not pass 100% of edge cases. The implementation prioritizes:
- Common use cases (Latin, CJK, emoji)
- Performance
- Integration with CSS Text Module
- Practical text layout needs

### UAX #29: Text Segmentation

**Status**: ✅ **Official Tests Included**

The `unicode/uax29` package includes official conformance test data:
- `GraphemeBreakTest.txt` - Grapheme cluster boundary tests
- `WordBreakTest.txt` - Word boundary tests
- `SentenceBreakTest.txt` - Sentence boundary tests

**Test Files**: Already included in `unicode/uax29/` directory
**Source**: [Unicode Text Segmentation Tests](https://www.unicode.org/Public/UCD/latest/ucd/auxiliary/)

### UAX #11: East Asian Width

**Status**: ✅ **Implemented** (Character property data embedded)

Uses official Unicode Character Database (UCD) data for East Asian Width properties.

**Source**: [UAX #11 Specification](https://www.unicode.org/reports/tr11/)

### UTS #51: Unicode Emoji

**Status**: ✅ **Implemented**

Handles emoji sequences, skin tones, ZWJ sequences based on Unicode Emoji specification.

**Source**: [UTS #51 Specification](https://www.unicode.org/reports/tr51/)

## CSS Text Module Conformance

### CSS Text Module Level 3

**Status**: ✅ **100% Feature Complete** (15/15 properties)

Official W3C test suite available:
- [CSS Text Level 3 Test Suite](https://test.csswg.org/suites/css-text-3_dev/nightly-unstable/)
- [Test Harness](https://test.csswg.org/harness/)

**Note**: CSS conformance tests are browser/renderer tests (HTML/CSS files). Since this is a pure text library (not a renderer), we cannot directly run the W3C test suite. Instead:

1. **Property Coverage**: All CSS Text Level 3 properties are implemented (see `CSS_TEXT_COVERAGE.md`)
2. **Unit Tests**: Comprehensive Go tests verify each property behaves according to spec
3. **Integration Tests**: Tests verify property combinations work correctly

**Browser Testing**: To test in a browser context:
- Implement a renderer using this library
- Run the [W3C CSS Text Test Suite](https://test.csswg.org/suites/css-text-3_dev/nightly-unstable/)
- Submit implementation report to W3C

### CSS Text Module Level 4

**Status**: ✅ **100% In-Scope Properties** (18/18 text concerns)

Official test suite status:
- Level 4 is still in Editor's Draft status
- Test suite: [CSS Text Level 4 Tests](https://github.com/web-platform-tests/wpt/tree/master/css/css-text)
- Many Level 4 features don't have conformance tests yet

**Testing Approach**: Same as Level 3 - unit and integration tests verify spec-compliant behavior.

## Running All Tests

### Text Library Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run with race detection
go test -race ./...

# Run specific conformance tests
go test -v -run Conformance ./...
go test -v -run Official ./...
```

### Unicode Library Tests

```bash
# Run tests in unicode packages
cd ../unicode
go test ./...

# Run official conformance tests (if test data present)
cd uax9
go test -v -run TestBidiConformance

cd ../uax14
go test -v -run TestOfficialUnicodeVectors

cd ../uax29
go test -v -run TestGraphemeBreak
go test -v -run TestWordBreak
go test -v -run TestSentenceBreak
```

## Conformance Test Data Sources

### Unicode Consortium

All Unicode conformance test data is available from the official Unicode website:

- **Latest Version (17.0.0)**: https://www.unicode.org/Public/17.0.0/ucd/
- **All Versions**: https://www.unicode.org/Public/
- **Test Data Repository**: https://github.com/unicode-org/conformance

Specific test files:
- **Bidi Tests**: https://www.unicode.org/Public/UCD/latest/ucd/BidiTest.txt
- **Line Break**: https://www.unicode.org/Public/UCD/latest/ucd/auxiliary/LineBreakTest.txt
- **Grapheme Break**: https://www.unicode.org/Public/UCD/latest/ucd/auxiliary/GraphemeBreakTest.txt
- **Word Break**: https://www.unicode.org/Public/UCD/latest/ucd/auxiliary/WordBreakTest.txt
- **Sentence Break**: https://www.unicode.org/Public/UCD/latest/ucd/auxiliary/SentenceBreakTest.txt

### W3C CSS Working Group

CSS conformance tests are part of the Web Platform Tests:

- **Web Platform Tests (WPT)**: https://github.com/web-platform-tests/wpt
- **CSS Text Tests**: https://github.com/web-platform-tests/wpt/tree/master/css/css-text
- **Test Suites**: https://test.csswg.org/
- **CSS Text Level 3**: https://test.csswg.org/suites/css-text-3_dev/nightly-unstable/

## CI/CD Integration

GitHub Actions automatically runs all tests on:
- Every push to `main` branch
- Every pull request
- Multiple Go versions (1.21, 1.22, 1.23)
- Multiple platforms (Ubuntu, macOS)

See `.github/workflows/ci.yml` for configuration.

## Contributing Conformance Tests

If you'd like to add more conformance testing:

1. **Download Test Data**: Get official test files from Unicode or W3C
2. **Add Test Parser**: Implement test file parser if needed
3. **Run Tests**: Create test functions that compare against official results
4. **Document**: Update this file with new conformance status

### Example: Adding UAX #14 Conformance

```bash
# 1. Download test file
cd unicode/uax14
curl -O https://www.unicode.org/Public/UCD/latest/ucd/auxiliary/LineBreakTest.txt

# 2. Run existing test (already implemented)
go test -v -run TestOfficialUnicodeVectors

# 3. Results will show pass rate
```

## Conformance vs. Practical Implementation

This library prioritizes **practical text layout** over 100% conformance in some cases:

### Why Not 100% Conformance?

1. **Complexity vs. Value**: Some edge cases in UAX #14 are extremely rare
2. **Performance**: Full conformance may require slower algorithms
3. **Integration**: CSS Text Module has different priorities than pure Unicode
4. **Use Case**: Terminal TUIs have different needs than browsers

### Our Approach

- ✅ **Core behaviors**: 100% correct for common cases
- ✅ **Real-world text**: Handles Latin, CJK, emoji, bidi correctly
- ✅ **CSS integration**: Matches CSS Text Module expectations
- ⚠️ **Edge cases**: May differ on obscure Unicode sequences
- ✅ **Tested**: Comprehensive unit and integration tests

## Conformance Badges

### Unicode Standards

- **UAX #9 (Bidi)**: ![100% Conformance](https://img.shields.io/badge/conformance-100%25-brightgreen)
- **UAX #11 (EA Width)**: ![Implemented](https://img.shields.io/badge/status-implemented-blue)
- **UAX #14 (Line Break)**: ![Partial](https://img.shields.io/badge/conformance-partial-yellow)
- **UAX #29 (Segmentation)**: ![Official Tests](https://img.shields.io/badge/tests-official-blue)
- **UTS #51 (Emoji)**: ![Implemented](https://img.shields.io/badge/status-implemented-blue)

### CSS Specifications

- **CSS Text Level 3**: ![100% Features](https://img.shields.io/badge/features-100%25-brightgreen)
- **CSS Text Level 4**: ![100% In-Scope](https://img.shields.io/badge/in--scope-100%25-brightgreen)

## References

- [Unicode Conformance](http://www.unicode.org/reports/tr41/)
- [W3C CSS Test Guidelines](https://www.w3.org/Style/CSS/Test/)
- [Web Platform Tests](https://web-platform-tests.org/)

## License

Conformance test data from Unicode Consortium and W3C are under their respective licenses. See individual test files for license information.
