//go:build js && wasm
// +build js,wasm

package handlers

import (
	"fmt"

	"github.com/dudubtw/receipt/constants"
	"github.com/dudubtw/receipt/front/jslayer"
	"github.com/dudubtw/receipt/renderer/pages"
	"github.com/dudubtw/receipt/wasm"
)

type GlobalDefiner struct{}

var Global = GlobalDefiner{}

func (GlobalDefiner) Error(err error) {
	var componentMessage = ""
	if wasm.Env == "development" {
		componentMessage = err.Error()
	}

	jslayer.Render(pages.Error(componentMessage), jslayer.Id(constants.IdRoot))
	jslayer.StopApp()
}

func GlobalFail() {
	fmt.Println("This page doesn't have access to this function!")
}
