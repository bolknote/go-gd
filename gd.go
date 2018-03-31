package gd

// Evgeny Stepanischev. 2011. http://bolknote.ru/ imbolk@gmail.com

/*
#include <gd.h>
#include <gdfx.h>
#include <gdfontt.h>
#include <gdfonts.h>
#include <gdfontmb.h>
#include <gdfontl.h>
#include <gdfontg.h>
#cgo LDFLAGS: -lgd -L/usr/local/lib/
#cgo CFLAGS: -I/usr/local/include/

// Avoid CGO bug https://github.com/golang/go/issues/19832
#pragma GCC diagnostic ignored "-Wincompatible-pointer-types"
*/
import "C"
import . "unsafe"
import "errors"

type Image struct {
	img      C.gdImagePtr
	getpixel func(p *Image, x, y int) Color
}
type Font struct{ fnt C.gdFontPtr }
type Color int
type Style int

type Point struct {
	X, Y int
}

const (
	ARCPIE   Style = 0
	ARCCHORD Style = 1 << iota
	ARCNOFILL
	ARCEDGED
)

const (
	FONTTINY = iota
	FONTSMALL
	FONTMEDIUMBOLD
	FONTLARGE
	FONTGIANT
)

func img(img C.gdImagePtr) *Image {
	if (int)((*img).trueColor) != 0 {
		return &Image{img: img, getpixel: gettruecolorpixel}
	}

	return &Image{img: img, getpixel: getpixel}
}

// http://php.net/manual/en/function.imagecreate.php
func Create(sx, sy int) *Image {
	return img(C.gdImageCreate(C.int(sx), C.int(sy)))
}

// http://php.net/manual/en/function.imagecreatetruecolor.php
func CreateTrueColor(sx, sy int) *Image {
	return img(C.gdImageCreateTrueColor(C.int(sx), C.int(sy)))
}

func CreateFromJpegPtr(imagebuffer []byte) *Image {
	return img(C.gdImageCreateFromJpegPtr(C.int(len(imagebuffer)), Pointer(&imagebuffer[0])))
}

func CreateFromPngPtr(imagebuffer []byte) *Image {
	return img(C.gdImageCreateFromPngPtr(C.int(len(imagebuffer)), Pointer(&imagebuffer[0])))
}

func CreateFromGifPtr(imagebuffer []byte) *Image {
	return img(C.gdImageCreateFromGifPtr(C.int(len(imagebuffer)), Pointer(&imagebuffer[0])))
}

func ImageToJpegBuffer(p *Image, quality int) []byte {
	var imgSize int
	pimgSize := (*C.int)(Pointer(&imgSize))

	buf := C.gdImageJpegPtr(p.img, pimgSize, C.int(quality))
	defer C.gdFree(buf)

	return C.GoBytes(buf, *pimgSize)
}

func ImageToPngBuffer(p *Image) []byte {
	var imgSize int
	pimgSize := (*C.int)(Pointer(&imgSize))

	buf := C.gdImagePngPtr(p.img, pimgSize)
	defer C.gdFree(buf)

	return C.GoBytes(buf, *pimgSize)
}

func ImageToGifBuffer(p *Image) []byte {
	var imgSize int
	pimgSize := (*C.int)(Pointer(&imgSize))

	buf := C.gdImageGifPtr(p.img, pimgSize)
	defer C.gdFree(buf)

	return C.GoBytes(buf, *pimgSize)
}

// http://php.net/manual/en/function.imagecreatefromjpeg.php
func CreateFromJpeg(infile string) *Image {
	name, mode := C.CString(infile), C.CString("rb")

	defer C.free(Pointer(name))
	defer C.free(Pointer(mode))

	file := C.fopen(name, mode)

	if file != nil {
		defer C.fclose(file)

		return img(C.gdImageCreateFromJpeg(file))
	}

	panic(errors.New("Error occurred while opening file."))
}

// http://php.net/manual/en/function.imagecreatefromgif.php
func CreateFromGif(infile string) *Image {
	name, mode := C.CString(infile), C.CString("rb")

	defer C.free(Pointer(name))
	defer C.free(Pointer(mode))

	file := C.fopen(name, mode)

	if file != nil {
		defer C.fclose(file)

		return img(C.gdImageCreateFromGif(file))
	}

	panic(errors.New("Error occurred while opening file."))
}

// http://php.net/manual/en/function.imagecreatefrompng.php
func CreateFromPng(infile string) *Image {
	name, mode := C.CString(infile), C.CString("rb")

	defer C.free(Pointer(name))
	defer C.free(Pointer(mode))

	file := C.fopen(name, mode)

	if file != nil {
		defer C.fclose(file)

		return img(C.gdImageCreateFromPng(file))
	}

	panic(errors.New("Error occurred while opening file."))
}

// http://php.net/manual/en/function.imagecreatefromwbmp.php
func CreateFromWbmp(infile string) *Image {
	name, mode := C.CString(infile), C.CString("rb")

	defer C.free(Pointer(name))
	defer C.free(Pointer(mode))

	file := C.fopen(name, mode)

	if file != nil {
		defer C.fclose(file)

		return img(C.gdImageCreateFromWBMP(file))
	}

	panic(errors.New("Error occurred while opening file."))
}

// http://php.net/manual/en/function.imagecreatefromxbm.php
func CreateImageFromXbm(infile string) *Image {
	name, mode := C.CString(infile), C.CString("rb")

	defer C.free(Pointer(name))
	defer C.free(Pointer(mode))

	file := C.fopen(name, mode)

	if file != nil {
		defer C.fclose(file)

		return img(C.gdImageCreateFromXbm(file))
	}

	panic(errors.New("Error occurred while opening file."))
}

// http://php.net/manual/en/function.imagecreatefromxpm.php
func CreateFromXpm(infile string) (im *Image) {
	defer func() {
		if e := recover(); e != nil {
			panic(errors.New("Error occurred while opening file."))
		}
	}()

	name := C.CString(infile)
	defer C.free(Pointer(name))

	return img(C.gdImageCreateFromXpm(name))
}

// http://php.net/manual/en/function.imagedestroy.php
func (p *Image) Destroy() {
	if p != nil && p.img != nil {
		C.gdImageDestroy(p.img)
	}
}

func (p *Image) SquareToCircle(radius int) *Image {
	return img(C.gdImageSquareToCircle(p.img, C.int(radius)))
}

// http://php.net/manual/en/function.imagejpeg.php
func (p *Image) Jpeg(out string, quality int) {
	name, mode := C.CString(out), C.CString("wb")

	defer C.free(Pointer(name))
	defer C.free(Pointer(mode))

	file := C.fopen(name, mode)

	if file != nil {
		defer C.fclose(file)

		C.gdImageJpeg(p.img, file, C.int(quality))
	} else {
		panic(errors.New("Error occurred while opening file for writing."))
	}
}

// http://php.net/manual/en/function.imagepng.php
func (p *Image) Png(out string) {
	name, mode := C.CString(out), C.CString("wb")

	defer C.free(Pointer(name))
	defer C.free(Pointer(mode))

	file := C.fopen(name, mode)

	if file != nil {
		defer C.fclose(file)

		C.gdImagePng(p.img, file)
	} else {
		panic(errors.New("Error occurred while opening file for writing."))
	}
}

// http://php.net/manual/en/function.imagegif.php
func (p *Image) Gif(out string) {
	name, mode := C.CString(out), C.CString("wb")

	defer C.free(Pointer(name))
	defer C.free(Pointer(mode))

	file := C.fopen(name, mode)

	if file != nil {
		defer C.fclose(file)

		C.gdImageGif(p.img, file)
	} else {
		panic(errors.New("Error occurred while opening file for writing."))
	}
}

// http://php.net/manual/en/function.imagewbmp.php
func (p *Image) Wbmp(out string, foreground Color) {
	name, mode := C.CString(out), C.CString("wb")

	defer C.free(Pointer(name))
	defer C.free(Pointer(mode))

	file := C.fopen(name, mode)

	if file != nil {
		defer C.fclose(file)

		C.gdImageWBMP(p.img, C.int(foreground), file)
	} else {
		panic(errors.New("Error occurred while opening file for writing."))
	}
}

// http://php.net/manual/en/function.imagecolortransparent.php
func (p *Image) ColorTransparent(color Color) {
	C.gdImageColorTransparent(p.img, C.int(color))
}

// http://php.net/manual/en/function.imagepalettecopy.php
func (p *Image) PaletteCopy(dst *Image) {
	C.gdImagePaletteCopy(dst.img, p.img)
}

// http://php.net/manual/en/function.imagecopyresampled.php
func (p *Image) CopyResampled(dst *Image, dstX, dstY, srcX, srcY, dstW, dstH, srcW, srcH int) {
	C.gdImageCopyResampled(dst.img, p.img, C.int(dstX), C.int(dstY), C.int(srcX), C.int(srcY),
		C.int(dstW), C.int(dstH), C.int(srcW), C.int(srcH))
}

// http://php.net/manual/en/function.imagecopyresized.php
func (p *Image) CopyResized(dst *Image, dstX, dstY, srcX, srcY, dstW, dstH, srcW, srcH int) {
	C.gdImageCopyResized(dst.img, p.img, C.int(dstX), C.int(dstY), C.int(srcX), C.int(srcY),
		C.int(dstW), C.int(dstH), C.int(srcW), C.int(srcH))
}

// http://php.net/manual/en/function.imagecopymerge.php
func (p *Image) CopyMerge(dst *Image, dstX, dstY, srcX, srcY, w, h, pct int) {
	C.gdImageCopyMerge(dst.img, p.img, C.int(dstX), C.int(dstY), C.int(srcX), C.int(srcY),
		C.int(w), C.int(h), C.int(pct))
}

// http://php.net/manual/en/function.imagecopymergegray.php
func (p *Image) CopyMergeGray(dst *Image, dstX, dstY, srcX, srcY, w, h, pct int) {
	C.gdImageCopyMergeGray(dst.img, p.img, C.int(dstX), C.int(dstY), C.int(srcX), C.int(srcY),
		C.int(w), C.int(h), C.int(pct))
}

// http://php.net/manual/en/function.imagecopy.php
func (p *Image) Copy(dst *Image, dstX, dstY, srcX, srcY, w, h int) {
	C.gdImageCopy(dst.img, p.img, C.int(dstX), C.int(dstY), C.int(srcX), C.int(srcY), C.int(w), C.int(h))
}

func (p *Image) CopyRotated(dst *Image, dstX, dstY, srcX, srcY, srcWidth, srcHeight, angle int) {
	C.gdImageCopyRotated(dst.img, p.img, C.double(dstX), C.double(dstY), C.int(srcX), C.int(srcY),
		C.int(srcWidth), C.int(srcHeight), C.int(angle))
}

// http://php.net/manual/en/function.imagecolorallocate.php
func (p *Image) ColorAllocate(r, g, b int) Color {
	return (Color)(C.gdImageColorAllocate(p.img, C.int(r), C.int(g), C.int(b)))
}

// http://php.net/manual/en/function.imagecolorallocatealpha.php
func (p *Image) ColorAllocateAlpha(r, g, b, a int) Color {
	return (Color)(C.gdImageColorAllocateAlpha(p.img, C.int(r), C.int(g), C.int(b), C.int(a)))
}

// http://php.net/manual/en/function.imagecolorclosest.php
func (p *Image) ColorClosest(r, g, b int) Color {
	return (Color)(C.gdImageColorClosest(p.img, C.int(r), C.int(g), C.int(b)))
}

// http://php.net/manual/en/function.imagecolorclosestalpha.php
func (p *Image) ColorClosestAlpha(r, g, b, a int) Color {
	return (Color)(C.gdImageColorClosestAlpha(p.img, C.int(r), C.int(g), C.int(b), C.int(a)))
}

// http://php.net/manual/en/function.imagecolorclosesthwb.php
func (p *Image) ColorClosestHWB(r, g, b int) Color {
	return (Color)(C.gdImageColorClosestHWB(p.img, C.int(r), C.int(g), C.int(b)))
}

// http://php.net/manual/en/function.imagecolorexact.php
func (p *Image) ColorExact(r, g, b int) Color {
	return (Color)(C.gdImageColorExact(p.img, C.int(r), C.int(g), C.int(b)))
}

// http://php.net/manual/en/function.imagecolorexactalpha.php
func (p *Image) ColorExactAlpha(r, g, b, a int) Color {
	return (Color)(C.gdImageColorExactAlpha(p.img, C.int(r), C.int(g), C.int(b), C.int(a)))
}

// http://php.net/manual/en/function.imagecolorresolve.php
func (p *Image) ColorResolve(r, g, b int) Color {
	return (Color)(C.gdImageColorResolve(p.img, C.int(r), C.int(g), C.int(b)))
}

// http://php.net/manual/en/function.imagecolorresolvealpha.php
func (p *Image) ColorResolveAlpha(r, g, b, a int) Color {
	return (Color)(C.gdImageColorResolveAlpha(p.img, C.int(r), C.int(g), C.int(b), C.int(a)))
}

// http://php.net/manual/en/function.imagecolordeallocate.php
func (p *Image) ColorDeallocate(color Color) {
	C.gdImageColorDeallocate(p.img, C.int(color))
}

// http://php.net/manual/en/function.imagefill.php
func (p *Image) Fill(x, y int, c Color) {
	C.gdImageFill(p.img, C.int(x), C.int(y), C.int(c))
}

// http://php.net/manual/en/function.imagefilledarc.php
func (p *Image) FilledArc(cx, cy, w, h, s, e, color Color, style Style) {
	C.gdImageFilledArc(p.img, C.int(cx), C.int(cy), C.int(w), C.int(h), C.int(s),
		C.int(e), C.int(color), C.int(style))
}

// http://php.net/manual/en/function.imagearc.php
func (p *Image) Arc(cx, cy, w, h, s, e int, color Color) {
	C.gdImageArc(p.img, C.int(cx), C.int(cy), C.int(w), C.int(h),
		C.int(s), C.int(e), C.int(color))
}

// http://php.net/manual/en/function.imagefilledellipse.php
func (p *Image) FilledEllipse(cx, cy, w, h int, color Color) {
	C.gdImageFilledEllipse(p.img, C.int(cx), C.int(cy), C.int(w), C.int(h), C.int(color))
}

// http://php.net/manual/en/function.imageellipse.php
func (p *Image) Ellipse(cx, cy, w, h int, color Color) {
	C.gdImageEllipse(p.img, C.int(cx), C.int(cy), C.int(w), C.int(h), C.int(color))
}

// http://php.net/manual/en/function.imagefilltoborder.php
func (p *Image) FillToBorder(x, y, border, color Color) {
	C.gdImageFillToBorder(p.img, C.int(x), C.int(y), C.int(border), C.int(color))
}

func (p *Image) Sharpen(pct int) {
	C.gdImageSharpen(p.img, C.int(pct))
}

// http://php.net/manual/en/function.imagesx.php
func (p *Image) Sx() int {
	return (int)((*p.img).sx)
}

// http://php.net/manual/en/function.imagesy.php
func (p *Image) Sy() int {
	return (int)((*p.img).sy)
}

func (p *Image) GetInterlaced() bool {
	return (int)((*p.img).interlace) != 0
}

// http://php.net/manual/en/function.imagecolorstotal.php
func (p *Image) ColorsTotal() int {
	return (int)((*p.img).colorsTotal)
}

func (p *Image) TrueColor() bool {
	return (int)((*p.img).trueColor) != 0
}

// http://php.net/manual/en/function.imagesetpixel.php
func (p *Image) SetPixel(x, y int, color Color) {
	C.gdImageSetPixel(p.img, C.int(x), C.int(y), C.int(color))
}

func getpixel(p *Image, x, y int) Color {
	return (Color)(C.gdImageGetPixel(p.img, C.int(x), C.int(y)))
}

func gettruecolorpixel(p *Image, x, y int) Color {
	return (Color)(C.gdImageGetPixel(p.img, C.int(x), C.int(y)))
}

// http://php.net/manual/en/function.imagecolorat.php
func (p *Image) ColorAt(x, y int) Color {
	return (*p).getpixel(p, x, y)
}

func (p *Image) AABlend() {
	C.gdImageAABlend(p.img)
}

// http://php.net/manual/en/function.imageline.php
func (p *Image) Line(x1, y1, x2, y2 int, color Color) {
	C.gdImageLine(p.img, C.int(x1), C.int(y1), C.int(x2), C.int(y2), C.int(color))
}

// http://php.net/manual/en/function.imagedashedline.php
func (p *Image) DashedLine(x1, y1, x2, y2 int, color Color) {
	C.gdImageDashedLine(p.img, C.int(x1), C.int(y1), C.int(x2), C.int(y2), C.int(color))
}

func (p *Image) Rectangle(x1, y1, x2, y2 int, color Color) {
	C.gdImageRectangle(p.img, C.int(x1), C.int(y1), C.int(x2), C.int(y2), C.int(color))
}

// http://php.net/manual/en/function.imagefilledrectangle.php
func (p *Image) FilledRectangle(x1, y1, x2, y2 int, color Color) {
	C.gdImageFilledRectangle(p.img, C.int(x1), C.int(y1), C.int(x2), C.int(y2), C.int(color))
}

// http://php.net/manual/en/function.imagesavealpha.php
func (p *Image) SaveAlpha(saveflag bool) {
	C.gdImageSaveAlpha(p.img, map[bool]C.int{true: 1, false: 0}[saveflag])
}

// http://php.net/manual/en/function.imagealphablending.php
func (p *Image) AlphaBlending(blendmode bool) {
	C.gdImageAlphaBlending(p.img, map[bool]C.int{true: 1, false: 0}[blendmode])
}

// http://php.net/manual/en/function.imageinterlace.php
func (p *Image) Interlace(interlacemode bool) {
	C.gdImageInterlace(p.img, map[bool]C.int{true: 1, false: 0}[interlacemode])
}

// http://php.net/manual/en/function.imagesetthickness.php
func (p *Image) SetThickness(thickness int) {
	C.gdImageSetThickness(p.img, C.int(thickness))
}

// http://php.net/manual/en/function.imagetruecolortopalette.php
func (p *Image) TrueColorToPalette(ditherFlag bool, colorsWanted int) {
	C.gdImageTrueColorToPalette(p.img, map[bool]C.int{true: 1, false: 0}[ditherFlag], C.int(colorsWanted))
}

// http://php.net/manual/en/function.imagesetstyle.php
func (p *Image) SetStyle(style ...Color) {
	C.gdImageSetStyle(p.img, (*C.int)(Pointer(&style)), C.int(len(style)))
}

func (p *Image) SetAntiAliased(c Color) {
	C.gdImageSetAntiAliased(p.img, C.int(c))
}

func (p *Image) SetAntiAliasedDontBlend(c Color, dont_blend bool) {
	C.gdImageSetAntiAliasedDontBlend(p.img, C.int(c), map[bool]C.int{true: 1, false: 0}[dont_blend])
}

// http://php.net/manual/en/function.imagesettile.php
func (p *Image) SetTile(tile Image) {
	C.gdImageSetTile(p.img, tile.img)
}

// http://php.net/manual/en/function.imagesetbrush.php
func (p *Image) SetBrush(brush Image) {
	C.gdImageSetBrush(p.img, brush.img)
}

func GetFont(size byte) *Font {
	switch size {
	case FONTTINY:
		return &Font{fnt: C.gdFontGetTiny()}
	case FONTSMALL:
		return &Font{fnt: C.gdFontGetSmall()}
	case FONTMEDIUMBOLD:
		return &Font{fnt: C.gdFontGetMediumBold()}
	case FONTLARGE:
		return &Font{fnt: C.gdFontGetLarge()}
	case FONTGIANT:
		return &Font{fnt: C.gdFontGetGiant()}
	}

	panic(errors.New("Invalid font size"))
}

// http://php.net/manual/en/function.imagechar.php
func (p *Image) Char(font *Font, x, y int, c string, color Color) {
	C.gdImageChar(p.img, (*font).fnt, C.int(x), C.int(y), C.int(([]byte(c))[0]), C.int(color))
}

// http://php.net/manual/en/function.imagecharup.php
func (p *Image) CharUp(font *Font, x, y int, c string, color Color) {
	C.gdImageCharUp(p.img, (*font).fnt, C.int(x), C.int(y), C.int(([]byte(c))[0]), C.int(color))
}

// http://php.net/manual/en/function.imagestring.php
func (p *Image) String(font *Font, x, y int, s string, color Color) {
	str := Pointer(C.CString(s))
	C.gdImageString(p.img, (*font).fnt, C.int(x), C.int(y), (*C.uchar)(str), C.int(color))
	C.free(str)
}

// http://php.net/manual/en/function.imagestringup.php
func (p *Image) StringUp(font *Font, x, y int, s string, color Color) {
	str := Pointer(C.CString(s))
	C.gdImageStringUp(p.img, (*font).fnt, C.int(x), C.int(y), (*C.uchar)(str), C.int(color))
	C.free(str)
}

func (p *Image) StringFT(fg Color, fontname string, ptsize, angle float64, x, y int, str string) (brect [8]int32) {
	C.gdFontCacheSetup()
	defer C.gdFontCacheShutdown()

	cfontname, cstr := C.CString(fontname), C.CString(str)

	C.gdImageStringFT(p.img, (*C.int)(Pointer(&brect)), C.int(fg), cfontname, C.double(ptsize),
		C.double(angle), C.int(x), C.int(y), cstr)

	C.free(Pointer(cfontname))
	C.free(Pointer(cstr))

	return
}

func pointsTogdPoints(points []Point) C.gdPointPtr {
	const maxlen = 1 << 30 // ought to be enough for anybody
	alen := len(points)

	if alen > maxlen {
		panic("Too many arguments")
	}

	mem := C.calloc(C.size_t(alen), C.sizeof_gdPoint)
	gdpoints := (*[maxlen]C.gdPoint)(mem)

	for i, v := range points {
		gdpoints[i] = C.gdPoint{C.int(v.X), C.int(v.Y)}
	}

	return C.gdPointPtr(mem)
}

// http://php.net/manual/en/function.imagepolygon.php
func (p *Image) Polygon(points []Point, c Color) {
	gdpoints := pointsTogdPoints(points)
	C.gdImagePolygon(p.img, gdpoints, C.int(len(points)), C.int(c))
	C.free(Pointer(gdpoints))
}

// http://php.net/manual/en/function.imageopenpolygon.php
func (p *Image) OpenPolygon(points []Point, c Color) {
	gdpoints := pointsTogdPoints(points)
	C.gdImageOpenPolygon(p.img, gdpoints, C.int(len(points)), C.int(c))
	C.free(Pointer(gdpoints))
}

// http://php.net/manual/en/function.imagefilledpolygon.php
func (p *Image) FilledPolygon(points []Point, c Color) {
	gdpoints := pointsTogdPoints(points)
	C.gdImageFilledPolygon(p.img, gdpoints, C.int(len(points)), C.int(c))
	C.free(Pointer(gdpoints))
}

// http://php.net/manual/en/function.imagecolorsforindex.php
func (p *Image) ColorsForIndex(index Color) map[string]int {
	if p.TrueColor() {
		return map[string]int{
			"alpha": (int(index) & 0x7F000000) >> 24,
			"red":   (int(index) & 0xFF0000) >> 16,
			"green": (int(index) & 0x00FF00) >> 8,
			"blue":  (int(index) & 0x0000FF),
		}
	}

	return map[string]int{
		"red":   (int)((*p.img).red[index]),
		"green": (int)((*p.img).green[index]),
		"blue":  (int)((*p.img).blue[index]),
		"alpha": (int)((*p.img).alpha[index]),
	}
}

