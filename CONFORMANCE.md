# Conformance Testing

This document describes the official conformance tests available for the standards implemented in this library and how to run them.

## Overview

This library implements CSS Text Module specifications and integrates with the `unicode` library for Unicode standard support. This document focuses on testing **this text library's implementation** - the Unicode standards themselves are tested in the `unicode` library.

## Unicode Standards Integration

This library integrates with the separate `unicode` library which provides:

- **UAX #9**: Bidirectional Algorithm (100% conformance - 513,494/513,494 tests)
- **UAX #11**: East Asian Width (character property data)
- **UAX #14**: Line Breaking Algorithm (practical implementation)
- **UAX #29**: Text Segmentation (official tests included)
- **UTS #51**: Unicode Emoji (sequence handling)

**Testing**: Unicode conformance is tested in the `unicode` library. This library's tests focus on:
- Correct integration with Unicode algorithms
- CSS Text Module behavior
- Text measurement and wrapping
- Bidirectional text handling

For Unicode conformance details, see the `unicode` library documentation.

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

## Running Tests

### Test Library Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run with race detection
go test -race ./...

# Run with verbose output
go test -v ./...
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
- Multiple Go versions (1.21, 1.23, 1.25)
- Platform: Ubuntu Latest
- Includes: tests, linting, builds, and benchmarks

See `.github/workflows/ci.yml` for configuration.

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

## Status Summary

### Text Library Implementation

- **CSS Text Level 3**: ![100% Features](https://img.shields.io/badge/features-100%25-brightgreen)
- **CSS Text Level 4**: ![100% In-Scope](https://img.shields.io/badge/in--scope-100%25-brightgreen)
- **Unicode Integration**: ![Fully Integrated](https://img.shields.io/badge/unicode-integrated-blue)
- **Test Coverage**: ![Comprehensive](https://img.shields.io/badge/tests-comprehensive-brightgreen)
- **CI/CD**: ![Active](https://img.shields.io/badge/ci-active-brightgreen)

## References

- [Unicode Conformance](http://www.unicode.org/reports/tr41/)
- [W3C CSS Test Guidelines](https://www.w3.org/Style/CSS/Test/)
- [Web Platform Tests](https://web-platform-tests.org/)

## License

This software is licensed under the BearWare 1.0 License (MIT-compatible).

See the LICENSE file for full text.
