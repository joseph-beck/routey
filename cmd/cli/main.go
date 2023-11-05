package main

import (
	"fmt"
	"net/http"

	routey "github.com/joseph-beck/routey/pkg/router"
)

func main() {
	r := routey.New()
	// Test route that returns some http.
	r.Route(routey.Route{
		Path:   "/",
		Params: "",
		Method: routey.Get,
		HandlerFunc: func(c *routey.Context) {
			b := "<h1>Healthy</h1>"
			c.Render(http.StatusOK, b)
		},
		DecoratorFunc: nil,
	})
	// Test error route
	r.Route(routey.Route{
		Path:   "/400",
		Params: "",
		Method: routey.Get,
		HandlerFunc: func(c *routey.Context) {
			c.Status(http.StatusBadRequest)
		},
		DecoratorFunc: nil,
	})
	// Test route that panics
	r.Route(routey.Route{
		Path:   "/panic",
		Params: "",
		Method: routey.Get,
		HandlerFunc: func(c *routey.Context) {
			panic("this is a panic")
		},
		DecoratorFunc: nil,
	})
	// Test route that has parameters
	r.Route(routey.Route{
		Path:   "/echo",
		Params: `/(?P<string>\w+)`,
		Method: routey.Get,
		HandlerFunc: func(c *routey.Context) {
			p, err := c.Param("string")
			if err != nil {
				c.Status(http.StatusBadRequest)
				return
			}
			c.Render(http.StatusOK, p)
		},
		DecoratorFunc: nil,
	})
	// Test param int
	r.Route(routey.Route{
		Path:   "/add",
		Params: `/(?P<one>\w+)/(?P<two>\w+)`,
		Method: routey.Get,
		HandlerFunc: func(c *routey.Context) {
			o, err := c.ParamInt("one")
			if err != nil {
				c.Status(http.StatusBadRequest)
				return
			}
			t, err := c.ParamInt("two")
			if err != nil {
				c.Status(http.StatusBadRequest)
				return
			}
			a := o + t

			c.Render(http.StatusOK, fmt.Sprintf("%d", a))
		},
		DecoratorFunc: nil,
	})

	r.Run()
}
