package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type HelloService struct{}

func (s HelloService) Add() []Route {
	return []Route{
		{
			Path:          "/hello",
			Params:        "",
			Method:        Get,
			HandlerFunc:   s.Get(),
			DecoratorFunc: nil,
		},
	}
}

func (s *HelloService) Get() HandlerFunc {
	return func(c *Context) {

	}
}

func TestAppNew(t *testing.T) {
	app := New()
	assert.NotNil(t, app)
	assert.Equal(t, app.port, ":8080")
	assert.True(t, app.debugMode)
}

func TestRoute(t *testing.T) {
	app := New()
	assert.Equal(t, 0, len(app.routes))
	route := Route{
		Path:          "/hello",
		Params:        "",
		Method:        Get,
		HandlerFunc:   func(c *Context) {},
		DecoratorFunc: nil,
	}
	app.Route(route)
	assert.Equal(t, 1, len(app.routes))
	assert.Equal(t, "/hello", app.routes[0].Path)
}

func TestService(t *testing.T) {
	app := New()
	assert.Equal(t, 0, len(app.routes))
	s := HelloService{}
	app.Service(s)
	assert.Equal(t, 1, len(app.routes))
	assert.Equal(t, "/hello", app.routes[0].Path)
}

func TestAdd(t *testing.T) {
	app := New()
	assert.Equal(t, 0, len(app.routes))
	app.Add(Get, "/hello", "/:name", func(c *Context) {}, nil)
	assert.Equal(t, 1, len(app.routes))
	assert.Equal(t, "/hello", app.routes[0].Path)
}

func TestGet(t *testing.T) {
	app := New()
	assert.Equal(t, 0, len(app.routes))
	app.Get("/hello", func(c *Context) {})
	assert.Equal(t, 1, len(app.routes))
	assert.Equal(t, "/hello", app.routes[0].Path)
}

func TestPost(t *testing.T) {
	app := New()
	assert.Equal(t, 0, len(app.routes))
	app.Post("/hello", func(c *Context) {})
	assert.Equal(t, 1, len(app.routes))
	assert.Equal(t, "/hello", app.routes[0].Path)
}

func TestPut(t *testing.T) {
	app := New()
	assert.Equal(t, 0, len(app.routes))
	app.Put("/hello", func(c *Context) {})
	assert.Equal(t, 1, len(app.routes))
	assert.Equal(t, "/hello", app.routes[0].Path)
}

func TestPatch(t *testing.T) {
	app := New()
	assert.Equal(t, 0, len(app.routes))
	app.Patch("/hello", func(c *Context) {})
	assert.Equal(t, 1, len(app.routes))
	assert.Equal(t, "/hello", app.routes[0].Path)
}

func TestDelete(t *testing.T) {
	app := New()
	assert.Equal(t, 0, len(app.routes))
	app.Delete("/hello", func(c *Context) {})
	assert.Equal(t, 1, len(app.routes))
	assert.Equal(t, "/hello", app.routes[0].Path)
}

func TestHead(t *testing.T) {
	app := New()
	assert.Equal(t, 0, len(app.routes))
	app.Head("/hello", func(c *Context) {})
	assert.Equal(t, 1, len(app.routes))
	assert.Equal(t, "/hello", app.routes[0].Path)
}

func TestOptions(t *testing.T) {
	app := New()
	assert.Equal(t, 0, len(app.routes))
	app.Options("/hello", func(c *Context) {})
	assert.Equal(t, 1, len(app.routes))
	assert.Equal(t, "/hello", app.routes[0].Path)
}
