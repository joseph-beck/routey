package router

import (
	"html/template"
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

type HTMLRenderer interface {
	Instance(string, any) Renderer
}

type HTMLDelims struct {
	Left  string
	Right string
}

type HTML struct {
	Template *template.Template
	Name     string
	Data     any
}

type HTMLRender struct {
	Template *template.Template
	Delims   HTMLDelims
}

type HTMLDebug struct {
	Files   []string
	Glob    string
	Delims  HTMLDelims
	FuncMap template.FuncMap
}

func (h HTML) Render(w http.ResponseWriter) error {
	h.WriteContentType(w)

	if h.Name == "" {
		return h.Template.Execute(w, h.Data)
	}
	return h.Template.ExecuteTemplate(w, h.Name, h.Data)
}

func (h HTML) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, htmlContentType)
}

func (h HTMLRender) Instance(n string, d any) Renderer {
	return HTML{
		Template: h.Template,
		Name:     n,
		Data:     d,
	}
}

func (h HTMLDebug) Instance(n string, d any) Renderer {
	var t *template.Template
	l := false

	if h.FuncMap == nil {
		h.FuncMap = template.FuncMap{}
	}
	if len(h.Files) > 0 {
		t = template.Must(template.New("").Delims(h.Delims.Left, h.Delims.Right).Funcs(h.FuncMap).ParseFiles(h.Files...))
		l = true
	}
	if h.Glob != "" {
		t = template.Must(template.New("").Delims(h.Delims.Left, h.Delims.Right).Funcs(h.FuncMap).ParseGlob(h.Glob))
		l = true
	}
	if !l {
		panic("tried to create an empty HTML render")
	}

	return HTML{
		Template: t,
		Name:     n,
		Data:     d,
	}
}
