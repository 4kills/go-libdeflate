package libdeflate

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"io"
	"testing"
)

var shortString = []byte("hello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\n")

/*---------------------
		UNIT TESTS
-----------------------*/

func TestNewCompressor(t *testing.T) {
	c, err := NewCompressor()
	if err != nil {
		t.Error(err)
	}
	defer c.Close()

	c, err = NewCompressorLevel(30)
	if err == nil {
		t.Fail()
	}
}

func TestCompressDEFLATE(t *testing.T) {
	c, _ := NewCompressor()
	defer c.Close()

	_, comp, err := c.Compress(shortString, nil, ModeDEFLATE)
	if err != nil {
		t.Error(err)
	}

	b := bytes.NewBuffer(comp)
	r := flate.NewReader(b)
	defer r.Close()

	dc := &bytes.Buffer{}
	io.Copy(dc, r)

	slicesEqual(shortString, dc.Bytes(), t)
}

func TestCompressGzip(t *testing.T) {
	c, _ := NewCompressor()
	defer c.Close()

	_, comp, err := c.Compress(shortString, nil, ModeGzip)
	if err != nil {
		t.Error(err)
	}

	b := bytes.NewBuffer(comp)
	r, err := gzip.NewReader(b)
	if err != nil {
		t.Error(err)
	}
	defer r.Close()

	dc := &bytes.Buffer{}
	io.Copy(dc, r)

	slicesEqual(shortString, dc.Bytes(), t)
}

func TestCompressZlibMaxComp(t *testing.T) {
	c, _ := NewCompressorLevel(MaxStdZlibCompressionLevel)
	defer c.Close()
	_, comp, err := c.CompressZlib(shortString, nil)
	if err != nil {
		t.Error(err)
	}

	r, _ := zlib.NewReader(bytes.NewBuffer(comp))
	defer r.Close()
	decomp := make([]byte, len(shortString))
	r.Read(decomp)

	slicesEqual(shortString, decomp, t)
}

func TestCompressZlib(t *testing.T) {
	c, _ := NewCompressor()
	defer c.Close()
	_, comp, err := c.CompressZlib(shortString, nil)
	if err != nil {
		t.Error(err)
	}

	r, _ := zlib.NewReader(bytes.NewBuffer(comp))
	defer r.Close()
	decomp := make([]byte, len(shortString))
	r.Read(decomp)

	slicesEqual([]byte(shortString), decomp, t)
}

// this test doesn't really say as much as TestCompress
func TestCompressMeta(t *testing.T) {
	c, _ := NewCompressor()
	defer c.Close()

	if _, _, err := c.CompressZlib(make([]byte, 0), nil); err == nil {
		t.Error("expected error")
	}

	n, out, err := c.CompressZlib(shortString, nil)
	if err != nil || n == 0 || n >= len(shortString) || n != len(out) {
		t.Error(err)
		t.Error(n)
	}

	out2 := make([]byte, len(shortString))
	n, _, err = c.CompressZlib(shortString, out2)
	if err != nil || n == 0 {
		t.Error(err)
		t.Error(n)
	}

	slicesEqual(out, out2[:n], t) //tests if rep produces same results
}

/*---------------------
	INTEGRATION TESTS
-----------------------*/

func TestCompressDecompress(t *testing.T) {
	c, _ := NewCompressor()
	defer c.Close()
	_, comp, err := c.CompressZlib(shortString, nil)
	if err != nil {
		t.Error(err)
	}

	out := make([]byte, len(shortString))
	dc, _ := NewDecompressor()
	defer dc.Close()
	if c, _, err := dc.DecompressZlib(comp, out); err != nil || c != len(comp){
		t.Error(err)
	}
	slicesEqual(shortString, out, t)
}

/*---------------------
		HELPER
-----------------------*/

func slicesEqual(expected, actual []byte, t *testing.T) {
	if len(expected) != len(actual) {
		t.Error("len of slices unequal")
		t.FailNow()
	}

	for i := range expected {
		if expected[i] != actual[i] {
			t.Errorf("slices differ at %d: want: %d, got: %d", i, expected[i], actual[i])
			t.FailNow()
		}
	}
}
