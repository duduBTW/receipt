//go:build js && wasm
// +build js,wasm

package handlers

import (
	"fmt"
	"syscall/js"

	"github.com/a-h/templ"
	"github.com/dudubtw/receipt/constants"
	jslayer "github.com/dudubtw/receipt/front/jslayer"
	"github.com/dudubtw/receipt/renderer/components"
)

func AppendAddCategoryComponent() {
	err := jslayer.AppendHTMLInside(jslayer.Id(constants.IdRoot), components.AddCategory())
	if err != nil {
		fmt.Println(err)
	}
	jslayer.CreateIcons()
}

func CreateModalSetup() func() {
	var modalState jslayer.StateProps[components.AddCategoryModalProps]

	var addCategoryButtonClickEvent = jslayer.EventListener{
		Selector:  jslayer.Id(constants.IdAddCategoryButton),
		EventType: "click",
		Listener: func(this js.Value, args []js.Value) {
			modalState.Set(components.AddCategoryModalProps{
				IsOpen: true,
			})
		},
	}

	modalState = jslayer.StateProps[components.AddCategoryModalProps]{
		Value:  components.AddCategoryModalDefaultProps,
		Target: jslayer.Id(constants.IdAddCategoryModal),
		RenderComponent: func(props components.AddCategoryModalProps) templ.Component {
			return components.AddCategoryModal(props)
		},
		OnMounted: func(value components.AddCategoryModalProps) {
			fmt.Println("Mounted", value)
			addCategoryButtonClickEvent.Add()
		},
	}

	AppendAddCategoryComponent()

	var eventList = []jslayer.EventListener{addCategoryButtonClickEvent}
	jslayer.AddEvents(eventList)
	return func() {
		jslayer.RemoveEvents(eventList)
	}
}
