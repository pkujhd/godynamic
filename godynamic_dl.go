//go:build linux && cgo
// +build linux,cgo

package godynamic

/*
#cgo linux LDFLAGS: -ldl
#include <dlfcn.h>
#include <limits.h>
#include <stdlib.h>
#include <stdint.h>

#include <stdio.h>

static uintptr_t pluginOpen(const char* path, char** err) {
	void* h = dlopen(path, RTLD_NOW|RTLD_GLOBAL);
	if (h == NULL) {
		*err = (char*)dlerror();
	}
	return (uintptr_t)h;
}

static void* pluginLookup(uintptr_t h, const char* name, char** err) {
	void* r = dlsym((void*)h, name);
	if (r == NULL) {
		*err = (char*)dlerror();
	}
	return r;
}

static int pluginClose(void* h, char** err){
	int r = dlclose(h);
	if (r != 0) {
		*err = (char*)dlerror();
	}
	return r;
}
*/
import "C"

import (
	"errors"
	"strconv"
	"sync"
	"unsafe"
)

var (
	plugins map[string]*Dynamiclib
)

func getaddrs() error {
	//because -buildmode=shared, go:linkname doesn't work as scheduled,
	//so lookup firstmoduledata/lastmoduledatap/itabTable/itabLock in libstd.so
	if firstmoduledatap != nil {
		return nil
	}
	var cErr *C.char
	name := make([]byte, len(itabLockName)+1)
	copy(name, itabLockName)
	addr := C.pluginLookup(0, (*C.char)(unsafe.Pointer(&name[0])), &cErr)
	if cErr != nil {
		return errors.New(`find ` + itabLockName + `failed!error message:` + C.GoString(cErr))
	}
	itabLock = (*mutex)(unsafe.Pointer(addr))
	name = make([]byte, len(itabTableName)+1)
	copy(name, itabTableName)
	addr = C.pluginLookup(0, (*C.char)(unsafe.Pointer(&name[0])), &cErr)
	if cErr != nil {
		return errors.New(`find ` + itabTableName + `failed!error message:` + C.GoString(cErr))
	}
	itabTable = *(**itabTableType)(unsafe.Pointer(addr))
	name = make([]byte, len(firstmoduledataName)+1)
	copy(name, firstmoduledataName)
	addr = C.pluginLookup(0, (*C.char)(unsafe.Pointer(&name[0])), &cErr)
	if cErr != nil {
		return errors.New(`find ` + firstmoduledataName + `failed!error message:` + C.GoString(cErr))
	}
	firstmoduledatap = (*moduledata)(unsafe.Pointer(addr))
	name = make([]byte, len(lastmoduledatapName)+1)
	copy(name, lastmoduledatapName)
	addr = C.pluginLookup(0, (*C.char)(unsafe.Pointer(&name[0])), &cErr)
	if cErr != nil {
		return errors.New(`find ` + lastmoduledatapName + `failed!error message:` + C.GoString(cErr))
	}
	lastmoduledatapp = (**moduledata)(unsafe.Pointer(addr))
	lastmoduledatap = *lastmoduledatapp
	return nil
}

func openplugin(name, libpath string) (*Dynamiclib, error) {
	cPath := make([]byte, C.PATH_MAX+1)
	cRelName := make([]byte, len(name)+1)
	copy(cRelName, name)
	if C.realpath(
		(*C.char)(unsafe.Pointer(&cRelName[0])),
		(*C.char)(unsafe.Pointer(&cPath[0]))) == nil {
		return nil, errors.New(`godynamic.Open("` + name + `"): realpath failed`)
	}

	filepath := C.GoString((*C.char)(unsafe.Pointer(&cPath[0])))

	pluginsMu.Lock()
	defer pluginsMu.Unlock()
	if plugins == nil {
		plugins = make(map[string]*Dynamiclib)
	}
	if p := plugins[libpath]; p != nil {
		if p.err != "" {
			return nil, errors.New(`godynamic.Open("` + name + `"): ` + p.err + ` (previous failure)`)
		}
		<-p.loaded
		return p, nil
	}
	var cErr *C.char

	h := C.pluginOpen((*C.char)(unsafe.Pointer(&cPath[0])), &cErr)
	if h == 0 {
		return nil, errors.New(`godynamic.Open("` + name + `"): ` + C.GoString(cErr))
	}
	if len(name) > 3 && name[len(name)-3:] == ".so" {
		name = name[:len(name)-3]
	}
	md, errstr := lastmoduleinit()
	if errstr != "" {
		plugins[libpath] = &Dynamiclib{
			filepath: filepath,
			err:      errstr,
		}
		return nil, errors.New(`godynamic.Open("` + name + `"): ` + errstr)
	}
	p := &Dynamiclib{
		md:       md,
		filepath: filepath,
		libpath:  libpath,
		loaded:   make(chan struct{}),
		h:        uintptr(h),
	}
	plugins[libpath] = p

	//call init functions
	initName := getInitFuncName(libpath)
	init := make([]byte, len(initName)+1)
	copy(init, initName)
	initTask := C.pluginLookup(h, (*C.char)(unsafe.Pointer(&init[0])), &cErr)
	if initTask != nil {
		doInitialize(initTask)
	}

	close(p.loaded)
	return p, nil
}

func lookup(p *Dynamiclib, symName string) (unsafe.Pointer, error) {
	var cErr *C.char
	newName := getSymbolName(symName)
	name := make([]byte, len(newName)+1)
	copy(name, newName)
	ptr := C.pluginLookup(C.ulong(p.h), (*C.char)(unsafe.Pointer(&name[0])), &cErr)
	if ptr != nil {
		return ptr, nil
	}
	return nil, errors.New("godynamic: symbol " + symName + " not found in plugin " + p.libpath)
}

func closeplugin(p *Dynamiclib) error {
	pluginsMu.Lock()
	defer pluginsMu.Unlock()
	delete(plugins, p.libpath)
	removeitabs(p.md)
	removeModule(p.md)
	modulesinit()
	p.md = nil
	var cErr *C.char
	if p.h != 0 {
		r := C.pluginClose(unsafe.Pointer(p.h), &cErr)
		if r != 0 {
			return errors.New("godynamic:" + p.libpath + " close failed. return:" + strconv.Itoa(int(r)))
		}
	}
	return nil
}

func lastmoduleinit() (md *moduledata, errstr string) {
	for module := firstmoduledatap; module != nil; module = module.next {
		if module.next == nil {
			md = module
			break
		}
	}

	if md == nil {
		panic("runtime: no plugin module data")
	}

	for _, pmd := range activeModules() {
		if pmd.pluginpath == md.pluginpath {
			continue
		}

		if inRange(pmd.text, pmd.etext, md.text, md.etext) ||
			inRange(pmd.bss, pmd.ebss, md.bss, md.ebss) ||
			inRange(pmd.data, pmd.edata, md.data, md.edata) ||
			inRange(pmd.types, pmd.etypes, md.types, md.etypes) {
			println("plugin: new module data overlaps with previous moduledata")
			println("\tpmd.text-etext=", hex(pmd.text), "-", hex(pmd.etext))
			println("\tpmd.bss-ebss=", hex(pmd.bss), "-", hex(pmd.ebss))
			println("\tpmd.data-edata=", hex(pmd.data), "-", hex(pmd.edata))
			println("\tpmd.types-etypes=", hex(pmd.types), "-", hex(pmd.etypes))
			println("\tmd.text-etext=", hex(md.text), "-", hex(md.etext))
			println("\tmd.bss-ebss=", hex(md.bss), "-", hex(md.ebss))
			println("\tmd.data-edata=", hex(md.data), "-", hex(md.edata))
			println("\tmd.types-etypes=", hex(md.types), "-", hex(md.etypes))
			panic("plugin: new module data overlaps with previous moduledata")
		}
	}
	for _, pkghash := range md.pkghashes {
		if pkghash.linktimehash != *pkghash.runtimehash {
			md.bad = true
			return nil, "dynamic library was built with a different version of package " + pkghash.modulename
		}
	}

	// Initialize the freshly loaded module.
	modulesinit()
	typelinksinit()

	moduledataverify1(md)
	additabs(md)

	return md, ""
}

//go:linkname inRange runtime.inRange
func inRange(r0, r1, v0, v1 uintptr) bool

var pluginsMu sync.Mutex
