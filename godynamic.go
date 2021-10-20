package godynamic

import (
	"unsafe"
)

// Dynamiclib is a loaded Go plugin.
type Dynamiclib struct {
	md       *moduledata
	libpath  string
	filepath string
	h        uintptr
	err      string        // set if plugin failed to load
	loaded   chan struct{} // closed when loaded
}

// Open opens a Go plugin.
// If a path has already been opened, then the existing *Dynamiclib is returned.
// It is safe for concurrent use by multiple goroutines.
func Open(path, libpath string) (*Dynamiclib, error) {
	if firstmoduledatap == nil {
		err := getaddrs()
		if err != nil {
			return nil, err
		}
	}
	return openplugin(path, libpath)
}

func (p *Dynamiclib) Close() error {
	return closeplugin(p)
}

// Lookup searches for a symbol named symName in plugin p.
// A symbol is any exported variable or function.
// It reports an error if the symbol is not found.
// It is safe for concurrent use by multiple goroutines.
func (p *Dynamiclib) Lookup(symName string) (unsafe.Pointer, error) {
	return lookup(p, symName)
}
