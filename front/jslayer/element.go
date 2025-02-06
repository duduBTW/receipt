//go:build js && wasm
// +build js,wasm

package jslayer

import "syscall/js"

type ElementSchema struct {
	Element js.Value
}

func (schema ElementSchema) TextContent() string {
	return schema.Element.Get("textContent").String()
}

func Element(selector string) (ElementSchema, error) {
	var schema ElementSchema
	element, err := QuerySelector(selector)
	if err != nil {
		return schema, err
	}

	schema.Element = element
	return schema, nil
}
