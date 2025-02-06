//go:build js && wasm
// +build js,wasm

package handlers

import (
	"fmt"
	"syscall/js"

	"github.com/a-h/templ"
	"github.com/dudubtw/receipt/constants"
	"github.com/dudubtw/receipt/front/jslayer"
	"github.com/dudubtw/receipt/renderer/components"
	"github.com/dudubtw/receipt/renderer/pages"
)

type TabState struct {
	Id          string
	DefaultData pages.TabsContentProps
	state       jslayer.StateProps[pages.TabsContentProps]
	events      []jslayer.EventListener
}

func (state *TabState) New(getBody func(pages.TabsContentProps) templ.Component) {
	isRerenderFromKeyboardAction := false
	tabItemsSelector := jslayer.Id(state.Id) + " " + jslayer.Id(constants.IdTabItem)

	var set = func(activeIndex, targetIndex int) {
		if activeIndex == targetIndex {
			return
		}

		isRerenderFromKeyboardAction = true
		// Removes old index
		state.state.Value.Items[activeIndex].IsActive = false
		// Sets new index
		state.state.Value.Items[targetIndex].IsActive = true

		// Rerenders page
		state.state.Set(state.state.Value)
	}

	tabClickEvent := jslayer.EventListener{
		Selector:  tabItemsSelector,
		EventType: "click",
		Listener: func(this js.Value, args []js.Value) {
			clickedTabValue := this.Get("dataset").Get(components.TabItemValueAttr).String()
			targetIndex := state.IndexFromValue(clickedTabValue)
			activeIndex := state.ActiveIndex()
			set(activeIndex, targetIndex)
		},
	}

	tabKeyDownEvent := jslayer.EventListener{
		Selector:  tabItemsSelector,
		EventType: "keydown",
		Listener: func(this js.Value, args []js.Value) {
			event := args[0]
			intiialDirection := 0
			direction := intiialDirection

			switch event.Get("key").String() {
			case "ArrowLeft":
				direction = -1
			case "ArrowRight":
				direction = 1
			}

			if direction == intiialDirection {
				return
			}

			activeIndex := state.ActiveIndex()
			targetIndex := Clamp(activeIndex+direction, 0, len(state.state.Value.Items)-1)
			set(activeIndex, targetIndex)
		},
	}

	state.events = append(state.events, tabKeyDownEvent, tabClickEvent)

	var focus = func() {
		if !isRerenderFromKeyboardAction {
			return
		}

		isRerenderFromKeyboardAction = false
		elements, err := jslayer.QuerySelectorAll(tabItemsSelector)
		if err != nil {
			fmt.Println("Could not focus tab!")
			return
		}

		jslayer.FocusElement(elements[state.ActiveIndex()])
	}

	state.state = jslayer.StateProps[pages.TabsContentProps]{
		Value:  state.DefaultData,
		Target: jslayer.Id(state.Id),
		RenderComponent: func(value pages.TabsContentProps) templ.Component {
			return pages.Tabs(pages.TabsProps{
				Items: value.Items,
				Body:  getBody(value),
			})
		},
		OnMounted: func(value pages.TabsContentProps) {
			jslayer.AddEvents(state.events)
			focus()
			jslayer.CreateIcons()
		},
	}

	jslayer.AddEvents(state.events)
}

func (state TabState) ActiveIndex() int {
	for index, tab := range state.DefaultData.Items {
		if tab.IsActive {
			return index
		}
	}

	return -1
}

func (state TabState) IndexFromValue(targetValue string) int {
	for index, tab := range state.DefaultData.Items {
		if tab.Value == targetValue {
			return index
		}
	}

	return -1
}

func Clamp(value, min, max int) int {
	if value > max {
		return max
	}

	if value < min {
		return min
	}

	return value
}

func GetActiveTab(tabs []components.TabItemProps) string {
	for _, tab := range tabs {
		if tab.IsActive {
			return tab.Value
		}
	}

	return ""
}
