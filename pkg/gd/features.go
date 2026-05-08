package gd

/*
#include <stdlib.h>
#include <gd.h>
*/
import "C"
import "unsafe"

type VersionInfo struct {
	Major   int
	Minor   int
	Release int
	Extra   string
	String  string
}

type Format string

const (
	FormatPNG  Format = "png"
	FormatJPEG Format = "jpg"
	FormatGIF  Format = "gif"
	FormatWebP Format = "webp"
	FormatWBMP Format = "wbmp"
	FormatBMP  Format = "bmp"
	FormatTGA  Format = "tga"
	FormatTIFF Format = "tiff"
	FormatGD   Format = "gd"
	FormatGD2  Format = "gd2"
	FormatHEIF Format = "heif"
	FormatAVIF Format = "avif"
	FormatXBM  Format = "xbm"
	FormatXPM  Format = "xpm"
)

func Version() VersionInfo {
	return VersionInfo{
		Major:   int(C.gdMajorVersion()),
		Minor:   int(C.gdMinorVersion()),
		Release: int(C.gdReleaseVersion()),
		Extra:   C.GoString(C.gdExtraVersion()),
		String:  C.GoString(C.gdVersionString()),
	}
}

func SupportsFormat(format Format, writing bool) bool {
	return SupportsFileType("x."+string(format), writing)
}

func SupportsFileType(filename string, writing bool) bool {
	cname := C.CString(filename)
	defer C.free(unsafe.Pointer(cname))
	return C.gdSupportsFileType(cname, boolInt(writing)) != 0
}

// Features reports which optional libgd codecs are available at runtime.
//
// FontConfig support is intentionally absent: probing it would require
// flipping libgd's global gdFTUseFontConfig flag and produce a process-wide
// side-effect just for a feature query. Callers that need FontConfig should
// invoke UseFontConfig explicitly.
type Features struct {
	PNG, JPEG, GIF, WebP, WBMP bool
	BMP, TGA, TIFF             bool
	GD, GD2                    bool
	HEIF, AVIF                 bool
	XBM, XPM                   bool
	FreeType                   bool
}

func RuntimeFeatures() Features {
	return Features{
		PNG:  SupportsFileType("x.png", false),
		JPEG: SupportsFileType("x.jpg", false),
		GIF:  SupportsFileType("x.gif", false),
		WebP: SupportsFileType("x.webp", false),
		WBMP: SupportsFileType("x.wbmp", false),
		BMP:  SupportsFileType("x.bmp", false),
		TGA:  SupportsFileType("x.tga", false),
		TIFF: SupportsFileType("x.tiff", false),
		GD:   SupportsFileType("x.gd", false),
		GD2:  SupportsFileType("x.gd2", false),
		HEIF: SupportsFileType("x.heif", false),
		AVIF: SupportsFileType("x.avif", false),
		XBM:  SupportsFileType("x.xbm", false),
		XPM:  SupportsFileType("x.xpm", false),
		// FreeType is a compile-time requirement of this package (text.go calls
		// gdImageStringFTEx), so it is always available when the package builds.
		FreeType: true,
	}
}
