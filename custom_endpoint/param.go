package custom_endpoint

import (
	"net/http"
)

func NewHandleFuncParam(
	path string,
	f func(http.ResponseWriter, *http.Request),
	method ...string,
) *handleFuncParam {
	return &handleFuncParam{
		path:   path,
		f:      f,
		method: method,
	}
}

type handleFuncParam struct {
	path   string
	f      func(http.ResponseWriter, *http.Request)
	method []string
}
