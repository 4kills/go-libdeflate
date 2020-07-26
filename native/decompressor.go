package native

/*
#include "libdeflate.h"
#include "helper.h"

typedef struct libdeflate_decompressor decomp;
typedef enum libdeflate_result res;
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

// Decompress decompresses the given data from in to out and returns the number of bytes
// written to out, out and an error if something went wrong.
// If error != nil then the data in out is undefined.
// If you pass nil as out this function will allocate a sufficient buffer.
func (dc *Decompressor) Decompress(in, out []byte) ([]byte, error) {
	if len(in) == 0 {
		return out, errorNoInput
	}

	panic("implement")
	return nil, nil
}

func (dc *Decompressor) decompress(in, out []byte) ([]byte, error) {

	return nil, nil
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
