//go:build js && wasm
// +build js,wasm

package handlers

type GlobalDefiner struct{}

var Global = GlobalDefiner{}
