//go:build js && wasm
// +build js,wasm

package jslayer

import "syscall/js"

func IsNil(value js.Value) bool {
	valueType := value.Type()
	return valueType == js.TypeNull || valueType == js.TypeUndefined
}
