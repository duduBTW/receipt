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
)

func ReceiptsSetup() func() {
	var isLoading = false
	var selectedRecepit models.Receipt
	var categorySelect = CategorySelect{
		Id:           constants.IdReceiptsSelectCategory,
		DefaultValue: jslayer.GetQueryParam(constants.ReceiptSearchParamCategory),
		OnValueChange: func(category string) {
			jslayer.SetQueryParam(constants.ReceiptSearchParamCategory, category)
		},
	}

	var modal = ReceiptModal(func(nr models.NewReceipt, v js.Value) (models.Receipt, error) {
		return service.UpdateRecepit(selectedRecepit.CopyNew(nr))
	})

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
	categorySelect.New()
	return func() {
		recepitClickHandler.Remove()
		categorySelect.Remove()
	}
}
