//go:build js && wasm
// +build js,wasm

package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"syscall/js"

	"github.com/dudubtw/receipt/front/jslayer"
	"github.com/dudubtw/receipt/models"
)

type Result struct {
	receipt models.Receipt
	err     error
}

func UploadRecepit(receipt models.NewReceipt, file js.Value) (models.Receipt, error) {
	ch := make(chan Result)

	formData := js.Global().Get("FormData").New()
	formData.Call("append", models.NewReceiptFormFieldsInstance.File, file)
	formData.Call("append", models.NewReceiptFormFieldsInstance.CategoryID, receipt.CategoryID)
	formData.Call("append", models.NewReceiptFormFieldsInstance.Date, receipt.Date)

	opts := map[string]interface{}{
		"method": "POST",
		"body":   formData,
	}

	go func() {
		promise := js.Global().Call("fetch", "/upload", js.ValueOf(opts))

		json := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			receiptJSON := args[0]
			var receipt models.Receipt
			if err := json.Unmarshal([]byte(jslayer.JSON{Value: receiptJSON}.Stringify()), &receipt); err != nil {
				ch <- Result{models.Receipt{}, err}
			}

			ch <- Result{receipt, nil}
			return nil
		})

		then := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			response := args[0]
			fmt.Println(response)
			if response.Get("ok").Bool() {
				response.Call("json").Call("then", json)
			} else {
				js.Global().Call("alert", "Failed to upload receipt")
				ch <- Result{models.Receipt{}, errors.New("Upload failed")}
			}
			return nil
		})

		catch := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			ch <- Result{models.Receipt{}, errors.New("Network error")}
			return nil
		})

		promise.Call("then", then).Call("catch", catch)
	}()

	result := <-ch
	return result.receipt, result.err
}
