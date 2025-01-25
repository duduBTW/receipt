//go:build js && wasm
// +build js,wasm

package handlers

import (
	"fmt"
	"strconv"
	"syscall/js"
	"time"

	"github.com/a-h/templ"
	"github.com/dudubtw/receipt/constants"
	"github.com/dudubtw/receipt/front/jslayer"
	"github.com/dudubtw/receipt/models"
	"github.com/dudubtw/receipt/renderer/components"
)

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

type ReceiptModalActions struct {
	Open  func(recepit models.Receipt)
	Close func()
}

func ReceiptModal(onSubmit func(models.NewReceipt, js.Value) (models.Receipt, error)) ReceiptModalActions {
	var modalState jslayer.StateProps[components.AddCategoryModalProps]
	var modalFocusTrapEvent = modalFocusTrapFactory()
	var selectedCategory = ""

	var categorySelect = CategorySelect{
		Id:           constants.IdAddCategoryCategorySelect,
		DefaultValue: selectedCategory,
		OnValueChange: func(s string) {
			selectedCategory = s
		},
	}

	var closeModal = func() {
		modalState.Set(components.AddCategoryModalProps{
			IsOpen: false,
		})
	}

	var openModal = func(recepit models.Receipt) {
		modalState.Set(components.AddCategoryModalProps{
			IsOpen:  true,
			Recepit: recepit,
		})
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
			fmt.Println("selectedCategory", selectedCategory)
			event := args[0]
			event.Call("preventDefault")

			categoryId, err := strconv.ParseInt(selectedCategory, 10, 64)
			if err != nil {
				fmt.Println(err)
				return
			}

			dateInput, err := jslayer.QuerySelector(jslayer.Id(constants.IdAddCategoryDateInput))
			if err != nil {
				fmt.Println("1")
				fmt.Println(err)
				return
			}

			fileInput, err := jslayer.QuerySelector(jslayer.Id(constants.IdAddCategoryFileInput))
			if err != nil {
				fmt.Println("2")
				fmt.Println(err)
				return
			}

			// Get the date value
			dateValue := dateInput.Get("value").String()
			date, err := time.Parse("2006-01-02", dateValue)
			if err != nil {
				fmt.Println("Invalid date")
				return
			}

			_, err = onSubmit(models.NewReceipt{
				CategoryID: categoryId,
				Date:       date.Format("2006-01-02"),
			}, fileInput.Get("files").Index(0))

			if err != nil {
				fmt.Println("Failed to submit modal", err)
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
					categorySelect.Set(strconv.FormatInt(value.Recepit.CategoryID, 10))
				}()

				dateInput, err := jslayer.QuerySelector(jslayer.Id(constants.IdAddCategoryDateInput))
				if err == nil {
					dateInput.Set("value", value.Recepit.Date)
				}

				closeModalKeyboardEvent.Add()
				closeAddCategoryButtonClickEvent.Add()
				modalFocusTrapEvent.Add()
				formSubmitHandler.Add()
				jslayer.CreateIcons()
				setDefaultDate()
				jslayer.Focus(jslayer.Id(constants.IdAddCategoryCloseButton))
			} else {
				formSubmitHandler.Remove()
				modalFocusTrapEvent.Remove()
				closeModalKeyboardEvent.Remove()
				closeAddCategoryButtonClickEvent.Remove()
				categorySelect.Remove()
			}
		},
	}

	return ReceiptModalActions{
		Open:  openModal,
		Close: closeModal,
	}
}
