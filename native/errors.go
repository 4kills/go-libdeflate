package native

import "errors"

var (
	errorOutOfMemory  = errors.New("libdeflate: native: out of memory")
	errorInvalidLevel = errors.New("libdeflate: native: illegal compression level")
	errorShortBuffer  = errors.New("libdeflate: native: short buffer")
	errorNoInput      = errors.New("libdeflate: native: empty input")
)
