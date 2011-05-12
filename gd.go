package gd
// #include <gd.h>
// #include <gdfx.h>
import "C"
import "os"

type Image struct {img C.gdImagePtr}
type Color int
type Style int

const (
    ARCPIE Style = 0
    ARCCHORD Style = 1 << iota
    ARCNOFILL
    ARCEDGED
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
    }

    panic(os.NewError("Error occurred while opening file for writing."))
}

func (p *Image) CopyResampled(dst Image, dstX, dstY, srcX, srcY, dstW, dstH, srcW, srcH int) {
    C.gdImageCopyResampled(dst.img, p.img, C.int(dstX), C.int(dstY), C.int(srcX), C.int(srcY),
        C.int(dstW), C.int(dstH), C.int(srcW), C.int(srcH))
}

func (p *Image) Copy(dst Image, dstX, dstY, srcX, srcY, w, h int) {
    C.gdImageCopy(dst.img, p.img, C.int(dstX), C.int(dstY), C.int(srcX), C.int(srcY), C.int(w), C.int(h))
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




