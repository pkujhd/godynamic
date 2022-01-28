//go:build go1.17 && !go1.19
// +build go1.17,!go1.19

package godynamic

func getSymbolName(name string) string {
	return name + ".abiinternal"
}
