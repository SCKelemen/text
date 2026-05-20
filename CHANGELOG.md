# Changelog

## [v1.2.0] - 2026-05-20

### Changed

- `github.com/SCKelemen/unicode` replaced with `github.com/SCKelemen/unicode/v6` (`v6.2.0`). This brings the v6 performance improvements (ASCII fast paths, memory optimization, rule-based state machines) to all text operations that use UAX #9, #11, #14, #29, #50, and UTS #51.
- `github.com/SCKelemen/units` bumped from `v1.1.0` to `v1.2.1` (`ParseLength` rejects scientific notation).

### Note

This is a pure dependency migration with no source-level API changes. The `text` public API remains identical. All UAX-based operations benefit from the v6 performance improvements transparently.
