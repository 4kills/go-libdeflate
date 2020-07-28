# libdeflate for go

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

This library wraps the [libdeflate](https://github.com/ebiggers/libdeflate) zlib-, gzip- and deflate-(de)compression library for golang, using cgo.

It is **siginificantly faster** than go's standard compress/zlib/gzip/flate libraries (see [benchmarks](#benchmarks)) at the expense of not being able to stream data. Therefore, this library is optimal for the special use case of (de)compressing whole-buffered in-memory data: If it fits into your RAM, this library can (de)compress it much faster than the standard libraries can. 

## Table of Contents

- [Features](#features)
  - [Completeness of the go Wrapper](#availability-of-the-original-libdeflate-api)

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
