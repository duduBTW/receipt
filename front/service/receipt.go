//go:build js && wasm
// +build js,wasm

package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"syscall/js"

	"github.com/dudubtw/receipt/constants"
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
		promise := js.Global().Call("fetch", constants.ApiUpload, js.ValueOf(opts))

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

type FetchReceiptChan struct {
	receipt models.Receipt
	err     error
}

func FetchReceipt(id string) (models.Receipt, error) {
	ch := make(chan FetchReceiptChan)

	go func() {
		response, err := http.Get(constants.ApiRecepipt + id)
		if err != nil {
			ch <- FetchReceiptChan{models.Receipt{}, err}
		}
		defer response.Body.Close()

		var recepit models.Receipt
		if err := json.NewDecoder(response.Body).Decode(&recepit); err != nil {
			ch <- FetchReceiptChan{models.Receipt{}, err}
			return
		}
		ch <- FetchReceiptChan{recepit, nil}
	}()

	result := <-ch
	return result.receipt, result.err
}

func UpdateRecepit(receipt models.Receipt) (models.Receipt, error) {
	ch := make(chan FetchReceiptChan)
	go func() {
		body, err := json.Marshal(receipt)
		if err != nil {
			ch <- FetchReceiptChan{models.Receipt{}, err}
			return
		}

		fmt.Println("body", string(body))

		response, err := http.Post(constants.ApiUpdateReceipt, "application/json", bytes.NewBuffer(body))
		if err != nil {
			ch <- FetchReceiptChan{models.Receipt{}, err}
			return
		}
		defer response.Body.Close()

		var recepit models.Receipt
		if err := json.NewDecoder(response.Body).Decode(&recepit); err != nil {
			ch <- FetchReceiptChan{models.Receipt{}, err}
			return
		}
		ch <- FetchReceiptChan{recepit, nil}
	}()

	result := <-ch
	return result.receipt, result.err
}
