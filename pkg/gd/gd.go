package gd

/*
#include <gd.h>

static int go_gd_image_sx(gdImagePtr im) { return gdImageSX(im); }
static int go_gd_image_sy(gdImagePtr im) { return gdImageSY(im); }
static int go_gd_image_true_color(gdImagePtr im) { return gdImageTrueColor(im); }
static int go_gd_image_colors_total(gdImagePtr im) { return gdImageColorsTotal(im); }
static int go_gd_image_transparent(gdImagePtr im) { return gdImageGetTransparent(im); }
static int go_gd_image_interlaced(gdImagePtr im) { return gdImageGetInterlaced(im); }
static int go_gd_image_resolution_x(gdImagePtr im) { return gdImageResolutionX(im); }
static int go_gd_image_resolution_y(gdImagePtr im) { return gdImageResolutionY(im); }
*/
import "C"
import "runtime"

// Image owns a gdImagePtr. Call Close when the image is no longer needed.
//
// Image values are not safe for concurrent use by multiple goroutines without
// external synchronization. In particular, two goroutines must not call any
// method on the same *Image (including Close) concurrently.
type Image struct {
	ptr C.gdImagePtr
	// tile and brush hold references to images registered via SetTile/SetBrush
	// so that libgd's borrowed pointer stays valid for as long as this image
	// keeps using it. Cleared by ClearTile/ClearBrush or Close.
	tile  *Image
	brush *Image
}

// Color is a libgd palette index, truecolor value, or special drawing color.
type Color int32

// RGBA stores libgd color channels. Alpha is in libgd's 0..127 range,
// where 0 is opaque and 127 is fully transparent.
type RGBA struct {
	R, G, B, A int
}

// Point represents an integer image coordinate.
type Point struct {
	X, Y int
}

// Rect represents an image rectangle.
type Rect struct {
	X, Y, Width, Height int
}

// ArcStyle controls filled arc rendering.
type ArcStyle int

const (
	ArcPie    ArcStyle = 0
	ArcChord  ArcStyle = 1
	ArcNoFill ArcStyle = 2
	ArcEdged  ArcStyle = 4
)

const (
	Styled        Color = -2
	Brushed       Color = -3
	StyledBrushed Color = -4
	Tiled         Color = -5
	Transparent   Color = -6
	AntiAliased   Color = -7
)

func wrapImage(ptr C.gdImagePtr) (*Image, error) {
	if ptr == nil {
		return nil, ErrImageCreate
	}
	return withImageFinalizer(&Image{ptr: ptr}), nil
}

func wrapDecodedImage(ptr C.gdImagePtr) (*Image, error) {
	if ptr == nil {
		return nil, ErrDecode
	}
	return withImageFinalizer(&Image{ptr: ptr}), nil
}

func withImageFinalizer(im *Image) *Image {
	runtime.SetFinalizer(im, func(im *Image) { _ = im.Close() })
	return im
}

func (im *Image) cptr() (C.gdImagePtr, error) {
	if im == nil || im.ptr == nil {
		return nil, ErrClosedImage
	}
	return im.ptr, nil
}

// NewPalette creates a palette image with up to 256 colors.
func NewPalette(width, height int) (*Image, error) {
	if width <= 0 || height <= 0 {
		return nil, ErrInvalidArgument
	}
	return wrapImage(C.gdImageCreate(C.int(width), C.int(height)))
}

// NewTrueColor creates a truecolor image.
func NewTrueColor(width, height int) (*Image, error) {
	if width <= 0 || height <= 0 {
		return nil, ErrInvalidArgument
	}
	return wrapImage(C.gdImageCreateTrueColor(C.int(width), C.int(height)))
}

// Close releases the underlying gdImagePtr.
//
// Close is idempotent. After Close, methods that read state (Width, Height,
// TrueColor, ...) return zero values; mutating methods return ErrClosedImage.
func (im *Image) Close() error {
	if im == nil || im.ptr == nil {
		return nil
	}
	runtime.SetFinalizer(im, nil)
	C.gdImageDestroy(im.ptr)
	im.ptr = nil
	im.tile = nil
	im.brush = nil
	return nil
}

// Destroy is a v1 compatibility alias for Close.
//
// Deprecated: use Close.
func (im *Image) Destroy() {
	_ = im.Close()
}

// Clone returns a deep copy of im.
func (im *Image) Clone() (*Image, error) {
	ptr, err := im.cptr()
	if err != nil {
		return nil, err
	}
	return wrapImage(C.gdImageClone(ptr))
}

// Width returns the image width in pixels, or 0 if the image is closed.
func (im *Image) Width() int {
	ptr, err := im.cptr()
	if err != nil {
		return 0
	}
	return int(C.go_gd_image_sx(ptr))
}

// Height returns the image height in pixels, or 0 if the image is closed.
func (im *Image) Height() int {
	ptr, err := im.cptr()
	if err != nil {
		return 0
	}
	return int(C.go_gd_image_sy(ptr))
}

// Bounds returns a rectangle covering the full image. Returns the zero Rect
// when the image is closed.
func (im *Image) Bounds() Rect {
	return Rect{Width: im.Width(), Height: im.Height()}
}

// TrueColor reports whether the image is truecolor rather than palette-based.
// Returns false for a closed image.
func (im *Image) TrueColor() bool {
	ptr, err := im.cptr()
	if err != nil {
		return false
	}
	return C.go_gd_image_true_color(ptr) != 0
}

// ColorsTotal returns the number of palette colors, or 0 for a closed image.
func (im *Image) ColorsTotal() int {
	ptr, err := im.cptr()
	if err != nil {
		return 0
	}
	return int(C.go_gd_image_colors_total(ptr))
}

// TransparentColor returns the transparent palette index or -1 (also for a
// closed image).
func (im *Image) TransparentColor() Color {
	ptr, err := im.cptr()
	if err != nil {
		return Color(-1)
	}
	return Color(C.go_gd_image_transparent(ptr))
}

// Interlaced reports whether interlacing is enabled. Returns false for a
// closed image.
func (im *Image) Interlaced() bool {
	ptr, err := im.cptr()
	if err != nil {
		return false
	}
	return C.go_gd_image_interlaced(ptr) != 0
}

// Resolution returns horizontal and vertical DPI, or (0, 0) for a closed
// image.
func (im *Image) Resolution() (x, y int) {
	ptr, err := im.cptr()
	if err != nil {
		return 0, 0
	}
	return int(C.go_gd_image_resolution_x(ptr)), int(C.go_gd_image_resolution_y(ptr))
}

// SetResolution sets horizontal and vertical DPI.
func (im *Image) SetResolution(x, y uint) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	C.gdImageSetResolution(ptr, C.uint(x), C.uint(y))
	return nil
}
