package gd

/*
#include <stdlib.h>
#include <stdio.h>
#include <gd.h>

static void go_gd_ctx_free(gdIOCtxPtr ctx) {
	if (ctx != NULL && ctx->gd_free != NULL) {
		ctx->gd_free(ctx);
	}
}
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// IOContext wraps libgd gdIOCtx. Close must be called unless ownership was
// consumed by Extract.
type IOContext struct {
	ptr       C.gdIOCtxPtr
	file      *C.FILE
	filePath  string
	inputData unsafe.Pointer
	dynamic   bool
	closed    bool
}

// NewDynamicReadContext creates a gdIOCtx that reads from a copy of data.
func NewDynamicReadContext(data []byte) (*IOContext, error) {
	if len(data) == 0 {
		return nil, ErrInvalidArgument
	}
	buf := C.CBytes(data)
	ctx := C.gdNewDynamicCtxEx(C.int(len(data)), buf, 0)
	if ctx == nil {
		C.free(buf)
		return nil, ErrInvalidArgument
	}
	return withIOContextFinalizer(&IOContext{ptr: ctx, inputData: buf}), nil
}

// NewDynamicWriteContext creates a gdIOCtx backed by libgd-managed memory
// that grows as encoders write to it. Use Extract to retrieve the bytes.
func NewDynamicWriteContext() (*IOContext, error) {
	buf := C.malloc(1)
	if buf == nil {
		return nil, ErrInvalidArgument
	}
	ctx := C.gdNewDynamicCtx(1, buf)
	if ctx == nil {
		C.free(buf)
		return nil, ErrInvalidArgument
	}
	return withIOContextFinalizer(&IOContext{ptr: ctx, dynamic: true}), nil
}

// NewFileContext creates a gdIOCtx around a C FILE opened for path with mode.
// Close on the returned context closes the underlying file.
func NewFileContext(path, mode string) (*IOContext, error) {
	file, err := fopen(path, mode)
	if err != nil {
		return nil, err
	}
	ctx := C.gdNewFileCtx(file)
	if ctx == nil {
		C.fclose(file)
		return nil, ErrInvalidArgument
	}
	return withIOContextFinalizer(&IOContext{ptr: ctx, file: file, filePath: path}), nil
}

func withIOContextFinalizer(ctx *IOContext) *IOContext {
	runtime.SetFinalizer(ctx, func(ctx *IOContext) { _ = ctx.Close() })
	return ctx
}

func (ctx *IOContext) cptr() (C.gdIOCtxPtr, error) {
	if ctx == nil || ctx.ptr == nil || ctx.closed {
		return nil, ErrInvalidArgument
	}
	return ctx.ptr, nil
}

// Close releases the libgd resources backing the context.
//
// For file-backed contexts the underlying FILE is also closed and any pending
// I/O error is reported. Calling Close more than once is a no-op.
func (ctx *IOContext) Close() error {
	if ctx == nil || ctx.closed {
		return nil
	}
	runtime.SetFinalizer(ctx, nil)
	ctx.closed = true
	if ctx.ptr != nil {
		C.go_gd_ctx_free(ctx.ptr)
		ctx.ptr = nil
	}
	var err error
	if ctx.file != nil {
		err = closeFile(ctx.filePath, ctx.file)
		ctx.file = nil
	}
	if ctx.inputData != nil {
		C.free(ctx.inputData)
		ctx.inputData = nil
	}
	return err
}

// Extract returns the bytes accumulated by a dynamic write context and
// consumes the context. The returned context is closed; Extract may not be
// called on read or file-backed contexts.
//
// libgd's gdDPExtractData transfers ownership of both the data buffer and
// the context itself, so we must not call gd_free on the context afterwards.
func (ctx *IOContext) Extract() ([]byte, error) {
	ptr, err := ctx.cptr()
	if err != nil {
		return nil, err
	}
	if !ctx.dynamic {
		return nil, ErrInvalidArgument
	}
	var size C.int
	data := C.gdDPExtractData(ptr, &size)
	runtime.SetFinalizer(ctx, nil)
	ctx.ptr = nil
	ctx.closed = true
	if data == nil || size <= 0 {
		return nil, ErrEncode
	}
	defer C.gdFree(data)
	return C.GoBytes(data, size), nil
}

func DecodePNGContext(ctx *IOContext) (*Image, error) {
	ptr, err := ctx.cptr()
	if err != nil {
		return nil, err
	}
	return wrapDecodedImage(C.gdImageCreateFromPngCtx(ptr))
}

func DecodeGIFContext(ctx *IOContext) (*Image, error) {
	ptr, err := ctx.cptr()
	if err != nil {
		return nil, err
	}
	return wrapDecodedImage(C.gdImageCreateFromGifCtx(ptr))
}

func DecodeWBMPContext(ctx *IOContext) (*Image, error) {
	ptr, err := ctx.cptr()
	if err != nil {
		return nil, err
	}
	return wrapDecodedImage(C.gdImageCreateFromWBMPCtx(ptr))
}

func DecodeJPEGContext(ctx *IOContext, opts *JPEGOptions) (*Image, error) {
	ptr, err := ctx.cptr()
	if err != nil {
		return nil, err
	}
	o := defaultJPEGOptions(opts)
	return wrapDecodedImage(C.gdImageCreateFromJpegCtxEx(ptr, boolInt(o.IgnoreWarnings)))
}

func DecodeWebPContext(ctx *IOContext) (*Image, error) {
	ptr, err := ctx.cptr()
	if err != nil {
		return nil, err
	}
	return wrapDecodedImage(C.gdImageCreateFromWebpCtx(ptr))
}

func DecodeHEIFContext(ctx *IOContext) (*Image, error) {
	ptr, err := ctx.cptr()
	if err != nil {
		return nil, err
	}
	return wrapDecodedImage(C.gdImageCreateFromHeifCtx(ptr))
}

func DecodeAVIFContext(ctx *IOContext) (*Image, error) {
	ptr, err := ctx.cptr()
	if err != nil {
		return nil, err
	}
	return wrapDecodedImage(C.gdImageCreateFromAvifCtx(ptr))
}

func DecodeTIFFContext(ctx *IOContext) (*Image, error) {
	ptr, err := ctx.cptr()
	if err != nil {
		return nil, err
	}
	return wrapDecodedImage(C.gdImageCreateFromTiffCtx(ptr))
}

func DecodeTGAContext(ctx *IOContext) (*Image, error) {
	ptr, err := ctx.cptr()
	if err != nil {
		return nil, err
	}
	return wrapDecodedImage(C.gdImageCreateFromTgaCtx(ptr))
}

func DecodeBMPContext(ctx *IOContext) (*Image, error) {
	ptr, err := ctx.cptr()
	if err != nil {
		return nil, err
	}
	return wrapDecodedImage(C.gdImageCreateFromBmpCtx(ptr))
}

func DecodeGDContext(ctx *IOContext) (*Image, error) {
	ptr, err := ctx.cptr()
	if err != nil {
		return nil, err
	}
	return wrapDecodedImage(C.gdImageCreateFromGdCtx(ptr))
}

func DecodeGD2Context(ctx *IOContext) (*Image, error) {
	ptr, err := ctx.cptr()
	if err != nil {
		return nil, err
	}
	return wrapDecodedImage(C.gdImageCreateFromGd2Ctx(ptr))
}

func DecodeGD2PartContext(ctx *IOContext, srcX, srcY, width, height int) (*Image, error) {
	ptr, err := ctx.cptr()
	if err != nil {
		return nil, err
	}
	return wrapDecodedImage(C.gdImageCreateFromGd2PartCtx(ptr, C.int(srcX), C.int(srcY), C.int(width), C.int(height)))
}

func (im *Image) EncodePNGContext(ctx *IOContext, opts *PNGOptions) error {
	img, err := im.cptr()
	if err != nil {
		return err
	}
	out, err := ctx.cptr()
	if err != nil {
		return err
	}
	o := defaultPNGOptions(opts)
	if o.Compression < 0 {
		C.gdImagePngCtx(img, out)
	} else {
		C.gdImagePngCtxEx(img, out, C.int(o.Compression))
	}
	return nil
}

func (im *Image) EncodeGIFContext(ctx *IOContext) error {
	img, err := im.cptr()
	if err != nil {
		return err
	}
	out, err := ctx.cptr()
	if err != nil {
		return err
	}
	C.gdImageGifCtx(img, out)
	return nil
}

func (im *Image) EncodeTIFFContext(ctx *IOContext) error {
	img, err := im.cptr()
	if err != nil {
		return err
	}
	out, err := ctx.cptr()
	if err != nil {
		return err
	}
	C.gdImageTiffCtx(img, out)
	return nil
}

func (im *Image) EncodeBMPContext(ctx *IOContext, opts *BMPOptions) error {
	img, err := im.cptr()
	if err != nil {
		return err
	}
	out, err := ctx.cptr()
	if err != nil {
		return err
	}
	compression := C.int(defaultBMPOptions(opts).Compression)
	C.gdImageBmpCtx(img, out, compression)
	return nil
}

func (im *Image) EncodeWBMPContext(ctx *IOContext, foreground Color) error {
	img, err := im.cptr()
	if err != nil {
		return err
	}
	out, err := ctx.cptr()
	if err != nil {
		return err
	}
	C.gdImageWBMPCtx(img, C.int(foreground), out)
	return nil
}

func (im *Image) EncodeJPEGContext(ctx *IOContext, opts *JPEGOptions) error {
	img, err := im.cptr()
	if err != nil {
		return err
	}
	out, err := ctx.cptr()
	if err != nil {
		return err
	}
	o := defaultJPEGOptions(opts)
	C.gdImageJpegCtx(img, out, C.int(o.Quality))
	return nil
}

func (im *Image) EncodeWebPContext(ctx *IOContext, opts *WebPOptions) error {
	img, err := im.cptr()
	if err != nil {
		return err
	}
	out, err := ctx.cptr()
	if err != nil {
		return err
	}
	o := defaultWebPOptions(opts)
	C.gdImageWebpCtx(img, out, C.int(o.Quality))
	return nil
}

func (im *Image) EncodeHEIFContext(ctx *IOContext, opts *HEIFOptions) error {
	img, err := im.cptr()
	if err != nil {
		return err
	}
	out, err := ctx.cptr()
	if err != nil {
		return err
	}
	o := defaultHEIFOptions(opts)
	chroma := C.CString(string(o.Chroma))
	defer C.free(unsafe.Pointer(chroma))
	C.gdImageHeifCtx(img, out, C.int(o.Quality), C.gdHeifCodec(o.Codec), C.gdHeifChroma(chroma))
	return nil
}

func (im *Image) EncodeAVIFContext(ctx *IOContext, opts *AVIFOptions) error {
	img, err := im.cptr()
	if err != nil {
		return err
	}
	out, err := ctx.cptr()
	if err != nil {
		return err
	}
	o := defaultAVIFOptions(opts)
	C.gdImageAvifCtx(img, out, C.int(o.Quality), C.int(o.Speed))
	return nil
}

func (im *Image) EncodeXBMContext(ctx *IOContext, name string, foreground Color) error {
	img, err := im.cptr()
	if err != nil {
		return err
	}
	out, err := ctx.cptr()
	if err != nil {
		return err
	}
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	C.gdImageXbmCtx(img, cname, C.int(foreground), out)
	return nil
}

// EncodeXBM serialises the image to an XBM-format byte slice. libgd only
// exposes XBM encoding through gdIOCtx, so this helper allocates a transient
// dynamic context internally.
func (im *Image) EncodeXBM(name string, foreground Color) ([]byte, error) {
	if _, err := im.cptr(); err != nil {
		return nil, err
	}
	ctx, err := NewDynamicWriteContext()
	if err != nil {
		return nil, err
	}
	if err := im.EncodeXBMContext(ctx, name, foreground); err != nil {
		_ = ctx.Close()
		return nil, err
	}
	return ctx.Extract()
}

// EncodeXBMFile is the file-backed counterpart to EncodeXBM.
func (im *Image) EncodeXBMFile(path, name string, foreground Color) error {
	if _, err := im.cptr(); err != nil {
		return err
	}
	ctx, err := NewFileContext(path, "wb")
	if err != nil {
		return err
	}
	if err := im.EncodeXBMContext(ctx, name, foreground); err != nil {
		_ = ctx.Close()
		return err
	}
	return ctx.Close()
}
