package native

/*
#include "libdeflate.h"
#include "helper.h"
#include <stddef.h>
#include <stdlib.h>

typedef struct libdeflate_decompressor decomp;
typedef enum libdeflate_result res;

size_t* sizePtr(int fit) {
	if (fit == 1) return NULL;
	return (size_t*) malloc(sizeof(size_t));
}
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

// Decompress decompresses the given data from in to out and returns out and an error if something went wrong.
// If error != nil, then the data in out is undefined.
// If you pass a buffer to out, the size of this buffer must exactly match the length of the decompressed data.
// If you pass nil as out, this function will allocate a sufficient buffer and return it.
func (dc *Decompressor) Decompress(in, out []byte) ([]byte, error) {
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

	f := 0
	if fit {
		f = 1
	}
	s := C.sizePtr(C.int(f))

	r := C.libdeflate_zlib_decompress(dc.dc,
		unsafe.Pointer(inAddr), intToInt64(len(in)),
		unsafe.Pointer(outAddr), intToInt64(len(out)),
		s,
	)
	err := parseResult(C.res(r))

	n := len(out)
	if !fit {
		n = int(*s)
	}
	C.free(unsafe.Pointer(s))

	return n, err
}

func parseResult(r C.res) error {
	switch r {
	case C.LIBDEFLATE_SUCCESS:
		return nil
	case C.LIBDEFLATE_BAD_DATA:
		return errorBadData
	case C.LIBDEFLATE_SHORT_OUTPUT:
		return errorShortOutput
	case C.LIBDEFLATE_INSUFFICIENT_SPACE:
		return errorInsufficientSpace
	default:
		return errorUnknown
	}
}

// Close frees the memory allocated by C objects
func (dc *Decompressor) Close() {
	C.libdeflate_free_decompressor(dc.dc)
}
