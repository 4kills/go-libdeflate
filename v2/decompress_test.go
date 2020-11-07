package libdeflate

import (
	"bytes"
	"compress/zlib"
	"testing"
)

/*---------------------
		UNIT TESTS
-----------------------*/
func TestDecompressZlibConvenience(t *testing.T) {
	// compress with go standard lib
	buf := &bytes.Buffer{}
	w := zlib.NewWriter(buf)
	w.Write([]byte(shortString))
	w.Close()
	in := buf.Bytes()

	// decompress with this lib

	out := make([]byte, len(shortString))
	dc, _ := NewDecompressor()
	defer dc.Close()
	if c, _, err := DecompressZlib(in, out); err != nil || c != len(in) {
		t.Error(err)
	}
	slicesEqual([]byte(shortString), out, t)

	c, out, err := dc.DecompressZlib(in, nil)
	if err != nil || c != len(in) {
		t.Error(err)
	}
	slicesEqual([]byte(shortString), out, t)
}

/*---------------------
		BENCHMARKS
-----------------------*/

func BenchmarkDecompressZlib(b *testing.B) {
	_, in, _ := CompressZlib(shortString, nil)
	out := make([]byte, len(shortString))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		DecompressZlib(in, out)
	}
}

func BenchmarkDecompressor_DecompressZlib(b *testing.B) {
	_, in, _ := CompressZlib(shortString, nil)
	out := make([]byte, len(shortString))

	dc, _ := NewDecompressor()
	defer dc.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		dc.DecompressZlib(in, out)
	}
}
