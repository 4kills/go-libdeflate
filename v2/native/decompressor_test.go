package native

import (
	"bytes"
	"compress/zlib"
	"testing"
)

func compressWithStdLib(toCompress []byte) []byte {
	buf := &bytes.Buffer{}
	w := zlib.NewWriter(buf)
	w.Write([]byte(toCompress))
	w.Close()
	compressed := buf.Bytes()
	return compressed
}

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

func TestDecompressNewDecompressorWithExtendedDecompression(t *testing.T) {
	in := compressWithStdLib(shortString)

	out := make([]byte, len(shortString))
	dc, _ := NewDecompressorWithExtendedDecompression(3)
	defer dc.Close()
	if c, _, err := dc.Decompress(in, out, DecompressZlib); err != nil || c != len(in) {
		t.Error(err)
	}
	slicesEqual(shortString, out, t)
}

func TestDecompress(t *testing.T) {
	in := compressWithStdLib(shortString)

	// decompress with this lib

	out := make([]byte, len(shortString))
	dc, _ := NewDecompressor()
	defer dc.Close()
	if c, _, err := dc.Decompress(in, out, DecompressZlib); err != nil || c != len(in) {
		t.Error(err)
	}
	slicesEqual(shortString, out, t)

	c, out, err := dc.Decompress(in, nil, DecompressZlib)
	if err != nil || c != len(in) {
		t.Error(err)
	}
	slicesEqual(shortString, out, t)
}

func TestDecompressOversizedInput(t *testing.T) {
	in := compressWithStdLib(shortString)

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
	if err != nil || c != len(in) {
		t.Error(err)
	}
	slicesEqual(shortString, out, t)
}
