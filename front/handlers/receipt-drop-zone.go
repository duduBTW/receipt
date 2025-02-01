//go:build js && wasm
// +build js,wasm

package handlers

import (
	"fmt"
	"syscall/js"

	"github.com/a-h/templ"
	"github.com/dudubtw/receipt/constants"
	"github.com/dudubtw/receipt/front/jslayer"
	"github.com/dudubtw/receipt/renderer/components"
)

type ReceiptDropZoneHooks struct {
	Disabled    bool
	OnDragStart func()
	OnDragEnd   func()
	OnDrop      func(file js.Value)
}

func ReceiptDropZone(hooks ReceiptDropZoneHooks) []jslayer.EventListener {
	dragCounter := 0

	dragZoneState := jslayer.StateProps[bool]{
		Value:  false,
		Target: jslayer.Id(constants.IdReceiptDropZone),
		RenderComponent: func(value bool) templ.Component {
			return components.ReceiptDropZone()
		},
	}

	dragEnterHandler := jslayer.EventListener{
		Selector:  jslayer.Window,
		EventType: "dragenter",
		Listener: func(this js.Value, args []js.Value) {
			event := args[0]
			jslayer.PreventDefault(event)

			if hooks.Disabled {
				return
			}

			dragCounter++
			dragZoneState.SetAttribute("data-active", "true")

			if hooks.OnDragStart != nil {
				hooks.OnDragStart()
			}
		},
	}

	dragLeaveHandler := jslayer.EventListener{
		Selector:  jslayer.Window,
		EventType: "dragleave",
		Listener: func(this js.Value, args []js.Value) {
			event := args[0]
			jslayer.PreventDefault(event)

			dragCounter--
			if dragCounter != 0 {
				return
			}

			dragZoneState.SetAttribute("data-active", "false")

			if hooks.OnDragEnd != nil {
				hooks.OnDragEnd()
			}
		},
	}

	dragOverHandler := jslayer.EventListener{
		Selector:  jslayer.Window,
		EventType: "dragover",
		Listener: func(this js.Value, args []js.Value) {
			event := args[0]
			jslayer.PreventDefault(event)
		},
	}

	dropHandler := jslayer.EventListener{
		Selector:  jslayer.Window,
		EventType: "drop",
		Listener: func(this js.Value, args []js.Value) {
			fmt.Println("drop")
			event := args[0]
			jslayer.PreventDefault(event)
			dragCounter = 0
			dragZoneState.SetAttribute("data-active", "false")

			if hooks.OnDrop != nil {
				// Extract files
				files := event.Get("dataTransfer").Get("files")
				if files.Length() < 0 {
					return
				}

				firstFile := files.Index(0)
				hooks.OnDrop(firstFile)
			}
		},
	}

	return []jslayer.EventListener{dragEnterHandler, dragOverHandler, dragLeaveHandler, dropHandler}
}
