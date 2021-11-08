//go:build go1.10 && !go1.17
// +build go1.10,!go1.17

package godynamic

func getSymbolName(name string) string {
	return name
}
