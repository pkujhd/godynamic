//go:build go1.10 && !go1.13
// +build go1.10,!go1.13

package godynamic

import (
	"unsafe"
)

const (
	_InitTaskSuffix = ".init"
)

func getInitFuncName(packagename string) string {
	return packagename + _InitTaskSuffix
}

func doInitialize(initTask unsafe.Pointer) error {
	funcPtrContainer := (uintptr)(unsafe.Pointer(&initTask))
	runFunc := *(*func())(unsafe.Pointer(&funcPtrContainer))
	runFunc()
	return nil
}
