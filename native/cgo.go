package native

/*
#cgo CFLAGS: -I${SRCDIR}/libs/
#cgo windows LDFLAGS: ${SRCDIR}/libs/libdeflatewin.a
#cgo linux LDFLAGS: ${SRCDIR}/libs/libdeflatelinux.a
*/
import "C"
