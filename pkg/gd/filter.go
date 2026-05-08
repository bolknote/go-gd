package gd

/*
#include <gd.h>
#include <gdfx.h>

static int go_gd_convolution(gdImagePtr im, float *values, float div, float offset) {
	float filter[3][3];
	for (int y = 0; y < 3; y++) {
		for (int x = 0; x < 3; x++) {
			filter[y][x] = values[y * 3 + x];
		}
	}
	return gdImageConvolution(im, filter, div, offset);
}

static int go_gd_image_scatter_ex(gdImagePtr im, int sub, int plus, int *colors, unsigned int num_colors, unsigned int seed) {
	gdScatter scatter;
	scatter.sub = sub;
	scatter.plus = plus;
	scatter.colors = colors;
	scatter.num_colors = num_colors;
	scatter.seed = seed;
	return gdImageScatterEx(im, &scatter);
}
*/
import "C"
import "unsafe"

type PixelateMode int

const (
	PixelateUpperLeft PixelateMode = iota
	PixelateAverage
)

type ScatterOptions struct {
	Sub, Plus int
	Colors    []Color
	Seed      uint
}

func (im *Image) Sharpen(percent int) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	C.gdImageSharpen(ptr, C.int(percent))
	return nil
}

func (im *Image) GaussianBlur() error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	if C.gdImageGaussianBlur(ptr) == 0 {
		return ErrUnsupportedFeature
	}
	return nil
}

func (im *Image) SelectiveBlur() error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	if C.gdImageSelectiveBlur(ptr) == 0 {
		return ErrUnsupportedFeature
	}
	return nil
}

func (im *Image) EdgeDetectQuick() error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	if C.gdImageEdgeDetectQuick(ptr) == 0 {
		return ErrUnsupportedFeature
	}
	return nil
}

func (im *Image) Emboss() error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	if C.gdImageEmboss(ptr) == 0 {
		return ErrUnsupportedFeature
	}
	return nil
}

func (im *Image) MeanRemoval() error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	if C.gdImageMeanRemoval(ptr) == 0 {
		return ErrUnsupportedFeature
	}
	return nil
}

func (im *Image) Smooth(weight float64) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	if C.gdImageSmooth(ptr, C.float(weight)) == 0 {
		return ErrUnsupportedFeature
	}
	return nil
}

func (im *Image) Negate() error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	if C.gdImageNegate(ptr) == 0 {
		return ErrUnsupportedFeature
	}
	return nil
}

func (im *Image) Grayscale() error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	if C.gdImageGrayScale(ptr) == 0 {
		return ErrUnsupportedFeature
	}
	return nil
}

func (im *Image) Brightness(value int) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	if C.gdImageBrightness(ptr, C.int(value)) == 0 {
		return ErrUnsupportedFeature
	}
	return nil
}

func (im *Image) Contrast(value float64) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	if C.gdImageContrast(ptr, C.double(value)) == 0 {
		return ErrUnsupportedFeature
	}
	return nil
}

func (im *Image) AdjustColor(r, g, b, a int) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	if C.gdImageColor(ptr, C.int(r), C.int(g), C.int(b), C.int(a)) == 0 {
		return ErrUnsupportedFeature
	}
	return nil
}

func (im *Image) Convolution(filter [3][3]float32, divisor, offset float32) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	values := [9]C.float{}
	for y := range filter {
		for x := range filter[y] {
			values[y*3+x] = C.float(filter[y][x])
		}
	}
	if C.go_gd_convolution(ptr, (*C.float)(unsafe.Pointer(&values[0])), C.float(divisor), C.float(offset)) == 0 {
		return ErrUnsupportedFeature
	}
	return nil
}

func (im *Image) Pixelate(blockSize int, mode PixelateMode) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	if C.gdImagePixelate(ptr, C.int(blockSize), C.uint(mode)) == 0 {
		return ErrUnsupportedFeature
	}
	return nil
}

func (im *Image) Scatter(sub, plus int) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	if C.gdImageScatter(ptr, C.int(sub), C.int(plus)) == 0 {
		return ErrUnsupportedFeature
	}
	return nil
}

func (im *Image) ScatterColor(sub, plus int, colors []Color) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	if len(colors) == 0 {
		return ErrInvalidArgument
	}
	ccolors := make([]C.int, len(colors))
	for i, c := range colors {
		ccolors[i] = C.int(c)
	}
	if C.gdImageScatterColor(ptr, C.int(sub), C.int(plus), (*C.int)(unsafe.Pointer(&ccolors[0])), C.uint(len(ccolors))) == 0 {
		return ErrUnsupportedFeature
	}
	return nil
}

func (im *Image) ScatterEx(opts ScatterOptions) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	var (
		colorsPtr *C.int
		numColors = C.uint(len(opts.Colors))
	)
	if len(opts.Colors) > 0 {
		ccolors := make([]C.int, len(opts.Colors))
		for i, c := range opts.Colors {
			ccolors[i] = C.int(c)
		}
		colorsPtr = (*C.int)(unsafe.Pointer(&ccolors[0]))
	}
	if C.go_gd_image_scatter_ex(ptr, C.int(opts.Sub), C.int(opts.Plus), colorsPtr, numColors, C.uint(opts.Seed)) == 0 {
		return ErrUnsupportedFeature
	}
	return nil
}
