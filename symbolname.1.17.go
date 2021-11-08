//go:build go1.17 && !go1.18
// +build go1.17,!go1.18

package godynamic

func getSymbolName(name string) string {
	return name + ".abiinternal"
}
