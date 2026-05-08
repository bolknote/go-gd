package gd

/*
#cgo darwin CFLAGS: -I/opt/homebrew/include -I/usr/local/include
#cgo darwin LDFLAGS: -L/opt/homebrew/lib -L/usr/local/lib -lgd
#cgo linux LDFLAGS: -lgd
#include <gd.h>
*/
import "C"

const maxPaletteColors = C.gdMaxColors

// MaxPaletteColors is the maximum number of colors in a libgd palette image.
const MaxPaletteColors = maxPaletteColors
