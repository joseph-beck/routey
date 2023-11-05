package router

import (
	"log"
	"net/http"
)

type App struct {
	routes []Route
	port   string
}

func New() *App {
	return &App{
		port: ":8080",
	}
}

func (a *App) Add(method string, path string, handler HandlerFunc) {
	m := parseMethod(method)
	a.routes = append(a.routes, Route{
		Path:        path,
		Method:      m,
		HandlerFunc: handler,
	})
}

func (a *App) Route(route Route) {
	a.routes = append(a.routes, route)
}

func (a *App) Get(path string, handler HandlerFunc) {
	a.routes = append(a.routes, Route{
		Path:        path,
		Method:      Get,
		HandlerFunc: handler,
	})
}

func (a *App) Post(path string, handler HandlerFunc) {
	a.routes = append(a.routes, Route{
		Path:        path,
		Method:      Post,
		HandlerFunc: handler,
	})
}

func (a *App) Put(path string, handler HandlerFunc) {
	a.routes = append(a.routes, Route{
		Path:        path,
		Method:      Put,
		HandlerFunc: handler,
	})
}

func (a *App) Patch(path string, handler HandlerFunc) {
	a.routes = append(a.routes, Route{
		Path:        path,
		Method:      Patch,
		HandlerFunc: handler,
	})
}

func (a *App) Delete(path string, handler HandlerFunc) {
	a.routes = append(a.routes, Route{
		Path:        path,
		Method:      Delete,
		HandlerFunc: handler,
	})
}

func (a *App) Head(path string, handler HandlerFunc) {
	a.routes = append(a.routes, Route{
		Path:        path,
		Method:      Head,
		HandlerFunc: handler,
	})
}

func (a *App) Options(path string, handler HandlerFunc) {
	a.routes = append(a.routes, Route{
		Path:        path,
		Method:      Options,
		HandlerFunc: handler,
	})
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		r := recover()
		if r != nil {
			log.Println("error occurred: ", r)
			http.Error(w, "Oh Dear", http.StatusInternalServerError)
		}
	}()

	for _, e := range a.routes {
		c := &Context{
			writer:  w,
			request: r,
		}

		m := e.Match(c)
		if !m {
			continue
		}

		e.HandlerFunc.Serve(c)
		return
	}

	http.NotFound(w, r)
}

func (a *App) Run() {
	http.ListenAndServe(a.port, a)
}
