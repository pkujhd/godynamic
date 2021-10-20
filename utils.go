package godynamic

import (
	"unsafe"
)

type hex uint64

//go:linkname add runtime.add
func add(p unsafe.Pointer, x uintptr) unsafe.Pointer

//go:linkname adduintptr runtime.add
func adduintptr(p uintptr, x int) unsafe.Pointer

//go:nosplit
//go:noinline
//see runtime.internal.atomic.Loadp
func loadp(ptr unsafe.Pointer) unsafe.Pointer {
	return *(*unsafe.Pointer)(ptr)
}
