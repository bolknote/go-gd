package gd

/*
#include <gd.h>
*/
import "C"

type GIFDisposal int

const (
	GIFDisposalUnknown GIFDisposal = iota
	GIFDisposalNone
	GIFDisposalRestoreBackground
	GIFDisposalRestorePrevious
)

type GIFFrameOptions struct {
	LocalColorMap bool
	LeftOffset    int
	TopOffset     int
	Delay         int
	Disposal      GIFDisposal
	Previous      *Image
}

func (im *Image) GIFAnimBegin(globalColorMap bool, loops int) ([]byte, error) {
	ptr, err := im.cptr()
	if err != nil {
		return nil, err
	}
	var size C.int
	return imageBytes(C.gdImageGifAnimBeginPtr(ptr, &size, boolInt(globalColorMap), C.int(loops)), size)
}

func (im *Image) GIFAnimAdd(opts GIFFrameOptions) ([]byte, error) {
	ptr, err := im.cptr()
	if err != nil {
		return nil, err
	}
	var prev C.gdImagePtr
	if opts.Previous != nil {
		prev, err = opts.Previous.cptr()
		if err != nil {
			return nil, err
		}
	}
	var size C.int
	return imageBytes(C.gdImageGifAnimAddPtr(ptr, &size, boolInt(opts.LocalColorMap), C.int(opts.LeftOffset), C.int(opts.TopOffset), C.int(opts.Delay), C.int(opts.Disposal), prev), size)
}

func GIFAnimEnd() ([]byte, error) {
	var size C.int
	return imageBytes(C.gdImageGifAnimEndPtr(&size), size)
}

func (im *Image) GIFAnimBeginFile(path string, globalColorMap bool, loops int) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	return withWriteFile(path, func(file *C.FILE) error {
		C.gdImageGifAnimBegin(ptr, file, boolInt(globalColorMap), C.int(loops))
		return nil
	})
}

func (im *Image) GIFAnimAddFile(path string, opts GIFFrameOptions) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	var prev C.gdImagePtr
	if opts.Previous != nil {
		prev, err = opts.Previous.cptr()
		if err != nil {
			return err
		}
	}
	return withAppendFile(path, func(file *C.FILE) error {
		C.gdImageGifAnimAdd(ptr, file, boolInt(opts.LocalColorMap), C.int(opts.LeftOffset), C.int(opts.TopOffset), C.int(opts.Delay), C.int(opts.Disposal), prev)
		return nil
	})
}

func GIFAnimEndFile(path string) error {
	return withAppendFile(path, func(file *C.FILE) error {
		C.gdImageGifAnimEnd(file)
		return nil
	})
}

func (im *Image) GIFAnimBeginContext(ctx *IOContext, globalColorMap bool, loops int) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	out, err := ctx.cptr()
	if err != nil {
		return err
	}
	C.gdImageGifAnimBeginCtx(ptr, out, boolInt(globalColorMap), C.int(loops))
	return nil
}

func (im *Image) GIFAnimAddContext(ctx *IOContext, opts GIFFrameOptions) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	out, err := ctx.cptr()
	if err != nil {
		return err
	}
	var prev C.gdImagePtr
	if opts.Previous != nil {
		prev, err = opts.Previous.cptr()
		if err != nil {
			return err
		}
	}
	C.gdImageGifAnimAddCtx(ptr, out, boolInt(opts.LocalColorMap), C.int(opts.LeftOffset), C.int(opts.TopOffset), C.int(opts.Delay), C.int(opts.Disposal), prev)
	return nil
}

func GIFAnimEndContext(ctx *IOContext) error {
	out, err := ctx.cptr()
	if err != nil {
		return err
	}
	C.gdImageGifAnimEndCtx(out)
	return nil
}
