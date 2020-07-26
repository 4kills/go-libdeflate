package libdeflate

import "github.com/4kills/libdeflate/native"

// Compressor compresses data at the specified compression level.
//
// A single compressor must not not be used across multiple threads concurrently.
// If you want to compress concurrently, create a compressor for each thread.
type Compressor struct {
	c *native.Compressor
	lvl int
}

// NewCompressor returns a new Compressor used to compress data with compression level DefaultCompressionLevel.
// Errors if out of memory.
// See NewCompressorLevel for custom compression level
func NewCompressor() (Compressor, error) {
	return NewCompressorLevel(DefaultCompressionLevel)
}

// NewCompressor returns a new Compressor used to compress data.
// Errors if out of memory or if an invalid compression level was passed.
//
// The compression level is legal if and only if:
// MinCompressionLevel <= level <= MaxCompressionLevel
func NewCompressorLevel(level int) (Compressor, error) {
	c, err := native.NewCompressor(level)
	return Compressor{c, level}, err
}

// CompressZlib compresses the data from in to out (in zlib format) and returns the number
// of bytes written to out, out (sliced to written) or an error if the out buffer was too short.
// If you pass nil for out, this function will allocate a fitting buffer and return it (not preferred though). o
//
// See c.Compress for further information.
func (c Compressor) CompressZlib(in, out []byte) (int, []byte, error) {
	return c.Compress(in, out, ModeZlib)
}

// Compress compresses the data from in to out and returns the number
// of bytes written to out, out (sliced to written) or an error if the out buffer was too short.
// If you pass nil for out, this function will allocate a fitting buffer and return it (not preferred though).
//
// m specifies which compression format should be used (e.g. ModeZlib)
//
// Notice that for extremely small or already highly compressed data,
// the compressed data could be larger than uncompressed.
// If out == nil: For a too large discrepancy (len(out) > 1000 + 2 * len(in)) Compress will error
func (c Compressor) Compress(in, out []byte, m Mode) (int, []byte, error) {
	switch m {
	case ModeZlib: return c.c.Compress(in, out)
	default: panic("libdeflate: compress: invalid mode")
	}
}

// Level returns the compression level at which this Compressor compresses.
// May be called after having closed a Compressor.
func (c Compressor) Level() int {
	return c.lvl
}

// Close closes the compressor and releases all occupied resources.
// It is the users responsibility to close compressors in order to free resources,
// as the underlying c objects are not subject to the go garbage collector. They have to be freed manually.
//
// After closing, the compressor must not be used anymore, as the methods will panic (except for the c.Level() method).
func (c Compressor) Close() {
	c.c.Close()
}