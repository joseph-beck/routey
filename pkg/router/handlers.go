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

// Middleware function, run before every request
type MiddlewareFunc func(c *Context)

// Serves the middleware function
func (m MiddlewareFunc) Serve(c *Context) {
	m(c)
}

// Func for shutting down the router
type ShutdownFunc func()
