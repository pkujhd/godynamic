//go:build go1.10
// +build go1.10

package godynamic

import "unsafe"

// Mutual exclusion locks.  In the uncontended case,
// as fast as spin locks (just a few user-level instructions),
// but on the contention path they sleep in the kernel.
// A zeroed Mutex is unlocked (no need to initialize each lock).
type mutex struct {
	// Futex-based impl treats it as uint32 key,
	// while sema-based impl as M* waitm.
	// Used to be a union, but unions break precise GC.
	key uintptr
}

//go:linkname lock runtime.lock
func lock(l *mutex)

//go:linkname unlock runtime.unlock
func unlock(l *mutex)

//go:linkname atomicstorep runtime.atomicstorep
func atomicstorep(ptr unsafe.Pointer, new unsafe.Pointer)

// layout of Itab known to compilers
// allocated in non-garbage-collected memory
// Needs to be in sync with
// ../cmd/compile/internal/gc/reflect.go:/^func.dumptabs.
type itab struct {
	inter *interfacetype
	_type *_type
	hash  uint32 // copy of _type.hash. Used for type switches.
	_     [4]byte
	fun   [1]uintptr // variable sized. fun[0]==0 means _type does not implement inter.
}

const itabInitSize = 512

// Note: change the formula in the mallocgc call in itabAdd if you change these fields.
type itabTableType struct {
	size    uintptr             // length of entries array. Always a power of 2.
	count   uintptr             // current number of filled entries.
	entries [itabInitSize]*itab // really [size] large
}

var itabLock *mutex = nil
var itabTable *itabTableType = nil

//go:linkname itabHashFunc runtime.itabHashFunc
func itabHashFunc(inter *interfacetype, typ *_type) uintptr

//go:linkname itabAdd runtime.itabAdd
func itabAdd(m *itab)

func additabs(module *moduledata) {
	lock(itabLock)
	for _, itab := range module.itablinks {
		itabAdd(itab)
	}
	unlock(itabLock)
}

func removeitabs(module *moduledata) bool {
	lock(itabLock)
	defer unlock(itabLock)
	for i := uintptr(0); i < itabTable.size; i++ {
		p := (**itab)(add(unsafe.Pointer(&itabTable.entries), i*PtrSize))
		m := (*itab)(loadp(unsafe.Pointer(p)))
		if m != nil {
			uintptrm := uintptr(unsafe.Pointer(m))
			inter := uintptr(unsafe.Pointer(m.inter))
			_type := uintptr(unsafe.Pointer(m._type))
			if (inter >= module.types && inter <= module.etypes) || (_type >= module.types && _type <= module.etypes) ||
				(uintptrm >= module.types && uintptrm <= module.etypes) {
				atomicstorep(unsafe.Pointer(p), unsafe.Pointer(nil))
				itabTable.count = itabTable.count - 1
			}
		}
	}
	return true
}
