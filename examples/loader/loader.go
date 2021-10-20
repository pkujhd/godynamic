package main

import (
	"flag"
	"fmt"
	"os"
	"unsafe"

	"github.com/pkujhd/godynamic"
)

func main() {
	library := flag.String("l", "", "load dynamic library file ")
	pluginpath := flag.String("p", "main", "pluginpath")
	run := flag.String("r", "Enter", "run function")
	times := flag.Int("times", 1, "run count")
	flag.Parse()
	for i := 0; i < *times; i++ {
		p, err := godynamic.Open(*library, *pluginpath)
		if err != nil {
			fmt.Println(err)
		} else {
			ptr, err := p.Lookup(*run)
			if err != nil {
				fmt.Println(err)
			} else {
				p := &ptr
				f := *(*func())(unsafe.Pointer(&p))
				f()
				os.Stdout.Sync()
			}
			p.Close()
		}
	}
}
