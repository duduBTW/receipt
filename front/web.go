//go:build js && wasm
// +build js,wasm

package main

import (
	"syscall/js"

	"github.com/dudubtw/receipt/constants"
	"github.com/dudubtw/receipt/front/handlers"
	jslayer "github.com/dudubtw/receipt/front/jslayer"
)

func main() {
	c := make(chan struct{}, 0)

	jslayer.RegisterFunction(constants.JsFunctionsCreateCategory, handlers.CreateModalSetup)

	loadCallback := js.FuncOf(func(this js.Value, args []js.Value) any {
		if !jslayer.IsNil(js.Global().Get("start")) {
			js.Global().Call("start")
		}
		return nil
	})
	defer loadCallback.Release()

	js.Global().Call("addEventListener", "load", loadCallback)

	<-c
}
