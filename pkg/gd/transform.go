package gd

/*
#include <gd.h>

static gdRect go_gd_rect(int x, int y, int w, int h) {
	gdRect r;
	r.x = x;
	r.y = y;
	r.width = w;
	r.height = h;
	return r;
}

static int go_gd_transform_affine_get_image(gdImagePtr *dst, gdImagePtr src, int x, int y, int w, int h, double *affine) {
	gdRect rect = go_gd_rect(x, y, w, h);
	return gdTransformAffineGetImage(dst, src, &rect, affine);
}

static int go_gd_transform_affine_copy(gdImagePtr dst, int dst_x, int dst_y, gdImagePtr src, int x, int y, int w, int h, double *affine) {
	gdRect rect = go_gd_rect(x, y, w, h);
	return gdTransformAffineCopy(dst, dst_x, dst_y, src, &rect, affine);
}

static int go_gd_affine_apply_to_point(double *out_x, double *out_y, double in_x, double in_y, double *affine) {
	gdPointF src;
	gdPointF dst;
	src.x = in_x;
	src.y = in_y;
	int ok = gdAffineApplyToPointF(&dst, &src, affine);
	*out_x = dst.x;
	*out_y = dst.y;
	return ok;
}

static int go_gd_transform_affine_bounding_box(int x, int y, int w, int h, double *affine, int *out) {
	gdRect src = go_gd_rect(x, y, w, h);
	gdRect bbox;
	int ok = gdTransformAffineBoundingBox(&src, affine, &bbox);
	out[0] = bbox.x;
	out[1] = bbox.y;
	out[2] = bbox.width;
	out[3] = bbox.height;
	return ok;
}
*/
import "C"
import "unsafe"

type InterpolationMethod int

const (
	InterpolationDefault InterpolationMethod = iota
	InterpolationBell
	InterpolationBessel
	InterpolationBilinearFixed
	InterpolationBicubic
	InterpolationBicubicFixed
	InterpolationBlackman
	InterpolationBox
	InterpolationBSpline
	InterpolationCatmullRom
	InterpolationGaussian
	InterpolationGeneralizedCubic
	InterpolationHermite
	InterpolationHamming
	InterpolationHanning
	InterpolationMitchell
	InterpolationNearestNeighbour
	InterpolationPower
	InterpolationQuadratic
	InterpolationSinc
	InterpolationTriangle
	InterpolationWeighted4
	InterpolationLinear
	InterpolationLanczos3
	InterpolationLanczos8
	InterpolationBlackmanBessel
	InterpolationBlackmanSinc
	InterpolationQuadraticBSpline
	InterpolationCubicSpline
	InterpolationCosine
	InterpolationWelsh
)

type Affine [6]float64

type PointF struct {
	X, Y float64
}

func (im *Image) FlipHorizontal() error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	C.gdImageFlipHorizontal(ptr)
	return nil
}

func (im *Image) FlipVertical() error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	C.gdImageFlipVertical(ptr)
	return nil
}

func (im *Image) FlipBoth() error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	C.gdImageFlipBoth(ptr)
	return nil
}

func (im *Image) Scale(width, height uint) (*Image, error) {
	ptr, err := im.cptr()
	if err != nil {
		return nil, err
	}
	out := C.gdImageScale(ptr, C.uint(width), C.uint(height))
	if out == nil {
		return nil, ErrTransform
	}
	return wrapImage(out)
}

func (im *Image) RotateInterpolated(angle float32, background Color) (*Image, error) {
	ptr, err := im.cptr()
	if err != nil {
		return nil, err
	}
	out := C.gdImageRotateInterpolated(ptr, C.float(angle), C.int(background))
	if out == nil {
		return nil, ErrTransform
	}
	return wrapImage(out)
}

func (im *Image) SetInterpolationMethod(method InterpolationMethod) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	if C.gdImageSetInterpolationMethod(ptr, C.gdInterpolationMethod(method)) == 0 {
		return ErrUnsupportedFeature
	}
	return nil
}

func (im *Image) InterpolationMethod() (InterpolationMethod, error) {
	ptr, err := im.cptr()
	if err != nil {
		return 0, err
	}
	return InterpolationMethod(C.gdImageGetInterpolationMethod(ptr)), nil
}

// AffineIdentity returns the identity affine transform.
func AffineIdentity() Affine {
	var m [6]C.double
	C.gdAffineIdentity((*C.double)(unsafe.Pointer(&m[0])))
	return affineFromC(m)
}

// AffineScale returns an affine transform that scales by (scaleX, scaleY).
func AffineScale(scaleX, scaleY float64) Affine {
	var m [6]C.double
	C.gdAffineScale((*C.double)(unsafe.Pointer(&m[0])), C.double(scaleX), C.double(scaleY))
	return affineFromC(m)
}

// AffineRotate returns an affine transform that rotates by angle (degrees).
func AffineRotate(angle float64) Affine {
	var m [6]C.double
	C.gdAffineRotate((*C.double)(unsafe.Pointer(&m[0])), C.double(angle))
	return affineFromC(m)
}

// AffineTranslate returns an affine transform that translates by (offsetX, offsetY).
func AffineTranslate(offsetX, offsetY float64) Affine {
	var m [6]C.double
	C.gdAffineTranslate((*C.double)(unsafe.Pointer(&m[0])), C.double(offsetX), C.double(offsetY))
	return affineFromC(m)
}

// ApplyToPoint applies the affine transform to src.
func (a Affine) ApplyToPoint(src PointF) PointF {
	matrix := a.c()
	var x, y C.double
	C.go_gd_affine_apply_to_point(&x, &y, C.double(src.X), C.double(src.Y), (*C.double)(unsafe.Pointer(&matrix[0])))
	return PointF{X: float64(x), Y: float64(y)}
}

// Flip returns a flipped copy of the affine transform.
func (a Affine) Flip(horizontal, vertical bool) Affine {
	src := a.c()
	var dst [6]C.double
	C.gdAffineFlip((*C.double)(unsafe.Pointer(&dst[0])), (*C.double)(unsafe.Pointer(&src[0])), boolInt(horizontal), boolInt(vertical))
	return affineFromC(dst)
}

// ShearHorizontal composes the receiver with a horizontal shear by angle (degrees).
// Returns ErrInvalidArgument if libgd rejects the shear angle (e.g. ±90°).
func (a Affine) ShearHorizontal(angle float64) (Affine, error) {
	shear, err := affineShearHorizontal(angle)
	if err != nil {
		return Affine{}, err
	}
	return a.Concat(shear), nil
}

// ShearVertical composes the receiver with a vertical shear by angle (degrees).
// Returns ErrInvalidArgument if libgd rejects the shear angle (e.g. ±90°).
func (a Affine) ShearVertical(angle float64) (Affine, error) {
	shear, err := affineShearVertical(angle)
	if err != nil {
		return Affine{}, err
	}
	return a.Concat(shear), nil
}

// Invert returns the inverse of the affine transform. Returns
// ErrInvalidArgument if the matrix is not invertible.
func (a Affine) Invert() (Affine, error) {
	src := a.c()
	var dst [6]C.double
	if C.gdAffineInvert((*C.double)(unsafe.Pointer(&dst[0])), (*C.double)(unsafe.Pointer(&src[0]))) == 0 {
		return Affine{}, ErrInvalidArgument
	}
	return affineFromC(dst), nil
}

// Concat returns the composition a ∘ other.
func (a Affine) Concat(other Affine) Affine {
	m1, m2 := a.c(), other.c()
	var dst [6]C.double
	C.gdAffineConcat((*C.double)(unsafe.Pointer(&dst[0])), (*C.double)(unsafe.Pointer(&m1[0])), (*C.double)(unsafe.Pointer(&m2[0])))
	return affineFromC(dst)
}

func (a Affine) Expansion() float64 {
	src := a.c()
	return float64(C.gdAffineExpansion((*C.double)(unsafe.Pointer(&src[0]))))
}

func (a Affine) Rectilinear() bool {
	src := a.c()
	return C.gdAffineRectilinear((*C.double)(unsafe.Pointer(&src[0]))) != 0
}

func (a Affine) Equal(other Affine) bool {
	m1, m2 := a.c(), other.c()
	return C.gdAffineEqual((*C.double)(unsafe.Pointer(&m1[0])), (*C.double)(unsafe.Pointer(&m2[0]))) != 0
}

// TransformAffine returns a new image with srcArea transformed by affine.
func (im *Image) TransformAffine(srcArea Rect, affine Affine) (*Image, error) {
	src, err := im.cptr()
	if err != nil {
		return nil, err
	}
	matrix := affine.c()
	var dst C.gdImagePtr
	if C.go_gd_transform_affine_get_image(&dst, src, C.int(srcArea.X), C.int(srcArea.Y), C.int(srcArea.Width), C.int(srcArea.Height), (*C.double)(unsafe.Pointer(&matrix[0]))) == 0 {
		return nil, ErrTransform
	}
	if dst == nil {
		return nil, ErrTransform
	}
	return wrapImage(dst)
}

// TransformAffineCopy applies affine to srcArea of src and copies the result
// into dst at (dstX, dstY).
func (dst *Image) TransformAffineCopy(src *Image, dstX, dstY int, srcArea Rect, affine Affine) error {
	dstPtr, err := dst.cptr()
	if err != nil {
		return err
	}
	srcPtr, err := src.cptr()
	if err != nil {
		return err
	}
	matrix := affine.c()
	if C.go_gd_transform_affine_copy(dstPtr, C.int(dstX), C.int(dstY), srcPtr, C.int(srcArea.X), C.int(srcArea.Y), C.int(srcArea.Width), C.int(srcArea.Height), (*C.double)(unsafe.Pointer(&matrix[0]))) == 0 {
		return ErrTransform
	}
	return nil
}

// TransformBoundingBox returns the bounding rect that contains src after
// applying the affine transform.
func (a Affine) TransformBoundingBox(src Rect) Rect {
	matrix := a.c()
	var out [4]C.int
	C.go_gd_transform_affine_bounding_box(C.int(src.X), C.int(src.Y), C.int(src.Width), C.int(src.Height), (*C.double)(unsafe.Pointer(&matrix[0])), (*C.int)(unsafe.Pointer(&out[0])))
	return Rect{X: int(out[0]), Y: int(out[1]), Width: int(out[2]), Height: int(out[3])}
}

func (a Affine) c() [6]C.double {
	var out [6]C.double
	for i, v := range a {
		out[i] = C.double(v)
	}
	return out
}

func affineFromC(in [6]C.double) Affine {
	var out Affine
	for i, v := range in {
		out[i] = float64(v)
	}
	return out
}

func affineShearHorizontal(angle float64) (Affine, error) {
	var dst [6]C.double
	if C.gdAffineShearHorizontal((*C.double)(unsafe.Pointer(&dst[0])), C.double(angle)) == 0 {
		return Affine{}, ErrInvalidArgument
	}
	return affineFromC(dst), nil
}

func affineShearVertical(angle float64) (Affine, error) {
	var dst [6]C.double
	if C.gdAffineShearVertical((*C.double)(unsafe.Pointer(&dst[0])), C.double(angle)) == 0 {
		return Affine{}, ErrInvalidArgument
	}
	return affineFromC(dst), nil
}
