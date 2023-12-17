package router

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

// App struct
//
//   - routes: array containg all Route structs
//
//   - port: string port
//
//   - logger: structured logging
//
//   - debugMode: if debugMode is enabled things such as HTML as less static, will be slower but easier to debug.
//
//   - htmlDelims: HTML Delimiters, these can be customized.
//
//   - htmlRender: HTML Renderer, an interface that renders the HTML to the user.
//
//   - funcMap: templates.FuncMap
type App struct {
	routes []Route
	port   string

	logger    *logrus.Logger
	debugMode bool

	htmlDelims HTMLDelims
	htmlRender HTMLRenderer
	funcMap    template.FuncMap
}

// Create a new default App
func New() *App {
	return &App{
		port: ":8080",

		logger:    logrus.New(),
		debugMode: true,

		htmlDelims: HTMLDelims{Left: "{{", Right: "}}"},
		funcMap:    template.FuncMap{},
	}
}

// Add a Route with the method, path, params, handler and decorator
func (a *App) Add(method Method, path string, params string, handler HandlerFunc, decorator DecoratorFunc) {
	p, err := parseParams(params)
	if err != nil {
		a.logger.WithFields(logrus.Fields{
			"Error": "Route",
		}).Error("Bad Params on route", path, "params", params)
	}

	a.routes = append(a.routes, Route{
		Path:          path,
		Params:        p,
		Method:        method,
		HandlerFunc:   handler,
		DecoratorFunc: decorator,
	})
}

// Add Route
func (a *App) Route(route Route) {
	a.routes = append(a.routes, route)
}

// Add Get route
func (a *App) Get(path string, handler HandlerFunc) {
	a.Route(Route{
		Path:        path,
		Method:      Get,
		HandlerFunc: handler,
	})
}

// Add Post route
func (a *App) Post(path string, handler HandlerFunc) {
	a.Route(Route{
		Path:        path,
		Method:      Post,
		HandlerFunc: handler,
	})
}

// Add Put route
func (a *App) Put(path string, handler HandlerFunc) {
	a.Route(Route{
		Path:        path,
		Method:      Put,
		HandlerFunc: handler,
	})
}

// Add Patch route
func (a *App) Patch(path string, handler HandlerFunc) {
	a.Route(Route{
		Path:        path,
		Method:      Patch,
		HandlerFunc: handler,
	})
}

// Add Delete route
func (a *App) Delete(path string, handler HandlerFunc) {
	a.Route(Route{
		Path:        path,
		Method:      Delete,
		HandlerFunc: handler,
	})
}

// Add Head route
func (a *App) Head(path string, handler HandlerFunc) {
	a.Route(Route{
		Path:        path,
		Method:      Head,
		HandlerFunc: handler,
	})
}

// Add Options route
func (a *App) Options(path string, handler HandlerFunc) {
	a.Route(Route{
		Path:        path,
		Method:      Options,
		HandlerFunc: handler,
	})
}

// Load a folder of HTML files.
func (a *App) LoadHTMLGlob(p string) {
	t := template.Must(template.New("").Delims(a.htmlDelims.Left, a.htmlDelims.Right).Funcs(a.funcMap).ParseGlob(p))

	if a.debugMode {
		a.htmlRender = HTMLDebug{
			Glob:    p,
			FuncMap: a.funcMap,
			Delims:  a.htmlDelims,
		}
		return
	}

	a.SetHTMLTemplate(t)
}

// Load a series of HTML files.
func (a *App) LoadHTMLFiles(f ...string) {
	if a.debugMode {
		a.htmlRender = HTMLDebug{
			Files:   f,
			FuncMap: a.funcMap,
			Delims:  a.htmlDelims,
		}
		return
	}

	t := template.Must(template.New("").Delims(a.htmlDelims.Left, a.htmlDelims.Right).Funcs(a.funcMap).ParseFiles(f...))
	a.SetHTMLTemplate(t)
}

func (a *App) SetHTMLTemplate(t *template.Template) {
	a.htmlRender = HTMLRender{Template: t.Funcs(a.funcMap)}
}

// ServerHTTP with ResponseWriter and Request
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		r := recover()
		if r != nil {
			a.logger.WithFields(logrus.Fields{
				"Error": "Panic",
			}).Error(r)
			http.Error(w, "Oh Dear", http.StatusInternalServerError)
		}
	}()

	for _, e := range a.routes {
		c := &Context{
			app: a,

			writer:  w,
			request: r,
		}

		m := e.Match(c)
		if !m {
			continue
		}

		if e.DecoratorFunc == nil {
			e.HandlerFunc.Serve(c)
			logRequest(a.logger, e)
			return
		}
		e.DecoratorFunc(e.HandlerFunc).Serve(c)
		logRequest(a.logger, e)
		return
	}

	http.NotFound(w, r)
}

// Run the App
func (a *App) Run() {
	fmt.Println(`

	 _____   ____  _    _ _______ ________     __
	|  __ \ / __ \| |  | |__   __|  ____\ \   / /
	| |__) | |  | | |  | |  | |  | |__   \ \_/ / 
	|  _  /| |  | | |  | |  | |  |  __|   \   /  
	| | \ \| |__| | |__| |  | |  | |____   | |   
	|_|  \_\\____/ \____/   |_|  |______|  |_|   
												 
												 
	`)

	a.logger.WithFields(logrus.Fields{
		"State": "Loading",
	}).Info("Loading app...")
	a.logger.WithFields(logrus.Fields{
		"State": "Routing",
	}).Info(fmt.Sprintf("Serving %d routes", len(a.routes)))
	http.ListenAndServe(a.port, a)
}

// Shutdown the App, should be ran as go Shutdown()
func (a *App) Shutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	fmt.Print("\r")

	a.logger.WithFields(logrus.Fields{
		"State": "Closing",
	}).Info("Closing app...")

	// Closing down stuffs

	a.logger.WithFields(logrus.Fields{
		"State": "Exit",
	}).Info("Closed app")

	os.Exit(0)
}

// Log a request with info
func logRequest(l *logrus.Logger, e Route) {
	l.WithFields(logrus.Fields{
		"Request": e.Method.String(),
	}).Info(e.Path + e.Params)
}
