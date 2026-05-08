package gd

/*
#include <gd.h>
*/
import "C"

// libgdConsts mirrors libgd's enum and #define values as Go ints so that the
// public Go enum constants can be cross-checked against the linked library
// at run time. Cgo is not allowed in *_test.go files, so the bridge lives
// here and the parity tests reference these fields.
type libgdConsts struct {
	InterpolationDefault          int
	InterpolationBell             int
	InterpolationBessel           int
	InterpolationBilinearFixed    int
	InterpolationBicubic          int
	InterpolationBicubicFixed     int
	InterpolationBlackman         int
	InterpolationBox              int
	InterpolationBSpline          int
	InterpolationCatmullRom       int
	InterpolationGaussian         int
	InterpolationGeneralizedCubic int
	InterpolationHermite          int
	InterpolationHamming          int
	InterpolationHanning          int
	InterpolationMitchell         int
	InterpolationNearestNeighbour int
	InterpolationPower            int
	InterpolationQuadratic        int
	InterpolationSinc             int
	InterpolationTriangle         int
	InterpolationWeighted4        int
	InterpolationLinear           int
	InterpolationLanczos3         int
	InterpolationLanczos8         int
	InterpolationBlackmanBessel   int
	InterpolationBlackmanSinc     int
	InterpolationQuadraticBSpline int
	InterpolationCubicSpline      int
	InterpolationCosine           int
	InterpolationWelsh            int

	ArcPie    int
	ArcChord  int
	ArcNoFill int
	ArcEdged  int

	CompareImage       int
	CompareNumColors   int
	CompareColor       int
	CompareSizeX       int
	CompareSizeY       int
	CompareTransparent int
	CompareBackground  int
	CompareInterlace   int
	CompareTrueColor   int

	Styled        int
	Brushed       int
	StyledBrushed int
	Tiled         int
	Transparent   int
	AntiAliased   int

	MaxPaletteColors int
}

var libgd = libgdConsts{
	InterpolationDefault:          int(C.GD_DEFAULT),
	InterpolationBell:             int(C.GD_BELL),
	InterpolationBessel:           int(C.GD_BESSEL),
	InterpolationBilinearFixed:    int(C.GD_BILINEAR_FIXED),
	InterpolationBicubic:          int(C.GD_BICUBIC),
	InterpolationBicubicFixed:     int(C.GD_BICUBIC_FIXED),
	InterpolationBlackman:         int(C.GD_BLACKMAN),
	InterpolationBox:              int(C.GD_BOX),
	InterpolationBSpline:          int(C.GD_BSPLINE),
	InterpolationCatmullRom:       int(C.GD_CATMULLROM),
	InterpolationGaussian:         int(C.GD_GAUSSIAN),
	InterpolationGeneralizedCubic: int(C.GD_GENERALIZED_CUBIC),
	InterpolationHermite:          int(C.GD_HERMITE),
	InterpolationHamming:          int(C.GD_HAMMING),
	InterpolationHanning:          int(C.GD_HANNING),
	InterpolationMitchell:         int(C.GD_MITCHELL),
	InterpolationNearestNeighbour: int(C.GD_NEAREST_NEIGHBOUR),
	InterpolationPower:            int(C.GD_POWER),
	InterpolationQuadratic:        int(C.GD_QUADRATIC),
	InterpolationSinc:             int(C.GD_SINC),
	InterpolationTriangle:         int(C.GD_TRIANGLE),
	InterpolationWeighted4:        int(C.GD_WEIGHTED4),
	InterpolationLinear:           int(C.GD_LINEAR),
	InterpolationLanczos3:         int(C.GD_LANCZOS3),
	InterpolationLanczos8:         int(C.GD_LANCZOS8),
	InterpolationBlackmanBessel:   int(C.GD_BLACKMAN_BESSEL),
	InterpolationBlackmanSinc:     int(C.GD_BLACKMAN_SINC),
	InterpolationQuadraticBSpline: int(C.GD_QUADRATIC_BSPLINE),
	InterpolationCubicSpline:      int(C.GD_CUBIC_SPLINE),
	InterpolationCosine:           int(C.GD_COSINE),
	InterpolationWelsh:            int(C.GD_WELSH),

	ArcPie:    int(C.gdPie),
	ArcChord:  int(C.gdChord),
	ArcNoFill: int(C.gdNoFill),
	ArcEdged:  int(C.gdEdged),

	CompareImage:       int(C.GD_CMP_IMAGE),
	CompareNumColors:   int(C.GD_CMP_NUM_COLORS),
	CompareColor:       int(C.GD_CMP_COLOR),
	CompareSizeX:       int(C.GD_CMP_SIZE_X),
	CompareSizeY:       int(C.GD_CMP_SIZE_Y),
	CompareTransparent: int(C.GD_CMP_TRANSPARENT),
	CompareBackground:  int(C.GD_CMP_BACKGROUND),
	CompareInterlace:   int(C.GD_CMP_INTERLACE),
	CompareTrueColor:   int(C.GD_CMP_TRUECOLOR),

	Styled:        int(C.gdStyled),
	Brushed:       int(C.gdBrushed),
	StyledBrushed: int(C.gdStyledBrushed),
	Tiled:         int(C.gdTiled),
	Transparent:   int(C.gdTransparent),
	AntiAliased:   int(C.gdAntiAliased),

	MaxPaletteColors: int(C.gdMaxColors),
}
