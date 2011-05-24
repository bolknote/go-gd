package gd
// Evgeny Stepanischev. 2011. http://bolknote.ru/ imbolk@gmail.com

// #include <gd.h>
// #include <gdfx.h>
// #include <gdfontt.h>
// #include <gdfonts.h>
// #include <gdfontmb.h>
// #include <gdfontl.h>
// #include <gdfontg.h>
import "C"
import "os"
import "path/filepath"
import "strings"
import "io/ioutil"
import . "unsafe"
//import "fmt"

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

func CreateImageFromXbm(infile string) *Image {
    file := C.fopen(C.CString(infile), C.CString("rb"))

    if file != nil {
        defer C.fclose(file)

        return &Image{img: C.gdImageCreateFromXbm(file)}
    }

    panic(os.NewError("Error occurred while opening file."))
}

func CreateImageFromXpm(infile string) (im *Image) {
    defer func() {
        if e := recover(); e != nil {
            panic(os.NewError("Error occurred while opening file."))
        }
    }()

    return &Image{img: C.gdImageCreateFromXpm(C.CString(infile))}
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

func (p *Image) PaletteCopy(dst *Image) {
    C.gdImagePaletteCopy(dst.img, p.img)
}

func (p *Image) CopyResampled(dst *Image, dstX, dstY, srcX, srcY, dstW, dstH, srcW, srcH int) {
    C.gdImageCopyResampled(dst.img, p.img, C.int(dstX), C.int(dstY), C.int(srcX), C.int(srcY),
        C.int(dstW), C.int(dstH), C.int(srcW), C.int(srcH))
}

func (p *Image) CopyResized(dst *Image, dstX, dstY, srcX, srcY, dstW, dstH, srcW, srcH int) {
    C.gdImageCopyResized(dst.img, p.img, C.int(dstX), C.int(dstY), C.int(srcX), C.int(srcY),
        C.int(dstW), C.int(dstH), C.int(srcW), C.int(srcH))
}

func (p *Image) CopyMerge(dst *Image, dstX, dstY, srcX, srcY, w, h, pct int) {
    C.gdImageCopyMerge(dst.img, p.img, C.int(dstX), C.int(dstY), C.int(srcX), C.int(srcY),
        C.int(w), C.int(h), C.int(pct))
}

func (p *Image) CopyMergeGray(dst *Image, dstX, dstY, srcX, srcY, w, h, pct int) {
    C.gdImageCopyMergeGray(dst.img, p.img, C.int(dstX), C.int(dstY), C.int(srcX), C.int(srcY),
        C.int(w), C.int(h), C.int(pct))
}

func (p *Image) Copy(dst *Image, dstX, dstY, srcX, srcY, w, h int) {
    C.gdImageCopy(dst.img, p.img, C.int(dstX), C.int(dstY), C.int(srcX), C.int(srcY), C.int(w), C.int(h))
}

func (p *Image) CopyRotated(dst *Image, dstX, dstY, srcX, srcY, srcWidth, srcHeight, angle int) {
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

func (p *Image) FilledEllipse(cx, cy, w, h int, color Color) {
    C.gdImageFilledEllipse(p.img, C.int(cx), C.int(cy), C.int(w), C.int(h), C.int(color))
}

/* NB: unable to import gdImageEllipse. Something wrong in CGO I think:

Undefined symbols:
  "_gdImageEllipse", referenced from:
        __cgo_6af175c8cf06_Cfunc_gdImageEllipse in gd.cgo2.o
             (maybe you meant: __cgo_6af175c8cf06_Cfunc_gdImageEllipse)
ld: symbol(s) not found
collect2: ld returned 1 exit status
*/
func (p *Image) Ellipse(cx, cy, w, h int, color Color) {
    C.gdImageArc(p.img, C.int(cx), C.int(cy), C.int(w), C.int(h), 360, 360, C.int(color))
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

func (p *Image) StringFT(fg Color, fontname string, ptsize, angle float64, x, y int, str string) (brect [8]int) {
    C.gdFontCacheSetup()
    defer C.gdFontCacheShutdown()

    C.gdImageStringFT(p.img, (*C.int)(Pointer(&brect)), C.int(fg), C.CString(fontname), C.double(ptsize),
        C.double(angle), C.int(x), C.int(y), C.CString(str))

    return
}

func (p *Image) Polygon(points [](struct{x, y int}), c Color) {
    C.gdImagePolygon(p.img, (C.gdPointPtr)(Pointer(&points)), C.int(len(points)), C.int(c))
}

func (p *Image) OpenPolygon(points [](struct{x, y int}), c Color) {
    C.gdImageOpenPolygon(p.img, (C.gdPointPtr)(Pointer(&points)), C.int(len(points)), C.int(c))
}

func (p *Image) FilledPolygon(points [](struct{x, y int}), c Color) {
    C.gdImageFilledPolygon(p.img, (C.gdPointPtr)(Pointer(&points)), C.int(len(points)), C.int(c))
}

func (p *Image) ColorsForIndex(index Color) map[string]int {
    if p.TrueColor() {
        return map[string]int{
            "alpha":  (int(index) & 0x7F000000) >> 24,
            "red":(int(index) & 0xFF0000) >> 16,
            "green": (int(index) & 0x00FF00) >> 8,
            "blue":(int(index) & 0x0000FF),
        }
    }

    return map[string]int{
        "red":      (int)((*p.img).red[index]),
        "green":    (int)((*p.img).green[index]),
        "blue":     (int)((*p.img).blue[index]),
        "alpha":    (int)((*p.img).alpha[index]),
    }
}

func searchfonts(dir string) (out []string){
    files, e := ioutil.ReadDir(dir)
    if e == nil {
        for _, file := range files {
            switch {
            case file.IsDirectory():
                entry := filepath.Join(dir, file.Name)

                out = append(out, searchfonts(entry)...)

            case file.IsRegular():
                ext := strings.ToLower(filepath.Ext(file.Name)[1:])
                whitelist := []string{"ttf", "otf", "cid", "cff", "pcf", "fnt", "bdr", "pfr", "pfa", "pfb", "afm"}

                for _, wext := range whitelist {
                    if ext == wext {
                        out = append(out, file.Name)
                        break
                    }
                }
            }
        }
    }

    return
}

func GetFonts() (list []string) {
    for _, dir := range strings.Split(C.DEFAULT_FONTPATH, C.PATHSEPARATOR, -1) {
        list = append(list, searchfonts(dir)...)
    }

    return
}

func (p *Image) getpixelfunc() (func(p *Image, x, y int) Color) {
    if p.TrueColor() {
        return func(p *Image, x, y int) Color { return p.GetTrueColorPixel(x,y) }
    }

    return func(p *Image, x, y int) Color { return p.GetPixel(x,y) }
}

func (p *Image) filter(flt func(r, g, b, a, x, y int) (int, int, int, int)) {
    f := p.getpixelfunc()

    sx, sy := p.Sx(), p.Sy()
    for y := 0; y<sy; y++ {
        for x := 0; x<sx; x++ {
            rgba := p.ColorsForIndex(f(p, x, y))
            r, g, b, a := flt(rgba["red"], rgba["green"], rgba["blue"], rgba["alpha"], x, y)

            newpxl := p.ColorAllocateAlpha(r, g, b, a)
            if newpxl == -1 {
                newpxl = p.ColorClosestAlpha(r, g, b, a)
            }

            p.SetPixel(x, y, newpxl)
        }
    }
}

func (p *Image) GrayScale() {
    p.filter(func(r, g, b, a, x, y int) (int, int, int, int) {
        c := (int) (.299 * float64(r) + .587 * float64(g) + .114 * float64(b))

        return c, c, c, a
    })
}

func (p *Image) Negate() {
    p.filter(func(r, g, b, a, x, y int) (int, int, int, int) {
        r = 255 - r
        g = 255 - g
        b = 255 - b

        return r, g, b, a
    })
}

func min(n1, n2 int) int {
    if n1 < n2 {
        return n1
    }

    return n2
}

func max(n1, n2 int) int {
    if n1 > n2 {
        return n1
    }

    return n2
}

func (p *Image) Brightness(brightness int) {
    if brightness == 0 {
        return
    }

    p.filter(func(r, g, b, a, x, y int) (int, int, int, int) {
        r = min(255, max(r + brightness, 0))
        g = min(255, max(g + brightness, 0))
        b = min(255, max(b + brightness, 0))

        return r, g, b, a
    })
}

func (p *Image) Contrast(contrast float64) {
    contrast = (100.0 - contrast) / 100.0
    contrast *= contrast

    corr := func(c int, contrast float64) int {
        f := float64(c) / 255.0
        f -= .5
        f *= contrast
        f += .5
        f *= 255.0

        return min(255, max(0, int(f)))
    }

    p.filter(func(r, g, b, a, x, y int) (int, int, int, int) {
        r = corr(r, contrast)
        g = corr(g, contrast)
        b = corr(b, contrast)

        return r, g, b, a
    })
}

func (p *Image) Color(r, g, b, a int) {
    p.filter(func(ri, gi, bi, ai, x, y int) (int, int, int, int) {
        ri = max(0, min(255, r+ri))
        gi = max(0, min(255, g+gi))
        bi = max(0, min(255, b+bi))
        ai = max(0, min(255, a+ai))

        return ri, gi, bi, ai
    })
}

func (p *Image) Convolution(filter [3][3]float32, filter_div, offset float32) {
    sx, sy := p.Sx(), p.Sy()
    srcback := CreateTrueColor(sx, sy)
    defer srcback.Destroy()

    srcback.SaveAlpha(true)
    newpxl := srcback.ColorAllocateAlpha(0, 0, 0, 127)
    srcback.Fill(0, 0, newpxl)

    p.Copy(srcback, 0, 0, 0, 0, sx, sy)

    var af func(p *Image, c int) int

    if p.TrueColor() {
        af = func(p *Image, c int) int { return (c & 0x7F000000) >> 24 }
    } else {
        af = func(p *Image, c int) int { return (int)((*p.img).alpha[c]) }
    }

    f := p.getpixelfunc()

    for y := 0; y<sy; y++ {
        for x := 0; x<sx; x++ {
            newr, newg, newb := float32(0), float32(0), float32(0)

            for j := 0; j<3; j++ {
                yv := min(max(y - 1 + j, 0), sy - 1)

                for i := 0; i<3; i++ {
                    pxl := srcback.ColorsForIndex(f(srcback, min(max(x - 1 + i, 0), sx - 1), yv))

                    newr += float32(pxl["red"]) * filter[j][i]
                    newg += float32(pxl["green"]) * filter[j][i]
                    newb += float32(pxl["blue"]) * filter[j][i]
                }
            }

            newr = (newr / filter_div) + offset
            newg = (newg / filter_div) + offset
            newb = (newb / filter_div) + offset

            r := min(255, max(0, int(newr)))
            g := min(255, max(0, int(newg)))
            b := min(255, max(0, int(newb)))

            newa := af(srcback, int(f(srcback, x, y)))

            newpxl = p.ColorAllocateAlpha(r, g, b, newa)
            if newpxl == -1 {
                newpxl = p.ColorClosestAlpha(r, g, b, newa)
            }

            p.SetPixel(x, y, newpxl)
        }
    }
}

func (p *Image) GaussianBlur() {
    filter := [3][3]float32{{1.0, 2.0, 3.0}, {2.0, 4.0, 2.0}, {1.0, 2.0, 1.0}}
    p.Convolution(filter, 16, 0)
}

func (p *Image) EdgeDetectQuick() {
    filter := [3][3]float32{{-1.0, 0.0, -1.0}, {0.0, 4.0, 0.0}, {-1.0, 0.0, -1.0}}
    p.Convolution(filter, 1, 127)
}

func (p *Image) Emboss() {
    filter := [3][3]float32{{1.5, 0.0, 0.0}, {0.0, 0.0, 0.0}, {0.0, 0.0, -1.5}}
    p.Convolution(filter, 1, 127)
}

func (p *Image) MeanRemoval() {
    filter := [3][3]float32{{-1, -1, -1}, {-1, 9, -1}, {-1, -1, -1}}
    p.Convolution(filter, 1, 0)
}

func (p *Image) Smooth(weight float32) {
    filter := [3][3]float32{{1, 1, 1}, {1, weight, 1}, {1, 1, 1}}
    p.Convolution(filter, weight + 8, 0)
}

// Stack Blur Algorithm by Mario Klingemann <mario@quasimondo.com>
// "Go" language port by Evgeny Stepanischev http://bolknote.ru
func (img *Image) StackBlur(radius int, keepalpha bool) {
    if radius < 1 {
        return
    }

    pix := img.getpixelfunc()

    w, h := int(img.Sx()), int(img.Sy())
    wm, hm, wh, div := w-1, h-1, w * h, radius * 2 + 1

    len := 4
    if keepalpha {
        len = 3
    }

    rgba := make([][]byte, len)
    for i := 0; i<len; i++ {
        rgba[i] = make([]byte, wh)
    }

    vmin := make([]int, max(w, h))

    var x, y, i, yp, yi, yw, stackpointer, stackstart, rbs int
    var sir *[4]byte

    divsum := (div + 1) >> 1
    divsum *= divsum

    dv := make([]byte, 256 * divsum)

    for i = 0; i<256 * divsum; i++ {
        dv[i] = byte(i / divsum)
    }

    yw, yi = 0, 0
    stack := make([][4]byte, div)
    r1 := radius + 1

    for y = 0; y<h; y++ {
        sum    := make([]int, len)
        insum  := make([]int, len)
        outsum := make([]int, len)

        for i = -radius; i<=radius; i++ {
            coords := yi + min(wm, max(i, 0))
            yc := coords / w
            xc := coords % w

            p := img.ColorsForIndex(pix(img, xc, yc))

            sir = &stack[i + radius]
            sir[0] = (byte)(p["red"])
            sir[1] = (byte)(p["green"])
            sir[2] = (byte)(p["blue"])
            sir[3] = (byte)(p["alpha"])

            rbs = r1 - abs(i)
            for i := 0; i<len; i++ {
                sum[i] += int(sir[i]) * rbs
            }

            if i > 0 {
                for i := 0; i<len; i++ {
                    insum[i] += int(sir[i])
                }
            } else {
                for i := 0; i<len; i++ {
                    outsum[i] += int(sir[i])
                }
            }
        }

        stackpointer = radius

        for x = 0; x<w; x++ {
            for i := 0; i<len; i++ {
                rgba[i][yi] = dv[sum[i]]
                sum[i] -= outsum[i]
            }

            stackstart = stackpointer - radius + div
            sir = &stack[stackstart % div]

            for i := 0; i<len; i++ {
                outsum[i] -= int(sir[i])
            }

            if y == 0 {
                vmin[x] = min(x + radius + 1, wm)
            }

            coords := yw + vmin[x]
            yc := coords / w
            xc := coords % w

            p := img.ColorsForIndex(pix(img, xc, yc))

            sir[0] = byte(p["red"])
            sir[1] = byte(p["green"])
            sir[2] = byte(p["blue"])
            sir[3] = byte(p["alpha"])

            for i := 0; i<len; i++ {
                insum[i] += int(sir[i])
                sum[i] += insum[i]
            }

            stackpointer = (stackpointer + 1) % div
            sir = &stack[stackpointer % div]

            for i := 0; i<len; i++ {
                outsum[i] += int(sir[i])
                insum[i] -= int(sir[i])
            }

            yi++
        }

        yw += w
    }

    for x = 0; x<w; x++ {
        sum    := make([]int, len)
        insum  := make([]int, len)
        outsum := make([]int, len)

        yp = -radius * w

        for i = -radius; i <= radius; i++ {
            yi = max(0, yp) + x

            sir = &stack[i + radius]

            for i := 0; i<len; i++ {
                sir[i] = rgba[i][yi]
            }
            rbs = r1 - abs(i)

            for i := 0; i<len; i++ {
                sum[i] += int(rgba[i][yi]) * rbs
            }

            if i > 0 {
                for i := 0; i<len; i++ {
                    insum[i] += int(sir[i])
                }
            } else {
                for i := 0; i<len; i++ {
                    outsum[i] += int(sir[i])
                }
            }

            if i < hm {
                yp += w
            }
        }

        yi = x

        stackpointer = radius

        for y = 0; y < h; y++ {
            var alpha byte

            if keepalpha {
                pxl := img.ColorsForIndex(pix(img, yi % w, yi / w))
                alpha = byte(pxl["alpha"])
            } else {
                alpha = dv[sum[3]]
            }

            newpxl := img.ColorAllocateAlpha(int(dv[sum[0]]), int(dv[sum[1]]), int(dv[sum[2]]), int(alpha))
            if newpxl == -1 {
                newpxl = img.ColorClosestAlpha(int(dv[sum[0]]), int(dv[sum[1]]), int(dv[sum[2]]), int(alpha))
            }

            img.SetPixel(yi % w, yi / w, newpxl)

            for i := 0; i<len; i++ {
                sum[i] -= outsum[i]
            }

            stackstart = stackpointer - radius + div
            sir = &stack[stackstart % div]

            for i := 0; i<len; i++ {
                outsum[i] -= int(sir[i])
            }

            if x == 0 {
                vmin[y] = min(y + r1, hm) * w
            }

            p := x + vmin[y]

            for i := 0; i<len; i++ {
                sir[i] = rgba[i][p]
                insum[i] += int(sir[i])
                sum[i] += insum[i]
            }

            stackpointer = (stackpointer + 1) % div
            sir = &stack[stackpointer]

            for i := 0; i<len; i++ {
                outsum[i] += int(sir[i])
                insum[i] -= int(sir[i])
            }

            yi += w
        }
    }
}

func abs(i int) int {
    if i < 0 {
        return -i
    }

    return i
}

