
# GoDynamic

GoDynamic can load and run Golang dynamic library compiled by -buildmode=shared -linkshared

## How does it work?

GoDynamic works like a dynamic linker: use dl loads an .so libraray and lookup symbol, can unload.

Please note that GoDynamic is not a scripting engine. All features of Go are supported, and run just as fast and lightweight as native Go code.

## Comparison with plugin

GoDynamic reuses runtime by libstd.so, which makes it much smaller. And code loaded by GoDynamic is unloadable.

## Build

**Make sure you're using go >= 1.10.**

## Examples

```
export GO111MODULE=auto
go install -buildmode=shared std
go build -linkshared -v github.com/pkujhd/godynamic/examples/loader

go install -v -buildmode=shared -linkshared github.com/pkujhd/godynamic/examples/ubase
./loader -l $GOPATH/pkg/linux_`go env GOARCH`_dynlink/libgithub.com-pkujhd-godynamic-examples-ubase.so -p github.com/pkujhd/godynamic/examples/ubase.so -r github.com/pkujhd/godynamic/examples/ubase.Enter -times 10

go install -v -buildmode=shared -linkshared github.com/pkujhd/godynamic/examples/uschedule
./loader -l $GOPATH/pkg/linux_`go env GOARCH`_dynlink/libgithub.com-pkujhd-godynamic-examples-uschedule.so -p github.com/pkujhd/godynamic/examples/uschedule.so -r github.com/pkujhd/godynamic/examples/uschedule.Enter -times 10

go install -v -buildmode=shared -linkshared github.com/pkujhd/godynamic/examples/uhttp
./loader -l $GOPATH/pkg/linux_`go env GOARCH`_dynlink/libgithub.com-pkujhd-godynamic-examples-uhttp.so -p github.com/pkujhd/godynamic/examples/uhttp.so -r github.com/pkujhd/godynamic/examples/uhttp.Enter

```

## Warning
golang buildmode(-buildmode=shared) will be not support after golang 1.18
This has currently only been tested and developed on:
Golang 1.10-1.16 (x64/x86, linux)
