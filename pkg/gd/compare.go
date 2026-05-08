package gd

/*
#include <gd.h>
*/
import "C"

type CompareFlag int

const (
	CompareImage       CompareFlag = 1
	CompareNumColors   CompareFlag = 2
	CompareColor       CompareFlag = 4
	CompareSizeX       CompareFlag = 8
	CompareSizeY       CompareFlag = 16
	CompareTransparent CompareFlag = 32
	CompareBackground  CompareFlag = 64
	CompareInterlace   CompareFlag = 128
	CompareTrueColor   CompareFlag = 256
)

// Has reports whether all bits in flag are present. Has(0) is always false so
// that callers can safely test individual constants without special-casing
// the zero value.
func (f CompareFlag) Has(flag CompareFlag) bool {
	if flag == 0 {
		return false
	}
	return f&flag == flag
}

// Compare compares two images and returns a bitmask of CompareFlag values.
func Compare(a, b *Image) (CompareFlag, error) {
	aptr, err := a.cptr()
	if err != nil {
		return 0, err
	}
	bptr, err := b.cptr()
	if err != nil {
		return 0, err
	}
	return CompareFlag(C.gdImageCompare(aptr, bptr)), nil
}
