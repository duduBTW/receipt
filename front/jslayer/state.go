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
}

func NewState[T any](props StateProps[T]) StateProps[T] {
	fmt.Println(props.Value)
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
		fmt.Println("Error setting inner html: ", err)
	}

	if state.OnMounted != nil {
		state.OnMounted(value)
	}
}

func HTMLFromComponent(component templ.Component) string {
	componentHTML := new(strings.Builder)
	component.Render(context.Background(), componentHTML)
	return componentHTML.String()
}
