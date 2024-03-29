package main

import (
	"fmt"

	routey "github.com/joseph-beck/routey/pkg/router"
	"github.com/joseph-beck/routey/pkg/status"
)

func health() routey.HandlerFunc {
	return func(c *routey.Context) {
		c.Render(status.OK, "health")
	}
}

func ping() routey.HandlerFunc {
	type Thing struct {
		Name  string `json:"name"`
		Other string `json:"other"`
	}

	f := func(c *routey.Context) {
		var t Thing
		c.ShouldBindJSON(&t)
		fmt.Println(t)
	}

	return func(c *routey.Context) {
		f(c)

		var t Thing
		c.ShouldBindJSON(&t)
		fmt.Println(t)
		p := c.Protocol()

		c.Render(status.OK, "pong "+p)
	}
}

func hello() routey.HandlerFunc {
	return func(c *routey.Context) {
		b, err := c.Body()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(b) + "\nbody")

		c.Status(status.OK)
	}
}

func middleware() routey.MiddlewareFunc {
	return func(c *routey.Context) {
		fmt.Println("middleware")
	}
}

func shutdown() routey.ShutdownFunc {
	return func() {
		fmt.Println("shutting down")
	}
}

func main() {
	c := routey.Config{
		Port:  ":3000",
		Debug: true,
		CORS:  false,
	}
	r := routey.New(c)
	r.Use(middleware())
	r.Add(routey.Get, "/api/health", "", health(), nil)
	r.Add(routey.Get, "/api/ping", "", ping(), nil)
	r.Add(routey.Get, "/api/hello", "", hello(), nil)
	go r.Shutdown(shutdown())
	r.Run()
}
