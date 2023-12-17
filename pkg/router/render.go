package router

import (
	"net/http"
)

var (
	jsonContentType = []string{"application/json; charset=utf-8"}
	htmlContentType = []string{"text/html; charset=utf-8"}
)

func writeContentType(w http.ResponseWriter, v []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = v
	}
}

type Renderer interface {
	Render(http.ResponseWriter) error
	WriteContentType(http.ResponseWriter)
}
