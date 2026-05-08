package gd

/*
#include <stdlib.h>
#include <gd.h>

static gdPointPtr go_gd_points_from_coords(int *coords, int n) {
	gdPointPtr points = (gdPointPtr)calloc((size_t)n, sizeof(gdPoint));
	if (points == NULL) {
		return NULL;
	}
	for (int i = 0; i < n; i++) {
		points[i].x = coords[i * 2];
		points[i].y = coords[i * 2 + 1];
	}
	return points;
}

static int go_gd_image_polygon(gdImagePtr im, int *coords, int n, int color) {
	gdPointPtr points = go_gd_points_from_coords(coords, n);
	if (points == NULL) {
		return 0;
	}
	gdImagePolygon(im, points, n, color);
	free(points);
	return 1;
}

static int go_gd_image_open_polygon(gdImagePtr im, int *coords, int n, int color) {
	gdPointPtr points = go_gd_points_from_coords(coords, n);
	if (points == NULL) {
		return 0;
	}
	gdImageOpenPolygon(im, points, n, color);
	free(points);
	return 1;
}

static int go_gd_image_filled_polygon(gdImagePtr im, int *coords, int n, int color) {
	gdPointPtr points = go_gd_points_from_coords(coords, n);
	if (points == NULL) {
		return 0;
	}
	gdImageFilledPolygon(im, points, n, color);
	free(points);
	return 1;
}
*/
import "C"
import "unsafe"

func (im *Image) SetPixel(x, y int, color Color) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	C.gdImageSetPixel(ptr, C.int(x), C.int(y), C.int(color))
	return nil
}

func (im *Image) Pixel(x, y int) (Color, error) {
	ptr, err := im.cptr()
	if err != nil {
		return 0, err
	}
	return Color(C.gdImageGetPixel(ptr, C.int(x), C.int(y))), nil
}

func (im *Image) TrueColorPixel(x, y int) (Color, error) {
	ptr, err := im.cptr()
	if err != nil {
		return 0, err
	}
	return Color(C.gdImageGetTrueColorPixel(ptr, C.int(x), C.int(y))), nil
}

func (im *Image) BoundsSafe(x, y int) bool {
	ptr, err := im.cptr()
	if err != nil {
		return false
	}
	return C.gdImageBoundsSafe(ptr, C.int(x), C.int(y)) != 0
}

func (im *Image) AABlend() error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	C.gdImageAABlend(ptr)
	return nil
}

func (im *Image) Line(x1, y1, x2, y2 int, color Color) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	C.gdImageLine(ptr, C.int(x1), C.int(y1), C.int(x2), C.int(y2), C.int(color))
	return nil
}

func (im *Image) DashedLine(x1, y1, x2, y2 int, color Color) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	C.gdImageDashedLine(ptr, C.int(x1), C.int(y1), C.int(x2), C.int(y2), C.int(color))
	return nil
}

func (im *Image) Rectangle(x1, y1, x2, y2 int, color Color) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	C.gdImageRectangle(ptr, C.int(x1), C.int(y1), C.int(x2), C.int(y2), C.int(color))
	return nil
}

func (im *Image) FilledRectangle(x1, y1, x2, y2 int, color Color) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	C.gdImageFilledRectangle(ptr, C.int(x1), C.int(y1), C.int(x2), C.int(y2), C.int(color))
	return nil
}

func (im *Image) Arc(cx, cy, width, height, start, end int, color Color) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	C.gdImageArc(ptr, C.int(cx), C.int(cy), C.int(width), C.int(height), C.int(start), C.int(end), C.int(color))
	return nil
}

func (im *Image) FilledArc(cx, cy, width, height, start, end int, color Color, style ArcStyle) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	C.gdImageFilledArc(ptr, C.int(cx), C.int(cy), C.int(width), C.int(height), C.int(start), C.int(end), C.int(color), C.int(style))
	return nil
}

func (im *Image) Ellipse(cx, cy, width, height int, color Color) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	C.gdImageEllipse(ptr, C.int(cx), C.int(cy), C.int(width), C.int(height), C.int(color))
	return nil
}

func (im *Image) FilledEllipse(cx, cy, width, height int, color Color) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	C.gdImageFilledEllipse(ptr, C.int(cx), C.int(cy), C.int(width), C.int(height), C.int(color))
	return nil
}

func (im *Image) Fill(x, y int, color Color) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	C.gdImageFill(ptr, C.int(x), C.int(y), C.int(color))
	return nil
}

func (im *Image) FillToBorder(x, y int, border, color Color) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	C.gdImageFillToBorder(ptr, C.int(x), C.int(y), C.int(border), C.int(color))
	return nil
}

func cPoints(points []Point) ([]C.int, error) {
	if len(points) == 0 {
		return nil, ErrInvalidArgument
	}
	out := make([]C.int, len(points)*2)
	for i, p := range points {
		out[i*2] = C.int(p.X)
		out[i*2+1] = C.int(p.Y)
	}
	return out, nil
}

func (im *Image) Polygon(points []Point, color Color) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	cpoints, err := cPoints(points)
	if err != nil {
		return err
	}
	if C.go_gd_image_polygon(ptr, (*C.int)(unsafe.Pointer(&cpoints[0])), C.int(len(points)), C.int(color)) == 0 {
		return ErrInvalidArgument
	}
	return nil
}

func (im *Image) OpenPolygon(points []Point, color Color) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	cpoints, err := cPoints(points)
	if err != nil {
		return err
	}
	if C.go_gd_image_open_polygon(ptr, (*C.int)(unsafe.Pointer(&cpoints[0])), C.int(len(points)), C.int(color)) == 0 {
		return ErrInvalidArgument
	}
	return nil
}

func (im *Image) FilledPolygon(points []Point, color Color) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	cpoints, err := cPoints(points)
	if err != nil {
		return err
	}
	if C.go_gd_image_filled_polygon(ptr, (*C.int)(unsafe.Pointer(&cpoints[0])), C.int(len(points)), C.int(color)) == 0 {
		return ErrInvalidArgument
	}
	return nil
}

func (im *Image) SetClip(rect Rect) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	if rect.Width <= 0 || rect.Height <= 0 {
		return ErrInvalidArgument
	}
	C.gdImageSetClip(ptr, C.int(rect.X), C.int(rect.Y), C.int(rect.X+rect.Width-1), C.int(rect.Y+rect.Height-1))
	return nil
}

func (im *Image) Clip() (Rect, error) {
	ptr, err := im.cptr()
	if err != nil {
		return Rect{}, err
	}
	var x1, y1, x2, y2 C.int
	C.gdImageGetClip(ptr, &x1, &y1, &x2, &y2)
	return Rect{X: int(x1), Y: int(y1), Width: int(x2 - x1 + 1), Height: int(y2 - y1 + 1)}, nil
}

func (im *Image) SaveAlpha(save bool) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	C.gdImageSaveAlpha(ptr, boolInt(save))
	return nil
}

func (im *Image) AlphaBlending(blend bool) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	C.gdImageAlphaBlending(ptr, boolInt(blend))
	return nil
}

func (im *Image) Interlace(enabled bool) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	C.gdImageInterlace(ptr, boolInt(enabled))
	return nil
}

func (im *Image) SetThickness(thickness int) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	C.gdImageSetThickness(ptr, C.int(thickness))
	return nil
}

func (im *Image) SetStyle(style []Color) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	if len(style) == 0 {
		return ErrInvalidArgument
	}
	cstyle := make([]C.int, len(style))
	for i, c := range style {
		cstyle[i] = C.int(c)
	}
	C.gdImageSetStyle(ptr, (*C.int)(unsafe.Pointer(&cstyle[0])), C.int(len(cstyle)))
	return nil
}

func (im *Image) SetAntiAliased(color Color) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	C.gdImageSetAntiAliased(ptr, C.int(color))
	return nil
}

func (im *Image) SetAntiAliasedDontBlend(color Color, dontBlend bool) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	C.gdImageSetAntiAliasedDontBlend(ptr, C.int(color), boolInt(dontBlend))
	return nil
}

// SetTile registers tile as the source image for libgd's tile-fill operations.
//
// libgd does not copy tile internally, so the caller must keep it alive for
// the lifetime of any drawing call that uses it. This binding keeps a
// reference to tile from im to prevent GC from collecting tile while it is
// still in use; call ClearTile or Close im to release that reference.
func (im *Image) SetTile(tile *Image) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	tilePtr, err := tile.cptr()
	if err != nil {
		return err
	}
	C.gdImageSetTile(ptr, tilePtr)
	im.tile = tile
	return nil
}

// ClearTile detaches any tile previously registered via SetTile.
func (im *Image) ClearTile() error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	C.gdImageSetTile(ptr, nil)
	im.tile = nil
	return nil
}

// SetBrush registers brush as the source image for libgd's brush-stroke
// operations. See SetTile for the lifetime contract.
func (im *Image) SetBrush(brush *Image) error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	brushPtr, err := brush.cptr()
	if err != nil {
		return err
	}
	C.gdImageSetBrush(ptr, brushPtr)
	im.brush = brush
	return nil
}

// ClearBrush detaches any brush previously registered via SetBrush.
func (im *Image) ClearBrush() error {
	ptr, err := im.cptr()
	if err != nil {
		return err
	}
	C.gdImageSetBrush(ptr, nil)
	im.brush = nil
	return nil
}
