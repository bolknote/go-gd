package gd

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func newTrueColor(t *testing.T, w, h int) *Image {
	t.Helper()
	im, err := NewTrueColor(w, h)
	if err != nil {
		t.Fatal(err)
	}
	return im
}

func newPalette(t *testing.T, w, h int) *Image {
	t.Helper()
	im, err := NewPalette(w, h)
	if err != nil {
		t.Fatal(err)
	}
	return im
}

func allocColor(t *testing.T, im *Image, r, g, b int) Color {
	t.Helper()
	c, err := im.AllocateColor(r, g, b)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestEncodePNGFileRoundTrip(t *testing.T) {
	if !SupportsFormat(FormatPNG, true) || !SupportsFormat(FormatPNG, false) {
		t.Skip("linked libgd has no PNG support")
	}
	im := newTrueColor(t, 8, 8)
	defer func() { _ = im.Close() }()

	white := allocColor(t, im, 255, 255, 255)
	if err := im.Fill(0, 0, white); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "out.png")
	if err := im.EncodePNGFile(path, nil); err != nil {
		t.Fatalf("EncodePNGFile: %v", err)
	}
	decoded, err := DecodePNGFile(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = decoded.Close() }()
	if decoded.Width() != 8 || decoded.Height() != 8 {
		t.Fatalf("unexpected decoded size: %dx%d", decoded.Width(), decoded.Height())
	}
}

func TestEncodeFileToUnwritablePathReturnsError(t *testing.T) {
	if !SupportsFormat(FormatPNG, true) {
		t.Skip("linked libgd has no PNG support")
	}
	im := newTrueColor(t, 2, 2)
	defer func() { _ = im.Close() }()
	err := im.EncodePNGFile("/nonexistent-dir-go-gd-tests/out.png", nil)
	if err == nil {
		t.Fatalf("expected error writing to unwritable path")
	}
}

func TestCopyAndTransform(t *testing.T) {
	src := newTrueColor(t, 8, 8)
	defer func() { _ = src.Close() }()
	red := allocColor(t, src, 255, 0, 0)
	if err := src.Fill(0, 0, red); err != nil {
		t.Fatal(err)
	}

	dst := newTrueColor(t, 16, 16)
	defer func() { _ = dst.Close() }()
	if err := dst.CopyFrom(src, 0, 0, 0, 0, 8, 8); err != nil {
		t.Fatal(err)
	}
	if err := dst.CopyResizedFrom(src, 8, 0, 0, 0, 8, 8, 8, 8); err != nil {
		t.Fatal(err)
	}
	if err := dst.CopyResampledFrom(src, 0, 8, 0, 0, 8, 8, 8, 8); err != nil {
		t.Fatal(err)
	}
	if err := dst.CopyMergeFrom(src, 8, 8, 0, 0, 8, 8, 50); err != nil {
		t.Fatal(err)
	}
}

func TestRotateInterpolated(t *testing.T) {
	im := newTrueColor(t, 8, 8)
	defer func() { _ = im.Close() }()
	red := allocColor(t, im, 255, 0, 0)
	if err := im.Fill(0, 0, red); err != nil {
		t.Fatal(err)
	}
	if err := im.SetInterpolationMethod(InterpolationBilinearFixed); err != nil {
		t.Fatal(err)
	}
	rotated, err := im.RotateInterpolated(45, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = rotated.Close() }()
	if rotated.Width() <= 0 || rotated.Height() <= 0 {
		t.Fatalf("unexpected rotated size: %dx%d", rotated.Width(), rotated.Height())
	}
}

func TestPaletteAndTrueColorRoundTrip(t *testing.T) {
	im := newPalette(t, 8, 8)
	defer func() { _ = im.Close() }()
	if im.TrueColor() {
		t.Fatal("NewPalette must not be true-color")
	}
	if err := im.PaletteToTrueColor(); err != nil {
		t.Fatalf("PaletteToTrueColor: %v", err)
	}
	if !im.TrueColor() {
		t.Fatal("expected true-color after PaletteToTrueColor")
	}

	tc := newTrueColor(t, 8, 8)
	defer func() { _ = tc.Close() }()
	pal, err := tc.CreatePaletteFromTrueColor(true, 16)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = pal.Close() }()
	if pal.TrueColor() {
		t.Fatal("expected palette image, got true-color")
	}
}

func TestCropValidation(t *testing.T) {
	im := newTrueColor(t, 4, 4)
	defer func() { _ = im.Close() }()
	if _, err := im.CropThreshold(-1, 1.0); !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("expected ErrInvalidArgument for negative color, got %v", err)
	}
	if _, err := im.Crop(Rect{}); !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("expected ErrInvalidArgument for empty crop rect, got %v", err)
	}
	white := allocColor(t, im, 255, 255, 255)
	red := allocColor(t, im, 255, 0, 0)
	if err := im.Fill(0, 0, white); err != nil {
		t.Fatal(err)
	}
	if err := im.FilledRectangle(1, 1, 2, 2, red); err != nil {
		t.Fatal(err)
	}
	cropped, err := im.CropAuto(CropDefault)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = cropped.Close() }()
	if cropped.Width() <= 0 || cropped.Height() <= 0 {
		t.Fatalf("unexpected auto-crop size: %dx%d", cropped.Width(), cropped.Height())
	}
}

func TestStringEmptyAccepted(t *testing.T) {
	font, err := BuiltinFont(FontTiny)
	if err != nil {
		t.Fatal(err)
	}
	im := newPalette(t, 8, 8)
	defer func() { _ = im.Close() }()
	black := allocColor(t, im, 0, 0, 0)
	if err := im.String(font, 0, 0, "", black); err != nil {
		t.Fatalf("empty String must be a no-op, got %v", err)
	}
	if err := im.String16(font, 0, 0, nil, black); err != nil {
		t.Fatalf("empty String16 must be a no-op, got %v", err)
	}
	if err := im.StringUp16(font, 0, 0, []uint16{}, black); err != nil {
		t.Fatalf("empty StringUp16 must be a no-op, got %v", err)
	}
}

func TestBitmapTextAndFreeTypeErrors(t *testing.T) {
	font, err := BuiltinFont(FontSmall)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := BuiltinFont(FontSize(99)); !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("expected invalid font size, got %v", err)
	}

	im := newPalette(t, 32, 32)
	defer func() { _ = im.Close() }()
	black := allocColor(t, im, 0, 0, 0)
	if err := im.Char(font, 1, 1, 'A', black); err != nil {
		t.Fatal(err)
	}
	if err := im.CharUp(font, 10, 20, 'B', black); err != nil {
		t.Fatal(err)
	}
	if err := im.StringUp(font, 20, 20, "gd", black); err != nil {
		t.Fatal(err)
	}

	_, err = im.StringFT(black, "/definitely/not/a/font.ttf", 12, 0, 1, 1, "gd")
	var ftErr FreeTypeError
	if !errors.As(err, &ftErr) || ftErr.Error() == "" {
		t.Fatalf("expected FreeTypeError from invalid font, got %T: %v", err, err)
	}
	if err := im.StringFTCircle(16, 16, 10, 8, 0.5, "/definitely/not/a/font.ttf", 12, "top", "bottom", black); err == nil {
		t.Fatal("expected StringFTCircle error for invalid font")
	}
}

func TestSetTileKeepsTileAlive(t *testing.T) {
	im := newTrueColor(t, 16, 16)
	defer func() { _ = im.Close() }()
	tile := newTrueColor(t, 4, 4)
	if err := im.SetTile(tile); err != nil {
		t.Fatal(err)
	}
	if im.tile != tile {
		t.Fatal("SetTile did not retain reference to tile")
	}
	if err := im.ClearTile(); err != nil {
		t.Fatal(err)
	}
	if im.tile != nil {
		t.Fatal("ClearTile did not drop reference")
	}
	if err := tile.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestEncodeXBMRoundTrip(t *testing.T) {
	im := newPalette(t, 8, 8)
	defer func() { _ = im.Close() }()
	black := allocColor(t, im, 0, 0, 0)
	white := allocColor(t, im, 255, 255, 255)
	if err := im.Fill(0, 0, white); err != nil {
		t.Fatal(err)
	}
	if err := im.SetPixel(0, 0, black); err != nil {
		t.Fatal(err)
	}

	data, err := im.EncodeXBM("test_image", black)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(data, []byte("test_image_width")) {
		t.Fatalf("XBM output missing expected marker: %q", data)
	}

	path := filepath.Join(t.TempDir(), "out.xbm")
	if err := im.EncodeXBMFile(path, "tx", black); err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeXBMFile(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = decoded.Close() }()
	if decoded.Width() != 8 || decoded.Height() != 8 {
		t.Fatalf("unexpected decoded XBM size: %dx%d", decoded.Width(), decoded.Height())
	}
}

func TestColorRGBARoundTrip(t *testing.T) {
	im := newTrueColor(t, 2, 2)
	defer func() { _ = im.Close() }()
	c := TrueColorAlpha(10, 20, 30, 40)
	rgba, err := im.ColorRGBA(c)
	if err != nil {
		t.Fatal(err)
	}
	if rgba.R != 10 || rgba.G != 20 || rgba.B != 30 || rgba.A != 40 {
		t.Fatalf("ColorRGBA roundtrip mismatch: %+v", rgba)
	}
}

func TestClosedImageGettersZero(t *testing.T) {
	im := newTrueColor(t, 4, 4)
	if err := im.Close(); err != nil {
		t.Fatal(err)
	}
	if im.Width() != 0 || im.Height() != 0 {
		t.Fatalf("closed image must report zero size")
	}
	if im.TrueColor() || im.Interlaced() {
		t.Fatalf("closed image must report false flags")
	}
	if im.TransparentColor() != Color(-1) {
		t.Fatalf("closed image must report TransparentColor=-1")
	}
	if x, y := im.Resolution(); x != 0 || y != 0 {
		t.Fatalf("closed image must report zero resolution, got %d,%d", x, y)
	}
	if err := im.Close(); err != nil {
		t.Fatalf("double Close should be no-op, got %v", err)
	}
}

func TestImageConcurrentUseWithExternalSync(t *testing.T) {
	im := newTrueColor(t, 16, 16)
	defer func() { _ = im.Close() }()
	black := allocColor(t, im, 0, 0, 0)

	var (
		mu sync.Mutex
		wg sync.WaitGroup
	)
	for i := 0; i < 16; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()
			if err := im.SetPixel(i%16, i%16, black); err != nil {
				t.Errorf("SetPixel: %v", err)
				return
			}
			_, _ = im.Width(), im.Height()
		}()
	}
	wg.Wait()
}

func TestIOContextConcurrentUseWithExternalSync(t *testing.T) {
	ctx, err := NewDynamicWriteContext()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ctx.Close() }()

	var (
		mu sync.Mutex
		wg sync.WaitGroup
	)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()
			if _, err := ctx.cptr(); err != nil {
				t.Errorf("cptr: %v", err)
			}
		}()
	}
	wg.Wait()
}

func TestArgumentValidationEdgeCases(t *testing.T) {
	if _, err := NewTrueColor(0, 10); !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("NewTrueColor(0, 10): expected ErrInvalidArgument, got %v", err)
	}
	if _, err := NewTrueColor(10, 0); !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("NewTrueColor(10, 0): expected ErrInvalidArgument, got %v", err)
	}
	if _, err := NewPalette(-1, 8); !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("NewPalette(-1, 8): expected ErrInvalidArgument, got %v", err)
	}

	im := newTrueColor(t, 4, 4)
	defer func() { _ = im.Close() }()

	if err := im.SetClip(Rect{X: 0, Y: 0, Width: 0, Height: 1}); !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("SetClip zero-width: expected ErrInvalidArgument, got %v", err)
	}
	if err := im.SetStyle(nil); !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("SetStyle(nil): expected ErrInvalidArgument, got %v", err)
	}
	if err := im.ReplaceColorArray([]Color{1}, []Color{2, 3}); !errors.Is(err, ErrInvalidArgument) {
		t.Fatalf("ReplaceColorArray mismatched lengths: expected ErrInvalidArgument, got %v", err)
	}
}

func TestWriteToReadOnlyDirectoryFails(t *testing.T) {
	if !SupportsFormat(FormatPNG, true) {
		t.Skip("linked libgd has no PNG write support")
	}
	im := newTrueColor(t, 2, 2)
	defer func() { _ = im.Close() }()

	dir := t.TempDir()
	if err := os.Chmod(dir, 0o500); err != nil {
		t.Fatalf("chmod temp dir: %v", err)
	}
	defer func() { _ = os.Chmod(dir, 0o700) }()

	err := im.EncodePNGFile(filepath.Join(dir, "out.png"), nil)
	if err == nil {
		t.Fatalf("expected write error for read-only directory")
	}
}

func TestRuntimeFeaturesMatchSupportsFormatReadSide(t *testing.T) {
	f := RuntimeFeatures()
	cases := []struct {
		name string
		got  bool
		want bool
	}{
		{"PNG", f.PNG, SupportsFormat(FormatPNG, false)},
		{"JPEG", f.JPEG, SupportsFormat(FormatJPEG, false)},
		{"GIF", f.GIF, SupportsFormat(FormatGIF, false)},
		{"WebP", f.WebP, SupportsFormat(FormatWebP, false)},
		{"WBMP", f.WBMP, SupportsFormat(FormatWBMP, false)},
		{"BMP", f.BMP, SupportsFormat(FormatBMP, false)},
		{"TGA", f.TGA, SupportsFormat(FormatTGA, false)},
		{"TIFF", f.TIFF, SupportsFormat(FormatTIFF, false)},
		{"GD", f.GD, SupportsFormat(FormatGD, false)},
		{"GD2", f.GD2, SupportsFormat(FormatGD2, false)},
		{"HEIF", f.HEIF, SupportsFormat(FormatHEIF, false)},
		{"AVIF", f.AVIF, SupportsFormat(FormatAVIF, false)},
		{"XBM", f.XBM, SupportsFormat(FormatXBM, false)},
		{"XPM", f.XPM, SupportsFormat(FormatXPM, false)},
	}
	for _, tc := range cases {
		if tc.got != tc.want {
			t.Fatalf("%s mismatch: RuntimeFeatures=%v, SupportsFormat=%v", tc.name, tc.got, tc.want)
		}
	}
}
