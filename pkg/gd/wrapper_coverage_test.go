package gd

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func pixelRGBA(t *testing.T, im *Image, x, y int) RGBA {
	t.Helper()
	c, err := im.TrueColorPixel(x, y)
	if err != nil {
		t.Fatal(err)
	}
	return DecomposeTrueColor(c)
}

func palettePixelRGBA(t *testing.T, im *Image, x, y int) RGBA {
	t.Helper()
	c, err := im.Pixel(x, y)
	if err != nil {
		t.Fatal(err)
	}
	rgba, err := im.ColorRGBA(c)
	if err != nil {
		t.Fatal(err)
	}
	return rgba
}

func assertRGB(t *testing.T, got RGBA, r, g, b int) {
	t.Helper()
	if got.R != r || got.G != g || got.B != b {
		t.Fatalf("unexpected RGB: got %+v, want %d,%d,%d", got, r, g, b)
	}
}

func TestDrawingPrimitivesAndPixels(t *testing.T) {
	im := newTrueColor(t, 32, 32)
	defer func() { _ = im.Close() }()

	white := allocColor(t, im, 255, 255, 255)
	black := allocColor(t, im, 0, 0, 0)
	red := allocColor(t, im, 255, 0, 0)
	blue := allocColor(t, im, 0, 0, 255)
	green := allocColor(t, im, 0, 255, 0)

	if err := im.Fill(0, 0, white); err != nil {
		t.Fatal(err)
	}
	if !im.BoundsSafe(31, 31) || im.BoundsSafe(32, 32) {
		t.Fatal("unexpected BoundsSafe result")
	}
	if err := im.SetPixel(1, 1, black); err != nil {
		t.Fatal(err)
	}
	assertRGB(t, pixelRGBA(t, im, 1, 1), 0, 0, 0)

	if err := im.Rectangle(2, 2, 10, 10, red); err != nil {
		t.Fatal(err)
	}
	assertRGB(t, pixelRGBA(t, im, 2, 2), 255, 0, 0)

	if err := im.FilledRectangle(12, 2, 20, 10, blue); err != nil {
		t.Fatal(err)
	}
	assertRGB(t, pixelRGBA(t, im, 15, 5), 0, 0, 255)

	if err := im.DashedLine(0, 15, 31, 15, black); err != nil {
		t.Fatal(err)
	}
	if err := im.Arc(16, 16, 12, 10, 0, 180, red); err != nil {
		t.Fatal(err)
	}
	if err := im.FilledArc(16, 16, 8, 8, 0, 270, green, ArcPie); err != nil {
		t.Fatal(err)
	}
	if err := im.Ellipse(16, 16, 20, 10, black); err != nil {
		t.Fatal(err)
	}
	if err := im.FilledEllipse(25, 25, 6, 6, red); err != nil {
		t.Fatal(err)
	}
	assertRGB(t, pixelRGBA(t, im, 25, 25), 255, 0, 0)

	if err := im.FillToBorder(30, 30, red, green); err != nil {
		t.Fatal(err)
	}
	if err := im.Polygon([]Point{{1, 20}, {8, 20}, {8, 28}, {1, 28}}, black); err != nil {
		t.Fatal(err)
	}
	if err := im.OpenPolygon([]Point{{10, 20}, {16, 20}, {16, 28}}, red); err != nil {
		t.Fatal(err)
	}
	if err := im.FilledPolygon([]Point{{20, 20}, {28, 20}, {24, 28}}, blue); err != nil {
		t.Fatal(err)
	}
	assertRGB(t, pixelRGBA(t, im, 24, 23), 0, 0, 255)

	if err := im.SetThickness(2); err != nil {
		t.Fatal(err)
	}
	if err := im.SetStyle([]Color{red, Transparent, blue}); err != nil {
		t.Fatal(err)
	}
	if err := im.Line(0, 31, 31, 31, Styled); err != nil {
		t.Fatal(err)
	}
	if err := im.SetAntiAliased(black); err != nil {
		t.Fatal(err)
	}
	if err := im.SetAntiAliasedDontBlend(black, true); err != nil {
		t.Fatal(err)
	}
	if err := im.Line(0, 0, 31, 31, AntiAliased); err != nil {
		t.Fatal(err)
	}
	if err := im.AABlend(); err != nil {
		t.Fatal(err)
	}
}

func TestClipAlphaInterlaceResolutionAndBounds(t *testing.T) {
	im := newTrueColor(t, 10, 10)
	defer func() { _ = im.Close() }()

	if got := im.Bounds(); got.Width != 10 || got.Height != 10 {
		t.Fatalf("unexpected bounds: %+v", got)
	}
	if err := im.SetClip(Rect{X: 2, Y: 3, Width: 4, Height: 5}); err != nil {
		t.Fatal(err)
	}
	clip, err := im.Clip()
	if err != nil {
		t.Fatal(err)
	}
	if clip != (Rect{X: 2, Y: 3, Width: 4, Height: 5}) {
		t.Fatalf("unexpected clip: %+v", clip)
	}
	if err := im.SaveAlpha(true); err != nil {
		t.Fatal(err)
	}
	if err := im.AlphaBlending(false); err != nil {
		t.Fatal(err)
	}
	if err := im.Interlace(true); err != nil {
		t.Fatal(err)
	}
	if !im.Interlaced() {
		t.Fatal("expected interlaced image")
	}
	if err := im.SetResolution(72, 96); err != nil {
		t.Fatal(err)
	}
	if x, y := im.Resolution(); x != 72 || y != 96 {
		t.Fatalf("unexpected resolution: %d,%d", x, y)
	}
	if got, err := im.InterpolationMethod(); err != nil || got < 0 {
		t.Fatalf("unexpected interpolation method: %d, %v", got, err)
	}
	im.Destroy()
	if im.Width() != 0 {
		t.Fatalf("Destroy should close image")
	}
}

func TestPaletteColorOperations(t *testing.T) {
	im := newPalette(t, 8, 8)
	defer func() { _ = im.Close() }()

	red := allocColor(t, im, 255, 0, 0)
	greenAlpha, err := im.AllocateColorAlpha(0, 255, 0, 20)
	if err != nil {
		t.Fatal(err)
	}
	if im.ColorsTotal() < 2 {
		t.Fatalf("expected allocated palette colors, got %d", im.ColorsTotal())
	}

	if c, err := im.ExactColor(255, 0, 0); err != nil || c != red {
		t.Fatalf("ExactColor: got %d, %v; want %d", c, err, red)
	}
	if c, err := im.ExactColorAlpha(0, 255, 0, 20); err != nil || c != greenAlpha {
		t.Fatalf("ExactColorAlpha: got %d, %v; want %d", c, err, greenAlpha)
	}
	if _, err := im.ClosestColor(254, 1, 1); err != nil {
		t.Fatal(err)
	}
	if _, err := im.ClosestColorAlpha(0, 250, 0, 20); err != nil {
		t.Fatal(err)
	}
	if _, err := im.ClosestColorHWB(255, 0, 0); err != nil {
		t.Fatal(err)
	}
	if _, err := im.ResolveColor(1, 2, 3); err != nil {
		t.Fatal(err)
	}
	if _, err := im.ResolveColorAlpha(4, 5, 6, 7); err != nil {
		t.Fatal(err)
	}

	if err := im.Fill(0, 0, red); err != nil {
		t.Fatal(err)
	}
	assertRGB(t, palettePixelRGBA(t, im, 0, 0), 255, 0, 0)

	if err := im.SetTransparentColor(red); err != nil {
		t.Fatal(err)
	}
	if im.TransparentColor() != red {
		t.Fatalf("unexpected transparent color: %d", im.TransparentColor())
	}

	blue := allocColor(t, im, 0, 0, 255)
	if err := im.ReplaceColor(red, blue); err != nil {
		t.Fatal(err)
	}
	assertRGB(t, palettePixelRGBA(t, im, 0, 0), 0, 0, 255)

	yellow := allocColor(t, im, 255, 255, 0)
	if err := im.ReplaceColorThreshold(blue, yellow, 0.1); err != nil {
		t.Fatal(err)
	}
	assertRGB(t, palettePixelRGBA(t, im, 0, 0), 255, 255, 0)

	if err := im.ReplaceColorArray([]Color{yellow}, []Color{red}); err != nil {
		t.Fatal(err)
	}
	assertRGB(t, palettePixelRGBA(t, im, 0, 0), 255, 0, 0)

	other := newPalette(t, 8, 8)
	defer func() { _ = other.Close() }()
	if err := other.PaletteCopyFrom(im); err != nil {
		t.Fatal(err)
	}
	if other.ColorsTotal() != im.ColorsTotal() {
		t.Fatalf("palette copy mismatch: got %d, want %d", other.ColorsTotal(), im.ColorsTotal())
	}
	if err := im.DeallocateColor(greenAlpha); err != nil {
		t.Fatal(err)
	}
}

func TestPaletteQuantizationAndColorMatch(t *testing.T) {
	a := newTrueColor(t, 8, 8)
	defer func() { _ = a.Close() }()
	b := newTrueColor(t, 8, 8)
	defer func() { _ = b.Close() }()
	red := allocColor(t, a, 255, 0, 0)
	if err := a.Fill(0, 0, red); err != nil {
		t.Fatal(err)
	}
	if err := b.Fill(0, 0, red); err != nil {
		t.Fatal(err)
	}

	if err := a.SetPaletteQuality(50, 90); err != nil {
		t.Fatal(err)
	}
	if err := a.TrueColorToPaletteWithMethod(QuantDefault, 0); err != nil && !errors.Is(err, ErrUnsupportedFeature) {
		t.Fatal(err)
	}
	if a.TrueColor() {
		if err := a.TrueColorToPalette(false, 16); err != nil {
			t.Fatal(err)
		}
	}
	if a.TrueColor() {
		t.Fatal("expected palette image after quantization")
	}
	if err := a.ColorMatch(b); err != nil && !errors.Is(err, ErrColorMatch) {
		t.Fatal(err)
	}
}

func TestBrushTileAndCopyVariants(t *testing.T) {
	dst := newTrueColor(t, 24, 24)
	defer func() { _ = dst.Close() }()
	src := newTrueColor(t, 8, 8)
	defer func() { _ = src.Close() }()
	red := allocColor(t, src, 255, 0, 0)
	blue := allocColor(t, src, 0, 0, 255)
	if err := src.Fill(0, 0, red); err != nil {
		t.Fatal(err)
	}
	if err := src.FilledRectangle(2, 2, 5, 5, blue); err != nil {
		t.Fatal(err)
	}

	if err := dst.CopyMergeGrayFrom(src, 0, 0, 0, 0, 8, 8, 50); err != nil {
		t.Fatal(err)
	}
	if err := dst.CopyRotatedFrom(src, 12, 12, 0, 0, 8, 8, 30); err != nil {
		t.Fatal(err)
	}

	brush := newTrueColor(t, 2, 2)
	defer func() { _ = brush.Close() }()
	black := allocColor(t, brush, 0, 0, 0)
	if err := brush.Fill(0, 0, black); err != nil {
		t.Fatal(err)
	}
	if err := dst.SetBrush(brush); err != nil {
		t.Fatal(err)
	}
	if dst.brush != brush {
		t.Fatal("SetBrush did not retain brush")
	}
	if err := dst.Line(0, 23, 23, 23, Brushed); err != nil {
		t.Fatal(err)
	}
	if err := dst.ClearBrush(); err != nil {
		t.Fatal(err)
	}
	if dst.brush != nil {
		t.Fatal("ClearBrush did not release brush")
	}

	tile := newTrueColor(t, 2, 2)
	defer func() { _ = tile.Close() }()
	green := allocColor(t, tile, 0, 255, 0)
	if err := tile.Fill(0, 0, green); err != nil {
		t.Fatal(err)
	}
	if err := dst.SetTile(tile); err != nil {
		t.Fatal(err)
	}
	if err := dst.FilledRectangle(16, 16, 23, 23, Tiled); err != nil {
		t.Fatal(err)
	}
	if err := dst.ClearTile(); err != nil {
		t.Fatal(err)
	}
}

func TestTransformsAndAffineImages(t *testing.T) {
	src := newTrueColor(t, 12, 12)
	defer func() { _ = src.Close() }()
	red := allocColor(t, src, 255, 0, 0)
	if err := src.Fill(0, 0, red); err != nil {
		t.Fatal(err)
	}

	if err := src.FlipHorizontal(); err != nil {
		t.Fatal(err)
	}
	if err := src.FlipVertical(); err != nil {
		t.Fatal(err)
	}
	if err := src.FlipBoth(); err != nil {
		t.Fatal(err)
	}

	m := AffineTranslate(1, 1).Concat(AffineScale(1, 1)).Concat(AffineRotate(0))
	if !m.Rectilinear() {
		t.Fatal("expected rectilinear affine")
	}
	if m.Expansion() <= 0 {
		t.Fatalf("unexpected affine expansion: %f", m.Expansion())
	}
	if _, err := m.Invert(); err != nil {
		t.Fatal(err)
	}
	bbox := m.TransformBoundingBox(Rect{X: 0, Y: 0, Width: 4, Height: 4})
	if bbox.Width <= 0 || bbox.Height <= 0 {
		t.Fatalf("unexpected affine bbox: %+v", bbox)
	}
	out, err := src.TransformAffine(Rect{X: 0, Y: 0, Width: 8, Height: 8}, m)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = out.Close() }()

	dst := newTrueColor(t, 16, 16)
	defer func() { _ = dst.Close() }()
	if err := dst.TransformAffineCopy(src, 1, 1, Rect{X: 0, Y: 0, Width: 8, Height: 8}, m); err != nil {
		t.Fatal(err)
	}
}

func TestFilterWrappers(t *testing.T) {
	ops := []struct {
		name string
		fn   func(*Image) error
	}{
		{"GaussianBlur", func(im *Image) error { return im.GaussianBlur() }},
		{"SelectiveBlur", func(im *Image) error { return im.SelectiveBlur() }},
		{"EdgeDetectQuick", func(im *Image) error { return im.EdgeDetectQuick() }},
		{"Emboss", func(im *Image) error { return im.Emboss() }},
		{"MeanRemoval", func(im *Image) error { return im.MeanRemoval() }},
		{"Smooth", func(im *Image) error { return im.Smooth(1.0) }},
		{"Negate", func(im *Image) error { return im.Negate() }},
		{"Brightness", func(im *Image) error { return im.Brightness(10) }},
		{"Contrast", func(im *Image) error { return im.Contrast(10) }},
		{"AdjustColor", func(im *Image) error { return im.AdjustColor(1, 2, 3, 0) }},
		{"Convolution", func(im *Image) error {
			return im.Convolution([3][3]float32{{0, 0, 0}, {0, 1, 0}, {0, 0, 0}}, 1, 0)
		}},
		{"Pixelate", func(im *Image) error { return im.Pixelate(2, PixelateAverage) }},
		{"Scatter", func(im *Image) error { return im.Scatter(1, 2) }},
		{"ScatterColor", func(im *Image) error {
			red := allocColor(t, im, 255, 0, 0)
			return im.ScatterColor(1, 2, []Color{red})
		}},
	}

	for _, op := range ops {
		t.Run(op.name, func(t *testing.T) {
			im := newTrueColor(t, 8, 8)
			defer func() { _ = im.Close() }()
			red := allocColor(t, im, 255, 0, 0)
			if err := im.Fill(0, 0, red); err != nil {
				t.Fatal(err)
			}
			if err := op.fn(im); err != nil && !errors.Is(err, ErrUnsupportedFeature) {
				t.Fatal(err)
			}
		})
	}
}

func TestVersionAndRuntimeFeatures(t *testing.T) {
	v := Version()
	if v.Major == 0 || v.String == "" {
		t.Fatalf("unexpected libgd version: %+v", v)
	}
	f := RuntimeFeatures()
	if !f.PNG && !f.JPEG && !f.GIF && !f.BMP {
		t.Fatalf("unexpected empty feature set: %+v", f)
	}
}

func TestFormatRoundTrips(t *testing.T) {
	base := newTrueColor(t, 10, 10)
	defer func() { _ = base.Close() }()
	white := allocColor(t, base, 255, 255, 255)
	black := allocColor(t, base, 0, 0, 0)
	if err := base.Fill(0, 0, white); err != nil {
		t.Fatal(err)
	}
	if err := base.Line(0, 0, 9, 9, black); err != nil {
		t.Fatal(err)
	}

	type formatCase struct {
		name   string
		format Format
		encode func(*Image) ([]byte, error)
		decode func([]byte) (*Image, error)
		file   func(*Image, string) error
		read   func(string) (*Image, error)
	}
	cases := []formatCase{
		{"JPEG", FormatJPEG, func(im *Image) ([]byte, error) { return im.EncodeJPEG(nil) }, func(b []byte) (*Image, error) { return DecodeJPEG(b, nil) }, func(im *Image, p string) error { return im.EncodeJPEGFile(p, nil) }, func(p string) (*Image, error) { return DecodeJPEGFile(p, nil) }},
		{"GIF", FormatGIF, func(im *Image) ([]byte, error) { return im.EncodeGIF() }, DecodeGIF, func(im *Image, p string) error { return im.EncodeGIFFile(p) }, DecodeGIFFile},
		{"BMP", FormatBMP, func(im *Image) ([]byte, error) { return im.EncodeBMP(nil) }, DecodeBMP, func(im *Image, p string) error { return im.EncodeBMPFile(p, nil) }, DecodeBMPFile},
		{"TIFF", FormatTIFF, func(im *Image) ([]byte, error) { return im.EncodeTIFF() }, DecodeTIFF, func(im *Image, p string) error { return im.EncodeTIFFFile(p) }, DecodeTIFFFile},
		{"GD", FormatGD, func(im *Image) ([]byte, error) { return im.EncodeGD() }, DecodeGD, func(im *Image, p string) error { return im.EncodeGDFile(p) }, DecodeGDFile},
		{"GD2", FormatGD2, func(im *Image) ([]byte, error) { return im.EncodeGD2(nil) }, DecodeGD2, func(im *Image, p string) error { return im.EncodeGD2File(p, nil) }, DecodeGD2File},
		{"WebP", FormatWebP, func(im *Image) ([]byte, error) { return im.EncodeWebP(nil) }, DecodeWebP, func(im *Image, p string) error { return im.EncodeWebPFile(p, nil) }, DecodeWebPFile},
		{"HEIF", FormatHEIF, func(im *Image) ([]byte, error) { return im.EncodeHEIF(nil) }, DecodeHEIF, func(im *Image, p string) error { return im.EncodeHEIFFile(p, nil) }, DecodeHEIFFile},
		{"AVIF", FormatAVIF, func(im *Image) ([]byte, error) { return im.EncodeAVIF(nil) }, DecodeAVIF, func(im *Image, p string) error { return im.EncodeAVIFFile(p, nil) }, DecodeAVIFFile},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if !SupportsFormat(tc.format, true) || !SupportsFormat(tc.format, false) {
				t.Skipf("linked libgd has no %s read/write support", tc.name)
			}
			data, err := tc.encode(base)
			if err != nil {
				if errors.Is(err, ErrEncode) || errors.Is(err, ErrUnsupportedFeature) {
					t.Skipf("%s encode advertised by libgd but unavailable at runtime: %v", tc.name, err)
				}
				t.Fatal(err)
			}
			decoded, err := tc.decode(data)
			if err != nil {
				if errors.Is(err, ErrDecode) || errors.Is(err, ErrUnsupportedFeature) {
					t.Skipf("%s decode advertised by libgd but unavailable at runtime: %v", tc.name, err)
				}
				t.Fatal(err)
			}
			defer func() { _ = decoded.Close() }()
			if decoded.Width() != 10 || decoded.Height() != 10 {
				t.Fatalf("unexpected decoded size: %dx%d", decoded.Width(), decoded.Height())
			}

			path := filepath.Join(t.TempDir(), "out."+string(tc.format))
			if tc.name == "JPEG" {
				path = filepath.Join(t.TempDir(), "out.jpg")
			}
			if err := tc.file(base, path); err != nil {
				if errors.Is(err, ErrEncode) || errors.Is(err, ErrUnsupportedFeature) {
					t.Skipf("%s file encode advertised by libgd but unavailable at runtime: %v", tc.name, err)
				}
				t.Fatal(err)
			}
			fromFile, err := tc.read(path)
			if err != nil {
				if errors.Is(err, ErrDecode) || errors.Is(err, ErrUnsupportedFeature) {
					t.Skipf("%s file decode advertised by libgd but unavailable at runtime: %v", tc.name, err)
				}
				t.Fatal(err)
			}
			defer func() { _ = fromFile.Close() }()
			if fromFile.Width() != 10 || fromFile.Height() != 10 {
				t.Fatalf("unexpected file decoded size: %dx%d", fromFile.Width(), fromFile.Height())
			}
		})
	}
}

func TestWBMPAndGenericFileRoundTrips(t *testing.T) {
	pal := newPalette(t, 8, 8)
	defer func() { _ = pal.Close() }()
	white := allocColor(t, pal, 255, 255, 255)
	black := allocColor(t, pal, 0, 0, 0)
	if err := pal.Fill(0, 0, white); err != nil {
		t.Fatal(err)
	}
	if err := pal.SetPixel(1, 1, black); err != nil {
		t.Fatal(err)
	}

	if SupportsFormat(FormatWBMP, true) && SupportsFormat(FormatWBMP, false) {
		data, err := pal.EncodeWBMP(black)
		if err != nil {
			t.Fatal(err)
		}
		decoded, err := DecodeWBMP(data)
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = decoded.Close() }()
		if decoded.Width() != 8 || decoded.Height() != 8 {
			t.Fatalf("unexpected WBMP size: %dx%d", decoded.Width(), decoded.Height())
		}

		path := filepath.Join(t.TempDir(), "out.wbmp")
		if err := pal.EncodeWBMPFile(path, black); err != nil {
			t.Fatal(err)
		}
		fileDecoded, err := DecodeWBMPFile(path)
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = fileDecoded.Close() }()
	}

	if SupportsFormat(FormatPNG, true) && SupportsFormat(FormatPNG, false) {
		path := filepath.Join(t.TempDir(), "generic.png")
		if err := pal.EncodeFile(path); err != nil {
			t.Fatal(err)
		}
		decoded, err := DecodeFile(path)
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = decoded.Close() }()
		if decoded.Width() != 8 || decoded.Height() != 8 {
			t.Fatalf("unexpected generic decoded size: %dx%d", decoded.Width(), decoded.Height())
		}
	}
}

func TestContextEncoders(t *testing.T) {
	im := newTrueColor(t, 8, 8)
	defer func() { _ = im.Close() }()
	white := allocColor(t, im, 255, 255, 255)
	if err := im.Fill(0, 0, white); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name   string
		format Format
		encode func(*Image, *IOContext) error
		decode func(*IOContext) (*Image, error)
	}{
		{"JPEG", FormatJPEG, func(im *Image, ctx *IOContext) error { return im.EncodeJPEGContext(ctx, nil) }, func(ctx *IOContext) (*Image, error) { return DecodeJPEGContext(ctx, nil) }},
		{"GIF", FormatGIF, func(im *Image, ctx *IOContext) error { return im.EncodeGIFContext(ctx) }, DecodeGIFContext},
		{"BMP", FormatBMP, func(im *Image, ctx *IOContext) error { return im.EncodeBMPContext(ctx, nil) }, DecodeBMPContext},
		{"TIFF", FormatTIFF, func(im *Image, ctx *IOContext) error { return im.EncodeTIFFContext(ctx) }, DecodeTIFFContext},
		{"WebP", FormatWebP, func(im *Image, ctx *IOContext) error { return im.EncodeWebPContext(ctx, nil) }, DecodeWebPContext},
		{"HEIF", FormatHEIF, func(im *Image, ctx *IOContext) error { return im.EncodeHEIFContext(ctx, nil) }, DecodeHEIFContext},
		{"AVIF", FormatAVIF, func(im *Image, ctx *IOContext) error { return im.EncodeAVIFContext(ctx, nil) }, DecodeAVIFContext},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if !SupportsFormat(tc.format, true) || !SupportsFormat(tc.format, false) {
				t.Skipf("linked libgd has no %s read/write support", tc.name)
			}
			out, err := NewDynamicWriteContext()
			if err != nil {
				t.Fatal(err)
			}
			if err := tc.encode(im, out); err != nil {
				_ = out.Close()
				if errors.Is(err, ErrEncode) || errors.Is(err, ErrUnsupportedFeature) {
					t.Skipf("%s context encode advertised by libgd but unavailable at runtime: %v", tc.name, err)
				}
				t.Fatal(err)
			}
			data, err := out.Extract()
			if err != nil {
				t.Fatal(err)
			}
			in, err := NewDynamicReadContext(data)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = in.Close() }()
			decoded, err := tc.decode(in)
			if err != nil {
				if errors.Is(err, ErrDecode) || errors.Is(err, ErrUnsupportedFeature) {
					t.Skipf("%s context decode advertised by libgd but unavailable at runtime: %v", tc.name, err)
				}
				t.Fatal(err)
			}
			defer func() { _ = decoded.Close() }()
			if decoded.Width() != 8 || decoded.Height() != 8 {
				t.Fatalf("unexpected decoded size: %dx%d", decoded.Width(), decoded.Height())
			}
		})
	}
}

func TestGIFAnimationFileAndCaches(t *testing.T) {
	im := newPalette(t, 8, 8)
	defer func() { _ = im.Close() }()
	black := allocColor(t, im, 0, 0, 0)
	if err := im.SetPixel(0, 0, black); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(t.TempDir(), "anim.gif")
	if err := im.GIFAnimBeginFile(path, true, 0); err != nil {
		t.Fatal(err)
	}
	if err := im.GIFAnimAddFile(path, GIFFrameOptions{Delay: 1, Disposal: GIFDisposalNone}); err != nil {
		t.Fatal(err)
	}
	if err := GIFAnimEndFile(path); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(data, []byte("GIF")) {
		t.Fatalf("unexpected GIF animation output prefix: %q", data[:min(len(data), 8)])
	}

	FontCacheSetup()
	FontCacheShutdown()
	_ = UseFontConfig(false)
}
