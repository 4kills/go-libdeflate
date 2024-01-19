package libdeflate

import "runtime"

// panicFreeCloser is a type that can be closed without panicking if it is already closed. This is used for AutoClose functionality.
type panicFreeCloser interface {
	PanicFreeClose()
}

// attachAutoClose attaches a finalizer to the given panicFreeCloser that calls PanicFreeClose() when the object is garbage collected.
// Used for AutoClose functionality. Do not call this function directly.
func attachAutoClose(c panicFreeCloser) {
	runtime.SetFinalizer(c, func(finalized panicFreeCloser) {
		finalized.PanicFreeClose()
	})
}
