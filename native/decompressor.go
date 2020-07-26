package native

/*
#include "libdeflate.h"

typedef struct libdeflate_decompressor decomp;

int isNull(decomp* c) {
	if(!c) {
		return 1;
	}
	return 0;
}
*/
import "C"

// Decompressor decompresses any DEFLATE, zlib or gzip compressed data at any level
type Decompressor struct {
	dc *C.decomp
}

// NewDecompressor returns a new Decompressor or and error if out of memory
func NewDecompressor() (*Decompressor, error) {
	dc := C.libdeflate_alloc_decompressor()
	if C.isNull(dc) == 1 {
		return nil, errorOutOfMemory
	}

	return &Decompressor{dc}
}
