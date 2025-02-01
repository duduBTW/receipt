//go:build js && wasm
// +build js,wasm

package jslayer

import (
	"context"
	"fmt"
	"strings"

	"github.com/a-h/templ"
)

type StateProps[T any] struct {
	Value           T
	Target          string
	RenderComponent func(value T) templ.Component
	OnMounted       func(value T)
	attributes      map[string]string
}

func NewState[T any](props StateProps[T]) StateProps[T] {
	fmt.Println(props.Value)
	props.attributes = make(map[string]string) // Initialize attributes map
	props.Set(props.Value)
	return props
}

func Render(component templ.Component, target string) error {
	componentHTML := HTMLFromComponent(component)
	return ReplaceWithHTML(target, componentHTML)
}

func (state *StateProps[T]) Set(value T) {
	state.Value = value
	err := Render(state.RenderComponent(value), state.Target)
	if err != nil {
		fmt.Println("Error setting inner HTML: ", err)
	} else {
		// Reapply stored attributes
		for key, val := range state.attributes {
			state.SetAttribute(key, val)
		}
	}

	if state.OnMounted != nil {
		state.OnMounted(value)
	}
}

// Function to set and store attributes
func (state *StateProps[T]) SetAttribute(attribute, value string) {
	if state.attributes == nil {
		state.attributes = make(map[string]string)
	}

	// Store the attribute in the map
	state.attributes[attribute] = value

	// Update the attribute on the rendered element
	if state.Target != "" {
		err := UpdateAttribute(state.Target, attribute, value)
		if err != nil {
			fmt.Println("Error updating attribute:", err)
		}
	}
}

// Helper function to extract HTML from a templ.Component
func HTMLFromComponent(component templ.Component) string {
	componentHTML := new(strings.Builder)
	component.Render(context.Background(), componentHTML)
	return componentHTML.String()
}

// Simulated function to update attributes of a DOM element
func UpdateAttribute(target, attribute, value string) error {
	element, err := QuerySelector(target)
	if err != nil {
		return err
	}

	SetAttr(element, attribute, value)
	return nil
}
