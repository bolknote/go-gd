package gd

import "testing"

// TestEnumParityWithLibGD asserts that the Go-side enum constants match the
// underlying libgd values exposed via libgd_consts.go. The Go binding
// enumerates these constants manually via iota; if libgd ever reorders an
// enum, this test will catch the divergence.
func TestEnumParityWithLibGD(t *testing.T) {
	t.Run("InterpolationMethod", func(t *testing.T) {
		cases := []struct {
			name string
			got  InterpolationMethod
			want int
		}{
			{"Default", InterpolationDefault, libgd.InterpolationDefault},
			{"Bell", InterpolationBell, libgd.InterpolationBell},
			{"Bessel", InterpolationBessel, libgd.InterpolationBessel},
			{"BilinearFixed", InterpolationBilinearFixed, libgd.InterpolationBilinearFixed},
			{"Bicubic", InterpolationBicubic, libgd.InterpolationBicubic},
			{"BicubicFixed", InterpolationBicubicFixed, libgd.InterpolationBicubicFixed},
			{"Blackman", InterpolationBlackman, libgd.InterpolationBlackman},
			{"Box", InterpolationBox, libgd.InterpolationBox},
			{"BSpline", InterpolationBSpline, libgd.InterpolationBSpline},
			{"CatmullRom", InterpolationCatmullRom, libgd.InterpolationCatmullRom},
			{"Gaussian", InterpolationGaussian, libgd.InterpolationGaussian},
			{"GeneralizedCubic", InterpolationGeneralizedCubic, libgd.InterpolationGeneralizedCubic},
			{"Hermite", InterpolationHermite, libgd.InterpolationHermite},
			{"Hamming", InterpolationHamming, libgd.InterpolationHamming},
			{"Hanning", InterpolationHanning, libgd.InterpolationHanning},
			{"Mitchell", InterpolationMitchell, libgd.InterpolationMitchell},
			{"NearestNeighbour", InterpolationNearestNeighbour, libgd.InterpolationNearestNeighbour},
			{"Power", InterpolationPower, libgd.InterpolationPower},
			{"Quadratic", InterpolationQuadratic, libgd.InterpolationQuadratic},
			{"Sinc", InterpolationSinc, libgd.InterpolationSinc},
			{"Triangle", InterpolationTriangle, libgd.InterpolationTriangle},
			{"Weighted4", InterpolationWeighted4, libgd.InterpolationWeighted4},
			{"Linear", InterpolationLinear, libgd.InterpolationLinear},
			{"Lanczos3", InterpolationLanczos3, libgd.InterpolationLanczos3},
			{"Lanczos8", InterpolationLanczos8, libgd.InterpolationLanczos8},
			{"BlackmanBessel", InterpolationBlackmanBessel, libgd.InterpolationBlackmanBessel},
			{"BlackmanSinc", InterpolationBlackmanSinc, libgd.InterpolationBlackmanSinc},
			{"QuadraticBSpline", InterpolationQuadraticBSpline, libgd.InterpolationQuadraticBSpline},
			{"CubicSpline", InterpolationCubicSpline, libgd.InterpolationCubicSpline},
			{"Cosine", InterpolationCosine, libgd.InterpolationCosine},
			{"Welsh", InterpolationWelsh, libgd.InterpolationWelsh},
		}
		for _, tc := range cases {
			if int(tc.got) != tc.want {
				t.Errorf("InterpolationMethod %s: Go=%d, libgd=%d", tc.name, tc.got, tc.want)
			}
		}
	})

	t.Run("ArcStyle", func(t *testing.T) {
		cases := []struct {
			name string
			got  ArcStyle
			want int
		}{
			{"Pie", ArcPie, libgd.ArcPie},
			{"Chord", ArcChord, libgd.ArcChord},
			{"NoFill", ArcNoFill, libgd.ArcNoFill},
			{"Edged", ArcEdged, libgd.ArcEdged},
		}
		for _, tc := range cases {
			if int(tc.got) != tc.want {
				t.Errorf("ArcStyle %s: Go=%d, libgd=%d", tc.name, tc.got, tc.want)
			}
		}
	})

	t.Run("CompareFlag", func(t *testing.T) {
		cases := []struct {
			name string
			got  CompareFlag
			want int
		}{
			{"Image", CompareImage, libgd.CompareImage},
			{"NumColors", CompareNumColors, libgd.CompareNumColors},
			{"Color", CompareColor, libgd.CompareColor},
			{"SizeX", CompareSizeX, libgd.CompareSizeX},
			{"SizeY", CompareSizeY, libgd.CompareSizeY},
			{"Transparent", CompareTransparent, libgd.CompareTransparent},
			{"Background", CompareBackground, libgd.CompareBackground},
			{"Interlace", CompareInterlace, libgd.CompareInterlace},
			{"TrueColor", CompareTrueColor, libgd.CompareTrueColor},
		}
		for _, tc := range cases {
			if int(tc.got) != tc.want {
				t.Errorf("CompareFlag %s: Go=%d, libgd=%d", tc.name, tc.got, tc.want)
			}
		}
	})

	t.Run("SpecialColors", func(t *testing.T) {
		cases := []struct {
			name string
			got  Color
			want int
		}{
			{"Styled", Styled, libgd.Styled},
			{"Brushed", Brushed, libgd.Brushed},
			{"StyledBrushed", StyledBrushed, libgd.StyledBrushed},
			{"Tiled", Tiled, libgd.Tiled},
			{"Transparent", Transparent, libgd.Transparent},
			{"AntiAliased", AntiAliased, libgd.AntiAliased},
		}
		for _, tc := range cases {
			if int(tc.got) != tc.want {
				t.Errorf("Color %s: Go=%d, libgd=%d", tc.name, tc.got, tc.want)
			}
		}
	})

	t.Run("MaxPaletteColors", func(t *testing.T) {
		if MaxPaletteColors != libgd.MaxPaletteColors {
			t.Errorf("MaxPaletteColors: Go=%d, libgd=%d", MaxPaletteColors, libgd.MaxPaletteColors)
		}
	})
}
