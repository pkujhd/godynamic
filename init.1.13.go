//go:build go1.13 && !go1.20
// +build go1.13,!go1.20

package godynamic

import "C"

import (
	"unsafe"
)

const (
	_InitTaskSuffix = "..inittask"
)

func getInitFuncName(packagename string) string {
	return packagename + _InitTaskSuffix
}

// doInit is defined in package runtime
//go:linkname doInit runtime.doInit
func doInit(t unsafe.Pointer) // t should be a *runtime.initTask

func doInitialize(initTask unsafe.Pointer) error {
	doInit(initTask)
	return nil
}
