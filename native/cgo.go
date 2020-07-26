package native

/*
#cgo CFLAGS: -I${SRCDIR}/libs/
#cgo windows LDFLAGS: ${SRCDIR}/libs/libdeflatewin.a
#cgo linux LDFLAGS: ${SRCDIR}/libs/libdeflatelinux.a
#cgo darwin LDFLAGS: ${SRCDIR}/libs/libdeflatedarwin.a
*/
import "C"
