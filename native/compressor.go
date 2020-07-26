package native

/*
#include "libdeflate.h"
#include "helper.h"

typedef struct libdeflate_compressor comp;
*/
import "C"
import "unsafe"

// Compressor compresses data to zlib format at the specified level
type Compressor struct {
	c   *C.comp
	lvl int
}

// NewCompressor returns a new Compressor used to compress data.
// Errors if out of memory or invalid lvl
func NewCompressor(lvl int) (*Compressor, error) {
	if lvl < minLevel || lvl > maxLevel {
		return nil, errorInvalidLevel
	}

	c := C.libdeflate_alloc_compressor(C.int(lvl))
	if C.isNull(unsafe.Pointer(c)) == 1 {
		return nil, errorOutOfMemory
	}

	return &Compressor{c, lvl}, nil
}

// Compress compresses the data from in to out and returns the number
// of bytes written to out, out and an error if the out buffer was too short.
// If you pass nil for out, this function will allocate a fitting buffer and return it.
func (c *Compressor) Compress(in, out []byte) (int, []byte, error) {
	if len(in) == 0 {
		return 0, out, errorNoInput
	}

	if out != nil {
		return c.compress(in, out)
	}

	out = make([]byte, len(in))
	n, out, err := c.compress(in, out)
	if err != nil {
		copy(out, in)
		return len(in), out[:len(in)], nil
	}

	return n, out[:n], nil
}

func (c *Compressor) compress(in, out []byte) (int, []byte, error) {
	inAddr := startMemAddress(in)
	outAddr := startMemAddress(out)

	written := int(C.libdeflate_zlib_compress(c.c,
		unsafe.Pointer(inAddr), intToInt64(len(in)),
		unsafe.Pointer(outAddr), intToInt64(len(out)),
	))

	if written == 0 {
		return written, out, errorShortBuffer
	}
	return written, out, nil
}

// Close frees the memory allocated by C objects
func (c *Compressor) Close() {
	C.libdeflate_free_compressor(c.c)
}
