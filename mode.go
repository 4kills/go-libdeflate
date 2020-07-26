package libdeflate

type Mode int

const(
	ModeDEFLATE Mode = iota
	ModeZlib
	ModeGzip
)
