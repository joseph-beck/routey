package main

import (
	"net/http"

	routey "github.com/joseph-beck/routey/pkg/router"
)

func health() routey.HandlerFunc {
	return func(c *routey.Context) {
		c.Render(http.StatusOK, "health")
	}
}

func ping() routey.HandlerFunc {
	return func(c *routey.Context) {
		c.Render(http.StatusOK, "pong")
	}
}

func main() {
	c := routey.Config{
		Port:  ":3000",
		Debug: true,
		CORS:  false,
	}
	r := routey.New(c)
	r.Add(routey.Get, "/api", "/:id", health(), nil)
	r.Add(routey.Get, "/api/ping", "", ping(), nil)
	r.Run()
}
