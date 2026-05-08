// Package gd provides Go bindings for libgd.
//
// The package targets libgd 2.3.x and uses cgo. Advanced formats such as
// AVIF, HEIF, TIFF, WebP, FreeType, RAQM, and libimagequant depend on how the
// linked libgd library was built. Use Version and SupportsFileType to inspect
// the runtime library before enabling optional paths.
//
// Encode methods that write to a file return errors for Go-side validation and
// file open failures. Some libgd file encoders report write/codec failures only
// through libgd's global error callback.
//
// Image values are not safe for concurrent use by multiple goroutines without
// external synchronization.
package gd
