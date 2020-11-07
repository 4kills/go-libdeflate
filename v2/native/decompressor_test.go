package native

import (
	"bytes"
	"compress/zlib"
	"testing"
)

/*---------------------
		UNIT TESTS
-----------------------*/

func TestParseResult(t *testing.T) {
	if err := parseResult(1); err != errorBadData {
		t.Fail()
	}
	if err := parseResult(200); err != errorUnknown {
		t.Fail()
	}
	if err := parseResult(0); err != nil {
		t.Fail()
	}
}

func TestDecompress(t *testing.T) {
	// compress with go standard zlib
	buf := &bytes.Buffer{}
	w := zlib.NewWriter(buf)
	w.Write([]byte(shortString))
	w.Close()
	in := buf.Bytes()

	// decompress with this lib

	out := make([]byte, len(shortString))
	dc, _ := NewDecompressor()
	defer dc.Close()
	if c, _, err := dc.Decompress(in, out, DecompressZlib); err != nil || c != len(in) {
		t.Error(err)
	}
	slicesEqual(shortString, out, t)

	c, out, err := dc.Decompress(in, nil, DecompressZlib)
	if err != nil || c != len(in){
		t.Error(err)
	}
	slicesEqual(shortString, out, t)
}

func TestDecompressOversizedInput(t *testing.T) {
	// compress with go standard zlib
	buf := &bytes.Buffer{}
	w := zlib.NewWriter(buf)
	w.Write([]byte(shortString))
	w.Close()
	in := buf.Bytes()

	// decompress with this lib

	oversized := append(in, in...)
	out := make([]byte, len(shortString))
	dc, _ := NewDecompressor()
	defer dc.Close()
	if c, _, err := dc.Decompress(oversized, out, DecompressZlib); err != nil || c != len(in) {
		t.Error(err)
	}
	slicesEqual(shortString, out, t)

	c, out, err := dc.Decompress(oversized, nil, DecompressZlib)
	if err != nil || c != len(in){
		t.Error(err)
	}
	slicesEqual(shortString, out, t)
}
