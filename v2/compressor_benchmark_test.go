package libdeflate

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

/*---------------------
		BENCHMARKS
-----------------------*/

// real world data benchmarks

const decompressedMcPacketsLoc = "https://raw.githubusercontent.com/4kills/zlib_benchmark/master/decompressed_mc_packets.json"

var decompressedMcPackets [][]byte

func BenchmarkCompressZlibAllMcPacketsLibdeflate(b *testing.B) {
	compressZlibAllMcPacketsLibdeflateLevel(DefaultCompressionLevel, b)
}

func BenchmarkCompressZlibAllMcPacketsStdLib(b *testing.B) {
	compressZlibAllMcPacketsStdLibLevel(DefaultCompressionLevel, b)
}

func BenchmarkCompressZlibAllMcPacketsFastestLibdeflate(b *testing.B) {
	compressZlibAllMcPacketsLibdeflateLevel(MinCompressionLevel, b)
}

func BenchmarkCompressZlibAllMcPacketsFastestStdLib(b *testing.B) {
	compressZlibAllMcPacketsStdLibLevel(MinCompressionLevel, b)
}

func BenchmarkCompressZlib1McPacketLibdeflate(b *testing.B) {
	compressZlib1McPacketLibdeflateLevel(DefaultCompressionLevel, b)
}

func BenchmarkCompressZlib1McPacketStdLib(b *testing.B) {
	compressZlib1McPacketStdLibLevel(DefaultCompressionLevel, b)
}

//-----------------------------------------------------------------

func compressZlibAllMcPacketsLibdeflateLevel(level int, b *testing.B) {
	loadPacketsIfNil(&decompressedMcPackets, decompressedMcPacketsLoc)
	c, _ := NewCompressorLevel(level)
	defer c.Close()

	b.ResetTimer()

	reportBytesPerIteration(decompressedMcPackets, b)

	for i := 0; i < b.N; i++ {
		for _, v := range decompressedMcPackets {
			b.StopTimer()
			out := make([]byte, len(v))
			b.StartTimer()

			n, _, _ := c.CompressZlib(v, out)

			b.ReportMetric(float64(len(v))/float64(n), "ratio")
		}
	}
}

func compressZlibAllMcPacketsStdLibLevel(level int, b *testing.B) {
	loadPacketsIfNil(&decompressedMcPackets, decompressedMcPacketsLoc)
	w, _ := zlib.NewWriterLevel(&bytes.Buffer{}, level)
	defer w.Close()

	b.ResetTimer()

	reportBytesPerIteration(decompressedMcPackets, b)

	for i := 0; i < b.N; i++ {
		for _, v := range decompressedMcPackets {
			b.StopTimer()
			buf := bytes.NewBuffer(make([]byte, 0, len(v)))
			w.Reset(buf)
			b.StartTimer()

			w.Write(v)
			w.Flush() // has to be called too, because this lib's compressor always flushes

			b.ReportMetric(float64(len(v))/float64(buf.Len()), "ratio")
		}
	}
}

func compressZlib1McPacketLibdeflateLevel(level int, b *testing.B) {
	loadPacketsIfNil(&decompressedMcPackets, decompressedMcPacketsLoc)
	c, _ := NewCompressorLevel(level)
	defer c.Close()
	out := make([]byte, len(decompressedMcPackets[1]))

	b.ResetTimer()

	reportBytesPerIteration(decompressedMcPackets[1:2], b)

	for i := 0; i < b.N; i++ {
		n, _, _ := c.CompressZlib(decompressedMcPackets[1], out)

		b.ReportMetric(float64(len(decompressedMcPackets[1]))/float64(n), "ratio")
	}
}

func compressZlib1McPacketStdLibLevel(level int, b *testing.B) {
	loadPacketsIfNil(&decompressedMcPackets, decompressedMcPacketsLoc)
	w, _ := zlib.NewWriterLevel(&bytes.Buffer{}, level)
	defer w.Close()
	out := make([]byte, 0, len(decompressedMcPackets[1]))

	b.ResetTimer()

	reportBytesPerIteration(decompressedMcPackets[1:2], b)

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		buf := bytes.NewBuffer(out)
		w.Reset(buf)
		b.StartTimer()

		w.Write(decompressedMcPackets[1])
		w.Flush() // has to be called too, because this lib's compressor always flushes

		b.ReportMetric(float64(len(decompressedMcPackets[1]))/float64(buf.Len()), "ratio")
	}
}

/*---------------------
		HELPER
-----------------------*/

func reportBytesPerIteration(input [][]byte, b *testing.B) {
	b.StopTimer()
	numOfBytes := 0
	for _, v := range input {
		numOfBytes += len(v)
	}
	b.ReportMetric(float64(numOfBytes), "bytes_processed/op")
	b.StartTimer()
}

func loadPacketsIfNil(packets *[][]byte, loc string) {
	if *packets != nil {
		return
	}
	*packets = loadPackets(loc)
}

func loadPackets(loc string) [][]byte {
	b, err := downloadFile(loc)
	if err != nil {
		panic(err)
	}

	return unmarshal(b)
}

func unmarshal(b *bytes.Buffer) [][]byte {
	var out [][]byte

	byteValue, _ := ioutil.ReadAll(b)
	json.Unmarshal(byteValue, &out)
	return out
}

func downloadFile(url string) (*bytes.Buffer, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	b := &bytes.Buffer{}

	_, err = io.Copy(b, r.Body)
	return b, err
}
