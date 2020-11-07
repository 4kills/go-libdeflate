package libdeflate

import (
	"bytes"
	"compress/zlib"
	"testing"
)

/*---------------------
		BENCHMARKS
-----------------------*/

// real world data benchmarks

const compressedMcPacketsLoc = "https://raw.githubusercontent.com/4kills/zlib_benchmark/master/compressed_mc_packets.json"

var compressedMcPackets [][]byte

func BenchmarkDecompressZlibMcPacketsLibdeflate(b *testing.B) {
	loadPacketsIfNil(&compressedMcPackets, compressedMcPacketsLoc)
	loadPacketsIfNil(&decompressedMcPackets, decompressedMcPacketsLoc)
	dc, _ := NewDecompressor()
	defer dc.Close()

	b.ResetTimer()

	reportBytesPerIteration(compressedMcPackets, b)

	for i := 0; i < b.N; i++ {
		for j, v := range compressedMcPackets {
			b.StopTimer()
			decompressed := make([]byte, len(decompressedMcPackets[j]))
			b.StartTimer()

			dc.DecompressZlib(v, decompressed)
		}
	}
}

func BenchmarkDecompressZlibMcPacketsStdLib(b *testing.B) {
	loadPacketsIfNil(&compressedMcPackets, compressedMcPacketsLoc)
	loadPacketsIfNil(&decompressedMcPackets, decompressedMcPacketsLoc)

	buf := bytes.NewBuffer(compressedMcPackets[0]) // the std library needs this or else I can't create a reader
	r, _ := zlib.NewReader(buf)
	defer r.Close()

	b.ResetTimer()

	reportBytesPerIteration(compressedMcPackets, b)

	for i := 0; i < b.N; i++ {
		for j, v := range compressedMcPackets {
			b.StopTimer()
			res, _ := r.(zlib.Resetter)
			res.Reset(bytes.NewBuffer(v), nil) // to make the std reader work for new data
			decompressed := make([]byte, 0, len(decompressedMcPackets[j]))
			b.StartTimer()

			r.Read(decompressed)
		}
	}
}

func BenchmarkDecompressZlib1McPacketLibdeflate(b *testing.B) {
	loadPacketsIfNil(&compressedMcPackets, compressedMcPacketsLoc)
	loadPacketsIfNil(&decompressedMcPackets, decompressedMcPacketsLoc)
	dc, _ := NewDecompressor()
	defer dc.Close()
	decompressed := make([]byte, len(decompressedMcPackets[0]))

	b.ResetTimer()

	reportBytesPerIteration(compressedMcPackets[1:2], b)

	for i := 0; i < b.N; i++ {
		dc.DecompressZlib(compressedMcPackets[1], decompressed)
	}
}

func BenchmarkDecompressZlib1McPacketStdLib(b *testing.B) {
	loadPacketsIfNil(&compressedMcPackets, compressedMcPacketsLoc)
	loadPacketsIfNil(&decompressedMcPackets, decompressedMcPacketsLoc)

	buf := bytes.NewBuffer(compressedMcPackets[0]) // the std library needs this or else I can't create a reader
	r, _ := zlib.NewReader(buf)
	defer r.Close()

	b.ResetTimer()

	reportBytesPerIteration(compressedMcPackets[1:2], b)

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		res, _ := r.(zlib.Resetter)
		res.Reset(bytes.NewBuffer(compressedMcPackets[1]), nil) // to make the std reader work for new data
		decompressed := make([]byte, 0, len(decompressedMcPackets[1]))
		b.StartTimer()

		r.Read(decompressed)
	}
}
