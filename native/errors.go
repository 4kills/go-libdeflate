package native

import "errors"

var (
	ErrorOutOfMemory                     = errors.New("libdeflate: native: out of memory")
	ErrorInvalidLevel                    = errors.New("libdeflate: native: illegal compression level")
	ErrorShortBuffer                     = errors.New("libdeflate: native: short buffer")
	ErrorNoInput                         = errors.New("libdeflate: native: empty input")
	ErrorBadData                         = errors.New("libdeflate: native: bad data: data was corrupted, invalid or unsupported")
	ErrorUnknown                         = errors.New("libdeflate: native: unknown error code from c library")
	ErrorShortOutput                     = errors.New("libdeflate: native: buffer too long: decompressed to fewer bytes than expected, indicating possible error in decompression. Make sure your out buffer has the exact length of the decompressed data or pass nil for out")
	ErrorAlreadyClosed                   = errors.New("libdeflate: native: (de-)compressor already closed. It must not be used anymore")
	ErrorInsufficientDecompressionFactor = errors.New("libdeflate: native: your compressed data seems to be extraordinarily large when decompressed. " +
		"However, this could also indicate corrupted data. The current maximum decompression factor does not allow for larger decompression, try to increase it")

	// checked error (in native)
	ErrorInsufficientSpace = errors.New("libdeflate: native: buffer too short. Retry with larger buffer")
)
