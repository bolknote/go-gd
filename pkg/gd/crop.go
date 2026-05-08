package gd

/*
#include <gd.h>

static gdRect go_gd_crop_rect(int x, int y, int w, int h) {
	gdRect r;
	r.x = x;
	r.y = y;
	r.width = w;
	r.height = h;
	return r;
}
*/
import "C"

type CropMode int

const (
	CropDefault CropMode = iota
	CropTransparent
	CropBlack
	CropWhite
	CropSides
	CropThreshold
)

// Crop returns a copy of the image cropped to rect.
func (im *Image) Crop(rect Rect) (*Image, error) {
	ptr, err := im.cptr()
	if err != nil {
		return nil, err
	}
	if rect.Width <= 0 || rect.Height <= 0 {
		return nil, ErrInvalidArgument
	}
	crop := C.go_gd_crop_rect(C.int(rect.X), C.int(rect.Y), C.int(rect.Width), C.int(rect.Height))
	out := C.gdImageCrop(ptr, &crop)
	if out == nil {
		return nil, ErrCrop
	}
	return wrapImage(out)
}

// CropAuto crops the image automatically using one of libgd's heuristics.
func (im *Image) CropAuto(mode CropMode) (*Image, error) {
	ptr, err := im.cptr()
	if err != nil {
		return nil, err
	}
	out := C.gdImageCropAuto(ptr, C.uint(mode))
	if out == nil {
		return nil, ErrCrop
	}
	return wrapImage(out)
}

// CropThreshold crops away pixels that match color within the given
// threshold. color must be a non-negative truecolor or palette index.
func (im *Image) CropThreshold(color Color, threshold float32) (*Image, error) {
	ptr, err := im.cptr()
	if err != nil {
		return nil, err
	}
	if color < 0 {
		return nil, ErrInvalidArgument
	}
	out := C.gdImageCropThreshold(ptr, C.uint(color), C.float(threshold))
	if out == nil {
		return nil, ErrCrop
	}
	return wrapImage(out)
}
