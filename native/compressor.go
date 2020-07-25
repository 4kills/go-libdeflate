package native

/*
#include "libdeflate.h"

typedef struct libdeflate_compressor comp;

int isNull(comp* c) {
	if(!c) {
		return 1;
	}
	return 0;
}
*/
import "C"

type Compressor struct {
	c   *C.comp
	lvl int
}

func NewCompressor(lvl int) (*Compressor, error) {
	c := C.libdeflate_alloc_compressor(C.int(lvl))
	if C.isNull(c) == 1 {
		return nil, errorOutOfMemory
	}

	return &Compressor{c, lvl}, nil
}
