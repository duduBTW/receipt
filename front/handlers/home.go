//go:build js && wasm
// +build js,wasm

package handlers

import (
	"syscall/js"

	"github.com/dudubtw/receipt/constants"
	"github.com/dudubtw/receipt/front/jslayer"
)

func HomeSetup() func() {
	categoryCardClickHandler := jslayer.EventListener{
		Selector:  jslayer.Id(constants.IdCategoryCard),
		EventType: "click",
		Listener: func(this js.Value, args []js.Value) {
			elesments, err := jslayer.ElementQuerySelectorAll(this, "a")
			if err != nil || len(elesments) == 0 {
				return
			}

			jslayer.Click(elesments[0])
		},
	}

	categoryCardClickHandler.Add()
	return func() {
		categoryCardClickHandler.Remove()
	}
}
