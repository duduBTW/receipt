//go:build js && wasm
// +build js,wasm

package jslayer

import (
	"syscall/js"

	"github.com/a-h/templ"
)

type ElementSchema struct {
	Element  js.Value
	selector string
}

func (schema ElementSchema) TextContent() string {
	return schema.Element.Get("textContent").String()
}

func (schema ElementSchema) SetAttr(name, value string) {
	SetAttr(schema.Element, name, value)
}

// Removes the element from the DOM.
func (schema ElementSchema) Remove() {
	schema.Element.Call("remove")
}

func (schema ElementSchema) AppendHTMLInside(component templ.Component) error {
	AppendHTMLInside(schema.selector, component)
	return nil
}

func Element(selector string) (ElementSchema, error) {
	var schema ElementSchema
	element, err := QuerySelector(selector)
	if err != nil {
		return schema, err
	}

	schema.selector = selector
	schema.Element = element
	return schema, nil
}
