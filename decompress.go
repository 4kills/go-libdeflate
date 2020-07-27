package libdeflate

// Decompress decompresses the given data from in (zlib formatted) to out and returns out
// or an error if something went wrong.
//
// If you pass a buffer to out, the size of this buffer must exactly match the length of the decompressed data.
// If you pass nil to out, this function will allocate a sufficient buffer and return it.
//
// IF YOU WANT TO DECOMPRESS MORE THAN ONCE, PLEASE REFER TO NewDecompressor(),
// as this function creates a new Decompressor which is then closed at the end of the function.
//
// If error != nil, the data in out is undefined.
func DecompressZlib(in, out []byte) ([]byte, error) {
	return Decompress(in, out, ModeZlib)
}

// Decompress decompresses the given data from in to out and returns out or an error if something went wrong.
// Mode m specifies the format (e.g. zlib) of the data within in.
//
// If you pass a buffer to out, the size of this buffer must exactly match the length of the decompressed data.
// If you pass nil to out, this function will allocate a sufficient buffer and return it.
//
// IF YOU WANT TO DECOMPRESS MORE THAN ONCE, PLEASE REFER TO NewDecompressor(),
// as this function creates a new Decompressor which is then closed at the end of the function.
//
// If error != nil, the data in out is undefined.
func Decompress(in, out []byte, m Mode) ([]byte, error) {
	dc, err := NewDecompressor()
	if err != nil {
		return out, err
	}
	defer dc.Close()

	return dc.Decompress(in, out, m)
}