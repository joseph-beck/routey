package router

import "net/http"

// HandlerFunc takes a routey.Context pointer
type HandlerFunc func(c *Context)

// Serves the HandlerFunc
func (f HandlerFunc) Serve(c *Context) {
	f(c)
}

// Decorate a HandlerFunc
type DecoratorFunc func(f HandlerFunc) HandlerFunc

// Wrap a standard http Handler in a routey HandlerFunc
func Wrap(f http.HandlerFunc) HandlerFunc {
	return func(c *Context) {
		f(c.writer, c.request)
	}
}
