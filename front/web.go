//go:build js && wasm
// +build js,wasm

package main

import (
	"syscall/js"

	"github.com/dudubtw/receipt/constants"
	"github.com/dudubtw/receipt/front/handlers"
	jslayer "github.com/dudubtw/receipt/front/jslayer"
)

var appChan = make(chan struct{})

func stopApp() func() {
	close(appChan)
	return func() {}
}

func main() {
	// Global
	jslayer.RegisterFunction(constants.JsFunctionsImageModal, handlers.ImageModalSetup)
	jslayer.RegisterFunction(constants.JsStopApp, stopApp)

	// Page specific stuff
	jslayer.RegisterFunction(constants.JsFunctionsCreateCategory, handlers.CreateModalSetup)
	jslayer.RegisterFunction(constants.JsFunctionsReceipts, handlers.ReceiptsSetup)
	jslayer.RegisterFunction(constants.JsFunctionsHome, handlers.HomeSetup)

	loadCallback := js.FuncOf(func(this js.Value, args []js.Value) any {
		if !jslayer.IsNil(js.Global().Get("start")) {
			js.Global().Call("start")
		}
		return nil
	})
	defer loadCallback.Release()

	js.Global().Call("addEventListener", "load", loadCallback)
	<-appChan
}
