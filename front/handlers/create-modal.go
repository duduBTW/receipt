//go:build js && wasm
// +build js,wasm

package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"syscall/js"
	"time"

	"github.com/a-h/templ"
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
	err := jslayer.AppendHTMLInside(jslayer.Id(constants.IdRoot), components.AddCategory())
	if err != nil {
		fmt.Println(err)
	}
	jslayer.CreateIcons()
}

func modalFocusTrapFactory() jslayer.EventListener {
	var modalId = jslayer.Id(constants.IdAddCategoryModal)
	focusableElementSelector := modalId + "a[href]:not([tabindex='-1']), button:not([tabindex='-1']), input:not([tabindex='-1']), textarea:not([tabindex='-1']), select:not([tabindex='-1']), details:not([tabindex='-1']), [tabindex]:not([tabindex='-1'])"

	return jslayer.EventListener{
		Selector:  modalId,
		EventType: "keydown",
		Listener: func(this js.Value, args []js.Value) {
			event := args[0]
			key := event.Get("key").String()
			if key != "Tab" {
				return
			}

			modalElement, err := jslayer.QuerySelector(modalId)
			if err != nil {
				return
			}

			focusableElements, err := jslayer.ElementQuerySelectorAll(modalElement, focusableElementSelector)
			if err != nil {
				return
			}

			var focusableElementsLength = len(focusableElements)
			var firstFocusableElement = focusableElements[0]
			var lastFocusableElement = focusableElements[focusableElementsLength-1]

			if event.Get("shiftKey").Bool() && jslayer.IsFocused(firstFocusableElement) {
				event.Call("preventDefault")
				lastFocusableElement.Call("focus")
				return
			}

			if jslayer.IsFocused(lastFocusableElement) {
				event.Call("preventDefault")
				firstFocusableElement.Call("focus")
			}
		},
	}
}

func CreateModalSetup() func() {
	var modalState jslayer.StateProps[components.AddCategoryModalProps]
	var modalFocusTrapEvent = modalFocusTrapFactory()
	var selectedCategory = ""

	var categorySelect = CategorySelect{
		Id:           constants.IdAddCategoryCategorySelect,
		DefaultValue: selectedCategory,
		OnValueChange: func(s string) {

			fmt.Println("Selected category:", s)
			selectedCategory = s
		},
	}

	var closeModal = func() {
		modalState.Set(components.AddCategoryModalProps{
			IsOpen: false,
		})
	}

	var openModal = func(fileUrl string) {
		modalState.Set(components.AddCategoryModalProps{
			IsOpen:          true,
			ReceiptImageUrl: fileUrl,
		})
	}

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
			openModal(fileUrl)
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

	var closeAddCategoryButtonClickEvent = jslayer.EventListener{
		Selector:  jslayer.Id(constants.IdAddCategoryCloseButton),
		EventType: "click",
		Listener: func(this js.Value, args []js.Value) {
			closeModal()
		},
	}

	var closeModalKeyboardEvent = jslayer.EventListener{
		Selector:  jslayer.Window,
		EventType: "keydown",
		Listener: func(this js.Value, args []js.Value) {
			event := args[0]
			key := event.Get("key").String()

			if key == "Escape" {
				closeModal()
			}
		},
	}

	var formSubmitHandler = jslayer.EventListener{
		Selector:  jslayer.Id(constants.IdAddCategoryForm),
		EventType: "submit",
		Listener: func(this js.Value, args []js.Value) {
			event := args[0]
			event.Call("preventDefault")

			categoryId, err := strconv.ParseInt(selectedCategory, 10, 64)
			fmt.Println(err)
			if err != nil {
				return
			}

			dateInput, err := jslayer.QuerySelector(jslayer.Id(constants.IdAddCategoryDateInput))
			if err != nil {
				return
			}

			fileInput, err := jslayer.QuerySelector(jslayer.Id(constants.IdAddCategoryFileInput))
			if err != nil {
				return
			}
			file := fileInput.Get("files").Index(0)
			if jslayer.IsNil(file) {
				return
			}

			// Get the date value
			dateValue := dateInput.Get("value").String()
			date, err := time.Parse("2006-01-02", dateValue)
			if err != nil {
				fmt.Println("Invalid date")
				return
			}

			_, err = service.UploadRecepit(models.NewReceipt{
				CategoryID: categoryId,
				Date:       date.Format("2006-01-02"),
			}, file)

			if err != nil {
				fmt.Println("Error uploading file:", err)
				return
			}

			jslayer.Redirect(constants.ReceiptRoute, constants.ReceiptSearchParamCategory, selectedCategory)
		},
	}

	modalState = jslayer.StateProps[components.AddCategoryModalProps]{
		Value:  components.AddCategoryModalDefaultProps,
		Target: jslayer.Id(constants.IdAddCategoryModal),
		RenderComponent: func(props components.AddCategoryModalProps) templ.Component {
			return components.AddCategoryModal(props)
		},
		OnMounted: func(value components.AddCategoryModalProps) {
			if value.IsOpen {
				go func() {
					categorySelect.New()
				}()

				closeModalKeyboardEvent.Add()
				closeAddCategoryButtonClickEvent.Add()
				modalFocusTrapEvent.Add()
				formSubmitHandler.Add()
				jslayer.CreateIcons()
				setDefaultDate()
			} else {
				formSubmitHandler.Remove()
				modalFocusTrapEvent.Remove()
				closeModalKeyboardEvent.Remove()
				closeAddCategoryButtonClickEvent.Remove()
				categorySelect.Remove()
			}
		},
	}

	AppendAddCategoryComponent()

	var eventList = []jslayer.EventListener{addCategoryButtonClickEvent, fileChangeHandler}
	jslayer.AddEvents(eventList)
	return func() {
		jslayer.RemoveEvents(eventList)
	}
}
