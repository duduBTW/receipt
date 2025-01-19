//go:build js && wasm
// +build js,wasm

package jslayer

import (
	"net/url"
	"syscall/js"
)

func getCurrentUrl() string {
	return js.Global().Get("location").Get("href").String()
}

func navigate(to string) {
	js.Global().Get("location").Call("replace", to)
}

func GetQueryParam(queryParamName string) string {
	currentUrl, err := url.Parse(getCurrentUrl())
	if err != nil {
		return ""
	}

	return currentUrl.Query().Get(queryParamName)
}

func SetQueryParam(queryParamName, queryParamValue string) {
	currentUrl, err := url.Parse(getCurrentUrl())
	if err != nil {
		return
	}

	query := currentUrl.Query()
	query.Set(queryParamName, queryParamValue)
	currentUrl.RawQuery = query.Encode()
	navigate(currentUrl.String())
}

func Redirect(target, queryParamName, queryParamValue string) {
	currentUrl, err := url.Parse(getCurrentUrl())
	if err != nil {
		return
	}

	currentUrl.Path = target
	query := currentUrl.Query()
	query.Set(queryParamName, queryParamValue)
	currentUrl.RawQuery = query.Encode()
	navigate(currentUrl.String())
}
