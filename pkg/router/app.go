package router

import "net/http"

type App struct {
	routes []Route
	port   string
}

func New() App {
	return App{
		port: ":8080",
	}
}

func (a *App) Route(r Route) {
	a.routes = append(a.routes, r)
}

func (a *App) Get(p string, h HandlerFunc) {
	a.routes = append(a.routes, Route{
		Path:        p,
		Method:      Get,
		HandlerFunc: h,
	})
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, e := range a.routes {
		m := e.Match(r)
		if !m {
			continue
		}

		c := &Context{
			w: w,
			r: r,
		}
		e.HandlerFunc.Serve(c)
		return
	}

	http.NotFound(w, r)
}

func (a *App) Run() {
	http.ListenAndServe(a.port, a)
}
