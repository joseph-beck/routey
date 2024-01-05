package router

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"strings"
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
//   - corsMode: if localhost, 127.0.0.1 or no origin do not allow this request.
//
//   - htmlDelims: HTML Delimiters, these can be customized.
//
//   - htmlRender: HTML Renderer, an interface that renders the HTML to the user.
//
//   - funcMap: templates.FuncMap
type App struct {
	routes     []Route
	middleware []MiddlewareFunc
	port       string

	logger    *logrus.Logger
	debugMode bool
	corsMode  bool

	htmlDelims HTMLDelims
	htmlRender HTMLRenderer
	funcMap    template.FuncMap
}

// Configure routey
//
//   - Port: what port do you want routey to use?
//
//   - Debug: do you want routey to run in debug mode?
//
//   - CORS: do you want to run this in local only mode?
type Config struct {
	Port  string
	Debug bool
	CORS  bool
}

// Create a new default App
func New(c ...Config) *App {
	a := App{
		port: ":8080",

		logger:    logrus.New(),
		debugMode: true,
		corsMode:  false,

		htmlDelims: HTMLDelims{Left: "{{", Right: "}}"},
		funcMap:    template.FuncMap{},
	}

	if c != nil || len(c) != 0 {
		a.port = c[0].Port
		a.debugMode = c[0].Debug
		a.corsMode = c[0].CORS
	}

	return &a
}

func (a *App) Use(m ...MiddlewareFunc) {
	if a.middleware == nil {
		a.middleware = make([]MiddlewareFunc, 0)
	}

	if len(m) <= 0 {
		return
	}

	a.middleware = append(a.middleware, m...)
}

// Adds a Route to the App
func (a *App) Route(r Route) {
	err := r.Format()
	if err != nil {
		logError(a.logger, err.Error()+" "+r.Path+r.Params, "ROUTE")
	}

	a.routes = append(a.routes, r)
}

// Adds a service to the App
func (a *App) Service(s Service) {
	r := s.Add()
	for _, rs := range r {
		a.Route(rs)
	}
}

// Add a Route with the method, path, params, handler and decorator
func (a *App) Add(method Method, path string, params string, handler HandlerFunc, decorator DecoratorFunc) {
	a.Route(Route{
		Path:          path,
		Params:        params,
		Method:        method,
		HandlerFunc:   handler,
		DecoratorFunc: decorator,
	})
}

// Add Get route
func (a *App) Get(path string, params string, handler HandlerFunc) {
	a.Route(Route{
		Path:        path,
		Params:      params,
		Method:      Get,
		HandlerFunc: handler,
	})
}

// Add Post route
func (a *App) Post(path string, params string, handler HandlerFunc) {
	a.Route(Route{
		Path:        path,
		Params:      params,
		Method:      Post,
		HandlerFunc: handler,
	})
}

// Add Put route
func (a *App) Put(path string, params string, handler HandlerFunc) {
	a.Route(Route{
		Path:        path,
		Params:      params,
		Method:      Put,
		HandlerFunc: handler,
	})
}

// Add Patch route
func (a *App) Patch(path string, params string, handler HandlerFunc) {
	a.Route(Route{
		Path:        path,
		Params:      params,
		Method:      Patch,
		HandlerFunc: handler,
	})
}

// Add Delete route
func (a *App) Delete(path string, params string, handler HandlerFunc) {
	a.Route(Route{
		Path:        path,
		Params:      params,
		Method:      Delete,
		HandlerFunc: handler,
	})
}

// Add Head route
func (a *App) Head(path string, params string, handler HandlerFunc) {
	a.Route(Route{
		Path:        path,
		Params:      params,
		Method:      Head,
		HandlerFunc: handler,
	})
}

// Add Options route
func (a *App) Options(path string, params string, handler HandlerFunc) {
	a.Route(Route{
		Path:        path,
		Params:      params,
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

// Set the current HTML renderer.
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

	if a.corsMode {
		o := r.Header.Get("Origin")
		if strings.HasPrefix(o, "http://localhost") || strings.HasPrefix(o, "http://127.0.0.1") || o == "" {
			logWarn(a.logger, fmt.Sprintf("Origin violating CORS, %s", o), "CORS")
			http.Error(w, "CORS restricted for localhost/local addresses when is CORS mode", http.StatusForbidden)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", o)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	}

	for _, e := range a.routes {
		c := &Context{
			app:   a,
			route: &e,

			writer:  w,
			request: r,
			state:   Healthy,
		}

		m := e.Match(c)
		if !m {
			continue
		}

		if len(a.middleware) > 0 {
			for _, f := range a.middleware {
				f(c)
			}
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

	for _, re := range a.routes {
		logRoute(a.logger, re)
	}

	a.logger.WithFields(logrus.Fields{
		"State": "Loading",
	}).Info("Loading app...")

	a.logger.WithFields(logrus.Fields{
		"State": "Routing",
	}).Info(fmt.Sprintf("Serving %d routes, on port %s", len(a.routes), a.port))

	if a.debugMode {
		logWarn(a.logger, "Currently using Debug Mode", "DEBUG")
	}
	if !a.corsMode {
		logWarn(a.logger, "Currently not using CORS Mode", "CORS")
	}

	http.ListenAndServe(a.port, a)
}

// Shutdown the App, should be ran as go Shutdown()
func (a *App) Shutdown(f ...ShutdownFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	fmt.Print("\r")

	a.logger.WithFields(logrus.Fields{
		"State": "Closing",
	}).Info("Closing app...")

	// Closing down stuff

	if len(f) > 0 {
		for _, e := range f {
			e()
		}
	}

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

// Log a route that is being used
func logRoute(l *logrus.Logger, e Route) {
	l.WithFields(logrus.Fields{
		"Route": e.Method.String(),
	}).Info("added route " + e.Path + e.Params)
}

// Log error
func logError(l *logrus.Logger, m string, v string) {
	l.WithFields(logrus.Fields{
		"Error": v,
	}).Error(m)
}

// Log a warning
func logWarn(l *logrus.Logger, m string, v string) {
	l.WithFields(logrus.Fields{
		"Warn": v,
	}).Warn(m)
}
