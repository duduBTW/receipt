//go:build js && wasm
// +build js,wasm

package handlers

import (
	"syscall/js"

	"github.com/a-h/templ"
	"github.com/dudubtw/receipt/constants"
	"github.com/dudubtw/receipt/front/jslayer"
	"github.com/dudubtw/receipt/models"
	"github.com/dudubtw/receipt/renderer/pages"
)

func HomeSetup() func() {
	categories, err := jslayer.JsonData[[]models.Category](constants.IdHomeData)
	if err != nil {
		Global.Error(err)
	}

	tabsState := TabState{
		Id: constants.IdHomeTabs,
		DefaultData: pages.TabsContentProps{
			Items: pages.HomeTabs,
		},
	}

	categoryCardClickHandler := jslayer.EventListener{
		Selector:  jslayer.Id(constants.IdCategoryCard),
		EventType: "click",
		Listener: func(this js.Value, args []js.Value) {
			elesments, err := jslayer.ElementQuerySelectorAll(this, "a")
			if err != nil || len(elesments) == 0 {
				return
			}

			jslayer.Click(elesments[0])
		},
	}

	tabsState.New(func(props pages.TabsContentProps) templ.Component {
		activeTab := GetActiveTab(props.Items)

		switch activeTab {
		case "categories":
			return pages.CategoriesComponent(pages.CategoriesComponentProps{
				Categories: categories,
			})
		case "date":
			return pages.TestDate()
		}

		return nil
	})

	categoryCardClickHandler.Add()
	return func() {
		categoryCardClickHandler.Remove()
	}
}
