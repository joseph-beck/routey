package router

import (
	"html/template"
	"net/http"
)

// Interface for Rendering HTML
type HTMLRenderer interface {
	Instance(string, any) Renderer
}

// Delimiters used in the HTML file, defaults to {{ ... }}
type HTMLDelims struct {
	Left  string
	Right string
}

// Default HTML renderer.
type HTML struct {
	Template *template.Template
	Name     string
	Data     any
}

// HTML Renderer implementor
type HTMLRender struct {
	Template *template.Template
	Delims   HTMLDelims
}

// Debug HTML renderer, allows for editing of HTML files at runtime.
type HTMLDebug struct {
	Files   []string
	Glob    string
	Delims  HTMLDelims
	FuncMap template.FuncMap
}

// Render the HTML
func (h HTML) Render(w http.ResponseWriter) error {
	h.WriteContentType(w)

	if h.Name == "" {
		return h.Template.Execute(w, h.Data)
	}
	return h.Template.ExecuteTemplate(w, h.Name, h.Data)
}

// Write the content type
func (h HTML) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, htmlContentType)
}

// Create an instance of the HTMLRenderer
func (h HTMLRender) Instance(n string, d any) Renderer {
	return HTML{
		Template: h.Template,
		Name:     n,
		Data:     d,
	}
}

// Create an instance of the HTMLRenderer, for the debugging Renderer
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
