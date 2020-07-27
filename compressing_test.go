package libdeflate

import (
	"bytes"
	"compress/zlib"
	"io"
	"testing"
)

/*---------------------
		UNIT TESTS
-----------------------*/

func TestCompressZlibConvenience(t *testing.T) {
	_, comp, err := CompressZlib(shortString, nil)
	if err != nil {
		t.Error(err)
	}

	b := bytes.NewBuffer(comp)
	r, err := zlib.NewReader(b)
	if err != nil {
		t.Error(err)
	}
	defer r.Close()

	dc := &bytes.Buffer{}
	io.Copy(dc, r)

	slicesEqual(shortString, dc.Bytes(), t)
}

/*---------------------
		BENCHMARKS
-----------------------*/

func BenchmarkCompressZlib(b *testing.B) {
	out := make([]byte, len(shortString))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		CompressZlib(shortString, out)
	}
}

func BenchmarkCompressor_CompressZlib(b *testing.B) {
	c, _ := NewCompressor()
	defer c.Close()
	out := make([]byte, len(shortString))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		c.CompressZlib(shortString, out)
	}
}
