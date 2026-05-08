package gd

import "errors"

var (
	ErrClosedImage        = errors.New("gd: image is closed")
	ErrDecode             = errors.New("gd: decode failed")
	ErrEncode             = errors.New("gd: encode failed")
	ErrImageCreate        = errors.New("gd: image creation failed")
	ErrInvalidArgument    = errors.New("gd: invalid argument")
	ErrUnsupportedFeature = errors.New("gd: unsupported by linked libgd")
	ErrPaletteFull        = errors.New("gd: palette is full (256 colors allocated)")
	ErrColorMatch         = errors.New("gd: color match failed")
	ErrCrop               = errors.New("gd: crop failed")
	ErrTransform          = errors.New("gd: transform failed")
)
