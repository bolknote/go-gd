package gd

/*
#include <stdlib.h>
#include <gd.h>
#include <gdfx.h>

static void go_gd_discard_error(int priority, const char *format, va_list args) {
	(void)priority;
	(void)format;
	(void)args;
}

static void go_gd_set_error_method_discard(void) {
	gdSetErrorMethod(go_gd_discard_error);
}
*/
import "C"
import (
	"unsafe"
)

type AlphaEffect int

const (
	EffectReplace    AlphaEffect = 0
	EffectAlphaBlend AlphaEffect = 1
	EffectNormal     AlphaEffect = 2
	EffectOverlay    AlphaEffect = 3
	EffectMultiply   AlphaEffect = 4
)

func AlphaBlend(dest, src Color) Color {
	return Color(C.gdAlphaBlend(C.int(dest), C.int(src)))
}

func LayerOverlay(dest, src Color) Color {
	return Color(C.gdLayerOverlay(C.int(dest), C.int(src)))
}

func LayerMultiply(dest, src Color) Color {
	return Color(C.gdLayerMultiply(C.int(dest), C.int(src)))
}

func SetErrorMethodDiscard() {
	C.go_gd_set_error_method_discard()
}

func ClearErrorMethod() {
	C.gdClearErrorMethod()
}

func (im *Image) AlphaBlendingEffect(effect AlphaEffect) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	C.gdImageAlphaBlending(ptr, C.int(effect))
	return nil
}

func (im *Image) NeuQuant(maxColors, sampleFactor int) (*Image, error) {
	ptr, err := im.cptr()
	if err != nil {
		return nil, err
	}
	return wrapImage(C.gdImageNeuQuant(ptr, C.int(maxColors), C.int(sampleFactor)))
}

func (im *Image) CopyGaussianBlurred(radius int, sigma float64) (*Image, error) {
	ptr, err := im.cptr()
	if err != nil {
		return nil, err
	}
	return wrapImage(C.gdImageCopyGaussianBlurred(ptr, C.int(radius), C.double(sigma)))
}

func (im *Image) SquareToCircle(radius int) (*Image, error) {
	ptr, err := im.cptr()
	if err != nil {
		return nil, err
	}
	return wrapImage(C.gdImageSquareToCircle(ptr, C.int(radius)))
}

func (im *Image) StringFTCircle(cx, cy int, radius, textRadius, fillPortion float64, font string, points float64, top, bottom string, fg Color) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	cfont := C.CString(font)
	ctop := C.CString(top)
	cbottom := C.CString(bottom)
	defer C.free(unsafe.Pointer(cfont))
	defer C.free(unsafe.Pointer(ctop))
	defer C.free(unsafe.Pointer(cbottom))

	msg := C.gdImageStringFTCircle(ptr, C.int(cx), C.int(cy), C.double(radius), C.double(textRadius), C.double(fillPortion), cfont, C.double(points), ctop, cbottom, C.int(fg))
	if msg != nil {
		return FreeTypeError{Message: C.GoString(msg)}
	}
	return nil
}
