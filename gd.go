package gd
// #include <gd.h>
// #include <gdfx.h>
// #include <gdfontt.h>
// #include <gdfonts.h>
// #include <gdfontmb.h>
// #include <gdfontl.h>
// #include <gdfontg.h>
import "C"
import "os"
import . "unsafe"
//import "utf16"

type Image struct {img C.gdImagePtr}
type Font  struct {fnt C.gdFontPtr}
type Color int
type Style int

const (
    ARCPIE Style = 0
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

func Create(sx, sy int) *Image {
    return &Image{img: C.gdImageCreate(C.int(sx), C.int(sy))}
}

func CreateTrueColor(sx, sy int) *Image {
    return &Image{img: C.gdImageCreateTrueColor(C.int(sx), C.int(sy))}
}

func CreateFromJpeg(infile string) *Image {
    file := C.fopen(C.CString(infile), C.CString("rb"))

    if file != nil {
        defer C.fclose(file)

        return &Image{img: C.gdImageCreateFromJpeg(file)}
    }

    panic(os.NewError("Error occurred while opening file."))
}


func CreateFromGif(infile string) *Image {
    file := C.fopen(C.CString(infile), C.CString("rb"))

    if file != nil {
        defer C.fclose(file)

        return &Image{img: C.gdImageCreateFromGif(file)}
    }

    panic(os.NewError("Error occurred while opening file."))
}

func CreateFromPng(infile string) *Image {
    file := C.fopen(C.CString(infile), C.CString("rb"))

    if file != nil {
        defer C.fclose(file)

        return &Image{img: C.gdImageCreateFromPng(file)}
    }

    panic(os.NewError("Error occurred while opening file."))
}

func CreateImageFromWbmp(infile string) *Image {
    file := C.fopen(C.CString(infile), C.CString("rb"))

    if file != nil {
        defer C.fclose(file)

        return &Image{img: C.gdImageCreateFromWBMP(file)}
    }

    panic(os.NewError("Error occurred while opening file."))
}

func (p *Image) Destroy() {
    if p != nil && p.img != nil {
        C.gdImageDestroy(p.img)
    }
}

func (p *Image) SquareToCircle(radius int) *Image {
    return &Image{img: C.gdImageSquareToCircle(p.img, C.int(radius))}
}

func (p *Image) Jpeg(out string, quality int) {
    file := C.fopen(C.CString(out), C.CString("wb"))

    if file != nil {
        defer C.fclose(file)

        C.gdImageJpeg(p.img, file, C.int(quality))
    } else {
        panic(os.NewError("Error occurred while opening file for writing."))
    }
}

func (p *Image) Png(out string) {
    file := C.fopen(C.CString(out), C.CString("wb"))

    if file != nil {
        defer C.fclose(file)

        C.gdImagePng(p.img, file)
    } else {
        panic(os.NewError("Error occurred while opening file for writing."))
    }
}

func (p *Image) Gif(out string) {
    file := C.fopen(C.CString(out), C.CString("wb"))

    if file != nil {
        defer C.fclose(file)

        C.gdImageGif(p.img, file)
    } else {
        panic(os.NewError("Error occurred while opening file for writing."))
    }
}

func (p *Image) Wbmp(out string, foreground Color) {
    file := C.fopen(C.CString(out), C.CString("wb"))

    if file != nil {
        defer C.fclose(file)

        C.gdImageWBMP(p.img, C.int(foreground), file)
    } else {
        panic(os.NewError("Error occurred while opening file for writing."))
    }
}

func (p *Image) ColorTransparent(color Color) {
    C.gdImageColorTransparent(p.img, C.int(color))
}

func (p *Image) PaletteCopy(dst Image) {
    C.gdImagePaletteCopy(dst.img, p.img)
}

func (p *Image) CopyResampled(dst Image, dstX, dstY, srcX, srcY, dstW, dstH, srcW, srcH int) {
    C.gdImageCopyResampled(dst.img, p.img, C.int(dstX), C.int(dstY), C.int(srcX), C.int(srcY),
        C.int(dstW), C.int(dstH), C.int(srcW), C.int(srcH))
}

func (p *Image) CopyResized(dst Image, dstX, dstY, srcX, srcY, dstW, dstH, srcW, srcH int) {
    C.gdImageCopyResized(dst.img, p.img, C.int(dstX), C.int(dstY), C.int(srcX), C.int(srcY),
        C.int(dstW), C.int(dstH), C.int(srcW), C.int(srcH))
}

func (p *Image) CopyMerge(dst Image, dstX, dstY, srcX, srcY, w, h, pct int) {
    C.gdImageCopyMerge(dst.img, p.img, C.int(dstX), C.int(dstY), C.int(srcX), C.int(srcY),
        C.int(w), C.int(h), C.int(pct))
}

func (p *Image) CopyMergeGray(dst Image, dstX, dstY, srcX, srcY, w, h, pct int) {
    C.gdImageCopyMergeGray(dst.img, p.img, C.int(dstX), C.int(dstY), C.int(srcX), C.int(srcY),
        C.int(w), C.int(h), C.int(pct))
}

func (p *Image) Copy(dst Image, dstX, dstY, srcX, srcY, w, h int) {
    C.gdImageCopy(dst.img, p.img, C.int(dstX), C.int(dstY), C.int(srcX), C.int(srcY), C.int(w), C.int(h))
}

func (p *Image) CopyRotated(dst Image, dstX, dstY, srcX, srcY, srcWidth, srcHeight, angle int) {
    C.gdImageCopyRotated(dst.img, p.img, C.double(dstX), C.double(dstY), C.int(srcX), C.int(srcY),
        C.int(srcWidth), C.int(srcHeight), C.int(angle))
}

func (p *Image) ColorAllocate(r, g, b int) Color {
    return (Color)(C.gdImageColorAllocate(p.img, C.int(r), C.int(g), C.int(b)))
}

func (p *Image) ColorAllocateAlpha(r, g, b, a int) Color {
    return (Color)(C.gdImageColorAllocateAlpha(p.img, C.int(r), C.int(g), C.int(b), C.int(a)))
}

func (p *Image) ColorClosest(r, g, b int) Color {
    return (Color)(C.gdImageColorClosest(p.img, C.int(r), C.int(g), C.int(b)))
}

func (p *Image) ColorClosestAlpha(r, g, b, a int) Color {
    return (Color)(C.gdImageColorClosestAlpha(p.img, C.int(r), C.int(g), C.int(b), C.int(a)))
}

func (p *Image) ColorClosestHWB(r, g, b int) Color {
    return (Color)(C.gdImageColorClosestHWB(p.img, C.int(r), C.int(g), C.int(b)))
}

func (p *Image) ColorExact(r, g, b int) Color {
    return (Color)(C.gdImageColorExact(p.img, C.int(r), C.int(g), C.int(b)))
}

func (p *Image) ColorExactAlpha(r, g, b, a int) Color {
    return (Color)(C.gdImageColorExactAlpha(p.img, C.int(r), C.int(g), C.int(b), C.int(a)))
}

func (p *Image) ColorResolve(r, g, b int) Color {
    return (Color)(C.gdImageColorResolve(p.img, C.int(r), C.int(g), C.int(b)))
}

func (p *Image) ColorResolveAlpha(r, g, b, a int) Color {
    return (Color)(C.gdImageColorResolveAlpha(p.img, C.int(r), C.int(g), C.int(b), C.int(a)))
}

func (p *Image) ColorDeallocate(color Color) {
    C.gdImageColorDeallocate(p.img, C.int(color))
}

func (p *Image) Fill(x, y int, c Color) {
    C.gdImageFill(p.img, C.int(x), C.int(y), C.int(c))
}

func (p *Image) FilledArc(cx, cy, w, h, s, e, color Color, style Style) {
    C.gdImageFilledArc(p.img, C.int(cx), C.int(cy), C.int(w), C.int(h), C.int(s),
        C.int(e), C.int(color), C.int(style))
}

func (p *Image) Arc(cx, cy, w, h, s, e int, color Color) {
    C.gdImageArc(p.img, C.int(cx), C.int(cy), C.int(w), C.int(h),
        C.int(s), C.int(e), C.int(color))
}

func (p *Image) FilledEllipse(cx, cy, w, h, color Color) {
    C.gdImageFilledEllipse(p.img, C.int(cx), C.int(cy), C.int(w), C.int(h), C.int(color))
}

func (p *Image) FillToBorder(x, y, border, color Color) {
    C.gdImageFillToBorder(p.img, C.int(x), C.int(y), C.int(border), C.int(color))
}

func (p *Image) Sharpen(pct int) {
    C.gdImageSharpen(p.img, C.int(pct))
}

func (p *Image) Sx() int {
    return (int)((*p.img).sx)
}

func (p *Image) Sy() int {
    return (int)((*p.img).sy)
}

func (p *Image) GetInterlaced() bool {
    return (int)((*p.img).interlace) != 0
}

func (p *Image) ColorsTotal() int {
    return (int)((*p.img).colorsTotal)
}

func (p *Image) TrueColor() bool {
    return (int)((*p.img).trueColor) != 0
}

func (p *Image) SetPixel(x, y int, color Color) {
    C.gdImageSetPixel(p.img, C.int(x), C.int(y), C.int(color))
}

func (p *Image) GetPixel(x, y int) Color {
    return (Color)(C.gdImageGetPixel(p.img, C.int(x), C.int(y)))
}

func (p *Image) GetTrueColorPixel(x, y int) Color {
    return (Color)(C.gdImageGetTrueColorPixel(p.img, C.int(x), C.int(y)))
}

func (p *Image) AABlend() {
    C.gdImageAABlend(p.img)
}

func (p *Image) Line(x1, y1, x2, y2 int, color Color) {
    C.gdImageLine(p.img, C.int(x1), C.int(y1), C.int(x2), C.int(y2), C.int(color))
}

func (p *Image) DashedLine(x1, y1, x2, y2 int, color Color) {
    C.gdImageDashedLine(p.img, C.int(x1), C.int(y1), C.int(x2), C.int(y2), C.int(color))
}

func (p *Image) Rectangle(x1, y1, x2, y2 int, color Color) {
    C.gdImageRectangle(p.img, C.int(x1), C.int(y1), C.int(x2), C.int(y2), C.int(color))
}

func (p *Image) FilledRectangle(x1, y1, x2, y2 int, color Color) {
    C.gdImageFilledRectangle(p.img, C.int(x1), C.int(y1), C.int(x2), C.int(y2), C.int(color))
}

func (p *Image) SaveAlpha(saveflag bool) {
    C.gdImageSaveAlpha(p.img, map[bool]C.int{true: 1, false: 0}[saveflag])
}

func (p *Image) AlphaBlending(blendmode bool) {
    C.gdImageAlphaBlending(p.img, map[bool]C.int{true: 1, false: 0}[blendmode])
}

func (p *Image) Interlace(interlacemode bool) {
    C.gdImageInterlace(p.img, map[bool]C.int{true: 1, false: 0}[interlacemode])
}

func (p *Image) SetThickness(thickness int) {
    C.gdImageSetThickness(p.img, C.int(thickness))
}

func (p *Image) TrueColorToPalette(ditherFlag bool, colorsWanted int) {
    C.gdImageTrueColorToPalette(p.img, map[bool]C.int{true: 1, false: 0}[ditherFlag], C.int(colorsWanted))
}

func (p *Image) SetStyle(style... Color) {
    C.gdImageSetStyle(p.img, (*C.int)(Pointer(&style)), C.int(len(style)))
}

func (p *Image) SetAntiAliased(c Color) {
     C.gdImageSetAntiAliased(p.img, C.int(c))
}

func (p *Image) SetAntiAliasedDontBlend(c Color, dont_blend bool) {
    C.gdImageSetAntiAliasedDontBlend(p.img, C.int(c), map[bool]C.int{true: 1, false: 0}[dont_blend])
}

func (p *Image) SetTile(tile Image) {
    C.gdImageSetTile(p.img, tile.img)
}

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

    panic(os.NewError("Invalid font size"))
}

func (p *Image) Char(font *Font, x, y int, c string, color Color) {
    C.gdImageChar(p.img, (*font).fnt, C.int(x), C.int(y), C.int(([]byte(c))[0]), C.int(color))
}

func (p *Image) CharUp(font *Font, x, y int, c string, color Color) {
    C.gdImageCharUp(p.img, (*font).fnt, C.int(x), C.int(y), C.int(([]byte(c))[0]), C.int(color))
}

func (p *Image) String(font *Font, x, y int, s string, color Color) {
    C.gdImageString(p.img, (*font).fnt, C.int(x), C.int(y), (*C.uchar)(Pointer(C.CString(s))), C.int(color))
}

func (p *Image) StringUp(font *Font, x, y int, s string, color Color) {
    C.gdImageStringUp(p.img, (*font).fnt, C.int(x), C.int(y), (*C.uchar)(Pointer(C.CString(s))), C.int(color))
}

func (p *Image) StringFT(fg Color, fontname string, ptsize, angle, x, y int, str string) (brect [8]int) {
    C.gdFontCacheSetup()
    defer C.gdFontCacheShutdown()

    C.gdImageStringFT(p.img, (*C.int)(Pointer(&brect)), C.int(fg), C.CString(fontname), C.double(ptsize),
        C.double(angle), C.int(x), C.int(y), C.CString(str))

    return
}




