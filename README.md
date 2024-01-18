# libdeflate (optimized zlib) for go

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

This ultra fast Go zlib library wraps the [libdeflate](https://github.com/ebiggers/libdeflate) zlib-, gzip- and deflate-(de)compression library for Go, using cgo.

It is **significantly faster** (4-5 times) than Go's standard compress/zlib/gzip/flate libraries (see [benchmarks](#benchmarks)) at the expense of not being able to stream data (e.g. from disk). Therefore, this library is optimal for the use case of **(de)compressing whole-buffered in-memory data**: If it fits into your memory, this library can (de)compress it much faster than the standard libraries (or even C zlib) can. 

```diff
+ If you start using this library, use V2.
```
For a better user experience right from the start use V2 via `go get github.com/4kills/go-libdeflate/v2`. 

## Table of Contents

- [Features](#features)
  - [Completeness of the Go Wrapper](#availability-of-the-original-libdeflate-api)
- [Installation](#installation)
  - [Prerequisites (cgo)](#prerequisites-working-cgo)
  - [Download and Installation](#download-and-installation)
- [Usage](#usage)
  - [Compress](#compress)
  - [Decompress](#decompress)
- [Notes](#notes)
- [Benchmarks](#benchmarks)
  - [Compression](#compression)
  - [Decompression](#decompression)
  - [Compression Ratio](#compression-ratio)
- [License](#license)
- [Attribution](#attribution)

# Features

- Super fast zlib, gzip, and raw deflate compression / decompression
- Convenience functions for quicker one-time compression / decompression 
- More (zlib/gzip/flate compatible) compression levels for better compression ratios than with standard zlib/gzip/flate
- Simple and clean API 

### Availability of the original [libdeflate](https://github.com/ebiggers/libdeflate) API:
   - [x] zlib/gzip/deflate compression
   - [x] zlib/gzip/deflate decompression
   - [x] Definite upper bound of compressed size
   - [x] Decompression w/ info about number of consumed bytes
   - [x] adler32 and crc32 checksums
   - [ ] ~~Custom memory allocator~~ : *no implementation planned, due to too little relevance for a high level Go API*
   
# Installation

For the library to work, you need cgo, libedflate (which is used by this library under the hood), and pkg-config (linker):

## Install [cgo](https://golang.org/cmd/cgo/)

**TL;DR**: Get **[cgo](https://golang.org/cmd/cgo/)** working.

In order to use this library with your Go source code, you must be able to use the Go tool **[cgo](https://golang.org/cmd/cgo/)**, which, in turn, requires a **GCC compiler**.

If you are on **Linux**, there is a good chance you already have GCC installed, otherwise just get it with your favorite package manager.

If you are on **MacOS**, Xcode - for instance - supplies the required tools.

If you are on **Windows**, you will need to install GCC.
I can recommend [tdm-gcc](https://jmeubank.github.io/tdm-gcc/) which is based
off of MinGW. Please note that [cgo](https://golang.org/cmd/cgo/) requires the 64-bit version (as stated [here](https://github.com/golang/go/wiki/cgo#windows)). 

For **any other** the procedure should be about the same. Just google. 

## Install [pkg-config](https://www.freedesktop.org/wiki/Software/pkg-config/) and [libdeflate](https://github.com/ebiggers/libdeflate)

This SDK uses [libdeflate](https://github.com/ebiggers/libdeflate) under the hood. For the SDK to work, you need to install `libdeflate` on your system which is super easy! 
Additionally we require [pkg-config](https://www.freedesktop.org/wiki/Software/pkg-config/) which facilitates linking `libdeflate` with this (cgo) SDK. 
How exactly you install these two packages depends on your operating system.

#### MacOS (HomeBrew):
```sh
brew install libdeflate
brew install pkg-config
```

#### Linux:
Use the package manager available on your distro to install the required packages. 

#### Windows (MinGW/WSL2):
Here, you can either use [WSL2](https://learn.microsoft.com/en-us/windows/wsl/install)
or [MinGW](https://www.mingw-w64.org/) and from there install the required packages. 

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

You can also pass nil for out, and the function will allocate a fitting buffer by itself:

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
// must be exactly the size of the output, if unknown, pass nil for out(see below)
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

There are also convenience functions that allow one-time compression to be easier, as well as functions to directly compress to zlib format.

# Notes

- **Do NOT use the <ins>same</ins> Compressor / Decompressor across multiple threads <ins>simultaneously</ins>.** However, you can create as many of them as you like, so if you want to parallelize your application, just create a compressor / decompressor for each thread. (See Memory Usage down below for more info)

- **Always `Close()` your Compressor / Decompressor when you are done with it** - especially if you create a new compressor/decompressor for each compression/decompression you undertake (which is generally discouraged anyway). As the C-part of this library is not subject to the Go garbage collector, the memory allocated by it must be released manually (by a call to `Close()`) to avoid memory leakage.

- Memory Usage: `Compressing` requires at least ~32 KiB of additional memory during execution, while `Decompressing` also requires at least ~32 KiB of additional memory during execution. 

# Benchmarks

These benchmarks were conducted with "real-life-type data" to ensure that these tests are most representative for an actual use case in a practical production environment.
As the zlib standard has been traditionally used for compressing smaller chunks of data, I have decided to follow suite by opting for Minecraft client-server communication packets, as they represent the optimal use case for this library. 

To that end, I have recorded 930 individual Minecraft packets, totalling 11,445,993 bytes in uncompressed data and 1,564,159 bytes in compressed data.
These packets represent actual client-server communication and were recorded using [this](https://github.com/haveachin/infrared) software.

The benchmarks were executed on different hardware and operating systems, including AMD and Intel processors, as well as all the out-of-the-box supported operating systems (Windows, Linux, MacOS). All the benchmarked functions/methods were executed hundreds of times, and the numbers you are about to see are the averages over all these executions.

The data was compressed using compression level 6 (current default of zlib). 

These benchmarks compare this library (blue) to the Go standard library (yellow) and show that this library performs **way** better in all cases. 

- <details>
  
    <summary> (A note regarding testing on your machine) </summary>
  
    Please note that you will need an Internet connection for some benchmarks to function. This is because these benchmarks will download the mc packets from [here](https://github.com/4kills/zlib_benchmark) and temporarily store them in memory for the duration of the benchmark tests, so this repository won't have to include the data in order save space on your machine and to make it a lightweight library.
  
  </details>

## Compression

![compression total](https://i.imgur.com/4eEFh5o.png)

This chart shows how long it took for the methods of this library (blue), and the standard library (yellow) to compress **all** of the 930 packets (~11.5 MB) on different systems in milliseconds. Note that the two rightmost data points were tested on **exactly the same** hardware in a dual-boot setup and that Linux seems to generally perform better than Windows.

![compression relative](https://i.imgur.com/K68qwFF.png)

This chart shows the time it took for this library's `Compress` (blue) to compress the data in nanoseconds, as well as the time it took for the standard library's `Write` (WriteStd, yellow) to compress the data in nanoseconds. The vertical axis shows percentages relative to the time needed by the standard library, thus you can see how much faster this library is. 

For example: This library only needed ~29% of the time required by the standard library to compress the packets on an Intel Core i5-6600K on Windows. 
That makes the standard library a substantial **~244.8% slower** than this library. 

## Decompression

![compression total](https://i.imgur.com/aW9L4cx.png)

This chart shows how long it took for the methods of this library (blue), and the standard library (yellow) to decompress **all** of the 930 packets (~1.5 MB) on different systems in milliseconds. Note that the two rightmost data points were tested on **exactly the same** hardware in a dual-boot setup and that Linux seems to generally perform better than Windows.

![decompression relative](https://i.imgur.com/lVhGfWO.png)

This chart shows the time it took for this library's `Decompress` (blue) to decompress the data in nanoseconds, as well as the time it took for the standard library's `Read` (ReadStd, Yellow) to decompress the data in nanoseconds. The vertical axis shows percentages relative to the time needed by the standard library, thus you can see how much faster this library is. 

For example: This library only needed ~34% of the time required by the standard library to decompress the packets on an Intel Core i5-6600K on Windows. 
That makes the standard library a substantial **~194.1% slower** than this library.

## Compression Ratio

Across all the benchmarks on all the different hardware / operating systems the compression ratios were consistent: 
This library had a compression ratio of **5.77** while the standard library had a compression ratio of **5.75**, which is a negligible difference. 

The compression ratio r is calculated as r = ucs / cs, where ucs = uncompressed size and cs = compressed size. 

# License

```txt
MIT License

Copyright (c) 2020 Dominik Ochs 

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

# Attribution

This library heavily depends on the C library [libdeflate](https://github.com/ebiggers/libdeflate), so everything in the folder native/libs is licensed under:

```
MIT License 
[for the license text of the MIT License, see LICENSE]

Copyright 2016 Eric Biggers
```

(See also: native/libs/LICENSE) 
