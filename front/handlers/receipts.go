//go:build js && wasm
// +build js,wasm

package handlers

import (
	"fmt"
	"syscall/js"

	"github.com/dudubtw/receipt/constants"
	"github.com/dudubtw/receipt/front/jslayer"
	"github.com/dudubtw/receipt/front/service"
	"github.com/dudubtw/receipt/models"
	"github.com/dudubtw/receipt/renderer/components"
)

func FocusSelectOnOpen() {
	var wasRedirectedFromCategorySelect = jslayer.GetQueryParam(constants.QueryParamReceiptFromCategorySelect)
	fmt.Println(wasRedirectedFromCategorySelect)
	if wasRedirectedFromCategorySelect != constants.QueryParamReceiptFromCategorySelectTrueValue {
		return
	}

	jslayer.Focus(jslayer.Id(constants.IdReceiptsSelectCategory + " select"))
}

func ReceiptsSetup() func() {
	var isLoading = false
	var selectedRecepit models.Receipt
	var categorySelect = CategorySelect{
		Id:           constants.IdReceiptsSelectCategory,
		DefaultValue: jslayer.GetQueryParam(constants.ReceiptSearchParamCategory),
		OnValueChange: func(category string) {
			jslayer.SetQueryParam([][2]string{
				{
					constants.ReceiptSearchParamCategory, category,
				},
				{
					constants.QueryParamReceiptFromCategorySelect,
					constants.QueryParamReceiptFromCategorySelectTrueValue,
				},
			})
		},
		OnMounted: FocusSelectOnOpen,
	}

	var modal = ReceiptModal(ReceiptModalParams{
		DefaultProps: components.AddCategoryModalProps{
			IsOpen:      false,
			Title:       "Editar comprovante",
			ButtonLabel: "Salvar",
		},
		OnSubmit: func(nr models.NewReceipt) (models.Receipt, error) {
			return service.UpdateRecepit(selectedRecepit.CopyNew(nr))
		},
	})

	var openReceiptImage = func(element js.Value) {
		Global.OpenImage(element.Call("getAttribute", "src").String())
	}

	recepitImageClickHandler := jslayer.EventListener{
		Selector:  jslayer.Id(constants.IdReceiptCardImage),
		EventType: "click",
		Listener: func(this js.Value, args []js.Value) {
			event := args[0]
			jslayer.StopPropagation(event)
			openReceiptImage(this)
		},
	}

	recepitImageKeyDownHandler := jslayer.EventListener{
		Selector:  jslayer.Id(constants.IdReceiptCardImage),
		EventType: "keydown",
		Listener: func(this js.Value, args []js.Value) {
			event := args[0]
			key := event.Get("key").String()

			if key == "Enter" {
				openReceiptImage(this)

			}
		},
	}

	recepitClickHandler := jslayer.EventListener{
		Selector:  jslayer.Id(constants.IdReceiptCard),
		EventType: "click",
		Listener: func(this js.Value, args []js.Value) {
			if isLoading {
				return
			}

			isLoading = true
			rawId := this.Call("getAttribute", constants.DataRecepitId)
			if jslayer.IsNil(rawId) {
				fmt.Println("Failed to edit recepit!")
			}

			recepit, err := service.FetchReceipt(rawId.String())
			if err != nil {
				return
			}

			fmt.Println(recepit)
			isLoading = false
			selectedRecepit = recepit
			modal.Open(recepit)
		},
	}

	recepitClickHandler.Add()
	recepitImageClickHandler.Add()
	categorySelect.New()
	recepitImageKeyDownHandler.Add()
	return func() {
		recepitImageKeyDownHandler.Remove()
		recepitImageClickHandler.Remove()
		recepitClickHandler.Remove()
		categorySelect.Remove()
	}
}
