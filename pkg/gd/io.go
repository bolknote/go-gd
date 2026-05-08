package gd

/*
#include <errno.h>
#include <stdlib.h>
#include <stdio.h>
#include <gd.h>

static int go_gd_errno(void) { return errno; }
static int go_gd_ferror(FILE *f) { return ferror(f); }
*/
import "C"
import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"
)

type JPEGOptions struct {
	Quality        int
	IgnoreWarnings bool
}

type PNGOptions struct {
	Compression int
}

type WebPOptions struct {
	Quality int
}

type BMPOptions struct {
	Compression int
}

type GD2Options struct {
	ChunkSize int
	Format    int
}

type HEIFCodec int

const (
	HEIFCodecUnknown HEIFCodec = 0
	HEIFCodecHEVC    HEIFCodec = 1
	HEIFCodecAV1     HEIFCodec = 4
)

type HEIFChroma string

const (
	HEIFChroma420 HEIFChroma = "420"
	HEIFChroma422 HEIFChroma = "422"
	HEIFChroma444 HEIFChroma = "444"
)

type HEIFOptions struct {
	Quality int
	Codec   HEIFCodec
	Chroma  HEIFChroma
}

type AVIFOptions struct {
	Quality int
	Speed   int
}

func defaultJPEGOptions(opts *JPEGOptions) JPEGOptions {
	if opts == nil {
		return JPEGOptions{Quality: -1}
	}
	return *opts
}

func defaultPNGOptions(opts *PNGOptions) PNGOptions {
	if opts == nil {
		return PNGOptions{Compression: -1}
	}
	return *opts
}

func defaultWebPOptions(opts *WebPOptions) WebPOptions {
	if opts == nil {
		return WebPOptions{Quality: -1}
	}
	return *opts
}

// defaultHEIFOptions fills in the structural defaults (libgd codec/chroma
// enums must be non-zero). Quality is left untouched: callers that pass
// HEIFOptions{} get Quality=0, which libgd treats as "lowest quality"; pass
// a nil pointer to get the libgd default of -1.
func defaultHEIFOptions(opts *HEIFOptions) HEIFOptions {
	out := HEIFOptions{Quality: -1, Codec: HEIFCodecHEVC, Chroma: HEIFChroma444}
	if opts != nil {
		out = *opts
		if out.Codec == HEIFCodecUnknown {
			out.Codec = HEIFCodecHEVC
		}
		if out.Chroma == "" {
			out.Chroma = HEIFChroma444
		}
	}
	return out
}

func defaultAVIFOptions(opts *AVIFOptions) AVIFOptions {
	if opts == nil {
		return AVIFOptions{Quality: -1, Speed: -1}
	}
	return *opts
}

func defaultBMPOptions(opts *BMPOptions) BMPOptions {
	if opts == nil {
		return BMPOptions{}
	}
	return *opts
}

func defaultGD2Options(opts *GD2Options) GD2Options {
	if opts == nil {
		return GD2Options{ChunkSize: 128, Format: 2}
	}
	return *opts
}

// openFile opens path through C fopen and returns the FILE pointer along with
// a closer that flushes the stream and reports any deferred I/O error. The
// closer must be invoked exactly once unless an error is returned.
func openFile(path, mode string) (*C.FILE, func() error, error) {
	file, err := fopen(path, mode)
	if err != nil {
		return nil, func() error { return nil }, err
	}
	return file, func() error { return closeFile(path, file) }, nil
}

func closeFile(path string, file *C.FILE) error {
	hadError := C.go_gd_ferror(file) != 0
	rc := C.fclose(file)
	errno := C.go_gd_errno()
	if rc != 0 {
		if errno != 0 {
			return fmt.Errorf("gd: close %q: %w", path, syscall.Errno(errno))
		}
		return fmt.Errorf("gd: close %q: %w", path, ErrEncode)
	}
	if hadError {
		return fmt.Errorf("gd: write %q: %w", path, ErrEncode)
	}
	return nil
}

func fopen(path, mode string) (*C.FILE, error) {
	cpath := C.CString(path)
	cmode := C.CString(mode)
	file := C.fopen(cpath, cmode)
	errno := C.go_gd_errno()
	C.free(unsafe.Pointer(cpath))
	C.free(unsafe.Pointer(cmode))
	if file == nil {
		return nil, openFileError(path, errno)
	}
	return file, nil
}

func openFileError(path string, errno C.int) error {
	if errno != 0 {
		return fmt.Errorf("gd: open %q: %w", path, syscall.Errno(errno))
	}
	return fmt.Errorf("gd: open %q: %w", path, ErrInvalidArgument)
}

// withReadFile invokes fn with an opened-for-read C FILE and always closes it
// afterwards, returning fn's error and ignoring close errors (the data, if
// any, has already been consumed into Go memory by fn).
func withReadFile[T any](path string, fn func(*C.FILE) (T, error)) (T, error) {
	var zero T
	file, cleanup, err := openFile(path, "rb")
	if err != nil {
		return zero, err
	}
	out, ferr := fn(file)
	_ = cleanup()
	return out, ferr
}

// withWriteFile invokes fn with an opened-for-write C FILE and propagates
// close errors so that buffered I/O failures (ENOSPC, EIO, ...) reach the
// caller.
func withWriteFile(path string, fn func(*C.FILE) error) error {
	file, cleanup, err := openFile(path, "wb")
	if err != nil {
		return err
	}
	if ferr := fn(file); ferr != nil {
		return errors.Join(ferr, cleanup())
	}
	return cleanup()
}

// withAppendFile is the equivalent of withWriteFile but uses append mode.
func withAppendFile(path string, fn func(*C.FILE) error) error {
	file, cleanup, err := openFile(path, "ab")
	if err != nil {
		return err
	}
	if ferr := fn(file); ferr != nil {
		return errors.Join(ferr, cleanup())
	}
	return cleanup()
}

func decodeBuffer(data []byte, fn func(C.int, unsafe.Pointer) C.gdImagePtr) (*Image, error) {
	if len(data) == 0 {
		return nil, ErrInvalidArgument
	}
	return wrapDecodedImage(fn(C.int(len(data)), unsafe.Pointer(&data[0])))
}

func imageBytes(ptr unsafe.Pointer, size C.int) ([]byte, error) {
	if ptr == nil {
		return nil, ErrEncode
	}
	defer C.gdFree(ptr)
	if size <= 0 {
		return nil, ErrEncode
	}
	return C.GoBytes(ptr, size), nil
}

func DecodeFile(path string) (*Image, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	return wrapDecodedImage(C.gdImageCreateFromFile(cpath))
}

func (im *Image) EncodeFile(path string) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	if C.gdImageFile(ptr, cpath) == 0 {
		return ErrEncode
	}
	return nil
}

func DecodePNG(data []byte) (*Image, error) {
	return decodeBuffer(data, func(size C.int, ptr unsafe.Pointer) C.gdImagePtr {
		return C.gdImageCreateFromPngPtr(size, ptr)
	})
}

func DecodePNGFile(path string) (*Image, error) {
	return withReadFile(path, func(file *C.FILE) (*Image, error) {
		return wrapDecodedImage(C.gdImageCreateFromPng(file))
	})
}

func (im *Image) EncodePNG(opts *PNGOptions) ([]byte, error) {
	ptr, err := im.cptr()
	if err != nil {
		return nil, err
	}
	var size C.int
	o := defaultPNGOptions(opts)
	if o.Compression < 0 {
		return imageBytes(C.gdImagePngPtr(ptr, &size), size)
	}
	return imageBytes(C.gdImagePngPtrEx(ptr, &size, C.int(o.Compression)), size)
}

func (im *Image) EncodePNGFile(path string, opts *PNGOptions) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	o := defaultPNGOptions(opts)
	return withWriteFile(path, func(file *C.FILE) error {
		if o.Compression < 0 {
			C.gdImagePng(ptr, file)
		} else {
			C.gdImagePngEx(ptr, file, C.int(o.Compression))
		}
		return nil
	})
}

func DecodeJPEG(data []byte, opts *JPEGOptions) (*Image, error) {
	o := defaultJPEGOptions(opts)
	return decodeBuffer(data, func(size C.int, ptr unsafe.Pointer) C.gdImagePtr {
		return C.gdImageCreateFromJpegPtrEx(size, ptr, boolInt(o.IgnoreWarnings))
	})
}

func DecodeJPEGFile(path string, opts *JPEGOptions) (*Image, error) {
	o := defaultJPEGOptions(opts)
	return withReadFile(path, func(file *C.FILE) (*Image, error) {
		return wrapDecodedImage(C.gdImageCreateFromJpegEx(file, boolInt(o.IgnoreWarnings)))
	})
}

func (im *Image) EncodeJPEG(opts *JPEGOptions) ([]byte, error) {
	ptr, err := im.cptr()
	if err != nil {
		return nil, err
	}
	var size C.int
	o := defaultJPEGOptions(opts)
	return imageBytes(C.gdImageJpegPtr(ptr, &size, C.int(o.Quality)), size)
}

func (im *Image) EncodeJPEGFile(path string, opts *JPEGOptions) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	o := defaultJPEGOptions(opts)
	return withWriteFile(path, func(file *C.FILE) error {
		C.gdImageJpeg(ptr, file, C.int(o.Quality))
		return nil
	})
}

func DecodeGIF(data []byte) (*Image, error) {
	return decodeBuffer(data, func(size C.int, ptr unsafe.Pointer) C.gdImagePtr {
		return C.gdImageCreateFromGifPtr(size, ptr)
	})
}

func DecodeGIFFile(path string) (*Image, error) {
	return withReadFile(path, func(file *C.FILE) (*Image, error) {
		return wrapDecodedImage(C.gdImageCreateFromGif(file))
	})
}

func (im *Image) EncodeGIF() ([]byte, error) {
	ptr, err := im.cptr()
	if err != nil {
		return nil, err
	}
	var size C.int
	return imageBytes(C.gdImageGifPtr(ptr, &size), size)
}

func (im *Image) EncodeGIFFile(path string) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	return withWriteFile(path, func(file *C.FILE) error {
		C.gdImageGif(ptr, file)
		return nil
	})
}

func DecodeWebP(data []byte) (*Image, error) {
	return decodeBuffer(data, func(size C.int, ptr unsafe.Pointer) C.gdImagePtr {
		return C.gdImageCreateFromWebpPtr(size, ptr)
	})
}

func DecodeWebPFile(path string) (*Image, error) {
	return withReadFile(path, func(file *C.FILE) (*Image, error) {
		return wrapDecodedImage(C.gdImageCreateFromWebp(file))
	})
}

func (im *Image) EncodeWebP(opts *WebPOptions) ([]byte, error) {
	ptr, err := im.cptr()
	if err != nil {
		return nil, err
	}
	var size C.int
	o := defaultWebPOptions(opts)
	if o.Quality < 0 {
		return imageBytes(C.gdImageWebpPtr(ptr, &size), size)
	}
	return imageBytes(C.gdImageWebpPtrEx(ptr, &size, C.int(o.Quality)), size)
}

func (im *Image) EncodeWebPFile(path string, opts *WebPOptions) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	o := defaultWebPOptions(opts)
	return withWriteFile(path, func(file *C.FILE) error {
		if o.Quality < 0 {
			C.gdImageWebp(ptr, file)
		} else {
			C.gdImageWebpEx(ptr, file, C.int(o.Quality))
		}
		return nil
	})
}

func DecodeWBMP(data []byte) (*Image, error) {
	return decodeBuffer(data, func(size C.int, ptr unsafe.Pointer) C.gdImagePtr {
		return C.gdImageCreateFromWBMPPtr(size, ptr)
	})
}

func DecodeWBMPFile(path string) (*Image, error) {
	return withReadFile(path, func(file *C.FILE) (*Image, error) {
		return wrapDecodedImage(C.gdImageCreateFromWBMP(file))
	})
}

func (im *Image) EncodeWBMP(foreground Color) ([]byte, error) {
	ptr, err := im.cptr()
	if err != nil {
		return nil, err
	}
	var size C.int
	return imageBytes(C.gdImageWBMPPtr(ptr, &size, C.int(foreground)), size)
}

func (im *Image) EncodeWBMPFile(path string, foreground Color) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	return withWriteFile(path, func(file *C.FILE) error {
		C.gdImageWBMP(ptr, C.int(foreground), file)
		return nil
	})
}

func DecodeBMP(data []byte) (*Image, error) {
	return decodeBuffer(data, func(size C.int, ptr unsafe.Pointer) C.gdImagePtr {
		return C.gdImageCreateFromBmpPtr(size, ptr)
	})
}

func DecodeBMPFile(path string) (*Image, error) {
	return withReadFile(path, func(file *C.FILE) (*Image, error) {
		return wrapDecodedImage(C.gdImageCreateFromBmp(file))
	})
}

func (im *Image) EncodeBMP(opts *BMPOptions) ([]byte, error) {
	ptr, err := im.cptr()
	if err != nil {
		return nil, err
	}
	compression := C.int(defaultBMPOptions(opts).Compression)
	var size C.int
	return imageBytes(C.gdImageBmpPtr(ptr, &size, compression), size)
}

func (im *Image) EncodeBMPFile(path string, opts *BMPOptions) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	compression := C.int(defaultBMPOptions(opts).Compression)
	return withWriteFile(path, func(file *C.FILE) error {
		C.gdImageBmp(ptr, file, compression)
		return nil
	})
}

func DecodeTGA(data []byte) (*Image, error) {
	return decodeBuffer(data, func(size C.int, ptr unsafe.Pointer) C.gdImagePtr {
		return C.gdImageCreateFromTgaPtr(size, ptr)
	})
}

func DecodeTGAFile(path string) (*Image, error) {
	return withReadFile(path, func(file *C.FILE) (*Image, error) {
		return wrapDecodedImage(C.gdImageCreateFromTga(file))
	})
}

func DecodeTIFF(data []byte) (*Image, error) {
	return decodeBuffer(data, func(size C.int, ptr unsafe.Pointer) C.gdImagePtr {
		return C.gdImageCreateFromTiffPtr(size, ptr)
	})
}

func DecodeTIFFFile(path string) (*Image, error) {
	return withReadFile(path, func(file *C.FILE) (*Image, error) {
		return wrapDecodedImage(C.gdImageCreateFromTiff(file))
	})
}

func (im *Image) EncodeTIFF() ([]byte, error) {
	ptr, err := im.cptr()
	if err != nil {
		return nil, err
	}
	var size C.int
	return imageBytes(C.gdImageTiffPtr(ptr, &size), size)
}

func (im *Image) EncodeTIFFFile(path string) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	return withWriteFile(path, func(file *C.FILE) error {
		C.gdImageTiff(ptr, file)
		return nil
	})
}

func DecodeGD(data []byte) (*Image, error) {
	return decodeBuffer(data, func(size C.int, ptr unsafe.Pointer) C.gdImagePtr {
		return C.gdImageCreateFromGdPtr(size, ptr)
	})
}

func DecodeGDFile(path string) (*Image, error) {
	return withReadFile(path, func(file *C.FILE) (*Image, error) {
		return wrapDecodedImage(C.gdImageCreateFromGd(file))
	})
}

func (im *Image) EncodeGD() ([]byte, error) {
	ptr, err := im.cptr()
	if err != nil {
		return nil, err
	}
	var size C.int
	return imageBytes(C.gdImageGdPtr(ptr, &size), size)
}

func (im *Image) EncodeGDFile(path string) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	return withWriteFile(path, func(file *C.FILE) error {
		C.gdImageGd(ptr, file)
		return nil
	})
}

func DecodeGD2(data []byte) (*Image, error) {
	return decodeBuffer(data, func(size C.int, ptr unsafe.Pointer) C.gdImagePtr {
		return C.gdImageCreateFromGd2Ptr(size, ptr)
	})
}

func DecodeGD2Part(data []byte, srcX, srcY, width, height int) (*Image, error) {
	return decodeBuffer(data, func(size C.int, ptr unsafe.Pointer) C.gdImagePtr {
		return C.gdImageCreateFromGd2PartPtr(size, ptr, C.int(srcX), C.int(srcY), C.int(width), C.int(height))
	})
}

func DecodeGD2PartFile(path string, srcX, srcY, width, height int) (*Image, error) {
	return withReadFile(path, func(file *C.FILE) (*Image, error) {
		return wrapDecodedImage(C.gdImageCreateFromGd2Part(file, C.int(srcX), C.int(srcY), C.int(width), C.int(height)))
	})
}

func DecodeGD2File(path string) (*Image, error) {
	return withReadFile(path, func(file *C.FILE) (*Image, error) {
		return wrapDecodedImage(C.gdImageCreateFromGd2(file))
	})
}

func (im *Image) EncodeGD2(opts *GD2Options) ([]byte, error) {
	ptr, err := im.cptr()
	if err != nil {
		return nil, err
	}
	o := defaultGD2Options(opts)
	chunkSize, format := C.int(o.ChunkSize), C.int(o.Format)
	var size C.int
	return imageBytes(C.gdImageGd2Ptr(ptr, chunkSize, format, &size), size)
}

func (im *Image) EncodeGD2File(path string, opts *GD2Options) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	o := defaultGD2Options(opts)
	chunkSize, format := C.int(o.ChunkSize), C.int(o.Format)
	return withWriteFile(path, func(file *C.FILE) error {
		C.gdImageGd2(ptr, file, chunkSize, format)
		return nil
	})
}

func DecodeHEIF(data []byte) (*Image, error) {
	return decodeBuffer(data, func(size C.int, ptr unsafe.Pointer) C.gdImagePtr {
		return C.gdImageCreateFromHeifPtr(size, ptr)
	})
}

func DecodeHEIFFile(path string) (*Image, error) {
	return withReadFile(path, func(file *C.FILE) (*Image, error) {
		return wrapDecodedImage(C.gdImageCreateFromHeif(file))
	})
}

func (im *Image) EncodeHEIF(opts *HEIFOptions) ([]byte, error) {
	ptr, err := im.cptr()
	if err != nil {
		return nil, err
	}
	var size C.int
	o := defaultHEIFOptions(opts)
	chroma := C.CString(string(o.Chroma))
	defer C.free(unsafe.Pointer(chroma))
	return imageBytes(C.gdImageHeifPtrEx(ptr, &size, C.int(o.Quality), C.gdHeifCodec(o.Codec), C.gdHeifChroma(chroma)), size)
}

func (im *Image) EncodeHEIFFile(path string, opts *HEIFOptions) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	o := defaultHEIFOptions(opts)
	chroma := C.CString(string(o.Chroma))
	defer C.free(unsafe.Pointer(chroma))
	return withWriteFile(path, func(file *C.FILE) error {
		C.gdImageHeifEx(ptr, file, C.int(o.Quality), C.gdHeifCodec(o.Codec), C.gdHeifChroma(chroma))
		return nil
	})
}

func DecodeAVIF(data []byte) (*Image, error) {
	return decodeBuffer(data, func(size C.int, ptr unsafe.Pointer) C.gdImagePtr {
		return C.gdImageCreateFromAvifPtr(size, ptr)
	})
}

func DecodeAVIFFile(path string) (*Image, error) {
	return withReadFile(path, func(file *C.FILE) (*Image, error) {
		return wrapDecodedImage(C.gdImageCreateFromAvif(file))
	})
}

func (im *Image) EncodeAVIF(opts *AVIFOptions) ([]byte, error) {
	ptr, err := im.cptr()
	if err != nil {
		return nil, err
	}
	var size C.int
	o := defaultAVIFOptions(opts)
	return imageBytes(C.gdImageAvifPtrEx(ptr, &size, C.int(o.Quality), C.int(o.Speed)), size)
}

func (im *Image) EncodeAVIFFile(path string, opts *AVIFOptions) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	o := defaultAVIFOptions(opts)
	return withWriteFile(path, func(file *C.FILE) error {
		C.gdImageAvifEx(ptr, file, C.int(o.Quality), C.int(o.Speed))
		return nil
	})
}

func DecodeXBMFile(path string) (*Image, error) {
	return withReadFile(path, func(file *C.FILE) (*Image, error) {
		return wrapDecodedImage(C.gdImageCreateFromXbm(file))
	})
}

func DecodeXPMFile(path string) (*Image, error) {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	return wrapDecodedImage(C.gdImageCreateFromXpm(cpath))
}
