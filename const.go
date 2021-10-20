package godynamic

import "unsafe"

// size
const (
	PtrSize    = 4 << (^uintptr(0) >> 63)
	Uint32Size = int(unsafe.Sizeof(uint32(0)))
	IntSize    = int(unsafe.Sizeof(int(0)))
	UInt64Size = int(unsafe.Sizeof(uint64(0)))
)

const (
	itabLockName        = "runtime.itabLock"
	itabTableName       = "runtime.itabTable"
	firstmoduledataName = "runtime.firstmoduledata"
	lastmoduledatapName = "runtime.lastmoduledatap"
)
