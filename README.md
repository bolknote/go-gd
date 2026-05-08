# go-gd

Go bindings for [libgd](https://libgd.github.io/), targeting modern Go and libgd 2.3.x.

This branch is a breaking v2 rewrite. The module path is:

```sh
go get github.com/bolknote/go-gd/v2
```

The public package lives in `pkg/gd`:

```go
import gd "github.com/bolknote/go-gd/v2/pkg/gd"
```

## Requirements

- Go 1.22 or newer
- cgo enabled
- libgd 2.3.x headers and library installed

On macOS with Homebrew:

```sh
brew install gd
```

On Debian/Ubuntu:

```sh
sudo apt-get install libgd-dev
```

The bindings prefer a normal system libgd installation. Advanced formats depend on how libgd was built: AVIF, HEIF, TIFF, WebP, FreeType, FontConfig, RAQM, and libimagequant may not be available in every package.

## Quick Start

```go
package main

import gd "github.com/bolknote/go-gd/v2/pkg/gd"

func main() {
	img, err := gd.NewTrueColor(320, 180)
	if err != nil {
		panic(err)
	}
	defer img.Close()

	white, err := img.AllocateColor(255, 255, 255)
	if err != nil {
		panic(err)
	}
	black, err := img.AllocateColor(0, 0, 0)
	if err != nil {
		panic(err)
	}

	if err := img.Fill(0, 0, white); err != nil {
		panic(err)
	}
	if err := img.Line(0, 0, img.Width()-1, img.Height()-1, black); err != nil {
		panic(err)
	}
	if err := img.EncodePNGFile("out.png", nil); err != nil {
		panic(err)
	}
}
```

Runnable examples live under `examples/`:

```sh
go run ./examples/basic
```

## API Shape

v2 replaces the old PHP-style API with Go-style methods:

- Decode functions return `(*Image, error)`.
- Encode functions return `error` or `([]byte, error)`.
- `Image.Close()` releases the underlying `gdImagePtr`. `Destroy()` is kept as a deprecated alias for v1 compatibility.
- Colors use `Color` and `RGBA` value types instead of `map[string]int`.
- Format-specific options use structs such as `JPEGOptions`, `PNGOptions`, `WebPOptions`, `HEIFOptions`, and `AVIFOptions`.
- Built-in bitmap text APIs draw Latin-1 strings. Use the FreeType APIs (`StringFT`, `StringFTEx`) with a `.ttf` font path for UTF-8 text.

## Concurrency

`*Image` and `*IOContext` values are not safe for concurrent use by multiple goroutines without external synchronization. Each value can be used by one goroutine at a time, including its `Close` call.

`SetErrorMethodDiscard`, `ClearErrorMethod`, `UseFontConfig`, and the `FontCache*` helpers manipulate libgd's process-global state.

## Streaming I/O

For pipes, sockets, and arbitrary `io.Writer`/`io.Reader`-like sinks, use `IOContext`:

```go
ctx, err := gd.NewDynamicWriteContext()
if err != nil {
    panic(err)
}
if err := img.EncodePNGContext(ctx, nil); err != nil {
    panic(err)
}
data, err := ctx.Extract()
```

`NewFileContext(path, mode)` wraps a C `FILE*` and is the only way to encode XBM through libgd; the convenience helpers `EncodeXBM`/`EncodeXBMFile` build on top of it.

## Feature Detection

Use the linked libgd runtime information before enabling optional formats:

```go
version := gd.Version()
features := gd.RuntimeFeatures()

_ = version.String
_ = features.WebP
```

For a single format:

```go
if gd.SupportsFormat(gd.FormatAVIF, false) {
	// AVIF decode is available in this libgd build.
}
```

## Covered Areas

The v2 package tracks public `BGD_DECLARE` symbols from `gd.h` and `gdfx.h`.
Every public symbol must be either wrapped by the Go package or explicitly
documented as an internal/deprecated/unsafe exclusion in the API coverage test.

The v2 package wraps the libgd 2.3.x surface:

- Image lifecycle, dimensions, clone, resolution, clipping
- PNG, JPEG, GIF, WebP, WBMP, BMP, TGA, TIFF, GD/GD2, HEIF, AVIF, XBM, XPM
- Buffer, file, and `gdIOCtx` encode/decode APIs
- Drawing primitives, polygons, styles, brushes, tiles, alpha handling
- Palette and truecolor helpers, quantization, color replace/match
- Copy, resize, resample, rotate, scale, flip, interpolation and affine transforms
- Native libgd filters, effects, and `gdfx.h`
- Crop and compare APIs
- Built-in bitmap fonts and FreeType text rendering
- GIF animation buffer, file, and context APIs
- Version and feature probes
- Deprecated source/sink compatibility APIs are intentionally excluded; use `gdIOCtx`.

## Testing

Run:

```sh
make check
```

`make check` lists the linked libgd API, runs the coverage audit through `go test`, and runs `go vet`.
Tests skip optional paths when the linked libgd does not support a feature. JPEG byte-for-byte golden tests are intentionally avoided because compression output can vary between libgd builds.

The root contains project metadata only. The library package is under `pkg/gd`, examples are under `examples/`, and image fixtures are under `testdata/`.

## Migrating From v1

The rewrite intentionally breaks old names:

- `CreateTrueColor` -> `NewTrueColor`
- `CreateFromJpeg` -> `DecodeJPEGFile`
- `ImageToJpegBuffer` -> `EncodeJPEG`
- `Sx` / `Sy` -> `Width` / `Height`
- `Jpeg`, `Png`, `Gif`, `Webp` -> `EncodeJPEGFile`, `EncodePNGFile`, `EncodeGIFFile`, `EncodeWebPFile`
- `ColorsForIndex` -> `ColorRGBA`

The old `gdcompat.go` helpers were removed from the core package. Native libgd filters replace grayscale, brightness, contrast, convolution, gaussian blur, emboss, edge detect, mean removal, smoothing, color adjustment, and negation. Non-libgd helpers such as font directory scanning, stack blur, and custom smooth ellipse drawing are no longer part of the core API.
