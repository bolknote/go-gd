package gd
// #include <gd.h>
import "C"
import "os"

type Image struct {img C.gdImagePtr}

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

func (p *Image) Jpeg(out string, quality int) {
    file := C.fopen(C.CString(out), C.CString("wb"))
    if file != nil {
        defer C.fclose(file)

        C.gdImageJpeg(p.img, file, C.int(quality))
    }

    panic(os.NewError("Error occurred while opening file for writing."))
}

func (p *Image) CopyResampled(dst Image, dstX, dstY, srcX, srcY, dstW, dstH, srcW, srcH int) {
    C.gdImageCopyResampled(dst.img, p.img, C.int(dstX), C.int(dstY), C.int(srcX), C.int(srcY), C.int(dstW), C.int(dstH), C.int(srcW), C.int(srcH))
}
