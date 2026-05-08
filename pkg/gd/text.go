package gd

/*
#include <stdlib.h>
#include <string.h>
#include <gd.h>
#include <gdfontt.h>
#include <gdfonts.h>
#include <gdfontmb.h>
#include <gdfontl.h>
#include <gdfontg.h>

static char *go_gd_image_string_ft_ex(gdImagePtr im, int *brect, int fg, const char *fontlist,
	double ptsize, double angle, int x, int y, const char *string, int use_extra,
	int flags, double linespacing, int charmap, int hdpi, int vdpi, char **xshow, char **fontpath) {
	gdFTStringExtra extra;
	memset(&extra, 0, sizeof(extra));
	extra.flags = flags;
	extra.linespacing = linespacing;
	extra.charmap = charmap;
	extra.hdpi = hdpi;
	extra.vdpi = vdpi;
	char *msg = gdImageStringFTEx(im, brect, fg, fontlist, ptsize, angle, x, y, string, use_extra ? &extra : NULL);
	if (xshow != NULL) {
		*xshow = extra.xshow;
	}
	if (fontpath != NULL) {
		*fontpath = extra.fontpath;
	}
	return msg;
}
*/
import "C"
import (
	"unsafe"
)

// Font references one of libgd's built-in bitmap fonts. Built-in fonts are
// owned by libgd and must not be closed by callers.
type Font struct {
	ptr C.gdFontPtr
}

type FontSize int

const (
	FontTiny FontSize = iota
	FontSmall
	FontMediumBold
	FontLarge
	FontGiant
)

type FTOptions struct {
	LineSpacing           float64
	Charmap               int
	HDPI                  int
	VDPI                  int
	DisableKerning        bool
	UseFontPathName       bool
	UseFontConfig         bool
	ReturnFontPathName    bool
	ReturnHorizontalSpans bool
}

type FTResult struct {
	Bounds   [8]int
	XShow    string
	FontPath string
}

// FreeTypeError reports a libgd FreeType text-rendering error.
type FreeTypeError struct {
	Message string
}

func (e FreeTypeError) Error() string {
	return "gd: freetype: " + e.Message
}

func BuiltinFont(size FontSize) (*Font, error) {
	switch size {
	case FontTiny:
		return &Font{ptr: C.gdFontGetTiny()}, nil
	case FontSmall:
		return &Font{ptr: C.gdFontGetSmall()}, nil
	case FontMediumBold:
		return &Font{ptr: C.gdFontGetMediumBold()}, nil
	case FontLarge:
		return &Font{ptr: C.gdFontGetLarge()}, nil
	case FontGiant:
		return &Font{ptr: C.gdFontGetGiant()}, nil
	default:
		return nil, ErrInvalidArgument
	}
}

func (f *Font) cptr() (C.gdFontPtr, error) {
	if f == nil || f.ptr == nil {
		return nil, ErrInvalidArgument
	}
	return f.ptr, nil
}

func (im *Image) Char(font *Font, x, y int, r rune, color Color) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	fptr, err := font.cptr()
	if err != nil {
		return err
	}
	C.gdImageChar(ptr, fptr, C.int(x), C.int(y), C.int(r), C.int(color))
	return nil
}

func (im *Image) CharUp(font *Font, x, y int, r rune, color Color) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	fptr, err := font.cptr()
	if err != nil {
		return err
	}
	C.gdImageCharUp(ptr, fptr, C.int(x), C.int(y), C.int(r), C.int(color))
	return nil
}

// String draws Latin-1 text with a built-in bitmap font.
func (im *Image) String(font *Font, x, y int, s string, color Color) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	fptr, err := font.cptr()
	if err != nil {
		return err
	}
	cs, err := latin1CString(s)
	if err != nil {
		return err
	}
	defer C.free(unsafe.Pointer(cs))
	C.gdImageString(ptr, fptr, C.int(x), C.int(y), (*C.uchar)(unsafe.Pointer(cs)), C.int(color))
	return nil
}

// StringUp draws vertical Latin-1 text with a built-in bitmap font.
func (im *Image) StringUp(font *Font, x, y int, s string, color Color) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	fptr, err := font.cptr()
	if err != nil {
		return err
	}
	cs, err := latin1CString(s)
	if err != nil {
		return err
	}
	defer C.free(unsafe.Pointer(cs))
	C.gdImageStringUp(ptr, fptr, C.int(x), C.int(y), (*C.uchar)(unsafe.Pointer(cs)), C.int(color))
	return nil
}

func latin1CString(s string) (*C.char, error) {
	buf := make([]byte, 0, len(s)+1)
	for _, r := range s {
		if r > 0xff {
			return nil, ErrInvalidArgument
		}
		buf = append(buf, byte(r))
	}
	buf = append(buf, 0)
	return (*C.char)(C.CBytes(buf)), nil
}

// String16 draws horizontal Unicode (UCS-2) text with a built-in bitmap font.
// An empty s is a no-op rather than an error, matching the behaviour of
// String for empty strings.
func (im *Image) String16(font *Font, x, y int, s []uint16, color Color) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	fptr, err := font.cptr()
	if err != nil {
		return err
	}
	cs := make([]C.ushort, len(s)+1)
	for i, r := range s {
		cs[i] = C.ushort(r)
	}
	C.gdImageString16(ptr, fptr, C.int(x), C.int(y), (*C.ushort)(unsafe.Pointer(&cs[0])), C.int(color))
	return nil
}

// StringUp16 draws vertical Unicode (UCS-2) text with a built-in bitmap font.
func (im *Image) StringUp16(font *Font, x, y int, s []uint16, color Color) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	fptr, err := font.cptr()
	if err != nil {
		return err
	}
	cs := make([]C.ushort, len(s)+1)
	for i, r := range s {
		cs[i] = C.ushort(r)
	}
	C.gdImageStringUp16(ptr, fptr, C.int(x), C.int(y), (*C.ushort)(unsafe.Pointer(&cs[0])), C.int(color))
	return nil
}

func FontCacheSetup() {
	C.gdFontCacheSetup()
}

func FontCacheShutdown() {
	C.gdFontCacheShutdown()
}

func FreeFontCache() {
	C.gdFreeFontCache()
}

func UseFontConfig(enabled bool) bool {
	return C.gdFTUseFontConfig(boolInt(enabled)) != 0
}

// StringFT renders text via libgd's FreeType backend with default options.
// libgd exposes gdImageStringTTF as a macro alias for gdImageStringFT, so no
// separate StringTTF wrapper is provided.
func (im *Image) StringFT(color Color, fontPath string, pointSize, angle float64, x, y int, text string) (FTResult, error) {
	return im.StringFTEx(color, fontPath, pointSize, angle, x, y, text, nil)
}

func (im *Image) StringFTEx(color Color, fontPath string, pointSize, angle float64, x, y int, text string, opts *FTOptions) (FTResult, error) {
	ptr, err := im.cptr()
	if err != nil {
		return FTResult{}, err
	}
	cfont := C.CString(fontPath)
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(cfont))
	defer C.free(unsafe.Pointer(ctext))

	var brect [8]C.int
	useExtra := C.int(0)
	flags := C.int(0)
	lineSpacing := C.double(0)
	charmap := C.int(0)
	hdpi := C.int(0)
	vdpi := C.int(0)
	if opts != nil {
		useExtra = 1
		if opts.LineSpacing != 0 {
			flags |= C.gdFTEX_LINESPACE
			lineSpacing = C.double(opts.LineSpacing)
		}
		if opts.Charmap != 0 {
			flags |= C.gdFTEX_CHARMAP
			charmap = C.int(opts.Charmap)
		}
		if opts.HDPI != 0 || opts.VDPI != 0 {
			flags |= C.gdFTEX_RESOLUTION
			hdpi = C.int(opts.HDPI)
			vdpi = C.int(opts.VDPI)
		}
		if opts.DisableKerning {
			flags |= C.gdFTEX_DISABLE_KERNING
		}
		if opts.ReturnHorizontalSpans {
			flags |= C.gdFTEX_XSHOW
		}
		if opts.UseFontPathName {
			flags |= C.gdFTEX_FONTPATHNAME
		}
		if opts.UseFontConfig {
			flags |= C.gdFTEX_FONTCONFIG
		}
		if opts.ReturnFontPathName {
			flags |= C.gdFTEX_RETURNFONTPATHNAME
		}
	}

	var xshow, cFontPath *C.char
	msg := C.go_gd_image_string_ft_ex(ptr, (*C.int)(unsafe.Pointer(&brect[0])), C.int(color), cfont, C.double(pointSize), C.double(angle), C.int(x), C.int(y), ctext, useExtra, flags, lineSpacing, charmap, hdpi, vdpi, &xshow, &cFontPath)
	if msg != nil {
		return FTResult{}, FreeTypeError{Message: C.GoString(msg)}
	}

	var result FTResult
	for i, v := range brect {
		result.Bounds[i] = int(v)
	}
	if xshow != nil {
		result.XShow = C.GoString(xshow)
		C.gdFree(unsafe.Pointer(xshow))
	}
	if cFontPath != nil {
		result.FontPath = C.GoString(cFontPath)
		C.gdFree(unsafe.Pointer(cFontPath))
	}
	return result, nil
}
