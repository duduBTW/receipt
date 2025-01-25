//go:build js && wasm
// +build js,wasm

package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"syscall/js"
	"time"

	"github.com/dudubtw/receipt/constants"
	jslayer "github.com/dudubtw/receipt/front/jslayer"
	"github.com/dudubtw/receipt/front/service"
	"github.com/dudubtw/receipt/models"
	"github.com/dudubtw/receipt/renderer/components"
)

func setDefaultDate() {
	dateInput, err := jslayer.QuerySelector(jslayer.Id(constants.IdAddCategoryDateInput))
	if err != nil {
		return
	}

	dateInput.Set("value", time.Now().Format("2006-01-02"))
}

func AppendAddCategoryComponent() {
	err := jslayer.PrependHTMLInside(jslayer.Id(constants.IdGloabal), components.AddCategory())
	if err != nil {
		fmt.Println(err)
	}
	jslayer.CreateIcons()
}

func CreateModalSetup() func() {
	var modal = ReceiptModal(func(receipt models.NewReceipt, file js.Value) (models.Receipt, error) {
		return service.UploadRecepit(receipt, file)
	})

	var fileChangeHandler = jslayer.EventListener{
		Selector:  jslayer.Id(constants.IdAddCategoryFileInput),
		EventType: "change",
		Listener: func(this js.Value, args []js.Value) {
			fileInput := this
			file := fileInput.Get("files").Index(0)

			if file.IsNull() || !strings.HasPrefix(file.Get("type").String(), "image/") {
				return
			}

			// Create a file  url
			fileUrl := js.Global().Get("URL").Call("createObjectURL", file).String()
			modal.Open(models.Receipt{
				ImageName: fileUrl,
				CategoryID: func() int64 {
					idStr := jslayer.GetQueryParam(constants.ReceiptSearchParamCategory)
					id, err := strconv.ParseInt(idStr, 10, 64)
					if err != nil {
						return 0
					}
					return id
				}(),
			})
		},
	}

	var addCategoryButtonClickEvent = jslayer.EventListener{
		Selector:  jslayer.Id(constants.IdAddCategoryButton),
		EventType: "click",
		Listener: func(this js.Value, args []js.Value) {
			fileInput, err := jslayer.QuerySelector(jslayer.Id(constants.IdAddCategoryFileInput))
			if err != nil {
				return
			}

			fileInput.Call("click")
		},
	}

	AppendAddCategoryComponent()

	var eventList = []jslayer.EventListener{addCategoryButtonClickEvent, fileChangeHandler}
	jslayer.AddEvents(eventList)
	return func() {
		jslayer.RemoveEvents(eventList)
	}
}
