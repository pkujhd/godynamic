//go:build go1.17 && !go1.20
// +build go1.17,!go1.20

package godynamic

func getSymbolName(name string) string {
	return name + ".abiinternal"
}
