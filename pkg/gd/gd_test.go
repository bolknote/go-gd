package gd

import (
	"math"
	"path/filepath"
	"testing"
)

func TestCreateDrawEncodeDecodePNG(t *testing.T) {
	if !SupportsFormat(FormatPNG, true) || !SupportsFormat(FormatPNG, false) {
		t.Skip("linked libgd has no PNG support")
	}

	im, err := NewTrueColor(16, 16)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = im.Close() }()

	if im.Width() != 16 || im.Height() != 16 {
		t.Fatalf("unexpected size: %dx%d", im.Width(), im.Height())
	}

	white, err := im.AllocateColor(255, 255, 255)
	if err != nil {
		t.Fatal(err)
	}
	black, err := im.AllocateColor(0, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if err := im.Fill(0, 0, white); err != nil {
		t.Fatal(err)
	}
	if err := im.Line(0, 0, 15, 15, black); err != nil {
		t.Fatal(err)
	}

	data, err := im.EncodePNG(nil)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodePNG(data)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = decoded.Close() }()
	if decoded.Width() != 16 || decoded.Height() != 16 {
		t.Fatalf("unexpected decoded size: %dx%d", decoded.Width(), decoded.Height())
	}
}

func TestDecodeFixtureJPEG(t *testing.T) {
	if !SupportsFormat(FormatJPEG, false) {
		t.Skip("linked libgd has no JPEG support")
	}
	img, err := DecodeJPEGFile(filepath.Join("..", "..", "testdata", "images", "source.jpg"), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = img.Close() }()
	if img.Width() <= 0 || img.Height() <= 0 {
		t.Fatalf("unexpected fixture size: %dx%d", img.Width(), img.Height())
	}
}

func TestFiltersTransformCropCompare(t *testing.T) {
	im, err := NewTrueColor(10, 10)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = im.Close() }()

	red, err := im.AllocateColor(255, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if err := im.Fill(0, 0, red); err != nil {
		t.Fatal(err)
	}
	if err := im.Grayscale(); err != nil {
		t.Fatal(err)
	}

	scaled, err := im.Scale(5, 5)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = scaled.Close() }()
	if scaled.Width() != 5 || scaled.Height() != 5 {
		t.Fatalf("unexpected scaled size: %dx%d", scaled.Width(), scaled.Height())
	}

	cropped, err := im.Crop(Rect{X: 1, Y: 1, Width: 4, Height: 4})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = cropped.Close() }()
	if cropped.Width() != 4 || cropped.Height() != 4 {
		t.Fatalf("unexpected cropped size: %dx%d", cropped.Width(), cropped.Height())
	}

	clone, err := im.Clone()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = clone.Close() }()
	flags, err := Compare(im, clone)
	if err != nil {
		t.Fatal(err)
	}
	if flags != 0 {
		t.Fatalf("clone differs from source: %d", flags)
	}
}

func TestTextAndGIFAnimation(t *testing.T) {
	font, err := BuiltinFont(FontSmall)
	if err != nil {
		t.Fatal(err)
	}
	im, err := NewPalette(20, 20)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = im.Close() }()

	black, err := im.AllocateColor(0, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if err := im.String(font, 1, 1, "gd", black); err != nil {
		t.Fatal(err)
	}

	begin, err := im.GIFAnimBegin(true, 0)
	if err != nil {
		t.Fatal(err)
	}
	frame, err := im.GIFAnimAdd(GIFFrameOptions{Delay: 10, Disposal: GIFDisposalNone})
	if err != nil {
		t.Fatal(err)
	}
	end, err := GIFAnimEnd()
	if err != nil {
		t.Fatal(err)
	}
	if len(begin) == 0 || len(frame) == 0 || len(end) == 0 {
		t.Fatal("empty gif animation chunk")
	}
}

func TestIOContextRoundTrip(t *testing.T) {
	if !SupportsFormat(FormatPNG, true) || !SupportsFormat(FormatPNG, false) {
		t.Skip("linked libgd has no PNG support")
	}

	im, err := NewTrueColor(8, 8)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = im.Close() }()
	white, err := im.AllocateColor(255, 255, 255)
	if err != nil {
		t.Fatal(err)
	}
	if err := im.Fill(0, 0, white); err != nil {
		t.Fatal(err)
	}

	out, err := NewDynamicWriteContext()
	if err != nil {
		t.Fatal(err)
	}
	if err := im.EncodePNGContext(out, nil); err != nil {
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
	decoded, err := DecodePNGContext(in)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = decoded.Close() }()
	if decoded.Width() != 8 || decoded.Height() != 8 {
		t.Fatalf("unexpected ctx decoded size: %dx%d", decoded.Width(), decoded.Height())
	}
}

func TestAdditionalCoreWrappers(t *testing.T) {
	im, err := NewTrueColor(12, 12)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = im.Close() }()

	red := TrueColorAlpha(255, 0, 0, 0)
	blue := TrueColorAlpha(0, 0, 255, 20)
	_ = AlphaBlend(red, blue)
	_ = LayerOverlay(red, blue)
	_ = LayerMultiply(red, blue)

	SetErrorMethodDiscard()
	ClearErrorMethod()

	if err := im.AlphaBlendingEffect(EffectOverlay); err != nil {
		t.Fatal(err)
	}
	if err := im.ScatterEx(ScatterOptions{Sub: 1, Plus: 2, Seed: 1}); err != nil {
		t.Fatal(err)
	}
	if _, err := im.CopyGaussianBlurred(1, 1.0); err != nil {
		t.Fatal(err)
	}
	if _, err := im.SquareToCircle(6); err != nil {
		t.Fatal(err)
	}
	if err := im.Sharpen(10); err != nil {
		t.Fatal(err)
	}
	if _, err := im.NeuQuant(16, 1); err != nil {
		t.Fatal(err)
	}
}

func TestText16AndAffineExtras(t *testing.T) {
	font, err := BuiltinFont(FontSmall)
	if err != nil {
		t.Fatal(err)
	}
	im, err := NewPalette(32, 32)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = im.Close() }()
	black, err := im.AllocateColor(0, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if err := im.String16(font, 1, 1, []uint16{'g', 'd'}, black); err != nil {
		t.Fatal(err)
	}
	if err := im.StringUp16(font, 20, 20, []uint16{'g', 'd'}, black); err != nil {
		t.Fatal(err)
	}
	FreeFontCache()

	identity := AffineIdentity()
	point := identity.ApplyToPoint(PointF{X: 2, Y: 3})
	if point.X != 2 || point.Y != 3 {
		t.Fatalf("unexpected affine point: %+v", point)
	}
	_ = identity.Flip(true, false)
	if _, err := identity.ShearHorizontal(0.1); err != nil {
		t.Fatal(err)
	}
	if _, err := identity.ShearVertical(0.1); err != nil {
		t.Fatal(err)
	}
	_ = identity.TransformBoundingBox(Rect{X: 0, Y: 0, Width: 4, Height: 4})
}

func TestAffineShearUsesReceiver(t *testing.T) {
	scale := AffineScale(2, 3)

	horizontalShear, err := affineShearHorizontal(0.25)
	if err != nil {
		t.Fatal(err)
	}
	expectedHorizontal := scale.Concat(horizontalShear)
	actualHorizontal, err := scale.ShearHorizontal(0.25)
	if err != nil {
		t.Fatal(err)
	}
	if !actualHorizontal.Equal(expectedHorizontal) {
		t.Fatalf("horizontal shear did not compose with receiver: got %+v, want %+v", actualHorizontal, expectedHorizontal)
	}

	verticalShear, err := affineShearVertical(math.Pi / 12)
	if err != nil {
		t.Fatal(err)
	}
	expectedVertical := scale.Concat(verticalShear)
	actualVertical, err := scale.ShearVertical(math.Pi / 12)
	if err != nil {
		t.Fatal(err)
	}
	if !actualVertical.Equal(expectedVertical) {
		t.Fatalf("vertical shear did not compose with receiver: got %+v, want %+v", actualVertical, expectedVertical)
	}
}

func TestGIFAnimationContext(t *testing.T) {
	im, err := NewPalette(8, 8)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = im.Close() }()

	ctx, err := NewDynamicWriteContext()
	if err != nil {
		t.Fatal(err)
	}
	if err := im.GIFAnimBeginContext(ctx, true, 0); err != nil {
		t.Fatal(err)
	}
	if err := im.GIFAnimAddContext(ctx, GIFFrameOptions{Delay: 5, Disposal: GIFDisposalNone}); err != nil {
		t.Fatal(err)
	}
	if err := GIFAnimEndContext(ctx); err != nil {
		t.Fatal(err)
	}
	data, err := ctx.Extract()
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Fatal("empty gif animation ctx output")
	}
}

func TestNegativeCases(t *testing.T) {
	if _, err := NewTrueColor(0, 1); err != ErrInvalidArgument {
		t.Fatalf("expected invalid argument, got %v", err)
	}
	if _, err := DecodePNG(nil); err != ErrInvalidArgument {
		t.Fatalf("expected invalid argument, got %v", err)
	}
	im, err := NewTrueColor(1, 1)
	if err != nil {
		t.Fatal(err)
	}
	if err := im.Close(); err != nil {
		t.Fatal(err)
	}
	if err := im.SetPixel(0, 0, 0); err != ErrClosedImage {
		t.Fatalf("expected closed image, got %v", err)
	}
}

func TestOptionDefaultsAndValidation(t *testing.T) {
	if MaxPaletteColors != 256 {
		t.Fatalf("unexpected palette size: %d", MaxPaletteColors)
	}

	pngDefaults := defaultPNGOptions(nil)
	if pngDefaults.Compression != -1 {
		t.Fatalf("unexpected png default: %+v", pngDefaults)
	}

	heifDefaults := defaultHEIFOptions(nil)
	if heifDefaults.Quality != -1 || heifDefaults.Codec != HEIFCodecHEVC || heifDefaults.Chroma != HEIFChroma444 {
		t.Fatalf("unexpected heif defaults for nil opts: %+v", heifDefaults)
	}
	heifFromZero := defaultHEIFOptions(&HEIFOptions{})
	if heifFromZero.Quality != 0 || heifFromZero.Codec != HEIFCodecHEVC || heifFromZero.Chroma != HEIFChroma444 {
		t.Fatalf("unexpected heif defaults for HEIFOptions{}: %+v", heifFromZero)
	}

	font, err := BuiltinFont(FontSmall)
	if err != nil {
		t.Fatal(err)
	}
	im, err := NewPalette(8, 8)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = im.Close() }()
	black, err := im.AllocateColor(0, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if err := im.String(font, 0, 0, "привет", black); err != ErrInvalidArgument {
		t.Fatalf("expected invalid argument for non-Latin-1 text, got %v", err)
	}
	if err := im.SetClip(Rect{Width: 0, Height: 1}); err != ErrInvalidArgument {
		t.Fatalf("expected invalid clip, got %v", err)
	}
}

func TestIOContextExtractRejectsReadContext(t *testing.T) {
	ctx, err := NewDynamicReadContext([]byte("not an image"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ctx.Extract(); err != ErrInvalidArgument {
		t.Fatalf("expected invalid extract, got %v", err)
	}
	if err := ctx.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestCompareFlagsHas(t *testing.T) {
	flags := CompareImage | CompareColor
	if !flags.Has(CompareImage) || !flags.Has(CompareColor) {
		t.Fatalf("expected flags to contain image and color: %d", flags)
	}
	if flags.Has(CompareSizeX) {
		t.Fatalf("did not expect size flag in %d", flags)
	}
	if flags.Has(0) {
		t.Fatalf("Has(0) must be false even on non-zero flags: %d", flags)
	}
	if CompareFlag(0).Has(0) {
		t.Fatalf("Has(0) must be false on zero flags")
	}
}
