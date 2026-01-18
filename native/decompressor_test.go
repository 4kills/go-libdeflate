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
	if err := parseResult(1); err != ErrorBadData {
		t.Fail()
	}
	if err := parseResult(200); err != ErrorUnknown {
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
	if _, err := dc.Decompress(in, out, DecompressZlib); err != nil {
		t.Error(err)
	}
	slicesEqual([]byte(shortString), out, t)

	out, err := dc.Decompress(in, nil, DecompressZlib)
	if err != nil {
		t.Error(err)
	}
	slicesEqual([]byte(shortString), out, t)
}
