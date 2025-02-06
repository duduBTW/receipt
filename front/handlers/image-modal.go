//go:build js && wasm
// +build js,wasm

package handlers

import (
	"syscall/js"

	"github.com/dudubtw/receipt/constants"
	"github.com/dudubtw/receipt/front/jslayer"
	"github.com/dudubtw/receipt/renderer/components"
)

var openModalFn = func(imageUrl string) {
	GlobalFail()
}
var closeModalFn = GlobalFail

func (GlobalDefiner) OpenImage(imageUrl string) {
	openModalFn(imageUrl)
}
func (GlobalDefiner) CloseImage() {
	closeModalFn()
}

func ImageModalSetup() func() {
	var closeModalKeyboardEvent jslayer.EventListener
	var closeModalClickHandler jslayer.EventListener

	var modal = jslayer.StateProps[components.ImageModalProps]{
		Value:           components.DefaultImageModalProps,
		Target:          jslayer.Id(constants.IdImageModal),
		RenderComponent: components.ImageModal,
		OnMounted: func(value components.ImageModalProps) {
			if !value.IsOpen {
				closeModalClickHandler.Remove()
				return
			}

			closeModalClickHandler.Add()
			jslayer.CreateIcons()
		},
	}

	closeModalKeyboardEvent = jslayer.EventListener{
		Selector:  jslayer.Window,
		EventType: "keydown",
		Listener: func(this js.Value, args []js.Value) {
			event := args[0]
			key := event.Get("key").String()

			if key == "Escape" {
				closeModalFn()
			}
		},
	}

	closeModalClickHandler = jslayer.EventListener{
		Selector:  jslayer.Id(constants.IdImageModalCloseButton),
		EventType: "click",
		Listener: func(this js.Value, args []js.Value) {
			closeModalFn()
		},
	}

	openModalFn = func(imageUrl string) {
		modal.Set(components.ImageModalProps{
			IsOpen:   true,
			ImageUrl: imageUrl,
		})
	}

	closeModalFn = func() {
		modal.Set(components.DefaultImageModalProps)
	}

	closeModalKeyboardEvent.Add()
	return func() {
		closeModalKeyboardEvent.Remove()
	}
}
