# Changelog

All notable changes to this project will be documented in this file.

## Unreleased

### Added

- Added the v2 module at `github.com/bolknote/go-gd/v2`.
- Added the public package under `pkg/gd`.
- Added Go-style image lifecycle, drawing, color, transform, filter, text, GIF animation, and format encode/decode APIs.
- Added `IOContext` wrappers for libgd `gdIOCtx` read/write flows.
- Added runtime feature and version helpers.
- Added API coverage, enum parity, wrapper, roundtrip, and pixel-level tests.
- Added GitHub Actions CI, `Makefile`, lint configuration, and MIT license.
- Added `examples/basic`.

### Changed

- Reworked the project as a breaking v2 rewrite.
- Replaced the old PHP-style package surface with Go-style constructors, methods, options structs, typed errors, and explicit `Close` ownership.
- Moved library code from the repository root to `pkg/gd`.
- Moved the sample image fixture to `testdata/images`.
- Replaced the old `README` with `README.md`.

### Removed

- Removed the v1 root-level `gd.go` and `gdcompat.go` APIs.
- Removed the old `example/sample.go`.
- Removed non-libgd helper APIs from the core package, including custom image helpers that are not backed by libgd.

### Migration Notes

- Import v2 as `github.com/bolknote/go-gd/v2/pkg/gd`.
- Use `NewTrueColor` / `NewPalette` instead of v1-style create helpers.
- Use `Decode*` and `Encode*` functions/methods instead of PHP-style names.
- Call `Image.Close()` when an image is no longer needed. `Destroy()` remains as a deprecated alias.
- Use `Color`, `RGBA`, `Rect`, `Point`, and format-specific option structs instead of untyped maps or loosely typed parameters.
