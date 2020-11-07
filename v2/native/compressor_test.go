package native

import (
	"bytes"
	"compress/zlib"
	"testing"
)

var shortString = []byte("hello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\nhello, world\n")

/*---------------------
		UNIT TESTS
-----------------------*/

func TestNewCompressor(t *testing.T) {
	c, err := NewCompressor(DefaultCompressionLevel)
	if err != nil {
		t.Error(err)
	}
	defer c.Close()

	c, err = NewCompressor(30)
	if err == nil {
		t.Fail()
	}
}

func TestCompressMaxComp(t *testing.T) {
	c, _ := NewCompressor(MaxStdZlibCompressionLevel)
	defer c.Close()
	_, comp, err := c.Compress(shortString, nil, CompressZlib)
	if err != nil {
		t.Error(err)
	}

	r, _ := zlib.NewReader(bytes.NewBuffer(comp))
	defer r.Close()
	decomp := make([]byte, len(shortString))
	r.Read(decomp)

	slicesEqual([]byte(shortString), decomp, t)
}

func TestCompress(t *testing.T) {
	c, _ := NewCompressor(DefaultCompressionLevel)
	defer c.Close()
	_, comp, err := c.Compress(shortString, nil, CompressZlib)
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
	c, _ := NewCompressor(DefaultCompressionLevel)
	defer c.Close()

	if _, _, err := c.Compress(make([]byte, 0), nil, CompressZlib); err == nil {
		t.Error("expected error")
	}

	n, out, err := c.Compress(shortString, nil, CompressZlib)
	if err != nil || n == 0 || n >= len(shortString) || n != len(out) {
		t.Error(err)
		t.Error(n)
	}

	out2 := make([]byte, len(shortString))
	n, _, err = c.Compress(shortString, out2, CompressZlib)
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
	c, _ := NewCompressor(DefaultCompressionLevel)
	defer c.Close()
	_, comp, err := c.Compress(shortString, nil, CompressZlib)
	if err != nil {
		t.Error(err)
	}

	out := make([]byte, len(shortString))
	dc, _ := NewDecompressor()
	defer dc.Close()
	if c, _, err := dc.Decompress(comp, out, DecompressZlib); err != nil || c != len(comp) {
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
