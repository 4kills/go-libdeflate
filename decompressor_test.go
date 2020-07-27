package libdeflate

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"testing"
)

/*---------------------
		UNIT TESTS
-----------------------*/

func TestDecompressDEFLATE(t *testing.T) {
	// compress with go standard lib
	buf := &bytes.Buffer{}
	w, _ := flate.NewWriter(buf, flate.DefaultCompression)
	w.Write([]byte(shortString))
	w.Close()
	in := buf.Bytes()

	// decompress with this lib

	out := make([]byte, len(shortString))
	dc, _ := NewDecompressor()
	defer dc.Close()
	if _, err := dc.Decompress(in, out, ModeDEFLATE); err != nil {
		t.Error(err)
	}
	slicesEqual([]byte(shortString), out, t)

	out, err := dc.Decompress(in, nil, ModeDEFLATE)
	if err != nil {
		t.Error(err)
	}
	slicesEqual([]byte(shortString), out, t)
}

func TestDecompressGzip(t *testing.T) {
	// compress with go standard lib
	buf := &bytes.Buffer{}
	w := gzip.NewWriter(buf)
	w.Write([]byte(shortString))
	w.Close()
	in := buf.Bytes()

	// decompress with this lib

	out := make([]byte, len(shortString))
	dc, _ := NewDecompressor()
	defer dc.Close()
	if _, err := dc.Decompress(in, out, ModeGzip); err != nil {
		t.Error(err)
	}
	slicesEqual([]byte(shortString), out, t)

	out, err := dc.Decompress(in, nil, ModeGzip)
	if err != nil {
		t.Error(err)
	}
	slicesEqual([]byte(shortString), out, t)
}

func TestDecompressZlib(t *testing.T) {
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
	if _, err := dc.DecompressZlib(in, out); err != nil {
		t.Error(err)
	}
	slicesEqual([]byte(shortString), out, t)

	out, err := dc.DecompressZlib(in, nil)
	if err != nil {
		t.Error(err)
	}
	slicesEqual([]byte(shortString), out, t)
}
