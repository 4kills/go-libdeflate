# libdeflate for go

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

This library wraps the [libdeflate](https://github.com/ebiggers/libdeflate) zlib-, gzip- and deflate-(de)compression library for golang, using cgo.

It is **siginificantly faster** than go's standard compress/zlib/gzip/flate libraries (see [benchmarks](#benchmarks)) at the expense of not being able to stream data. Therefore, this library is optimal for the special use case of (de)compressing whole-buffered in-memory data: If it fits into your RAM, this library can (de)compress it much faster than the standard libraries can. 

## Table of Contents

- [Features](#features)
  - [Completeness of the go Wrapper](#availability-of-the-original-libdeflate-api  
- [Installation](#installation)
  - [Prerequisites (cgo)](#prerequisites-working-cgo)
  - [Download and Installation](#download-and-installation)
- [Usage](#usage)
  - [Compress](#compress)
  - [Decompress](#decompress)

# Features

- Super fast zlib, gzip, and raw deflate compression / decompression
- Convenience functions for quicker one-time compression / decompression 
- More (zlib/gzip/flate compatible) compression levels for better compression ratios than with standard zlib/gzip/flate
- Simple and clean API 

### Availability of the original [libdeflate](https://github.com/ebiggers/libdeflate) API:
   - [x] zlib/gzip/deflate compression
   - [x] zlib/gzip/deflate decompression
   - [ ] Definite upper bound of compressed size
   - [ ] Decompression w/ info about number of consumed bytes
   - [ ] adler32 and crc32 checksums
   - [ ] ~~Custom memory allocator~~ : *no implementation planned, due to too little relevance for a high level go API*
   
# Installation

## Prerequisites (working [cgo](https://golang.org/cmd/cgo/))

In order to use this library with your go source code, you must be able to use the go tool **[cgo](https://golang.org/cmd/cgo/)**, which in turn requires a **GCC compiler**.

If you are on **Linux**, there is a good chance you already have GCC installed, otherwise just get it with your favorite package manager.

If you are on **MacOS**, Xcode - for instance - supplies the required tools.

If you are on **Windows**, you will need to install GCC.
I can recommend [tdm-gcc](https://jmeubank.github.io/tdm-gcc/) which is based
off of MinGW. Please note that [cgo](https://golang.org/cmd/cgo/) requires the 64-bit version (as stated [here](https://github.com/golang/go/wiki/cgo#windows)). 

For **any other** the procedure should be about the same. Just google. 

**TL;DR**: Get **[cgo](https://golang.org/cmd/cgo/)** working.

## Download and Installation

If you want to build for `$GOARCH=amd64` and either **Windows, Linux** or **MacOS** just go get this library and everything will work right away. 

(You may also use go modules (available since go 1.11) to get the version of a specific branch or tag, if you want to try out or use experimental features. However, beware that these versions are not necessarily guaranteed to be stable or thoroughly tested.)

<details>

<summary> Instructions for building for different GOOS,GOARCH </summary>


First of all, it is not encouraged to build for non-64-bit archs, as this library works best for 64-bit systems. 

A list of possible GOOS,GOARCH combinations can be viewed [here](https://golang.org/doc/install/source#environment). 

**Instructions:**

1. You will need to compile and build the C [libdeflate](https://github.com/ebiggers/libdeflate) library for your target system by cloning the [repository](https://github.com/ebiggers/libdeflate) and executing the Makefile (specifying the build). *You should always use GCC for compilation as this produces the fastest libraries.* 

2. Step 1 will yield compiled library files. You are going to want to use the static library (usually ending with .a \[in case of windows, rename .lib to .a\]), give it an adequate name (like `libdeflate_GOOS_GOARCH.a`) and copy it to the native/libs folder of this library.

3. Go to the native/cgo.go file, which should roughly look like this: 
```go
package native

/*
#cgo CFLAGS: -I${SRCDIR}/libs/
#cgo windows,amd64 LDFLAGS: ${SRCDIR}/libs/libdeflate_windows_amd64.a 
#cgo linux,amd64 LDFLAGS: ${SRCDIR}/libs/libdeflate_linux_amd64.a
#cgo darwin,amd64 LDFLAGS: ${SRCDIR}/libs/libdeflate_darwin_amd64.a
*/
import "C"
```
Now you want to add your build of libdeflate to the cgo directives, more specifically, to the linker flags, like this (omit the '+'): 
```diff
package native

/*
#cgo CFLAGS: -I${SRCDIR}/libs/
+#cgo GOOS,GOARCH LDFLAGS: ${SRCDIR}/libs/libdeflate_GOOS_GOARCH.a
#cgo windows,amd64 LDFLAGS: ${SRCDIR}/libs/libdeflate_windows_amd64.a
#cgo linux,amd64 LDFLAGS: ${SRCDIR}/libs/libdeflate_linux_amd64.a
#cgo darwin,amd64 LDFLAGS: ${SRCDIR}/libs/libdeflate_darwin_amd64.a
*/
import "C"
```

That's it! It should work now!

</details>

# Usage

## Compress

First, you need to create a compressor that can be used for any type of compression. 

You can also specify a level of compression for which holds true: The higher the level, the higher the compression at the expense of speed. 
`-> lower level = fast, bad compression; higher level = slow, good compression`. Test what works best for your application but generally the `DefaultCompressionLevel` is fine most of the time. 

```go
// Compressor with default compression level. Errors if out of memory
c, err := libdeflate.NewCompressor()

// Compressor with custom compression level. Errors if out of memory or if an illegal level was passed. 
c, err = libdeflate.NewCompressorLevel(2)
```

Then you can compress the actual data with a given mode of compression (currently supported: zlib, gzip, raw deflate): 

```go 
decomp := []byte(`Some data to compress: May be anything,  
    but it might be a good idea to only compress data that exceeds a certain threshold in size, 
    as compressed data can become larger (due to overhead)`)
comp := make([]byte, len(decomp)) // supplying a fitting buffer is in all cases the fastest approach

n, _, err := c.Compress(decomp, comp, libdeflate.ModeZlib) // Errors if buffer was too short
comp = comp[:n]
```

You can also pass nil for out and the function will allocate a fitting buffer by itself:

```go
_, comp, err = c.Compress(decomp, nil, libdeflate.ModeZlib)
```

After you are done with the compressor, do not forget to close it to free c-allocated-memory:

```go
c.Close()
```

## Decompress 

As with compression, you need to create a decompressor which can also be used for any type of decompression at any compression level:

```go
// Doesn't need a compression level; works universally. Errors if out of memory.
dc, err := libdeflate.NewDecompressor() 
```

Then you can decompress the actual data with a given mode of compression (currently supported: zlib, gzip, raw deflate): 

```go 
decompressed := make([]byte, len(decomp)) 

_, err = dc.Decompress(comp, decompressed, ModeZlib) 
```

Just like with compress you can also pass nil and get a fitting buffer:

```go
decompressed, err = dc.Decompress(comp, nil, ModeZlib)
```

After you are done with the decompressor, do not forget to close it to free c-allocated-memory:

```go
dc.Close()
```

There are also convencience methods that allow one-time compression to be more easy as well as directly compress to zlib format.
