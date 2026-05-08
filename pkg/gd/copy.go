package gd

/*
#include <gd.h>
*/
import "C"

func (im *Image) CopyFrom(src *Image, dstX, dstY, srcX, srcY, width, height int) error {
	dstPtr, err := im.cptr()
	if err != nil {
		return err
	}
	srcPtr, err := src.cptr()
	if err != nil {
		return err
	}
	C.gdImageCopy(dstPtr, srcPtr, C.int(dstX), C.int(dstY), C.int(srcX), C.int(srcY), C.int(width), C.int(height))
	return nil
}

func (im *Image) CopyResizedFrom(src *Image, dstX, dstY, srcX, srcY, dstW, dstH, srcW, srcH int) error {
	dstPtr, err := im.cptr()
	if err != nil {
		return err
	}
	srcPtr, err := src.cptr()
	if err != nil {
		return err
	}
	C.gdImageCopyResized(dstPtr, srcPtr, C.int(dstX), C.int(dstY), C.int(srcX), C.int(srcY), C.int(dstW), C.int(dstH), C.int(srcW), C.int(srcH))
	return nil
}

func (im *Image) CopyResampledFrom(src *Image, dstX, dstY, srcX, srcY, dstW, dstH, srcW, srcH int) error {
	dstPtr, err := im.cptr()
	if err != nil {
		return err
	}
	srcPtr, err := src.cptr()
	if err != nil {
		return err
	}
	C.gdImageCopyResampled(dstPtr, srcPtr, C.int(dstX), C.int(dstY), C.int(srcX), C.int(srcY), C.int(dstW), C.int(dstH), C.int(srcW), C.int(srcH))
	return nil
}

func (im *Image) CopyMergeFrom(src *Image, dstX, dstY, srcX, srcY, width, height, pct int) error {
	dstPtr, err := im.cptr()
	if err != nil {
		return err
	}
	srcPtr, err := src.cptr()
	if err != nil {
		return err
	}
	C.gdImageCopyMerge(dstPtr, srcPtr, C.int(dstX), C.int(dstY), C.int(srcX), C.int(srcY), C.int(width), C.int(height), C.int(pct))
	return nil
}

func (im *Image) CopyMergeGrayFrom(src *Image, dstX, dstY, srcX, srcY, width, height, pct int) error {
	dstPtr, err := im.cptr()
	if err != nil {
		return err
	}
	srcPtr, err := src.cptr()
	if err != nil {
		return err
	}
	C.gdImageCopyMergeGray(dstPtr, srcPtr, C.int(dstX), C.int(dstY), C.int(srcX), C.int(srcY), C.int(width), C.int(height), C.int(pct))
	return nil
}

func (im *Image) CopyRotatedFrom(src *Image, dstX, dstY float64, srcX, srcY, srcWidth, srcHeight, angle int) error {
	dstPtr, err := im.cptr()
	if err != nil {
		return err
	}
	srcPtr, err := src.cptr()
	if err != nil {
		return err
	}
	C.gdImageCopyRotated(dstPtr, srcPtr, C.double(dstX), C.double(dstY), C.int(srcX), C.int(srcY), C.int(srcWidth), C.int(srcHeight), C.int(angle))
	return nil
}
