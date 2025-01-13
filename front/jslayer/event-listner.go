//go:build js && wasm
// +build js,wasm

package jslayer

import (
	"syscall/js"
)

type EventListener struct {
	Selector    string
	EventType   string
	Listener    func(this js.Value, args []js.Value)
	__listeners []RegisterdEventListener
}

type RegisterdEventListener struct {
	EventListener
	jsListener js.Func
}

func (event *EventListener) Add() {
	event.Remove()

	registerdEventListener := RegisterdEventListener{
		EventListener: *event,
		jsListener: js.FuncOf(func(this js.Value, args []js.Value) any {
			event.Listener(this, args)
			return nil
		}),
	}

	__listen("addEventListener", registerdEventListener)
	event.__listeners = append(event.__listeners, registerdEventListener)
}

func (event *EventListener) Remove() {
	for _, listener := range event.__listeners {
		__listen("removeEventListener", listener)
	}

	event.__listeners = nil
}

func __listen(function string, event RegisterdEventListener) {
	elements, err := QuerySelectorAll(event.Selector)
	if err != nil {
		return
	}

	for _, element := range elements {
		element.Call(function, event.EventType, event.jsListener)
	}
}

func AddEvents(events []EventListener) {
	for _, event := range events {
		event.Add()
	}
}

func RemoveEvents(events []EventListener) {
	for _, event := range events {
		event.Remove()
	}
}
