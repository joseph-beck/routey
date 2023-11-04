package main

import (
	"net/http"

	routey "github.com/joseph-beck/routey/pkg/router"
)

func main() {
	r := routey.New()
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

	r.Run()
}
