//go:build js && wasm
// +build js,wasm

package handlers

import (
	"fmt"
	"strconv"
	"syscall/js"
	"time"

	"github.com/dudubtw/receipt/constants"
	"github.com/dudubtw/receipt/front/jslayer"
	"github.com/dudubtw/receipt/renderer/components"
	"github.com/dudubtw/receipt/wasm"
)

var snackbar = func(status components.SnackbarStatus, message string) {
	GlobalFail()
}

func (GlobalDefiner) Snackbar(status components.SnackbarStatus, message string, devError error) {
	if wasm.Env == "development" {
		fmt.Println(devError.Error())
	}

	snackbar(status, message)
}

func SnackbarSetup() func() {
	var closeFns = make(map[string]func())

	snackbar = func(status components.SnackbarStatus, message string) {
		id := "id-" + strconv.FormatInt(time.Now().UnixMilli(), 10)
		fmt.Println(id)
		containerElement, err := jslayer.Element(jslayer.Id(constants.IdSnackbarContainer))
		if err != nil {
			Global.Error(err)
		}

		containerElement.AppendHTMLInside(components.Snackbar(components.SnackbarProps{
			Status: status,
			Label:  message,
			Id:     id,
		}))

		jslayer.CreateIcons()

		snackbarElement, err := jslayer.Element(jslayer.Id(id))
		if err != nil {
			Global.Error(err)
		}

		var closeBtnEventListner jslayer.EventListener
		var close = func() {
			closeBtnEventListner.Remove()
			delete(closeFns, id)

			snackbarElement.SetAttr("data-closing", "true")
			var animationEndEventListner jslayer.EventListener
			animationEndEventListner = jslayer.EventListener{
				Selector:  jslayer.Id(id),
				EventType: "transitionend",
				Listener: func(this js.Value, args []js.Value) {
					snackbarElement.Remove()
					animationEndEventListner.Remove()
				},
			}
			animationEndEventListner.Add()
		}

		closeBtnEventListner = jslayer.EventListener{
			Selector:  jslayer.Id(id) + " " + jslayer.Id(constants.IdSnackbarCloseButton),
			EventType: "click",
			Listener: func(this js.Value, args []js.Value) {
				close()
			},
		}
		closeBtnEventListner.Add()

		closeFns[id] = close

		time.AfterFunc(3*time.Second, close)
	}

	return func() {}
}
