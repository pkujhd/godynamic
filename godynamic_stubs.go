// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !linux || !cgo
// +build !linux !cgo

package godynamic

import (
	"errors"
	"unsafe"
)

func getaddrs() error {
	return errors.New("godynamic: not implemented")
}

func lookup(p *Dynamiclib, symName string) (unsafe.Pointer, error) {
	return unsafe.Pointer(uintptr(0)), errors.New("godynamic: not implemented")
}

func openplugin(name, libpath string) (*Dynamiclib, error) {
	return nil, errors.New("godynamic: not implemented")
}

func closeplugin(p *Dynamiclib) error {
	return errors.New("godynamic: not implemented")
}
