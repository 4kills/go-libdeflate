package native

/*
#include "libdeflate.h"
#include "helper.h"
#include <stddef.h>
#include <stdlib.h>
#include <stdint.h>

typedef struct libdeflate_decompressor decomp;
*/
import "C"
import "unsafe"

// Decompressor decompresses any DEFLATE, zlib or gzip compressed data at any level
type Decompressor struct {
	dc *C.decomp
	isClosed bool
}

// NewDecompressor returns a new Decompressor or and error if out of memory
func NewDecompressor() (*Decompressor, error) {
	dc := C.libdeflate_alloc_decompressor()
	if C.isNull(unsafe.Pointer(dc)) == 1 {
		return nil, errorOutOfMemory
	}

	return &Decompressor{dc, false}, nil
}

// Decompress decompresses the given data from in to out and returns out and an error if something went wrong.
// If error != nil, then the data in out is undefined.
// If you pass a buffer to out, the size of this buffer must exactly match the length of the decompressed data.
// If you pass nil as out, this function will allocate a sufficient buffer and return it.
func (dc *Decompressor) Decompress(in, out []byte) ([]byte, error) {
	if dc.isClosed {
		panic(errorAlreadyClosed)
	}
	if len(in) == 0 {
		return out, errorNoInput
	}

	if out != nil {
		_, err := dc.decompress(in, out, true)
		return out, err
	}

	n := 0
	inc := 6
	err := errorInsufficientSpace
	for err == errorInsufficientSpace {
		out = make([]byte, len(in)*inc)
		n, err = dc.decompress(in, out, false)
		if inc >= 16 {
			inc += 3
			continue
		}
		inc += 5
	}

	return out[:n], err
}

func (dc *Decompressor) decompress(in, out []byte, fit bool) (int, error) {
	inAddr := startMemAddr(in)
	outAddr := startMemAddr(out)

	var s int64
	sPtr := uintptr(unsafe.Pointer(&s))
	if fit {
		sPtr = 0
	}

	err := parseResult(C.res(C.libdeflate_zlib_decompress(dc.dc,
		unsafe.Pointer(inAddr), intToInt64(len(in)),
		unsafe.Pointer(outAddr), intToInt64(len(out)),
		C.mkPtr(C.size_t(sPtr)),
	)))

	n := len(out)
	if !fit {
		n = int(s)
	}

	return n, err
}

// Close frees the memory allocated by C objects
func (dc *Decompressor) Close() {
	if dc.isClosed {
		panic(errorAlreadyClosed)
	}
	C.libdeflate_free_decompressor(dc.dc)
	dc.isClosed = true
}
