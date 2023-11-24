package router

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	logr "github.com/sirupsen/logrus"
)

// App struct
//
//   - routes: array containg all Route structs
//
//   - port: string port
//
//   - logger: structured logging
type App struct {
	routes []Route
	port   string

	logger *logr.Logger
}

// Create a new default App
func New() *App {
	return &App{
		port: ":8080",

		logger: logr.New(),
	}
}

// Add a Route with the method, path, params, handler and decorator
func (a *App) Add(method Method, path string, params string, handler HandlerFunc, decorator DecoratorFunc) {
	p, err := parseParams(params)
	if err != nil {
		a.logger.WithFields(logr.Fields{
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

// ServerHTTP with ResponseWriter and Request
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		r := recover()
		if r != nil {
			a.logger.WithFields(logr.Fields{
				"Error": "Panic",
			}).Error(r)
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

	a.logger.WithFields(logr.Fields{
		"State": "Loading",
	}).Info("Loading app...")
	a.logger.WithFields(logr.Fields{
		"State": "Routing",
	}).Info(fmt.Sprintf("Serving %d routes", len(a.routes)))
	http.ListenAndServe(a.port, a)
}

// Shutdown the App, should be ran like go Shutdown()
func (a *App) Shutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	fmt.Print("\r")

	a.logger.WithFields(logr.Fields{
		"State": "Closing",
	}).Info("Closing app...")

	// Closing down stuffs

	a.logger.WithFields(logr.Fields{
		"State": "Exit",
	}).Info("Closed app")

	os.Exit(0)
}

// Log a request with info
func logRequest(l *logr.Logger, e Route) {
	l.WithFields(logr.Fields{
		"Request": e.Method.String(),
	}).Info(e.Path + e.Params)
}
