package libdeflate

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"encoding/hex"
	"strings"
	"testing"
)

/*---------------------
		UNIT TESTS
-----------------------*/

func TestDecompressDEFLATE(t *testing.T) {
	// compress with go standard lib
	buf := &bytes.Buffer{}
	w, _ := flate.NewWriter(buf, flate.DefaultCompression)
	w.Write(shortString)
	w.Close()
	in := buf.Bytes()

	// decompress with this lib

	out := make([]byte, len(shortString))
	dc, _ := NewDecompressor()
	defer dc.Close()
	if c, _, err := dc.Decompress(in, out, ModeDEFLATE); err != nil || c != len(in){
		t.Error(err)
	}
	slicesEqual(shortString, out, t)

	c, out, err := dc.Decompress(in, nil, ModeDEFLATE)
	if err != nil || c != len(in) {
		t.Error(err)
	}
	slicesEqual(shortString, out, t)
}

func TestDecompressGzip(t *testing.T) {
	// compress with go standard lib
	buf := &bytes.Buffer{}
	w := gzip.NewWriter(buf)
	w.Write(shortString)
	w.Close()
	in := buf.Bytes()

	// decompress with this lib

	out := make([]byte, len(shortString))
	dc, _ := NewDecompressor()
	defer dc.Close()
	if c, _, err := dc.Decompress(in, out, ModeGzip); err != nil || c != len(in){
		t.Error(err)
	}
	slicesEqual(shortString, out, t)

	c, out, err := dc.Decompress(in, nil, ModeGzip)
	if err != nil || c != len(in) {
		t.Error(err)
	}
	slicesEqual(shortString, out, t)
}

func TestDecompressZlib(t *testing.T) {
	// compress with go standard lib
	buf := &bytes.Buffer{}
	w := zlib.NewWriter(buf)
	w.Write(shortString)
	w.Close()
	in := buf.Bytes()

	// decompress with this lib

	out := make([]byte, len(shortString))
	dc, _ := NewDecompressor()
	defer dc.Close()
	if c, _, err := dc.DecompressZlib(in, out); err != nil || c != len(in) {
		t.Error(err)
	}
	slicesEqual(shortString, out, t)

	c, out, err := dc.DecompressZlib(in, nil)
	if err != nil || c != len(in) {
		t.Error(err)
	}
	slicesEqual(shortString, out, t)
}

var (
	deadlyCode = "789c7d90316f83301085f7fc0ac45cac3b1b63cc0685542c55a4264b1784825bd1024686284851fe7b8104c454c" +
		"b8bdff7de9def6e3b6b3cf6db298dedc0b26f02384014858efb2af78ecb1573c224e24e92301e4919d158c0dd7e7" +
		"9a40a5de765f39ed76aca2e6a364caf5a17cacdbe2a7dcd8c6ab5e9d790eafab49d2ce87ac4f308724a5c5c70553" +
		"6bfe1b93f968fa214287570bcfc882c6034a0fc736bfde873f38f99aee666fccfb3ad0f043d4604122917dc1abdc" +
		"c6178deaa81a8416d60afcfba9ae8293e2cf2a5536609753f04194881be106c1de6d2eec7058c9c3180cdfc87691" +
		"f814505305c77796d9e66042e372dd2622ec001c143ff09bec7bea151f90c05f3e916cc22f57cc976f73fc74374"
	comp, _ = hex.DecodeString(deadlyCode)
)

func TestDecompressWithDeadlyCode(t *testing.T) {
	dc, _ := NewDecompressor()
	defer dc.Close()
	_, _, err := dc.Decompress(comp, nil, ModeZlib)
	if err == nil || !strings.Contains(err.Error(), "maximum decompression factor") {
		t.Fail()
		return
	}
}
