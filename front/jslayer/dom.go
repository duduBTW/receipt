//go:build js && wasm
// +build js,wasm

package jslayer

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"syscall/js"

	"github.com/a-h/templ"
)

const Window = "string"

func GetElementById(id string) (js.Value, error) {
	element := js.Global().Get("document").Call("getElementById", id)
	if IsNil(element) {
		return element, errors.New("Element not found")
	}

	return element, nil
}

func QuerySelector(selector string) (js.Value, error) {
	if selector == Window {
		return js.Global(), nil
	}

	element := js.Global().Get("document").Call("querySelector", selector)
	if IsNil(element) {
		return element, errors.New("Element not found")
	}

	return element, nil
}

func QuerySelectorAll(selector string) ([]js.Value, error) {
	return ElementQuerySelectorAll(js.Global().Get("document"), selector)
}
func ElementQuerySelectorAll(element js.Value, selector string) ([]js.Value, error) {
	if selector == Window {
		return []js.Value{js.Global()}, nil
	}

	nodeList := element.Call("querySelectorAll", selector)
	if IsNil(nodeList) {
		return nil, errors.New("Element not found")
	}

	var elements []js.Value
	length := nodeList.Get("length").Int()
	for i := 0; i < length; i++ {
		element := nodeList.Index(i)
		if IsNil(element) {
			continue
		}

		elements = append(elements, element)
	}

	return elements, nil
}

func Id(id string) string {
	return "#" + id
}

func GetStringAttr(selector string, attr string) (string, error) {
	element := js.Global().Get("document").Call("querySelector", selector)
	if IsNil(element) {
		return "", errors.New("Element not found")
	}
	return element.Get("dataset").Get(attr).String(), nil
}

func removeFirstAndLastChar(s string) string {
	if len(s) <= 1 {
		return s
	}
	return s[1 : len(s)-2]
}

func GetJsonData[T any](id string) (T, error) {
	var data T

	element := js.Global().Get("document").Call("getElementById", id)
	if IsNil(element) {
		return data, errors.New("Element not found")
	}
	jsonData := element.Get("innerText").String()
	jsonData = strings.Trim(strings.ReplaceAll(removeFirstAndLastChar(jsonData), `\"`, `"`), " ")
	fmt.Println(jsonData)

	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

func SetInnerText(selector string, text string) {
	element := js.Global().Get("document").Call("querySelector", selector)
	if IsNil(element) {
		fmt.Println("Element not found", selector)
		return
	}
	element.Set("innerText", text)
}

func CopyToClipboard(text string) {
	clipboard := js.Global().Get("navigator").Get("clipboard")
	if IsNil(clipboard) {
		fmt.Println("Clipboard API not supported")
		return
	}
	clipboard.Call("writeText", text)
}

func SetInnerHTML(selector string, html string) error {
	element := js.Global().Get("document").Call("querySelector", selector)
	if IsNil(element) {
		return errors.New("Element not found")
	}
	element.Set("innerHTML", html)
	return nil
}

func ReplaceWithHTML(selector string, html string) error {
	element := js.Global().Get("document").Call("querySelector", selector)
	if IsNil(element) {
		return errors.New("Element not found")
	}
	element.Call("replaceWith", js.Global().Get("document").Call("createRange").Call("createContextualFragment", html))
	return nil
}

func AppendHTMLInside(selector string, component templ.Component) error {
	element := js.Global().Get("document").Call("querySelector", selector)
	if IsNil(element) {
		return errors.New("Element not found")
	}
	html := HTMLFromComponent(component)
	fragment := js.Global().Get("document").Call("createRange").Call("createContextualFragment", html)
	element.Call("appendChild", fragment)
	return nil
}

func PrependHTMLInside(selector string, component templ.Component) error {
	element := js.Global().Get("document").Call("querySelector", selector)
	if IsNil(element) {
		return errors.New("Element not found")
	}
	html := HTMLFromComponent(component)
	fragment := js.Global().Get("document").Call("createRange").Call("createContextualFragment", html)
	firstChild := element.Get("firstChild")
	element.Call("insertBefore", fragment, firstChild)
	return nil
}

// https://developer.mozilla.org/pt-BR/docs/Web/API/Element/setAttribute
func SetAttr(element js.Value, name, value string) {
	element.Call("setAttribute", name, value)
}

// https://developer.mozilla.org/en-US/docs/Web/API/CSSStyleDeclaration/setProperty
func SetCssVar(element js.Value, name, value string) {
	element.Get("style").Call("setProperty", "--"+name, value)
}

func RegisterFunction(name string, function func() func()) {
	js.Global().Set(name, js.FuncOf(func(this js.Value, args []js.Value) any {
		go func() {
			function()
		}()
		return nil
	}))
}

func CreateIcons() {
	js.Global().Get("lucide").Call("createIcons")
}

func IsFocused(element js.Value) bool {
	activeElement := js.Global().Get("document").Get("activeElement")
	return activeElement.Equal(element)
}

func FocusElement(element js.Value) {
	element.Call("focus")
}

func Focus(selector string) error {
	element, err := QuerySelector(selector)
	if err != nil {
		return err
	}
	FocusElement(element)
	return nil
}

func Click(element js.Value) {
	element.Call("click")
}

func StopPropagation(element js.Value) {
	element.Call("stopPropagation")
}

type JSON struct {
	Value js.Value
}

func (json JSON) Stringify() string {
	return js.Global().Get("JSON").Call("stringify", json.Value).String()
}

func PreventDefault(event js.Value) {
	event.Call("preventDefault")
}
