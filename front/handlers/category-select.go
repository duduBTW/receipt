//go:build js && wasm
// +build js,wasm

package handlers

import (
	"syscall/js"

	"github.com/dudubtw/receipt/front/jslayer"
	"github.com/dudubtw/receipt/front/service"
	"github.com/dudubtw/receipt/models"
	"github.com/dudubtw/receipt/renderer/components"
)

type CategorySelect struct {
	Id            string
	categories    []models.Category
	state         jslayer.StateProps[components.CategorySelectComponentProps]
	OnValueChange func(string)
	DefaultValue  string
	changeEvent   jslayer.EventListener
	OnMounted     func()
}

func (cSelect *CategorySelect) fetch() {
	if len(cSelect.categories) > 0 {
		return
	}

	categories, err := service.FetchCategories()
	if err != nil {
		return
	}

	cSelect.categories = categories
}

func (cSelect *CategorySelect) New() {
	cSelect.fetch()

	// Remove previous state
	cSelect.changeEvent.Remove()
	cSelect.changeEvent = jslayer.EventListener{
		Selector:  jslayer.Id(cSelect.Id),
		EventType: "change",
		Listener: func(this js.Value, args []js.Value) {
			if cSelect.OnValueChange == nil {
				return
			}
			event := args[0]
			value := event.Get("target").Get("value")
			if jslayer.IsNil(value) {
				return
			}

			cSelect.OnValueChange(value.String())
		},
	}

	// Create state
	cSelect.state = jslayer.NewState(jslayer.StateProps[components.CategorySelectComponentProps]{
		Target: jslayer.Id(cSelect.Id),
		Value: components.CategorySelectComponentProps{
			Id:         cSelect.Id,
			Categories: cSelect.categories,
		},
		RenderComponent: components.CategorySelectComponent,
		OnMounted: func(value components.CategorySelectComponentProps) {
			jslayer.CreateIcons()

			if cSelect.OnMounted != nil {
				cSelect.OnMounted()
			}
		},
	})

	// Listen to event listeners
	cSelect.changeEvent.Add()

	// Set default value
	cSelect.Set(cSelect.DefaultValue)

}

func (cSelect *CategorySelect) Remove() {
	cSelect.changeEvent.Remove()
}

func (cSelect *CategorySelect) Set(value string) error {
	selectElement, err := jslayer.QuerySelector(jslayer.Id(cSelect.Id) + " select")
	if err != nil {
		return err
	}

	selectElement.Set("value", value)
	return nil
}
