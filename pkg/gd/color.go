package gd

/*
#include <gd.h>

static int go_gd_true_color_alpha(int r, int g, int b, int a) { return gdTrueColorAlpha(r, g, b, a); }
static int go_gd_true_color_get_alpha(int c) { return gdTrueColorGetAlpha(c); }
static int go_gd_true_color_get_red(int c) { return gdTrueColorGetRed(c); }
static int go_gd_true_color_get_green(int c) { return gdTrueColorGetGreen(c); }
static int go_gd_true_color_get_blue(int c) { return gdTrueColorGetBlue(c); }
static int go_gd_image_red(gdImagePtr im, int c) { return gdImageRed(im, c); }
static int go_gd_image_green(gdImagePtr im, int c) { return gdImageGreen(im, c); }
static int go_gd_image_blue(gdImagePtr im, int c) { return gdImageBlue(im, c); }
static int go_gd_image_alpha(gdImagePtr im, int c) { return gdImageAlpha(im, c); }
*/
import "C"
import "unsafe"

type QuantizationMethod int

const (
	QuantDefault QuantizationMethod = iota
	QuantJQuant
	QuantNeuQuant
	QuantLIQ
)

// TrueColorAlpha composes a truecolor value using libgd channel semantics.
func TrueColorAlpha(r, g, b, a int) Color {
	return Color(C.go_gd_true_color_alpha(C.int(r), C.int(g), C.int(b), C.int(a)))
}

// DecomposeTrueColor decomposes a truecolor value.
func DecomposeTrueColor(color Color) RGBA {
	c := C.int(color)
	return RGBA{
		R: int(C.go_gd_true_color_get_red(c)),
		G: int(C.go_gd_true_color_get_green(c)),
		B: int(C.go_gd_true_color_get_blue(c)),
		A: int(C.go_gd_true_color_get_alpha(c)),
	}
}

func (im *Image) AllocateColor(r, g, b int) (Color, error) {
	ptr, err := im.cptr()
	if err != nil {
		return 0, err
	}
	c := C.gdImageColorAllocate(ptr, C.int(r), C.int(g), C.int(b))
	if c < 0 {
		return 0, ErrPaletteFull
	}
	return Color(c), nil
}

func (im *Image) AllocateColorAlpha(r, g, b, a int) (Color, error) {
	ptr, err := im.cptr()
	if err != nil {
		return 0, err
	}
	c := C.gdImageColorAllocateAlpha(ptr, C.int(r), C.int(g), C.int(b), C.int(a))
	if c < 0 {
		return 0, ErrPaletteFull
	}
	return Color(c), nil
}

func (im *Image) ClosestColor(r, g, b int) (Color, error) {
	ptr, err := im.cptr()
	if err != nil {
		return 0, err
	}
	return Color(C.gdImageColorClosest(ptr, C.int(r), C.int(g), C.int(b))), nil
}

func (im *Image) ClosestColorAlpha(r, g, b, a int) (Color, error) {
	ptr, err := im.cptr()
	if err != nil {
		return 0, err
	}
	return Color(C.gdImageColorClosestAlpha(ptr, C.int(r), C.int(g), C.int(b), C.int(a))), nil
}

func (im *Image) ClosestColorHWB(r, g, b int) (Color, error) {
	ptr, err := im.cptr()
	if err != nil {
		return 0, err
	}
	return Color(C.gdImageColorClosestHWB(ptr, C.int(r), C.int(g), C.int(b))), nil
}

func (im *Image) ExactColor(r, g, b int) (Color, error) {
	ptr, err := im.cptr()
	if err != nil {
		return 0, err
	}
	return Color(C.gdImageColorExact(ptr, C.int(r), C.int(g), C.int(b))), nil
}

func (im *Image) ExactColorAlpha(r, g, b, a int) (Color, error) {
	ptr, err := im.cptr()
	if err != nil {
		return 0, err
	}
	return Color(C.gdImageColorExactAlpha(ptr, C.int(r), C.int(g), C.int(b), C.int(a))), nil
}

func (im *Image) ResolveColor(r, g, b int) (Color, error) {
	ptr, err := im.cptr()
	if err != nil {
		return 0, err
	}
	return Color(C.gdImageColorResolve(ptr, C.int(r), C.int(g), C.int(b))), nil
}

func (im *Image) ResolveColorAlpha(r, g, b, a int) (Color, error) {
	ptr, err := im.cptr()
	if err != nil {
		return 0, err
	}
	return Color(C.gdImageColorResolveAlpha(ptr, C.int(r), C.int(g), C.int(b), C.int(a))), nil
}

func (im *Image) DeallocateColor(color Color) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	C.gdImageColorDeallocate(ptr, C.int(color))
	return nil
}

func (im *Image) ColorRGBA(color Color) (RGBA, error) {
	ptr, err := im.cptr()
	if err != nil {
		return RGBA{}, err
	}
	if im.TrueColor() {
		return DecomposeTrueColor(color), nil
	}
	c := C.int(color)
	return RGBA{
		R: int(C.go_gd_image_red(ptr, c)),
		G: int(C.go_gd_image_green(ptr, c)),
		B: int(C.go_gd_image_blue(ptr, c)),
		A: int(C.go_gd_image_alpha(ptr, c)),
	}, nil
}

func (im *Image) SetTransparentColor(color Color) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	C.gdImageColorTransparent(ptr, C.int(color))
	return nil
}

func (im *Image) PaletteCopyFrom(src *Image) error {
	dstPtr, err := im.cptr()
	if err != nil {
		return err
	}
	srcPtr, err := src.cptr()
	if err != nil {
		return err
	}
	C.gdImagePaletteCopy(dstPtr, srcPtr)
	return nil
}

func (im *Image) PaletteToTrueColor() error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	if C.gdImagePaletteToTrueColor(ptr) == 0 {
		return ErrUnsupportedFeature
	}
	return nil
}

func (im *Image) TrueColorToPalette(dither bool, colorsWanted int) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	if C.gdImageTrueColorToPalette(ptr, boolInt(dither), C.int(colorsWanted)) == 0 {
		return ErrInvalidArgument
	}
	return nil
}

func (im *Image) TrueColorToPaletteWithMethod(method QuantizationMethod, speed int) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	if C.gdImageTrueColorToPaletteSetMethod(ptr, C.int(method), C.int(speed)) == 0 {
		return ErrUnsupportedFeature
	}
	return nil
}

func (im *Image) SetPaletteQuality(minQuality, maxQuality int) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	C.gdImageTrueColorToPaletteSetQuality(ptr, C.int(minQuality), C.int(maxQuality))
	return nil
}

func (im *Image) CreatePaletteFromTrueColor(dither bool, colorsWanted int) (*Image, error) {
	ptr, err := im.cptr()
	if err != nil {
		return nil, err
	}
	return wrapImage(C.gdImageCreatePaletteFromTrueColor(ptr, boolInt(dither), C.int(colorsWanted)))
}

func (im *Image) ColorMatch(other *Image) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	otherPtr, err := other.cptr()
	if err != nil {
		return err
	}
	if C.gdImageColorMatch(ptr, otherPtr) != 0 {
		return ErrColorMatch
	}
	return nil
}

func (im *Image) ReplaceColor(src, dst Color) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	if C.gdImageColorReplace(ptr, C.int(src), C.int(dst)) < 0 {
		return ErrInvalidArgument
	}
	return nil
}

func (im *Image) ReplaceColorThreshold(src, dst Color, threshold float32) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	if C.gdImageColorReplaceThreshold(ptr, C.int(src), C.int(dst), C.float(threshold)) < 0 {
		return ErrInvalidArgument
	}
	return nil
}

func (im *Image) ReplaceColorArray(src, dst []Color) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	if len(src) != len(dst) {
		return ErrInvalidArgument
	}
	if len(src) == 0 {
		return nil
	}
	srcC := make([]C.int, len(src))
	dstC := make([]C.int, len(dst))
	for i := range src {
		srcC[i], dstC[i] = C.int(src[i]), C.int(dst[i])
	}
	if C.gdImageColorReplaceArray(ptr, C.int(len(srcC)), (*C.int)(unsafe.Pointer(&srcC[0])), (*C.int)(unsafe.Pointer(&dstC[0]))) < 0 {
		return ErrInvalidArgument
	}
	return nil
}

func boolInt(v bool) C.int {
	if v {
		return 1
	}
	return 0
}
