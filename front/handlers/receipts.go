//go:build js && wasm
// +build js,wasm

package handlers

import (
	"github.com/dudubtw/receipt/constants"
	"github.com/dudubtw/receipt/front/jslayer"
)

func ReceiptsSetup() func() {
	var categorySelect = CategorySelect{
		Id:           constants.IdReceiptsSelectCategory,
		DefaultValue: jslayer.GetQueryParam(constants.ReceiptSearchParamCategory),
		OnValueChange: func(category string) {
			jslayer.SetQueryParam(constants.ReceiptSearchParamCategory, category)
		},
	}

	categorySelect.New()
	return func() {
		categorySelect.Remove()
	}
}
