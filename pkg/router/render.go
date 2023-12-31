package router

import (
	"net/http"
)

var (
	jsonContentType  = []string{"application/json; charset=utf-8"}
	tomlContentType  = []string{"application/toml; charset=utf-8"}
	yamlContentType  = []string{"application/x-yaml; charset=utf-8"}
	xmlContentType   = []string{"application/xml; charset=utf-8"}
	htmlContentType  = []string{"text/html; charset=utf-8"}
	plainContentType = []string{"text/plain; charset=utf-8"}
)

// Writes the content type to the response header
func writeContentType(w http.ResponseWriter, v []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = v
	}
}

// An interface for Rendering in a response
type Renderer interface {
	Render(http.ResponseWriter) error
	WriteContentType(http.ResponseWriter)
}
