package native

/*
#include "libdeflate.h"
#include "helper.h"

typedef struct libdeflate_decompressor decomp;
*/
import "C"
import "unsafe"

// Decompressor decompresses any DEFLATE, zlib or gzip compressed data at any level
type Decompressor struct {
	dc *C.decomp
}

// NewDecompressor returns a new Decompressor or and error if out of memory
func NewDecompressor() (*Decompressor, error) {
	dc := C.libdeflate_alloc_decompressor()
	if C.isNull(unsafe.Pointer(dc)) == 1 {
		return nil, errorOutOfMemory
	}

	return &Decompressor{dc}, nil
}
